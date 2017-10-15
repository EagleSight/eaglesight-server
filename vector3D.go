package main

type vector3D struct {
	x, y, z float64
}

func (v *vector3D) multiplyByMatrix3(m *matrix3) vector3D {

	var p vector3D

	p.x = m._11*v.x + m._12*v.y + m._13*v.z
	p.y = m._21*v.x + m._22*v.y + m._23*v.z
	p.z = m._31*v.x + m._32*v.y + m._33*v.z

	return p
}
