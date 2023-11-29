package lib

// Graphics is a wrapper for raylib graphics library
import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Graphics struct {
	WindowTitle string
	ShouldClose bool
	Image       rl.Texture2D
	KeyLeft     bool
	KeyRight    bool
	KeyUp       bool
	KeyDown     bool
}

// Init initializes the graphics context
// x, y - window size
// title - window title

func Init(x, y int, title string) (*Graphics, error) {
	rl.InitWindow(int32(x), int32(y), title)
	rl.SetTargetFPS(60)

	return &Graphics{WindowTitle: title, ShouldClose: false}, nil
}

// DrawImage draws an image on the graphics context
// imagePath - path to the image
// x, y - position of the image on the screen
func (g *Graphics) DrawImage(x, y int32, imagePath string) {
	img := rl.LoadTexture(imagePath)

	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// Vykreslení obrázku na grafický kontext
	rl.DrawTexture(img, x, y, rl.RayWhite)

	rl.EndDrawing()

}
func (g *Graphics) SetFPS(x int) {
	rl.SetTargetFPS(int32(x))
}

// Update updates the graphics context
func (g *Graphics) Update() {
	g.ShouldClose = rl.WindowShouldClose()
	g.KeyLeft = rl.IsKeyDown(rl.KeyLeft)
	g.KeyRight = rl.IsKeyDown(rl.KeyRight)
	g.KeyUp = rl.IsKeyDown(rl.KeyUp)
	g.KeyDown = rl.IsKeyDown(rl.KeyDown)
}

// CloseWindow closes the graphics context
func (g *Graphics) CloseWindow() {
	rl.UnloadTexture(g.Image)
	rl.CloseWindow()
}
