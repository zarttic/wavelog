package wavelog

import (
	"testing"
)

func TestLog(t *testing.T) {
	Trace("trace", "test")
	Debug("debug", "test")
	Warning("waring", "test")
	Info("info", "test")
	Error("error", "test")
	Fatal("fatal", "test")

}
