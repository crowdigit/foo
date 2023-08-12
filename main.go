package main

import (
	"github.com/crowdigit/foo/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 768
	SCREEN_DEPTH  = 1
)

type RenderParameter struct {
	screenWidth  int
	screenHeight int
	screenDepth  int
	vao          uint32
}

func (p RenderParameter) Proj() mgl32.Mat4 {
	return mgl32.Mat4{
		2.0 / float32(p.screenWidth), 0, 0, 0,
		0, 2.0 / float32(p.screenHeight), 0, 0,
		0, 0, float32(p.screenDepth), 0,
		-1, -1, -1, 1,
	}
}

func (p RenderParameter) VAO() uint32 {
	return p.vao
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err := img.Init(img.INIT_PNG); err != nil {
		panic(err)
	}
	defer img.Quit()

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

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1)

	vao, err := initVAO()
	if err != nil {
		panic(err)
	}

	renderParameter := RenderParameter{
		screenWidth:  SCREEN_WIDTH,
		screenHeight: SCREEN_HEIGHT,
		screenDepth:  SCREEN_DEPTH,
		vao:          vao,
	}

	colorProgram, err := shader.LoadColorShader(renderParameter)
	if err != nil {
		panic(err)
	}

	textureProgram, err := shader.LoadTextureShader(renderParameter)
	if err != nil {
		panic(err)
	}

	renderer := RendererImpl{
		ColorProgram:   colorProgram,
		TextureProgram: textureProgram,
	}

	texID, err := loadTexture("./data/player.png")
	if err != nil {
		panic(err)
	}
	defer gl.DeleteTextures(1, &texID)

	scene, err := loadScene("./data/scene_test.json")
	if err != nil {
		panic(err)
	}

	keyboard := Keyboard{}

	player := &Player{
		Touch: &TouchImpl{},
		pos:   mgl32.Vec2{0, 60},
		size:  mgl32.Vec2{20, 20},
		texID: texID,
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
