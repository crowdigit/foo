package main

type Renderer interface{}

type Object interface {
	Position() (int, int)
	Render(Renderer)
}

type Player struct {
	x int
	y int

	r, g, b uint8
	speed   int
}

func (p Player) Position() (int, int) {
	return p.x, p.y
}

func (p Player) Render(renderer Renderer) {
	/*
		rect := sdl.Rect{
			X: int32(p.x)*MAP_GRID_SIZE + RENDER_CENTER_OFFSET_X,
			Y: int32(-p.y+MAP_HEIGHT-1)*MAP_GRID_SIZE + RENDER_CENTER_OFFSET_Y,
			W: MAP_GRID_SIZE,
			H: MAP_GRID_SIZE,
		}
		renderer.SetDrawColor(p.r, p.g, p.b, 0)
		if err := renderer.FillRect(&rect); err != nil {
			panic(err)
		}
	*/
}
