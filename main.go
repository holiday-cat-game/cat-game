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

//go:embed bg_graveyard.png
var bgGraveyardPNG []byte

//go:embed bg_house.png
var bgHousePNG []byte

//go:embed bg_combined.png
var bgCombinedPNG []byte

//go:embed ghost.png
var ghostPNG []byte

type Game struct {
	state       string
	x           float32
	y           float32
	outfits     []*ebiten.Image
	outfitNames []string
	outfit      int
	bgs         []*ebiten.Image
	bgNames     []string
	bg          int
	fishX       float32
	fishY       float32
	pumpkin     *ebiten.Image
	score       int
	timeLeft    int
	ghost       *ebiten.Image
	ghostX      float32
	ghostY      float32
}

func (g *Game) Update() error {
	if g.state == "menu" {
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			g.outfit = (g.outfit + 1) % len(g.outfits)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			g.bg = (g.bg + 1) % len(g.bgs)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.score = 0
			g.timeLeft = 30 * 60
			g.x = 150
			g.y = 110
			g.ghostX = 290
			g.ghostY = 20
			g.state = "playing"
		}
		return nil
	}

	if g.state == "gameover" {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = "menu"
		}
		return nil
	}

	g.timeLeft--
	if g.timeLeft <= 0 {
		g.state = "gameover"
		return nil
	}

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

	if g.ghostX < g.x {
		g.ghostX += 0.8
	}
	if g.ghostX > g.x {
		g.ghostX -= 0.8
	}
	if g.ghostY < g.y {
		g.ghostY += 0.8
	}
	if g.ghostY > g.y {
		g.ghostY -= 0.8
	}

	ghostCenterX := g.ghostX + 12
	ghostCenterY := g.ghostY + 12

	if ghostCenterX >= g.x && ghostCenterX <= g.x+20 &&
		ghostCenterY >= g.y && ghostCenterY <= g.y+20 {
		g.state = "gameover"
	}

	return nil
}

func drawCenteredAt(screen *ebiten.Image, text string, centerX, y int) {
	x := centerX - len(text)*6/2
	ebitenutil.DebugPrintAt(screen, text, x, y)
}

func drawCentered(screen *ebiten.Image, text string, y int) {
	drawCenteredAt(screen, text, 160, y)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.bgs[g.bg], &ebiten.DrawImageOptions{})

	if g.state == "menu" {
		lineColor := color.RGBA{150, 130, 180, 255}
		vector.DrawFilledRect(screen, 40, 34, 240, 150, color.RGBA{30, 15, 45, 215}, false)

		drawCentered(screen, "HALLOWEEN CAT", 38)
		vector.DrawFilledRect(screen, 52, 58, 216, 1, lineColor, false)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(62, 74)
		screen.DrawImage(g.outfits[g.outfit], op)

		drawCenteredAt(screen, "< "+g.outfitNames[g.outfit]+" >", 190, 68)
		drawCenteredAt(screen, "Press C to change", 190, 82)
		vector.DrawFilledRect(screen, 130, 103, 120, 1, lineColor, false)
		drawCenteredAt(screen, "< "+g.bgNames[g.bg]+" >", 190, 110)
		drawCenteredAt(screen, "Press B to change", 190, 124)

		vector.DrawFilledRect(screen, 52, 150, 216, 1, lineColor, false)
		drawCentered(screen, "Press SPACE to start", 160)
		return
	}

	pop := &ebiten.DrawImageOptions{}
	pop.GeoM.Translate(float64(g.fishX), float64(g.fishY))
	screen.DrawImage(g.pumpkin, pop)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.x), float64(g.y-6))
	screen.DrawImage(g.outfits[g.outfit], op)

	gop := &ebiten.DrawImageOptions{}
	gop.GeoM.Translate(float64(g.ghostX), float64(g.ghostY))
	screen.DrawImage(g.ghost, gop)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Pumpkins: %d   Time: %d", g.score, g.timeLeft/60))

	if g.state == "gameover" {
		vector.DrawFilledRect(screen, 80, 90, 160, 60, color.RGBA{30, 15, 45, 230}, false)
		drawCentered(screen, "GAME OVER", 96)
		drawCentered(screen, fmt.Sprintf("You got %d pumpkins!", g.score), 112)
		drawCentered(screen, "Press SPACE for menu", 128)
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
		state: "menu",
		outfits: []*ebiten.Image{
			loadImage(catPlainPNG),
			loadImage(catWitchBlackPNG),
			loadImage(catWitchPurplePNG),
			loadImage(catPumpkinHatPNG),
		},
		outfitNames: []string{"Classic", "Witch", "Purple Witch", "Pumpkin"},
		bgs: []*ebiten.Image{
			loadImage(bgGraveyardPNG),
			loadImage(bgHousePNG),
			loadImage(bgCombinedPNG),
		},
		bgNames: []string{"Graveyard", "Haunted House", "Spooky Town"},
		fishX:   40,
		fishY:   40,
		pumpkin: loadImage(pumpkinPNG),
		ghost:   loadImage(ghostPNG),
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
