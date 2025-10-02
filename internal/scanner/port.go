package scanner

// Port represents information about a port and its associated process
type Port struct {
	Number      int    `json:"port"`
	PID         int    `json:"pid"`
	ProcessName string `json:"process"`
	Command     string `json:"command"`
	Protocol    string `json:"protocol"`
}
