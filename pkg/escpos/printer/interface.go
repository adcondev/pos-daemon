package printer

// Printer defines the interface for any ESC/POS compatible printer
type Printer interface {
	Initialize() error
	Disconnect() error
	Write(data []byte) (int, error)
	Print(text string) error
	Feed(lines int) error
	Cut() error
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
