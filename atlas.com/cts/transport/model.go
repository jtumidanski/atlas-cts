package transport

type Model struct {
	enabled     bool
	source      uint32
	departure   uint32
	transport   []uint32
	arrival     uint32
	destination uint32
	state       string
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

func (m Model) Transport() []uint32 {
	return m.transport
}

func (m Model) Arrival() uint32 {
	return m.arrival
}

func (m Model) Destination() uint32 {
	return m.destination
}

func (m Model) State() string {
	return m.state
}

func (m Model) updateState(state string) Model {
	return Model{
		enabled:     m.enabled,
		source:      m.source,
		departure:   m.departure,
		transport:   m.transport,
		arrival:     m.arrival,
		destination: m.destination,
		state:       state,
	}
}
