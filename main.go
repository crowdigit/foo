package main

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 768
	SCREEN_DEPTH  = 1
)

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

	renderer := RendererImpl{
		vao:             vao,
		program:         shader,
		pvmUniformLoc:   pvmUniformLoc,
		colorUniformLoc: colorUniformLoc,
	}

	player := &Player{pos: mgl32.Vec2{0, 20}, size: mgl32.Vec2{20, 20}}

	objects := []Object{
		Block{pos: mgl32.Vec2{0, 0}, size: mgl32.Vec2{400, 20}},
		Block{pos: mgl32.Vec2{400, 20}, size: mgl32.Vec2{400, 20}},
		Block{pos: mgl32.Vec2{800, 40}, size: mgl32.Vec2{400, 20}},
		player,
	}

	blocks := make([]Object, 0, 10)
	for _, object := range objects {
		if _, ok := object.(Block); ok {
			blocks = append(blocks, object)
		}
	}

	mov := mgl32.Vec2{0, 0}
	left, right, up, down := false, false, false, false

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
					case sdl.SCANCODE_UP:
						up = true
						break
					case sdl.SCANCODE_DOWN:
						down = true
						break
					}
				} else if event.Type == sdl.KEYUP {
					switch event.Keysym.Scancode {
					case sdl.SCANCODE_LEFT:
						left = false
						break
					case sdl.SCANCODE_RIGHT:
						right = false
						break
					case sdl.SCANCODE_UP:
						up = false
						break
					case sdl.SCANCODE_DOWN:
						down = false
						break
					}
				}
			}
		}

		mov = mgl32.Vec2{0, 0}

		if left {
			mov = mov.Add(mgl32.Vec2{-1, 0})
		}
		if right {
			mov = mov.Add(mgl32.Vec2{1, 0})
		}
		if up {
			mov = mov.Add(mgl32.Vec2{0, 1})
		}
		if down {
			mov = mov.Add(mgl32.Vec2{0, -1})
		}

		player.prevPos = player.pos
		player.pos = player.pos.Add(mov)

		for _, block := range blocks {
			if collides, conflict := CheckCollision(player, block); collides {
				ResolveCollision(player, block, conflict)
			}
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, object := range objects {
			object.Render(renderer)
		}

		window.GLSwap()
	}
}
