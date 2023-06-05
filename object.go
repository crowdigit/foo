package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object interface {
	Render(Renderer)
}

type Player struct {
	x, y    int
	r, g, b uint8
	speed   int
}

func (p Player) Position() (int, int) {
	return p.x, p.y
}

func (p Player) Render(renderer Renderer) {
	gl.BindVertexArray(renderer.VAO())

	x := float32(p.x*MAP_GRID_SIZE + RENDER_CENTER_OFFSET_X)
	y := float32(p.y*MAP_GRID_SIZE + RENDER_CENTER_OFFSET_Y)

	matModel := mgl32.Translate2D(x, y).Mul3(mgl32.Scale2D(MAP_GRID_SIZE, MAP_GRID_SIZE))
	matProjModel := renderer.Proj().Mul3(matModel)

	gl.UseProgram(renderer.Program())
	gl.UniformMatrix3fv(renderer.PVMUniformLoc(), 1, false, &matProjModel[0])

	// TODO render queue
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	gl.BindVertexArray(0)
}
