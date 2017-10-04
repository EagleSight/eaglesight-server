package main

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"
)

func BenchmarkUpdatePlane(b *testing.B) {
	b.StopTimer()
	buf := new(bytes.Buffer)
	plane := NewPlane(666)
	deltaT := time.Now()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		deltaT.Add(time.Second / 90)
		binary.Write(buf, binary.BigEndian, uint8(0x3))
		binary.Write(buf, binary.BigEndian, uint32(3))
		binary.Write(buf, binary.BigEndian, uint16(1))
		plane.UpdateIntoBuffer(buf, []byte{0x3, 5, 5, 5, 5}, deltaT)
		buf.Reset()
	}

}
