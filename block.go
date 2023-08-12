package main

import (
	"encoding/json"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
)

type Block struct {
	pos, size mgl32.Vec2
}

// Position implements main.Object
func (b Block) Position() mgl32.Vec2 {
	return b.pos
}

// Position implements main.Object
func (b Block) PrevPosition() mgl32.Vec2 {
	return b.pos
}

// Size implements main.Object
func (b Block) Size() mgl32.Vec2 {
	return b.size
}

// Render implements main.Object
func (b Block) Render(renderer Renderer) {
	pos := mgl32.Vec3{b.pos[0], b.pos[1], OFFSET_Z_BLOCK}
	renderer.RenderColoredRect(pos, b.size, 255, 0, 0)
}

// UnmarshalJSON implements json.Unmarshaler
func (b *Block) UnmarshalJSON(bytes []byte) error {
	arr := make([]float32, 0, 4)
	if err := json.Unmarshal(bytes, &arr); err != nil {
		return errors.Wrap(err, "failed to unmarshal bytes into block")
	} else if len(arr) != 4 {
		return errors.New("block definition must be list with length 4")
	}

	b.pos = mgl32.Vec2{arr[0], arr[1]}
	b.size = mgl32.Vec2{arr[2], arr[3]}

	return nil
}
