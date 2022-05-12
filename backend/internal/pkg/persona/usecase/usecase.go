package usecase

import (
	"RSOI/internal/models"
	"RSOI/internal/pkg/persona"
)

type PUsecase struct {
	repo persona.IRepository
}

func NewPUsecase(repo persona.IRepository) *PUsecase {
	return &PUsecase{repo: repo}
}

func (us *PUsecase) Create(persona *models.PersonaRequest) (uint, int) {
	return us.repo.Insert(persona)
}

func (us *PUsecase) Read(id uint) (*models.PersonaResponse, int) {
	return us.repo.Select(id)
}

func (us *PUsecase) ReadAll() ([]*models.PersonaResponse, int) {
	return us.repo.SelectAll()
}

func (us *PUsecase) Update(id uint, persona *models.PersonaRequest) int {
	persona.ID = id
	return us.repo.Update(persona)
}

func (us *PUsecase) Delete(id uint) int {
	return us.repo.Delete(id)
}
