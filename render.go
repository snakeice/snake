package main

import (
	"time"

	tm "github.com/buger/goterm"
)

var frame = 0
var fps = -1
var stopRender chan bool = make(chan bool, 1)

type ScreenBox struct {
	MaxX int
	MaxY int
	MinX int
	MinY int
}

var printBox *tm.Box = tm.NewBox(80|tm.PCT, 70|tm.PCT+1, 0)

var Box ScreenBox = ScreenBox{
	MinY: 3,
	MinX: 3,
	MaxX: printBox.Width,
	MaxY: printBox.Height,
}

func renderEngine(game *GameStruct) {
	tickTimeout := time.Second / 30
	go fpsCalc()
	for {
		frame++
		tm.Clear()
		tm.Print(tm.MoveTo(printBox.String(), 2, 2))
		tm.Background(" ", tm.BLUE)
		tm.MoveCursor(1, 1)
		select {
		case <-stopRender:
			return
		default:
			renderPanel(game)
			renderPlayer(game)
			renderFruits(game)

			tm.MoveCursor(tm.Width(), tm.Height())
			tm.Flush()
			time.Sleep(tickTimeout)
		}
	}
}

func fpsCalc() {
	for {
		initFrame := frame
		time.Sleep(time.Second)
		fps = frame - initFrame
	}
}

func renderPanel(game *GameStruct) {

	tm.Print(" Snake console super game :3")
	tm.Printf("\tPoints: %d", game.Points)
	tm.Printf("\tFPS: %d", fps)
	tm.Printf("\tPlayer: %v", game.Player.position[0])
	tm.Print("\n")
}

func renderPlayer(game *GameStruct) {
	for index, cell := range game.Player.position {
		tm.MoveCursor(cell.x, cell.y)

		if index%2 == 0 {
			tm.Print(tm.Color("█", tm.BLUE))
		} else {
			tm.Print(tm.Color("█", tm.CYAN))
		}
	}
}

func renderFruits(game *GameStruct) {
	for _, fruit := range game.Fruits {
		tm.MoveCursor(fruit.Position.x, fruit.Position.y)
		switch fruit.Type {
		case Sweet:
			tm.Print(tm.Color(tm.Bold("@"), tm.GREEN))
		case Poison:
			tm.Print(tm.Color(tm.Bold("X"), tm.RED))
		}
	}
}
