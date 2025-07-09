package ticket

type Wrapper struct {
	Data Ticket `json:"data"`
}

type Ticket struct {
	Identificador string `json:"identificador"`
	Vendedor      string `json:"vendedor"`
	Folio         string `json:"folio"`
	Serie         string `json:"serie"`
	FechaSistema  string `json:"fecha_sistema"`
	TipoOperacion string `json:"tipo_operacion"`
	Anulada       bool   `json:"anulada"`

	Cliente  ClienteInfo  `json:"cliente_info"`
	Sucursal SucursalInfo `json:"sucursal_info"`

	Descuento float64 `json:"descuento,string"`
	Total     float64 `json:"total,string"`
	Saldo     float64 `json:"saldo,string"`
	Pagado    float64 `json:"pagado,string"`
	Cambio    float64 `json:"cambio,string"`

	Conceptos      []Concepto      `json:"conceptos"`
	DocumentosPago []DocumentoPago `json:"documentos_pago"`
	Pagos          []Pago          `json:"pagos"`
}

type ClienteInfo struct {
	Nombre        string `json:"cliente"`
	RFC           string `json:"cliente_rfc"`
	CP            string `json:"cliente_cp"`
	UsoCFDI       string `json:"cliente_uso_cfdi"`
	RegimenFiscal string `json:"cliente_regimen_fiscal"`
	Calle         string `json:"cliente_calle"`
	NumeroExt     string `json:"cliente_numero_ext"`
	NumeroInt     string `json:"cliente_numero_int"`
	Colonia       string `json:"cliente_colonia"`
	Localidad     string `json:"cliente_localidad"`
	Municipio     string `json:"cliente_delegacion"`
}

type SucursalInfo struct {
	Nombre       string `json:"sucursal_nombre"`
	RFC          string `json:"sucursal_rfc"`
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

type Concepto struct {
	Clave              string     `json:"clave"`
	Descripcion        string     `json:"descripcion"`
	Cantidad           float64    `json:"cantidad,string"`
	Unidad             string     `json:"unidad"`
	PrecioVenta        float64    `json:"precio_venta,string"`
	Total              float64    `json:"total,string"`
	ClaveProdServicion string     `json:"clave_prod_servicion"`
	ClaveUnidadSAT     string     `json:"clave_unidad_sat"`
	VentaGranel        bool       `json:"venta_granel,string"`
	Impuestos          []Impuesto `json:"impuestos"`
}

type Impuesto struct {
	Factor  string  `json:"factor"`
	Base    float64 `json:"base,string"`
	Importe float64 `json:"importe,string"`
	Codigo  float64 `json:"impuestos"`
	Tasa    float64 `json:"tasa,string"`
	Entidad string  `json:"entidad"`
	Tipo    string  `json:"tipo"`
}

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
	Metodo        string  `json:"forma_pago"`
	Cantidad      float64 `json:"cantidad,string"`
	Identificador string  `json:"forma_pago_identificador"`
}

type Pago = FormaPago
