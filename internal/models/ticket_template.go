package models

// TicketTemplate representa la estructura raíz para las plantillas de tickets
type TicketTemplate struct {
	Data TicketTemplateData `json:"data"` // Datos de la plantilla
}

// TicketTemplateData representa la configuración para la impresión de tickets
type TicketTemplateData struct {
	// Configuración de dimensiones y tamaños
	TicketWidth     IntFlex `json:"ticket_width,string"`      // Ancho del ticket en caracteres
	RazonSocialSize int     `json:"razon_social_size,string"` // Tamaño de la razón social
	DatosSize       int     `json:"datos_size,string"`        // Tamaño de los datos generales
	LogoWidth       int     `json:"logo_width,string"`        // Ancho del logotipo

	// Elementos visibles - cabecera
	VerLogotipo      BoolFlex `json:"ver_logotipo"`       // Mostrar logotipo
	VerNombre        BoolFlex `json:"ver_nombre"`         // Mostrar nombre del negocio
	VerNombreC       BoolFlex `json:"ver_nombre_c"`       // Mostrar nombre comercial
	VerRFC           BoolFlex `json:"ver_rfc"`            // Mostrar RFC
	VerDom           BoolFlex `json:"ver_dom"`            // Mostrar domicilio
	VerLeyenda       BoolFlex `json:"ver_leyenda"`        // Mostrar leyenda
	VerRegimen       BoolFlex `json:"ver_regimen"`        // Mostrar régimen fiscal
	VerEmail         BoolFlex `json:"ver_email"`          // Mostrar email
	VerTelefono      BoolFlex `json:"ver_telefono"`       // Mostrar teléfono
	VerNombreCliente BoolFlex `json:"ver_nombre_cliente"` // Mostrar nombre del cliente
	VerFolio         BoolFlex `json:"ver_folio"`          // Mostrar folio
	VerFecha         BoolFlex `json:"ver_fecha"`          // Mostrar fecha
	VerTienda        BoolFlex `json:"ver_tienda"`         // Mostrar nombre de tienda

	// Elementos visibles - detalle
	VerPrecioU       BoolFlex `json:"ver_precio_u"`       // Mostrar precio unitario
	VerCantProductos BoolFlex `json:"ver_cant_productos"` // Mostrar cantidad de productos

	// Elementos visibles - impuestos
	IncluyeImpuestos  BoolFlex `json:"incluye_impuestos"`   // Mostrar impuestos incluidos
	VerImpuestos      BoolFlex `json:"ver_impuestos"`       // Mostrar desglose de impuestos
	VerImpuestosTotal BoolFlex `json:"ver_impuestos_total"` // Mostrar total de impuestos

	// Textos personalizados
	CambiarCabecera    string `json:"cambiar_cabecera"`    // Texto personalizado de cabecera
	CambiarReclamacion string `json:"cambiar_reclamacion"` // Texto para reclamaciones
	CambiarPie         string `json:"cambiar_pie"`         // Texto personalizado de pie
}
