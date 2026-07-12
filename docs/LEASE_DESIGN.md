# 代理 IP 租约 / 会话配额 / 冷却设计规格

状态：设计 only（对应 TODO #13；实现见 #14/#15/#16）  
范围：网关选路层（`selector` + `affinity` + 配置），不引入租约中心服务  
非目标：跨实例共享状态、计费、用户级配额、持久化租约表

---

## 1. 问题与现状

### 1.1 现状（代码事实）

| 能力 | 现状 | 位置 |
|------|------|------|
| Session 黏连 | 有：`session_id → {proxy_id, address, region, last_active}` | `affinity.Store` |
| 反向占用索引 | 无：无法 O(1) 回答「某 proxy 上挂了几个 session」 | — |
| 每节点 session 上限 | 无 | — |
| 节点冷却 CD | 无 | — |
| 首次绑定选点 | 地域内延迟 top-K 再按 session 哈希 | `selector.pickForSession` |
| 无 session 选点 | 地域内最低延迟（同延迟随机） | `selector.Pick` |
| 绑定 TTL | `SessionTTLMinutes`（默认 10），GC 每 1 分钟 | `affinity` + `main` |
| 失败 | 无可用节点 → `ErrNoNode`（可带 region） | `selector` |

### 1.2 痛点

1. 多个不同 `session-*` 可哈希/延迟策略撞到同一出口 IP，自动化场景（批量注册等）出口撞车。  
2. 某节点刚被 session A 占用后，session B 仍可立刻选中同一节点。  
3. 当前是 **best-effort 粘性网关**，不是 **租约中心**：无容量与冷却契约。

### 1.3 设计目标

在保持「地域 > 黏连 > 质量/分散」前提下，增加：

- **会话配额**：每 proxy 同时服务的活跃 session 数 ≤ `max_sessions_per_proxy`
- **节点冷却**：某 proxy 被 **新 session 首次绑定** 后，在 `proxy_cooldown_minutes` 内不对 **其他新 session** 开放
- 可测试、可配置、与现有 `Resolve`/`Pick`/affinity GC 语义可组合

---

## 2. 术语

| 术语 | 定义 |
|------|------|
| **Session** | DSL 中 `session-<id>` 的 id；空 session 表示无黏连请求 |
| **Binding / 黏连** | affinity 中 session→proxy 的映射；TTL 由 `SessionTTLMinutes` 控制 |
| **占用 (occupancy)** | 某 proxy 上当前未过期 binding 的数量 |
| **租约 (lease)** | 逻辑概念：binding 存续期间对该 proxy 的一次占用；**不**单独建租约表 |
| **冷却 (cooldown)** | proxy 在时间窗内禁止被 **新 session 首次绑定** 的状态 |
| **粘性命中** | 已有 binding 且节点可用、region 匹配、不在 excludes 中 → 复用 |
| **首次绑定** | 无有效 binding，需 `pickForSession` 并 `SetProxy` |
| **换绑 (rebind)** | 原 binding 失效（节点不可用 / region 不匹配 / exclude）后的首次绑定 |

---

## 3. 配置项

### 3.1 新增配置

| 字段 | 环境变量 | 默认 | 约束 | 热更新 |
|------|----------|------|------|--------|
| `MaxSessionsPerProxy` | `MAX_SESSIONS_PER_PROXY` | `0`（不限制） | ≥ 0；`0`=不限制；`≥1`=上限；负值拒绝保存 | 与 SessionTTL 同路径即时生效 |
| `ProxyCooldownMinutes` | `PROXY_COOLDOWN_MINUTES` | `0` | ≥ 0；`0`=关闭冷却 | 即时生效（仅影响后续判定） |

JSON 字段建议：

```json
{
  "max_sessions_per_proxy": 0,
  "proxy_cooldown_minutes": 0
}
```

自动化推荐（非发布默认）：

```json
{
  "max_sessions_per_proxy": 1,
  "proxy_cooldown_minutes": 5
}
```

### 3.2 与既有配置关系

| 配置 | 关系 |
|------|------|
| `SessionTTLMinutes` | 控制 binding 过期 → **释放占用**；与冷却无关 |
| `DefaultRegion` / 地域过滤 | 选路候选集不变；配额/冷却在候选集上再过滤 |
| `MaxRetry` | 入口重试 excludes 仍生效；配额满/冷却节点视为「不可选」，可进 excludes 或过滤层 |

