package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"snake/models"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newGame() *models.Game {
	g := &models.Game{
		Apple:       models.Position{X: 3 * models.GridSize, Y: 3 * models.GridSize},
		MoveTime:    4,
		SnakeBodyP1: make([]models.Position, 1),
		SnakeBodyP2: make([]models.Position, 1),
	}
	g.SnakeBodyP1[0].X = models.XNumInScreen / 2
	g.SnakeBodyP1[0].Y = models.YNumInScreen / 2
	g.SnakeBodyP2[0].X = models.XNumInScreen / 3
	g.SnakeBodyP2[0].Y = models.YNumInScreen / 3
	return g
}

func main() {
	ebiten.SetWindowSize(models.ScreenWidth, models.ScreenHeight)
	ebiten.SetWindowTitle("Snake (Ebiten Demo)")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
