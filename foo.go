package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	w = 600
	h = 400
)

// ------------------------------------------------------------------------------------------------

type Game struct{

	inited bool

	width int
	height int
	image *ebiten.Image

	audio_context *audio.Context
	audio_players []*audio.Player

	px int
	py int
	speedx int
	speedy int
	tick int
}

// ------------------------------------------------------------------------------------------------

func (self *Game) DrawSprite(x int, y int, img *ebiten.Image) {

	e_width, e_height := img.Size()

	opts := new(ebiten.DrawImageOptions)
	opts.GeoM.Translate(float64(x) - (float64(e_width) / 2), float64(y) - (float64(e_height) / 2.0))

	self.image.DrawImage(img, opts)
}

func (self *Game) PlaySound(s string) {

	soundbytes, ok := sounds[s]
	if (!ok) {
		fmt.Printf("No such sound: %v\n", s)
		return
	}

	wav_reader := bytes.NewReader(soundbytes)		// wav_reader satisfies io.Reader/Seeker. Relies on the WAV being 16 bit stereo.

	player, err := audio.NewPlayer(self.audio_context, wav_reader)
	if err != nil {
		return
	}

	self.audio_players = append(self.audio_players, player)

	player.Play()
}

// ------------------------------------------------------------------------------------------------

func (self *Game) Init() {

	self.audio_context = audio.NewContext(44100)

	self.width = w
	self.height = h
	self.image = ebiten.NewImage(self.width, self.height)

	self.speedx = 2
	self.speedy = 1

	self.inited = true
}

func (self *Game) PurgeAudio() {

	var active_players []*audio.Player

	for _, player := range self.audio_players {
		if player.IsPlaying() {
			active_players = append(active_players, player)
		} else {
			player.Close()		// Not sure if this is needed.
		}
	}
	self.audio_players = active_players
}

// ------------------------------------------------------------------------------------------------

func (self *Game) GameLogic() {

	self.tick++

	if (self.px < 0) {
		self.speedx = 2
		self.PlaySound("test.wav")
	}
	if (self.px >= self.width) {
		self.speedx = -2
		self.PlaySound("test.wav")
	}
	if (self.py < 0) {
		self.speedy = 1
		self.PlaySound("test.wav")
	}
	if (self.py >= self.height) {
		self.speedy = -1
		self.PlaySound("test.wav")
	}

	self.px += self.speedx
	self.py += self.speedy
}

// ------------------------------------------------------------------------------------------------

func (self *Game) Update() error {

	if (!self.inited) {
		self.Init()
	}

	self.PurgeAudio()

	self.GameLogic()

	return nil
}

func (self *Game) Draw(screen *ebiten.Image) {

	self.image.Clear()
	self.DrawSprite(self.px, self.py, sprites["powerup.png"])

	screen.DrawImage(self.image, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f -- FPS: %0.2f -- Players: %v", ebiten.CurrentTPS(), ebiten.CurrentFPS(), len(self.audio_players)))
}

func (self *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return w, h
}

// ------------------------------------------------------------------------------------------------

var sprites map[string]*ebiten.Image
var sounds map[string][]byte

func main() {

	load_sprites()
	load_sounds()

	g := new(Game)

	ebiten.SetWindowSize(w * 2, h * 2)
	ebiten.SetWindowTitle("Foo")

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}

func load_sprites() {

	sprites = make(map[string]*ebiten.Image)

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

func load_sounds() {

	sounds = make(map[string][]byte)

	files, err := ioutil.ReadDir("./sounds")
    if err != nil {
        panic(err)
    }

	for _, info := range files {
		f, err := os.Open("./sounds/" + info.Name())
		if err != nil {
			fmt.Println(err)
		} else {
			f.Seek(44, io.SeekStart)						// Skip WAV header
			sounds[info.Name()], _ = ioutil.ReadAll(f)
			f.Close()
		}
    }
}
