package persona

import "RSOI/internal/models"

type IUsecase interface {
	Create(persona *models.PersonaRequest) (uint, int)
	Read(id uint) (*models.PersonaResponse, int)
	ReadAll() ([]*models.PersonaResponse, int)
	Update(id uint, persona *models.PersonaRequest) int
	Delete(id uint) int
}
