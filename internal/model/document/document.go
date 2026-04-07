package document

import (
	"time"

	"gorm.io/gorm"

	"github.com/demoManito/pulse/internal/service"
)

// Document 文档
type Document struct {
	ID        int64     `gorm:"cloumn:id"`
	CreatedAt time.Time `gorm:"cloumn:created_at"`
	UpdatedAt time.Time `gorm:"cloumn:updated_at"`

	Source      string `gorm:"cloumn:source"` // 文档来源，如 "wx", "local" 等
	DocID       string `gorm:"cloumn:doc_id"`
	Title       string `gorm:"cloumn:title"`
	FilePath    string `gorm:"cloumn:file_path"`
	ContentHash string `gorm:"cloumn:content_hash"`
	Version     int64  `gorm:"cloumn:version"`
	Status      int    `gorm:"cloumn:status"`
}

func (Document) TableName() string {
	return "document"
}

func DB() *gorm.DB {
	return service.DB
}

func (d *Document) GetByDocID(docID string) (*Document, error) {
	var doc Document
	err := DB().Where("doc_id = ?", docID).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (d *Document) Upsert(doc *Document) error {
	return DB().Where("doc_id = ?", doc.DocID).
		Assign(Document{
			Title:       doc.Title,
			FilePath:    doc.FilePath,
			ContentHash: doc.ContentHash,
			Version:     doc.Version,
			Status:      doc.Status,
		}).
		FirstOrCreate(doc).Error
}
