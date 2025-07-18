package models

// NewTicket representa la estructura de un nuevo ticket a procesar
type NewTicket struct {
	Data NewTicketData `json:"data"` // Datos del nuevo ticket
}

// NewTicketData contiene los datos para crear un nuevo ticket
type NewTicketData struct {
	TicketData
	// Identificadores y referencias
	SerieIdentificador    string `json:"serie_identificador"`    // ID de la serie (codificado)
	Sucursal              string `json:"sucursal"`               // Código de sucursal
	ClienteIdentificador  string `json:"cliente_identificador"`  // ID del cliente (codificado)
	VendedorIdentificador string `json:"vendedor_identificador"` // ID del vendedor (codificado)

	// Metadatos del ticket
	Enviada BoolFlex `json:"enviada"` // Indica si fue enviado

	// Datos de la sucursal
	SucursalEmail    string `json:"sucursal_email"`
	SucursalLeyenda1 string `json:"sucursal_leyenda_1"` // Leyenda 1
	SucursalLeyenda2 string `json:"sucursal_leyenda_2"` // Leyenda 2

	// Comentarios y metadatos adicionales
	Comentario            interface{} `json:"comentario"`              // Comentario general (puede ser nulo)
	ComentarioInterno     interface{} `json:"comentario_interno"`      // Comentario interno (puede ser nulo)
	AlmacenID             string      `json:"almacen_id"`              // ID del almacén
	TipoConversionFactura string      `json:"tipo_conversion_factura"` // Tipo de conversión a factura

	// Montos y valores económicos
	Costo           string      `json:"costo"`            // Costo total
	CostoBruto      string      `json:"costo_bruto"`      // Costo bruto
	DescuentoMotivo string      `json:"descuento_motivo"` // Motivo del descuento
	MetodoPago      string      `json:"metodo_pago"`      // Método de pago (SAT)
	Abonado         interface{} `json:"abonado"`          // Monto abonado (puede ser nulo)

	// Datos de recetas (si aplica)
	Receta map[string]interface{} `json:"receta"` // Datos de receta (estructura variable)

	// Impuestos
	Impuestos []Impuesto `json:"impuestos"` // Lista de impuestos globales
}
