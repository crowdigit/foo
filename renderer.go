package main

import (
	"github.com/crowdigit/foo/shader"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer interface {
	RenderColoredRect(pos mgl32.Vec3, size mgl32.Vec2, r, g, b uint8)
	RenderTexturedRect(pos mgl32.Vec3, size mgl32.Vec2, texID uint32)
}

type RendererImpl struct {
	shader.ColorProgram
	shader.TextureProgram
}

func (i RendererImpl) RenderColoredRect(pos mgl32.Vec3, size mgl32.Vec2, r, g, b uint8) {
	i.ColorProgram.Render(pos, size, r, g, b)
}

func (i RendererImpl) RenderTexturedRect(pos mgl32.Vec3, size mgl32.Vec2, texID uint32) {
	i.TextureProgram.Render(pos, size, texID)
}
