package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	EMPTY  = 0
	PLAYER = 1

	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600

	MAP_WIDTH     = 10
	MAP_HEIGHT    = 10
	MAP_GRID_SIZE = 40

	RENDER_CENTER_OFFSET_X = SCREEN_WIDTH/2 - MAP_GRID_SIZE*MAP_WIDTH/2
	RENDER_CENTER_OFFSET_Y = SCREEN_HEIGHT/2 - MAP_GRID_SIZE*MAP_HEIGHT/2
)

type Object interface {
	Position() (int, int)
	Render(*sdl.Renderer)
}

type Player struct {
	x int
	y int

	r, g, b uint8
}

func (p Player) Position() (int, int) {
	return p.x, p.y
}

func (p Player) Render(renderer *sdl.Renderer) {
	rect := sdl.Rect{
		X: int32(p.x)*MAP_GRID_SIZE + RENDER_CENTER_OFFSET_X,
		Y: int32(-p.y+MAP_HEIGHT-1)*MAP_GRID_SIZE + RENDER_CENTER_OFFSET_Y,
		W: MAP_GRID_SIZE,
		H: MAP_GRID_SIZE,
	}
	renderer.SetDrawColor(p.r, p.g, p.b, 0)
	if err := renderer.FillRect(&rect); err != nil {
		panic(err)
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	player := Player{x: 0, y: 0, r: 255, g: 255, b: 255}
	objects := map[string]Object{
		"player":  &player,
		"monster": Player{x: 9, y: 9, r: 255, g: 0, b: 0},
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

		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		for _, object := range objects {
			object.Render(renderer)
		}

		renderer.SetDrawColor(255, 255, 255, 0)
		renderer.DrawRects(fieldRects)
		renderer.Present()
	}
}
