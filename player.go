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

func (p *Player) Update(keyboard Keyboard) {
	force := p.Force()

	g := mgl32.Vec2{0, -0.098}
	force = force.Add(g)

	if p.TouchingFloor() {
		if keyboard.Space {
			force = force.Add(mgl32.Vec2{0, 5})
		} else if keyboard.Left && keyboard.Right {
			// do nothing
		} else if keyboard.Left {
			force[0] -= 0.1
			force[0] = Max(force[0], -2)
		} else if keyboard.Right {
			force[0] += 0.1
			force[0] = Min(force[0], 2)
		} else if !keyboard.Left && !keyboard.Right {
			if force[0] > 0 {
				force[0] = Max(force[0]-0.08, 0)
			} else {
				force[0] = Min(force[0]+0.08, 0)
			}
		}
	}

	p.prevPos = p.pos
	p.force = force
	p.pos = p.pos.Add(p.force)

	p.ResetTouch()
}
