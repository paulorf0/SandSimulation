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

// Game implements ebiten.Game interface.
type Game struct {
	world *World
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	// Write your game's logical update.
	// segurando o botão, cria um retângulo a cada spawnInterval frames
	const spawnInterval = 2
	if d := inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft); d > 0 && (d-1)%spawnInterval == 0 {
		x, y := ebiten.CursorPosition()
		g.CreateEntity(float64(x), float64(y))
	}

	g.world.Update(1.0 / 60)
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.world.DrawEntities(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
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
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Your game's title")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
