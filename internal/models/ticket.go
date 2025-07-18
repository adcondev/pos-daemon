package models

// Ticket representa la estructura raíz de un ticket de venta
type Ticket struct {
	Data TicketData `json:"data"` // Datos del ticket
}

// TicketData contiene todos los datos de un ticket de venta
type TicketData struct {
	// Metadatos del ticket
	Identificador string   `json:"identificador"`  // ID único del ticket
	Vendedor      string   `json:"vendedor"`       // Nombre del vendedor
	Folio         string   `json:"folio"`          // Número de folio
	Serie         string   `json:"serie"`          // Serie del folio
	FechaSistema  string   `json:"fecha_sistema"`  // Fecha y hora (DD/MM/YYYY HH:MM:SS)
	TipoOperacion string   `json:"tipo_operacion"` // NOTA_VENTA, FACTURA, etc.
	Anulada       BoolFlex `json:"anulada"`        // Indica si el ticket está anulado

	// Montos del ticket
	Descuento         float64  `json:"descuento,string"`                        // Monto de descuento aplicado
	DescuentoNotaCred *float64 `json:"descuento_nota_credito,string,omitempty"` // Descuento por nota de crédito
	Total             float64  `json:"total,string"`                            // Monto total del ticket
	Saldo             float64  `json:"saldo,string"`                            // Saldo pendiente
	Pagado            float64  `json:"pagado,string"`                           // Monto pagado
	Cambio            float64  `json:"cambio,string"`                           // Cambio entregado al cliente

	Cliente

	Sucursal // Datos de la sucursal

	// Enlaces y códigos QR
	AutofacturaLink   string `json:"autofactura_link"`    // Enlace para autofacturación
	AutofacturaLinkQr string `json:"autofactura_link_qr"` // Enlace QR para autofacturación

	// Conceptos y pagos
	Conceptos      []Concepto      `json:"conceptos"`       // Lista de productos o servicios
	DocumentosPago []DocumentoPago `json:"documentos_pago"` // Documentos de pago
	Pagos          []Pago          `json:"pago"`            // Formas de pago utilizadas
}

type Cliente struct {
	// Datos del cliente
	ClienteNombre         string `json:"cliente"`                 // Nombre del cliente
	ClienteRFC            string `json:"cliente_rfc"`             // RFC del cliente
	ClienteCP             string `json:"cliente_cp"`              // Código postal del cliente
	ClienteUsoCFDI        string `json:"cliente_uso_cfdi"`        // Uso del CFDI (SAT)
	ClienteRegimenFiscal  string `json:"cliente_regimen_fiscal"`  // Régimen fiscal del cliente
	ClienteCalle          string `json:"cliente_calle"`           // Calle del domicilio
	ClienteNumeroExterior string `json:"cliente_numero_exterior"` // Número exterior
	ClienteNumeroInterior string `json:"cliente_numero_interior"` // Número interior
	ClienteColonia        string `json:"cliente_colonia"`         // Colonia
	ClienteLocalidad      string `json:"cliente_localidad"`       // Localidad
	ClienteMunicipio      string `json:"cliente_delegacion"`      // Municipio o delegación
	ClienteEstado         string `json:"cliente_estado"`          // Estado
	ClientePais           string `json:"cliente_pais"`            // País
	ClienteEmail          string `json:"cliente_emails"`          // Correos electrónicos
}

type Sucursal struct {
	// Datos de la sucursal
	SucursalRFC             string `json:"sucursal_rfc"`              // RFC de la sucursal
	SucursalNombre          string `json:"sucursal_nombre"`           // Razón social
	SucursalNombreComercial string `json:"sucursal_nombre_comercial"` // Nombre comercial
	SucursalTienda          string `json:"sucursal_tienda"`           // Nombre de la tienda
	SucursalRegimenClave    string `json:"sucursal_regimen_clave"`    // Clave del régimen fiscal
	SucursalRegimen         string `json:"sucursal_regimen"`          // Régimen fiscal descriptivo
	SucursalCalle           string `json:"sucursal_calle"`            // Calle
	SucursalNumero          string `json:"sucursal_numero"`           // Número exterior
	SucursalNumeroInt       string `json:"sucursal_numero_int"`       // Número interior
	SucursalColonia         string `json:"sucursal_colonia"`          // Colonia
	SucursalLocalidad       string `json:"sucursal_localidad"`        // Localidad
	SucursalMunicipio       string `json:"sucursal_municipio"`        // Municipio
	SucursalEstado          string `json:"sucursal_estado"`           // Estado
	SucursalCP              string `json:"sucursal_cp"`               // Código postal
	SucursalPais            string `json:"sucursal_pais"`             // País
	SucursalEmails          string `json:"sucursal_emails"`           // Correos electrónicos
	SucursalTelefono        string `json:"sucursal_telefono"`         // Teléfono
}

// Concepto representa un producto o servicio en el ticket
type Concepto struct {
	Clave                 string     `json:"clave"`                   // Código del producto
	Descripcion           string     `json:"descripcion"`             // Descripción del producto
	Cantidad              float64    `json:"cantidad,string"`         // Cantidad vendida
	Unidad                string     `json:"unidad"`                  // Unidad de medida
	PrecioVenta           float64    `json:"precio_venta,string"`     // Precio unitario
	Total                 float64    `json:"total,string"`            // Total por concepto
	ClaveProductoServicio string     `json:"clave_producto_servicio"` // Clave SAT
	ClaveUnidadSAT        string     `json:"clave_unidad_sat"`        // Clave de unidad SAT
	VentaGranel           BoolFlex   `json:"venta_granel"`            // Indica si es venta a granel
	Impuestos             []Impuesto `json:"impuestos"`               // Impuestos aplicados
	Series                []string   `json:"series,omitempty"`        // Números de serie (opcional)
}

// Impuesto representa un impuesto aplicado a un concepto
type Impuesto struct {
	Factor  string  `json:"factor"`         // Tasa o Cuota
	Base    float64 `json:"base,string"`    // Base del impuesto
	Importe float64 `json:"importe,string"` // Importe calculado
	Codigo  string  `json:"impuestos"`      // Código (001, 002, 003) o nombre
	Tasa    float64 `json:"tasa,string"`    // Tasa aplicada
	Entidad string  `json:"entidad"`        // Federal o Local
	Tipo    string  `json:"tipo"`           // T (trasladado) o R (retenido)
}

// DocumentoPago representa un documento de pago asociado al ticket
type DocumentoPago struct {
	Total      float64  `json:"total,string"`       // Total del documento
	TipoCambio float64  `json:"tipo_cambio,string"` // Tipo de cambio aplicado
	Saldo      float64  `json:"saldo,string"`       // Saldo pendiente
	Nota       string   `json:"nota"`               // Nota adicional
	Sistema    string   `json:"sistema"`            // Fecha y hora del sistema
	Anulado    BoolFlex `json:"anulado"`            // Indica si está anulado
	Cambio     float64  `json:"cambio,string"`      // Cambio entregado
	FechaPago  string   `json:"fecha_pago"`         // Fecha del pago
	FormasPago []Pago   `json:"formas_pago"`        // Formas de pago utilizadas
}

// FormaPago representa una forma de pago utilizada en el ticket
type Pago struct {
	FormaPago     string  `json:"forma_pago"`               // Descripción de la forma de pago
	Cantidad      float64 `json:"cantidad,string"`          // Cantidad pagada
	Identificador string  `json:"forma_pago_identificador"` // ID de la forma de pago
}
