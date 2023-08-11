package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

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

func (p *Player) Update(scene Scene, keyboard Keyboard) {
	for _, block := range scene.Blocks {
		if collides, _ := CheckCollision(p, block); collides {
			ResolveCollision(p, block)
		}
	}

	force := p.Force()

	gravity := mgl32.Vec2{0, -0.098}
	force = force.Add(gravity)

	const friction float32 = 0.08
	const accel float32 = 0.1
	const accelLimit float32 = 2

	if p.TouchingFloor() {
		if keyboard[sdl.SCANCODE_SPACE].Press {
			force = force.Add(mgl32.Vec2{0, 5})
		} else if keyboard[sdl.SCANCODE_LEFT].Press == keyboard[sdl.SCANCODE_RIGHT].Press {
			if force[0] > 0 {
				force[0] = Max(force[0]-friction, 0)
			} else {
				force[0] = Min(force[0]+friction, 0)
			}
		} else if keyboard[sdl.SCANCODE_LEFT].Press {
			force[0] = Max(force[0]-accel, -accelLimit)
		} else if keyboard[sdl.SCANCODE_RIGHT].Press {
			force[0] = Min(force[0]+accel, accelLimit)
		}
	}

	p.prevPos = p.pos
	p.force = force
	p.pos = p.pos.Add(p.force)

	p.ResetTouch()
}
