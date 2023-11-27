package lib

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	window *sdl.Window
}

// Init inicializuje okno s danou šířkou, výškou a názvem.
func Init(title string, width, height int32) (*Window, error) {
	err := sdl.Init(uint32(sdl.INIT_EVERYTHING))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SDL: %v", err)
	}

	window, err := sdl.CreateWindow(title, int32(sdl.WINDOWPOS_UNDEFINED), int32(sdl.WINDOWPOS_UNDEFINED), width, height, uint32(sdl.WINDOW_SHOWN))
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %v", err)
	}

	return &Window{window: window}, nil
}

// DrawImage vykreslí daný obrázek na okno.
func (w *Window) DrawImage(imagePath string) error {
	surface, err := sdl.LoadBMP(imagePath)
	if err != nil {
		return fmt.Errorf("failed to load image: %v", err)
	}
	defer surface.Free()

	renderer, err := w.window.GetRenderer()
	if err != nil {
		return fmt.Errorf("failed to get renderer: %v", err)
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("failed to create texture: %v", err)
	}
	defer texture.Destroy()

	renderer2, err := w.window.GetRenderer()
	if err != nil {
		log.Fatalf("Failed to get renderer: %v", err)
	}
	renderer2.Present()

	return nil
}

// Quit ukončí SDL a okno.
func (w *Window) Quit() {
	w.window.Destroy()
	sdl.Quit()
}
