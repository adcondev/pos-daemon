// Package ticket proporciona estructuras y funcionalidades para el manejo
// de tickets de venta del sistema POS.
//
// Este paquete define los modelos de datos para tickets de venta incluyendo
// información del cliente, sucursal, conceptos, impuestos y formas de pago.
// También incluye funcionalidades para parsear y serializar tickets desde
// y hacia JSON.
//
// Los tickets manejan información completa de ventas incluyendo:
//   - Metadatos del ticket (folio, serie, fecha, vendedor)
//   - Información del cliente y sucursal
//   - Conceptos/productos vendidos con sus impuestos
//   - Formas de pago y documentos de pago
//   - Cálculos de totales, descuentos e IVA
package ticket