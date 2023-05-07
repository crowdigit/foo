package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	EMPTY  = 0
	PLAYER = 1

	MAP_WIDTH  = 10
	MAP_HEIGHT = 10
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

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		render := [MAP_WIDTH][MAP_HEIGHT]byte{}
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
		case "u":
			player.y = Min(player.y+1, MAP_HEIGHT-1)
		case "d":
			player.y = Max(player.y-1, 0)
		case "l":
			player.x = Max(player.x-1, 0)
		case "r":
			player.x = Min(player.x+1, MAP_WIDTH-1)
		case "q":
			break loop
		}

		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	}
}
