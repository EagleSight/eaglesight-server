package world

import (
	"encoding/binary"
	"log"
	"math"
	"os"

	"github.com/eaglesight/eaglesight-backend/mathutils"
)

// Terrain ...
type Terrain struct {
	width    uint
	depth    uint
	distance float64
	points   []uint16
}

// LoadTerrain loads the terrain
func LoadTerrain() (Terrain, error) {

	reader, err := os.Open("./map.esmap")
	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	t := Terrain{
		width:    0,
		depth:    0,
		distance: 0.0,
	}

	// We read the header
	header := make([]byte, 2+2+2+4)

	_, err = reader.Read(header[:])

	if err != nil {
		return t, err
	}

	// Load the width
	t.width = uint(binary.LittleEndian.Uint16(header[0:2]))

	// Load the depth
	t.depth = uint(binary.LittleEndian.Uint16(header[2:4]))

	// Load the distance
	t.distance = float64(math.Float32frombits(binary.LittleEndian.Uint32(header[6:10])))

	// Here comes the me
	t.points = make([]uint16, t.width*t.depth)

	data := make([]byte, 2)
	for i := 0; i < len(t.points); i++ {

		_, err := reader.Read(data)

		if err != nil {
			return t, err
		}

		t.points[i] = binary.LittleEndian.Uint16(data)
	}

	return t, nil
}

// OverredTriangle find the triangle that is overred by the vector pos. Return a triangle made of 3 Vector3D
func (t *Terrain) OverredTriangle(pos mathutils.Vector3D) (s [3]mathutils.Vector3D) {

	// 0 1
	// 2 3

	col := uint(math.Ceil(pos.X / t.distance)) // X
	row := uint(math.Ceil(pos.Z / t.distance)) // Z

	// We check if we are out of bound
	if col < 0 || col >= t.width-1 || row < 0 || row >= t.depth-1 {
		s[0].X = math.NaN()
		return s // s[0] == NaN if out of bound
	}

	// UP LEFT
	index0 := row*t.width + col
	s[0].X = float64(index0%t.width) * t.distance
	s[0].Y = float64(t.points[index0])
	s[0].Z = math.Ceil(float64(index0/t.width)) * t.distance

	if math.Mod(pos.X, t.distance) > math.Mod(pos.Z, t.distance) {

		// DOWN RIGHT
		index1 := index0 + 1 + t.width
		s[1].X = s[0].X + t.distance
		s[1].Y = float64(t.points[index1])
		s[1].Z = s[0].Z + t.distance

		// UP RIGHT
		index2 := index0 + 1
		s[2].X = s[1].X
		s[2].Y = float64(t.points[index2])
		s[2].Z = s[0].Z

	} else {

		// DOWN LEFT
		index1 := index0 + t.width
		s[1].X = s[0].X
		s[1].Y = float64(t.points[index1])
		s[1].Z = s[0].Z + t.distance

		// DOWN RIGHT
		index2 := index0 + 1 + t.width
		s[2].X = s[0].X + t.distance
		s[2].Y = float64(t.points[index2])
		s[2].Z = s[1].Z

	}

	return s
}
