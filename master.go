package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net"
	"net/http"
)

// Master is an abstraction to the host that coordinate everything
type Master struct {
	url    string
	secret string
}

// NewMaster return a Master
func NewMaster(url string, secret string) (*Master, error) {

	// TODO : Double check the URL

	return &Master{
		url:    url,
		secret: secret,
	}, nil
}

// TEST THIS!
// Get this machines address, or 127.0.0.1
func (m *Master) getSlaveIP() string {
	// Shamelessly stolen from Stack Overlow -> https://stackoverflow.com/a/37382208
	conn, err := net.Dial("udp", "8.8.8.8:80") // TODO: find a way to ping the master
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

// TEST THIS!
func (m *Master) generateAuthHash(route string, secret string) string {

	// Generate the authentication hash (secret + | + slave's IP )
	authHash := sha1.Sum([]byte(route + "|" + secret + "|" + m.getSlaveIP()))

	// Encode it to base64
	return base64.StdEncoding.EncodeToString(authHash[:])
}

// IsReachable detect if there is a master to reach
func (m *Master) IsReachable() bool {

	if m.url == "" {
		return false
	}

	conn, err := net.Dial("udp", m.url)
	if err != nil {
		return false
	}
	conn.Close()

	return true
}

// TEST THIS!
// FetchParameters fetch the game parameters from the master
func (m *Master) FetchParameters() (io.ReadCloser, error) {

	reader, err := m.message("GET", "/game/parameters", []byte{})

	return reader, err
}

// TEST THIS!
// Ready notifies the master that the server is ready
func (m *Master) Ready() error {

	if m.IsReachable() {
		_, err := m.message("POST", "/game/ready", []byte{})
		return err
	}

	return nil

}

// TEST THIS!
// message send an authenticaded request to the master
func (m *Master) message(method string, route string, body []byte) (io.ReadCloser, error) {

	req, err := http.NewRequest(method, m.url+route, bytes.NewReader(body))

	if err != nil {
		return nil, err // TODO: Is nil the right thing to return?
	}

	req.Header.Add("Auth", m.generateAuthHash(route, m.secret))

	client := http.Client{}

	// We send the ip address + secret to the master
	resp, err := client.Do(req)

	if err != nil {
		resp.Body.Close()
		return nil, err // TODO: Is nil the right thing to return?
	}

	return resp.Body, nil
}

// Disengage tells the master to shutdown this droplet
func (m *Master) Disengage() {

}
