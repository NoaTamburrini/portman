package scanner

// Port represents information about a port and its associated process
type Port struct {
	Number      int    `json:"port"`
	PID         int    `json:"pid"`
	ProcessName string `json:"process"`
	Command     string `json:"command"`
	Protocol    string `json:"protocol"`
	State       string `json:"state"` // LISTEN, ESTABLISHED, or "" for udp/stateless
}

// IsListening reports whether this socket is a listening (server) socket.
func (p Port) IsListening() bool {
	return p.State == "LISTEN"
}
