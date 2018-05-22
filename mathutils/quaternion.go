package mathutils

// Quaternion ...
type Quaternion struct {
	W float64
	X float64
	Y float64
	Z float64
}

// Mul multiplies 2 quaternions
func (q *Quaternion) Mul(r *Quaternion) (s *Quaternion) {
	return &Quaternion{
		W: q.W*r.W - q.X*r.X - q.Y*r.Y - q.Z*r.Z,
		X: q.W*r.X + q.X*r.W + q.Y*r.Z - q.Z*r.Y,
		Y: q.W*r.Y - q.X*r.Z + q.Y*r.W + q.Z*r.X,
		Z: q.W*r.Z + q.X*r.Y - q.Y*r.X + q.Z*r.W,
	}
}
