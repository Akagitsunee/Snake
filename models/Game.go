package models

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"math/rand"
)

type Game struct {
	MoveDirectionP1 int
	MoveDirectionP2 int
	SnakeBodyP1     []Position
	SnakeBodyP2     []Position
	Apple           Position
	Timer           int
	MoveTime        int
	Score           int
	BestScore       int
	Level           int
}

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	GridSize     = 10
	XNumInScreen = ScreenWidth / GridSize
	YNumInScreen = ScreenHeight / GridSize
)

const (
	dirNone = iota
	dirLeft
	dirRight
	dirDown
	dirUp
)

func (g *Game) collidesWithApple(snakeBody *[]Position) bool {
	return (*snakeBody)[0].X == g.Apple.X &&
		(*snakeBody)[0].Y == g.Apple.Y
}

func (g *Game) collidesWithSelf(snakeBody *[]Position) bool {
	for _, v := range (*snakeBody)[1:] {
		if (*snakeBody)[0].X == v.X &&
			(*snakeBody)[0].Y == v.Y {
			return true
		}
	}
	return false
}

func (g *Game) collidesWithWall(snakeBody *[]Position) bool {
	return (*snakeBody)[0].X < 0 ||
		(*snakeBody)[0].Y < 0 ||
		(*snakeBody)[0].X >= XNumInScreen ||
		(*snakeBody)[0].Y >= YNumInScreen
}

func (g *Game) collidesWithOtherSnake(snakeBody *[]Position) bool {
	if snakeBody == &g.SnakeBodyP1 {
		return g.loopTroughSnakeBody(snakeBody, &g.SnakeBodyP2)
	} else if snakeBody == &g.SnakeBodyP2 {
		return g.loopTroughSnakeBody(snakeBody, &g.SnakeBodyP1)
	}
	return false
}

func (g *Game) loopTroughSnakeBody(snakeBodyCollide *[]Position, snakeBodyCollided *[]Position) bool {
	for _, v := range *snakeBodyCollided {
		if (*snakeBodyCollide)[0].X == v.X &&
			(*snakeBodyCollide)[0].Y == v.Y {
			return true
		}
	}
	return false
}

func (g *Game) needsToMoveSnake() bool {
	return g.Timer%g.MoveTime == 0
}

func (g *Game) reset() {
	g.Apple.X = 3 * GridSize
	g.Apple.Y = 3 * GridSize
	g.MoveTime = 4
	g.SnakeBodyP1 = g.SnakeBodyP1[:1]
	g.SnakeBodyP1[0].X = XNumInScreen / 2
	g.SnakeBodyP1[0].Y = YNumInScreen / 2
	g.SnakeBodyP2 = g.SnakeBodyP2[:1]
	g.SnakeBodyP2[0].X = XNumInScreen / 3
	g.SnakeBodyP2[0].Y = YNumInScreen / 3
	g.Score = 0
	g.Level = 1
	g.MoveDirectionP1 = dirNone
	g.MoveDirectionP2 = dirNone
}

func (g *Game) isKeyJustPressed(left ebiten.Key, right ebiten.Key, up ebiten.Key, down ebiten.Key, moveDirection *int) {
	if inpututil.IsKeyJustPressed(left) {
		if *moveDirection != dirRight {
			*moveDirection = dirLeft
		}
	} else if inpututil.IsKeyJustPressed(right) {
		if *moveDirection != dirLeft {
			*moveDirection = dirRight
		}
	} else if inpututil.IsKeyJustPressed(down) {
		if *moveDirection != dirUp {
			*moveDirection = dirDown
		}
	} else if inpututil.IsKeyJustPressed(up) {
		if *moveDirection != dirDown {
			*moveDirection = dirUp
		}
	}
}

func (g *Game) snakeMustMove(snake *[]Position, moveDirection *int) {
	if g.collidesWithWall(snake) || g.collidesWithSelf(snake) || g.collidesWithOtherSnake(snake) {
		g.reset()
	}

	if g.collidesWithApple(snake) {
		g.eatApple(snake)
	}

	for i := int64(len(*snake)) - 1; i > 0; i-- {
		(*snake)[i].X = (*snake)[i-1].X
		(*snake)[i].Y = (*snake)[i-1].Y
	}
	switch *moveDirection {
	case dirLeft:
		(*snake)[0].X--
	case dirRight:
		(*snake)[0].X++
	case dirDown:
		(*snake)[0].Y++
	case dirUp:
		(*snake)[0].Y--
	}
}

func (g *Game) Update() error {
	go g.isKeyJustPressed(ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS, &g.MoveDirectionP2)
	go g.isKeyJustPressed(ebiten.KeyArrowLeft, ebiten.KeyArrowRight, ebiten.KeyArrowUp, ebiten.KeyArrowDown, &g.MoveDirectionP1)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.reset()
	}

	if g.needsToMoveSnake() {
		go g.snakeMustMove(&g.SnakeBodyP1, &g.MoveDirectionP1)
		go g.snakeMustMove(&g.SnakeBodyP2, &g.MoveDirectionP2)
	}

	g.Timer++

	return nil
}

func (g *Game) eatApple(snakeBody *[]Position) {
	g.Apple.X = rand.Intn(XNumInScreen - 1)
	g.Apple.Y = rand.Intn(YNumInScreen - 1)
	*snakeBody = append(*snakeBody, Position{
		X: (*snakeBody)[len(*snakeBody)-1].X,
		Y: (*snakeBody)[len(*snakeBody)-1].Y,
	})
	if len(*snakeBody) > 10 && len(*snakeBody) < 20 {
		g.Level = 2
		g.MoveTime = 3
	} else if len(*snakeBody) > 20 {
		g.Level = 3
		g.MoveTime = 2
	} else {
		g.Level = 1
	}
	g.Score++
	if g.BestScore < g.Score {
		g.BestScore = g.Score
	}
}

func (g *Game) drawBody(screen *ebiten.Image, snake *[]Position, rgba color.RGBA) {
	for _, v := range *snake {
		ebitenutil.DrawRect(screen, float64(v.X*GridSize), float64(v.Y*GridSize), GridSize, GridSize, rgba)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	go g.drawBody(screen, &g.SnakeBodyP1, color.RGBA{G: 0xff, B: 0x79, A: 0xCC})
	go g.drawBody(screen, &g.SnakeBodyP2, color.RGBA{B: 0xff, A: 0xCC})

	go ebitenutil.DrawRect(screen, float64(g.Apple.X*GridSize), float64(g.Apple.Y*GridSize), GridSize, GridSize, color.RGBA{0xFF, 0x00, 0x00, 0xff})

	if g.MoveDirectionP1 == dirNone && g.MoveDirectionP2 == dirNone {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Press up/down/left/right to start"))
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f Level: %d Score: %d Best Score: %d", ebiten.CurrentFPS(), g.Level, g.Score, g.BestScore))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
