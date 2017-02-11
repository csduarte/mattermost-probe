package probe

// Probe represents a basic probe
type Probe interface {
	Setup() error
	Start() error
}
