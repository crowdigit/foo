package main

import (
	_ "embed"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pkg/errors"
)

//go:embed shader/rect.glsl
var vertSource string

//go:embed shader/frag.glsl
var fragSource string

func compileShaderSource(source string, xtype uint32) (uint32, error) {
	shader := gl.CreateShader(xtype)
	if shader == 0 {
		panic("returned shader ID is zero")
	}

	sources, free := gl.Strs(source + "\x00")
	defer free()

	gl.ShaderSource(shader, 1, sources, nil)
	gl.CompileShader(shader)

	var compileStatus int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &compileStatus)
	if compileStatus != gl.TRUE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)))
		gl.GetShaderInfoLog(shader, logLength, nil, log)
		return 0, errors.New(gl.GoStr(log))
	}

	return shader, nil
}

func linkProgram(vertShader, fragShader uint32) (uint32, error) {
	program := gl.CreateProgram()
	if program == 0 {
		return 0, errors.New("failed to create shader program")
	}

	gl.AttachShader(program, vertShader)
	defer gl.DetachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	defer gl.DetachShader(program, fragShader)

	gl.LinkProgram(program)

	var linkStatus int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkStatus)
	if linkStatus != gl.TRUE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)))
		gl.GetProgramInfoLog(program, logLength, nil, log)
		return 0, errors.New(gl.GoStr(log))
	}

	return program, nil
}

func loadShader() (uint32, error) {
	vertShader, err := compileShaderSource(vertSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, errors.Wrap(err, "failed to compile vertex shader source")
	}
	defer gl.DeleteShader(vertShader)

	fragShader, err := compileShaderSource(fragSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, errors.Wrap(err, "failed to compile fragment shader source")
	}
	defer gl.DeleteShader(fragShader)

	program, err := linkProgram(vertShader, fragShader)
	if err != nil {
		return 0, errors.Wrap(err, "failed to link shaders")
	}

	return program, nil
}

func initVAO() (uint32, error) {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	if vao == gl.INVALID_VALUE {
		return 0, errors.New("failed to create VAO")
	}

	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	if vbo == gl.INVALID_VALUE {
		return 0, errors.New("failed to create VBO")
	}

	vertices := []float32{
		0.0, 0.0, 1.0, 0.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 1.0, 0.0, 1.0,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 2*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return vao, nil
}
