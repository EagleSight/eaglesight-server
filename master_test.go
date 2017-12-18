package main

import (
	"testing"
)

func testIsReachable(t *testing.T) {

	master, _ := NewMaster("127.0.0.1", "")

	if master.IsReachable() {
		t.Errorf("The master should not be reachable!")
	}

}

func testNew(t *testing.T) {

	url := "http://www.example.com"
	secret := "secret"

	master, err := NewMaster(url, secret)

	if err != nil {
		t.Error(err)
	}

	if master.url != url {
		t.Errorf("master.url is not as expected: %s != %s", master.url, url)
	}

	if master.secret != secret {
		t.Errorf("master.secret is not as expected: %s != %s", master.secret, secret)
	}

}
