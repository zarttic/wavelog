## WaveLog
轻量化，易上手的异步日志库
## 食用指南
### 安装
`go get github.com/zarttic/wavelog`
### 使用
```go
	Trace("trace", "test")
	Debug("debug", "test")
	Warning("waring", "test")
	Info("info", "test")
	Error("error", "test")
	Fatal("fatal", "test")
```
