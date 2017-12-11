package main

import (
	"testing"
)

func testIsReachable(t *testing.T) {

	master, _ := NewMaster("", "")

	if master.IsReachable() {
		t.Errorf("The master should not be reachable!")
	}

}
