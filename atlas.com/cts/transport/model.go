package transport

type Model struct {
	enabled     bool
	source      uint32
	departure   uint32
	transport   uint32
	arrival     uint32
	destination uint32
}

func (m Model) Enabled() bool {
	return m.enabled
}

func (m Model) Source() uint32 {
	return m.source
}

func (m Model) Departure() uint32 {
	return m.departure
}

func (m Model) Transport() uint32 {
	return m.transport
}

func (m Model) Arrival() uint32 {
	return m.arrival
}

func (m Model) Destination() uint32 {
	return m.destination
}