### 3.3 默认值选择理由

- `max_sessions_per_proxy=0`（发布默认）：**兼容优先**，不改变现网多 session 可撞同一节点的行为；自动化场景显式设 `1`。  
- 推荐值 `1`：一 session 一出口，满足批量注册隔离。  
- `proxy_cooldown_minutes=0`：默认不冷却；需要「分配后静默期」时再打开。

---

## 4. 数据模型

### 4.1 Affinity 扩展（内存，进程内）

保持 `session_id → Binding` 为主存储；增加 **反向索引** 与 **冷却表**（可同包 `affinity` 或 selector 旁路结构，实现阶段二选一并单点拥有）。

**Binding（既有 + 可选字段）**

```text
Binding {
  ProxyID     int64
  NodeAddress string
  Region      string
  LastActive  time.Time
  // 可选：BoundAt time.Time  // 首次绑定时刻；冷却可用 LastActive 或单独 BoundAt
}
```

**反向占用索引**

```text
proxy_id → set(session_id)   // 仅含未过期 binding
```

不变量：

1. `SetProxy` / `Get`（刷新）/ `Remove` / GC 删除时，正反向索引同时更新。  
2. `CountByProxy(proxyID) == len(reverse[proxyID])`。  
3. 过期 binding 不计入占用（与 `List`/`Count` 一致：以 TTL 判定）。

**冷却表**

```text
proxy_id → cooldown_until time.Time
```

- 写入时机：见 §6。  
- 读：`now < cooldown_until` 则节点处于冷却。  
- GC：扫描时删除 `cooldown_until <= now` 的条目（可选；读时惰性判断即可）。

### 4.2 不落库

- 占用与冷却 **默认纯内存**（与现有 affinity 一致；重启清空可接受）。  
- 不在 SQLite `proxies` 表持久化 `cooldown_until`（避免与健康检查/订阅刷新状态纠缠）。  
- 若未来多实例，需独立设计共享层；**本规格单进程**。

### 4.3 时钟

- 与 affinity 一致：可注入 `now func() time.Time`，便于单测。

---

## 5. 与 Affinity 的关系

```text
                    ┌─────────────────────────────────────┐
                    │           affinity.Store            │
  session 请求 ──►  │  正向: session → Binding            │
                    │  反向: proxy_id → {session...}      │
                    │  冷却: proxy_id → cooldown_until    │
                    └──────────────┬──────────────────────┘
                                   │
                    selector.Resolve / Pick
                                   │
              过滤: 可用 ∧ 未 exclude ∧ 占用未满 ∧ (新绑定时未冷却)
```

| 规则 | 说明 |
|------|------|
| 租约 ≡ 活跃 binding | 不另建 lease 对象；占用 = 反向索引大小 |
| TTL 释放占用 | GC/`Get` 发现过期 → 删 binding → 反向索引减一 |
| 粘性命中 **不** 受冷却阻挡 | 已绑定 session 在冷却期内仍可继续用该节点 |
| 粘性命中 **不** 增加占用 | 已在反向索引中则计数不变；`Get` 仅刷新 LastActive |
| 换绑 | 先 `Remove` 旧 binding（释放旧 proxy 占用），再选新节点并 `SetProxy` |
| 地域优先级 | 不变：region 不匹配仍强制换绑；配额/冷却只在目标 region 候选上生效 |
| 空 session | 不写 affinity；见 §6.3 对配额/冷却的适用规则 |

---

## 6. Resolve / Pick 变更点

### 6.1 决策优先级（有 session）

```text
1. 地域约束（候选仅限 requested/default region 语义，与现网一致）
2. 粘性命中（binding 有效 ∧ 节点可用 ∧ region 匹配 ∧ 非 exclude）
     → 直接返回；不检查冷却；不改变占用计数（已占用）
3. 需要首次绑定 / 换绑：
     a. 取 region 内可用节点（现有 available + excludes）
     b. 过滤：occupancy(proxy) < max_sessions_per_proxy
     c. 过滤：not in_cooldown(proxy)   // 仅对新绑定
     d. 在剩余集合上执行现有 pickForSession（top-K + hash）或等价策略
     e. SetProxy + reverse 占用 +1
     f. 若 proxy_cooldown_minutes > 0：设置 cooldown_until = now + CD
        （策略见 §6.2：默认「每次新 session 绑定刷新 CD」）
4. 过滤后为空 → ErrNoNode（可扩展错误信息，见 §7）
```

