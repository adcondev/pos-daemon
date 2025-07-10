# POS Daemon - Sistema de Impresión para Puntos de Venta

[![Go Version](https://img.shields.io/badge/Go-1.24.4+-blue.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey.svg)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Sistema completo de impresión para puntos de venta que soporta impresoras térmicas ESC/POS en Windows. Incluye una biblioteca Go completa y comandos de demostración para imprimir tickets de venta.

## 🚀 Características

- **Biblioteca ESC/POS completa**: Soporte para comandos estándar ESC/POS
- **Conector Windows**: Integración directa con el spooler de Windows
- **Codificación de caracteres**: Soporte para CP858 con caracteres latinos
- **Múltiples métodos de impresión**: API de alto nivel y acceso RAW directo
- **Manejo de tickets**: Estructuras para tickets de venta complejos
- **Documentación completa**: Toda la API documentada en español

### Funcionalidades de impresión soportadas

- ✅ Texto con formato (negrita, subrayado, justificación)
- ✅ Códigos de barras (UPC-A, UPC-E, EAN13, EAN8, Code39, ITF, Codabar, Code93, Code128)
- ✅ Códigos QR con múltiples niveles de corrección
- ✅ Códigos PDF417 (estándar y truncado)
- ✅ Control de papel (avance, retroceso, corte)
- ✅ Apertura de cajón portamonedas
- 🔧 Impresión de imágenes (en desarrollo)

## 📋 Requisitos del sistema

- **Sistema operativo**: Windows 10/11
- **Go**: Versión 1.24.4 o superior
- **Impresora**: Compatible con ESC/POS (ej: EC-PM-80250)
- **Permisos**: Acceso de administrador para instalación de drivers

## 🛠 Instalación

### 1. Instalación de drivers de impresora

Para la impresora [EC-PM-80250](https://eclinepos.com/Producto.php?categoria=Impresoras&&buscar=EC-PM-80250):

1. Descarga los drivers: [Driver-2022.zip](https://eclinepos.com/Descargas/ControladoresZip/Impresoras/EC-PM-80250/Driver-2022.zip)
2. Ejecuta el instalador como administrador
3. Selecciona **Install USB Virtual Serial Port Driver**
4. Instala **Install Printer Driver (N)** → **POS Printer 300DPI Series**
5. En **Dispositivos e impresoras**, configura:
   - Puerto: **USB001 – UnknownPrinter**
   - Nombre: **EC-PM-80250**

### 2. Instalación del proyecto

```bash
git clone https://github.com/AdConDev/pos-daemon.git
cd pos-daemon
go mod download
```

### 3. Compilación

```bash
# Compilar todos los comandos
go build ./cmd/print
go build ./cmd/raw_print

# O compilar para Windows desde otro sistema
GOOS=windows go build ./cmd/print
GOOS=windows go build ./cmd/raw_print
```

## 🎯 Uso rápido

### Comando print (recomendado)

Utiliza la biblioteca ESC/POS de alto nivel:

```bash
go run cmd/print/main.go -printer "EC-PM-80250" -debug
```

### Comando raw_print

Acceso directo a la API de Windows para máximo control:

```bash
go run cmd/raw_print/main.go
```

### Usando la biblioteca en tu código

```go
package main

import (
    "log"
    "pos-daemon.adcon.dev/internal/platform/windows"
    "pos-daemon.adcon.dev/pkg/escpos"
)

func main() {
    // Crear conector Windows
    connector, err := windows.NewWindowsPrintConnector("EC-PM-80250")
    if err != nil {
        log.Fatal(err)
    }
    defer connector.Close()

    // Crear impresora
    printer, err := escpos.NewPrinter(connector, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Imprimir ticket
    printer.SetJustification(escpos.JUSTIFY_CENTER)
    printer.Text("Mi Empresa\n")
    printer.SetJustification(escpos.JUSTIFY_LEFT)
    printer.Text("Producto: $10.00\n")
    printer.Cut(escpos.CUT_FULL, 0)
}
```

## 📚 Estructura del proyecto

```
pos-daemon/
├── cmd/
│   ├── print/          # Comando de impresión con biblioteca ESC/POS
│   └── raw_print/      # Comando de impresión RAW directa
├── internal/
│   ├── config/         # Configuración y parseo de flags
│   ├── ticket/         # Modelos de datos de tickets
│   └── platform/
│       └── windows/    # Conector específico para Windows
├── pkg/
│   └── escpos/         # Biblioteca ESC/POS principal
└── docs/               # Documentación adicional
```

### Paquetes principales

- **`escpos`**: Biblioteca principal ESC/POS con soporte completo de comandos
- **`windows`**: Conector que utiliza la API del spooler de Windows
- **`ticket`**: Modelos para manejo de tickets de venta con JSON
- **`config`**: Utilidades de configuración y parseo de flags

## 🔧 Configuración avanzada

### Perfiles de capacidad

```go
profile := &escpos.CapabilityProfile{
    SupportsBarcodeB:   true,
    SupportsQrCode:     true,
    SupportsPdf417Code: true,
    CodePages: map[int]string{
        0: "CP437",
        3: "CP858",
    },
}

printer, err := escpos.NewPrinter(connector, profile)
```

### Manejo de tickets complejos

```go
ticket := &ticket.Ticket{
    Identificador: "TKT001",
    Vendedor:      "Juan Pérez",
    Total:         150.50,
    Conceptos: []ticket.Concepto{
        {
            Descripcion: "Producto A",
            Cantidad:    2.0,
            PrecioVenta: 75.25,
            Total:       150.50,
        },
    },
}

// Convertir a JSON
data, err := ticket.ToBytes()
```

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests con cobertura
go test -cover ./...

# Tests específicos
go test ./pkg/escpos -v
go test ./internal/config -v
```

## 📖 Documentación de la API

Toda la API está documentada en español siguiendo las convenciones de godoc:

```bash
# Generar documentación local
go doc -all ./pkg/escpos
go doc -all ./internal/platform/windows
```

### Ejemplos de uso

- [Impresión básica](pkg/escpos/printer_test.go)
- [Configuración avanzada](internal/config/flags_test.go)
- [Manejo de tickets](internal/ticket/parse_test.go)

## 🐛 Resolución de problemas

### Errores comunes

1. **"Error al abrir la impresora"**: Verificar que el nombre coincida exactamente
2. **"Caracteres extraños en la impresión"**: Verificar codificación CP858
3. **"Acceso denegado"**: Ejecutar como administrador

### Debugging

```bash
# Habilitar logs detallados
go run cmd/print/main.go -printer "EC-PM-80250" -debug

# Verificar estado de impresoras Windows
wmic printer list brief
```

## 🤝 Contribución

Las contribuciones son bienvenidas. Por favor:

1. Fork del proyecto
2. Crear rama para tu feature (`git checkout -b feature/nueva-caracteristica`)
3. Commit de cambios (`git commit -am 'Añadir nueva característica'`)
4. Push a la rama (`git push origin feature/nueva-caracteristica`)
5. Crear Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 🔗 Enlaces útiles

- [Documentación ESC/POS](https://reference.epson-biz.com/modules/ref_escpos/)
- [Drivers EC-PM-80250](https://eclinepos.com/Producto.php?categoria=Impresoras&&buscar=EC-PM-80250)
- [API Windows Spooler](https://docs.microsoft.com/en-us/windows/win32/printdocs/printing-and-print-spooler-functions)

---

## 📌 Instalación detallada de drivers EC-PM-80250

### Requisitos previos

- Conexión USB (o adaptador serial virtual)
- Windows 10 o Windows 11
- Permisos de administrador en el equipo

### Pasos de instalación y configuración

#### 1. Descarga e instalación de drivers

1. Descarga el paquete de drivers desde: [descarga](https://eclinepos.com/Descargas/ControladoresZip/Impresoras/EC-PM-80250/Driver-2022.zip)
2. Ejecuta el instalador como administrador.
3. Selecciona **Install USB Virtual Serial Port Driver**.
4. Verifica en el **Administrador de dispositivos** que aparezca un puerto COM virtual.
5. Si la impresora no se detecta:
    - Vuelve a ejecutar el instalador.
    - Selecciona **Install Printer Driver (N)** y elige **POS Printer 300DPI Series**.
6. Al finalizar, deberían estar instaladas ambas versiones del driver:
    - POS Printer 203DPI Series
    - POS Printer 300DPI Series

#### 2. Configuración del puerto y nombre

1. Abre **Panel de control > Hardware y sonido > Dispositivos e impresoras**.
2. Localiza **POS Printer 203DPI Series**, haz clic derecho y selecciona **Propiedades de impresora**.
3. En la pestaña **Puertos**:
    - Por defecto estará en un puerto LPT.
    - Selecciona el puerto virtual **USB001 – UnknownPrinter**.
4. Ve a la pestaña **General** y cambia el nombre del dispositivo a **EC-PM-80250**.
5. Haz clic en **Aplicar** y **Aceptar**.

#### 3. Prueba de impresión

- Dentro de **Propiedades de impresora**, haz clic en **Imprimir página de prueba**.
- Alternativamente, abre una ventana de comando y ejecuta:
  ```bat
  echo Prueba EC-PM-80250 > \\?\USB#VID_xxxx&PID_yyyy#…\{GUID}\Printer
  ```

#### Observaciones

- En Windows 11 no suele ser necesario deshabilitar la validación de firma de drivers.
- Si trabajas en Windows 10 y necesitas instalar un driver sin certificado, consulta estos tutoriales:
    - [Video 1](https://www.youtube.com/watch?v=dEx-A-1ti_8&&ab_channel=SolucionesPOS)
    - [Video 2](https://www.youtube.com/watch?v=DtAIu2Is1nE&&t=320s&&ab_channel=INTSTORE)