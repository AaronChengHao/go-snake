package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"snake/game"
)

func main() {
	ebiten.SetWindowSize(game.WindowW, game.WindowH)
	ebiten.SetWindowTitle("睿曦与睿懿的故事")
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
