package main

import "github.com/go-gl/mathgl/mgl32"

type Player struct {
	Touch

	pos, prevPos mgl32.Vec2
	size         mgl32.Vec2
	force        mgl32.Vec2
}

func (p *Player) Position() mgl32.Vec2 {
	return p.pos
}

func (p *Player) PrevPosition() mgl32.Vec2 {
	return p.prevPos
}

func (p *Player) Size() mgl32.Vec2 {
	return p.size
}

func (p *Player) Render(renderer Renderer) {
	DrawRectColor(renderer, p.pos, p.size, 0, 255, 0)
}

func (p *Player) SetPosition(pos mgl32.Vec2) {
	p.pos = pos
}

func (p *Player) Force() mgl32.Vec2 {
	return p.force
}

func (p *Player) SetForce(force mgl32.Vec2) {
	p.force = force
}
