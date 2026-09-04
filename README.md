# CATendar

CATendar 是一个可常驻桌面、支持窗口置顶的本地月历。应用使用 Wails 将 Vue 3 前端、Go 业务逻辑和 SQLite 数据库打包进同一个桌面进程，不依赖云服务或本地 HTTP 服务。

## 功能

- 月视图、上月/下月导航和一键返回今天
- 创建、编辑、确认删除事件
- 全天事件或带起止时间的事件
- 事件颜色、标题强调和可点击的描述链接
- 单日侧栏，以及可切换左右位置的侧栏
- `100%–150%` 日历缩放；缩放后可滚动或右键拖动画布
- 浅色/深色主题、事件时间显隐，设置会保存在本机
- 窗口置顶、关闭后驻留系统托盘
- 键盘操作和清晰的焦点状态
- AI Sync：自动调用本机 Codex 或 Claude Code，从最近 30 天邮件中整理并批量更新日历

## 快捷键

| 操作 | 快捷键 |
|---|---|
| 上一个月 / 下一个月 | `←` / `→` |
| 返回今天 | `T` |
| 新建事件 | `N` |
| 缩放日历 | `Ctrl`/`Cmd` + 滚轮 |

日期格支持键盘聚焦。点击日期查看当天安排；悬停日期后点击 `+` 可直接新建事件。

## 技术栈

- Wails 2.12
- Go 1.24
- Vue 3、TypeScript、Pinia、Naive UI、Vite
- SQLite（`go-sqlite3`）

仓库只保留一套 Wails 主线。旧 Electron 前端、Echo REST 服务及其重复依赖已经移除。

## 本地开发

先安装 Go、Node.js、平台所需的 C 编译工具和 Wails CLI，然后执行：

```bash
cd frontend
npm ci
cd ..
wails dev
```

只预览前端布局时可以运行：

```bash
cd frontend
npm run dev
```

纯浏览器预览没有 Wails IPC，因此可以浏览界面和打开表单，但不会写入事件。

## AI Sync

AI Sync 由三层组成，Codex 和 Claude Code 只是在最外层使用不同的启动适配器：

1. `core/email` 在 Go Core 内连接邮箱，读取最近 30 天邮件；邮箱密码或应用专用密码不会交给 Agent。
2. `core/ai_sync/cli` 将本次同步允许使用的邮件读取、日历查询和批量写入能力封装成统一的 CATendar CLI。
3. `core/ai_sync/agent` 自动启动用户选择的 Codex 或 Claude Code，并给两者同一份任务约束和同一套 CLI。

在工具栏打开 AI Sync 设置后：

- 配置一个 IMAP 邮箱账户并测试连接；首版支持通用 IMAP，内置 QQ、Gmail、Outlook 和 iCloud 参数预设。
- 选择 Codex 或 Claude Code。CATendar 会自动探测本机命令，也可以填写可执行文件的绝对路径。
- 点击工具栏的 `AI Sync`。应用会在后台自动调用 Agent，无需用户另外打开终端或复制提示词。

QQ 邮箱使用 `imap.qq.com:993` 和 TLS。需要先在 QQ 邮箱网页设置中开启 IMAP/SMTP 服务并生成授权码，然后把授权码填入 CATendar 的 `Password / app authorization code`，不要填写 QQ 登录密码。

Agent 只会获得一次性、本次运行有效的 CATendar CLI 会话。它不能读取邮箱凭据，也不能任意访问 CATendar 数据库；日历写入由 Go Core 校验、去重并执行。邮件生成的事件会在描述中包含 `catendar://email/...` 来源链接，点击后由应用重新读取并展示原邮件。

每次同步最多列出最近 30 天的 200 封邮件。列表只暴露邮件 ID、Message-ID、主题、发件人、邮箱服务器收件时间 `receivedAt`、发件头时间 `sentAt` 和内部来源链接；Agent 再按候选批量读取正文，每次最多 20 封，可以继续分批读取。正文会清理为纯文本并限制为每封 256 KiB，附件不会交给 Agent。单封邮件解析失败会以 `readError` 返回，不再中断整批同步。

两种 Agent 共用一份高召回筛选和时间推导规则：面试、笔试、会议、截止日期等强信号邮件必须读取正文；“收到后”等相对时间以 `receivedAt` 为锚点，明确写“发送后”时才使用 `sentAt`。最终仍只写入有邮件证据且能确定或计算日期的事件。

CATendar CLI 是 Agent 使用的内部能力界面；它仅在 AI Sync 启动的受控会话中可用：

```text
CATendar cli email list --cursor 0 --limit 50
CATendar cli email read --ids <email-id>[,<email-id>]
CATendar cli calendar list --from 2026-09-01 --to 2026-10-01
CATendar cli calendar batch-upsert --input events.json
```

## 构建

```bash
wails build
```

Wails 会先执行前端类型检查与 Vite 构建，再将 `frontend/dist` 嵌入可执行文件。产物位于 `build/bin/`，不会提交到 Git。

单独验证两端时：

```bash
cd frontend && npm run build
cd .. && go build ./...
```

## 本地数据与凭据

事件数据库保存在系统用户配置目录下的 `CATendar/calendar.db`：

- Windows：`%AppData%\CATendar\calendar.db`
- macOS：`~/Library/Application Support/CATendar/calendar.db`
- Linux：`$XDG_CONFIG_HOME/CATendar/calendar.db`，未设置时通常为 `~/.config/CATendar/calendar.db`

SQLite 只保存日历、邮箱账户的非敏感配置和 AI Provider 设置。邮箱密码或应用专用密码使用系统钥匙串保存：macOS Keychain、Windows Credential Manager，或 Linux Secret Service。它们不会写入 SQLite、日志或 Git。

## 目录结构

```text
.
├── main.go                       # 进程入口、Wails 启动和 CLI 模式分流
├── app.go                        # 日历能力的 Wails 适配层
├── ai_sync_app.go                # AI Sync 的 Wails 适配层
├── tray*.go                      # 系统托盘平台实现
├── core/                         # 不依赖 UI 的 Go Core
│   ├── core.go                   # 组装数据库、日历、邮箱和 AI Sync 服务
│   ├── database.go               # SQLite 初始化和迁移
│   ├── calendar/                 # 日历 CRUD、批量写入、去重和冲突检查
│   ├── email/                    # 邮箱配置、系统钥匙串和 IMAP 读取
│   └── ai_sync/
│       ├── ai_sync.go            # 同步任务状态与能力编排
│       ├── session.go            # 一次性授权的本地 Core/CLI 会话
│       ├── cli/                  # 统一 CATendar CLI 表现层
│       └── agent/                # Agent Runner、Prompt、Codex/Claude 适配器
├── frontend/                     # Vue 3 界面及 Wails 生成绑定
├── build/                        # 打包图标和平台元数据
└── wails.json                    # Wails 构建配置
```
