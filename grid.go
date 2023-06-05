package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Grid struct {
	//width shows how many columns this grid has
	width int
	//height show how many rows this grid has
	height int

	// Rendering parameters
	lineWidth float32
	r, g, b   uint8
	gridSize  float32
}

func (o Grid) Render(renderer Renderer) {
	gl.BindVertexArray(renderer.VAO())

	gridScaleX := o.gridSize * float32(o.width)
	gridScaleY := o.gridSize * float32(o.height)
	r, g, b := normalizeColor(o.r, o.g, o.b)

	gl.UseProgram(renderer.Program())
	gl.Uniform3f(renderer.ColorUniformLoc(), r, g, b)

	// draw vertical lines
	{
		scaleX := o.lineWidth
		scaleY := gridScaleY + o.lineWidth
		translateY := 1 / float32(2) * (SCREEN_HEIGHT - gridScaleY - o.lineWidth)

		for i := 0; i < o.width+1; i += 1 {
			translateX := 1/float32(2)*(SCREEN_WIDTH-gridScaleX-o.lineWidth) + o.gridSize*float32(i)

			matModel := mgl32.Translate3D(translateX, translateY, OFFSET_Z_GRID).
				Mul4(mgl32.Scale3D(scaleX, scaleY, 1))
			matProjModel := renderer.Proj().Mul4(matModel)
			gl.UniformMatrix4fv(renderer.PVMUniformLoc(), 1, false, &matProjModel[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
		}
	}

	// draw horizontal lines
	{
		scaleX := gridScaleX + o.lineWidth
		scaleY := o.lineWidth
		translateX := 1 / float32(2) * (SCREEN_WIDTH - gridScaleX - o.lineWidth)

		for i := 0; i < o.height+1; i += 1 {
			translateY := 1/float32(2)*(SCREEN_HEIGHT-gridScaleY-o.lineWidth) + o.gridSize*float32(i)

			matModel := mgl32.Translate3D(translateX, translateY, OFFSET_Z_GRID).Mul4(mgl32.Scale3D(scaleX, scaleY, 1))
			matProjModel := renderer.Proj().Mul4(matModel)

			gl.UniformMatrix4fv(renderer.PVMUniformLoc(), 1, false, &matProjModel[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
		}
	}

	gl.BindVertexArray(0)
}
