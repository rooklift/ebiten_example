package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rooklift/ebiten_example/game"
)

const (
	w = 600
	h = 400
)

func main() {

	game.LoadResources("sprites", "sounds")
	g := game.NewGame(w, h)

	ebiten.SetWindowSize(w * 2, h * 2)
	ebiten.SetWindowTitle("Foo")

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
