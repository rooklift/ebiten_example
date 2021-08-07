package game

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var sprites map[string]*ebiten.Image
var sounds map[string][]byte

func LoadResources(sprites_path string, sounds_path string) {
	load_sprites(sprites_path)
	load_sounds(sounds_path)
}

func load_sprites(sprites_path string) {

	sprites = make(map[string]*ebiten.Image)

	files, err := ioutil.ReadDir(sprites_path)
	if err != nil {
		panic(err)
	}

	for _, info := range files {
		f, err := os.Open(filepath.Join(sprites_path, info.Name()))
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

func load_sounds(sounds_path string) {

	sounds = make(map[string][]byte)

	files, err := ioutil.ReadDir(sounds_path)
	if err != nil {
		panic(err)
	}

	for _, info := range files {
		f, err := os.Open(filepath.Join(sounds_path, info.Name()))
		if err != nil {
			fmt.Println(err)
		} else {
			f.Seek(44, io.SeekStart)						// Skip WAV header
			sounds[info.Name()], _ = ioutil.ReadAll(f)
			f.Close()
		}
	}
}

// ------------------------------------------------------------------------------------------------

type Game struct {

	inited bool

	width int
	height int
	image *ebiten.Image

	audio_context *audio.Context
	audio_players []*audio.Player

	entities []*Entity
}

func NewGame(width int, height int) *Game {
	ret := new(Game)
	ret.width = width
	ret.height = height
	return ret
}

func (self *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return self.width, self.height
}

func (self *Game) Draw(screen *ebiten.Image) {

	self.image.Clear()

	for _, ent := range self.entities {
		ent.Draw()
	}

	screen.DrawImage(self.image, nil)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f -- FPS: %0.2f -- Sounds: %v", ebiten.CurrentTPS(), ebiten.CurrentFPS(), len(self.audio_players)))
}

func (self *Game) Update() error {

	if (!self.inited) {
		self.audio_context = audio.NewContext(44100)
		self.image = ebiten.NewImage(self.width, self.height)
		self.entities = append(self.entities, NewEntity(self, PLAYER, 16, 16, 0, 0, "ship.png"))
		self.inited = true
	}

	self.PurgeAudio()

	for _, ent := range self.entities {
		ent.Behave()
	}

	return nil
}

// ------------------------------------------------------------------------------------------------

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

type EntityType int

const (
	PLAYER				EntityType = iota
)

type Entity struct {
	game				*Game
	t					EntityType
	x					float64
	y					float64
	speedx				float64
	speedy				float64
	sprite_string		string
}

func NewEntity(game *Game, t EntityType, x float64, y float64, speedx float64, speedy float64, sprite_string string) *Entity {

	if sprites[sprite_string] == nil {
		panic("NewEntity: unknown sprite")
	}

	ret := &Entity{
		game: game,
		t: t,
		x: x,
		y: y,
		speedx: speedx,
		speedy: speedy,
		sprite_string: sprite_string,
	}

	return ret
}

func (self *Entity) Draw() {

	img := sprites[self.sprite_string]
	e_width, e_height := img.Size()
	opts := new(ebiten.DrawImageOptions)
	opts.GeoM.Translate(self.x - (float64(e_width) / 2), self.y - (float64(e_height) / 2.0))

	self.game.image.DrawImage(img, opts)
}

func (self *Entity) Behave() {

	switch self.t {

	case PLAYER:

		if ebiten.IsKeyPressed(ebiten.KeyD) {
			self.speedx += 0.1
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			self.speedx -= 0.1
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			self.speedy += 0.1
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			self.speedy -= 0.1
		}

		if (self.x < 0) {
			self.speedx = math.Abs(self.speedx)
			self.x = 0
			self.game.PlaySound("test.wav")
		}
		if (self.x > float64(self.game.width)) {
			self.speedx = math.Abs(self.speedx) * -1
			self.x = float64(self.game.width)
			self.game.PlaySound("test.wav")
		}
		if (self.y < 0) {
			self.speedy = math.Abs(self.speedy)
			self.y = 0
			self.game.PlaySound("test.wav")
		}
		if (self.y > float64(self.game.height)) {
			self.speedy = math.Abs(self.speedy) * -1
			self.y = float64(self.game.height)
			self.game.PlaySound("test.wav")
		}

		self.x += self.speedx
		self.y += self.speedy
	}
}

