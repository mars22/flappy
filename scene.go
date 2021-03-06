package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time  int
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
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

	pipes, err := newPipes(r)
	if err != nil {
		return nil, fmt.Errorf("can't load bird %v", err)
	}
	scene := &scene{bg: bgTexture, bird: bird, pipes: pipes}
	return scene, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()
				if s.bird.isDead() {
					time.Sleep(1 * time.Second)
					if err := drawTitle(r, "Game Over"); err != nil {
						errc <- err
					}
					time.Sleep(time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}

		}
	}()
	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
		return false
	default:
		// log.Printf("unknown event: %T", e)
		return false
	}

}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
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
	if err := s.pipes.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}
