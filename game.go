package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	world *World
}

func (g *Game) Update(screen *ebiten.Image) error {
	const spawnInterval = 2
	if d := inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft); d > 0 && (d-1)%spawnInterval == 0 {
		x, y := ebiten.CursorPosition()
		g.CreateEntity(float64(x), float64(y))
	}

	g.world.Update(1.0 / 60)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.DrawEntities(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Game() {

}

func (g *Game) CreateEntity(x, y float64) {
	entity := NewRect(x, y, 3, 3, 10)
	g.world.entities = append(g.world.entities, *entity)
}

func main() {
	game := &Game{
		world: &World{},
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Sand Simulation")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
