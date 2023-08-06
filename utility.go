package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pkg/errors"
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
