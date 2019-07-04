package logs

import (
	"testing"
	"time"
)

func TestNewMyLogger(t *testing.T) {
	logger := NewMyLogger(WARNING)
	logger.Register(AddConsole(), AddFile(true, `{"path":"runtime", "save_name":"test", "log_file_ext":"log"}`))
	logger.DEBUG("123", "456")
	time.Sleep(40 * time.Second)
	logger.DEBUG("678, 3445")
	time.Sleep(time.Minute)
}
