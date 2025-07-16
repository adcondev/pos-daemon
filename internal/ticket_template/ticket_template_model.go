package ticket_template

import (
	adt "pos-daemon.adcon.dev/internal"
)

type TicketTemplate struct {
	Data TicketTemplateData `json:"data"`
}

// TicketTemplateData represents the configuration for ticket printing
type TicketTemplateData struct {
	TicketWidth        adt.IntFlex  `json:"ticket_width,string"`
	RazonSocialSize    int          `json:"razon_social_size,string"` // Falta implementar
	DatosSize          int          `json:"datos_size,string"`        // Falta implementar
	VerLogotipo        adt.BoolFlex `json:"ver_logotipo"`
	LogoWidth          int          `json:"logo_width,string"` // Falta implmentar
	VerNombre          adt.BoolFlex `json:"ver_nombre"`
	VerNombreC         adt.BoolFlex `json:"ver_nombre_c"`
	VerRFC             adt.BoolFlex `json:"ver_rfc"`
	VerDom             adt.BoolFlex `json:"ver_dom"`
	VerLeyenda         adt.BoolFlex `json:"ver_leyenda"` // Falta saber bien que es la leyenda
	VerRegimen         adt.BoolFlex `json:"ver_regimen"`
	VerEmail           adt.BoolFlex `json:"ver_email"`
	VerNombreCliente   adt.BoolFlex `json:"ver_nombre_cliente"`
	VerFolio           adt.BoolFlex `json:"ver_folio"`
	VerFecha           adt.BoolFlex `json:"ver_fecha"`
	VerTienda          adt.BoolFlex `json:"ver_tienda"`
	CambiarCabecera    string       `json:"cambiar_cabecera"`
	VerPrecioU         adt.BoolFlex `json:"ver_precio_u"`
	IncluyeImpuestos   adt.BoolFlex `json:"incluye_impuestos"` // Falta saber que es
	VerCantProductos   adt.BoolFlex `json:"ver_cant_productos"`
	CambiarReclamacion string       `json:"cambiar_reclamacion"`
	VerTelefono        adt.BoolFlex `json:"ver_telefono"`
	CambiarPie         string       `json:"cambiar_pie"`
	VerImpuestos       adt.BoolFlex `json:"ver_impuestos"`       // Saber en que difiere con IncluyeImpuestos
	VerImpuestosTotal  adt.BoolFlex `json:"ver_impuestos_total"` // Saber en que difiere con IncluyeImpuestos y VerImpuestos
}
