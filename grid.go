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

func (p Grid) Render(renderer Renderer) {
	gl.BindVertexArray(renderer.VAO())

	gridScaleX := p.gridSize * float32(p.width)
	gridScaleY := p.gridSize * float32(p.height)

	gl.UseProgram(renderer.Program())

	// draw vertical lines
	{
		scaleX := p.lineWidth
		scaleY := gridScaleY + p.lineWidth
		translateY := 1 / float32(2) * (SCREEN_HEIGHT - gridScaleY - p.lineWidth)

		for i := 0; i < p.width+1; i += 1 {
			translateX := 1/float32(2)*(SCREEN_WIDTH-gridScaleX-p.lineWidth) + p.gridSize*float32(i)

			matModel := mgl32.Translate2D(translateX, translateY).Mul3(mgl32.Scale2D(scaleX, scaleY))
			matProjModel := renderer.Proj().Mul3(matModel)
			gl.UniformMatrix3fv(renderer.PVMUniformLoc(), 1, false, &matProjModel[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
		}
	}

	// draw horizontal lines
	{
		scaleX := gridScaleX + p.lineWidth
		scaleY := p.lineWidth
		translateX := 1 / float32(2) * (SCREEN_WIDTH - gridScaleX - p.lineWidth)

		for i := 0; i < p.height+1; i += 1 {
			translateY := 1/float32(2)*(SCREEN_HEIGHT-gridScaleY-p.lineWidth) + p.gridSize*float32(i)

			matModel := mgl32.Translate2D(translateX, translateY).Mul3(mgl32.Scale2D(scaleX, scaleY))
			matProjModel := renderer.Proj().Mul3(matModel)

			gl.UniformMatrix3fv(renderer.PVMUniformLoc(), 1, false, &matProjModel[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
		}
	}

	gl.BindVertexArray(0)
}
