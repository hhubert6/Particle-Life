package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	simulation "particle-life/internal"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	SCREEN_WIDTH     = 1440
	SCREEN_HEIGHT    = 760
	PARTICLE_RADIUS  = float32(1.5)
	MAX_ZOOM         = 2
	MIN_ZOOM         = 0.3
	ZOOM_IN_FACTOR   = 1.01
	ZOOM_OUT_FACTOR  = 0.99
	DEFAULT_ZOOM     = float64(0.8)
	NUM_OF_PARTICLES = 2000
	NUM_OF_COLORS    = 8
	OFFSET_STEP      = 5
)

var zoom = DEFAULT_ZOOM
var cameraOffsetX = float64(0)
var cameraOffsetY = float64(0)

type Game struct {
	simulation simulation.Simulation
}

func NewGame() *Game {
	return &Game{
		simulation.NewParticleSimulation(
			NUM_OF_PARTICLES,
			createRandomForces(NUM_OF_COLORS),
		),
	}
}

func createRandomForces(numOfColors int) *[][]float64 {
	forceMatrix := make([][]float64, numOfColors)
	for i := range forceMatrix {
		forceMatrix[i] = make([]float64, numOfColors)

		for j := range forceMatrix[i] {
			forceMatrix[i][j] = rand.Float64()*2 - 1
		}
	}
	return &forceMatrix
}

func (g *Game) Update() error {
	g.simulation.Update()

	if ebiten.IsKeyPressed(ebiten.KeyEqual) && zoom <= MAX_ZOOM {
		zoom *= ZOOM_IN_FACTOR
		cameraOffsetX *= ZOOM_IN_FACTOR
		cameraOffsetY *= ZOOM_IN_FACTOR
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) && zoom >= MIN_ZOOM {
		zoom *= ZOOM_OUT_FACTOR
		cameraOffsetY *= ZOOM_OUT_FACTOR
		cameraOffsetX *= ZOOM_OUT_FACTOR
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		cameraOffsetY += OFFSET_STEP
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		cameraOffsetY -= OFFSET_STEP
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		cameraOffsetX -= OFFSET_STEP
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		cameraOffsetX += OFFSET_STEP
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.2f", ebiten.ActualFPS()))

	for _, particle := range g.simulation.Particles() {
		clr := colorful.Hsl(float64(particle.Color)*(360/NUM_OF_COLORS), 0.9, 0.5)
		x, y := getScreenPosition(particle.Position.X, particle.Position.Y)

		if 0 <= x && x <= SCREEN_WIDTH && 0 <= y && y <= SCREEN_HEIGHT {
			DrawCircle(screen, x, y, clr)
		}
	}

	// for k := range g.simulation.Buckets() {
	// 	sx, sy := getScreenPosition(k.X*simulation.R_MAX, k.Y*simulation.R_MAX)
	// 	width := simulation.R_MAX * SCREEN_HEIGHT * zoom
	// 	vector.StrokeRect(screen, sx, sy, float32(width), float32(width), 1, color.White, false)
	// }
}

func getScreenPosition(x, y float64) (float32, float32) {
	centerOffsetX := SCREEN_WIDTH/2 - SCREEN_WIDTH*zoom/2
	centerOffsetY := SCREEN_HEIGHT/2 - SCREEN_HEIGHT*zoom/2

	screenX := x*SCREEN_HEIGHT*zoom + centerOffsetX + cameraOffsetX
	screenY := y*SCREEN_HEIGHT*zoom + centerOffsetY + cameraOffsetY

	return float32(screenX), float32(screenY)
}

func DrawCircle(screen *ebiten.Image, x, y float32, clr color.Color) {
	vector.DrawFilledCircle(screen, x, y, PARTICLE_RADIUS*float32(zoom), clr, true)
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
