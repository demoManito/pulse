package doc

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/demoManito/pulse/internal/model/document"
	"github.com/demoManito/pulse/internal/service"
	"github.com/demoManito/pulse/pkg/logger"
	"github.com/demoManito/pulse/pkg/wecom"
)

const (
	StatusActive  = 1
	StatusDeleted = 2
)

// HandleEvent 处理企业微信文档回调事件。
// eventData 是解密后的 XML 消息体。
func HandleEvent(ctx context.Context, eventData []byte) error {
	var event wecom.DocEvent
	if err := xml.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("doc sync: 解析事件 XML 失败: %w", err)
	}

	if event.DocID == "" {
		logger.Infof("doc sync: 忽略非文档事件, MsgType=%s Event=%s", event.MsgType, event.Event)
		return nil
	}

	logger.Infof("doc sync: 收到文档事件, DocID=%s Event=%s ChangeType=%s", event.DocID, event.Event, event.ChangeType)

	switch event.ChangeType {
	case "delete_doc":
		return handleDelete(ctx, event.DocID)
	default:
		// create_doc / update_doc 等都走同步流程
		return handleSync(ctx, event.DocID)
	}
}

// handleSync 处理文档新建/更新：拉取内容 → 转 Markdown → 写入 Git → 更新数据库。
func handleSync(ctx context.Context, docID string) error {
	// 1. 获取文档信息
	docInfo, err := service.WeCom.GetDocInfo(docID)
	if err != nil {
		return fmt.Errorf("doc sync: 获取文档信息失败: %w", err)
	}

	// 2. 过滤：如果配置了 wiki_id，只同步属于该知识库的文档
	wikiID := service.WeCom.WikiID()
	if wikiID != "" {
		if !strings.HasPrefix(docInfo.WikiFileID, wikiID+"_") {
			logger.Infof("doc sync: 文档不属于目标知识库, DocID=%s WikiFileID=%s", docID, docInfo.WikiFileID)
			return nil
		}
	}

	// 3. 获取文档内容
	contentResp, err := service.WeCom.GetDocContent(docID)
	if err != nil {
		return fmt.Errorf("doc sync: 获取文档内容失败: %w", err)
	}

	// 3. 转换为 Markdown
	markdown, err := HTMLToMarkdown(contentResp.Data.Doc.DocContent)
	if err != nil {
		return fmt.Errorf("doc sync: HTML 转 Markdown 失败: %w", err)
	}

	// 4. 计算内容哈希，判断是否有变更
	hash := ContentHash(markdown)
	doc := &document.Document{}
	existing, err := doc.GetByDocID(docID)
	if err == nil && existing.ContentHash == hash {
		logger.Infof("doc sync: 文档内容未变更, DocID=%s", docID)
		return nil
	}

	// 5. 确定文件路径：以 doc_id 为目录，文档标题为文件名
	//    结构：<doc_id>/<title>.md
	fileName := sanitizeFileName(docInfo.DocName) + ".md"
	relPath := filepath.Join(docID, fileName)

	// 如果数据库中已有记录，检查标题是否变更（需要重命名旧文件）
	if existing != nil && existing.FilePath != "" && existing.FilePath != relPath {
		oldAbsPath := filepath.Join(service.Git.LocalPath(), existing.FilePath)
		if err := os.Remove(oldAbsPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("doc sync: 删除旧文件失败: %w", err)
		}
		// stage 旧文件的删除
		if err := service.Git.Add(existing.FilePath); err != nil {
			return fmt.Errorf("doc sync: git add 旧文件失败: %w", err)
		}
	}

	// 6. 写入本地 Git 仓库
	absPath := filepath.Join(service.Git.LocalPath(), relPath)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return fmt.Errorf("doc sync: 创建目录失败: %w", err)
	}
	if err := os.WriteFile(absPath, []byte(markdown), 0o644); err != nil {
		return fmt.Errorf("doc sync: 写入文件失败: %w", err)
	}

	// 7. Git add → commit → push
	if err := service.Git.Add(relPath); err != nil {
		return fmt.Errorf("doc sync: git add 失败: %w", err)
	}
	commitMsg := fmt.Sprintf("sync: 更新文档「%s」", docInfo.DocName)
	if err := service.Git.Commit(commitMsg); err != nil {
		return fmt.Errorf("doc sync: git commit 失败: %w", err)
	}
	if err := service.Git.Push(ctx); err != nil {
		return fmt.Errorf("doc sync: git push 失败: %w", err)
	}

	// 8. 更新数据库
	record := &document.Document{
		Source:      "wx",
		DocID:       docID,
		Title:       docInfo.DocName,
		FilePath:    relPath,
		ContentHash: hash,
		Version:     docInfo.ModifyTime,
		Status:      StatusActive,
	}
	if err := doc.Upsert(record); err != nil {
		return fmt.Errorf("doc sync: 更新数据库失败: %w", err)
	}

	logger.Infof("doc sync: 文档同步完成, DocID=%s FilePath=%s", docID, relPath)
	return nil
}

// handleDelete 处理文档删除：从 Git 仓库删除文件 → 更新数据库。
func handleDelete(ctx context.Context, docID string) error {
	doc := &document.Document{}
	existing, err := doc.GetByDocID(docID)
	if err != nil {
		logger.Warnf("doc sync: 删除事件但数据库中未找到文档, DocID=%s: %v", docID, err)
		return nil
	}

	// 1. 删除本地文件
	absPath := filepath.Join(service.Git.LocalPath(), existing.FilePath)
	if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("doc sync: 删除文件失败: %w", err)
	}

	// 2. Git add → commit → push
	if err := service.Git.Add(existing.FilePath); err != nil {
		return fmt.Errorf("doc sync: git add 失败: %w", err)
	}
	commitMsg := fmt.Sprintf("sync: 删除文档「%s」", existing.Title)
	if err := service.Git.Commit(commitMsg); err != nil {
		return fmt.Errorf("doc sync: git commit 失败: %w", err)
	}
	if err := service.Git.Push(ctx); err != nil {
		return fmt.Errorf("doc sync: git push 失败: %w", err)
	}

	// 3. 标记数据库记录为已删除
	existing.Status = StatusDeleted
	if err := doc.Upsert(existing); err != nil {
		return fmt.Errorf("doc sync: 更新数据库失败: %w", err)
	}

	logger.Infof("doc sync: 文档删除完成, DocID=%s FilePath=%s", docID, existing.FilePath)
	return nil
}

// sanitizeFileName 将文档标题转换为安全的文件名。
func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_", "|", "_",
	)
	name = replacer.Replace(name)
	if name == "" {
		name = "untitled"
	}
	return name
}
