# CodexKeeper 上游同步台账

本项目是 `vpromise/CodexKeeper`，上游为 `Willxup/cpa-usage-keeper`。采用选择性适配，仅保留原生 Codex／Claude；审查边界不等于完整合并或全部实现。时间均为 `Asia/Shanghai`。

## 当前检查点

| 项目 | 已核实状态 |
| --- | --- |
| 本次检查时间 | 2026-11-08 05:09:06 +08:00 |
| 发布目标 | `https://github.com/vpromise/CodexKeeper` 的 `main`，origin fetch/push 均为此地址 |
| 本次隔离候选基线 | [`206182621355bc98575e6470c5bc84e7da2b5985`](https://github.com/vpromise/CodexKeeper/commit/206182621355bc98575e6470c5bc84e7da2b5985)，与开始时远端、本地 main 相同 |
| 最新正式发布／实际集成边界 | 最新为 [v2.0.3](https://github.com/Willxup/cpa-usage-keeper/releases/tag/v2.0.3)，`9db74cfc380d11064af80cf1626efd6e3d068217`；本地仍为原生精简 v1 语义，历史选择性适配加 #670 原生部分／#672／#676，**不是 v2 升级** |
| 本次完整审查范围 | [`9db74cfc380d11064af80cf1626efd6e3d068217..aa687529e392833e17ad0eb7ccfc4dec8eda0982`](https://github.com/Willxup/cpa-usage-keeper/compare/9db74cfc380d11064af80cf1626efd6e3d068217...aa687529e392833e17ad0eb7ccfc4dec8eda0982)：38 个新增提交，全部未发布；正式版新增 0；另复核历史待决 |
| 主分支已审查检查点 | `aa687529e392833e17ad0eb7ccfc4dec8eda0982`；已全部分类，不代表全部实现或完成未采纳项的集成验证 |
| 本次实际改动 | 仅更新完整范围与逐项决定的台账；无应用代码或新依赖适配。上次三组修复已发布于 `20618262` |
| 待确认事项 | 原 9 项全部延续，没有新增用户决策；未发布的独立修复另按技术延期跟踪 |

本轮重新查询官方 Releases 列表与 latest：最新仍是 2026-10-30 17:13:11 +08:00 的 v2.0.3，非 draft、非 prerelease；没有新增中间正式发布。Git ref API 与独立 fetch 均解析为 commit `9db74cfc380d11064af80cf1626efd6e3d068217`，没有使用 `target_commitish`。main 为 `aa687529e392833e17ad0eb7ccfc4dec8eda0982`，v2.0.3 是其祖先；38 个增量均只在 main。独立 refs 为 `refs/upstream-sync/cpa/main`、`refs/upstream-sync/cpa/releases/v2.0.3`，不覆盖 fork tags、不改 origin。历史选择性边界不变，当前决定以 2026-11-08 小节为准。

## 有来源说明的历史基线

运行前没有 `docs/upstream-sync.md`。以下为历史选择性集成证据，本次没有重新逐提交审计全部早期历史，也没有追认未执行的历史测试。

- [`f7ea9f75b56abb904d7dbedbcb19f3158928fde5`](https://github.com/vpromise/CodexKeeper/commit/f7ea9f75b56abb904d7dbedbcb19f3158928fde5) 建立原生 Codex／Claude 精简范围。`internal/providers/providers.go` 当前只认可这两家，显式不支持的 provider 不因模型或 executor 名称被改认；provider 缺失时仅允许已知原生 executor。
- [`012609760c426382cec70015dabc043b4718779b`](https://github.com/vpromise/CodexKeeper/commit/012609760c426382cec70015dabc043b4718779b) 的提交、两份 README 和实际实现共同支持“相关更新选择性适配至 v1.15.10”。它没有把 v1.15.10 的全部 ancestry 合并到 fork。
- 历史审查聊天（`01a11995-5112-7353-97e7-88db16164f42`）记载 v1.15.7 `dab6760c66c80b3d32a7b649cfcfe3bf1d7ee4a1` 至 v1.15.10 的 23 项主线审查；用户随后要求“按你的分析开始实施同步更新和修复”。这只是继承的审查／授权证据，本次新完成的完整审查范围见上表。

| 历史来源 | 当前实现证据及保留的选择 |
| --- | --- |
| [#577](https://github.com/Willxup/cpa-usage-keeper/pull/577)、[#583](https://github.com/Willxup/cpa-usage-keeper/pull/583) | `internal/quota/codex.go`、订阅到期并发测试、`internal/repository/usage_identities.go` 和 metadata diff 测试；保留到期日期防旧快照覆盖、别名／统计字段与原生 provider 范围 |
| [#584](https://github.com/Willxup/cpa-usage-keeper/pull/584)、[#593](https://github.com/Willxup/cpa-usage-keeper/pull/593)、[#587](https://github.com/Willxup/cpa-usage-keeper/pull/587)、[#611](https://github.com/Willxup/cpa-usage-keeper/pull/611) | Claude 5h／7d 历史、部分窗口合并、Max 5x／20x 标识，以及已更新的 Codex／Claude 请求头均在当前 quota 实现和测试中；本次不再次改变客户端指纹 |
| [#575](https://github.com/Willxup/cpa-usage-keeper/pull/575)、[#581](https://github.com/Willxup/cpa-usage-keeper/pull/581)、[#582](https://github.com/Willxup/cpa-usage-keeper/pull/582)、[#602](https://github.com/Willxup/cpa-usage-keeper/pull/602) | 统一编辑、焦点恢复、原生筛选 URL 已在 credentials 中；优先级仅限认证文件，API Key 只提供别名／启停。CodexProxy API Key PATCH 无 `priority`，此限制继续保留，非待批准的跨项目扩展 |
| [#578](https://github.com/Willxup/cpa-usage-keeper/pull/578)、[#579](https://github.com/Willxup/cpa-usage-keeper/pull/579)、[#580](https://github.com/Willxup/cpa-usage-keeper/pull/580)、[#585](https://github.com/Willxup/cpa-usage-keeper/pull/585) | 已有堆叠趋势、Top5／Others 色带、`latency_scatter` 和 1000 点采样优化；本地采样迁移保留备份、事务、计数／最大值／sketch，不能以整文件覆盖丢掉保护 |
| [#573](https://github.com/Willxup/cpa-usage-keeper/pull/573) | `internal/service/redis_usage.go` 已校验 `resolved_client_ip`，不合法时退回 `client_ip`；当前 CodexProxy 队列尚无该字段，属于已有兼容性预留，并未自行信任 X-Forwarded-For |
| [#603](https://github.com/Willxup/cpa-usage-keeper/pull/603)、[#613](https://github.com/Willxup/cpa-usage-keeper/pull/613) | 归档期限能力已实现，默认 `0` 永久保留；只有至少 90 的配置启用归档删除。Claude reset grant 必须显式选择确认，不自动重试结果不明的兑换。本次没有改变保留配置或消费权益 |
| [#574](https://github.com/Willxup/cpa-usage-keeper/pull/574)、[#592](https://github.com/Willxup/cpa-usage-keeper/pull/592)、[#604](https://github.com/Willxup/cpa-usage-keeper/pull/604)、[#612](https://github.com/Willxup/cpa-usage-keeper/pull/612) | 当前有标题／徽标修复、Alpine 3.24 runtime、构建上下文排除与虚拟时间聚合防抖测试；本次未调整 Docker 或调度逻辑 |
| [#609](https://github.com/Willxup/cpa-usage-keeper/pull/609)、[#589](https://github.com/Willxup/cpa-usage-keeper/pull/589) | 国际 Kimi 整体排除；上游完整部署方案、非原生界面／截图不覆盖本地 README。社区排行、用户登录、非原生凭证及旧采集降级继续排除 |

以上具体方案已落实在当前代码和历史 README。未发现需要承接的 Keeper 未提交适配或未完成的明确批准事项；同聊天中的 CodexProxy 冷却／插件决策不属于 Keeper 台账。

历史聊天还报告 `01260976` 曾在 dmit／dmit2／bawg 部署并保留备份／回滚资料；本次没有连接服务器，未独立核实这些历史运行状态，不把此次台账推送视作部署。

## 2026-10-10 新范围的逐项决策

审阅了完整提交列表、实体／解析／迁移／归档／benchmark 列契约、共享 quota hooks、native callers、xAI 条件路径，以及相应 PR 理由和新增回归。合并提交与原始提交成组去重；此范围没有 revert。下表列出全部 7 个 SHA，避免把功能分支提交与 merge 重复计算为新功能。

| 上游 SHA／PR | 分类、当前行为、依赖与证据 | 本次决定及重新评估条件 |
| --- | --- | --- |
| [`5853b7ee2079f63f123fae84b1e1c796e68b4a60`](https://github.com/Willxup/cpa-usage-keeper/commit/5853b7ee2079f63f123fae84b1e1c796e68b4a60)、[`6c3b7d71a84f41b3e6a5567987dc94fc10a96d81`](https://github.com/Willxup/cpa-usage-keeper/commit/6c3b7d71a84f41b3e6a5567987dc94fc10a96d81)、merge [`2c09d65d6a90b1dc9802fe7b26704a1d95d0c909`](https://github.com/Willxup/cpa-usage-keeper/commit/2c09d65d6a90b1dc9802fe7b26704a1d95d0c909)；[#619](https://github.com/Willxup/cpa-usage-keeper/pull/619) | **适用，延期至正式发布。** 当前 `useCredentialsTabData.ts` 的 `staleFailedState` 会在存在旧 quota 时压制刷新失败；当前缓存与手动刷新分开存储。上游合并为同一快照、用 revision／working 集合防迟到缓存覆盖，并清理离页缓存错误。需成组适配 `useQuotaCache`、`useQuotaRefreshTasks`、tab wiring 和 recovery 测试；`6c3b7d71` 为依赖的离页修正。无新后端接口或刷新频率变化，原生两家均受益 | 尚未正式发布，按默认稳定版策略本次不移植。正式发布后重新核对最新 follow-up、为现有原生 hooks 适配错误／恢复／换页／缓存竞态测试。收益是可靠显示认证失败；风险是任务状态和缓存清理回归。没有待批准的 provider、统计或架构变化 |
| [`08ec76afd994bd67d6493c5241a87167b807f692`](https://github.com/Willxup/cpa-usage-keeper/commit/08ec76afd994bd67d6493c5241a87167b807f692)、merge [`41a9b90cfcee50ea5a43f6ad9cfb25ac4d5f1bd0`](https://github.com/Willxup/cpa-usage-keeper/commit/41a9b90cfcee50ea5a43f6ad9cfb25ac4d5f1bd0)；[#623](https://github.com/Willxup/cpa-usage-keeper/pull/623) | **原生路径可适用，延期。** Keeper 当前仅解析已有请求／会话字段，未保存 `execution_id`、`trace_id`、`node_kind`、`is_fork`、`is_compaction`。上游为 hot/archive 增加 nullable 列，更新 storage columns、解析、migration 注册及 benchmark 参数；旧行保留 NULL、false 保留为 false，不改变 event identity、token、聚合或 retention。显示／导出未在此 patch 中完成。当前 CodexProxy `internal/redisqueue/plugin.go` 的 JSON 序列化未输出这 5 个字段，所以现在添加列不会给实际采集带来新元数据 | 未发布且实际生产者无字段，延期到正式发布并核实对应 CodexProxy 输出后再评估。只考虑兼容性增列与保存已经提供的可选字段；不跨项目修改生产者、不补造历史、不更改去重／费用定义。后续必须验证旧／部分 schema、索引和 trigger 保留、hot/archive/replay、UTC／上海时区和同 request 的多个事件；采用旧库备份及回退评估。单纯可选字段保存与兼容性增列不自动构成需要用户批准的破坏性迁移 |
| [`51071027b28f00f4782d4bd508d2a55f8a662553`](https://github.com/Willxup/cpa-usage-keeper/commit/51071027b28f00f4782d4bd508d2a55f8a662553)、merge [`778b5f8e7196f3fb1eae8f11b63f51de33ea80ed`](https://github.com/Willxup/cpa-usage-keeper/commit/778b5f8e7196f3fb1eae8f11b63f51de33ea80ed)；[#625](https://github.com/Willxup/cpa-usage-keeper/pull/625) | **排除。** xAI 产品百分比改为共享周额度明细，旧产品行不参与耗尽判断，新增 `usageBreakdown` 与 tooltip。当前无 xAI provider／quota registry／凭证。虽触及共享 `QuotaRow`、`PortalTooltip`、native 共用 quota panel，但新增分支仅在 `usageBreakdown` 有值时运行，当前 Codex／Claude 正常响应无该字段。共享 helper 的变动没有独立的原生修复 | 精简范围内永久排除此功能与仅服务它的共享类型／tooltip 改动；不恢复 xAI。只有用户另行修改 provider 范围或后续出现原生 usage breakdown 的具体需求，才重新评估，不把整包 patch 自动引入 |

## 2026-10-10 CodexProxy 兼容性核实（历史）

只读检查 `/Users/vpromisex/Files/Code/RPCodex/CodexProxy`，HEAD 为 `2663a7456ca284d9d05daf8963dd66bf2cd8bd42`，工作区干净。本次没有修改、fetch 或推送 CodexProxy。

- `internal/api/redis_queue_protocol.go` 支持 `AUTH`、`SUBSCRIBE usage`、`LPOP usage`，与 Keeper `internal/cpa/endpoints.go` 一致；当前单采集器、先订阅后补收、失败重试既有批次等行为未变。
- `internal/redisqueue/plugin.go` 输出 native provider／executor、token breakdown、请求头及请求 ID；Claude `response_headers` 通道和现有 accounting 输入继续可用，但上述五个执行／节点字段和 `resolved_client_ip` 未包含在队列 JSON 中。
- `internal/api/server_management.go` 保留 `api-call`、`reset-quota`、Codex／Claude key 路由、`auth-files/fields`。`config_lists.go` 的两家 API Key PATCH 仍无 `priority`，现有认证文件独享优先级编辑限制仍正确。
- 结论为当前本地源码契约兼容；没有把源码核对称作服务器端到端验证，没有读取线上配置或调用真实管理写接口。

## 2026-10-10 验证与输入身份（历史）

本次发布只增加文档；为首次运行检验工作流，额外对未改变的应用基线实际执行完整本地前后端验证。所有测试使用现有临时 SQLite／mock fixtures；隔离 worktree 无生产 `.env`、数据库或备份，不连接真实采集队列，不写真实凭据、不使用额度重置券。

环境：macOS／darwin arm64，Go `1.27.1`（模块要求 1.26），Node `26.8.1`，npm `11.19.0`，Apple clang `21.0.0`，`CGO_ENABLED=1`。因现有 README 记载的 Node 26 localStorage／happy-dom 冲突，Vitest 使用 `--no-experimental-webstorage`；没有跳过 Vitest 或修改测试断言。

| 实际命令 | 退出码／结果 |
| --- | --- |
| `npm --prefix web ci` | 0，使用 checked-in lockfile，239 个依赖包；未升级依赖。npm audit 提示 8 项既有漏洞（2 moderate／4 high／2 critical），本次没有进行依赖修复或完整漏洞研判 |
| `NODE_OPTIONS=--no-experimental-webstorage npm --prefix web run test -- --reporter=json --outputFile=<本次 work>/keeper-vitest.json` | 0，155 个文件，1186 tests passed，0 failed／pending |
| `npm --prefix web run lint` | 0 |
| `npm --prefix web run typecheck` | 0；这不能替代对测试文件的 Vitest 执行 |
| `npm --prefix web run build` | 0，真实候选前端 208 modules；产生 HTML、CSS、JS；有既有 >500 kB chunk 提示，未改构建设置 |
| `CGO_ENABLED=1 TZ=Asia/Shanghai go test -count=1 -json ./cmd/... ./internal/...` | 0，67 个包中 29 个有测试包通过、38 个无测试包；2417 tests/subtests passed，0 failed，1 项现有平台 skip；没有复用 cached 结果 |
| `CGO_ENABLED=1 go build -o <本次临时目录>/cpa-usage-keeper ./cmd/server` | 0，使用上述真实 `web/dist`，非占位前端、非 CGO-disabled 构建 |
| `<本次临时目录>/cpa-usage-keeper --version` | 0，输出 `dev`；只验证可执行文件启动，不加载服务配置 |
| `git diff --check`、显式暂存文档与发布历史检查 | 发布前执行，拒绝无关源代码、build 产物或 user 改动进入提交 |

Go test 内唯一跳过项是 `internal/logging/test.TestDedicatedErrorMaintenanceClosesStaleDailyFile`：该现有测试依赖 `/proc/self/fd`，macOS 不提供。SQLite、migration、archive、quota、native provider、inbox mock tests 均包含于实际后端检查。没有为本轮添加测试排除。

构建删除了 tracked `web/dist/.gitkeep`；已原样恢复该空文件，并在最终 assets 状态下重新执行 Go 全套测试及构建。没有将该构建副作用提交。最终应用源码与测试对应基线 tree `afd82fb4c241b807415c934fdc1f781b15989989`；唯一拟发布差异为本文件。

本次保留的验证 manifest 为 `<本次 work>/keeper-validation-manifest.json`（本次 work：`/Users/vpromisex/Files/Codex/2026-10-10/new-chat/work`），逐文件列出 1061 个源码／测试／依赖锁／构建配置输入和 4 个 embed 文件 SHA-256。源码 manifest digest：`ad58b222c8cf00ae69203b40c3f374d877bf56efa95ad45eaecd2c9b40c88933`；assets manifest digest：`1daa24dacbf6bea057c015c6abd0e4b25cee72e32c1fe3ff88d445b14f347817`。二者按 path 排序后以紧凑 JSON 编码并做 SHA-256，不含环境文件、runtime 数据或缓存。

| 嵌入文件 | SHA-256 |
| --- | --- |
| `web/dist/.gitkeep` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `web/dist/index.html` | `299d069bc9a5fbc69554c32986b91b03543d0ad24bdd679c24b2b05d115cea50` |
| `web/dist/assets/index-6LjvyvW3.css` | `74a9df69009b44a3ff8d8543b7ab62d6f56504512aae0f6866235bd8d03a505b` |
| `web/dist/assets/index-BhmepEaZ.js` | `fd691eddb00a01300687e754f2533f65df5f2b91c727b57495b1fc37146db95f` |

未覆盖：Linux／Windows runtime、该 `/proc` 日志检查、Docker build（未改容器相关输入）、race（无并发代码修改）、浏览器视觉交互（无 UI 修改）、线上端到端／真实账号／真实额度重置，以及未采纳 patch 的集成验证。没有以历史测试代替此次执行；历史 deployment 状态也未验证。

## 2026-10-10 发布与运行状态（历史）

- 仓库锁为 resolved common Git directory 下 `codexkeeper-upstream-sync.lock`，由本次长生命周期进程持有，记录 host、PID、run ID、上海时区时间与隔离 worktree。开始时没有其他 Keeper 同步活跃聊天或既有锁；用户 checkout 无未提交改动、无远端以外本地提交。
- 本次使用独立分支 `sync/keeper-20261010` 和 disposable worktree；未 stash、reset、clean 用户 checkout，没有导入上游发布历史。
- 当前 workflow 定义中，main push 仅触发 CI；二进制／镜像发布由版本 tags 或手动 dispatch 触发。本次只正常推送台账文档，不创建 tag／Release，不 dispatch、不改 workflow 或仓库策略。
- 开始时 `01260976` 的 GitHub Actions run 查询为 0；不能把“CI workflow active”视作已通过。此次文档提交后的远端 SHA、check-runs、status 与 Actions 查询结果在当次运行报告记录，避免为自引用 SHA 创建循环台账提交。
- 本次没有服务器访问或部署。推送台账不会更新运行中的 Keeper 服务。

## 后续运行规则与待确认事项

后续检查新正式版、上游 main 和本表延期项；即使主分支 checkpoint 已推进，也必须重新处理 #619／#623，直到实现或形成新的有依据决定。没有新值得采纳的内容时，不创建空提交或仅更新时间的提交。

2026-10-10 时无待确认事项；当时 #619／#623 等待正式发布，xAI／Kimi 等为范围外排除。后续正式发布与当前待确认事项以以下运行记录为准，不把这些历史延期判断误当作已完成。

## 2026-10-18：选择性适配 v1.15.11

### 仓库、覆盖范围与独立批次

本轮开始时 Keeper 为干净 main，无无关本地提交、未提交改动、既有同步锁或另一项活跃 Keeper 同步。读取了最新 skill、三份 references、两份 README、完整台账；路径上的各层及仓库内未发现适用 AGENTS.md。另一项 CodexProxy 同步聊天活跃，但属于不同仓库，本轮只读核实其接口。

由长生命周期 PID 9749 持有 common Git directory 下的原子 lock，run ID `keeper-20261018-01a1242c`，host `Wills-Mac-Studio.local`，开始于 05:04:51 +08:00。在 `sync/keeper-20261018` 的 disposable worktree 实施，基线为 `ca727641`。没有将上游 ancestry 合并进 fork，也没有借用 CodexProxy 的发布授权。

完整范围为 v1.15.10 `b3592bff` 至上游 main `aa6056bf`，23 个实际提交；其中 v1.15.11 正式版包含 16 个，之后未发布 7 个。相对台账旧检查点 `778b5f8e`，本轮新增审查 16 个，并重新评估此前 7 个提交。所有 merge 与其功能分支去重，合并后的差异也已核实；本范围没有 revert。以下十项逻辑变化覆盖全部范围：

| 来源与提交 | 决策、依赖及本地行为 |
| --- | --- |
| [#619](https://github.com/Willxup/cpa-usage-keeper/pull/619)：`5853b7ee2079f63f123fae84b1e1c796e68b4a60`、`6c3b7d71a84f41b3e6a5567987dc94fc10a96d81`、merge `2c09d65d6a90b1dc9802fe7b26704a1d95d0c909` | **已适配，解除等待正式发布。** 统一 quota 数据与任务状态，以 revision 防迟到缓存覆盖手动刷新，清理离页的缓存来源错误；保留任务、瞬时错误与成功快照的不同生命周期。保持一分钟轮询和本地 Claude reset、认证文件独享 priority 边界。原 helper 的旧快照压制错误断言改为失败可见断言，并以 26 项实际 hook／panel 回归覆盖错误、恢复、换页、缓存缺项和竞态 |
| [#623](https://github.com/Willxup/cpa-usage-keeper/pull/623)：`08ec76afd994bd67d6493c5241a87167b807f692`、merge `41a9b90cfcee50ea5a43f6ad9cfb25ac4d5f1bd0` | **仍延期，正式版条件已满足，生产者条件未满足。** 当前 CodexProxy 队列不输出 execution／trace／node／fork／compaction 字段；本轮无收益的增列不入批次。实体、storage columns、归档／replay、benchmark、nullable 幂等迁移与 false／NULL 语义须成组适配，不能只复制 decoder。等待实际支持字段后重新评估，不补造历史，不改变去重／费用 |
| [#625](https://github.com/Willxup/cpa-usage-keeper/pull/625)：`51071027b28f00f4782d4bd508d2a55f8a662553`、merge `778b5f8e7196f3fb1eae8f11b63f51de33ea80ed` | **继续排除。** xAI shared weekly product breakdown 无原生 quota caller；共享类型、tooltip 的新增分支也只服务该字段。保留当前 Codex／Claude registry，不引入 xAI 专属 UI 或 usageBreakdown |
| [#630](https://github.com/Willxup/cpa-usage-keeper/pull/630)：`d51aaa776ee0edf51bd05aa641e6fa9240901935`、merge `2cb53b46deb4fc06cc6b8182676101f5285ade48` | **待确认，未纳入代码。** 上游把分布图桌面 300→240 px、移动 250→200 px，去掉独立日期行与图例，把汇总条目改为曲线开关；同时改变点击时 tooltip 生命周期。属于可选布局／交互方案而非必要数据修复，与独立四项无依赖。具体选择见下文 |
| [#631](https://github.com/Willxup/cpa-usage-keeper/pull/631)：`2956fd9533d025a9f22a6c155b505edc8a612667`、notice `c4d67f5214e10d2e057de841712ff2683712677f`、merge `9baffb0add6003adbab7ed84d5f7ac3f1003313c` | **已适配到现有请求表和凭证请求历史。** 流式图标依赖明确 `stream=true`；false／缺失显示非流式／未知，绝不从 SSE／WS 文本推断。可聚焦 tooltip 展示原始完整端点，Effort 仅在凭证详情中显示，虚拟行卸载清除 tooltip。附带 Microsoft Fluent 图标 MIT notice，并验证 public→dist→Go embed 一致；未恢复用户视图或改变页面整体布局 |
| [#632](https://github.com/Willxup/cpa-usage-keeper/pull/632)：`5564099e690d5d97862d7dc310d2834e125d2747`、merge `1b9875f0758c22fc63b960289c2a7d34c5bb060b` | **已适配。** 在 User-Agent 截断文本旁增加复制完整值及成功／失败反馈，缺失值不显示按钮；抽取现有日志复制 helper，保留 HTTP textarea fallback 与焦点恢复。与 #631 共用请求卡，但语义独立；对应客户端 metadata 和既有日志复制测试均通过 |
| [#637](https://github.com/Willxup/cpa-usage-keeper/pull/637)：`d967a44b53cb26e685e9b99553c16d5ff68d266d`、merge `71a41cc2d1cce8d75b42473adf31f2616ad8085e` | **按原生 handler 适配。** 明细全零且报有非零 Total、未发生负数裁剪时，保留 Total。合成 2120-token 图片事件证实此前入库为 0，修复后事件和 Overview 保留 2120，已配置单价时成本仍为 0，不补造输入／输出、不回填历史。负数损坏规则、cached-only 本地兼容规则继续保留；处理器测试只覆盖三种已支持的 Codex executor，不恢复 xAI／Meta registry |
| [#640](https://github.com/Willxup/cpa-usage-keeper/pull/640)：`adabf45b52fb6da39ed1a76c2baaf36499c7323d`、merge `8c5ab2e6458286c475fc137a72d075a389481a8d` | **未发布，排除。** 为上游仓库增加 Star History 外部图表，与本项目原生能力无关，不覆盖 fork README 或增加外部图表请求 |
| [#641](https://github.com/Willxup/cpa-usage-keeper/pull/641)：`6d6e69aa1bcee4a1040061dce13276238a5ee95e`、merge `6b8e33c2476d6cbb8818ec653331979824681cc3` | **未发布，随 #640 排除。** 只改变 Star History 的 timeline/date 参数，依赖未采纳的图表，无独立原生修复 |
| [#644](https://github.com/Willxup/cpa-usage-keeper/pull/644)：`a4696c80651e97cbcb2103080489a6fc70715fd8`、`6e4cef80c742333312640a62b41a73cb3c5dcf02`、merge `aa6056bf34667f28549691ef893d75a6fc0d47de` | **未发布且与当前协议不兼容，延期。** 改为强制 execution_id，缺少 request_id 时用它作 event_key；同时按短事务扫描 hot/archive 回填空 event_key、将旧 request-ID 错误 inbox 重新排队。当前生产者不提供 execution_id，整包移植会拒绝现有原生消息。需要 #623、实际生产者支持、保留旧消息兼容的明确方案，以及迁移／重放防重复和回滚评估后才有采用条件。当前维持 request_id 协议，不执行历史恢复，不把发布或用户沉默当作允许修改协议／历史的授权 |

### 当前 CodexProxy 源码契约

只读检查当前 `codex/user-provider-isolation`，HEAD `605896e1b7a3637665d31e1df93adf8ecdcc67fa`，工作区干净。接口快照中 `internal/redisqueue/plugin.go` blob `d91c8a2820cdea5a5212b1bb431d7effe25a8a18`，仍输出 request_id、provider／executor、token breakdown、response_headers，未输出 #623 五字段、resolved_client_ip 或 stream。实际缺失 stream 的记录在新 UI 中保持“未知”；预览中的 true／false 为合成数据，不代表线上已支持该字段。

AUTH／SUBSCRIBE usage／LPOP usage 及 api-call、reset-quota、auth-files/fields 路由仍可对齐。两家 API Key PATCH 仍无 priority，本地 `setAuthFilePriority` 及认证文件专用写入边界已保留，未被上游的泛化 `setCredentialPriority` 恢复。没有加载 CodexProxy 配置、访问生产接口或修改该仓库。

### 实际验证及输入身份

环境为 Go 1.27.1／Node 26.8.1／npm 11.19.0、darwin arm64，所有 Go 检查显式 `CGO_ENABLED=1`。隔离 worktree 无生产环境文件和数据库。两项实质 bug 均先复现后修复：Total-only 的三种 Codex executor 和 SQLite 入库回归先失败；quota helper 的旧快照压制失败先复现，再通过真实 hooks 的完整 recovery 场景确认修复。没有为发布改变负数保护或隐藏失败。

| 检查 | 实际结果 |
| --- | --- |
| `npm --prefix web ci` | exit 0，使用原 lockfile，未升级依赖；仍提示既有 8 项 audit 问题（2 moderate、4 high、2 critical），本轮未进行安全研判或依赖修复 |
| 定向 Vitest（quota recovery／wiring／reset、请求 cell、UA、凭证请求历史） | exit 0，6 文件／97 tests；其中 recovery 26、request cell 22 |
| `NODE_OPTIONS=--no-experimental-webstorage npm --prefix web run test -- --reporter=json --outputFile=<work>/keeper-vitest-20261018.json` | exit 0，157 文件、1237 tests passed、0 failed／pending；沿用已验证的 Node 26 happy-dom 调整 |
| `npm --prefix web run lint`、`npm --prefix web run typecheck` | 均 exit 0 |
| `npm --prefix web run build` | exit 0，212 modules，真实候选 HTML／CSS／JS 与 notice；已有大 chunk 提示未隐藏 |
| `CGO_ENABLED=1 TZ=Asia/Shanghai go test -count=1 -json ./cmd/... ./internal/...` | exit 0，29 个有测试包，2429 tests/subtests passed、0 failed；无 cached 结果，1 项现有 `/proc/self/fd` 日志测试在 macOS skip |
| `CGO_ENABLED=1 TZ=Asia/Shanghai go test -race -count=1 -json ./internal/service/tokenprocessor/... ./internal/service/... -run 'TestResponses\|TestProcessRedisUsageInboxPreservesImageTotalOnlyUsage\|TestUsageIngestionKeepsOnlyNativeProviders'` | exit 0，15 tests/subtests passed、0 failed／skip，无 cached 结果；为定向检查，不声称全量 race |
| `CGO_ENABLED=1 go build -o <temporary>/cpa-usage-keeper ./cmd/server`，随后 `--version` | 均 exit 0，输出 dev；编译前已真实构建前端并原样恢复空 `.gitkeep`，使用最终 embed 输入 |
| public／dist notice `cmp`，`CGO_ENABLED=1 go run <work>/keeper18-embed-probe.go` | exit 0，Go 的 `web.Static` 读取到与原 notice 完全一致的 1268 bytes |
| 本地合成预览 | macOS 内置浏览器浅色／深色、390×844 窄屏，完整端点 focus tooltip、14／13 px 图标、error／recovery 与 mock clipboard 成功反馈；clientWidth／scrollWidth 均 382，tooltip 在视口内；console errors 为空。无真实后端请求或系统剪贴板写入 |
| `git diff --check`、选定路径和暂存一致性 | 发布前再次检查；预览 HTML／TSX 已删除，Vite 已停止，未纳入发布；分布图、provider、schema／migration 和采集协议文件没有更改 |

完整 JSON 输出及 `<work>/keeper-validation-manifest-20261018.json` 保存在本聊天 work 中，1069 个源码／测试／lock／构建配置文件包含新增文件；按路径排序的紧凑 JSON SHA-256 为 `5b5c8edf550e55478facd38a547c0f3f603337dc054a10f0b1458f500388bde7`。最终 5 个 embed 文件 manifest SHA-256 为 `eaa0ab064e2955692e1f1c6955c282460b921186b538fe774db5e6cc72f99835`。候选身份为 `ca727641` 加该完整源码 manifest 与本次选择性差异；README／台账后续编辑不改变已验证 runtime 输入。最终暂存检查移除了上游 clipboard helper 的 EOF 空行，随后重新跑前端全部检查和构建；五个 embed 文件与 Go 检查时逐字节一致，Go 构建也再次通过。

| 最终 embed 文件 | SHA-256 |
| --- | --- |
| `web/dist/.gitkeep` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `web/dist/index.html` | `11d62b4ada1f83fa33a765aeb7854669f345f98ba3b7d927edc0a63df0cb1d82` |
| `web/dist/assets/index-C0dEZpUG.js` | `a10028244492143e611a7b6f148ae0846f4e6a5c4eab29dbf7e210a114ddb118` |
| `web/dist/assets/index-OFLChr_n.css` | `834e7a12dca577c92a818e690432f1397ddb55180e128afa8cef8cd995fa2af6` |
| `web/dist/third-party-notices.txt` | `3405117543d2887d5244eab04fe389d046a9b6c89a7cafa2863fed95931ad22f` |

未覆盖：Linux／Windows runtime、macOS 跳过的日志文件描述符检查、Docker build（无容器配置／依赖变更）、全量 race、真实 CPA／账号／生产队列／数据库／重置券，以及未采纳 #623／#630／#644 的集成和历史恢复。没有执行历史重算、删除或真实凭据写入。

### 发布边界与待确认事项

仅发布上述四项适配、必要测试、图标 notice 与 README／台账。正常推送到已核实的 fork main，提交 SHA、远端回读和对应 CI 结果写入当次运行报告，避免台账自引用提交号。现有 workflow 的 main push 为 CI；tag／手动发布工作流未触发，未改策略、创建 Release、发布镜像或部署。历史服务器记录保留，不能将本轮发布当作服务器升级。

**待确认 1：#630 用量分布图布局与交互。** 当前为桌面 300／移动 250 px，保留独立日期行、独立图例，汇总条目只打开明细；上游建议桌面 240／移动 200 px，移除日期行和重复图例，汇总条目同时切换曲线并显示 tooltip。收益是更紧凑、减少重复元素；风险是日期范围不再独立显示，原本只查看汇总的点击会隐藏／显示曲线。建议先保持现有界面，待用户明确选择后再成组采用上游方案；若批准，需要三语言、触摸／键盘切换、数据／维度更新、浅深色及桌面窄屏回归。没有 DB 迁移，回退为恢复这组前端文件并重新构建。需要用户决定：保留当前布局，还是批准上述完整布局／交互变更。

该项遵循 `codexkeeper-upstream-sync/SKILL.md` 引用的 `references/fork-contracts.md`：**“Preserve the existing layout, theme and responsive behavior unless the change itself is approved.”** 没有将上次分析、台账记录或未回复视作批准；此项不阻塞四项独立修复的发布。

#623 是已知生产者依赖的延期，不需要现在决定；#644 未发布且不能兼容当前生产者，等待正式版及保留旧消息兼容的具体方案，届时如涉及协议／历史恢复，须按 skill 的讨论边界另行确认。xAI 和上游 Star History 为明确排除，不列为待批准项。除 #630 外无其他待确认事项。

## 2026-10-25：审查 v2.0.0／v2.0.1，独立修复焦点与归档

### 仓库与完整范围

开始时本地及远端 main 均为干净的 `f21197c8ec2fa8f020f3a3cb202bd39070434ef3`，ahead/behind 0/0，没有 Keeper 未提交适配、无关本地提交或另一项活跃 Keeper 同步。已读取最新 skill、三份 references、README、README.zh、完整台账；仓库和适用父目录无 AGENTS.md。以 resolved common Git directory 内的原子锁串行运行：PID 25246、host `Wills-Mac-Studio.local`、run `keeper-20261025-01a1242c`、05:04:15 +08:00 开始。长生命周期 owner 不是已退出的 shell；实施在隔离 `sync/keeper-20261025` worktree，不 stash/reset/clean 用户 checkout。

本轮完整枚举旧 main 检查点 `aa6056bf..5343b110` 的 **220 个新提交**，159 个进入 v2.0.0、随后 56 个进入 v2.0.1、其后 5 个仅在 main。v1.15.11 `71a41cc2..5343b110` 共 227 个，差额 7 个是上次已审查的 Star History／execution recovery；本次连同 #623／#630 一起复核。已查看完整提交列表、发布合并树、各逻辑组差异及 PR 理由／测试依赖，并区分共享原生代码和被移除路径。没有实际 revert；搜索命中的 pill-button 提交只是说明可独立回退，后续 #649 确实删除了 #642 的独立 restore 工具，故按最终树决定，不能恢复中途产物。

发布标签直接为 commit，GitHub Git ref API 与独立 fetch refs 一致，未覆盖 fork tags。v2.0.0 promotion 的树等于第二父 `4a9a14ab`，无额外整合改动。所有 merge／子分支按下表去重；`ad911c7a`、`ca63512b`、`555f8da1`、`66718473`、`263f5899`、`85da0230`、`dd7a7216`、`dc8dfedf`、`77f5eaae`、`e305c59d` 均归入对应功能依赖，不另算新功能。全部提交及父 SHA 保存在本轮 manifest 的 review inventory；“审查完成”只表示范围已分类，不表示 v2 全量集成或未采纳代码已通过本地集成测试。

| 来源／不可变合并 SHA | 当前决定、依赖及原生适用性 |
| --- | --- |
| [#595](https://github.com/Willxup/cpa-usage-keeper/pull/595)／[`44c2e39007f348bc8c3f5ab064765a14df41030d`](https://github.com/Willxup/cpa-usage-keeper/commit/44c2e39007f348bc8c3f5ab064765a14df41030d) | 待确认；整组费用持久化、分支定价、30 天重算和启动迁移改变费用历史语义。包含 benchmark／恢复相关提交，不直接移植；见待确认 2。 |
| [#605](https://github.com/Willxup/cpa-usage-keeper/pull/605)／[`0f1f47f4f447d8ca0b7562a778163de031da2bbe`](https://github.com/Willxup/cpa-usage-keeper/commit/0f1f47f4f447d8ca0b7562a778163de031da2bbe) | 延期；仅 .gitignore／.dockerignore 的 agent/worktree 排除。无运行时原生修复，当前候选在仓库外；Docker 未安装，未验证容器输入变更，待有 Docker 的独立批次。 |
| [#606](https://github.com/Willxup/cpa-usage-keeper/pull/606)／[`555ce3ef8e561f06494ec4f20b14a62b7afe1df7`](https://github.com/Willxup/cpa-usage-keeper/commit/555ce3ef8e561f06494ec4f20b14a62b7afe1df7) | 已部分适配；57253edb6acc2f4459c0338e7b171daa3790750b 的 Modal 焦点保护及三项通用测试；不引入不存在的 v2 PricingSettings 测试和计价界面。 |
| [#608](https://github.com/Willxup/cpa-usage-keeper/pull/608)／[`29d98cf86a7cf9c97a94a0204616e2a48fd83feb`](https://github.com/Willxup/cpa-usage-keeper/commit/29d98cf86a7cf9c97a94a0204616e2a48fd83feb) | 前端／服务移除意图已等效；本地早已无排名与 Key 用户页。排除 DROP 排名表／列及删除设置的迁移，不为代码精简删除历史数据，也不恢复 Key 登录。 |
| [#617](https://github.com/Willxup/cpa-usage-keeper/pull/617)／[`9246be3a1ae63ffa423c8558793db95d437a8a24`](https://github.com/Willxup/cpa-usage-keeper/commit/9246be3a1ae63ffa423c8558793db95d437a8a24) | 集成组去重；原生 Claude reset／partial headers／provider filters／retention、请求头／Alpine 更新已在历史选择性适配；Kimi 排除。v2 费用协调和测试依赖 #595，延期。 |
| [#620](https://github.com/Willxup/cpa-usage-keeper/pull/620)／[`3fa5e1e63e93eb825c83d8b5b1d1c815a00cf940`](https://github.com/Willxup/cpa-usage-keeper/commit/3fa5e1e63e93eb825c83d8b5b1d1c815a00cf940) | 生产者依赖延期；整组换成 /v8/management、分组配置文档写回及新版 api-call。实际 Proxy 为 /v0/management；继续旧接口及认证文件专用 priority，不跨项目更新 Proxy。 |
| [#621](https://github.com/Willxup/cpa-usage-keeper/pull/621)／[`ecad3013a335954193aa8facd350e207e2935ab5`](https://github.com/Willxup/cpa-usage-keeper/commit/ecad3013a335954193aa8facd350e207e2935ab5) | 依赖 #595，待确认；从仍保留的明细重建 Overview 会降低已删除明细对应的历史请求／Token 总量。禁止照搬；费用迁移方案必须保留旧聚合，见待确认 2。 |
| [#626](https://github.com/Willxup/cpa-usage-keeper/pull/626)／[`c9ed114d743e3fede1342162e5e7fd3eca51edaa`](https://github.com/Willxup/cpa-usage-keeper/commit/c9ed114d743e3fede1342162e5e7fd3eca51edaa) | 集成组去重；#619 已实现；#623 仍等真实字段。计价重算失效协调及迁移排序仅在 #595 下适用，未独立移植。 |
| [#627](https://github.com/Willxup/cpa-usage-keeper/pull/627)／[`d1d56533db5bf7c6347b140876a79dccccfb92a8`](https://github.com/Willxup/cpa-usage-keeper/commit/d1d56533db5bf7c6347b140876a79dccccfb92a8) | 待确认；周期摘要常显、明细默认折叠、历史周期选择和返回图表会改变现有 quota history 交互，见待确认 3。 |
| [#628](https://github.com/Willxup/cpa-usage-keeper/pull/628)／[`d579f06a0a64bf69fab438b43017d733a45aab44`](https://github.com/Willxup/cpa-usage-keeper/commit/d579f06a0a64bf69fab438b43017d733a45aab44) | 排除；通过 version.txt 下载次数估算上游 v2 安装量，不是 fork 的版本检查修复。不新增该外部请求、app_settings 或发布资产。 |
| [#629](https://github.com/Willxup/cpa-usage-keeper/pull/629)／[`e5b5f04b3764d1ffeb654c0df59ef45634f69249`](https://github.com/Willxup/cpa-usage-keeper/commit/e5b5f04b3764d1ffeb654c0df59ef45634f69249) | 与 #635 同组待确认；速度 sum/count 进入 rollup 并依赖 v2 初次迁移回填，不单独添加当前没有消费者的列，见待确认 4。 |
| [#634](https://github.com/Willxup/cpa-usage-keeper/pull/634)／[`7e7edddd8530b8d94f471e3984533251b0139e6e`](https://github.com/Willxup/cpa-usage-keeper/commit/7e7edddd8530b8d94f471e3984533251b0139e6e) | 集成组去重；#631/#632 已实现且 clipboard EOF 空行已清理；xAI #625 排除，#630 仍待确认。不引入 v2 计价和速度依赖。 |
| [#635](https://github.com/Willxup/cpa-usage-keeper/pull/635)／[`1b30badb334caac1b99044f9884b88ca35bcbab2`](https://github.com/Willxup/cpa-usage-keeper/commit/1b30badb334caac1b99044f9884b88ca35bcbab2) | 待确认；新增 Analysis／Realtime 模型速度图及选择器，按每请求速度算术平均，依赖 #629，见待确认 4。内含 c5dbed11 的 midnight DST 修复与 #656 一并评估。 |
| [#638](https://github.com/Willxup/cpa-usage-keeper/pull/638)／[`1cce5e5b46ecf780db42a2aef0655bf915d529f9`](https://github.com/Willxup/cpa-usage-keeper/commit/1cce5e5b46ecf780db42a2aef0655bf915d529f9) | 待确认；图表保留／刷新生命周期与 tooltip 交互成组变化，涉及 #627 选中周期、#635 速度图及 v2 页面 wiring；见待确认 5，不以复制文件覆盖现有布局。 |
| [#639](https://github.com/Willxup/cpa-usage-keeper/pull/639)／[`ea4924ae075694a9ac7f6506711adcc7466adc8e`](https://github.com/Willxup/cpa-usage-keeper/commit/ea4924ae075694a9ac7f6506711adcc7466adc8e) | 等效；本地 f21197c8 已保存 total-only 原生 Codex 总量并有入库回归；不引入上游为 v2 测试增加的 pricing catalog 或其他 provider。 |
| [#642](https://github.com/Willxup/cpa-usage-keeper/pull/642)／[`b8d0485de5808791c502996787140ebce8382d44`](https://github.com/Willxup/cpa-usage-keeper/commit/b8d0485de5808791c502996787140ebce8382d44) | 主要依赖 #595 的启动迁移／备份指纹／恢复工具／readiness，暂不引入整套启动架构；只给不适用的 v8 版本发警告也不采纳。Calendar 依赖与 #656 同组审查。 |
| [#645](https://github.com/Willxup/cpa-usage-keeper/pull/645)／[`a978a323b78c474694dd40fcffefbf30c075c77e`](https://github.com/Willxup/cpa-usage-keeper/commit/a978a323b78c474694dd40fcffefbf30c075c77e) | 集成组去重；#644 已在稳定版，但当前缺 execution_id；#646/#647 属本轮排除的发布策略。没有独立的原生可用增量。 |
| [#646](https://github.com/Willxup/cpa-usage-keeper/pull/646)／[`eb854ff34f9ae44310f551b8cbe0ad70c6474a12`](https://github.com/Willxup/cpa-usage-keeper/commit/eb854ff34f9ae44310f551b8cbe0ad70c6474a12) | 排除；移除 main push CI、改为 tag 全量 gate／可复用发布工作流会改变仓库发布策略；本次明确禁止策略变更，保留 fork 当前 main CI。 |
| [#647](https://github.com/Willxup/cpa-usage-keeper/pull/647)／[`5e369f948970d68685da00fd93f5fbc38e212252`](https://github.com/Willxup/cpa-usage-keeper/commit/5e369f948970d68685da00fd93f5fbc38e212252) | 排除；release/v1、latest、GitHub Latest 和 Homebrew 渠道路由属于上游发布策略，不创建分支/tag/Release，不改 tap 或镜像策略。 |
| [#649](https://github.com/Willxup/cpa-usage-keeper/pull/649)／[`4a9a14aba4ad9c99c19bdcd5ec6aea9d4c726b57`](https://github.com/Willxup/cpa-usage-keeper/commit/4a9a14aba4ad9c99c19bdcd5ec6aea9d4c726b57) | 部分等效／其余依赖延期；临时 pricing-restore 工具后来被删除，不能从早期 #642 恢复它；迁移回退提示只适用未采纳的 #595。quota cache 多余依赖已在本地 f21197c8 清理，无额外修改。 |
| [#650](https://github.com/Willxup/cpa-usage-keeper/pull/650)／[`fc8cbde97fe6baa35189e51115716865ce377ec2`](https://github.com/Willxup/cpa-usage-keeper/commit/fc8cbde97fe6baa35189e51115716865ce377ec2) | 发布集成标记；fc8cbde9 的树与第二父 4a9a14ab 一致，没有另一个独立实现。按上述各组选择，不全量合并 v2。 |
| [#653](https://github.com/Willxup/cpa-usage-keeper/pull/653)／[`b35a04436ffb7be42b64c32a083014aeb31a688b`](https://github.com/Willxup/cpa-usage-keeper/commit/b35a04436ffb7be42b64c32a083014aeb31a688b) | 排除；仅 OpenAI Compatible 组成员过滤，当前 registry 只有 Codex／Claude，相关 source 已移除；共享 Fetch 的非 OpenAI 行为不变，不恢复其他 provider。 |
| [#654](https://github.com/Willxup/cpa-usage-keeper/pull/654)／[`57d3af409d148eb4c018e117f2355bb1ad4498c0`](https://github.com/Willxup/cpa-usage-keeper/commit/57d3af409d148eb4c018e117f2355bb1ad4498c0) | 排除；Antigravity 专属 UA／并行 quota 和 subscription，无原生 caller，不改 Codex／Claude 指纹或查询。 |
| [#655](https://github.com/Willxup/cpa-usage-keeper/pull/655)／[`59354601562d2f466d5816eab1faf96eb505d993`](https://github.com/Willxup/cpa-usage-keeper/commit/59354601562d2f466d5816eab1faf96eb505d993) | 待确认；新增 OIDC 管理员登录、session schema、启动撤销和依赖，扩展现有管理员认证边界；见待确认 6。不恢复用户 Key 登录。 |
| [#656](https://github.com/Willxup/cpa-usage-keeper/pull/656)／[`e153b8d8e97d5edc2190b3288bb15ea7e1571a78`](https://github.com/Willxup/cpa-usage-keeper/commit/e153b8d8e97d5edc2190b3288bb15ea7e1571a78) | 待确认；原有时间桶只读修正、DST 日界／原始事件读取、缓存、标签和实际日长分母需明确范围；存在与已删明细的 rollup-only 数据及 v2 费用投影交织，见待确认 7。当前未修改固定 1440 分母，已知非上海 DST 边界问题尚未修复，不宣称等效。 |
| [#657](https://github.com/Willxup/cpa-usage-keeper/pull/657)／[`e21b436f520d8e22a3d53531019fdcb329152436`](https://github.com/Willxup/cpa-usage-keeper/commit/e21b436f520d8e22a3d53531019fdcb329152436) | 延期；Ubuntu runner 与 artifact Action 运行时切换依赖上游新发布工作流，当前未修改 CI／发布文件，不把上游通过的 CI 当作 fork 的验证。 |
| [#658](https://github.com/Willxup/cpa-usage-keeper/pull/658)／[`f298d0654c04b24ab26b17c827aaf4d2a81024ed`](https://github.com/Willxup/cpa-usage-keeper/commit/f298d0654c04b24ab26b17c827aaf4d2a81024ed) | 待确认；无兼容别名地重命名 /auth-files、API 字段和列表偏好会破坏现有消费者；中文文案无收益，见待确认 8。 |
| [#660](https://github.com/Willxup/cpa-usage-keeper/pull/660)／[`d1f6d5f2bddc30674a8ccc974ccf1cc69f9f6aa1`](https://github.com/Willxup/cpa-usage-keeper/commit/d1f6d5f2bddc30674a8ccc974ccf1cc69f9f6aa1) | 验证阻塞，延期；Go 1.26.10／Node 24.22.1 与 Docker builder／workflow 版本对齐需成组验证。当前无 Docker CLI，不能证明新容器构建可用；本轮不改 go.mod／Dockerfile／依赖／CI。应用检查用现有 Go 1.27.1／Node 26.8.1，不冒充目标构建验证。 |
| [#661](https://github.com/Willxup/cpa-usage-keeper/pull/661)／[`5604768d352a1a85dd5c087b4e85bd657a159a02`](https://github.com/Willxup/cpa-usage-keeper/commit/5604768d352a1a85dd5c087b4e85bd657a159a02) | 已适配；729764e5c63b5fd28e3f229626380488a58c9ba2 的本地日期截止比较，保留聚合水位门禁与事务；0／90／180、ns 边界和春秋跨 DST 回归。只修热表误归档，不改期限、不回迁旧记录、不重算汇总。 |
| [#664](https://github.com/Willxup/cpa-usage-keeper/pull/664)／[`5343b110b6dda8c4849de3758600446fd422b7a2`](https://github.com/Willxup/cpa-usage-keeper/commit/5343b110b6dda8c4849de3758600446fd422b7a2) | 未发布，延期；仅五个 CI／build 提交。工具链系列追踪、Docker build args／pull、OS runner 和 paths-filter v4 变更均未同步；等待稳定发布及适用的 Docker／workflow 验证，不改变本轮选定应用代码。 |

### 历史待决复核与实际 CodexProxy

- #630 继续未批准，详见待确认 1；用户沉默和周期间隔均不构成批准。
- #623（`08ec76afd994bd67d6493c5241a87167b807f692`，merge `41a9b90cfcee50ea5a43f6ad9cfb25ac4d5f1bd0`）已正式发布，但实际生产者仍无五个执行／节点字段；保持技术延期，而非已实现。
- #644（`a4696c80651e97cbcb2103080489a6fc70715fd8`、`6e4cef80c742333312640a62b41a73cb3c5dcf02`，merge `aa6056bf34667f28549691ef893d75a6fc0d47de`）现在进入 v2.0.0，**发布条件已解除，兼容性条件未解除**。最终 decoder 仍对所有消息强制 execution_id；当前无此字段，不能移植。历史恢复／重放也未运行。等实际生产者兼容后再提出保留旧 request_id 消息、去重与可回退迁移的具体方案，不能在 Keeper 任务修改 Proxy。
- #625、#640／#641 仍排除；后两项现已发布，但上游 Star History 不应覆盖 fork README。

另一项 Proxy 同步正在运行，只读观察其 `codex/user-provider-isolation` 的已提交 HEAD `5e986f5e432f96abd8abda69a2783cd5abc65b07`，它的台账有该任务自己的未提交改动，未读取为本轮授权、未触碰。producer blob `d91c8a2820cdea5a5212b1bb431d7effe25a8a18`、RESP blob `33e1dd64696fc4273f2f1fbf6e1446f8fe978dba`、management routing blob `d2397f42c3c7acee48f97d532ecf28f7d80c49ba` 与上轮相同；key PATCH blob 为 `b0b4fd4463bc3a17bf809b682595e1c599b00a5a`，仍无 priority。

当前管理组仍为 `/v0/management`，支持现有 auth-files／Codex／Claude key、api-call／reset-quota；AUTH／SUBSCRIBE usage／LPOP usage 未变。队列继续输出 request_id、原有 native tokens、response_headers，附加 accounting_version／token_breakdown 不改这些原生 flat 字段；SDK helper 仅附加 breakdown（Total 为 0 时才兜底 Total），本轮不切换 accounting 定义。没有 execution_id、trace_id、node_kind、is_fork、is_compaction、resolved_client_ip 或 stream，UI 缺 stream 仍显示未知。结论限于本地源码契约，不是线上端到端验证；未加载 .env、连接真实队列或使用真实管理写接口。

### 实际验证与候选身份

只改四个应用／测试文件，README／台账另计。#606 三项新测试中，未修复时两项焦点竞态失败；修复后 9 项 Modal 测试全部通过。#661 新回归未修复时 UTC 六种场景都出现 archived=3 而不是 1，cutoff 后 1ns／500ms 被提前归档；修复后通过。没有降低断言、恢复非原生功能或创建新 schema migration；已有归档门禁、原子移动和默认永久保留不变。

环境：darwin arm64，Go 1.27.1、Node 26.8.1、npm 11.19.0、Apple clang 21.0.0，Go 检查全部显式 CGO_ENABLED=1。隔离候选没有生产环境文件／数据；SQLite／队列／凭据／额度测试均使用临时 fixture／mock。

| 实际命令／检查 | 结果 |
| --- | --- |
| `npm --prefix web ci` | exit 0，239 包，未改 lockfile；npm 报告 11 项 audit（2 moderate／7 high／2 critical），原 lock 未变，上轮为 8 项。另运行 `npm --prefix web audit --json` exit 1 并留证；这是依赖告警，不是应用测试失败，也不代表已完成漏洞可利用性研判 |
| `NODE_OPTIONS=--no-experimental-webstorage npm --prefix web test -- --reporter=json --outputFile=<work>/keeper-vitest-20261025.json` | exit 0，157 文件、1240 passed、0 failed／pending；沿用 README 的 Node 26 happy-dom 兼容选项 |
| `npm --prefix web run lint`、`npm --prefix web run typecheck` | 各自 exit 0 |
| `npm --prefix web run build` | exit 0，212 modules，真实 HTML／CSS／JS／notice；既有 >500 kB chunk 提示。构建后原样恢复 tracked 零字节 .gitkeep，才开始 Go 检查 |
| `CGO_ENABLED=1 TZ=Asia/Shanghai go test -count=1 -json ./cmd/... ./internal/...` | exit 0，29 个有测试包，2436 tests/subtests passed、0 failed、无 cached；1 项既有 macOS skip：`internal/logging/test.TestDedicatedErrorMaintenanceClosesStaleDailyFile` 依赖 /proc/self/fd |
| `CGO_ENABLED=1 TZ=<zone> go test -count=1 ./internal/repository/test -run 'TestArchive\|TestCleanup'` | UTC／Asia/Shanghai／America/Los_Angeles 均 exit 0 |
| `CGO_ENABLED=1 TZ=UTC go test -race -count=1 -json ./internal/repository/test -run 'TestArchive\|TestCleanup'` | exit 0，35 tests/subtests passed、0 failed／skip，无 cached；定向 archive/repository 检查，不是全量 race |
| `CGO_ENABLED=1 go build -o <temporary>/cpa-usage-keeper ./cmd/server`，随后 `--version` | 均 exit 0、输出 dev；使用本轮真实 embed，二进制 SHA-256 `08e0870c036c886726917f33a3db21b34f990c4329cc7352ebfb6fda26aab4e0`，没有启动生产服务 |
| 本地真实 Modal 合成预览 | 桌面浅色 1280 px／窄屏深色 390×844，延迟初始化后焦点仍在 invalid input，关闭动画结束后恢复 opener，窄屏 clientWidth/scrollWidth 均 390，console errors=[]；不是完整页面视觉验收，无后端请求。临时 fixture 已移除，Vite 停止，浏览器 viewport 恢复 |

完整源码／tests／lock／构建配置 1069 文件，包含 intended untracked 输入的枚举规则与上轮相同。manifest `<work>/keeper-validation-manifest-20261025.json` 保留逐文件 SHA-256、220 提交 inventory 和实际 JSON 测试摘要。source digest `bcf94534e2e0b72bbb0f74603e31cc7e8621ab6da0bd3d034a6f62a62559ee5c`；候选身份为 `f21197c8` 加该完整选择性差异。README／台账不影响 runtime 输入，最终暂存时再次比较 manifest。5 个 embed 文件 digest `a993c977ef6a5e6e591ba115e108a47176bbb1e0840bafa70304b44bfe14f969`：

- `web/dist/.gitkeep`：`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
- `web/dist/index.html`：`65ab42eda8f49f96f63cbc5aa9f3b0bb7be4cddc50a4e6096c02aee2c07e0d27`
- `web/dist/assets/index-CkJYsm3B.js`：`a9d13387574340fdd7d00c8a322e3b2a39f93717cee5a0ae655fe6d5aef3624a`
- `web/dist/assets/index-OFLChr_n.css`：`834e7a12dca577c92a818e690432f1397ddb55180e128afa8cef8cd995fa2af6`
- `web/dist/third-party-notices.txt`：`3405117543d2887d5244eab04fe389d046a9b6c89a7cafa2863fed95931ad22f`

本轮 work 为 `/Users/vpromisex/Files/Codex/2026-10-10/new-chat/work`，报告与 JSON 留在聊天工作目录，不把日志／二进制／node_modules／生成 dist 提交。未覆盖 Linux／Windows runtime、上述 /proc 检查、Docker（CLI 不存在，相关候选不入批次）、完整页面交互、全量 race、线上端到端及所有未采纳的 v2 集成／数据迁移。依赖安全修复需要单独针对实际暴露面与跨大版本测试评估；本轮不执行 npm audit fix、修改凭据或消费重置券。

### 本轮发布边界

拟发布仅 #606／#661 四个源／测试文件和 README／台账。已核实现有 main push 只触发 CI；tag／手动发布工作流不触发。提交、正常 push、远端 readback 和 exact-SHA CI 结果留在当次运行报告，不制造台账自引用提交链。上一提交 f21197c8 的 Actions runs 仍为 0，不能称作 CI 通过。本轮不改 workflow／repo policy，不创建 tag／Release，不发布镜像，不部署，也不修改或推送 Proxy。清理前须确认已推送，若失败则保留唯一候选；原 checkout 只在再次确认干净未变后 fast-forward。

### 全部待确认事项（8 项，均未实施）

本 skill 明确要求：**“Before changing provider/login scope, the collection protocol, billing/statistical definitions, destructive migration or retention policy, automatic quota-reset-grant use, or substantial architecture/dependencies, discuss a concrete plan unless that plan is already authorized.”** 其 fork contract 同时要求：**“Preserve the existing layout, theme and responsive behavior unless the change itself is approved.”** 依据为 `/Users/vpromisex/.codex/skills/codexkeeper-upstream-sync/SKILL.md` 与 `references/fork-contracts.md`，并与本任务的统计／历史／界面保护要求一致。以下是具体决策，不把审核建议、上游 PR 或未回复当作授权。

1. **#630 用量分布布局（延续）。** 当前桌面 300／移动 250 px、有独立日期行和图例，汇总点击只看明细；上游缩为 240／200 px、去掉两行并让汇总点击切换曲线。建议保留现状，若用户批准则成组更换这些元素与 tooltip。收益：紧凑；风险：日期可见性与原点击语义变化。需要决定保留还是采用该完整方案；批准后做三语言、键盘／触摸、明暗／窄屏、维度／数据变化回归；无迁移，回退前端文件及重建。
2. **#595 费用持久化及 #621／#642／#649 迁移链。** 当前费用按读取时价格计算，改价可影响历史估算；上游改为入库记账、时段／上下文／周末分支与最近 30 天主动重算，而且 #621 从现存明细重建会减少已删明细对应的历史请求／Token。建议本轮保留现状；若要切换，采用独立的原生兼容方案：备份并核验指纹、固定旧价快照、分批回填保留明细的费用、保存原有历史总量（没有明细的旧费用标记未知，不能丢掉旧汇总）、原子更新新增费用字段／缓存、仅显式操作才重算，并为新旧汇总分开验证。**不批准删除旧聚合或采用上游减量重建。** 收益：历史费用稳定、定价灵活；风险：迁移耗时／磁盘、未知历史费用表达与费用定义变化。需要决定继续读时估价还是批准这套“不减历史总量”的兼容费用方案；批准后补 fixtures、故障恢复／备份失败／重启／回填／新旧总量及费用一致性测试，再交付可恢复升级产物。回退须经已验证备份恢复，不让旧二进制直接打开迁移库；服务器实施仍另行授权。
3. **#627 配额历史交互。** 当前保持已有周期／详情呈现；上游摘要常显、明细默认折叠、点击历史周期切换单图、返回图表平滑滚动避开 sticky tabs。建议保留现状，除非需要此整组改版。收益：长历史易浏览；风险：默认可见内容与选择／删除／滚动交互改变。需要决定是否采用整组；无 DB 改动，三语言、不同周期长度、删除、键盘／触摸、reduced motion、窄屏／明暗回归；回退相关组件即可。
4. **#629／#635 模型输出速度比较。** 当前没有该比较图；上游为每请求 output/latency 及 output/(latency-TTFT) 的速度存 sum/count，再算均值，新增两张图与模型选择。建议暂不扩展，若批准则只实现 Codex／Claude、增量保存速度聚合、可恢复回填保留明细且不重建／减少旧请求与 Token 总量，并保留现有页面布局至单独确认图表位置。收益：模型速度可比较；风险：均值不等于总 Token/总耗时、历史样本缺失、额外迁移与布局。需要决定是否增加这些指标／图表及接受其均值口径；批准后 old-schema／回填／重试、零或无效 TTFT、加权组、缓存／rollup 与 UI 选择回归。回退新 UI/API，库回退用备份。
5. **#638 图表生命周期与 tooltip。** 当前沿用现有刷新／点击／关闭方式；上游保留 chart instance/data、成功手动刷新才入场动画，tooltip 跟指针且去掉关闭／捕获，并绑定 #627／#635 的新图。建议先保留交互；可批准一个受限方案：只为现有原生图表保留成功数据、修正刷新动画触发，继续现有 tooltip 关闭与触摸／键盘语义，不带入未批准图表。收益：减少刷新闪烁；风险：旧数据展示误导、query scope／failed response／tooltip 生命周期竞态。需要决定是否批准该受限方案（而非整包改交互）；测试成功／失败／自动／手动刷新、查询切换、明暗窄屏与焦点，前端可回退，无 DB。
6. **#655 OIDC 管理员登录。** 当前只有本地管理员密码／会话；上游增加外部 OIDC、PKCE／state／nonce／Discovery、白名单和启动撤销，空 allowlist 会让该客户端所有已认证用户获得管理员权限。建议保持密码登录；可选方案为默认关闭、启用时必须配置精确 subject 白名单及可信 issuer／client，不允许空白名单，保留密码备用／现有会话生命周期，只增加管理员 OIDC，不恢复 Key 用户或嵌入认证。收益：集中身份管理；风险：新的管理员信任边界／依赖、错误配置越权和降级残留 session。需要决定是否加入受限 OIDC 及允许的 IdP／主体；批准后 mock IdP、拒绝／重放／CSRF／撤销／反代子路径／限流与登录视觉测试。降级前撤销 OIDC 会话，schema 回退需评估备份。
7. **#656 时间查询与统计分母。** 当前自然日均值按固定 1440 分钟（`usage.go`／`service/usage.go`），DST 区间的 timestamp 文本比较／缓存启动／异常午夜存在已知漏取或错位风险；上海常规日没有 DST。上游只读解释原桶键、按时区段读事件与缓存、修正日界和重复标签，并把日均 RPM/TPM 改用真实 1380／1500 等分钟数，牵涉已有统计口径，不能把整个 50 文件补丁当纯显示修复。建议批准在独立批次成组修正只读日界／原桶键、缓存与标签，并明确日均速率改按实际日长；不改存储归桶、不迁移／重算历史、不引入 v2 存储费用／ModelSpeed。收益：正确区分重复小时与短／长日；风险：DST 日速率和部分日粒度选择与旧结果不同、旧非标准桶兼容性。需要决定是否允许实际日长分母及这组原生时间修复；批准后 UTC／上海／纽约／柏林／Santiago／Lord Howe、23／25h、跳过午夜、亚秒边界、旧 rollup-only／归档／缓存生命周期和前端标签完整回归。无写库迁移，可退回旧读取实现；当前尚未验证这组在 v1 动态计价下的集成。
8. **#658 Credentials 命名与破坏性路由改名。** 当前保留 /auth-files 和现有 API 字段、书签／列表偏好；上游无 alias 地改为 /credentials，并丢弃旧偏好。建议不做破坏性重命名；可选仅原生方案为英文改名并兼容旧路由／响应字段／偏好迁移，中文现有名称不变。收益：OAuth/API key 总称准确；风险：外部 API、URL 和存储偏好兼容性及较大无功能重构。需要决定保留现状还是批准兼容重命名；批准后新旧 URL、API shape、偏好、stats/reset／quota／request-log 及三语言测试。无需删除历史实体，可回退显示与兼容层。

除上述 8 项外，无其他需当前用户选择的事项。#623／#644／#620 等待实际生产者契约，#605／#660 等待验证条件，不因解除“正式发布”条件而自动采纳；未发布 #664 与明确范围外排除不是待批准的需求。

## 2026-11-01：选择性适配 v2.0.2 原生修复，审查至 v2.0.3

### 身份、完整范围及历史待决

本轮开始时本地 main／origin/main／GitHub 均为干净的 `8753d2ee38c71dc1529d6f0f571d79c584152a9b`，无未提交 Keeper 适配或无关本地提交。重读最新 skill、三份 references、两份 README、当前检查点与全部待决；已提交台账与上轮 SHA 相同，历史基线未改写。各适用父目录及仓库没有 AGENTS.md。活跃任务列表只有当前 Keeper 同步，另一项活跃任务属于 Proxy，未向它发送消息或借用授权。common Git directory 的原子锁由长生命周期 PID 63870 持有，run `keeper-20261101-01a1242c`、host `Wills-Mac-Studio.local`、上海时间 05:03:50 开始；隔离 worktree／分支 `sync/keeper-20261101` 的 fork base 为 8753d2ee。没有 stash/reset/clean 用户 checkout。

完整范围 `5343b110b6dda8c4849de3758600446fd422b7a2..9db74cfc380d11064af80cf1626efd6e3d068217` 共 **109** 个新增提交，v2.0.2 包含 46 个、v2.0.3 再包含 63 个；最新 main 正好等于 v2.0.3，无未发布增量。上轮末尾 5 个 #664 提交本轮已正式发布，但发布策略／Docker 验证条件仍不满足，不沿用“未发布”作为延期理由。官方 release、tag Git object、独立 fetch ref 一致，两个 tag 均直接为 commit；没有用 target_commitish、改 origin 或覆盖 fork tags。

完整枚举了提交与父 SHA，审查 14 个 PR 的理由、差异、原生 caller、schema／价格／前端依赖与测试，merge 与子提交去重。额外 integration merges `85a104be`、`222ec4b3`、`0d76e421`、`2c2b9b56` 纳入其对应组；没有 revert。#684 是对 #669 新增提示的后续撤销，按最终树检查，不复制该中途 UI。109 个提交 inventory 随 manifest 留证；审查覆盖不表示已验证未采纳 v2 集成。

| 来源／不可变合并 SHA | 当前分类、原生行为与依赖 |
| --- | --- |
| [#668](https://github.com/Willxup/cpa-usage-keeper/pull/668)／[`61160e1ce1657a00b4cd9ae01db415e95e87c3f3`](https://github.com/Willxup/cpa-usage-keeper/commit/61160e1ce1657a00b4cd9ae01db415e95e87c3f3) | 排除可选界面改版；Top Models 工具条删除会改变现有摘要，ModelSpeed 的 hover／focus／selection 修复没有本地已启用组件。保持现有布局、摘要和图表，不恢复未批准 #635；上游 runtime settings mock 的测试清理不引入新页面。 |
| [#669](https://github.com/Willxup/cpa-usage-keeper/pull/669)／[`caaf19428a474a67cfd6192a797a6d0ca80e3f76`](https://github.com/Willxup/cpa-usage-keeper/commit/caaf19428a474a67cfd6192a797a6d0ca80e3f76) | 生产者依赖延期；effective_service_tier 是独立新事实，需 hot/archive 增列、decoder／DTO／API／CSV 及 UI 成组适配；当前 Proxy 不输出该字段，不能用 request/response tier 推断或回填。旧 raw tier 展示保留；待生产者明确提供后重新评估。 |
| [#670](https://github.com/Willxup/cpa-usage-keeper/pull/670)／[`d1687fc9cbe5d51a82e31a693f4eb27b13a90fa3`](https://github.com/Willxup/cpa-usage-keeper/commit/d1687fc9cbe5d51a82e31a693f4eb27b13a90fa3) | 部分适配；采用进程内版本＋短锁协调原生 auth-file priority/status、原生 key status、批量停用／删除和 metadata sync；远端成功但本地失败也失效旧快照，部分批量成功请求新同步，管理 API Key 独立来源仍可持久化。排除 v8 分组配置 ETag／If-Match／412 UI、API Key priority 及 OpenAI group 支持；v0 PATCH 协议不变。 |
| [#671](https://github.com/Willxup/cpa-usage-keeper/pull/671)／[`ee62c96a863f0743874f34d4f076512b00bd9a1f`](https://github.com/Willxup/cpa-usage-keeper/commit/ee62c96a863f0743874f34d4f076512b00bd9a1f) | 技术延期；完整 v2 breakdown 可服务未登记 executor，但当前支持的 Claude/Codex/Websocket/Auto 四个 native executors 已登记，上游明确不改变这些 handler。需实际原生未知 executor 的契约及规范化／费用一致性证据后再评估，不为插件／其他 provider 扩展 fallback，也不把现有 producer 含 breakdown 等同于必须改统计定义。 |
| [#672](https://github.com/Willxup/cpa-usage-keeper/pull/672)／[`c378b9b3794802a79b2cb5f6454abd402e26e9a6`](https://github.com/Willxup/cpa-usage-keeper/commit/c378b9b3794802a79b2cb5f6454abd402e26e9a6) | 已适配；共享精确五别名 parser，两个请求列表一致，其他路径不再随意去 /v1 前缀；原始 event.endpoint、完整 tooltip、exports 不变，POST/GET 标签和 stream 元数据语义保留。只改显示文本，不增加 provider／API。 |
| [#673](https://github.com/Willxup/cpa-usage-keeper/pull/673)／[`b4f91c66a698e5504051ae114531627e55befe21`](https://github.com/Willxup/cpa-usage-keeper/commit/b4f91c66a698e5504051ae114531627e55befe21) | 依赖延期；后端新查询直接读取 cost_usd/cost_available，并依赖 #656 的 usageEventTimeReads，当前动态计价版本无此存储合同；前端优化围绕未批准 #627 周期选择结构。保留本地读时费用及界面，等待 #595／#656 方案后按原生范围拆分验证，不能把上游性能百分比当本地收益。 |
| [#675](https://github.com/Willxup/cpa-usage-keeper/pull/675)／[`930a0d97138e6ead69b7a222e80b2bdb7296efe0`](https://github.com/Willxup/cpa-usage-keeper/commit/930a0d97138e6ead69b7a222e80b2bdb7296efe0) | 新增待确认 9；按活动页面／时间范围收窄 API Key 候选，保持已选值并惰性缓存刷新。当前全量 key 下拉可选项会改变，需用户选择，不影响三组独立修复。详情见下文，不恢复 Key 用户登录。 |
| [#676](https://github.com/Willxup/cpa-usage-keeper/pull/676)／[`46fe6278db0c47a5460e882fa1742f15ab49bb14`](https://github.com/Willxup/cpa-usage-keeper/commit/46fe6278db0c47a5460e882fa1742f15ab49bb14) | 已适配；解析显式 true/false，nil／非法保留旧值；主标记传播到缺失主窗口但不污染 Additional，并按标记自身时间防乱序或覆盖较新主动刷新。保留现有 Claude 部分窗口／刷新频率／计价；不引入 v2 成本重算 generation 参数。 |
| [#678](https://github.com/Willxup/cpa-usage-keeper/pull/678)／[`bb3f38f296fd8c86e836f0ae4af608c2477248e7`](https://github.com/Willxup/cpa-usage-keeper/commit/bb3f38f296fd8c86e836f0ae4af608c2477248e7) | 排除整包上游文档／截图／发布打包变更；本 fork README 已说明原生、v0、CGO 和数据限制，不能换成 v2／OIDC／第三方教程或发布 workflow。仅补充本轮真实选择性适配说明。 |
| [#679](https://github.com/Willxup/cpa-usage-keeper/pull/679)／[`5f8a59e1b3bbee02dd0d4c24e466d2752d08fce8`](https://github.com/Willxup/cpa-usage-keeper/commit/5f8a59e1b3bbee02dd0d4c24e466d2752d08fce8) | 排除第三方 Key 身份发现 worker／分组维度／settings cards 与 schema；实际队列没有 is_native_key／access_provider 前置字段，也不扩展原生凭据范围或改历史 api_group_key。附带 session 日期／字体调整是可选外观变化，不作为原生修复引入。 |
| [#680](https://github.com/Willxup/cpa-usage-keeper/pull/680)／[`4cfffa1c2feea8e672dd2b739777f69a593b250f`](https://github.com/Willxup/cpa-usage-keeper/commit/4cfffa1c2feea8e672dd2b739777f69a593b250f) | 依赖未采纳的 #675；只修 scopedApiKeys hook callback 的依赖警告，本地无该 hook，不产生独立修复；现有 lint 为 0。 |
| [#682](https://github.com/Willxup/cpa-usage-keeper/pull/682)／[`852a2d123b2031ae52fac829dffddbb116f7d5ce`](https://github.com/Willxup/cpa-usage-keeper/commit/852a2d123b2031ae52fac829dffddbb116f7d5ce) | 排除可选副标题删除与 13 组三语言清理；不影响正确性，保留现有 Overview／Analysis／设置卡说明。遵守当前界面保护，不凭上游减法改 fork UI。 |
| [#683](https://github.com/Willxup/cpa-usage-keeper/pull/683)／[`55bd58a09debd25d48702aba4c45f459b23bdd4d`](https://github.com/Willxup/cpa-usage-keeper/commit/55bd58a09debd25d48702aba4c45f459b23bdd4d) | 排除可选 badge／计数样式简化；现有计数并未算错，保留 title counts、filter counters 与当前视觉规则；不复制 mixins／credential sections 整包改版。 |
| [#684](https://github.com/Willxup/cpa-usage-keeper/pull/684)／[`9db74cfc380d11064af80cf1626efd6e3d068217`](https://github.com/Willxup/cpa-usage-keeper/commit/9db74cfc380d11064af80cf1626efd6e3d068217) | 当前无该 mismatch indicator，等效无需修改；上游移除的是 #669 后才新增的 effective-vs-response 提示，当前未采纳该新字段，也不人为新增再删除。保留已有 request／response tier 显示。 |

历史复核：

- #630／#595（含 #621/#642/#649）／#627／#629/#635／#638／#655／#656／#658 原 8 项均没有用户批准，具体方案、当前行为、收益／风险与迁移回退仍以 2026-10-25 的逐项说明为准；未把本轮执行请求当作批准。
- #623／#644／#620 仍受实际生产者限制。当前 Proxy 只读 HEAD `5e986f5e432f96abd8abda69a2783cd5abc65b07`，分支 `codex/user-provider-isolation`，其 docs/upstream-sync.md 为另一任务的未提交改动，完全保留。producer／RESP／management route blobs 仍为 `d91c8a2820cdea5a5212b1bb431d7effe25a8a18`／`33e1dd64696fc4273f2f1fbf6e1446f8fe978dba`／`d2397f42c3c7acee48f97d532ecf28f7d80c49ba`，key PATCH blob `b0b4fd4463bc3a17bf809b682595e1c599b00a5a`。继续 /v0/management、AUTH／SUBSCRIBE usage／LPOP usage、request_id 和 response_headers；没有 execution/trace/node/fork/compaction、effective_service_tier、stream、is_native_key 或 access_provider 新字段。原生 flat tokens 及附加 breakdown 不变。
- 因此 #676 可直接消费已有 response_headers 通道，#670 只适配原生已有管理写入；没有导入 v8 配置 PUT、ETag、API Key priority 或其他 provider，也没有声称线上端到端已核实。#625、Star History、排名／Key 登录等原范围外选择继续保留。
- #605／#657／#660／#664 仍无本地 Docker 构建证明，且不在本轮修改 CI／发布策略；Docker CLI 仍未安装。没有把新 tag 发布误当验证通过。
- #668／#682／#683 的纯外观减法在“保留现有界面”的明确要求下本轮排除，不新增未经请求的外观改版待办；若用户另行要求再评估。

### 实际适配、复现与边界

1. **#670 原生子集。** 来源 `e539772b9e4c2fa270b71ea2c246fa80082ebbfd`、`94bf740b466d9109f685c1596111d70496360ee7` 及 merge；不采用 `0d095a874271abe712ca95d49cab17a3890ec6b5` 的 v8 ETag 写回与相关 412 UI。新增共享零值可用 coordinator，metadata 在网络读取前取版本，仅在版本仍匹配时写凭据来源；单条成功写入（含本地落库失败）和每项批量成功先失效旧快照，再由原 runner 刷新。独立 management key 来源不跳过，网络／聚合／刷新通知均在短锁外。App 为所有实际入口注入同一个实例，旧构造器作为 nil coordinator 的兼容封装保留。认证文件专用 priority、原生 API key status、批量部分失败、207 删除语义不变，未修改 Proxy。
2. **#672。** 来源 `a926691573ecf46446f3d66d8420de930efa7a7d` 及 merge。两个请求列表共享五个精确 alias，其余路径含版本前缀保留；原始 event、tooltip、exports 不改。新组件回归在修改前 10 项失败，原有 4 处“去掉 /v1”的显示期望按新明确合同更新，同时保留完整 tooltip、未知 stream、焦点与卸载断言。没有套用上游 effective tier 测试或弱化导出校验。
3. **#676。** 来源 `e531ca000be52af2d7dbf39a2fd2f0d47052a6cc`、`2235c140c3a676c5c9f3136a92633b91d53a24fe`、`e7a3e170d2f25d13306a26128687d0a79e4b99ce` 及 merge。显式限流 true/false 保留独立观察时间；缺失／非法不覆盖，主额度缺失窗口只同步主标记，Additional 不变。旧 pending 标记不能借新窗口时间覆盖较新主动刷新。先复现 clear/set/primary-only 四场景及 pending 显式值丢失，再通过回归；保留 Claude 部分窗口和现有计价，未带入 v2 成本重算 generation API。

#670 的本地临时 SQLite＋httptest 回归在修复前稳定复现 10 个失败场景：auth-file status／priority 与 native key status 的成功／本地失败，以及批量 status／delete 的成功／部分成功。远端失败不失效有用快照，后续新一轮仍能正常对齐，独立管理 key 同步继续，失败路径没有通过改断言隐藏。新用例未恢复被移除的 provider priority／OpenAI group 测试。

### 本轮验证及最终输入身份

所有数据库、管理写请求和 queue／quota fixtures 均为临时 SQLite／mock；无生产 .env、历史用户库、真实 queue、真实凭据或 reset grant。环境 darwin arm64，Go 1.27.1／Node 26.8.1／npm 11.19.0／Apple clang 21.0.0；所有 Go 检查显式 CGO_ENABLED=1。真实前端先构建，恢复 tracked 零字节 .gitkeep 后才跑最终 Go 检查。

| 检查／实际命令 | 结果 |
| --- | --- |
| `npm --prefix web ci` | exit 0，239 包，原 lockfile；audit 23 项（14 moderate／7 high／2 critical），上轮 11 项，依赖未改。另 `npm --prefix web audit --json` exit 1 留证，不宣称漏洞安全或已修复 |
| `NODE_OPTIONS=--no-experimental-webstorage npm --prefix web test -- --reporter=json --outputFile=<work>/keeper-vitest-20261101.json` | exit 0，158 文件／1268 passed，0 failed／pending；沿用 README 的 Node 26 happy-dom 兼容选项 |
| `npm --prefix web run lint`、`npm --prefix web run typecheck` | 各自 exit 0，无警告 |
| `npm --prefix web run build` | exit 0，213 modules；真实候选 HTML／JS／CSS／notice；保留既有大 chunk 提示，不改构建设置 |
| `CGO_ENABLED=1 TZ=Asia/Shanghai go test -count=1 -json ./cmd/... ./internal/...` | exit 0，29 个有测试包、2462 tests/subtests passed，0 failed、无 cached；1 项既有 macOS skip：`internal/logging/test.TestDedicatedErrorMaintenanceClosesStaleDailyFile` 依赖 /proc/self/fd |
| 三组新增 metadata／bulk／header 定向回归 | UTC exit 0；上海／洛杉矶各 exit 0，分别 26 tests/subtests passed，0 failed／skip。使用 `-run 'TestMetadataSnapshotDoesNotOverwriteCompletedCredentialMutation\|TestCredentialBulkMutationsInvalidateOldMetadataOnAnySuccess\|TestUsageHeader(LimitReachedUpdatesCacheOnlyWhenPresent\|PendingPreservesLatestExplicitLimitReached\|PendingDoesNotOverrideNewerManualLimitReached)'` |
| `CGO_ENABLED=1 TZ=UTC go test -race -count=1 -json ./internal/service/... ./internal/quota/... -run 'Metadata\|Credential\|AuthFile\|UsageHeader\|Priority\|ClaudeHeader'` | exit 0，4 个有匹配测试的包、171 tests/subtests passed，0 failed／skip，无 cached；定向 race，不称为全项目 race |
| `CGO_ENABLED=1 go build -o <temporary>/cpa-usage-keeper ./cmd/server`，随后 `--version` | 均 exit 0，dev；二进制 SHA-256 `f4331f37065a1542d107bbe79ceeb7b3dfb6a587ac53cdce73effb994d87a065`，真实 Go embed，无生产服务启动 |
| 真实组件合成预览 | 两个列表、English／简中／繁中，桌面浅色和 390×844 深色；显示规范别名，focus tooltip 保留原始 backend-api 路径和 Unknown／未知，窄屏 clientWidth/scrollWidth 均 382，tooltip x=64.99..317.01 在视口内，console errors=[]。不是完整登录／业务 E2E，无真实后端请求，fixture／Vite／临时 browser 已清理且 viewport 恢复 |
| 最终 diff／暂存／输入一致性 | 发布前比较全部源码和 embed 的逐文件 digest；不包含临时 preview、dist、node_modules、测试日志、env／runtime 数据；README／台账不改变已验证 runtime |

本轮 manifest `<work>/keeper-validation-manifest-20261101.json` 包含 1074 个源码／测试／lock／构建输入（含 5 个 intended untracked 新文件）、5 个 embed 文件、109 提交 inventory 与实际测试摘要；候选为 base `8753d2ee` 加完整选择性 diff。source digest **`830cadd867c58db3c1496187d136725f3bace9a5048ed9f1ea111f9216636294`**，assets digest **`d8268d1989f557a4373377f95b35647a227020f4c015e0d8ebdfe99a5eccfceb`**。

- `web/dist/.gitkeep`：`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
- `web/dist/index.html`：`a455761322b752aabcb330530e795806174d2931c58c4e265b87771204beb172`
- `web/dist/assets/index-CCwzxyoR.js`：`afff9992ca72600c5561c87079beefd088a5245acb568f2bf362e0eb9c28a52a`
- `web/dist/assets/index-OFLChr_n.css`：`834e7a12dca577c92a818e690432f1397ddb55180e128afa8cef8cd995fa2af6`
- `web/dist/third-party-notices.txt`：`3405117543d2887d5244eab04fe389d046a9b6c89a7cafa2863fed95931ad22f`

work 为 `/Users/vpromisex/Files/Codex/2026-10-10/new-chat/work`。未覆盖 Linux／Windows runtime、上述 /proc 日志检查、Docker（仍无 CLI，容器相关批次延期）、全项目 race、生产数据库／账号／queue／reset grant、未采纳 v2 模块。没有通过升级依赖／排除测试／伪造 embed 绕过失败。

### 发布状态与全部待确认（9 项）

仅正常提交／推送三组修复、必要回归与 README／台账；保留 main push CI，既有 tag／manual 的发布触发器未改变。远端 SHA、exact-SHA Actions／check-runs／statuses 的实际结果及清理状态保存于当次运行报告，避免循环台账提交。没有 tag／Release／镜像发布、仓库策略变更、服务器访问或部署，不修改／推送 Proxy。原 checkout 只有再次确认干净且 HEAD 未变后才 fast-forward；推送失败保留唯一候选。

本 skill（`/Users/vpromisex/.codex/skills/codexkeeper-upstream-sync/SKILL.md`） 要求 **“Do not interpret silence, a previous explanation, or an item's presence in the ledger as approval.”** 并通过 fork-contracts 要求 **“Preserve the existing layout, theme and responsive behavior unless the change itself is approved.”** 这让以下依赖项保持待确认，但不阻塞三组独立修复。前 8 项的完整当前行为／方案／收益风险／测试和回退说明仍保留于上一节，本轮逐项状态：

1. **#630 分布图**：仍保留 300／250 px、独立日期／图例与原点击；建议保留现状。需决定是否采用 240／200 px、删行及汇总切曲线整组。收益紧凑，风险日期可见性与点击语义；前端三语言／触摸／键盘／窄屏回归，无 DB，可回退组件。
2. **#595/#621/#642/#649 费用历史**：仍读时估价，建议保留；若批准，采用旧价快照＋备份指纹＋可恢复回填且不减少历史总量、缺历史费用保留未知的兼容方案。收益费用稳定，风险迁移／费用口径；需决定是否切换，完整故障恢复／历史总量检查，回退用已验证备份。
3. **#627 quota 历史**：仍现有周期呈现，建议保留；需决定是否摘要常显、默认折叠并选周期切图。收益易浏览，风险操作变化；不改 DB，测试周期／删除／滚动／焦点，可回退前端。
4. **#629/#635 速度比较**：当前无新比较图，建议不扩展；需决定是否采用“每请求速度均值”及原生 sum/count 可恢复回填／图表。收益性能对比，风险均值口径、历史样本和布局；不减少旧 totals，迁移验证与备份回退。
5. **#638 图表刷新**：仍当前动画／tooltip；可批准受限方案，仅保留成功数据并修动画，不改现有 tooltip 手势。收益少闪烁，风险 query scope／旧数据竞态；需决定该受限范围，成功／失败／手动／自动测试，无 DB。
6. **#655 OIDC**：仍管理员密码／会话；建议保持，可批准默认关闭且启用时必填精确 subject 白名单、保留密码备用、绝不恢复 Key 登录的方案。收益集中登录，风险信任边界／依赖与残留会话；需决定是否加入及 IdP/主体，mock 安全回归，降级先撤销 OIDC session。
7. **#656 时间读取／日均分母**：仍固定 1440 分钟；原生只读桶／事件／缓存日界修复与实际日长分母方案未批。建议明确批准后独立适配，不写库重算；收益 DST 正确，风险 DST 日读数改变。#673 的存储费用投影与 UI 不自动随之获批，需另外保持动态计价／旧界面的兼容拆分与时区 fixtures。
8. **#658 命名／路由**：仍旧 /auth-files／API 字段／偏好；建议拒绝破坏性改名，可批准新旧兼容的英文 Credentials 方案。收益称谓准确，风险 API／URL／偏好兼容；需决定，测试新旧 URL/shape 与三语言，无历史删除，可退回显示层。
9. **新增 #675 按页面和时间限定 API Key 候选。** 来源 `930a0d97138e6ead69b7a222e80b2bdb7296efe0` 及 `6b01c00c698aede39e80940bf49a51e213fb0e96`、`e72a616f064f636bcb231926febb559dc0c4ad67`、`308a41b8fb2148c917eeceac111d56ad277b8914`、`854d750cd5a24e2aef018abd8e37c0a4f4d57b84`。当前下拉来自全量 active roster，不按当前使用窗口缩窄。建议默认保留；如批准，只在现有下拉框中为 Overview／Realtime／Analysis／Request Events 提供只读范围候选，保留已选 key 和无范围 roster 兼容，静默／切范围只失效缓存并在下次打开时加载，手动刷新立即更新；不引入第三方 discovery、Key 登录或额外统计分母改变。收益是长列表更相关；风险是没在该范围用过的 key 不再可见、时间边界／过期候选错误。需要决定是否改变现有候选语义；批准后验证四页面、已选但不在候选、冷／热缓存、手动／自动、UTC／上海／DST 范围和键盘／窄屏交互。无需 DB migration，可退回原 filter hook／API 调用。

除这 9 项外无其他当前待用户决定项。生产者／工具依赖延期与范围外排除不等于等待用户批准；仅解除正式发布条件不会自动引入它们。

## 2026-11-08：审查 main 的 38 个未发布提交，无新增正式版适配

### 仓库、串行执行与覆盖范围

本地 main、origin/main 与 GitHub main 开始时均为 `206182621355bc98575e6470c5bc84e7da2b5985`，工作区干净、ahead/behind 0/0，无未提交 Keeper 适配或无关本地提交。读取了最新 skill、三份 references、两份 README 与完整台账；适用父目录和仓库无 AGENTS.md。活跃任务列表只有当前 Keeper 同步，另有独立 Proxy 同步，不向其他任务发送消息、不跨项目写入。

开始时 common Git directory 内无同步锁。原子锁由长生命周期 PID `35188` 持有，host `Wills-Mac-Studio.local`，run `keeper-20261108-7e51ab39`，开始于 `2026-11-08 05:05:40 +08:00`，并以进程表核实 owner。隔离分支 `sync/keeper-20261108`、worktree `/Users/vpromisex/Files/Codex/2026-10-10/new-chat/work/keeper-sync-20261108` 从上述 fork base 创建；未 stash/reset/clean 用户 checkout，也未把上游 ancestry 合并进 fork。

Releases 列表、latest、tag Git ref API 与 fetch 重新确认最新正式版仍为 v2.0.3／`9db74cfc`。完整新增范围为 `9db74cfc380d11064af80cf1626efd6e3d068217..aa687529e392833e17ad0eb7ccfc4dec8eda0982`，**38 个提交、4 个主线 PR、72 个最终差异文件**，全部未发布；正式版新增为 0。按实际祖先可达性枚举，不能按分支上提交的作者日期排除较早创建但本次才合入 main 的修复。

已审查四组的 PR 理由、最终差异、共享原生调用与回归依赖。组内提交及其父 SHA 全量留证，去重后的新增贡献分别为 #686 的 4、#687 的 16、#690 的 14、#692 的 4。六个 merge 均检查 remerge-diff：五个无额外解决差异；`611196258745f631b5ee7d95a1572ee173d9666d` 只解决 `CredentialDetailTriggers.test.tsx` 的 mock 冲突，保留 API mock、移除旧 i18n mock，不是新的运行时修复。没有实际 revert；trial 字号提交中“可回退”的描述不是 revert。主线合并后的最终树是判断依据，不恢复中途实现或已移除功能。

### 四组新增决定及重评条件

| 来源／精确 merge SHA | 当前行为、分类、依赖及决定 |
| --- | --- |
| [#686](https://github.com/Willxup/cpa-usage-keeper/pull/686)／[`4f4d3dd857bb1fec834aa60b36c43a1e37c20a53`](https://github.com/Willxup/cpa-usage-keeper/commit/4f4d3dd857bb1fec834aa60b36c43a1e37c20a53) | **未发布；排除可选外观扩展。** 新增 WebGL 登录光带、CSS fallback、半透明表面、隐藏页暂停和 reduced-motion 等生命周期；当前静态登录及管理员密码功能不依赖它。没有独立认证修复。不引入动画／透明主题，也不恢复 OIDC、Key 用户登录或 CPAMC embed。只有用户另行要求此外观方案才重新评估，不新增未经请求的改版待办。 |
| [#687](https://github.com/Willxup/cpa-usage-keeper/pull/687)／[`7b2b670f631335c0e2223af769c5abc04d731d8a`](https://github.com/Willxup/cpa-usage-keeper/commit/7b2b670f631335c0e2223af769c5abc04d731d8a) | **部分行为可适用，但未发布；整套卡片改版排除，独立修复延期。** 当前是行式凭据、quota/health 切换及原字体／启停／操作位置；上游改为同时呈现的紧凑卡片并删除旧显示偏好，不能覆盖现有界面。`27be13bd064b814d29473f2c9ab2869961033d9f` 的原生详情入口消失时回到稳定筛选控件、编辑框卸载后微任务判断入口仍连接，可能适用于当前行式 UI；当前编辑框已有同步 fallback，但详情无该 fallback，不能称已等效。待正式发布后，在原行式 UI 中单独复现／适配并验证保存／筛选移走入口、取消、键盘／触摸及关闭动画，不依赖卡片改版。`1e38cba3ae0a2c835c3b3314388042394379ba4f` 的 `trackAnchor` 为 opt-in，只接入本地没有的紧凑 `QuotaResetCountdown`，不因此引入逐帧追踪；后续须有既有 caller 的布局移动问题证据才抽取。 |
| [#690](https://github.com/Willxup/cpa-usage-keeper/pull/690)／[`d0e98252965d7ff34a400847c8cbf2d46182d6ec`](https://github.com/Willxup/cpa-usage-keeper/commit/d0e98252965d7ff34a400847c8cbf2d46182d6ec) | **未发布；排除新增 Models 页、价格区迁移及其他品牌资源；共享 guard 技术延期。** `/api/v1/models` 是管理员保护、no-store 的目录代理；本地已有内部 `/v1/models` fetch，但没有对外目录页，不能将新增功能视为缺失的接口修复。模型家族前缀优先只用于该目录的展示分组，绝不移入原生用量身份判断。价格搜索基于未采纳的 v2 PricingSettings／重算 UI，当前读时计价无相同组件。`3d3d81527f34eec1d2d95836dde26b3d791bd47a` 同时给共享 `ProviderBrandIcon` 映射加 `Object.hasOwn`；本地原生映射也用普通对象索引，待正式发布后单独验证 constructor/toString 等未知值无图标、Codex/Claude 正常且不扩展映射，不因新图标资源整包排除而遗忘这个 guard。 |
| [#692](https://github.com/Willxup/cpa-usage-keeper/pull/692)／[`aa687529e392833e17ad0eb7ccfc4dec8eda0982`](https://github.com/Willxup/cpa-usage-keeper/commit/aa687529e392833e17ad0eb7ccfc4dec8eda0982) | **未发布；当前没有同一弹框路径，外观变化排除。** `2d9e0e8f5ea61a4778afdaee6ac344022fe33b2f` 保留保存中的 alias opener 节点，修复的是 #679 后的 `ApiKeyIdentityCard`／alias modal；本地仍是行内 Input＋Save，无此 modal/opener 被移除路径，不恢复第三方身份发现以套用修复。标题／Key tooltip 删除与 hover 浮起是可选视觉变化，现有复制／显隐和提示保留。若以后用户明确采用卡片编辑路径，须连同 pending success／failure、关闭防护和回焦回归重评。 |

上游 PR 内测试与浏览器结果只用于理解覆盖意图，不能作为本 fork 的通过证据。#687 独立焦点修复与 #690 共享 guard 只等待稳定发布和本地验证，不需要当前用户批准新布局；本轮没有实施未发布补丁。

### 历史待决与 CodexProxy 只读兼容性

- #670 原生协调器／#672 精确端点显示／#676 限流恢复已在 fork `20618262`，不重复移植。#619、#606、#661、total-only 原生处理等现有选择性修复仍保留。
- 只读 Proxy 分支 `codex/user-provider-isolation`、HEAD `5e986f5e432f96abd8abda69a2783cd5abc65b07`；其 `docs/upstream-sync.md` 为其他任务的未提交修改，完全保留。producer／RESP／management routes／key PATCH blobs 仍为 `d91c8a2820cdea5a5212b1bb431d7effe25a8a18`／`33e1dd64696fc4273f2f1fbf6e1446f8fe978dba`／`d2397f42c3c7acee48f97d532ecf28f7d80c49ba`／`b0b4fd4463bc3a17bf809b682595e1c599b00a5a`。
- 继续支持 `/v0/management`、原生 auth-files/fields、Codex/Claude API Key PATCH、AUTH／SUBSCRIBE usage／LPOP usage；key PATCH 仍没有 priority。队列有 request_id、response_headers、flat tokens 和 token_breakdown，没有 execution_id／trace_id／node_kind／is_fork／is_compaction、effective_service_tier、stream、is_native_key／access_provider。没有连接真实管理接口或队列，没有修改／推送 Proxy。
- #623／#644／#620／#669 生产者前提仍不满足；继续延期，不导入强制 execution_id、v8 配置写回或补造历史字段。#671 当前四个原生 executor 仍有明确 handler；未知 executor fallback 仍须实际原生契约与统计／计价证据，不扩展插件或其他 provider。
- #673 仍依赖 #595 的持久化费用、#656 时间读取及 #627 交互，不能直接套用 cost_usd/cost_available 查询或声称性能收益。#675 的候选语义尚未批准，#680 对不存在 hook 的依赖清理仍无本地修复。
- #605／#657／#660／#664 的 Docker／工作流验证条件未解除：`command -v docker` 仍未找到 CLI，不能把已正式发布当作容器验证。没有修改 CI、工具链、依赖或仓库策略。
- xAI／Kimi／Antigravity、排名、Key 用户登录、第三方身份发现、Star History、#668/#682/#683 的可选外观变化仍按既有范围排除，不是等待用户批准的新增需求。

### 文档候选验证、发布边界与未覆盖项

本轮只有台账变化，记录新的完整 38 提交范围、具体原生候选与重评条件，不是空提交或仅更新时间。依照 validation.md 对有意义的文档批次要求，执行 Git 差异／链接对象／范围一致性检查；**未重新执行应用测试、npm ci、前端构建、CGO SQLite 测试或服务端构建**，不将上轮通过数冒充本轮新结果。没有 UI 改动，因此没有新浏览器验收。

逐文件核对候选全部 **1074 个源码／测试／lock／构建配置输入**，与 `keeper-validation-manifest-20261101.json` 完全一致，source digest 仍为 `830cadd867c58db3c1496187d136725f3bace9a5048ed9f1ea111f9216636294`。旧 manifest 的 1268 前端／2462 后端／171 定向 race 通过及 macOS `/proc/self/fd` skip 是 2026-11-01 的历史结果；旧 embed/binary 已随上一轮清理，本轮不声称重建或产生新二进制。依赖审计没有重跑，上一轮 23 项是历史发现，未宣称已修复或当前数量不变。

证据目录仍为 `/Users/vpromisex/Files/Codex/2026-10-10/new-chat/work`：`keeper-review-api-20261108.json` 保存官方 release／tag／main／PR 元数据，`keeper-review-manifest-20261108.json` 保存完整提交父链、分组、remerge 结果与源码指纹；四份 upstream diff 及唯一 merge-resolution diff 一并保留。本轮未执行未采纳代码的集成／迁移／视觉测试，也没有 Linux／Windows／Docker／生产端到端覆盖。

发布只包含本文件，正常推送至已核实的 fork main；当前 main push 为 CI，binary/container 发布仍由 tag 或 manual 触发，无额外自动部署 hook。发布前检查远端与原 checkout 未变、显式暂存路径和完整发布历史。最终提交号、远端回读、exact-SHA CI、原 checkout 更新及清理结果保存在当次运行报告，避免自引用 SHA 的循环提交。推送失败时保留唯一候选；成功后仅清理自有 worktree、分支和 lock。**本轮没有服务器访问、部署、tag／Release／镜像发布或工作流 dispatch。**

### 全部待确认事项（9 项，延续、均未实施）

本轮没有新的人工批准。skill `/Users/vpromisex/.codex/skills/codexkeeper-upstream-sync/SKILL.md` 明确要求 **“Do not interpret silence, a previous explanation, or an item's presence in the ledger as approval.”** 对费用／统计／管理员认证及迁移方案适用其 **“Before changing provider/login scope, the collection protocol, billing/statistical definitions, destructive migration or retention policy, automatic quota-reset-grant use, or substantial architecture/dependencies, discuss a concrete plan unless that plan is already authorized.”**；对界面方案适用 fork-contracts 的 **“Preserve the existing layout, theme and responsive behavior unless the change itself is approved.”**。因此以下具体方案继续待确认；技术延期和明确范围排除不混入此列表。

1. **[#630](https://github.com/Willxup/cpa-usage-keeper/pull/630) 用量分布布局。** 当前 300/250 px、日期行／图例、汇总只看详情。建议保留；可选完整方案是 240/200 px、去两行、汇总切曲线及相应 tooltip。收益紧凑，风险日期可见性和点击语义。需决定保留或采用整组；无迁移，三语言／键盘／触摸／明暗窄屏回归，回退组件并重建。
2. **[#595](https://github.com/Willxup/cpa-usage-keeper/pull/595)／[#621](https://github.com/Willxup/cpa-usage-keeper/pull/621)／[#642](https://github.com/Willxup/cpa-usage-keeper/pull/642)／[#649](https://github.com/Willxup/cpa-usage-keeper/pull/649) 费用历史。** 当前读时估价，改价影响历史估算。建议保持；可选原生方案是验证备份指纹、旧价快照、可恢复费用回填、缺明细费用保留未知，绝不减少旧请求／Token 聚合，只有显式操作才重算。收益历史费用稳定／灵活规则；风险定义、时间与磁盘、迁移恢复。需决定是否采用该不减历史总量方案；验证备份失败／中断恢复／历史总量与费用一致性，回退用验证备份，不能让旧二进制直接开迁移库。
3. **[#627](https://github.com/Willxup/cpa-usage-keeper/pull/627) 配额历史交互。** 当前既有周期呈现。建议保持；可选摘要常显、明细默认折叠、选择周期切图、返回图表滚动。收益易浏览长历史，风险操作与默认内容改变。需决定整组是否采用；无 DB，周期／删除／焦点／触摸／reduced-motion／窄屏回归，回退前端。
4. **[#629](https://github.com/Willxup/cpa-usage-keeper/pull/629)／[#635](https://github.com/Willxup/cpa-usage-keeper/pull/635) 模型速度。** 当前无新增比较图。建议暂不扩展；可选原生每请求 output/latency 与 output/(latency-TTFT) 的 sum/count 均值、可恢复回填及经确认位置的图表，保留旧 totals。收益诊断性能，风险均值口径／缺样本／schema 与布局。需决定是否接受指标、均值及图表；旧 schema／回填重试／无效 TTFT／缓存与 UI 回归，库回退用备份。
5. **[#638](https://github.com/Willxup/cpa-usage-keeper/pull/638) 图表刷新。** 当前动画和 tooltip 手势不变。可批准有限方案：只保留成功数据并修刷新动画，不改现有关闭／触摸／键盘行为、不新增图表。收益少闪烁，风险过期数据和 query scope 竞态。需确认此有限范围；自动／手动、成功／失败、切查询／焦点及明暗窄屏测试，无 DB，可回退前端。
6. **[#655](https://github.com/Willxup/cpa-usage-keeper/pull/655) OIDC。** 当前管理员密码与会话。建议保持；可选默认关闭、启用必填可信 issuer/client 和精确 subject 白名单、保留密码备用，不恢复 Key 用户登录。收益集中认证，风险管理员信任、依赖及残留 session。需决定加入与否及 IdP/主体；mock IdP 拒绝／重放／CSRF／撤销／子路径／限流回归，降级先撤销 OIDC session，schema 回退评估备份。
7. **[#656](https://github.com/Willxup/cpa-usage-keeper/pull/656) 时间读取／日均分母。** 当前固定 1440 分钟，DST/异常午夜读取及缓存边界问题仍在。建议批准后独立适配只读原桶解释、事件／缓存日界和标签，并明确采用实际 1380/1500 等日长分母；不重算历史、不改存储、不带入 v2 费用。收益 DST 正确，风险日读数改变／旧桶兼容。需决定是否接受该范围和实际日长；UTC／上海／纽约／柏林／Santiago／Lord Howe、亚秒／归档／rollup-only／缓存测试，无写库迁移，可回退 reader。#673 不自动获批。
8. **[#658](https://github.com/Willxup/cpa-usage-keeper/pull/658) Credentials 命名／路由。** 当前 /auth-files、旧 API 字段和偏好。建议拒绝破坏性改名；可选英文新称谓并兼容旧 URL／字段／偏好，中文不变。收益称谓准确，风险消费者兼容和重构。需选择保留或兼容方案；新旧 URL/API shape、偏好、quota/reset/stats/request-log 与三语言回归，不删历史，可退显示层。
9. **[#675](https://github.com/Willxup/cpa-usage-keeper/pull/675) 按页面／时间筛 API Key 候选。** 当前全量 active roster。建议保持；可选现有四页下拉仅列窗口已用 key，保留不在候选的已选 key 与无范围 roster，自动／切范围只失效缓存、下次打开加载，手动立即更新。收益列表相关，风险未用 Key 隐藏／边界／缓存过期。需决定是否改候选语义；无迁移，四页、冷热缓存、已选缺失、手动自动、时区/DST、键盘窄屏回归，可退原 hook/API；不增加第三方 discovery 或 Key 登录。

以上九项均不因本轮定时执行或无回复而获批；没有新增待确认事项。#687 焦点与 #690 guard 等未发布原生候选由后续稳定版检查自动复核，不需要用户决定是否破坏现有布局。
