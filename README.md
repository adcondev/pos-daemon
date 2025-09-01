![CI Status](https://github.com/AdConDev/pos-daemon/actions/workflows/ci.yml/badge.svg)
![License](https://img.shields.io/github/license/AdConDev/pos-daemon)

# POS Daemon 🖨️

A high-performance, protocol-agnostic Point of Sale printing daemon designed for modern retail environments. Built with Go for reliability and efficiency.

## Potential Features

- **Multi-Protocol Support**: ESC/POS, StarPRN, and custom protocols
- **Queue Management**: Intelligent print job queuing and prioritization
- **Network Printing**: Support for network, USB, and serial connections
- **REST API**: Simple HTTP endpoints for integration
- **Hot Reload**: Configuration changes without service restart
- **Error Recovery**: Automatic retry logic and graceful degradation

# Instalación de drivers para impresora [EC-PM-80250](https://eclinepos.com/Producto.php?categoria=Impresoras&&buscar=EC-PM-80250) en Windows 10/11

Este documento describe de forma clara y estructurada los pasos necesarios para instalar y configurar los drivers de la impresora térmica **EC‑PM‑80250** en Windows 10 y Windows 11.

---

## 📋 Requisitos previos

- Conexión USB (o adaptador serial virtual)
- Windows 10 o Windows 11
- Permisos de administrador en el equipo

---

## 🛠 Pasos de instalación y configuración

### 1. Descarga e instalación de drivers

1. Descarga el paquete de drivers desde: [descarga](https://eclinepos.com/Descargas/ControladoresZip/Impresoras/EC-PM-80250/Driver-2022.zip)
2. Ejecuta el instalador como administrador.
3. Selecciona **Install USB Virtual Serial Port Driver**.
4. Verifica en el **Administrador de dispositivos** que aparezca un puerto COM virtual.
5. Si la impresora no se detecta:
    - Vuelve a ejecutar el instalador.
    - Selecciona **Install Printer Driver (N)** y elige **POS Printer 300DPI Series**.
6. Al finalizar, deberían estar instaladas ambas versiones del driver:
    - POS Printer 203DPI Series
    - POS Printer 300DPI Series

### 2. Configuración del puerto y nombre

1. Abre **Panel de control > Hardware y sonido > Dispositivos e impresoras**.
2. Localiza **POS Printer 203DPI Series**, haz clic derecho y selecciona **Propiedades de impresora**.
3. En la pestaña **Puertos**:
    - Por defecto estará en un puerto LPT.
    - Selecciona el puerto virtual **USB001 – UnknownPrinter**.
4. Ve a la pestaña **General** y cambia el nombre del dispositivo a **EC-PM-80250**.
5. Haz clic en **Aplicar** y **Aceptar**.

### 3. Prueba de impresión

- Dentro de **Propiedades de impresora**, haz clic en **Imprimir página de prueba**.
- Alternativamente, abre una ventana de comando y ejecuta:
  ```bat
  echo Prueba EC-PM-80250 > \\?\USB#VID_xxxx&PID_yyyy#…\{GUID}\Printer
  ```

---

## 📌 Observaciones

- En Windows 11 no suele ser necesario deshabilitar la validación de firma de drivers.
- Si trabajas en Windows 10 y necesitas instalar un driver sin certificado, consulta estos tutoriales:
    - [Video 1](https://www.youtube.com/watch?v=dEx-A-1ti_8&&ab_channel=SolucionesPOS)
    - [Video 2](https://www.youtube.com/watch?v=DtAIu2Is1nE&&t=320s&&ab_channel=INTSTORE)

## 🛠️ Contribuir

¿Interesado en contribuir? Por favor lee:

- [CONTRIBUTING.md](CONTRIBUTING.md) - Flujo de trabajo y convenciones
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) - Reglas de participación

El proyecto utiliza Conventional Commits y SemVer para versionado automático.

