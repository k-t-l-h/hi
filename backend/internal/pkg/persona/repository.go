package persona

import "RSOI/internal/models"

type IRepository interface {
	Insert(persona *models.PersonaRequest) (uint, int)
	Select(id uint) (*models.PersonaResponse, int)
	SelectAll() ([]*models.PersonaResponse, int)
	Update(persona *models.PersonaRequest) int
	Delete(id uint) int
}
