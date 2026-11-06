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
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
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

//go:embed music.wav
var musicWAV []byte

//go:embed cat_santa.png
var catSantaPNG []byte

//go:embed cat_reindeer.png
var catReindeerPNG []byte

//go:embed cat_elf.png
var catElfPNG []byte

//go:embed cat_beanie.png
var catBeaniePNG []byte

//go:embed bg_village.png
var bgVillagePNG []byte

//go:embed bg_tree.png
var bgTreePNG []byte

//go:embed bg_forest_snow.png
var bgForestSnowPNG []byte

//go:embed gift.png
var giftPNG []byte

//go:embed cookie.png
var cookiePNG []byte

//go:embed snowman.png
var snowmanPNG []byte

//go:embed music_xmas.wav
var musicXmasWAV []byte

//go:embed sfx_pickup.wav
var sfxPickupWAV []byte

//go:embed sfx_ghost.wav
var sfxGhostWAV []byte

//go:embed sfx_gameover.wav
var sfxGameOverWAV []byte

const sampleRate = 44100

type Theme struct {
	title       string
	outfits     []*ebiten.Image
	outfitNames []string
	bgs         []*ebiten.Image
	bgNames     []string
	targets     []*ebiten.Image
	targetLabel string
	enemy       *ebiten.Image
	music       *audio.Player
}

type Game struct {
	state       string
	themes      []Theme
	theme       int
	outfit      int
	bg          int
	diff        int
	speeds      []float32
	diffNames   []string
	x           float32
	y           float32
	fishX       float32
	fishY       float32
	target      int
	ghostX      float32
	ghostY      float32
	score       int
	timeLeft    int
	sfxPickup   *audio.Player
	sfxGhost    *audio.Player
	sfxGameOver *audio.Player
	waitSound   bool
}

