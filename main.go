package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	paddleWidth  = 15
	paddleHeight = 100
	ballSize     = 15
	paddleSpeed  = 4
	ballSpeed    = 3
	enemyPaddleSpeed = 2.5
)

type Game struct {
	player1Y, player2Y float64
	ballX, ballY       float64
	ballVelocityX      float64
	ballVelocityY      float64
	player1Score       int
	player2Score       int
}

func NewGame() *Game {
	g := &Game{
		player1Y:     float64(screenHeight-paddleHeight) / 2,
		player2Y:     float64(screenHeight-paddleHeight) / 2,
		ballVelocityX: ballSpeed,
		ballVelocityY: ballSpeed,
	}
	g.resetBall()
	return g
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.handleInput()

	// Update ball position
	g.ballX += g.ballVelocityX
	g.ballY += g.ballVelocityY

	// Ball collision with walls
	if g.ballY < 0 {
		g.ballVelocityY = -g.ballVelocityY
	} else if g.ballY > screenHeight-ballSize {
		g.ballVelocityY = -g.ballVelocityY
	}
	collisionWithLeftPaddle := g.ballX < paddleWidth && g.ballY > g.player1Y && g.ballY < g.player1Y+paddleHeight
	collisionWithRightPaddle := g.ballX > screenWidth-paddleWidth-ballSize && g.ballY > g.player2Y && g.ballY < g.player2Y+paddleHeight

	// Ball collision with paddles
	if collisionWithLeftPaddle || collisionWithRightPaddle {
		g.ballVelocityX = -g.ballVelocityX
	}

	// Scoring
	if g.ballX < 0 {
		g.player2Score++
		g.resetBall()
	}
	if g.ballX > screenWidth-ballSize {
		g.player1Score++
		g.resetBall()
	}

	// AI for player2
	if g.ballVelocityX > 0 || g.ballVelocityY < 0 { 
		if g.ballY < g.player2Y+paddleHeight/2 {
			g.player2Y -= enemyPaddleSpeed
		} else if g.ballY > g.player2Y+paddleHeight/2 {
			g.player2Y += enemyPaddleSpeed
		}
	}

	g.draw(screen)
	return nil
}

func (g *Game) handleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.player1Y > 0 {
		g.player1Y -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.player1Y < screenHeight-paddleHeight {
		g.player1Y += paddleSpeed
	}
}

func (g *Game) resetBall() {
	g.ballX = float64(screenWidth) / 2
	g.ballY = float64(screenHeight) / 2

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		g.ballVelocityX = -ballSpeed
	} else {
        g.ballVelocityX = ballSpeed 
    }
	if rand.Intn(2) == 0 {
        g.ballVelocityY = -ballSpeed 
    } else {
        g.ballVelocityY = ballSpeed 
    }
}

func (g *Game) draw(screen *ebiten.Image) {
    // Draw paddles and ball 
    ebitenutil.DrawRect(screen, paddleWidth, g.player1Y, paddleWidth, paddleHeight, color.White)
    ebitenutil.DrawRect(screen, screenWidth-paddleWidth-ballSize, g.player2Y, paddleWidth, paddleHeight, color.White)
    ebitenutil.DrawRect(screen, g.ballX, g.ballY, ballSize, ballSize, color.White)

    // Draw scores 
    ebitenutil.DebugPrint(screen, fmt.Sprintf("Player 1: %d", g.player1Score))
    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Player 2: %d", g.player2Score), 0, 20)
}

func main() {
	game := NewGame()
	if err := ebiten.Run(game.Update, screenWidth, screenHeight, 1, "Pong"); err != nil {
        log.Fatal(err)
    }
}
