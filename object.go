package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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
	gl.BindVertexArray(renderer.VAO())

	matModel := mgl32.Translate2D(0, 0).Mul3(mgl32.Scale2D(100, 100))
	matProjModel := renderer.Proj().Mul3(matModel)

	gl.UseProgram(renderer.Program())
	gl.UniformMatrix3fv(renderer.PVMUniformLoc(), 1, true, &matProjModel[0])

	// TODO render queue
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

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
	gl.BindVertexArray(0)
}
