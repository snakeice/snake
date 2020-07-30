package main

import (
	"math/rand"
	"os"
	"time"

	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
)

type DirectionType uint
type FruitType uint

const (
	Left  DirectionType = iota
	Down  DirectionType = iota + 1
	Right DirectionType = iota + 1
	Up    DirectionType = iota + 1
)

const (
	Sweet  FruitType = 1
	Poison FruitType = 2
)

type Cell struct {
	x int
	y int
}

type GameStruct struct {
	KeyBuff chan keyboard.Key
	Player  *PlayerControl
	Points  int64
	Fruits  []*Fruit
}

type PlayerControl struct {
	position []*Cell
	speedX   int8
	speedY   int8
}

func (g *GameStruct) Tick() {
	tickInterval := time.Second / 10
	for {
		g.Player.tickPlayer()
		g.eat()
		time.Sleep(tickInterval)
	}
}

func (g *GameStruct) randTime() {
	tickInterval := time.Second
	for {
		if x := rand.Intn(3000-1) + 1; x < 300 {
			g.newFruitRand()
		}
		time.Sleep(tickInterval)
	}
}

func (g *GameStruct) newFruitRand() {
	ty := rand.Intn(3-1) + 1
	g.newFruit(FruitType(ty))
}

func (g *GameStruct) newFruit(fruitType FruitType) {
	x := rand.Intn(Box.MaxX-Box.MinX+1) + Box.MinX
	y := rand.Intn(Box.MaxY-Box.MinY+1) + Box.MinY
	g.Fruits = append(g.Fruits, &Fruit{
		Type: fruitType,
		Position: Cell{
			x: x,
			y: y,
		},
	})

}

func (g *GameStruct) eat() {
	for index, fruit := range g.Fruits {
		if fruit.Position.x == g.Player.position[0].x && fruit.Position.y == g.Player.position[0].y {
			if fruit.Type == Sweet {
				g.Points++
				g.Fruits = append(g.Fruits[:index], g.Fruits[index+1:]...)
				g.newFruit(Sweet)
				g.Player.position = append(g.Player.position, g.Player.position[len(g.Player.position)-1:]...)
			} else {
				stopRender <- true
				tm.Clear()
				tm.MoveCursor(0, 0)
				tm.Printf("You lose!\n\nPoints: %d", g.Points)
				tm.Flush()
				os.Exit(0)
			}
		}
	}
}

func MakeGame() *GameStruct {
	rand.Seed(time.Now().UnixNano())
	game := &GameStruct{
		KeyBuff: make(chan keyboard.Key, 30),
		Player:  makePlayer(),
		Points:  0,
		Fruits:  []*Fruit{},
	}

	game.newFruit(Poison)
	game.newFruit(Sweet)

	go game.listenerKeyboard()
	return game
}

func makePlayer() *PlayerControl {
	player := PlayerControl{
		speedX: 1,
		speedY: 0,
		position: []*Cell{
			{x: 10, y: 10},
			{x: 9, y: 10},
			{x: 8, y: 10},
		},
	}

	return &player
}

func (g *GameStruct) listenerKeyboard() {
	for {
		key := <-g.KeyBuff
		switch key {
		case keyboard.KeyArrowDown:
			g.Player.changeDirection(Down)
		case keyboard.KeyArrowLeft:
			g.Player.changeDirection(Left)
		case keyboard.KeyArrowRight:
			g.Player.changeDirection(Right)
		case keyboard.KeyArrowUp:
			g.Player.changeDirection(Up)

		}
	}
}

func (p *PlayerControl) changeDirection(dir DirectionType) {
	p.speedX = 0
	p.speedY = 0

	switch dir {
	case Left:
		p.speedX = -1
	case Right:
		p.speedX = 1
	case Down:
		p.speedY = 1
	case Up:
		p.speedY = -1
	}
}

func (p *PlayerControl) tickPlayer() {
	cell := &Cell{
		x: p.position[0].x,
		y: p.position[0].y,
	}

	if p.speedX > 0 {
		if cell.x < Box.MaxX {
			cell.x += int(p.speedX)
		} else {
			cell.x = Box.MinX
		}
	}

	if p.speedX < 0 {
		if cell.x > Box.MinX {
			cell.x += int(p.speedX)
		} else {
			cell.x = Box.MaxX
		}
	}

	if p.speedY > 0 {
		if cell.y < Box.MaxY {
			cell.y += int(p.speedY)
		} else {
			cell.y = Box.MinY
		}
	}

	if p.speedY < 0 {
		if cell.y > Box.MinY {
			cell.y += int(p.speedY)
		} else {
			cell.y = Box.MaxY
		}
	}

	p.position = append([]*Cell{cell}, p.position[:len(p.position)-1]...)
}

type Fruit struct {
	Type     FruitType
	Position Cell
}
