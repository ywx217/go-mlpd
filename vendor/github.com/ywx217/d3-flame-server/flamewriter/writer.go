package flamewriter

// FlameWriter flame writer interface
type FlameWriter interface {
	Write(root *Record) error
}
