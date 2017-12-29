package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

}

func run() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("couldn't init %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not init ttf %v", err)
	}

	defer ttf.Quit()

	flags := uint32(sdl.WINDOW_SHOWN | sdl.WINDOW_ALWAYS_ON_TOP)
	window, renderer, err := sdl.CreateWindowAndRenderer(800, 600, flags)
	if err != nil {
		return fmt.Errorf("couldn't create window and renderer %v", err)
	}
	defer window.Destroy()

	if err := drawTitle(renderer); err != nil {
		return fmt.Errorf("couldn't drawTitle %v", err)
	}
	time.Sleep(5 * time.Second)

	if err := drawBackground(renderer); err != nil {
		return fmt.Errorf("couldn't drawBackground %v", err)
	}
	time.Sleep(5 * time.Second)
	return nil

}

func drawBackground(r *sdl.Renderer) error {
	r.Clear()
	texture, err := img.LoadTexture(r, "resources/images/background.png")
	if err != nil {
		return fmt.Errorf("can't load texture %v", err)
	}
	defer texture.Destroy()

	if err := r.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("can't copy texture to the current rendering target%v", err)
	}
	r.Present()
	return nil
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()
	font, err := ttf.OpenFont("resources/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("can't open font %v", err)
	}
	defer font.Close()

	surface, err := font.RenderUTF8Solid("Flappy", sdl.Color{R: 255, G: 100, B: 0, A: 255})
	if err != nil {
		return fmt.Errorf("can't render font %v", err)
	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("can't create texture %v", err)
	}
	defer texture.Destroy()

	if err := r.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("can't copy texture to the current rendering target%v", err)
	}
	r.Present()
	return nil
}
