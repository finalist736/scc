package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"scc/termm"
	"syscall"
)

func main() {
	flag.Parse()

	interruptChannel := make(chan os.Signal)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	t := termm.Terminal{}
	go func() {
		err := t.Start()
		if err != nil {
			log.Println(err)
		}
		interruptChannel <- syscall.SIGTERM
	}()
	for {
		<-interruptChannel
		t.Stop()
		break
	}
}
