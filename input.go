package main

import (
	"github.com/eiannone/keyboard"
)

func InitInput(buff chan keyboard.Key) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		_, key, err := keyboard.GetKey()
		buff <- key

		if err != nil {
			panic(err)
		}
	}
}
