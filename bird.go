package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity   = 0.2
	jumpSpeed = -5
)

type bird struct {
	time     int
	textures []*sdl.Texture
	y, speed float64
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i < 5; i++ {
		file := fmt.Sprintf("resources/images/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(r, file)
		if err != nil {
			return nil, fmt.Errorf("can't load texture %v", err)
		}
		textures = append(textures, texture)
	}
	return &bird{textures: textures, y: 300, speed: 0}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.time++
	b.y -= b.speed
	b.speed += gravity
	if b.y < 0 {
		b.y = 0
		b.speed = -b.speed
	}
	maxY := float64(600 - 43/2)
	if b.y >= maxY {
		b.y = maxY
	}

	// we gonna animate birds 10 time slower then rest of the scene
	i := b.time / 10 % len(b.textures)
	rec := sdl.Rect{X: 0, Y: (600 - int32(b.y)) - 43/2, W: 50, H: 43}
	if err := r.Copy(b.textures[i], nil, &rec); err != nil {
		return fmt.Errorf("can't copy texture to the current rendering target %v", err)
	}
	return nil
}

func (b *bird) jump() {
	b.speed = jumpSpeed
}

func (b *bird) destroy() {
	for _, texture := range b.textures {
		texture.Destroy()
	}
}
