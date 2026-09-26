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
	w.Collision()
}

func (w *World) DrawEntities(screen *ebiten.Image) {
	for _, e := range w.entities {
		ebitenutil.DrawRect(screen, e.pos.X, e.pos.Y, float64(e.width), float64(e.height), yellow)
	}
}

func (w *World) Collision() {
	n := len(w.entities)
	for i := range w.entities {
		e := &w.entities[i]
		floor := float64(screenHeight - e.height)
		if e.pos.Y > floor {
			e.pos.Y = floor
			e.vel.Y = 0
		}
	}

	for p := 0; p < n; p++ {
		e1 := &w.entities[p]
		for q := p + 1; q < n; q++ {
			e2 := &w.entities[q]

			if mtv, ok := e1.Penetration(*e2); ok {
				if mtv.Y > 0 {
					// e1 está embaixo: quem sobe é e2, para não afundar e1 no chão
					e2.pos = e2.pos.Sum(mtv.Opposite())
					e2.vel.Y = 0
				} else {
					e1.pos = e1.pos.Sum(mtv)
					if mtv.Y != 0 {
						e1.vel.Y = 0
					}
				}
			}
		}
	}
}
