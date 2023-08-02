package main

type Player struct {
	x, y int
	w, h int
}

func (p Player) Position() (int, int) {
	return p.x, p.y
}

func (p Player) Render(renderer Renderer) {
	DrawRectColor(renderer, p.x, p.y, p.w, p.h, 0, 255, 0)
}