func (g *Game) Update() error {
	t := g.themes[g.theme]

	if g.state == "menu" {
		if inpututil.IsKeyJustPressed(ebiten.KeyT) {
			t.music.Pause()
			g.theme = (g.theme + 1) % len(g.themes)
			g.outfit = 0
			g.bg = 0
			playSound(g.themes[g.theme].music)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			g.outfit = (g.outfit + 1) % len(t.outfits)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			g.bg = (g.bg + 1) % len(t.bgs)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.diff = (g.diff + 1) % len(g.speeds)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.score = 0
			g.timeLeft = 30 * 60
			g.x = 150
			g.y = 110
			g.ghostX = 290
			g.ghostY = 20
			g.target = rand.Intn(len(t.targets))
			g.waitSound = false
			g.state = "playing"
		}
		return nil
	}

	if g.state == "gameover" {
		if g.waitSound && !g.sfxGhost.IsPlaying() {
			playSound(g.sfxGameOver)
			g.waitSound = false
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = "menu"
		}
		return nil
	}

	g.timeLeft--
	if g.timeLeft <= 0 {
		g.state = "gameover"
		playSound(g.sfxGameOver)
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
		playSound(g.sfxPickup)
		g.fishX = rand.Float32() * 308
		g.fishY = rand.Float32() * 228
		g.target = rand.Intn(len(t.targets))
	}

	speed := g.speeds[g.diff]

	if g.ghostX < g.x {
		g.ghostX += speed
	}
	if g.ghostX > g.x {
		g.ghostX -= speed
	}
	if g.ghostY < g.y {
		g.ghostY += speed
	}
	if g.ghostY > g.y {
		g.ghostY -= speed
	}

	ghostCenterX := g.ghostX + 12
	ghostCenterY := g.ghostY + 12

	if ghostCenterX >= g.x && ghostCenterX <= g.x+20 &&
		ghostCenterY >= g.y && ghostCenterY <= g.y+20 {
		g.state = "gameover"
		playSound(g.sfxGhost)
		g.waitSound = true
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
	t := g.themes[g.theme]
	screen.DrawImage(t.bgs[g.bg], &ebiten.DrawImageOptions{})

	if g.state == "menu" {
		lineColor := color.RGBA{150, 130, 180, 255}
		vector.DrawFilledRect(screen, 40, 24, 240, 182, color.RGBA{30, 15, 45, 215}, false)

		drawCentered(screen, "< "+t.title+" >", 26)
		drawCentered(screen, "Press T to change theme", 38)
		vector.DrawFilledRect(screen, 52, 58, 216, 1, lineColor, false)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(62, 90)
		screen.DrawImage(t.outfits[g.outfit], op)

		drawCenteredAt(screen, "< "+t.outfitNames[g.outfit]+" >", 190, 62)
		drawCenteredAt(screen, "Press C to change", 190, 74)
		vector.DrawFilledRect(screen, 130, 94, 120, 1, lineColor, false)

		drawCenteredAt(screen, "< "+t.bgNames[g.bg]+" >", 190, 98)
		drawCenteredAt(screen, "Press B to change", 190, 110)
		vector.DrawFilledRect(screen, 130, 130, 120, 1, lineColor, false)

		drawCenteredAt(screen, "< "+g.diffNames[g.diff]+" >", 190, 134)
		drawCenteredAt(screen, "Press D to change", 190, 146)

		vector.DrawFilledRect(screen, 52, 170, 216, 1, lineColor, false)
		drawCentered(screen, "Press SPACE to start", 178)
		return
	}

	pop := &ebiten.DrawImageOptions{}
	pop.GeoM.Translate(float64(g.fishX), float64(g.fishY))
	screen.DrawImage(t.targets[g.target], pop)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.x), float64(g.y-6))
	screen.DrawImage(t.outfits[g.outfit], op)

	gop := &ebiten.DrawImageOptions{}
	gop.GeoM.Translate(float64(g.ghostX), float64(g.ghostY))
	screen.DrawImage(t.enemy, gop)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%s: %d   Time: %d", t.targetLabel, g.score, g.timeLeft/60))

	if g.state == "gameover" {
		vector.DrawFilledRect(screen, 80, 90, 160, 60, color.RGBA{30, 15, 45, 230}, false)
		drawCentered(screen, "GAME OVER", 96)
		drawCentered(screen, fmt.Sprintf("You got %d %s!", g.score, strings.ToLower(t.targetLabel)), 112)
		drawCentered(screen, "Press SPACE for menu", 128)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func loadSound(ctx *audio.Context, data []byte) *audio.Player {
	stream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	player, err := ctx.NewPlayer(stream)
	if err != nil {
		log.Fatal(err)
	}
	return player
}

func loadMusic(ctx *audio.Context, data []byte) *audio.Player {
	stream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	player, err := ctx.NewPlayer(loop)
	if err != nil {
		log.Fatal(err)
	}
	player.SetVolume(0.5)
	return player
}

func playSound(p *audio.Player) {
	p.Rewind()
	p.Play()
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
	ebiten.SetWindowTitle("Holiday Cat")
	audioCtx := audio.NewContext(sampleRate)

	halloween := Theme{
		title: "HALLOWEEN CAT",
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
		bgNames:     []string{"Graveyard", "Haunted House", "Spooky Town"},
		targets:     []*ebiten.Image{loadImage(pumpkinPNG)},
		targetLabel: "Pumpkins",
		enemy:       loadImage(ghostPNG),
		music:       loadMusic(audioCtx, musicWAV),
	}

	christmas := Theme{
		title: "CHRISTMAS CAT",
		outfits: []*ebiten.Image{
			loadImage(catPlainPNG),
			loadImage(catSantaPNG),
			loadImage(catReindeerPNG),
			loadImage(catElfPNG),
			loadImage(catBeaniePNG),
		},
		outfitNames: []string{"Classic", "Santa", "Reindeer", "Elf", "Beanie"},
		bgs: []*ebiten.Image{
			loadImage(bgVillagePNG),
			loadImage(bgTreePNG),
			loadImage(bgForestSnowPNG),
		},
		bgNames:     []string{"Snowy Village", "Christmas Tree", "Snowy Forest"},
		targets:     []*ebiten.Image{loadImage(giftPNG), loadImage(cookiePNG)},
		targetLabel: "Treats",
		enemy:       loadImage(snowmanPNG),
		music:       loadMusic(audioCtx, musicXmasWAV),
	}

	game := &Game{
		state:       "menu",
		themes:      []Theme{halloween, christmas},
		speeds:      []float32{0.5, 0.8, 1.0},
		diffNames:   []string{"Easy", "Medium", "Hard"},
		diff:        1,
		fishX:       40,
		fishY:       40,
		sfxPickup:   loadSound(audioCtx, sfxPickupWAV),
		sfxGhost:    loadSound(audioCtx, sfxGhostWAV),
		sfxGameOver: loadSound(audioCtx, sfxGameOverWAV),
	}
	game.themes[0].music.Play()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
