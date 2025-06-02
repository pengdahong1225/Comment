package signalHandler

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type SignalHandler struct {
	reloadFunc  func()
	onlineFunc  func()
	offlineFunc func()
}

func NewSignalHandler(reload, online, offline func()) *SignalHandler {
	return &SignalHandler{
		reloadFunc:  reload,
		onlineFunc:  online,
		offlineFunc: offline,
	}
}
func (h *SignalHandler) Listen() {
	sigChan := make(chan os.Signal, 2)
	// 自定义信号
	signal.Notify(sigChan,
		syscall.Signal(34+1), // 35 -> reload
		syscall.Signal(34+2), // 36 -> online
		syscall.Signal(34+3), // 37 -> offline
		syscall.Signal(34+4), // 38 -> open debug log
		syscall.Signal(34+5), // 39 -> close debug log
	)

	for sig := range sigChan {
		switch sig {
		case syscall.Signal(34+1):
			h.reloadFunc()
		case syscall.Signal(34+2):
			h.onlineFunc()
		case syscall.Signal(34+3):
			h.offlineFunc()
		default:
			logrus.Errorf("received signal: %s, undo\n", sig)
		}
	}
}
