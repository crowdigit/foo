package main

type Block struct {
	x, y int
	w, h int
}

func (b Block) Position() (int, int) {
	return b.x, b.y
}

func (b Block) Render(renderer Renderer) {
	DrawRectColor(renderer, b.x, b.y, b.w, b.h, 255, 0, 0)
}
