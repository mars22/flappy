package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipes struct {
	mu      sync.RWMutex
	texture *sdl.Texture
	pipes   []*pipe
	speed   int32
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	var texture *sdl.Texture
	texture, err := img.LoadTexture(r, "resources/images/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("can't load texture %v", err)
	}
	pipes := &pipes{texture: texture, speed: 1}

	go func() {
		for {
			pipes.mu.Lock()
			pipes.pipes = append(pipes.pipes, newPipe())
			pipes.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return pipes, nil
}

func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var visibale []*pipe

	for _, p := range ps.pipes {
		p.update(ps.speed)
		if p.x+p.w > 0 {
			visibale = append(visibale, p)
		}
	}

	ps.pipes = visibale
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}
	return nil
}

func (ps *pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.texture.Destroy()
}

func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if p.touch(b) {
			b.setDead()
			return
		}
	}
}

type pipe struct {
	mu      sync.RWMutex
	x, w, h int32

	inverted bool
}

func newPipe() *pipe {
	return &pipe{
		x:        800 + rand.Int31n(50), //most right side of the screen
		h:        100 + rand.Int31n(200),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}
}

func (p *pipe) update(speed int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.x -= speed

}

func (p *pipe) paint(r *sdl.Renderer, t *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rec := sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rec.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := r.CopyEx(t, nil, &rec, 0, nil, flip); err != nil {
		return fmt.Errorf("can't copy texture to the current rendering target %v", err)
	}
	return nil
}

func (p *pipe) touch(b *bird) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.inverted {
		return b.x+b.w > p.x && (b.y+b.h/3) > 600-p.h
	}

	return b.x+b.w > p.x && 600-(b.y-b.h/3) > 600-p.h

}
