# Guía de Contribución y Desarrollo

## 📚 Introducción

Este repositorio utiliza un sistema automatizado de CI/CD basado en GitHub Actions. El flujo de trabajo está diseñado
para ser simple pero efectivo, permitiendo releases automáticos basados en conventional commits sin intervención manual.
Esta guía explica cómo funciona todo y cómo trabajar con el sistema.

## 🚀 Flujo de Trabajo

### Desarrollo Individual

```bash
# 1. Clonar el repositorio
git clone https://github.com/AdConDev/pos-daemon.git
cd pos-daemon

# 2. Instalar dependencias
go mod download

# 3. Hacer cambios y commits siguiendo conventional commits
git add .
git commit -m "feat: add new printer driver for Brand X"
# o para fixes
git commit -m "fix: resolve connection timeout issue"

# 4. Push directo a main (solo si trabajas solo)
git push origin master
```

Al hacer push a `master`, el sistema:

1. Ejecutará pruebas y linting
2. Detectará el tipo de cambio (feat, fix, etc.)
3. Actualizará automáticamente la versión y el changelog
4. Creará una nueva release en GitHub

### Desarrollo en Equipo

```bash
# 1. Crear una rama para tu característica/fix
git checkout -b feat/new-feature

# 2. Hacer cambios y commits siguiendo conventional commits
git add .
git commit -m "feat: add new feature"

# 3. Push a tu rama
git push origin feat/new-feature

# 4. Crear Pull Request a través de GitHub UI
# El título del PR debe seguir conventional commits
```

## 📝 Conventional Commits

El proyecto utiliza [Conventional Commits](https://www.conventionalcommits.org/) para automatizar la generación de
versiones y changelog.

### Tipos de Commit Principales

| Tipo       | Descripción                       | ¿Genera Release? |
|------------|-----------------------------------|------------------|
| `feat`     | Nuevas características            | Minor (0.X.0)    |
| `fix`      | Correcciones de bugs              | Patch (0.0.X)    |
| `feat!`    | Cambios que rompen compatibilidad | Major (X.0.0)    |
| `docs`     | Solo documentación                | No               |
| `refactor` | Refactorización de código         | No               |
| `test`     | Añadir/modificar tests            | No               |
| `chore`    | Tareas de mantenimiento           | No               |
| `deps`     | Actualizaciones de dependencias   | Patch (0.0.X)    |

### Ejemplos

```
feat: add support for Epson TM-T88VI printer
fix: prevent connection timeout on slow networks
feat!: change printer configuration API format
docs: update installation instructions
deps: update golang.org/x/text to v0.14.0
```

## 🔍 Linters y Calidad de Código

El proyecto usa [golangci-lint](https://golangci-lint.run/) con una configuración simplificada para mantener la calidad
del código.

### Linters Habilitados

- **errcheck**: Detecta errores no manejados
- **govet**: Encuentra bugs potenciales
- **staticcheck**: Analizador estático general
- **ineffassign**: Variables asignadas pero no usadas
- **gosec**: Detección de problemas de seguridad
- **unused**: Detecta código no utilizado

### Cómo Ejecutar el Linter Localmente

```bash
# Instalar golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Ejecutar linting
golangci-lint run
```

### Notas Importantes sobre Linting

- La configuración está en `.golangci.yml` y es deliberadamente minimalista
- Solo fallarán issues importantes, no cuestiones de estilo subjetivas
- Si necesitas ignorar una regla específica en una línea, usa: `//nolint:lintername`

## 🏷️ Versionado y Releases

El proyecto utiliza [SemVer](https://semver.org/) con releases automáticos:

- **Major (1.0.0)**: Cambios incompatibles - commits con `!` o `BREAKING CHANGE`
- **Minor (0.1.0)**: Nuevas características - commits con `feat:`
- **Patch (0.0.1)**: Correcciones y mejoras - commits con `fix:` o `deps:`

### Generación Automática de Changelog

El changelog se genera automáticamente basado en los mensajes de commit. Solo aparecerán en el changelog los tipos:

- ✨ Features
- 🐛 Bug Fixes
- ⚡ Performance
- 📦 Dependencies
- ⏪ Reverts
- ✅ Tests
- 🤖 Continuous Integration
- 🏗️ Build System

Los otros tipos (docs, refactor, etc.) se ocultan para mantener el changelog enfocado en cambios relevantes para
usuarios.

## ⚠️ Qué Evitar

1. **No hacer bypass de las protecciones de rama**: Las reglas están ahí por una razón
2. **No incluir contraseñas o tokens en el código**: Usa variables de entorno
3. **No hacer push directo a `main` si hay más colaboradores**: Siempre usa PR
4. **No ignorar los errores del linter**: Arregla los problemas reales
5. **No crear releases manualmente**: Deja que el sistema las genere automáticamente

## ✅ Mejores Prácticas

1. **Usar branches por característica**: `feat/nombre`, `fix/problema`
2. **Commits atómicos**: Un commit por cambio lógico
3. **PR pequeños**: Más fáciles de revisar, menos propensos a errores
4. **Tests para todo**: Mantén el coverage alto
5. **Documentar APIs**: Comenta funciones exportadas siguiendo las convenciones de Go

## 🔄 Dependabot

El proyecto tiene Dependabot configurado para:

- Actualizar dependencias Go semanalmente (lunes)
- Actualizar GitHub Actions mensualmente
- Auto-merge de actualizaciones patch seguras
- Agrupar actualizaciones de golang.org/x/* para minimizar PRs

Si una actualización falla los tests, **no la mergees manualmente** sin resolver los problemas.

## 🧰 Estructura de Workflows

| Workflow                   | Propósito                                      |
|----------------------------|------------------------------------------------|
| `ci.yml`                   | Ejecuta tests y linting en push/PRs            |
| `release.yml`              | Genera releases automáticas basadas en commits |
| `dependabot-automerge.yml` | Auto-merge para actualizaciones seguras        |
| `pr-validation.yml`        | Valida y etiqueta PRs automáticamente          |

## 📄 Branch Protection

La rama `main` está protegida con:

- Revisión obligatoria de PRs
- Tests y linting pasando
- Firma de commits requerida
- No push directo (excepto bots)

## 🤝 Notas para Equipos

Esta configuración es independiente del código Go del proyecto, pudiéndose aplicar a cualquier proyecto Go con mínimos
ajustes. Las ventajas de esta aproximación son:

1. **Consistencia**: Mismo flujo de trabajo en todos los proyectos
2. **Automatización**: Menos trabajo manual, menos errores
3. **Trazabilidad**: Historial claro de cambios
4. **Bajo mantenimiento**: Una vez configurado, funciona sin intervención

## 🔍 Resolución de Problemas Comunes

### El Release No Se Genera

- Verifica que tus commits sigan el formato correcto
- Asegúrate de que hay al menos un commit `feat:` o `fix:` desde el último release
- Revisa los logs del workflow `release.yml`

### Linting Falla

- Ejecuta `golangci-lint run` localmente para reproducir
- Revisa los errores específicos en el log de CI
- Recuerda que errores de linting son problemas reales, no solo estéticos

### Tests Fallan en CI pero Pasan Localmente

- Revisa si hay dependencias en entornos o configuración local
- Asegúrate de que no hay race conditions (`go test -race`)
- Verifica la versión de Go (CI usa Go 1.24)