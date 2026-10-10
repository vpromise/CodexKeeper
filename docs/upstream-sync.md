# CodexKeeper 上游同步台账

本项目是 `vpromise/CodexKeeper`，上游为 `Willxup/cpa-usage-keeper`。采用选择性适配，仅保留原生 Codex／Claude；审查边界不等于完整合并或全部实现。时间均为 `Asia/Shanghai`。

## 当前检查点

| 项目 | 已核实状态 |
| --- | --- |
| 本次检查时间 | 2026-10-10 13:03:47 +08:00 |
| 发布目标 | `https://github.com/vpromise/CodexKeeper` 的 `main`，origin fetch/push 均为此地址 |
| 本次隔离候选基线 | [`012609760c426382cec70015dabc043b4718779b`](https://github.com/vpromise/CodexKeeper/commit/012609760c426382cec70015dabc043b4718779b)，与开始时远端、本地 main 相同 |
| 最新正式发布／已选择性集成边界 | [v1.15.10](https://github.com/Willxup/cpa-usage-keeper/releases/tag/v1.15.10)，`b3592bff1e342c246d1324bcd99ba9b77232f315` |
| 本次完整审查范围 | [`b3592bff1e342c246d1324bcd99ba9b77232f315..778b5f8e7196f3fb1eae8f11b63f51de33ea80ed`](https://github.com/Willxup/cpa-usage-keeper/compare/b3592bff1e342c246d1324bcd99ba9b77232f315...778b5f8e7196f3fb1eae8f11b63f51de33ea80ed)：7 个提交，3 项合并后的逻辑变化；全部未发布 |
| 主分支已审查检查点 | `778b5f8e7196f3fb1eae8f11b63f51de33ea80ed`，延期／排除项保留在下表，不代表已实现 |
| 本次实际改动 | 首次建立本台账；没有移植未发布的应用代码，没有升级依赖或改变统计、采集、迁移、界面行为 |
| 待确认事项 | 无待确认事项；等待正式发布和生产者字段支持属于重新评估条件，不是用户批准请求 |

通过 GitHub Releases API 核实最新 release 非 draft、非 prerelease，发布于 2026-10-07 15:46:47 +08:00。标签对象为 commit，直接解析到上述不可变 SHA；没有使用 `target_commitish` 代替版本 SHA。不存在基线之后的新正式版或中间正式版。上游 main 与 release 分别 fetch 到 `refs/upstream-sync/cpa/main` 和 `refs/upstream-sync/cpa/releases/v1.15.10`，未覆盖 fork tags 或修改 origin。release 是当前 main 的祖先。

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

## CodexProxy 兼容性核实

只读检查 `/Users/vpromisex/Files/Code/RPCodex/CodexProxy`，HEAD 为 `2663a7456ca284d9d05daf8963dd66bf2cd8bd42`，工作区干净。本次没有修改、fetch 或推送 CodexProxy。

- `internal/api/redis_queue_protocol.go` 支持 `AUTH`、`SUBSCRIBE usage`、`LPOP usage`，与 Keeper `internal/cpa/endpoints.go` 一致；当前单采集器、先订阅后补收、失败重试既有批次等行为未变。
- `internal/redisqueue/plugin.go` 输出 native provider／executor、token breakdown、请求头及请求 ID；Claude `response_headers` 通道和现有 accounting 输入继续可用，但上述五个执行／节点字段和 `resolved_client_ip` 未包含在队列 JSON 中。
- `internal/api/server_management.go` 保留 `api-call`、`reset-quota`、Codex／Claude key 路由、`auth-files/fields`。`config_lists.go` 的两家 API Key PATCH 仍无 `priority`，现有认证文件独享优先级编辑限制仍正确。
- 结论为当前本地源码契约兼容；没有把源码核对称作服务器端到端验证，没有读取线上配置或调用真实管理写接口。

## 验证与输入身份

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

## 发布与运行状态

- 仓库锁为 resolved common Git directory 下 `codexkeeper-upstream-sync.lock`，由本次长生命周期进程持有，记录 host、PID、run ID、上海时区时间与隔离 worktree。开始时没有其他 Keeper 同步活跃聊天或既有锁；用户 checkout 无未提交改动、无远端以外本地提交。
- 本次使用独立分支 `sync/keeper-20261010` 和 disposable worktree；未 stash、reset、clean 用户 checkout，没有导入上游发布历史。
- 当前 workflow 定义中，main push 仅触发 CI；二进制／镜像发布由版本 tags 或手动 dispatch 触发。本次只正常推送台账文档，不创建 tag／Release，不 dispatch、不改 workflow 或仓库策略。
- 开始时 `01260976` 的 GitHub Actions run 查询为 0；不能把“CI workflow active”视作已通过。此次文档提交后的远端 SHA、check-runs、status 与 Actions 查询结果在当次运行报告记录，避免为自引用 SHA 创建循环台账提交。
- 本次没有服务器访问或部署。推送台账不会更新运行中的 Keeper 服务。

## 后续运行规则与待确认事项

后续检查新正式版、上游 main 和本表延期项；即使主分支 checkpoint 已推进，也必须重新处理 #619／#623，直到实现或形成新的有依据决定。没有新值得采纳的内容时，不创建空提交或仅更新时间的提交。

**无待确认事项。** #619 等待正式发布；#623 等待正式发布和实际生产者支持；xAI／Kimi 等为当前明确范围外排除。任何未来改变 provider／登录范围、采集协议、统计／计费定义、破坏性迁移／retention 或自动消费重置权益的具体方案，另行按照 skill 的讨论边界处理。
