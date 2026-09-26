package main

type Rect struct {
	pos Vec2d

	width  int
	height int

	weight int

	force Vec2d
	vel   Vec2d
}

func NewRect(x, y float64, width, height, weight int) *Rect {
	return &Rect{pos: Vec2d{x, y}, width: width, height: height, weight: weight, force: Vec2d{0, 0}}
}

func (r *Rect) ApplyForce(f Vec2d) {
	r.force = r.force.Sum(f)
}

func (r *Rect) Update(dt float64) {
	f := NatureForces(float64(r.weight))
	r.ApplyForce(f)

	var k float64 = 1 / float64(r.weight)
	acc := r.force.Scalar(k)
	r.vel = r.vel.Sum(acc.Scalar(dt))

	ds := r.vel.Scalar(dt)
	r.pos = r.pos.Sum(ds)

	r.force.Clear()
}
