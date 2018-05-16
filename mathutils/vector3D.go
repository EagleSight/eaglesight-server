package mathutils

// Vector3D ...
type Vector3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// MultiplyByMatrix3 multiply a Vector3D by a Matrix3
func (v *Vector3D) MultiplyByMatrix3(m *Matrix3) (p Vector3D) {
	p.X = m._11*v.X + m._12*v.Y + m._13*v.Z
	p.Y = m._21*v.X + m._22*v.Y + m._23*v.Z
	p.Z = m._31*v.X + m._32*v.Y + m._33*v.Z
	return p
}

// Sub substract 2 vectors
func (v *Vector3D) Sub(b Vector3D) (c Vector3D) {
	c.X = v.X - b.X
	c.Y = v.Y - b.Y
	c.Z = v.Z - b.Z
	return c
}

// Add adds two vectors together
func (v *Vector3D) Add(b Vector3D) (c Vector3D) {
	c.X = v.X + b.X
	c.Y = v.Y + b.Y
	c.Z = v.Z + b.Z
	return c
}

// MulScalar scale the vector by a sigle number
func (v *Vector3D) MulScalar(x float64) (c Vector3D) {
	c.X = v.X * x
	c.Y = v.Y * x
	c.Z = v.Z * x
	return c
}

// DivScalar divide the vector by a sigle number
func (v *Vector3D) DivScalar(x float64) (c Vector3D) {
	c.X = v.X / x
	c.Y = v.Y / x
	c.Z = v.Z / x
	return c
}

// CrossProduct returns the cross product of u and v
func CrossProduct(u *Vector3D, v *Vector3D) (d Vector3D) {
	d.X = u.Y*v.Z - u.Z*v.Y
	d.Y = u.Z*v.X - u.X*v.Z
	d.Z = u.X*v.Y - u.Y*v.X
	return d
}

// HeightOnTriangle calculate height of point on a triangle
func HeightOnTriangle(p Vector3D, t *[3]Vector3D) float64 {
	v := t[1].Sub(t[0])
	w := t[2].Sub(t[0])
	n := CrossProduct(&v, &w)
	y := (n.X*(t[0].X-p.X) + n.Z*(t[0].Z-p.Z) + n.Y*(t[0].Y)) / n.Y
	return y
}

// HighestInTriangle Find the height of the heighest point in triangle
func HighestInTriangle(triangle *[3]Vector3D) (h float64) {
	h = triangle[0].Y

	for _, point := range triangle[1:] {
		if point.Y > h {
			h = point.Y
		}
	}
	return h
}
