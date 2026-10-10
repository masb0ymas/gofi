package dtos

// Health is the liveness payload: process identity plus a real database ping.
// It is returned unwrapped (no {data} envelope) because it is a probe endpoint,
// and the web UI reads the fields directly.
type Health struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}
