package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rooklift/ebiten_example/game"
)

const (
	w = 640
	h = 360
)

func main() {

	game.LoadResources("sprites", "sounds")
	g := game.NewGame(w, h)

	// ebiten.SetWindowTitle("Foo")
	// ebiten.SetWindowSize(w * 2, h * 2)
	ebiten.SetFullscreen(true)

	err := ebiten.RunGame(g)
	if err != nil && err != game.USER_QUIT {
		fmt.Printf("%v\n", err)
	}
}
