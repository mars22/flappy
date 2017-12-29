package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time  int
	bg    *sdl.Texture
	birds []*sdl.Texture
}

func loadTexture(r *sdl.Renderer, file string) (*sdl.Texture, error) {
	texture, err := img.LoadTexture(r, "resources/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("can't load texture %v", err)
	}
	return texture, nil
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bgTexture, err := img.LoadTexture(r, "resources/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("can't load texture %v", err)
	}

	var birds []*sdl.Texture
	for i := 1; i < 5; i++ {
		file := fmt.Sprintf("resources/images/bird_frame_%d.png", i)
		bird, err := img.LoadTexture(r, file)
		if err != nil {
			return nil, fmt.Errorf("can't load texture %v", err)
		}
		birds = append(birds, bird)
	}

	scene := &scene{bg: bgTexture, birds: birds}
	return scene, nil
}

func (s *scene) run(ctx context.Context, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}

		}
	}()
	return errc
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("can't copy texture to the current rendering target%v", err)
	}

	// we gonna animate birds 10 time slower then rest of the scene
	i := s.time / 10 % len(s.birds)
	rec := sdl.Rect{X: 0, Y: 300 - 43/2, W: 50, H: 43}
	if err := r.Copy(s.birds[i], nil, &rec); err != nil {
		return fmt.Errorf("can't copy texture to the current rendering target%v", err)
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	for _, bird := range s.birds {
		bird.Destroy()
	}
}
