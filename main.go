package main

import (
	"fmt"
)

const (
	EMPTY  = 0
	PLAYER = 1
)

type Object interface {
	Position() (int, int)
}

type Player struct {
	x int
	y int
}

func (p Player) Position() (int, int) {
	return p.x, p.y
}

func main() {
	tiles := []byte{'.', 'x'}

	player := Player{x: 0, y: 0}

	objects := map[string]Object{
		"player":  &player,
		"monster": Player{x: 9, y: 9},
	}

loop:
	for {
		render := [10][10]byte{}
		for y, row := range render {
			for x := range row {
				render[y][x] = tiles[EMPTY]
			}
		}

		for _, object := range objects {
			x, y := object.Position()
			render[y][x] = tiles[PLAYER]
		}

		for y := len(render) - 1; y >= 0; y -= 1 {
			for _, cell := range render[y] {
				fmt.Printf("%c", cell)
			}
			println()
		}

		var cmd string
		if n, err := fmt.Scanf("%s", &cmd); err != nil || n != 1 {
			panic(err)
		}

		switch string(cmd) {
		case "up":
			player.y += 1
		case "down":
			player.y -= 1
		case "left":
			player.x -= 1
		case "right":
			player.x += 1
		case "exit":
			break loop
		}
	}
}
