# Pulse

企业微信文档 → Git 仓库实时同步工具。

当企业微信文档发生新建、更新或删除操作时，Pulse 通过回调事件实时感知变更，将文档内容转换为 Markdown 格式并自动同步到指定的 Git 仓库。

## 功能特性

- 实时监听企业微信文档变更（基于回调事件）
- 文档内容自动转换为 Markdown 格式存储
- 支持按知识库 ID 过滤，只同步指定知识库的文档
- 通过 GORM 管理文档元数据（文档 ID、文件路径、内容哈希等映射关系）
- 自动 Git commit & push，保留完整变更历史
- 支持 HTTP / SSH 两种 Git 认证方式
- 自动从 Git URL 提取仓库名，支持 `~` 路径展开

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
│   │   └── doc/                     # 文档同步业务逻辑（事件→拉取→转换→提交）；HTML→Markdown 转换
│   ├── model/
│   │   └── document/                # 文档模型
│   ├── server/
│   │   └── httpsvc/                 # HTTP 服务
│   ├── service/                     # 服务层（全局 DB、WeCom、Git 实例）
│   └── test/                        # 测试
├── pkg/
│   ├── database/                    # 数据库连接封装
│   ├── handler/                     # 通用接口封装
│   ├── logger/                      # 日志接口（支持 zap / logrus）
│   ├── wecom/                       # 企业微信 SDK（客户端、加解密、类型定义）
│   └── gitops/                      # Git 操作（clone/pull/commit/push）
```

## 依赖

| 依赖 | 用途 |
|---|---|
| `github.com/gin-gonic/gin` | HTTP 框架 |
| `gopkg.in/yaml.v3` | 配置文件解析 |
| `gorm.io/gorm` + `gorm.io/driver/mysql` | ORM + MySQL 驱动 |
| `github.com/go-git/go-git/v5` | Git 操作 |
| `github.com/JohannesKaufmann/html-to-markdown/v2` | HTML 转 Markdown |
| `go.uber.org/zap` / `github.com/sirupsen/logrus` | 日志 |

## 快速开始

### 1. 创建企业微信应用

1. 登录 [企业微信管理后台](https://work.weixin.qq.com/wework_admin/frame)
2. 进入「应用管理」→「自建」→ 创建应用
3. 记录以下信息：
   - **CorpID**：在「我的企业」→「企业信息」页面获取
   - **Secret**：在应用详情页获取
4. 确保应用拥有「文档」相关 API 权限

### 2. 准备数据库

```sql
CREATE DATABASE pulse CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 配置

编辑 `config/config.yaml`：

```yaml
http:
  port: 9784
  address: "0.0.0.0"

database:
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "your_password"
  database: "pulse"
  charset: "utf8mb4"
  parse_time: true
  loc: "Local"
  max_open_conns: 100
  max_idle_conns: 50
  conn_max_lifetime: 300
  conn_max_idle_time: 60

wecom:
  corp_id: "your_corp_id"
  corp_secret: "your_app_secret"
  callback_token: "your_callback_token"
  callback_encoding_aes_key: "your_43char_base64_aes_key"
  agent_id: 1000001
  sync:
    wiki_id: ""  # 要同步的知识库 ID，为空则同步所有文档

gitops:
  url: "https://git.example.com/your-org/docs.git"
  local_path: "~/DocRepos"   # 支持 ~ 展开，自动追加仓库名
  branch: "main"
  auth:
    type: "http"             # "ssh" 或 "http"
    token: "your_token"      # HTTP Token（type: "http"）
    ssh_key: ""              # SSH 私钥路径（type: "ssh"）
    password: ""             # SSH 私钥密码（type: "ssh"）
```

配置说明：
- `local_path` 支持 `~` 自动展开为用户主目录，并自动从 URL 提取仓库名拼接。例如 `~/DocRepos` + `https://git.example.com/org/prd.git` → `~/DocRepos/prd`
- `wiki_id` 填写知识库 ID（从知识库 URL 中获取），只有属于该知识库的文档才会触发同步。留空则同步所有收到的文档事件

### 4. 运行

```bash
go run main.go
```

### 5. 配置企业微信回调

1. 在企业微信管理后台，进入应用设置 →「接收消息」
2. 设置回调 URL 为 `https://your-domain:9784/callback`
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
  过滤：检查文档是否属于配置的知识库（wiki_id）
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
| `source` | 文档来源（`wx`） |
| `doc_id` | 企业微信文档唯一 ID |
| `title` | 文档标题 |
| `file_path` | Git 仓库中的相对路径（`<doc_id>/<title>.md`） |
| `content_hash` | 文档内容 SHA256，用于判断是否真正变更 |
| `version` | 文档版本号（取自 modify_time） |
| `status` | `1`=active / `2`=deleted |

## License

MIT
