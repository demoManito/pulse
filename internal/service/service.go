package service

import (
	"github.com/demoManito/pulse/pkg/database"
	"github.com/demoManito/pulse/pkg/gitops"
	"github.com/demoManito/pulse/pkg/wecom"
	"gorm.io/gorm"

	"github.com/demoManito/pulse/config"
)

var (
	DB    *gorm.DB
	WeCom *wecom.Client
	Git   *gitops.GitOps
)

func Init(cfg *config.Config) (err error) {
	DB, err = database.NewClient(cfg.Database)
	if err != nil {
		return err
	}
	WeCom, err = wecom.NewClient(cfg.WeCom)
	if err != nil {
		return err
	}
	Git, err = gitops.New(cfg.GitOps)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}