### 6.2 冷却触发策略（定稿）

**触发**：任意 **新 session 首次成功绑定** 到 proxy 时（含换绑到新 proxy）。  
**不触发**：

- 同一 session 粘性命中  
- 空 session 的 `Pick`（见 §6.3）  
- binding 仅刷新 LastActive  

**时长**：`proxy_cooldown_minutes`；`0` 表示永不写入冷却表。  
**刷新语义**：再次有 **其他** session 成功新绑定同一 proxy 时，`cooldown_until` 取 `max(原 until, now+CD)`（实现可简化为直接 `now+CD` 覆盖，因覆盖 ≥ 原值当间隔 < CD）。  
**与 max_sessions 关系**：冷却与配额独立；配额允许多 session 时，冷却仍限制「新绑定」——若 `max_sessions_per_proxy>1` 且节点在冷却中，**其他 session 仍不可新绑到该节点**，直到冷却结束（粘性 session 不受影响）。  
说明：当 `max_sessions_per_proxy>1` 且需要「配额内允许多 session 同时新绑」时，冷却会与该意图冲突；本规格明确 **冷却优先于「同节点多新 session」**。若产品要「仅 N=1 时冷却有意义」，实现可将冷却在 `max_sessions_per_proxy>1` 时自动忽略，但必须在配置文档写明；**默认实现：冷却始终过滤新绑定**。

### 6.3 无 session 的 Pick

| 选项 | 行为 | 选用 |
|------|------|------|
| A | Pick 完全忽略占用与冷却 | **默认** |
| B | Pick 也过滤占用满与冷却 | 备选 |

**定稿 A**：无 session 请求不建立租约，不应消耗「自动化隔离」资源；批量无 session 流量仍走最低延迟。  
若运维需要全局打散，应强制客户端带 session，或后续另开「全局并发连接上限」需求（非本规格）。

### 6.4 代码触点（实现清单，非本任务改代码）

| 触点 | 变更 |
|------|------|
| `config.Config` / `savedConfig` / env | 两字段 Load/Save/校验 |
| `affinity.Store` | 反向索引；`CountByProxy`；`SetProxy`/`Remove`/GC 同步；冷却 API |
| `selector.Resolve` | 新绑定路径过滤占用+冷却；绑定成功后登记冷却 |
| `selector.Pick` | 默认不变（策略 A） |
| `selector.pickForSession` / 候选构建 | 输入改为「已过滤列表」或增加 filter 钩子 |
| `proxy/server.go` / `socks5_server.go` | 无逻辑变更若仍只调 `Resolve`；错误映射见 §7 |
| WebUI config | 设置项（可随 #16 或本实现一并）；只读展示可选 |
| 测试 | 见 §8 |

### 6.5 伪代码

```text
func Resolve(store, sessions, route, excludes):
  if route.Session == "":
    return Pick(store, route.Region, excludes)  // 忽略配额/冷却

  proxy, rebindRegion := resolveBoundProxy(...)  // 现有粘性逻辑
  if proxy != nil:
    return proxy, nil

  candidates := available in rebindRegion minus excludes
  candidates = filter occupancy < MaxSessionsPerProxy
  if ProxyCooldownMinutes > 0:
    candidates = filter not cooldown(proxy)
  if candidates empty:
    return ErrNoNode (with reason)

  picked := pickForSessionFrom(candidates, route.Session)
  sessions.SetProxy(...)
  if ProxyCooldownMinutes > 0:
    sessions.SetCooldown(picked.ID, now+CD)
  return picked
```

---

## 7. 失败语义

### 7.1 错误类型

| 条件 | 错误 | 对外表现（HTTP / SOCKS5） |
|------|------|---------------------------|
| 地域无可用节点（现网） | `ErrNoNode` / `ErrNoNode for region: X` | HTTP 503；SOCKS5 通用失败（与现网一致） |
| 有节点但全部占用满 | **同一** `ErrNoNode` 族；消息可区分 | 同上；**不**返回 429（避免与认证/限流混淆） |
| 有节点但全部在冷却 | 同上 | 同上 |
| 占用满 + 冷却叠加 | 同上 | 同上 |
| 粘性节点 exclude 后重绑失败 | 同上 | 同上 |

### 7.2 错误消息（实现建议）

