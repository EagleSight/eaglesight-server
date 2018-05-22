package mathutils

import "math"

// Matrix3 is a 3x3 matrix
type Matrix3 struct {
	_11, _12, _13 float64
	_21, _22, _23 float64
	_31, _32, _33 float64
}

// NewMatrix3 create a new Matrix3
func NewMatrix3() Matrix3 {

	return Matrix3{
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

// Inverse put the inverse of m into n
func (m *Matrix3) Inverse(n *Matrix3) {

	t11 := m._33*m._22 - m._23*m._32
	t12 := m._23*m._31 - m._33*m._21
	t13 := m._32*m._21 - m._22*m._31
	det := m._11*t11 + m._12*t12 + m._13*t13

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
}

// MakeMatrix3X generate the X component of a rotation matrix
func MakeMatrix3X(angle float64) (mat Matrix3) {

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

// MakeMatrix3Y generate the Y component of a rotation matrix
func MakeMatrix3Y(angle float64) (mat Matrix3) {

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

// MakeMatrix3Z generate the Z component of a rotation matrix
func MakeMatrix3Z(angle float64) (mat Matrix3) {

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

// Mul multiplies a Matrix3 by m2 (also a Matrix3)
func (m *Matrix3) Mul(m2 Matrix3) (r Matrix3) {

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

// ToEulerAngle convert a Matrix3 to Euler's angles
func (m *Matrix3) ToEulerAngle() (v Vector3D) {

	v.Z = math.Asin(math.Max(-1, math.Min(1, m._21)))

	if math.Abs(m._21) < 0.999999 {
		v.X = math.Atan2(-m._23, m._22)
		v.Y = math.Atan2(-m._31, m._11)

	} else {
		v.X = 0
		v.Y = math.Atan2(m._13, m._33)
	}
	return v
}

// ToQuaternion convert a rotation matrix3 to quaternion
func (m *Matrix3) ToQuaternion() (q Quaternion) {
	tr := m._11 + m._22 + m._33

	if tr > 0 {
		S := math.Sqrt(tr+1.0) * 2 // S=4*qw
		q.W = 0.25 * S
		q.X = (m._32 - m._23) / S
		q.Y = (m._13 - m._31) / S
		q.Z = (m._21 - m._12) / S
	} else if (m._11 > m._22) && (m._11 > m._33) {
		S := math.Sqrt(1.0+m._11-m._22-m._33) * 2 // S=4*q.X
		q.W = (m._32 - m._23) / S
		q.X = 0.25 * S
		q.Y = (m._12 + m._21) / S
		q.Z = (m._13 + m._31) / S
	} else if m._22 > m._33 {
		S := math.Sqrt(1.0+m._22-m._11-m._33) * 2 // S=4*q.Y
		q.W = (m._13 - m._31) / S
		q.X = (m._12 + m._21) / S
		q.Y = 0.25 * S
		q.Z = (m._23 + m._32) / S
	} else {
		S := math.Sqrt(1.0+m._33-m._11-m._22) * 2 // S=4*q.Z
		q.W = (m._21 - m._12) / S
		q.X = (m._13 + m._31) / S
		q.Y = (m._23 + m._32) / S
		q.Z = 0.25 * S
	}
	return
}
