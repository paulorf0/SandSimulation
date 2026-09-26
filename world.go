package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// gravidade em px/s²
const gravity = 300.0

var (
	green  = color.RGBA{R: 0, G: 128, B: 0, A: 128}
	yellow = color.RGBA{R: 128, G: 128, B: 0, A: 128}
)

// pixel é uma imagem 1x1 branca, escalada e rotacionada para desenhar os retângulos
var pixel *ebiten.Image

func init() {
	pixel, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	pixel.Fill(color.White)
}

type World struct {
	entities []Rect

	// grade espacial reaproveitada entre frames (ver forEachNearbyPair)
	cells    [][]int
	cellOfID []int
}

func NatureForces(mass float64) Vec2d {
	var f Vec2d = Vec2d{0, 0}
	g := Vec2d{0, gravity}

	f = f.Sum(g.Scalar(mass))

	return f
}

func (w *World) Update(dt float64) {
	for i := range w.entities {
		w.entities[i].Update(dt)
	}
	w.Collision(dt)
}

func (w *World) DrawEntities(screen *ebiten.Image) {
	c := color.NRGBAModel.Convert(yellow).(color.NRGBA)
	for _, e := range w.entities {
		wd, ht := float64(e.width), float64(e.height)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(wd, ht)
		op.GeoM.Translate(-wd/2, -ht/2)
		op.GeoM.Rotate(e.angle)
		op.GeoM.Translate(e.pos.X+wd/2, e.pos.Y+ht/2)
		op.ColorM.Scale(float64(c.R)/255, float64(c.G)/255, float64(c.B)/255, float64(c.A)/255)
		screen.DrawImage(pixel, op)
	}
}

// collisionIterations é quantas vezes por frame as sobreposições são resolvidas.
// Empurrar um bloco pode enfiá-lo em outro; as passadas seguintes corrigem isso.
const collisionIterations = 4

func (w *World) Collision(dt float64) {
	n := len(w.entities)

	// trecho em X onde cada entidade está apoiada neste frame, somando todos
	// os blocos embaixo dela
	supLeft := make([]float64, n)
	supRight := make([]float64, n)
	for i := range n {
		supLeft[i], supRight[i] = math.Inf(1), math.Inf(-1)
	}
	addSupport := func(i int, top, bottom *Rect) {
		supLeft[i] = min(supLeft[i], max(top.pos.X, bottom.pos.X))
		supRight[i] = max(supRight[i], min(top.pos.X+float64(top.width), bottom.pos.X+float64(bottom.width)))
	}

	for range collisionIterations {
		w.collideFloor()
		w.forEachNearbyPair(func(p, q int) {
			e1, e2 := &w.entities[p], &w.entities[q]
			switch resolve(e1, e2) {
			case firstOnTop:
				addSupport(p, e1, e2)
			case secondOnTop:
				addSupport(q, e2, e1)
			}
		})
	}

	// a rotação acumula velocidade angular, então é aplicada uma vez por frame
	for i := range n {
		if supLeft[i] > supRight[i] {
			continue // sem apoio: caindo ou no chão
		}
		e := &w.entities[i]
		if !e.ApplyRotation(supLeft[i], supRight[i], dt) {
			e.Settle()
		}
	}

	w.collideFloor()
}

// forEachNearbyPair chama fn(p, q), com p < q, para cada par de blocos que
// pode estar se tocando. A tela é dividida em células do tamanho do maior
// bloco: blocos que se tocam estão sempre na mesma célula ou em células
// vizinhas, então só esses pares são testados, em vez de todos contra todos.
func (w *World) forEachNearbyPair(fn func(p, q int)) {
	cell := 1.0
	for _, e := range w.entities {
		cell = max(cell, float64(e.width), float64(e.height))
	}
	cols := int(math.Ceil(screenWidth/cell)) + 1
	rows := int(math.Ceil(screenHeight/cell)) + 1

	if len(w.cells) < cols*rows {
		w.cells = make([][]int, cols*rows)
	}
	cells := w.cells[:cols*rows]
	for i := range cells {
		cells[i] = cells[i][:0]
	}

	// blocos fora da tela ficam na célula da borda mais próxima
	w.cellOfID = w.cellOfID[:0]
	for i, e := range w.entities {
		cx := min(max(int(math.Floor(e.pos.X/cell)), 0), cols-1)
		cy := min(max(int(math.Floor(e.pos.Y/cell)), 0), rows-1)
		c := cy*cols + cx
		cells[c] = append(cells[c], i)
		w.cellOfID = append(w.cellOfID, c)
	}

	for p := range w.entities {
		cx, cy := w.cellOfID[p]%cols, w.cellOfID[p]/cols
		for y := max(cy-1, 0); y <= min(cy+1, rows-1); y++ {
			for x := max(cx-1, 0); x <= min(cx+1, cols-1); x++ {
				for _, q := range cells[y*cols+x] {
					if q > p {
						fn(p, q)
					}
				}
			}
		}
	}
}

func (w *World) collideFloor() {
	for i := range w.entities {
		e := &w.entities[i]
		floor := float64(screenHeight - e.height)
		if e.pos.Y > floor {
			e.pos.Y = floor
			e.vel.Y = 0
			e.Settle()
		}
	}
}

type contact int

const (
	noContact contact = iota
	sideContact
	firstOnTop  // e1 está apoiado em e2
	secondOnTop // e2 está apoiado em e1
)

// resolve separa e1 e e2 se estiverem sobrepostos. No contato vertical, quem
// está em cima sobe; no lateral, os dois se empurram de acordo com o peso.
func resolve(e1, e2 *Rect) contact {
	mtv, ok := e1.Penetration(*e2)
	if !ok {
		return noContact
	}

	if mtv.Y != 0 {
		// quem está embaixo não é empurrado, para não afundar no chão
		top, bottom, push, result := e1, e2, mtv, firstOnTop
		if mtv.Y > 0 {
			top, bottom, push, result = e2, e1, mtv.Opposite(), secondOnTop
		}
		top.pos = top.pos.Sum(push)
		if top.vel.Y > bottom.vel.Y {
			top.vel.Y = bottom.vel.Y
		}
		return result
	}

	// contato lateral: cada um anda uma parte proporcional ao peso do outro
	m1, m2 := float64(e1.weight), float64(e2.weight)
	e1.pos = e1.pos.Sum(mtv.Scalar(m2 / (m1 + m2)))
	e2.pos = e2.pos.Sum(mtv.Opposite().Scalar(m1 / (m1 + m2)))

	// se estão se aproximando, passam a andar juntos (conserva o momento)
	if (e1.vel.X-e2.vel.X)*mtv.X < 0 {
		vx := (m1*e1.vel.X + m2*e2.vel.X) / (m1 + m2)
		e1.vel.X, e2.vel.X = vx, vx
	}

	// quem tomba na direção do vizinho fica escorado nele e para de girar
	if e1.angVel*mtv.X < 0 {
		e1.Settle()
	}
	if e2.angVel*mtv.X > 0 {
		e2.Settle()
	}
	return sideContact
}
