package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_GROUPS = 3

type Config struct {
	winWidth               int32
	winHeight              int32
	winTitle               string
	fps                    int32
	background             rl.Color
	particleRadius         float32
	particlesPerGroup      int
	particleVelocityFactor float64
	particleActionDistance float64
}

type Position struct {
	X, Y float64
}

func (p Position) AsInt() (int32, int32) {
	return int32(p.X), int32(p.Y)
}

func RandomPosition(config *Config) Position {
	x := rand.Intn(int(config.winWidth))
	y := rand.Intn(int(config.winHeight))

	return Position{
		X: float64(x),
		Y: float64(y),
	}
}

type Particle struct {
	color    rl.Color
	position Position
	vx       float64
	vy       float64
}

func NewParticle(config *Config, color rl.Color) Particle {
	return Particle{
		color:    color,
		position: RandomPosition(config),
		vx:       0.0,
		vy:       0.0,
	}
}

func (p *Particle) Render(config *Config) {
	x, y := p.position.AsInt()
	rl.DrawCircle(x, y, config.particleRadius, p.color)
}

type ParticleDistance struct {
	xDelta       float64
	yDelta       float64
	displacement float64
}

func ParticleDistanceCalculate(one, two *Particle) ParticleDistance {
	dx := one.position.X - two.position.X
	dy := one.position.Y - two.position.Y

	displacement := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	return ParticleDistance{
		xDelta:       dx,
		yDelta:       dy,
		displacement: displacement,
	}
}

type InteractivityMatrix [MAX_GROUPS][MAX_GROUPS]float64
type ParticleGroup []Particle
type ParticleGroups [MAX_GROUPS]ParticleGroup

func NewParticleGroup(config *Config, color rl.Color) ParticleGroup {
	group := ParticleGroup{}
	for i := 0; i < config.particlesPerGroup; i++ {
		particle := NewParticle(config, color)
		group = append(group, particle)
	}

	return group
}

func ParticleGroupsUpdate(gOne ParticleGroup, gTwo ParticleGroup, gravity float64, config *Config) {
	max := len(gOne)

	for i := 0; i < max; i++ {
		fx := 0.0
		fy := 0.0

		for j := 0; j < max; j++ {
			a := &gOne[i]
			b := &gTwo[j]

			distance := ParticleDistanceCalculate(a, b)
			if distance.displacement > 0.0 && distance.displacement < config.particleActionDistance {
				force := gravity / distance.displacement
				fx += force * distance.xDelta
				fy += force * distance.yDelta
			}

			a.vx = (a.vx + fx) * config.particleVelocityFactor
			a.vy = (a.vy * fy) * config.particleVelocityFactor

			a.position.X += a.vx
			a.position.Y += a.vy

			if a.position.X < 0.0 || a.position.X > float64(config.winWidth) {
				a.vx *= -1.0
			}

			if a.position.Y < 0.0 || a.position.Y > float64(config.winHeight) {
				a.vy *= -1.0
			}
		}
	}
}

func (groups ParticleGroups) ApplyInteractivity(interactivity *InteractivityMatrix, config *Config) {
	for i := 0; i < MAX_GROUPS; i++ {
		for j := 0; j < MAX_GROUPS; j++ {
			ParticleGroupsUpdate(groups[i], groups[j], interactivity[i][j], config)
		}
	}
}

func (groups ParticleGroups) Render(config *Config) {
	for _, group := range groups {
		for _, particle := range group {
			particle.Render(config)
		}
	}
}

func main() {
	config := &Config{
		winWidth:               800,
		winHeight:              450,
		winTitle:               "Game",
		fps:                    60,
		background:             rl.Black,
		particleRadius:         2.0,
		particlesPerGroup:      300,
		particleVelocityFactor: 0.01,
		particleActionDistance: 75.0,
	}

	matrix := &InteractivityMatrix{
		{0.02, -0.5, -0.02}, // red
		{0.02, -0.04, -0.7}, // yellow
		{0.02, 0.02, 0.02},  // blue
	}

	particleGroups := ParticleGroups{}
	particleGroups[0] = NewParticleGroup(config, rl.Red)
	particleGroups[1] = NewParticleGroup(config, rl.Yellow)
	particleGroups[2] = NewParticleGroup(config, rl.Blue)

	rl.InitWindow(config.winWidth, config.winHeight, config.winTitle)
	defer rl.CloseWindow()
	rl.SetTargetFPS(config.fps)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(config.background)

		particleGroups.ApplyInteractivity(matrix, config)
		particleGroups.Render(config)

		rl.EndDrawing()
	}
}