保持 `errors.Is(err, ErrNoNode)` 稳定，便于入口判断：

```text
no available node
no available node for region: us
no available node for region: us (capacity)
no available node for region: us (cooldown)
```

- 后缀仅用于日志与排障；客户端契约仍是「无节点」。  
- **禁止**因配额/冷却返回认证失败或静默落到其他 region。

### 7.3 不失败的情形

| 情形 | 行为 |
|------|------|
| 粘性命中且节点在冷却中 | **成功**，继续该节点 |
| 粘性命中且占用已达上限（本 session 已计入） | **成功** |
| 换绑时旧节点占用释放后选新节点 | 旧 proxy 占用 -1 后参与其他 session 竞争 |
| `max_sessions_per_proxy` 热更新变小 | 已有超额 binding **不主动踢**；仅阻止新绑定直到自然过期 |
| `proxy_cooldown_minutes` 热更新变 0 | 不再写入新冷却；已有 until 可立即视为无效或保留至到期（实现选「读时若配置为 0 则忽略冷却」更简单） |

### 7.4 重试交互

入口 `tried` excludes：

- 拨号失败 exclude 某 proxy 后，`Resolve` 再选：粘性被 exclude → 走新绑定路径 → 受配额/冷却约束。  
- 不得在 exclude 列表耗尽后跨 region 回退。

---

## 8. 测试矩阵

### 8.1 单元（affinity）

| ID | 场景 | 期望 |
|----|------|------|
| A1 | SetProxy 两 session 同 proxy | CountByProxy=2 |
| A2 | Remove 一 session | CountByProxy=1 |
| A3 | TTL 过期后 Get/GC | 占用降为 0 |
| A4 | 并发 Set/Remove/Get race | `-race` PASS |
| A5 | SetCooldown + 注入时钟未到期 | InCooldown=true |
| A6 | 时钟越过 until | InCooldown=false |
| A7 | 配置 CD=0 | SetCooldown 为 no-op 或读侧忽略 |

### 8.2 单元（selector）

| ID | 场景 | 期望 |
|----|------|------|
| S1 | N=1，session-a 绑 proxy1 后 session-b 同 region | b **不得** 得 proxy1（有其他节点时） |
| S2 | N=1，仅一节点且已被占 | session-b → ErrNoNode |
| S3 | N=2，两 session 可同节点 | 允许同 proxy |
| S4 | 粘性：session-a 再请求 | 仍 proxy1；冷却/满载不影响 |
| S5 | CD>0：a 绑定后 b 在 CD 内 | b 不选 a 的节点（有替代时） |
| S6 | CD 到期后 b | 可选回该节点（占用也允许时） |
| S7 | 粘性 + 节点 exclude 换绑 | 释放旧占用；新节点遵守 N 与 CD |
| S8 | 无 session Pick | 忽略 N 与 CD（策略 A） |
| S9 | region 隔离 | us 占用不影响 jp 选路 |
| S10 | top-K 哈希仍稳定 | 在过滤后候选集上哈希；同输入同结果 |
| S11 | 全部冷却且 N 未满 | ErrNoNode（cooldown 类消息可选） |
| S12 | max 热更新 N:2→1 | 已有 2 binding 保留；第三 session 失败 |

### 8.3 集成 / 入口

| ID | 场景 | 期望 |
|----|------|------|
| I1 | HTTP 带两 session 撞容量 | 第二请求 503 + no available node |
| I2 | SOCKS5 同上 | 失败码与现网无节点一致 |
| I3 | 同 session 连续请求 | 出口稳定（粘性） |
| I4 | GC 释放后第二 session 成功 | 占用回收可观测 |

### 8.4 回归（不得破坏）

| ID | 场景 | 期望 |
|----|------|------|
| R1 | 无 session 最低延迟 | 行为与现网一致 |
| R2 | region 无节点不跨区 | ErrNoNode for region |
| R3 | 节点 fail 换绑同 region | 现有测试仍过 |
| R4 | SessionTTL 热更新 | affinity TTL 仍生效 |
| R5 | `go test -race ./affinity ./selector` | PASS |

### 8.5 实现顺序建议（对应 TODO）

1. **#14**：反向索引 + N 上限 + S1/S2/S3/S4/S7/A*  
2. **#15**：冷却 + S5/S6/S11  
3. **#16**：可观测 API（每节点 occupancy / cooldown_until）

---

