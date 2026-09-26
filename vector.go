package main

type Vec2d struct {
	X float64
	Y float64
}

func (v Vec2d) Sum(u Vec2d) Vec2d {
	return Vec2d{v.X + u.X, v.Y + u.Y}
}

func (v Vec2d) Scalar(u float64) Vec2d {
	return Vec2d{v.X * u, v.Y * u}
}

func (v Vec2d) Opposite() Vec2d {
	return Vec2d{-v.X, -v.Y}
}

func (v *Vec2d) Clear() {
	v.X = 0
	v.Y = 0
}
