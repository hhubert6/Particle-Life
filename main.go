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
	PARTICLE_RADIUS = float32(1)
)

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
			(particle.Position.X*SCREEN_HEIGHT+SCREEN_WIDTH/2-SCREEN_HEIGHT/2)*0.8,
			(particle.Position.Y*SCREEN_HEIGHT)*0.8,
			clr,
		)
	}
}

func DrawCircle(screen *ebiten.Image, x, y float64, clr color.Color) {
	vector.DrawFilledCircle(
		screen,
		float32(x),
		float32(y),
		PARTICLE_RADIUS,
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
