# Pulse

企业微信文档 → Git 仓库实时同步工具。

当企业微信文档发生新建、更新或删除操作时，Pulse 通过回调事件实时感知变更，将文档内容转换为 Markdown 格式并自动同步到指定的 Git 仓库。

## 功能特性

- 实时监听企业微信文档变更（基于回调事件）
- 文档内容自动转换为 Markdown 格式存储
- 通过 GORM 管理文档元数据（文档 ID、文件路径、内容哈希等映射关系）
- 自动 Git commit & push，保留完整变更历史
- 支持配置目标 Git 仓库地址和分支

## 项目结构

```
pulse/
├── main.go                          # 入口
├── cmd/
│   └── http/                        # HTTP 服务启动
├── config/                          # 配置加载
├── internal/
│   ├── app/
│   │   └── controller/              # 控制器（封装回调接口等）
│   │       └── api/                 # API 路由
│   ├── biz/
│   │   └── doc/                     # 文档同步业务逻辑（事件 → 拉取 → 转换 → 提交）；HTML → Markdown 转换
│   ├── model/
│   │   └── document/                # 文档模型
│   ├── server/
│   │   └── httpsvc/                 # HTTP 服务
│   ├── service/                     # 服务层
│   └── test/                        # 测试
├── pkg/
│   ├── database/                    # 数据库连接和操作封装
│   ├── handler/                     # 通用接口封装
│   ├── logger/                      # 日志接口
│   ├── wecom/                       # 企业微信 SDK（客户端、加解密）
│   └── gitops/                      # Git 操作业务逻辑（clone/pull/commit/push）
```

## 依赖

| 依赖 | 用途 |
|---|---|
| `gopkg.in/yaml.v3` | 配置文件解析 |
| `gorm.io/gorm` | ORM |
| `gorm.io/driver/mysql` | MySQL 驱动 |
| `github.com/go-git/go-git/v5` | Git 操作 |
| `github.com/JohannesKaufmann/html-to-markdown/v2` | HTML 转 Markdown |

## 快速开始

### 1. 创建企业微信应用

1. 登录 [企业微信管理后台](https://work.weixin.qq.com/wework_admin/frame)
2. 进入「应用管理」→「自建」→ 创建应用
3. 记录以下信息：
   - **CorpID**：在「我的企业」→「企业信息」页面获取
   - **Secret**：在应用详情页获取
4. 确保应用拥有「文档」相关 API 权限

### 2. 准备数据库

创建 MySQL 数据库：

```sql
CREATE DATABASE pulse CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 配置

复制配置文件并填写：

```bash
cp config/config.yaml.example config.yaml
```

```yaml
server:
  port: 8080

wecom:
  corp_id: "your_corp_id"
  corp_secret: "your_app_secret"
  callback_token: "your_callback_token"
  callback_encoding_aes_key: "your_43char_base64_aes_key"

database:
  driver: "mysql"
  dsn: "user:pass@tcp(127.0.0.1:3306)/pulse?charset=utf8mb4&parseTime=True"

git:
  repo_url: "git@github.com:your-org/your-docs.git"
  local_path: "./docs-repo"
  branch: "main"
  commit_author: "Pulse Bot"
  commit_email: "pulse@example.com"
```

### 4. 运行

```bash
go run main.go
```

### 5. 配置企业微信回调

1. 在企业微信管理后台，进入应用设置 →「接收消息」
2. 设置回调 URL 为 `https://your-domain:8080/callback`
3. 填入配置文件中的 Token 和 EncodingAESKey
4. 保存时企业微信会发送验证请求，验证通过即配置成功

> 注意：回调 URL 需要公网可访问。开发阶段可使用 ngrok 等工具暴露本地服务。

## 工作流程

```
企业微信文档变更
       ↓
  回调事件推送（POST /callback）
       ↓
  解密事件 → 解析事件类型
       ↓
  ┌─ 新建/更新 ──────────────────┐
  │  调用 API 获取文档内容         │
  │  转换为 Markdown              │
  │  对比 ContentHash 判断是否变更  │
  │  写入 Git 仓库 → commit → push │
  │  更新数据库元数据              │
  └──────────────────────────────┘
  ┌─ 删除 ──────────────────────┐
  │  查询数据库获取文件路径        │
  │  从 Git 仓库删除文件          │
  │  commit → push               │
  │  标记数据库记录为 deleted      │
  └──────────────────────────────┘
```

## 数据模型

| 字段 | 说明 |
|---|---|
| `doc_id` | 企业微信文档唯一 ID |
| `title` | 文档标题 |
| `file_path` | Git 仓库中的相对路径 |
| `content_hash` | 文档内容 SHA256，用于判断是否真正变更 |
| `version` | 文档版本号 |
| `status` | `active` / `deleted` |

## License

MIT
