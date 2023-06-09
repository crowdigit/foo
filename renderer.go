package main

import "github.com/go-gl/mathgl/mgl32"

type Renderer interface {
	VAO() uint32
	Proj() mgl32.Mat4
	Program() uint32
	PVMUniformLoc() int32
	ColorUniformLoc() int32
}

type RectRenderer struct {
	vao             uint32
	program         uint32
	pvmUniformLoc   int32
	colorUniformLoc int32
}

func (r RectRenderer) VAO() uint32 {
	return r.vao
}

func (r RectRenderer) Proj() mgl32.Mat4 {
	return mgl32.Mat4{
		2.0 / SCREEN_WIDTH, 0, 0, 0,
		0, 2.0 / SCREEN_HEIGHT, 0, 0,
		0, 0, SCREEN_DEPTH, 0,
		-1, -1, -1, 1,
	}
}

func (r RectRenderer) Program() uint32 {
	return r.program
}

func (r RectRenderer) PVMUniformLoc() int32 {
	return r.pvmUniformLoc
}

func (r RectRenderer) ColorUniformLoc() int32 {
	return r.colorUniformLoc
}
