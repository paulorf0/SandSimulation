package main

import "math"

type Rect struct {
	pos Vec2d

	width  int
	height int

	weight int

	force Vec2d
	vel   Vec2d

	angle  float64 // radianos; positivo = sentido horário na tela
	angVel float64
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

	r.angle += r.angVel * dt

	r.force.Clear()
}

// Settle é chamado quando r está parado sobre um apoio: zera o giro e o
// deslizamento e arredonda o ângulo para o múltiplo de 90° mais próximo.
// Se r caiu de lado, largura e altura são trocadas mantendo o centro.
func (r *Rect) Settle() {
	r.vel.X = 0
	r.angVel = 0
	if r.angle == 0 {
		return
	}

	k := int(math.Round(r.angle / (math.Pi / 2)))
	if k%2 != 0 {
		cx := r.pos.X + float64(r.width)/2
		cy := r.pos.Y + float64(r.height)/2
		r.width, r.height = r.height, r.width
		r.pos = Vec2d{cx - float64(r.width)/2, cy - float64(r.height)/2}
	}
	r.angle = 0
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

// ApplyRotation faz r tombar quando o centro de massa está fora do trecho
// [left, right] em X onde r está apoiado (a união de todos os seus apoios):
// r gira em torno da quina mais próxima e o centro passa a se mover ao redor
// dela. Retorna true se r está tombando.
func (r *Rect) ApplyRotation(left, right, dt float64) bool {
	// já girou 90°: está deitado de lado e não tomba mais
	if math.Abs(r.angle) >= math.Pi/2 {
		return false
	}

	w, h := float64(r.width), float64(r.height)

	cx := r.pos.X + w/2

	// distância horizontal do centro de massa até a quina de apoio
	var d float64
	switch {
	case cx > right:
		d = cx - right
	case cx < left:
		d = cx - left
	default:
		return false
	}

	// torque da gravidade (m·g·d) dividido pelo momento de inércia em torno
	// da quina (eixos paralelos); a massa se cancela
	inertia := (w*w+h*h)/12 + d*d + h*h/4
	r.angVel += gravity * d / inertia * dt

	// velocidade do centro girando em torno da quina, com y para baixo
	r.vel = Vec2d{r.angVel * h / 2, r.angVel * d}
	return true
}
