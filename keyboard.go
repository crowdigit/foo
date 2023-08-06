package main

import "github.com/veandco/go-sdl2/sdl"

type Keyboard struct {
	Q bool

	Left, Right bool
	Space       bool
}

func (k *Keyboard) Update(event *sdl.KeyboardEvent) {
	if event.Type == sdl.KEYDOWN {
		switch event.Keysym.Scancode {
		case sdl.SCANCODE_Q:
			k.Q = true
			break
		case sdl.SCANCODE_LEFT:
			k.Left = true
			break
		case sdl.SCANCODE_RIGHT:
			k.Right = true
			break
		case sdl.SCANCODE_SPACE:
			k.Space = true
			break
		}
	} else if event.Type == sdl.KEYUP {
		switch event.Keysym.Scancode {
		case sdl.SCANCODE_LEFT:
			k.Left = false
		case sdl.SCANCODE_RIGHT:
			k.Right = false
		case sdl.SCANCODE_SPACE:
			k.Space = false
		}
	}
}
