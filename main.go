package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 768
	SCREEN_DEPTH  = 1
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	configureOpenGL()

	glContext, err := window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	defer sdl.GLDeleteContext(glContext)

	renderer, err := initRenderer()
	if err != nil {
		panic(err)
	}

	scene, err := loadScene("./data/scene_test.json")
	if err != nil {
		panic(err)
	}

	keyboard := Keyboard{}

	player := &Player{
		Touch: &TouchImpl{},
		pos:   mgl32.Vec2{0, 60}, size: mgl32.Vec2{20, 20},
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event := event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.KeyboardEvent:
				keyboard[event.Keysym.Scancode].Update(event)
				if keyboard[sdl.SCANCODE_Q].Press {
					running = false
				}
				break
			}
		}

		player.Update(scene, keyboard)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		player.Render(renderer)
		scene.RenderBlocks(renderer)

		window.GLSwap()
	}
}
