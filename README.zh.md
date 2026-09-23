# CodexKeeper

[English](README.md)

基于 [CPA Usage Keeper](https://github.com/Willxup/cpa-usage-keeper) 精简的独立统计服务，面向 CodexProxy，只采集原生 **Codex 和 Claude**。保留原版页面布局、图表、浅色／深色主题和响应式界面。

## 保留与移除

保留总览、实时监控、分析、请求明细、导出、Token／缓存／延迟／错误／费用统计、账号名称与别名、Codex／Claude 额度和订阅信息，以及 Codex 额度历史。管理员密码登录、会话管理、SQLite 存储和备份继续保留。

移除社区与本地排行、排行上报和资料、API Key 用户登录及独立视图、其他供应商的凭证与额度支持、旧队列名称探测、HTTP 采集降级。通过其他供应商转发的 Claude／GPT 模型不会计入原生 Codex／Claude。

直接使用当前 CodexProxy 的管理接口和 RESP 数据流，不需要再向 CodexProxy 集成代码。价格目录仍作为通用数据源使用，供应商专属价格选择规则仅保留 OpenAI／Anthropic。

## Docker 内网部署

每个 CPA 实例只运行一个 Keeper 采集器。CodexProxy 需要开启 `usage-statistics-enabled: true`，允许 Docker 网络访问管理接口，并配置管理密钥。这里使用的是管理密钥，不是调用模型的客户端 API Key。

```sh
cp .env.example .env
# 在 .env 中填写 CPA_MANAGEMENT_KEY 和私有 LOGIN_PASSWORD。
# 如果 Docker 网络名称不同，设置 CPA_DOCKER_NETWORK。
docker compose --project-directory . -f deploy/docker-compose.example.yml up -d --build
```

示例加入已有的 `sub2api-production_default` 网络，通过 `cpa:8317` 访问 CodexProxy。容器名或端口不同，请修改 Compose 中的地址。镜像从本项目本地构建，不会拉取原版 Keeper。`./data` 持久化保存数据库、日志和备份。

建议使用新的数据目录。本项目不提供旧 Keeper 历史导入，也不会自动清空已有数据库；保留历史数据库迁移以便安全升级，因此复用旧数据库可能仍会展示其中已有的其他供应商数据。

**采集不需要域名。** 页面默认只映射到宿主机 `127.0.0.1:8080`，可用 SSH 隧道访问，也可通过 HTTPS 反向代理提供远程访问。独立域名是可选项；使用 `/keeper` 子路径时设置 `APP_BASE_PATH=/keeper`，反向代理应保留此前缀。`CPA_PUBLIC_URL` 用于页面中的“返回 CPA”链接。

普通 HTTP 反向代理不能转发 RESP 数据流，采集地址保持 Docker 内网的 `REDIS_QUEUE_ADDR=cpa:8317`。

如需与 CodexProxy 使用相同管理密码，可将 `LOGIN_PASSWORD` 留空，把 CodexProxy 的 bcrypt `remote-management.secret-key` 填入 `LOGIN_PASSWORD_HASH`。在 `.env` 中用单引号包住哈希值，避免 `$` 被替换。两种密码配置只能选择一种；以后修改 CPA 密码时，也需要更新 Keeper 的哈希。

## 采集机制

先订阅 `usage`，再用 `LPOP usage` 补收积压记录；断线后重复此流程。有效记录先写入 SQLite 收件箱，再做统计。数据库写入失败时重试已经收到的同一批记录，初始间隔由 `REDIS_QUEUE_RETRY_INTERVAL` 控制，默认 1 秒。

上游协议没有持久化确认机制，进程在接收与落库之间崩溃仍可能丢失记录。错误事件只有实时推送，没有积压补收。请保持采集器持续运行。

## 本地开发与检查

需要 Go 1.26+、支持 CGO／SQLite 的 C 编译器和 Node.js 24。

```sh
npm --prefix web ci
npm --prefix web run build
go run ./cmd/server --env .env --host 127.0.0.1

go test ./cmd/... ./internal/...
go build -o /tmp/codexkeeper ./cmd/server
npm --prefix web run typecheck
npm --prefix web run lint
npm --prefix web test
```

Node.js 26 的内置 localStorage 与现有 happy-dom 测试环境冲突，测试时使用：

```sh
NODE_OPTIONS=--no-experimental-webstorage npm --prefix web test
```

为减少无关改动，Go 模块名和内部二进制名称仍为 `cpa-usage-keeper`；项目链接和版本检查已指向 `vpromise/CodexKeeper`。

MIT 协议，保留原项目贡献者版权，见 [LICENSE](LICENSE)。
