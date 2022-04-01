package transport

import (
	"atlas-cts/configuration"
	"testing"
	"time"
)

// TestGetState1 tests getState
func TestGetState1(t *testing.T) {
	tc := configuration.TransportConfiguration{
		Enabled:            true,
		Source:             101000000,
		Departure:          101000301,
		Transport:          []uint32{200090010, 200090011},
		Arrival:            200000100,
		Destination:        200000000,
		OpenGateDuration:   240000,
		ClosedGateDuration: 60000,
		RideDuration:       600000,
	}
	tm := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	result := getState(tm, tc)
	if result != StatePreparing {
		t.Fatalf("GetState1 expected=%s, got=%s", StatePreparing, result)
	}
}

// TestGetState2 tests getState
func TestGetState2(t *testing.T) {
	tc := configuration.TransportConfiguration{
		Enabled:            true,
		Source:             101000000,
		Departure:          101000301,
		Transport:          []uint32{200090010, 200090011},
		Arrival:            200000100,
		Destination:        200000000,
		OpenGateDuration:   240000,
		ClosedGateDuration: 60000,
		RideDuration:       600000,
	}
	tm := time.Date(2009, 11, 17, 20, 32, 58, 651387237, time.UTC)
	result := getState(tm, tc)
	if result != StateBoarding {
		t.Fatalf("TestGetState2 expected=%s, got=%s", StateBoarding, result)
	}
}

// TestGetState3 tests getState
func TestGetState3(t *testing.T) {
	tc := configuration.TransportConfiguration{
		Enabled:            true,
		Source:             101000000,
		Departure:          101000301,
		Transport:          []uint32{200090010, 200090011},
		Arrival:            200000100,
		Destination:        200000000,
		OpenGateDuration:   240000,
		ClosedGateDuration: 60000,
		RideDuration:       600000,
	}
	tm := time.Date(2009, 11, 17, 20, 37, 58, 651387237, time.UTC)
	result := getState(tm, tc)
	if result != StateInProgress {
		t.Fatalf("TestGetState3 expected=%s, got=%s", StateInProgress, result)
	}
}
