package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed cat.png
var catPNG []byte

type Game struct {
	x   float32
	y   float32
	cat *ebiten.Image
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
	screen.Fill(color.RGBA{200, 230, 250, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.x), float64(g.y))
	screen.DrawImage(g.cat, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(catPNG))
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Cat")
	game := &Game{x: 150, y: 110, cat: ebiten.NewImageFromImage(img)}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
