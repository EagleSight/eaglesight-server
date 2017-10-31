package main

import (
	"encoding/binary"
	"math"
	"os"
)

// ArenaMap ...
type ArenaMap struct {
	width    uint
	depth    uint
	distance float64
	points   []uint16
}

// LoadArenaMap loads the map from a file
func LoadArenaMap() ArenaMap {

	am := ArenaMap{
		width:    0,
		depth:    0,
		distance: 0.0,
	}

	f, err := os.Open("./map.esmap")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// We read the header
	header := make([]byte, 2+2+2+4)

	_, err = f.Read(header[:])

	if err != nil {
		panic(err)
	}

	// Load the width
	am.width = uint(binary.LittleEndian.Uint16(header[0:2]))

	// Load the depth
	am.depth = uint(binary.LittleEndian.Uint16(header[2:4]))

	// Load the distance
	am.distance = float64(math.Float32frombits(binary.LittleEndian.Uint32(header[6:10])))

	// Here comes the me
	am.points = make([]uint16, am.width*am.depth)

	data := make([]byte, 2)
	for i := 0; i < len(am.points); i++ {

		_, err := f.Read(data)

		if err != nil {
			panic(err)
		}

		am.points[i] = binary.LittleEndian.Uint16(data)
	}

	return am

}

// OverredTriangle find the triangle that is overred by the vector pos. Return a triangle made of 3 Vector3D
func (am *ArenaMap) OverredTriangle(pos Vector3D) (s [3]Vector3D) {

	// 0 1
	// 2 3

	col := uint(math.Ceil(pos.x / am.distance)) // X
	row := uint(math.Ceil(pos.z / am.distance)) // Z

	// We check if we are out of bound
	if col < 0 || col >= am.width-1 || row < 0 || row >= am.depth-1 {
		s[0].x = math.NaN()
		return s // s[0] == NaN if out of bound
	}

	// UP LEFT
	index0 := row*am.width + col
	s[0].x = float64(index0%am.width) * am.distance
	s[0].y = float64(am.points[index0])
	s[0].z = math.Ceil(float64(index0/am.width)) * am.distance

	if math.Mod(pos.x, am.distance) > math.Mod(pos.z, am.distance) {

		// DOWN RIGHT
		index1 := index0 + 1 + am.width
		s[1].x = s[0].x + am.distance
		s[1].y = float64(am.points[index1])
		s[1].z = s[0].z + am.distance

		// UP RIGHT
		index2 := index0 + 1
		s[2].x = s[1].x
		s[2].y = float64(am.points[index2])
		s[2].z = s[0].z

	} else {

		// DOWN LEFT
		index1 := index0 + am.width
		s[1].x = s[0].x
		s[1].y = float64(am.points[index1])
		s[1].z = s[0].z + am.distance

		// DOWN RIGHT
		index2 := index0 + 1 + am.width
		s[2].x = s[0].x + am.distance
		s[2].y = float64(am.points[index2])
		s[2].z = s[1].z

	}

	return s
}
