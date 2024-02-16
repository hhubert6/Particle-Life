package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"particle_life/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SCREEN_WIDTH    = 920
	SCREEN_HEIGHT   = 720
	PARTICLE_RADIUS = float32(2)
	MAX_ZOOM        = 2
	MIN_ZOOM        = 0.3
)

var zoom = float64(1)
var cameraOffsetX = float64(0)
var cameraOffsetY = float64(0)

type Game struct {
	simulation simulation.Simulation
}

func NewGame() *Game {
	return &Game{
		simulation.NewParticleSimulation(900, createRandomForces(5)),
	}
}

func createRandomForces(numOfColors int) [][]float64 {
	forceMatrix := make([][]float64, numOfColors)
	for i := range forceMatrix {
		forceMatrix[i] = make([]float64, numOfColors)

		for j := range forceMatrix[i] {
			forceMatrix[i][j] = rand.Float64()*2 - 1
		}
	}
	return forceMatrix
}

func (g *Game) Update() error {
	g.simulation.Update()

	if ebiten.IsKeyPressed(ebiten.KeyEqual) && zoom <= MAX_ZOOM {
		zoom *= 1.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) && zoom >= MIN_ZOOM {
		zoom *= 0.99
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		cameraOffsetY += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		cameraOffsetY -= 2 
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		cameraOffsetX -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		cameraOffsetX += 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.2f", ebiten.ActualFPS()))

	var clr color.Color

	for _, particle := range g.simulation.Particles() {
		switch particle.Color {
		case 0:
			clr = color.RGBA{255, 0, 0, 255}
		case 1:
			clr = color.RGBA{0, 255, 0, 255}
		case 2:
			clr = color.RGBA{0, 0, 255, 255}
		case 3:
			clr = color.RGBA{255, 0, 255, 255}
		case 4:
			clr = color.RGBA{255, 255, 0, 255}
		}
		DrawCircle(
			screen,
			(particle.Position.X*SCREEN_HEIGHT*zoom + SCREEN_WIDTH/2 - SCREEN_HEIGHT*zoom/2) + cameraOffsetX,
			(particle.Position.Y*SCREEN_HEIGHT*zoom + SCREEN_HEIGHT/2 - SCREEN_HEIGHT*zoom/2) + cameraOffsetY,
			clr,
		)
	}
}

func DrawCircle(screen *ebiten.Image, x, y float64, clr color.Color) {
	vector.DrawFilledCircle(
		screen,
		float32(x),
		float32(y),
		PARTICLE_RADIUS*float32(zoom),
		clr,
		true,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Particle artificial life")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
