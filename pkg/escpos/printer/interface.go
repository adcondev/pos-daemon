package printer

// Printer defines the interface for any ESC/POS compatible printer
type Printer interface {
	Initialize() error
	Pulse(int, int, int) error
	Close() error
	Status() (Status, error)
}

// Status represents the printer status
type Status struct {
	Online      bool
	PaperStatus PaperStatus
	DrawerOpen  bool
	// Additional status fields
}

type PaperStatus int

const (
	PaperOK PaperStatus = iota
	PaperLow
	PaperOut
)
