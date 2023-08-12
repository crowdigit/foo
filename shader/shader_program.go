package shader

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Parameter interface {
	Proj() mgl32.Mat4
	VAO() uint32
}

type ColorProgram interface {
	Render(pos mgl32.Vec3, size mgl32.Vec2, r, g, b uint8)
}

type ColorProgramImpl struct {
	parameter       Parameter
	id              uint32
	pvmUniformLoc   int32
	colorUniformLoc int32
}

type TextureProgram interface {
	Render(pos mgl32.Vec3, size mgl32.Vec2, texID uint32)
}

type TextureProgramImpl struct {
	parameter         Parameter
	id                uint32
	pvmUniformLoc     int32
	samplerUniformLoc int32
}

func (p ColorProgramImpl) Render(pos mgl32.Vec3, size mgl32.Vec2, r, g, b uint8) {
	r_, g_, b_ := normalizeColor(r, g, b)

	gl.BindVertexArray(p.parameter.VAO())

	matModel := mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()).
		Mul4(mgl32.Scale3D(size.X(), size.Y(), 1))
	matProjModel := p.parameter.Proj().Mul4(matModel)

	gl.UseProgram(p.id)
	gl.UniformMatrix4fv(p.pvmUniformLoc, 1, false, &matProjModel[0])
	gl.Uniform3f(p.colorUniformLoc, r_, g_, b_)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	gl.BindVertexArray(0)
}

func normalizeUint8(a uint8) float32 {
	return float32(a) / 256
}

func normalizeColor(r, g, b uint8) (float32, float32, float32) {
	return normalizeUint8(r), normalizeUint8(g), normalizeUint8(b)
}

func (p TextureProgramImpl) Render(pos mgl32.Vec3, size mgl32.Vec2, texID uint32) {
	gl.BindVertexArray(p.parameter.VAO())

	matModel := mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()).
		Mul4(mgl32.Scale3D(size.X(), size.Y(), 1))
	matProjModel := p.parameter.Proj().Mul4(matModel)

	gl.UseProgram(p.id)
	gl.UniformMatrix4fv(p.pvmUniformLoc, 1, false, &matProjModel[0])

	gl.BindTexture(gl.TEXTURE_2D, texID)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	gl.BindVertexArray(0)
}
