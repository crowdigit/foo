package main

import "github.com/go-gl/mathgl/mgl32"

type Block struct {
	pos, size mgl32.Vec2
}

func (b Block) Position() mgl32.Vec2 {
	return b.pos
}

func (b Block) PrevPosition() mgl32.Vec2 {
	return b.pos
}

func (b Block) Size() mgl32.Vec2 {
	return b.size
}

func (b Block) Render(renderer Renderer) {
	DrawRectColor(renderer, b.pos, b.size, 255, 0, 0)
}

func (b Block) Move(mgl32.Vec2) {}