## 9. 迁移步骤

### 9.1 配置迁移

1. 代码增加字段默认值：`MaxSessionsPerProxy=1`，`ProxyCooldownMinutes=0`。  
2. 旧 `config.json` 缺字段 → Load 填默认（与现有 `omitempty` 模式一致）。  
3. **行为变化注意**：默认 N=1 会改变「多 session 可撞同一最快节点」的现状 → **破坏性**。  

**迁移开关策略（定稿）**：

| 阶段 | 默认 N | 说明 |
|------|--------|------|
| 合并实现时 | 环境/文档写明 | 若需零行为变化，可临时默认 N 极大（如 0 表示无限） |
| **本规格推荐** | `0` = 无限，`≥1` = 生效；**发布默认 `0` 保持兼容**；文档推荐自动化场景设 `1` | 与 TODO 文案「默认 1」冲突时以 **兼容优先** 为准 |

**兼容优先定稿**：

- `max_sessions_per_proxy`：`0` 或未配置 = **不限制**；`≥1` = 限制。  
- 产品文档示例与自动化推荐值：`1`。  
- `proxy_cooldown_minutes`：`0` = 关（已一致）。

（若产品坚持默认 N=1，须在 CHANGELOG 标 **breaking** 并给回滚：`MAX_SESSIONS_PER_PROXY=0`。）

### 9.2 运行时迁移

1. 滚动/重启后 affinity 为空 → 无历史占用，无需数据迁移。  
2. 先部署仅含指标/日志的版本（可选）→ 再打开 N/CD。  
3. WebUI 设置页增加两字段（可跟 #16）；保存路径复用 `config.Save` + 内存 cfg 更新。  
4. 回滚：设 N=0、CD=0 或回退二进制；无磁盘脏状态。

### 9.3 文档与发布

1. `README` / `.env.example` 增加两变量说明。  
2. `CHANGELOG` Unreleased：行为、默认、breaking 与否。  
3. 本文件随实现勾选「已实现」或拆到 PRD 附录。

### 9.4 实现切片（供 #14/#15）

```text
切片1 affinity 反向索引 + CountByProxy + 测试 A1–A4
切片2 selector 过滤占用 + 配置 N + 测试 S1–S4,S7,S9
切片3 冷却表 + 过滤 + 测试 S5,S6,S11,A5–A7
切片4 入口错误映射回归 I1–I2 + race
切片5（可选 #16）API 聚合占用/冷却
```

---

## 10. 可观测性（预告，归属 #16）

认证后 API 建议字段（设计预留，本规格不实现）：

```json
{
  "proxy_id": 12,
  "address": "…",
  "active_sessions": 1,
  "max_sessions": 1,
  "cooldown_remaining_seconds": 120
}
```

约束：无密码、无完整上游凭据；address 可与 sessions API 同样脱敏策略。

---

## 11. 决策摘要

| # | 决策 | 选择 |
|---|------|------|
| D1 | 租约存储 | 无独立表；affinity binding = 租约 |
| D2 | 占用计数 | proxy→session 反向索引 |
| D3 | N 默认 | 配置 `0`=无限（兼容）；推荐自动化 `1` |
| D4 | 冷却默认 | `0`=关 |
| D5 | 冷却 vs 粘性 | 粘性免疫冷却 |
| D6 | 冷却 vs 多配额 | 冷却仍挡新绑定 |
| D7 | 无 session | 不占配额、不触发冷却 |
| D8 | 失败 | 统一 ErrNoNode 族，不跨 region |
| D9 | 持久化 | 否 |
| D10 | 热更新 N 变小 | 不踢现有 session |

---

## 12. 非目标与后续

- 分布式租约 / Redis  
- 按用户名（认证 base user）的全局 session 配额  
- 连接级（TCP）并发限制（异于 session 配额）  
- 冷却写入 SQLite  
- 改变健康检查或 fail_count 语义  

---

## 13. 验收对照（TODO #13）

| 要求 | 本节 |
|------|------|
| 定义 max_sessions_per_proxy | §3、§6、D3 |
| 定义 proxy_cooldown_minutes / CD | §3、§6.2、D4 |
| 与 affinity 关系 | §5 |
| Resolve 变更点 | §6 |
| 失败语义 | §7 |
| 测试矩阵 | §8 |
| 迁移步骤 | §9 |
| 设计 only、默认可评审 | 全文无生产代码 |
