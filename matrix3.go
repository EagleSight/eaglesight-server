package main

import "math"

type matrix3 struct {
	_11, _12, _13 float64
	_21, _22, _23 float64
	_31, _32, _33 float64
}

func newMatrix3() matrix3 {
	return matrix3{
		_11: 1,
		_12: 0,
		_13: 0,
		_21: 0,
		_22: 1,
		_23: 0,
		_31: 0,
		_32: 0,
		_33: 1,
	}
}

func (m *matrix3) getInverse() matrix3 {

	n := *m

	t11 := n._33*n._22 - n._23*n._32
	t12 := n._23*n._31 - n._33*n._21
	t13 := n._32*n._21 - n._22*n._31

	det := n._11*t11 + n._12*t12 + n._13*t13

	if det == 0 {

		// Return the identity
		n._11 = 1
		n._12 = 0
		n._13 = 0

		n._21 = 0
		n._22 = 1
		n._23 = 0

		n._31 = 0
		n._32 = 0
		n._33 = 1

		return n
	}

	// Inverse the determinant
	detInv := 1 / det

	n._11 = t11 * detInv
	n._12 = (m._13*m._32 - m._33*m._12) * detInv
	n._13 = (m._23*m._12 - m._13*m._22) * detInv

	n._21 = t12 * detInv
	n._22 = (m._33*m._11 - m._13*m._31) * detInv
	n._23 = (m._13*m._21 - m._23*m._11) * detInv

	n._31 = t13 * detInv
	n._32 = (m._12*m._31 - m._32*m._11) * detInv
	n._33 = (m._22*m._11 - m._12*m._21) * detInv

	return n
}

func makeMatrix3X(angle float64) matrix3 {
	var mat matrix3

	mat._11 = 1
	mat._12 = 0
	mat._13 = 0

	mat._21 = 0
	mat._22 = math.Cos(angle)
	mat._23 = -math.Sin(angle)

	mat._31 = 0
	mat._32 = math.Sin(angle)
	mat._33 = math.Cos(angle)

	return mat
}

func makeMatrix3Y(angle float64) matrix3 {
	var mat matrix3

	mat._11 = math.Cos(angle)
	mat._12 = 0
	mat._13 = math.Sin(angle)

	mat._21 = 0
	mat._22 = 1
	mat._23 = 0

	mat._31 = -math.Sin(angle)
	mat._32 = 0
	mat._33 = math.Cos(angle)

	return mat
}

func makeMatrix3Z(angle float64) matrix3 {
	var mat matrix3

	mat._11 = math.Cos(angle)
	mat._12 = -math.Sin(angle)
	mat._13 = 0

	mat._21 = math.Sin(angle)
	mat._22 = math.Cos(angle)
	mat._23 = 0

	mat._31 = 0
	mat._32 = 0
	mat._33 = 1

	return mat
}

func (m *matrix3) Mul(m2 matrix3) matrix3 {

	var r matrix3

	r._11 = m._11*m2._11 + m._12*m2._21 + m._13*m2._31
	r._12 = m._11*m2._12 + m._12*m2._22 + m._13*m2._32
	r._13 = m._11*m2._13 + m._12*m2._23 + m._13*m2._33

	r._21 = m._21*m2._11 + m._22*m2._21 + m._23*m2._31
	r._22 = m._21*m2._12 + m._22*m2._22 + m._23*m2._32
	r._23 = m._21*m2._13 + m._22*m2._23 + m._23*m2._33

	r._31 = m._31*m2._11 + m._32*m2._21 + m._33*m2._31
	r._32 = m._31*m2._12 + m._32*m2._22 + m._33*m2._32
	r._33 = m._31*m2._13 + m._32*m2._23 + m._33*m2._33

	return r

}

func (m *matrix3) ToEulerAngle() (v Vector3D) {

	v.z = math.Asin(math.Max(-1, math.Min(1, m._21)))

	if math.Abs(m._21) < 0.999999 {
		v.x = math.Atan2(-m._23, m._22)
		v.y = math.Atan2(-m._31, m._11)
	} else {
		v.x = 0
		v.y = math.Atan2(m._13, m._33)
	}

	return v
}
