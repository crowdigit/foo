package main

const (
	OFFSET_Z_BLOCK  = 0
	OFFSET_Z_PLAYER = 1
)

type Object interface {
	Render(Renderer)
}

func normalizeUint8(a uint8) float32 {
	return float32(a) / 256
}

func normalizeColor(r, g, b uint8) (float32, float32, float32) {
	return normalizeUint8(r), normalizeUint8(g), normalizeUint8(b)
}
