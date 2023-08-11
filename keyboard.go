package main

import "github.com/veandco/go-sdl2/sdl"

type Key struct {
	Press     bool
	Timestamp uint32
}

func (k *Key) Update(event *sdl.KeyboardEvent) {
	k.Press = event.Type == sdl.KEYDOWN
	k.Timestamp = event.Timestamp
}

type Keyboard [512]Key
