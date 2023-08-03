package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

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

type TouchImpl struct {
	floor, ceiling, left, right bool
}

func (f *TouchImpl) ResetTouch() {
	f.floor = false
	f.ceiling = false
	f.left = false
	f.right = false
}

func (f *TouchImpl) TouchingFloor() bool {
	return f.floor
}

func (f *TouchImpl) TouchingCeiling() bool {
	return f.ceiling
}

func (f *TouchImpl) TouchingLeft() bool {
	return f.left
}

func (f *TouchImpl) TouchingRight() bool {
	return f.right
}

func (f *TouchImpl) TouchFloor() {
	f.floor = true
}

func (f *TouchImpl) TouchCeiling() {
	f.ceiling = true
}

func (f *TouchImpl) TouchLeft() {
	f.left = true
}

func (f *TouchImpl) TouchRight() {
	f.right = true
}

type Touch interface {
	ResetTouch()

	TouchingFloor() bool
	TouchingCeiling() bool
	TouchingLeft() bool
	TouchingRight() bool

	TouchFloor()
	TouchCeiling()
	TouchLeft()
	TouchRight()
}

type DynamicObject interface {
	Object
	Touch

	SetPosition(mgl32.Vec2)
	Force() mgl32.Vec2
	SetForce(mgl32.Vec2)
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
func CheckCollision(a, b Object) (bool, mgl32.Vec2) {
	aPos, aSize := a.Position(), a.Size()
	bPos, bSize := b.Position(), b.Size()

	xyes, x := Overlap1D(aPos.X(), aPos.X()+aSize.X(), bPos.X(), bPos.X()+bSize.X())
	yyes, y := Overlap1D(aPos.Y(), aPos.Y()+aSize.Y(), bPos.Y(), bPos.Y()+bSize.Y())

	return xyes && yyes, mgl32.Vec2{x, y}
}

// ResolveCollision resolves collision by moving object a
func ResolveCollision(a DynamicObject, b Object) {
	aPrevPos, bPrevPos := a.PrevPosition(), b.PrevPosition()
	aPos, bPos := a.Position(), b.Position()
	aSize, bSize := a.Size(), b.Size()

	if aPrevPos.Y() >= bPrevPos.Y()+bSize.Y() {
		// a was above b
		a.SetPosition(mgl32.Vec2{aPos.X(), bPos.Y() + bSize.Y()})
		a.SetForce(mgl32.Vec2{a.Force().X(), 0})
		a.TouchFloor()
		return
	} else if aPrevPos.Y()+aSize.Y() <= bPrevPos.Y() {
		// a was below b
		a.SetPosition(mgl32.Vec2{aPos.X(), bPos.Y() - aSize.Y()})
		a.SetForce(mgl32.Vec2{a.Force().X(), 0})
		a.TouchCeiling()
		return
	} else if aPrevPos.X()+aSize.X() <= bPrevPos.X() {
		// a was to the left of B
		a.SetPosition(mgl32.Vec2{bPos.X() - aSize.X(), aPos.Y()})
		a.SetForce(mgl32.Vec2{0, a.Force().Y()})
		a.TouchLeft()
		return
	} else if aPrevPos.X() >= bPrevPos.X()+bSize.X() {
		// a was to the right of B
		a.SetPosition(mgl32.Vec2{bPos.X() + bSize.Y(), aPos.Y()})
		a.SetForce(mgl32.Vec2{0, a.Force().Y()})
		a.TouchRight()
		return
	}
}

func normalizeUint8(a uint8) float32 {
	return float32(a) / 256
}

func normalizeColor(r, g, b uint8) (float32, float32, float32) {
	return normalizeUint8(r), normalizeUint8(g), normalizeUint8(b)
}
