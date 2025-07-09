package ticket

// ---------- envoltorio raíz ----------
type Wrapper struct {
	Data Ticket `json:"data"`
}

// ---------- ticket principal ----------
type Ticket struct {
	// metadatos
	Identificador string `json:"identificador"`
	Vendedor      string `json:"vendedor"`
	Folio         string `json:"folio"`
	Serie         string `json:"serie"`
	FechaSistema  string `json:"fecha_sistema"`  // "DD/MM/YYYY HH:MM:SS"
	TipoOperacion string `json:"tipo_operacion"` // NOTA_VENTA, FACTURA, etc.
	Anulada       bool   `json:"anulada,string"` // "0"/"1" → bool

	// montos
	Descuento         float64  `json:"descuento,string"`
	DescuentoNotaCred *float64 `json:"descuento_nota_credito,string,omitempty"`
	Total             float64  `json:"total,string"`
	Saldo             float64  `json:"saldo,string"`
	Pagado            float64  `json:"pagado,string"`
	Cambio            float64  `json:"cambio,string"`

	// datos relacionales
	Cliente  ClienteInfo  `json:"cliente_info,inline"`
	Sucursal SucursalInfo `json:"sucursal_info,inline"`

	// movimientos
	Conceptos      []Concepto      `json:"conceptos"`
	DocumentosPago []DocumentoPago `json:"documentos_pago"`
	Pagos          []FormaPago     `json:"pago"` // ← nombre exacto en JSON
}

// ---------- datos de cliente ----------
type ClienteInfo struct {
	Nombre         string `json:"cliente"`
	RFC            string `json:"cliente_rfc"`
	CP             string `json:"cliente_cp"`
	UsoCFDI        string `json:"cliente_uso_cfdi"`
	RegimenFiscal  string `json:"cliente_regimen_fiscal"`
	Calle          string `json:"cliente_calle"`
	NumeroExterior string `json:"cliente_numero_exterior"`
	NumeroInterior string `json:"cliente_numero_interior"`
	Colonia        string `json:"cliente_colonia"`
	Localidad      string `json:"cliente_localidad"`
	Municipio      string `json:"cliente_delegacion"`
	Estado         string `json:"cliente_estado"`
	Pais           string `json:"cliente_pais"`
	Emails         string `json:"cliente_emails"`
}

// ---------- datos de sucursal ----------
type SucursalInfo struct {
	RFC          string `json:"sucursal_rfc"`
	Nombre       string `json:"sucursal_nombre"`
	RegimenClave string `json:"sucursal_regimen_clave"`
	Calle        string `json:"sucursal_calle"`
	Numero       string `json:"sucursal_numero"`
	NumeroInt    string `json:"sucursal_numero_int"`
	Colonia      string `json:"sucursal_colonia"`
	Localidad    string `json:"sucursal_localidad"`
	Municipio    string `json:"sucursal_municipio"`
	Estado       string `json:"sucursal_estado"`
	CP           string `json:"sucursal_cp"`
	Pais         string `json:"sucursal_pais"`
}

// ---------- detalle de conceptos ----------
type Concepto struct {
	Clave                 string     `json:"clave"`
	Descripcion           string     `json:"descripcion"`
	Cantidad              float64    `json:"cantidad,string"`
	Unidad                string     `json:"unidad"`
	PrecioVenta           float64    `json:"precio_venta,string"`
	Total                 float64    `json:"total,string"`
	ClaveProductoServicio string     `json:"clave_producto_servicio"`
	ClaveUnidadSAT        string     `json:"clave_unidad_sat"`
	VentaGranel           bool       `json:"venta_granel,string"`
	Impuestos             []Impuesto `json:"impuestos"`
}

// ---------- impuestos (T = trasladado, R = retenido) ----------
type Impuesto struct {
	Factor  string  `json:"factor"` // Tasa / Cuota
	Base    float64 `json:"base,string"`
	Importe float64 `json:"importe,string"`
	Codigo  string  `json:"impuestos"` // 001, 002, 003 ó nombre local
	Tasa    float64 `json:"tasa,string"`
	Entidad string  `json:"entidad"` // Federal / Local
	Tipo    string  `json:"tipo"`    // T / R
}

// ---------- pagos ----------
type DocumentoPago struct {
	Total      float64     `json:"total,string"`
	TipoCambio float64     `json:"tipo_cambio,string"`
	Saldo      float64     `json:"saldo,string"`
	Nota       string      `json:"nota"`
	Sistema    string      `json:"sistema"`
	Anulado    bool        `json:"anulado,string"`
	Cambio     float64     `json:"cambio,string"`
	FechaPago  string      `json:"fecha_pago"`
	FormasPago []FormaPago `json:"formas_pago"`
}

type FormaPago struct {
	FormaPago     string  `json:"forma_pago"`
	Cantidad      float64 `json:"cantidad,string"`
	Identificador string  `json:"forma_pago_identificador"`
}
