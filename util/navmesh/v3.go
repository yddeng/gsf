package navmesh

import "math"

type V3 struct {
	X, Y, Z float64
}

func NewV3(x, y, z float64) *V3 {
	return &V3{X: x, Y: y, Z: z}
}

func (v *V3) Cross(v2 *V3) *V3 {
	return NewV3(v.Y*v2.Z-v.Z*v2.Y, v.Z*v2.X-v.X*v2.Z, v.X*v2.Y-v.Y*v2.X)
}

func (v *V3) Sub(v2 *V3) *V3 {
	return NewV3(v.X-v2.X, v.Y-v2.Y, v.Z-v2.Z)
}

func (v *V3) Distance(target *V3) float64 {
	dx := v.X - target.X
	dy := v.Y - target.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// 旋转一个弧度
func (v *V3) Rotate(radian float64) *V3 {
	sin := math.Sin(radian)
	cos := math.Cos(radian)
	v.X = cos*v.X + sin*v.Y
	v.Y = cos*v.Y - sin*v.X
	return v
}

// 目标点是否在半径为 distance 的圆内
func (v *V3) InRange(target *V3, distance float64) bool {
	x := target.X - v.X
	y := target.Y - v.Y
	return x*x+y*y <= distance*distance
}
