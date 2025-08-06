package profile

// Profile define todas las características físicas y capacidades de una impresora
type Profile struct {
	// Información básica
	Model       string
	Vendor      string
	Description string

	// Características físicas
	PaperWidth  float64 // en mm (58mm, 80mm, etc.)
	PaperHeight float64 // en mm (0 para rollo continuo)
	DPI         int     // Dots Per Inch (ej. 203, 300)
	DotsPerLine int     // Puntos por línea (ej. 384, 576)

	// Capacidades
	SupportsGraphics bool // Soporta gráficos (imágenes)
	SupportsBarcode  bool // Soporta códigos de barra nativos
	SupportsQR       bool // Soporta códigos QR nativos
	SupportsCutter   bool // Tiene cortador automático
	SupportsDrawer   bool // Soporta cajón de dinero
	SupportsColor    bool // Soporta impresión a color

	// Juegos de caracteres
	CharacterSets  []int // Códigos de página soportados
	DefaultCharSet int   // Código de página por defecto
	ActiveCharSet  int   // Código de página activo (para cambiar dinámicamente)

	// Configuración avanzada (opcional)
	FeedLinesAfterCut int // Líneas de avance después de cortar
	ImageThreshold    int // Umbral para conversión B/N (0-255)

	// Fuentes
	Fonts map[string]int // Lista de fuentes soportadas, nombre -> ancho (en puntos)

	// Extensible para características específicas
	// Usar un mapa genérico permite agregar características sin cambiar la estructura
	ExtendedFeatures map[string]interface{}
}

// ModelInfo devuelve una representación de string del modelo
func (p *Profile) ModelInfo() string {
	return p.Vendor + " " + p.Model
}

// GetCharWidth calcula el ancho físico de un caracter en milímetros
func (p *Profile) GetCharWidth(font string) int {
	return p.DotsPerLine / p.Fonts[font]
}

func CreateProfPT_210() *Profile {
	p := CreateProfile58mm()
	p.Model = "58mm GOOJPRT PT-210"
	p.CharacterSets = []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16, 17, 18, 19, 20, 21,
	}
	p.DefaultCharSet = 19 // CP858 para español
	p.SupportsQR = true  // Soporta QR nativo
	return p
}

func CreateProfGP_58N() *Profile {
	p := CreateProfile58mm()
	p.Model = "58mm GP-58N"
	p.CharacterSets = []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16, 17, 18, 19, 20, 21,
	}
	p.DefaultCharSet = 19 // CP858 para español
	return p
}

// CreateProfile58mm crea un perfil para impresora térmica de 58mm común
func CreateProfile58mm() *Profile {
	return &Profile{
		Model:       "Generic 58mm",
		Vendor:      "Generic",
		Description: "Impresora térmica genérica de 58mm",

		PaperWidth:  58,
		DPI:         203,
		DotsPerLine: 384, // Típico para 58mm a 203 DPI

		SupportsGraphics: true,
		SupportsBarcode:  true,
		SupportsQR:       false, // Muchas impresoras baratas no soportan QR nativo
		SupportsCutter:   false,
		SupportsDrawer:   false,
		SupportsColor:    false,

		CharacterSets: []int{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16, 17, 18, 19, 20, 21,
		}, // Más juegos de caracteres
		DefaultCharSet: 19, // CP858
		ActiveCharSet:  19, // CP858

		ExtendedFeatures: make(map[string]interface{}),
	}
}

func CreateProfEC_PM_80250() *Profile {
	p := CreateProfile80mm()
	p.Model = "80mm EC-PM-80250"
	p.CharacterSets = []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16, 17, 18, 19, 20, 21,
	}
	p.DefaultCharSet = 19 // CP858 para español
	return p
}

// CreateProfile80mm crea un perfil para impresora térmica de 80mm común
func CreateProfile80mm() *Profile {
	return &Profile{
		Model:       "Generic 80mm",
		Vendor:      "Generic",
		Description: "Impresora térmica genérica de 80mm",

		PaperWidth:  80,
		DPI:         203,
		DotsPerLine: 576, // Típico para 80mm a 203 DPI

		SupportsGraphics: true,
		SupportsBarcode:  true,
		SupportsQR:       true, // Las 80mm suelen tener más funciones
		SupportsCutter:   true,
		SupportsDrawer:   true,
		SupportsColor:    false,

		CharacterSets: []int{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16, 17, 18, 19, 20, 21,
		}, // Más juegos de caracteres
		DefaultCharSet: 19, // CP858
		ActiveCharSet:  19, // CP858

		FeedLinesAfterCut: 5,
		ImageThreshold:    128,

		ExtendedFeatures: make(map[string]interface{}),
	}
}
