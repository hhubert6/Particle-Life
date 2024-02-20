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
	PARTICLE_RADIUS  = float32(2)
	MAX_ZOOM         = 3
	MIN_ZOOM         = 0.4
	ZOOM_IN_FACTOR   = 1.01
	ZOOM_OUT_FACTOR  = 0.99
	DEFAULT_ZOOM     = float64(1)
	NUM_OF_PARTICLES = 5000
	NUM_OF_COLORS    = 2
	OFFSET_STEP      = 5
)

var zoom = DEFAULT_ZOOM
var cameraOffsetX = float64(0)
var cameraOffsetY = float64(0)

var paused = false
var spaceStagger = false
var spaceDelay = 0

var forces *[][]float64

type Game struct {
	simulation simulation.Simulation
}

func NewGame() *Game {
	forces = createForces()
	return &Game{
		simulation.NewParticleSimulation(
			NUM_OF_PARTICLES,
			forces,
		),
	}
}

func createForces() *[][]float64 {
	forceMatrix := make([][]float64, NUM_OF_COLORS)
	for i := range forceMatrix {
		forceMatrix[i] = make([]float64, NUM_OF_COLORS)
	}
	return &forceMatrix
}

func randomForces() {
	for i := range *forces {
		for j := range (*forces)[i] {
			(*forces)[i][j] = rand.Float64()*2 - 1
		}
	}
}

func (g *Game) Update() error {
	if !paused {
		g.simulation.Update()
	}

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
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !spaceStagger {
			paused = !paused
			spaceStagger = true
		}
	}
	if spaceStagger {
		spaceDelay += 1
		if spaceDelay > 30 {
			spaceStagger = false
			spaceDelay = 0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.2f", ebiten.ActualFPS()))

	for _, particle := range g.simulation.Particles() {
		clr := colorful.Hsl(float64(particle.Color)*(360/NUM_OF_COLORS), 0.9, 0.5)
		x, y := getScreenPosition(particle.Position.X, particle.Position.Y)

		if 0 <= x && x <= SCREEN_WIDTH && 0 <= y && y <= SCREEN_HEIGHT {
			DrawRect(screen, x, y, clr)
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

func DrawRect(screen *ebiten.Image, x, y float32, clr color.Color) {
	vector.DrawFilledRect(screen, x, y, PARTICLE_RADIUS*float32(zoom), PARTICLE_RADIUS*float32(zoom), clr, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	go func() {
		for {
			fmt.Println("Enter:")

			var input string
			_, err := fmt.Scanf("%s", &input)
			if err != nil {
				fmt.Println(err)
				continue
			}

			switch input {
			case "reset":
				for i := range *forces {
					for j := range (*forces)[i] {
						(*forces)[i][j] = 0
					}
				}
			case "random":
				randomForces()
			case "set":
				var row, col int
				var value float64
				_, err := fmt.Scanf("%d %d %f", &row, &col, &value)

				if err != nil || value < -1 || value > 1 || row < 0 || col < 0 || row >= len(*forces) || col >= len(*forces) {
					fmt.Println("Wrong input!", err)
					continue
				}

				(*forces)[row][col] = value
			default:
				fmt.Println("Wrong input!")
			}
		}
	}()

	game := NewGame()

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Particle artificial life")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
