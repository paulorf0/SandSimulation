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

func (r *Rect) ItsCollision(q Rect) bool {
	return r.pos.X < q.pos.X+float64(q.width) &&
		r.pos.X+float64(r.width) > q.pos.X &&
		r.pos.Y < q.pos.Y+float64(q.height) &&
		r.pos.Y+float64(r.height) > q.pos.Y
}

// Penetration retorna o deslocamento mínimo que r precisa fazer para sair de q,
// no eixo de menor sobreposição. Se o Y do vetor for negativo, r está em cima de q.
func (r *Rect) Penetration(q Rect) (Vec2d, bool) {
	rw, rh := float64(r.width), float64(r.height)
	qw, qh := float64(q.width), float64(q.height)

	overlapX := min(r.pos.X+rw, q.pos.X+qw) - max(r.pos.X, q.pos.X)
	overlapY := min(r.pos.Y+rh, q.pos.Y+qh) - max(r.pos.Y, q.pos.Y)
	if overlapX <= 0 || overlapY <= 0 {
		return Vec2d{}, false
	}

	if overlapX < overlapY {
		if r.pos.X+rw/2 < q.pos.X+qw/2 {
			return Vec2d{-overlapX, 0}, true
		}
		return Vec2d{overlapX, 0}, true
	}

	if r.pos.Y+rh/2 < q.pos.Y+qh/2 {
		return Vec2d{0, -overlapY}, true
	}
	return Vec2d{0, overlapY}, true
}
