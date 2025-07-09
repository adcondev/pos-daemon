// Package windows proporciona conectores específicos para impresoras
// en sistemas Windows utilizando la API del spooler de impresión.
//
// Este paquete implementa la interfaz PrintConnector para comunicarse
// con impresoras instaladas en Windows usando las funciones de la API
// de Windows Spooler. Permite enviar datos RAW directamente a la impresora
// sin procesamiento adicional del driver.
//
// El conector maneja automáticamente:
//   - Apertura y cierre de conexiones con la impresora
//   - Inicio y finalización de trabajos de impresión
//   - Envío de datos RAW al spooler
//   - Manejo de errores y limpieza de recursos
//
// Ejemplo de uso:
//
//	connector, err := windows.NewWindowsPrintConnector("EC-PM-80250")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer connector.Close()
//
//	data := []byte("Hola mundo")
//	_, err = connector.Write(data)
//	if err != nil {
//	    log.Fatal(err)
//	}
package windows