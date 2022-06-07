package models

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func (g *Game) collidesWithApple(snakeBody []Position) bool {
	return snakeBody[0].X == g.Apple.X &&
		snakeBody[0].Y == g.Apple.Y
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

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if g.MoveDirectionP1 != dirRight {
			g.MoveDirectionP1 = dirLeft
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if g.MoveDirectionP1 != dirLeft {
			g.MoveDirectionP1 = dirRight
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if g.MoveDirectionP1 != dirUp {
			g.MoveDirectionP1 = dirDown
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if g.MoveDirectionP1 != dirDown {
			g.MoveDirectionP1 = dirUp
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if g.MoveDirectionP2 != dirRight {
			g.MoveDirectionP2 = dirLeft
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if g.MoveDirectionP2 != dirLeft {
			g.MoveDirectionP2 = dirRight
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.MoveDirectionP2 != dirUp {
			g.MoveDirectionP2 = dirDown
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if g.MoveDirectionP2 != dirDown {
			g.MoveDirectionP2 = dirUp
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.reset()
	}

	if g.needsToMoveSnake() {
		if g.collidesWithWall(&g.SnakeBodyP1) || g.collidesWithWall(&g.SnakeBodyP2) || g.collidesWithSelf(&g.SnakeBodyP1) || g.collidesWithSelf(&g.SnakeBodyP2) || g.collidesWithOtherSnake(&g.SnakeBodyP1) || g.collidesWithOtherSnake(&g.SnakeBodyP2) {
			g.reset()
		}

		if g.collidesWithApple(g.SnakeBodyP1) {
			go g.eatApple(&g.SnakeBodyP1)
		} else if g.collidesWithApple(g.SnakeBodyP2) {
			g.eatApple(&g.SnakeBodyP2)
		}

		for i := int64(len(g.SnakeBodyP1)) - 1; i > 0; i-- {
			g.SnakeBodyP1[i].X = g.SnakeBodyP1[i-1].X
			g.SnakeBodyP1[i].Y = g.SnakeBodyP1[i-1].Y
		}
		switch g.MoveDirectionP1 {
		case dirLeft:
			g.SnakeBodyP1[0].X--
		case dirRight:
			g.SnakeBodyP1[0].X++
		case dirDown:
			g.SnakeBodyP1[0].Y++
		case dirUp:
			g.SnakeBodyP1[0].Y--
		}

		for i := int64(len(g.SnakeBodyP2)) - 1; i > 0; i-- {
			g.SnakeBodyP2[i].X = g.SnakeBodyP2[i-1].X
			g.SnakeBodyP2[i].Y = g.SnakeBodyP2[i-1].Y
		}
		switch g.MoveDirectionP2 {
		case dirLeft:
			g.SnakeBodyP2[0].X--
		case dirRight:
			g.SnakeBodyP2[0].X++
		case dirDown:
			g.SnakeBodyP2[0].Y++
		case dirUp:
			g.SnakeBodyP2[0].Y--
		}
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

func (g *Game) Draw(screen *ebiten.Image) {
	for _, v := range g.SnakeBodyP1 {
		go ebitenutil.DrawRect(screen, float64(v.X*GridSize), float64(v.Y*GridSize), GridSize, GridSize, color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	}
	ebitenutil.DrawRect(screen, float64(g.Apple.X*GridSize), float64(g.Apple.Y*GridSize), GridSize, GridSize, color.RGBA{0xFF, 0x00, 0x00, 0xff})

	if g.MoveDirectionP1 == dirNone {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Press up/down/left/right to start"))
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f Level: %d Score: %d Best Score: %d", ebiten.CurrentFPS(), g.Level, g.Score, g.BestScore))
	}

	for _, v := range g.SnakeBodyP2 {
		go ebitenutil.DrawRect(screen, float64(v.X*GridSize), float64(v.Y*GridSize), GridSize, GridSize, color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
