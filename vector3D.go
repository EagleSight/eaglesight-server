package main

// Vector3D ...
type Vector3D struct {
	x, y, z float64
}

func (v *Vector3D) multiplyByMatrix3(m *matrix3) Vector3D {

	var p Vector3D

	p.x = m._11*v.x + m._12*v.y + m._13*v.z
	p.y = m._21*v.x + m._22*v.y + m._23*v.z
	p.z = m._31*v.x + m._32*v.y + m._33*v.z

	return p
}

// Sub substract 2 vectors
func (v *Vector3D) Sub(b *Vector3D) (c Vector3D) {

	c.x = v.x - b.x
	c.y = v.y - b.y
	c.z = v.z - b.z

	return c
}

// Add adds two vectors together
func (v *Vector3D) Add(b *Vector3D) (c Vector3D) {

	c.x = v.x + b.x
	c.y = v.y + b.y
	c.z = v.z + b.z

	return c
}

// MulScalar scale the vector by a sigle number
func (v *Vector3D) MulScalar(b float64) (c Vector3D) {

	c.x = v.x * b
	c.y = v.y * b
	c.z = v.z * b

	return c
}

// CrossProduct returns the cross product of u and v
func CrossProduct(u *Vector3D, v *Vector3D) (d Vector3D) {

	d.x = u.y*v.z - u.z*v.y
	d.y = u.z*v.x - u.x*v.z
	d.z = u.x*v.y - u.y*v.x

	return d
}

func heightOnTriangle(p Vector3D, t *[3]Vector3D) float64 {

	v := t[1].Sub(&t[0])
	w := t[2].Sub(&t[0])

	n := CrossProduct(&v, &w)

	y := (n.x*(t[0].x-p.x) + n.z*(t[0].z-p.z) + n.y*(t[0].y)) / n.y

	return y

}

// Find the height of the heighest point in triangle
func highestInTriangle(triangle [3]Vector3D) (h float64) {

	h = 0

	for _, point := range triangle {
		if point.y > h {
			h = point.y
		}
	}

	return h

}
