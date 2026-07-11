# AGENTS.md — 本地环境与工作纪律（goproxy）

本文件记录本机固定事实与本项目工作纪律，避免每次重新发现。与 CLAUDE.md 并存：CLAUDE.md 讲架构，本文件讲“怎么在这台机器上干活”。

## 本地环境（已验证，2026-07）

- **Go**: `C:\Program Files\Go\bin\go.exe`，go1.26.4 windows/amd64（`go` 不在默认 PATH，需用全路径或临时加 PATH）。
- **gcc (cgo 必需)**: `C:\ProgramData\mingw64\mingw64\bin\gcc.exe`（mingw64 16.1.0）。
- **sing-box**: `C:\Users\babut\go\bin\sing-box.exe`（在 PATH 上，`sing-box` 直接可用），版本 1.13.14。
  - 其它副本：`D:\Program Files (x86)\Sing-box\sing-box.exe`、`D:\Program Files (x86)\v2rayN-With-Core\bin\sing_box\sing-box.exe`。
- **测试订阅地址（禁止写入任何会提交的文件；仅本地运行时使用）**：
  - sing-box: 存于本机私密处，不入库。
  - clash: 同上。
  - 红线：这些 URL 绝不写进代码/测试 fixture/配置/日志/AGENTS.md/memory/qdrant。

## 跑测试（必须开 CGO）

storage/custom/webui/selector 依赖 go-sqlite3，需 `CGO_ENABLED=1` + gcc 在 PATH。标准命令：

```pwsh
$env:CGO_ENABLED=1
$env:PATH = "C:\ProgramData\mingw64\mingw64\bin;C:\Program Files\Go\bin;$env:PATH"
go test ./...
```

- **不要**再以“CGO_ENABLED=0 无 gcc”为由跳过 storage/webui 测试——本机有 gcc，必须实跑。
- 真实 sing-box 集成测试**可以且必须**在本机跑（含 6000+ 节点规模验证），不得以“无 sing-box 二进制”为由回避。

## 工作纪律（本项目）

- **TDD 强制**：任何生产代码改动先写失败测试（RED），粘贴 RED 输出，再最小实现（GREEN），后 REFACTOR。测试首跑即过 = 无效，重来。
- **派 task 必须带 TDD 约束**，并要求粘贴 RED 证据；不符合验收标准打回重做，保留上下文 loop 直到通过严苛审计。
- **审计不采信 task 自述**：亲自复跑 build/vet/test、读 diff、查 git status。
- **验收硬门槛**：`go build ./...` EXIT 0；`go vet ./...` EXIT 0；相关测试全绿；gofmt 净；仅改动预期文件；无红线违规。
- **真实环境验证优先**：能在本机端到端证实的（真实 sing-box、6000 节点、cgo 全套测试），必须实测，不写“只能部署环境验证”的免责话。
- **提交需显式授权**：不擅自 commit/push/amend。
- **举一反三**：修一个 bug 时排查同类；改一处契约时全仓查消费方。

## 路 C 现状（sing-box 分片多进程，进行中）

- 已完成（TDD + 审计）：mixed 单端口(缺陷5)、ShardedSingBox 编排器、config.SingBoxShardCount(默认4可配)、接入 Manager、端口泄漏(缺陷4)+分片段超限保护。
- 均为**未提交**工作树状态。`stash@{0}` 存早期非 TDD 编排器备份。
- 待做：本机真实 sing-box + 6000 节点端到端验证（平滑重启：仅重启变化分片、其余连接不断）。
