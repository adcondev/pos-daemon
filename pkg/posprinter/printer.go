package posprinter

import (
	"pos-daemon.adcon.dev/pkg/posprinter/command"
	"pos-daemon.adcon.dev/pkg/posprinter/connector"
	"pos-daemon.adcon.dev/pkg/posprinter/protocol"
)

// GenericPrinter es una implementación que usa Protocol y Connector
type GenericPrinter struct {
	Protocol  protocol.Protocol
	Connector connector.WindowsPrintConnector
	// ...
}

// SetJustification establece la alineación usando el protocolo subyacente
func (p *GenericPrinter) SetJustification(alignment command.Alignment) error {
	cmd := p.Protocol.SetJustification(alignment)
	_, err := p.Connector.Write(cmd)
	return err
}

// Otros métodos...
