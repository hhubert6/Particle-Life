package simulation

import (
	"math"
	"math/rand"
	"sync"
)

const (
	DT                 = 0.016
	FRICTION_HALF_LIFE = 0.04
	R_MAX              = 0.15
	BETA               = 0.3
	FORCE_FACTOR       = 5
)

var frictionFactor = math.Pow(0.5, DT/FRICTION_HALF_LIFE)

type Vec2 struct {
	X, Y float64
}

type Particle struct {
	Position, Velocity Vec2
	Color              int
}

type Simulation interface {
	Update()
	Particles() []Particle
}

type ParticlesSimulation struct {
	particles []Particle
	forceM    [][]float64
}

func NewParticleSimulation(numOfParticles int, forceM [][]float64) *ParticlesSimulation {
	if len(forceM) == 0 {
		panic("Empty force matrix")
	}
	for i := range forceM {
		if len(forceM) != len(forceM[i]) {
			panic("Force matrix is not square matrix")
		}
	}

	return &ParticlesSimulation{
		particles: generateParticles(numOfParticles, len(forceM)),
		forceM:    forceM,
	}
}

// n - number of particles to generate, m - number of particles colors
func generateParticles(n, m int) []Particle {
	container := make([]Particle, n)

	for i := range container {
		x := rand.Float64() * 1.9
		y := rand.Float64()
		color := rand.Intn(m)

		container[i] = Particle{
			Position: Vec2{x, y},
			Velocity: Vec2{},
			Color:    color,
		}
	}

	return container
}

func (s *ParticlesSimulation) Update() {
	wg := sync.WaitGroup{}

	for i := range s.particles {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			totalForceX := float64(0)
			totalForceY := float64(0)

			for j := range s.particles {
				if i == j {
					continue
				}

				rx := s.particles[j].Position.X - s.particles[i].Position.X
				ry := s.particles[j].Position.Y - s.particles[i].Position.Y
				r := math.Hypot(rx, ry)

				if 0 < r && r < R_MAX {
					f := force(r/R_MAX, s.forceM[s.particles[i].Color][s.particles[j].Color])
					totalForceX += rx / r * f
					totalForceY += ry / r * f
				}
			}

			totalForceX *= R_MAX * FORCE_FACTOR
			totalForceY *= R_MAX * FORCE_FACTOR

			s.particles[i].Velocity.X *= frictionFactor
			s.particles[i].Velocity.Y *= frictionFactor

			s.particles[i].Velocity.X += totalForceX * DT
			s.particles[i].Velocity.Y += totalForceY * DT

			s.particles[i].Position.X += s.particles[i].Velocity.X * DT
			s.particles[i].Position.Y += s.particles[i].Velocity.Y * DT
		}(i)
	}

	wg.Wait()
}

func force(r, a float64) float64 {
	if r < BETA {
		return r/BETA - 1
	} else if r < 1 {
		return a * (1 - math.Abs(2*r-1-BETA)/(1-BETA))
	} else {
		return 0
	}
}

func (s *ParticlesSimulation) Particles() []Particle {
	return s.particles
}
