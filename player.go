package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	Touch

	pos, prevPos mgl32.Vec2
	size         mgl32.Vec2
	force        mgl32.Vec2

	texId uint32
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

func (p *Player) Render(renderer RectRenderer) {
	if renderer, ok := renderer.(TextureRenderer); ok {
		gl.BindVertexArray(renderer.VAO())

		matModel := mgl32.Translate3D(p.pos.X(), p.pos.Y(), OFFSET_Z_BLOCK).
			Mul4(mgl32.Scale3D(p.size.X(), p.size.Y(), 1))
		matProjModel := renderer.Proj().Mul4(matModel)

		gl.UseProgram(renderer.Program())
		gl.UniformMatrix4fv(renderer.PVMUniformLoc(), 1, false, &matProjModel[0])

		gl.BindTexture(gl.TEXTURE_2D, p.texId)

		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		gl.BindVertexArray(0)
	}
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
