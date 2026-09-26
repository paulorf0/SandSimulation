package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	green  = color.RGBA{R: 0, G: 128, B: 0, A: 128}
	yellow = color.RGBA{R: 128, G: 128, B: 0, A: 128}
)

type World struct {
	entities []Rect
}

func NatureForces(mass float64) Vec2d {
	var f Vec2d = Vec2d{0, 0}
	gravity := Vec2d{0, 9.8}

	f = f.Sum(gravity.Scalar(mass))

	return f
}

func (w *World) Update(dt float64) {
	for i := range w.entities {
		w.entities[i].Update(dt)
	}
}

func (w *World) DrawEntities(screen *ebiten.Image) {
	for _, e := range w.entities {
		ebitenutil.DrawRect(screen, e.pos.X, e.pos.Y, float64(e.width), float64(e.height), yellow)
	}
}
