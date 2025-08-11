# Guía de Desarrollo: CI/CD y Flujo de Trabajo

Este documento explica cómo funciona la infraestructura de CI/CD del proyecto, cómo gestionar versiones y cómo trabajar
con este repositorio sin depender de conocimiento tribal.

## 📋 Recomendaciones Adicionales

Además del archivo CONTRIBUTING.md, te sugiero:

1. **Añadir un Badge de CI en README.md**

# Configuración Rápida

## Requisitos

- Go 1.24 o superior
- Git configurado con firma de commits

## Primeros Pasos

1. Clona el repositorio
2. Ejecuta `go mod download`
3. Instala golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
