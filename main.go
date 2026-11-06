package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed cat.png
var catPNG []byte

type Game struct {
	x     float32
	y     float32
	cat   *ebiten.Image
	fishX float32
	fishY float32
	score int
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

	fishCenterX := g.fishX + 4
	fishCenterY := g.fishY + 4

	if fishCenterX >= g.x && fishCenterX <= g.x+20 &&
		fishCenterY >= g.y && fishCenterY <= g.y+20 {
		g.score++
		g.fishX = rand.Float32() * 312
		g.fishY = rand.Float32() * 232
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{110, 170, 90, 255})

	orange := color.RGBA{255, 150, 50, 255}
	vector.DrawFilledRect(screen, g.fishX, g.fishY, 8, 8, orange, false)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.x), float64(g.y))
	screen.DrawImage(g.cat, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))
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
	game := &Game{
		x:     150,
		y:     110,
		cat:   ebiten.NewImageFromImage(img),
		fishX: 40,
		fishY: 40,
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
