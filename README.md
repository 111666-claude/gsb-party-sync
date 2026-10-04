# party-sync

队伍状态版本向量：逐成员序号、幂等应用、乱序暂存补洞、差异检测与修复，只用 Go 标准库。

```
go run ./cmd/party-sync --sample=idem
go run ./cmd/party-sync --sample=gap
go run ./cmd/party-sync --sample=diff
go run ./cmd/party-sync --sample=missing
go run ./cmd/party-sync --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **版本向量**：每个成员各自维护一个从 1 开始、逐次加一的序号，成员之间互不比较。
- **幂等与顺序**：同一个成员同一个序号只应用一次；序号等于当前序号加一时当场应用并返回 `true`；
  已经应用过的序号或不小于 1 的序号回退返回 `false` 且不改状态。
- **乱序暂存**：序号大于当前序号加一时先暂存，不计入版本向量；缺口补齐后按序号连续级联应用，
  暂存清零。
- **差异检测**：`Diff` 只列出本地该成员序号大于对端的成员，`Missing` 只列出本地落后的成员，
  两个结果都按成员名升序，逐成员比较。
- **不变量**：同一组更新按任意顺序应用，最终状态与版本向量一致；重复应用不改变状态。
- **代价**：应用与差异检测不许扫全部历史更新，`Scanned` 不随历史更新数放大；
  10 万成员、每秒 10 万次应用与检测。

## 输出契约（不改格式）

```
state=.. second=..
state=.. version=.. pending=..
diff=..
missing=..
scanned=..
```

## 目录

```
sync.go               版本向量、幂等、暂存补洞与差异检测
sync_test.go          单元测试
cmd/party-sync/       命令行入口
```
