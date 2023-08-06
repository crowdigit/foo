package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func checkGLError() {
	err := gl.GetError()
	switch err {
	case gl.NO_ERROR:
		fmt.Println("no OpenGL error")
	case gl.INVALID_ENUM:
		fmt.Println("GL_INVALID_ENUM")
	case gl.INVALID_VALUE:
		fmt.Println("GL_INVALID_VALUE")
	case gl.INVALID_OPERATION:
		fmt.Println("GL_INVALID_OPERATION")
	default:
		fmt.Println("unknown OpenGL error")
	}
}

func loadScene(filepath string) (Scene, error) {
	sceneBytes, err := os.ReadFile(filepath)
	if err != nil {
		return Scene{}, errors.Wrap(err, "failed to read scene file")
	}

	scene := Scene{}
	if err := json.Unmarshal(sceneBytes, &scene); err != nil {
		return Scene{}, errors.Wrap(err, "failed to unmarshal scene file into scene type")
	}

	return scene, nil
}

func configureOpenGL() {
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
}

func initRenderer() (Renderer, error) {
	if err := gl.Init(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize OpenGL")
	}

	shader, err := loadShader()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load shader")
	}

	vao, err := initVAO()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize VAO")
	}

	gl.UseProgram(shader)
	pvmUniformLoc := gl.GetUniformLocation(shader, gl.Str("projModel"+"\x00"))
	colorUniformLoc := gl.GetUniformLocation(shader, gl.Str("color"+"\x00"))

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1)

	return RendererImpl{
		vao:             vao,
		program:         shader,
		pvmUniformLoc:   pvmUniformLoc,
		colorUniformLoc: colorUniformLoc,
	}, nil
}
