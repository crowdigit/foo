package shader

import (
	_ "embed"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pkg/errors"
)

//go:embed colorRect.glsl
var colorVertSource string

//go:embed colorFrag.glsl
var colorFragSource string

//go:embed rect.glsl
var textureVertSource string

//go:embed frag.glsl
var textureFragSource string

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

func loadShader(vertShader, fragShader string) (uint32, error) {
	colorVertShader, err := compileShaderSource(vertShader, gl.VERTEX_SHADER)
	if err != nil {
		return 0, errors.Wrap(err, "failed to compile vertex shader source")
	}
	defer gl.DeleteShader(colorVertShader)

	colorFragShader, err := compileShaderSource(fragShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, errors.Wrap(err, "failed to compile fragment shader source")
	}
	defer gl.DeleteShader(colorFragShader)

	program, err := linkProgram(colorVertShader, colorFragShader)
	if err != nil {
		return 0, errors.Wrap(err, "failed to link shaders")
	}

	return program, nil
}

func LoadColorShader(parameter Parameter) (ColorProgramImpl, error) {
	id, err := loadShader(colorVertSource, colorFragSource)
	if err != nil {
		return ColorProgramImpl{}, errors.Wrap(err, "failed to load shader")
	}

	gl.UseProgram(id)
	pvmUniformLoc := gl.GetUniformLocation(id, gl.Str("projModel"+"\x00"))
	colorUniformLoc := gl.GetUniformLocation(id, gl.Str("color"+"\x00"))

	return ColorProgramImpl{
		parameter:       parameter,
		id:              id,
		pvmUniformLoc:   pvmUniformLoc,
		colorUniformLoc: colorUniformLoc,
	}, nil
}

func LoadTextureShader(parameter Parameter) (TextureProgramImpl, error) {
	id, err := loadShader(textureVertSource, textureFragSource)
	if err != nil {
		return TextureProgramImpl{}, errors.Wrap(err, "failed to load shader")
	}

	gl.UseProgram(id)
	pvmUniformLoc := gl.GetUniformLocation(id, gl.Str("projModel"+"\x00"))
	samplerUniformLoc := gl.GetUniformLocation(id, gl.Str("sampler"+"\x00"))

	return TextureProgramImpl{
		parameter:         parameter,
		id:                id,
		pvmUniformLoc:     pvmUniformLoc,
		samplerUniformLoc: samplerUniformLoc,
	}, nil
}
