package main

type PlayerProfile struct {
	Token       uint32            `json:"token"`
	FlightProps *PlaneFlightProps `json:"flightProps"`
}

func DefaultPlayerProfile(uid uint32) PlayerProfile {
	return PlayerProfile{
		Token: uid,
		// TODO: Add default flight props
	}
}
