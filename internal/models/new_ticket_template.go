package models

type NewTicketTemplate struct {
	Data NewTicketTemplateData `json:"data"` // Datos de la plantilla
}

type NewTicketTemplateData struct {
	TicketTemplateData
	VerSeries BoolFlex `json:"ver_series"`
}
