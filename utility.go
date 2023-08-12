package main

import (
	"encoding/json"
	"fmt"
	"os"
	"unsafe"

	"github.com/crowdigit/foo/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/img"
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

func initShader(vertShader, fragShader string) (uint32, error) {
	shader, err := shader.LoadShader(vertShader, fragShader)
	if err != nil {
		return 0, errors.Wrap(err, "failed to load vertex shader")
	}
	return shader, nil
}

func newColorRenderer() (ColorRendererImpl, error) {
	colorShader, err := shader.LoadShader(shader.ColorVertSource, shader.ColorFragSource)
	if err != nil {
		return ColorRendererImpl{}, errors.Wrap(err, "failed to load shader")
	}

	vao, err := shader.InitVAO()
	if err != nil {
		return ColorRendererImpl{}, errors.Wrap(err, "failed to initialize VAO")
	}

	gl.UseProgram(colorShader)
	pvmUniformLoc := gl.GetUniformLocation(colorShader, gl.Str("projModel"+"\x00"))
	colorUniformLoc := gl.GetUniformLocation(colorShader, gl.Str("color"+"\x00"))

	return ColorRendererImpl{
		RectRendererImpl: RectRendererImpl{
			vao:           vao,
			program:       colorShader,
			pvmUniformLoc: pvmUniformLoc,
		},
		colorUniformLoc: colorUniformLoc,
	}, nil
}

func newTextureRenderer() (TextureRendererImpl, error) {
	textureShader, err := shader.LoadShader(shader.TextureVertSource, shader.TextureFragSource)
	if err != nil {
		return TextureRendererImpl{}, errors.Wrap(err, "failed to load shader")
	}

	vao, err := shader.InitVAO()
	if err != nil {
		return TextureRendererImpl{}, errors.Wrap(err, "failed to initialize VAO")
	}

	gl.UseProgram(textureShader)
	pvmUniformLoc := gl.GetUniformLocation(textureShader, gl.Str("projModel"+"\x00"))
	samplerUniformLoc := gl.GetUniformLocation(textureShader, gl.Str("sampler"+"\x00"))

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1)

	return TextureRendererImpl{
		RectRendererImpl: RectRendererImpl{
			vao:           vao,
			program:       textureShader,
			pvmUniformLoc: pvmUniformLoc,
		},
		samplerUniformLoc: samplerUniformLoc,
	}, nil
}

// reverseTexture reverses unsigned byte RGBA8888 image horizontally
func reverseTexture(data unsafe.Pointer, w, h int32) []uint8 {
	result := make([]uint8, w*h*4)
	for y := int32(0); y < h; y += 1 {
		for x := int32(0); x < w; x += 1 {
			result[4*(y*w+x)+0] = *(*uint8)(unsafe.Add(data, 4*((h-1-y)*w+x)+0))
			result[4*(y*w+x)+1] = *(*uint8)(unsafe.Add(data, 4*((h-1-y)*w+x)+1))
			result[4*(y*w+x)+2] = *(*uint8)(unsafe.Add(data, 4*((h-1-y)*w+x)+2))
			result[4*(y*w+x)+3] = *(*uint8)(unsafe.Add(data, 4*((h-1-y)*w+x)+3))
		}
	}
	return result
}

func loadTexture(filepath string) (uint32, error) {
	surf, err := img.Load(filepath)
	if err != nil {
		return 0, errors.Wrap(err, "failed to load texture image")
	}
	defer surf.Free()

	var texId uint32
	gl.GenTextures(1, &texId)

	reversed := gl.Ptr(reverseTexture(surf.Data(), surf.W, surf.H))

	gl.BindTexture(gl.TEXTURE_2D, texId)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, surf.W, surf.H, 0, gl.RGBA, gl.UNSIGNED_BYTE, reversed)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	return texId, nil
}
