package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	game := MakeGame()

	go renderEngine(game)
	go InitInput(game.KeyBuff)
	go game.Tick()
	go game.randTime()
	<-sigs
}
