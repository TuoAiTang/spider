package log

import "testing"

func TestInfo(t *testing.T) {
	Info("test:%v", 3)
}
