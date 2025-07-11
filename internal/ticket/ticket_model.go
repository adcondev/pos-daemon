package ticket

import (
	"encoding/json"
	"fmt"
)

// BoolFlexible Necesitado ya que Go no entiende 0 y 1 como true/false
type BoolFlexible bool

func (b *BoolFlexible) UnmarshalJSON(data []byte) error {
	// Intenta como bool
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		*b = BoolFlexible(boolVal)
		return nil
	}
	// Intenta como string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		if strVal == "1" {
			*b = true
		} else if strVal == "0" {
			*b = false
		} else if strVal == "true" {
			*b = true
		} else if strVal == "false" {
			*b = false
		} else {
			return fmt.Errorf("valor no soportado para BoolFlexible: %s", strVal)
		}
		return nil
	}
	return fmt.Errorf("no se pudo deserializar BoolFlexible: %s", string(data))
}

// Wrapper ---------- envoltorio raíz ----------
type Wrapper struct {
	Data Ticket `json:"data"`
}

// Ticket ---------- ticket principal ----------
type Ticket struct {
	// metadatos
	Identificador string       `json:"identificador"`
	Vendedor      string       `json:"vendedor"`
	Folio         string       `json:"folio"`
	Serie         string       `json:"serie"`
	FechaSistema  string       `json:"fecha_sistema"`  // "DD/MM/YYYY HH:MM:SS"
	TipoOperacion string       `json:"tipo_operacion"` // NOTA_VENTA, FACTURA, etc.
	Anulada       BoolFlexible `json:"anulada"`        // "0"/"1" o "true"/"false" → BoolFlexible

	// montos
	Descuento         float64  `json:"descuento,string"`
	DescuentoNotaCred *float64 `json:"descuento_nota_credito,string,omitempty"`
	Total             float64  `json:"total,string"`
	Saldo             float64  `json:"saldo,string"`
	Pagado            float64  `json:"pagado,string"`
	Cambio            float64  `json:"cambio,string"`

	// Datos relacionales de Cliente
	ClienteNombre         string `json:"cliente"`
	ClienteRFC            string `json:"cliente_rfc"`
	ClienteCP             string `json:"cliente_cp"`
	ClienteUsoCFDI        string `json:"cliente_uso_cfdi"`
	ClienteRegimenFiscal  string `json:"cliente_regimen_fiscal"`
	ClienteCalle          string `json:"cliente_calle"`
	ClienteNumeroExterior string `json:"cliente_numero_exterior"`
	ClienteNumeroInterior string `json:"cliente_numero_interior"`
	ClienteColonia        string `json:"cliente_colonia"`
	ClienteLocalidad      string `json:"cliente_localidad"`
	ClienteMunicipio      string `json:"cliente_delegacion"`
	ClienteEstado         string `json:"cliente_estado"`
	ClientePais           string `json:"cliente_pais"`
	ClienteEmails         string `json:"cliente_emails"`

	// Datos relacionales de Sucursal
	SucursalRFC             string `json:"sucursal_rfc"`
	SucursalNombre          string `json:"sucursal_nombre"`
	SucursalNombreComercial string `json:"sucursal_nombre_comercial"`
	SucursalTienda          string `json:"sucursal_tienda"`
	SucursalRegimenClave    string `json:"sucursal_regimen_clave"`
	SucursalRegimen         string `json:"sucursal_regimen"`
	SucursalCalle           string `json:"sucursal_calle"`
	SucursalNumero          string `json:"sucursal_numero"`
	SucursalNumeroInt       string `json:"sucursal_numero_int"`
	SucursalColonia         string `json:"sucursal_colonia"`
	SucursalLocalidad       string `json:"sucursal_localidad"`
	SucursalMunicipio       string `json:"sucursal_municipio"`
	SucursalEstado          string `json:"sucursal_estado"`
	SucursalCP              string `json:"sucursal_cp"`
	SucursalPais            string `json:"sucursal_pais"`
	SucursalEmail           string `json:"sucursal_email"`

	// movimientos
	Conceptos      []Concepto      `json:"conceptos"`
	DocumentosPago []DocumentoPago `json:"documentos_pago"`
	Pagos          []FormaPago     `json:"pago"` // ← nombre exacto en JSON
}

// Concepto ---------- detalle de conceptos ----------
type Concepto struct {
	Clave                 string       `json:"clave"`
	Descripcion           string       `json:"descripcion"`
	Cantidad              float64      `json:"cantidad,string"`
	Unidad                string       `json:"unidad"`
	PrecioVenta           float64      `json:"precio_venta,string"`
	Total                 float64      `json:"total,string"`
	ClaveProductoServicio string       `json:"clave_producto_servicio"`
	ClaveUnidadSAT        string       `json:"clave_unidad_sat"`
	VentaGranel           BoolFlexible `json:"venta_granel"`
	Impuestos             []Impuesto   `json:"impuestos"`
}

// Impuesto ---------- impuestos (T = trasladado, R = retenido) ----------
type Impuesto struct {
	Factor  string  `json:"factor"` // Tasa / Cuota
	Base    float64 `json:"base,string"`
	Importe float64 `json:"importe,string"`
	Codigo  string  `json:"impuestos"` // 001, 002, 003 ó nombre local
	Tasa    float64 `json:"tasa,string"`
	Entidad string  `json:"entidad"` // Federal / Local
	Tipo    string  `json:"tipo"`    // T / R
}

// DocumentoPago ---------- pagos ----------
type DocumentoPago struct {
	Total      float64      `json:"total,string"`
	TipoCambio float64      `json:"tipo_cambio,string"`
	Saldo      float64      `json:"saldo,string"`
	Nota       string       `json:"nota"`
	Sistema    string       `json:"sistema"`
	Anulado    BoolFlexible `json:"anulado"`
	Cambio     float64      `json:"cambio,string"`
	FechaPago  string       `json:"fecha_pago"`
	FormasPago []FormaPago  `json:"formas_pago"`
}

type FormaPago struct {
	FormaPago     string  `json:"forma_pago"`
	Cantidad      float64 `json:"cantidad,string"`
	Identificador string  `json:"forma_pago_identificador"`
}
