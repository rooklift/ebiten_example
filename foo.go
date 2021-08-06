package main

import (
	"fmt"
	"io/ioutil"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/rooklift/ebiten_example/wavreaderseeker"
)

const (
	w = 600
	h = 400
)

type Game struct{
	image *ebiten.Image

	inited bool
	width int
	height int
	px int
	py int
	speedx int
	speedy int
}

func (g *Game) DrawSprite(x int, y int, img *ebiten.Image) {

	e_width, e_height := img.Size()

	opts := new(ebiten.DrawImageOptions)
	opts.GeoM.Translate(float64(x) - (float64(e_width) / 2), float64(y) - (float64(e_height) / 2.0))

	g.image.DrawImage(img, opts)
}

func (g *Game) Update() error {

	if (!g.inited) {

		g.width = w
		g.height = h
		g.image = ebiten.NewImage(g.width, g.height)

		g.speedx = 2
		g.speedy = 1

		g.inited = true
	}

	if (g.px < 0) { g.speedx = 2 }
	if (g.px >= g.width) { g.speedx = -2 }
	if (g.py < 0) { g.speedy = 1 }
	if (g.py >= g.height) { g.speedy = -1 }

	g.px += g.speedx
	g.py += g.speedy

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.image.Clear()
	g.DrawSprite(g.px, g.py, sprites["powerup.png"])

	screen.DrawImage(g.image, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f -- FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return w, h
}

// ------------------------------------------------------------------------------------------------

var sprites map[string]*ebiten.Image
var sounds map[string]*wavreaderseeker.WavReaderSeeker

func init() {

	sprites = make(map[string]*ebiten.Image)
	sounds = make(map[string]*wavreaderseeker.WavReaderSeeker)

	files, err := ioutil.ReadDir("./sprites")
    if err != nil {
        panic(err)
    }

	for _, info := range files {
		f, err := os.Open("./sprites/" + info.Name())
		if err != nil {
			fmt.Println(err)
		} else {
			img, err := png.Decode(f)
			if err != nil {
				fmt.Println(err)
			} else {
				sprites[info.Name()] = ebiten.NewImageFromImage(img)
			}
			f.Close()
		}
    }
}

func main() {

	g := new(Game)

	ebiten.SetWindowSize(w * 2, h * 2)
	ebiten.SetWindowTitle("Foo")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
