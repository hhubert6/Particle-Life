package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SCREEN_WIDTH    = 1280
	SCREEN_HEIGHT   = 720
	PARTICLE_RADIUS = 2
)

var forceMatrix = [][]float32{
	{0.02, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}

type Vec2 struct {
	X, Y float32
}

type Atom struct {
	Position, Velocity Vec2
	Color              int
}

type Game struct {
	atoms1, atoms2, atoms3 []Atom
}

func NewGame() *Game {
	game := &Game{
		atoms1: generateAtoms(300, 0),
	}

	game.atoms1 = append(game.atoms1, generateAtoms(300, 1)...)
	game.atoms1 = append(game.atoms1, generateAtoms(300, 2)...)

	return game
}

func generateAtoms(n int, color int) (container []Atom) {
	container = make([]Atom, n)

	for i := range container {
		container[i] = Atom{
			Position: randomVec2(
				PARTICLE_RADIUS,
				PARTICLE_RADIUS,
				SCREEN_WIDTH-PARTICLE_RADIUS,
				SCREEN_HEIGHT-PARTICLE_RADIUS,
			),
			Velocity: Vec2{},
			Color:    color,
		}
	}

	return
}

func randomVec2(minX, minY, maxX, maxY float32) Vec2 {
	x := rand.Float32()*(maxX-minX) + minX
	y := rand.Float32()*(maxY-minY) + minY
	return Vec2{x, y}
}

func (g *Game) Update() error {
	updateAtoms(g.atoms1)
	updateAtoms(g.atoms2)
	updateAtoms(g.atoms3)

	g.Rule(g.atoms1, g.atoms1, 0.02, 200)

	return nil
}

func updateAtoms(atoms []Atom) {
	for i := range atoms {
		atoms[i].Position.X += atoms[i].Velocity.X
		atoms[i].Position.Y += atoms[i].Velocity.Y
		atoms[i].Velocity.X *= 0.7
		atoms[i].Velocity.Y *= 0.7

		if atoms[i].Position.X <= 0 || atoms[i].Position.X >= SCREEN_WIDTH {
			atoms[i].Velocity.X *= -2
		}
		if atoms[i].Position.Y <= 0 || atoms[i].Position.Y >= SCREEN_HEIGHT {
			atoms[i].Velocity.Y *= -2
		}
	}
}

func (g *Game) Rule(atoms1, atoms2 []Atom, force, forceRange float32) {
	for i := range atoms1 {
		forceX := float32(0)
		forceY := float32(0)

		for j := range atoms2 {
			d := distance(atoms1[i].Position, atoms2[j].Position)
			factor := float32(0)
			if 0 < d && d < forceRange/3 {
				factor += -0.1 / d
			}
			if 0 < d && d < 3*PARTICLE_RADIUS {
				factor += -1 / d
			} else if 0 < d && d < forceRange {
				factor += forceMatrix[atoms1[i].Color][atoms2[i].Color] / d
			}
			forceX += (atoms2[j].Position.X - atoms1[i].Position.X) * factor
			forceY += (atoms2[j].Position.Y - atoms1[i].Position.Y) * factor
		}

		atoms1[i].Velocity.X += forceX
		atoms1[i].Velocity.Y += forceY

	}
}

func distance(posA, posB Vec2) float32 {
	x := float64(posB.X - posA.X)
	y := float64(posB.Y - posA.Y)
	return float32(math.Sqrt(x*x + y*y))
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.2f", ebiten.ActualFPS()))

	var clr color.Color

	for _, atom := range g.atoms1 {
		switch atom.Color {
		case 0:
			clr = color.RGBA{255, 0, 0, 255}
		case 1:
			clr = color.RGBA{0, 255, 0, 255}
		case 2:
			clr = color.RGBA{0, 0, 255, 255}
		}
		DrawCircle(screen, atom.Position.X, atom.Position.Y, clr)
	}
}

func DrawCircle(screen *ebiten.Image, x, y float32, clr color.Color) {
	vector.DrawFilledCircle(screen, x, y, PARTICLE_RADIUS, clr, true)
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
