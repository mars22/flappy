package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg *sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bgTexture, err := img.LoadTexture(r, "resources/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("can't load texture %v", err)
	}
	scene := &scene{bg: bgTexture}
	return scene, nil
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("can't copy texture to the current rendering target%v", err)
	}
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
}
