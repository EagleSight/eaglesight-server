package game

// Connector represent an entry point to all the connections to a server
type Connector interface {
	Start(*Server) error
}
