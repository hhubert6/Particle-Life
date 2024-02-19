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
	Buckets() map[Vec2][]Particle
}

type ParticlesSimulation struct {
	particles []Particle
	forceM    *[][]float64
	buckets   map[Vec2][]Particle
}

func NewParticleSimulation(numOfParticles int, forceM *[][]float64) *ParticlesSimulation {
	if len(*forceM) == 0 {
		panic("Empty force matrix")
	}
	for i := range *forceM {
		if len(*forceM) != len((*forceM)[i]) {
			panic("Force matrix is not square matrix")
		}
	}

	simulation := &ParticlesSimulation{
		forceM:  forceM,
		buckets: make(map[Vec2][]Particle, 100),
	}

	simulation.particles = generateParticles(numOfParticles, len(*forceM), simulation.buckets)

	return simulation
}

// n - number of particles to generate, m - number of particles colors
func generateParticles(n, m int, buckets map[Vec2][]Particle) []Particle {
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

		key := getBucketKey(x, y)
		buckets[key] = append(buckets[key], container[i])
	}

	return container
}

func (s *ParticlesSimulation) Update() {
	wg := sync.WaitGroup{}

	for i := range s.particles {
		wg.Add(1)
		go s.updateParticle(i, &wg)
	}

	wg.Wait()

	s.updateBuckets()
}

func (s *ParticlesSimulation) updateParticle(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	p := s.particles[i]

	totalForceX := float64(0)
	totalForceY := float64(0)

	gridX, gridY := getGridPosition(p.Position.X, p.Position.Y)

	for offsetX := float64(-1); offsetX < 2; offsetX++ {
		for offsetY := float64(-1); offsetY < 2; offsetY++ {

			key := Vec2{X: gridX + offsetX, Y: gridY + offsetY}
			bucket, ok := s.buckets[key]

			if ok {
				for j := range bucket {
					rx := bucket[j].Position.X - p.Position.X
					ry := bucket[j].Position.Y - p.Position.Y
					r := math.Hypot(rx, ry)

					if 0 < r && r < R_MAX {
						f := force(r/R_MAX, (*s.forceM)[p.Color][bucket[j].Color])
						totalForceX += rx / r * f
						totalForceY += ry / r * f
					}
				}
			}
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

func (s *ParticlesSimulation) updateBuckets() {
	for v := range s.buckets {
		delete(s.buckets, v)
	}

	for i := range s.particles {
		key := getBucketKey(s.particles[i].Position.X, s.particles[i].Position.Y)
		s.buckets[key] = append(s.buckets[key], s.particles[i])
	}
}

func getBucketKey(x, y float64) Vec2 {
	gridX, gridY := getGridPosition(x, y)
	return Vec2{X: gridX, Y: gridY}
}

func getGridPosition(x, y float64) (gridX, gridY float64) {
	gridX = math.Floor(x / R_MAX)
	gridY = math.Floor(y / R_MAX)
	return
}

func (s *ParticlesSimulation) Particles() []Particle {
	return s.particles
}

func (s *ParticlesSimulation) Buckets() map[Vec2][]Particle {
	return s.buckets
}
