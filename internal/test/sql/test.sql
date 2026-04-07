CREATE TABLE IF NOT EXISTS `document` (
    `id`           bigint       NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `created_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `source`       varchar(255) NOT NULL DEFAULT '' COMMENT '文档来源',
    `doc_id`       varchar(255) NOT NULL DEFAULT '' COMMENT '文档唯一标识',
    `title`        varchar(255) NOT NULL DEFAULT '' COMMENT '文档标题',
    `file_path`    varchar(512) NOT NULL DEFAULT '' COMMENT '文件存储路径',
    `content_hash` varchar(255) NOT NULL DEFAULT '' COMMENT '内容哈希值，用于校验文档是否变更',
    `version`      bigint       NOT NULL DEFAULT 0 COMMENT '文档版本号',
    `status`       tinyint      NOT NULL DEFAULT 0 COMMENT '文档状态',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_docid` (`doc_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文档表';
