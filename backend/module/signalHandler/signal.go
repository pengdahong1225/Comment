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
		syscall.SIGRTMIN,   // 34 -> reload
		syscall.SIGRTMIN+1, // 35 -> online
		syscall.SIGRTMIN+2, // 36 -> offline
		syscall.SIGRTMIN+3, // 37 -> open debug log
		syscall.SIGRTMIN+4, // 38 -> close debug log
	)

	for sig := range sigChan {
		switch sig {
		case syscall.SIGRTMIN:
			h.reloadFunc()
		case syscall.SIGRTMIN + 1:
			h.onlineFunc()
		case syscall.SIGRTMIN + 1:
			h.offlineFunc()
		default:
			logrus.Errorf("received signal: %s, undo\n", sig)
		}
	}
}
