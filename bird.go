package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity   = 0.1
	jumpSpeed = -2
)

type bird struct {
	mu         sync.RWMutex
	time       int
	textures   []*sdl.Texture
	x, y, w, h int32
	speed      float64
	dead       bool
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
	return &bird{textures: textures, x: 10, y: 500, h: 43, w: 50, speed: 0}, nil
}

func (b *bird) isBellyTouching(y int32) bool {
	return b.y-b.h/3 < y
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.y -= int32(b.speed)
	b.speed += gravity
	if b.isBellyTouching(0) {
		b.dead = true
	}

	maxY := (600 - b.h/2)
	if b.y >= maxY {
		b.y = maxY
	}
}

func (b *bird) restart() {
	//Lock for write
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time = 0
	b.y = 300
	b.speed = 0
	b.dead = false

}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// we gonna animate birds 10 time slower then rest of the scene
	i := b.time / 10 % len(b.textures)
	rec := sdl.Rect{X: b.x, Y: (600 - b.y) - b.h/2, W: b.w, H: b.h}
	if err := r.Copy(b.textures[i], nil, &rec); err != nil {
		return fmt.Errorf("can't copy texture to the current rendering target %v", err)
	}
	return nil
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, texture := range b.textures {
		texture.Destroy()
	}
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dead
}

func (b *bird) setDead() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.dead = true
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = jumpSpeed
}
