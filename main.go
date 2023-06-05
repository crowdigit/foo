package main

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	EMPTY  = 0
	PLAYER = 1

	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600

	MAP_WIDTH           = 10
	MAP_HEIGHT          = 10
	MAP_GRID_SIZE       = 30
	MAP_GRID_LINE_WIDTH = 3

	RENDER_CENTER_OFFSET_X = SCREEN_WIDTH/2 - MAP_GRID_SIZE*MAP_WIDTH/2
	RENDER_CENTER_OFFSET_Y = SCREEN_HEIGHT/2 - MAP_GRID_SIZE*MAP_HEIGHT/2
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

	player := Player{x: 0, y: 0, r: 0, g: 0, b: 255, speed: 1}
	grid := Grid{
		width: MAP_WIDTH, height: MAP_HEIGHT,
		lineWidth: MAP_GRID_LINE_WIDTH,
		r:         255, g: 255, b: 255,
		gridSize: MAP_GRID_SIZE,
	}
	objects := map[string]Object{
		"player":  &player,
		"monster": Player{x: 6, y: 6, r: 255, g: 0, b: 0, speed: 1},
		"grid":    grid,
	}

	// Define field grids
	fieldRects := make([]sdl.Rect, MAP_HEIGHT*MAP_WIDTH)
	for i := 0; i < MAP_HEIGHT; i += 1 {
		for j := 0; j < MAP_WIDTH; j += 1 {
			fieldRects[i*MAP_WIDTH+j] = sdl.Rect{
				X: int32(j)*MAP_GRID_SIZE + RENDER_CENTER_OFFSET_X,
				Y: int32(i)*MAP_GRID_SIZE + RENDER_CENTER_OFFSET_Y,
				W: MAP_GRID_SIZE,
				H: MAP_GRID_SIZE,
			}
		}
	}

	gl.UseProgram(shader)
	pvmUniformLoc := gl.GetUniformLocation(shader, gl.Str("projModel"+"\x00"))
	colorUniformLoc := gl.GetUniformLocation(shader, gl.Str("color"+"\x00"))

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1)

	checkGLError()

	renderer := RectRenderer{
		vao:             vao,
		program:         shader,
		pvmUniformLoc:   pvmUniformLoc,
		colorUniformLoc: colorUniformLoc,
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
					case sdl.SCANCODE_LEFT:
						player.x = Max(player.x-1, 0)
					case sdl.SCANCODE_RIGHT:
						player.x = Min(player.x+1, MAP_WIDTH-1)
					case sdl.SCANCODE_UP:
						player.y = Min(player.y+1, MAP_HEIGHT-1)
					case sdl.SCANCODE_DOWN:
						player.y = Max(player.y-1, 0)
					case sdl.SCANCODE_Q:
						running = false
					}
				}
			}
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, object := range objects {
			object.Render(renderer)
		}

		window.GLSwap()
	}
}
