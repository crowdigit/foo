package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type RectRenderer interface {
	VAO() uint32
	Proj() mgl32.Mat4
	Program() uint32
	PVMUniformLoc() int32
}

type RectRendererImpl struct {
	vao           uint32
	program       uint32
	pvmUniformLoc int32
}

type ColorRenderer interface {
	RectRenderer
	ColorUniformLoc() int32
}

type ColorRendererImpl struct {
	RectRendererImpl
	colorUniformLoc int32
}

type TextureRenderer interface {
	RectRenderer
	SamplerLoc() int32
}

type TextureRendererImpl struct {
	RectRendererImpl
	samplerUniformLoc int32
}

func (r RectRendererImpl) VAO() uint32 {
	return r.vao
}

func (r RectRendererImpl) Proj() mgl32.Mat4 {
	return mgl32.Mat4{
		2.0 / SCREEN_WIDTH, 0, 0, 0,
		0, 2.0 / SCREEN_HEIGHT, 0, 0,
		0, 0, SCREEN_DEPTH, 0,
		-1, -1, -1, 1,
	}
}

func (r RectRendererImpl) Program() uint32 {
	return r.program
}

func (r RectRendererImpl) PVMUniformLoc() int32 {
	return r.pvmUniformLoc
}

func (r ColorRendererImpl) ColorUniformLoc() int32 {
	return r.colorUniformLoc
}

func (r TextureRendererImpl) SamplerLoc() int32 {
	return r.samplerUniformLoc
}

func DrawRectColor(renderer ColorRenderer, pos, size mgl32.Vec2, r, g, b uint8) {
	r_, g_, b_ := normalizeColor(r, g, b)

	gl.BindVertexArray(renderer.VAO())

	matModel := mgl32.Translate3D(pos.X(), pos.Y(), OFFSET_Z_BLOCK).
		Mul4(mgl32.Scale3D(size.X(), size.Y(), 1))
	matProjModel := renderer.Proj().Mul4(matModel)

	gl.UseProgram(renderer.Program())
	gl.UniformMatrix4fv(renderer.PVMUniformLoc(), 1, false, &matProjModel[0])
	gl.Uniform3f(renderer.ColorUniformLoc(), r_, g_, b_)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	gl.BindVertexArray(0)
}
