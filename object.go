package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	OFFSET_Z_PLAYER = 1
	OFFSET_Z_GRID   = 0
)

type Object interface {
	Render(Renderer)
}

func normalizeUint8(a uint8) float32 {
	return float32(a) / 256
}

func normalizeColor(r, g, b uint8) (float32, float32, float32) {
	return normalizeUint8(r), normalizeUint8(g), normalizeUint8(b)
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

	matModel := mgl32.Translate3D(x, y, OFFSET_Z_PLAYER).Mul4(mgl32.Scale3D(MAP_GRID_SIZE, MAP_GRID_SIZE, 1))
	matProjModel := renderer.Proj().Mul4(matModel)

	r, g, b := normalizeColor(p.r, p.g, p.b)

	gl.UseProgram(renderer.Program())
	gl.UniformMatrix4fv(renderer.PVMUniformLoc(), 1, false, &matProjModel[0])
	gl.Uniform3f(renderer.ColorUniformLoc(), r, g, b)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	gl.BindVertexArray(0)
}
