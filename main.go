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

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	glContext, err := window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	defer sdl.GLDeleteContext(glContext)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	shader, err := loadShader()
	if err != nil {
		panic(err)
	}

	vao, err := initVAO()
	if err != nil {
		panic(err)
	}

	gl.UseProgram(shader)
	pvmUniformLoc := gl.GetUniformLocation(shader, gl.Str("projModel"+"\x00"))
	colorUniformLoc := gl.GetUniformLocation(shader, gl.Str("color"+"\x00"))

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1)

	scene, err := loadScene("./data/scene_test.json")
	if err != nil {
		panic(err)
	}

	renderer := RendererImpl{
		vao:             vao,
		program:         shader,
		pvmUniformLoc:   pvmUniformLoc,
		colorUniformLoc: colorUniformLoc,
	}

	space := false
	left, right := false, false

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
				if event.Type == sdl.KEYDOWN {
					switch event.Keysym.Scancode {
					case sdl.SCANCODE_Q:
						running = false
						break
					case sdl.SCANCODE_LEFT:
						left = true
						break
					case sdl.SCANCODE_RIGHT:
						right = true
						break
					case sdl.SCANCODE_SPACE:
						space = true
						break
					}
				} else if event.Type == sdl.KEYUP {
					switch event.Keysym.Scancode {
					case sdl.SCANCODE_LEFT:
						left = false
					case sdl.SCANCODE_RIGHT:
						right = false
					case sdl.SCANCODE_SPACE:
						space = false
					}
				}
			}
		}

		for _, block := range scene.Blocks {
			if collides, _ := CheckCollision(player, block); collides {
				ResolveCollision(player, block)
			}
		}

		force := player.Force()
		force = force.Add(mgl32.Vec2{0, -0.098})

		if player.TouchingFloor() {
			if space {
				force = force.Add(mgl32.Vec2{0, 5})
			} else if left && right {
				// do nothing
			} else if left {
				force[0] -= 0.1
				force[0] = Max(force[0], -2)
			} else if right {
				force[0] += 0.1
				force[0] = Min(force[0], 2)
			} else if !left && !right {
				if force[0] > 0 {
					force[0] = Max(force[0]-0.08, 0)
				} else {
					force[0] = Min(force[0]+0.08, 0)
				}
			}
		}

		player.prevPos = player.pos
		player.force = force
		player.pos = player.pos.Add(player.force)

		player.ResetTouch()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		player.Render(renderer)
		scene.RenderBlocks(renderer)

		window.GLSwap()
	}
}
