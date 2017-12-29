package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time int
	bg   *sdl.Texture
	bird *bird
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

	bird, err := newBird(r)
	if err != nil {
		return nil, fmt.Errorf("can't load bird %v", err)
	}
	scene := &scene{bg: bgTexture, bird: bird}
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
		return fmt.Errorf("can't copy texture to the current rendering target %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
}
