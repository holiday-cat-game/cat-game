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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed cat_plain.png
var catPlainPNG []byte

//go:embed cat_witch_black.png
var catWitchBlackPNG []byte

//go:embed cat_witch_purple.png
var catWitchPurplePNG []byte

//go:embed cat_pumpkin_hat.png
var catPumpkinHatPNG []byte

//go:embed pumpkin.png
var pumpkinPNG []byte

type Game struct {
	x           float32
	y           float32
	outfits     []*ebiten.Image
	outfitNames []string
	outfit      int
	fishX       float32
	fishY       float32
	pumpkin     *ebiten.Image
	score       int
	timeLeft    int
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.outfit = (g.outfit + 1) % len(g.outfits)
	}

	if g.timeLeft <= 0 {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.score = 0
			g.timeLeft = 30 * 60
			g.x = 150
			g.y = 110
		}
		return nil
	}
	g.timeLeft--

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
	if g.y < 6 {
		g.y = 6
	}
	if g.y > 220 {
		g.y = 220
	}

	fishCenterX := g.fishX + 6
	fishCenterY := g.fishY + 6

	if fishCenterX >= g.x && fishCenterX <= g.x+20 &&
		fishCenterY >= g.y && fishCenterY <= g.y+20 {
		g.score++
		g.fishX = rand.Float32() * 308
		g.fishY = rand.Float32() * 228
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{110, 70, 150, 255})

	pop := &ebiten.DrawImageOptions{}
	pop.GeoM.Translate(float64(g.fishX), float64(g.fishY))
	screen.DrawImage(g.pumpkin, pop)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.x), float64(g.y-6))
	screen.DrawImage(g.outfits[g.outfit], op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Pumpkins: %d   Time: %d   Outfit: %s (C)",
		g.score, g.timeLeft/60, g.outfitNames[g.outfit]))

	if g.timeLeft <= 0 {
		vector.DrawFilledRect(screen, 80, 90, 160, 60, color.RGBA{30, 15, 45, 230}, false)
		ebitenutil.DebugPrintAt(screen, "GAME OVER", 133, 96)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("You got %d pumpkins!", g.score), 100, 112)
		ebitenutil.DebugPrintAt(screen, "Press R to play again", 97, 128)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func loadImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Halloween Cat")
	game := &Game{
		x: 150,
		y: 110,
		outfits: []*ebiten.Image{
			loadImage(catPlainPNG),
			loadImage(catWitchBlackPNG),
			loadImage(catWitchPurplePNG),
			loadImage(catPumpkinHatPNG),
		},
		outfitNames: []string{"Cat", "Witch", "Purple Witch", "Pumpkin"},
		fishX:       40,
		fishY:       40,
		pumpkin:     loadImage(pumpkinPNG),
		timeLeft:    30 * 60,
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
