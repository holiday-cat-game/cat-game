package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	x float32
	y float32
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.x += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.x -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.y += 2
	}

	if g.x < 0 {
		g.x = 0
	}
	if g.x > 300 {
		g.x = 300
	}
	if g.y < 0 {
		g.y = 0
	}
	if g.y > 220 {
		g.y = 220
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	orange := color.RGBA{255, 150, 50, 255}
	vector.DrawFilledRect(screen, g.x, g.y, 20, 20, orange, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Cat")
	game := &Game{x: 150, y: 110}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
