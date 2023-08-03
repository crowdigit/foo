package main

import "github.com/go-gl/mathgl/mgl32"

const (
	OFFSET_Z_BLOCK  = 0
	OFFSET_Z_PLAYER = 1
)

type Object interface {
	Render(Renderer)

	Position() mgl32.Vec2
	PrevPosition() mgl32.Vec2
	Size() mgl32.Vec2
}

func Overlap1D(amin, amax, bmin, bmax float32) (bool, float32) {
	aw := amax - amin
	bw := bmax - bmin
	min := Min(amin, bmin)
	max := Max(amax, bmax)
	cw := max - min
	return cw < aw+bw, aw + bw - cw
}

// CheckCollision checks if two object collides and returns true and maximum
// penetrating vector
func CheckCollision(a, b Object) (bool, float32, float32) {
	aPos, aSize := a.Position(), a.Size()
	bPos, bSize := b.Position(), b.Size()

	xyes, x := Overlap1D(aPos.X(), aPos.X()+aSize.X(), bPos.X(), bPos.X()+bSize.X())
	yyes, y := Overlap1D(aPos.Y(), aPos.Y()+aSize.Y(), bPos.Y(), bPos.Y()+bSize.Y())

	return xyes && yyes, x, y
}

func normalizeUint8(a uint8) float32 {
	return float32(a) / 256
}

func normalizeColor(r, g, b uint8) (float32, float32, float32) {
	return normalizeUint8(r), normalizeUint8(g), normalizeUint8(b)
}
