package wecom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/demoManito/pulse/pkg/logger"
)

const baseURL = "https://qyapi.weixin.qq.com"

// Client 企业微信 API 客户端，自动管理 access_token。
type Client struct {
	corpID     string
	corpSecret string

	mu          sync.RWMutex
	accessToken string
	expiresAt   time.Time

	httpClient *http.Client
}

// NewClient 创建企业微信 API 客户端。
func NewClient(cfg Config) (*Client, error) {
	return &Client{
		corpID:     cfg.CorpID,
		corpSecret: cfg.CorpSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// GetAccessToken 获取有效的 access_token，过期时自动刷新。
func (c *Client) GetAccessToken() (string, error) {
	c.mu.RLock()
	if c.accessToken != "" && time.Now().Before(c.expiresAt) {
		token := c.accessToken
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// 获取写锁后二次检查
	if c.accessToken != "" && time.Now().Before(c.expiresAt) {
		return c.accessToken, nil
	}

	url := fmt.Sprintf("%s/cgi-bin/gettoken?corpid=%s&corpsecret=%s", baseURL, c.corpID, c.corpSecret)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("wecom: get token request failed: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("wecom: decode token response failed: %w", err)
	}
	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("wecom: get token error: %d %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	c.accessToken = tokenResp.AccessToken
	// 提前 5 分钟刷新
	c.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn-300) * time.Second)
	logger.Infof("wecom: access token refreshed, expires at %s", c.expiresAt.Format(time.RFC3339))

	return c.accessToken, nil
}

// GetDocContent 根据文档 ID 获取文档内容。
func (c *Client) GetDocContent(docID string) (*DocContentResponse, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/cgi-bin/wedoc/document/get?access_token=%s", baseURL, token)
	body := fmt.Sprintf(`{"docid":"%s"}`, docID)
	resp, err := c.httpClient.Post(url, "application/json", io.NopCloser(
		strings.NewReader(body),
	))
	if err != nil {
		return nil, fmt.Errorf("wecom: get doc content request failed: %w", err)
	}
	defer resp.Body.Close()

	var result DocContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("wecom: decode doc content response failed: %w", err)
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wecom: get doc content error: %d %s", result.ErrCode, result.ErrMsg)
	}
	return &result, nil
}

// GetDocInfo 获取文档元信息。
func (c *Client) GetDocInfo(docID string) (*DocInfo, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/cgi-bin/wedoc/doc_info?access_token=%s", baseURL, token)
	body := fmt.Sprintf(`{"docid":"%s"}`, docID)
	resp, err := c.httpClient.Post(url, "application/json", io.NopCloser(
		strings.NewReader(body),
	))
	if err != nil {
		return nil, fmt.Errorf("wecom: get doc info request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode int     `json:"errcode"`
		ErrMsg  string  `json:"errmsg"`
		Data    DocInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("wecom: decode doc info response failed: %w", err)
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wecom: get doc info error: %d %s", result.ErrCode, result.ErrMsg)
	}
	return &result.Data, nil
}
