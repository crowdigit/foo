package main

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
	tiles := []string{".", "x"}

	objects := map[string]Object{
		"player":  Player{x: 0, y: 0},
		"monster": Player{x: 9, y: 9},
	}

	render := [10][10]string{}
	for y, row := range render {
		for x := range row {
			render[x][y] = tiles[EMPTY]
		}
	}

	for _, object := range objects {
		x, y := object.Position()
		render[x][y] = tiles[PLAYER]
	}

	for _, row := range render {
		for _, cell := range row {
			print(cell)
		}
		println()
	}
}
