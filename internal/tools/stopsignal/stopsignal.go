package stopsignal

import (
	"os"
	"os/signal"
	"syscall"
)

// Stop - функция блокирует выполнение горутины в которой вызвана пока не поступит сигнал прерывания
func Stop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
