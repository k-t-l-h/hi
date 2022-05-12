package repository

import (
	"RSOI/internal/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type PRepository struct {
	pool pgxpool.Pool
}

func NewPRepository(pl pgxpool.Pool) *PRepository {
	return &PRepository{pool: pl}
}

func (pr *PRepository) Insert(persona *models.PersonaRequest) (uint, int) {

	res := pr.pool.QueryRow(context.Background(), CREATEPERSONA,
		persona.Name, persona.Age, persona.Address, persona.Work)
	err := res.Scan(&persona.ID)
	if err != nil {
		return 0, models.NOTFOUND
	} else {
		return persona.ID, models.OKEY
	}
}

func (pr *PRepository) Select(id uint) (*models.PersonaResponse, int) {

	persona := models.PersonaResponse{ID: id}
	res := pr.pool.QueryRow(context.Background(), READPERSONA, persona.ID)

	err := res.Scan(&persona.Name, &persona.Age, &persona.Address, &persona.Work)
	if err != nil {
		log.Print(err)
		return &persona, models.NOTFOUND
	} else {
		return &persona, models.OKEY
	}
}

func (pr *PRepository) SelectAll() ([]*models.PersonaResponse, int) {
	tag, err := pr.pool.Query(context.Background(), READALLPERSONA)
	var personas []*models.PersonaResponse

	for tag.Next() {
		persona := models.PersonaResponse{}
		err = tag.Scan(&persona.ID, &persona.Name, &persona.Age, &persona.Address, &persona.Work)
		if err != nil {
			log.Print(err)
			break
		}
		personas = append(personas, &persona)
	}

	if err != nil || len(personas) == 0 {
		return personas, models.NOTFOUND
	} else {
		return personas, models.OKEY
	}
}

func (pr *PRepository) Update(persona *models.PersonaRequest) int {
	tag, err := pr.pool.Exec(context.Background(), UPDATEPERSONA,
		persona.Name, persona.Age, persona.Address, persona.Work, persona.ID)
	if err != nil {
		return models.BADREQUEST
	}
	if tag.RowsAffected() == 0 {
		return models.NOTFOUND
	}

	return models.OKEY

}

func (pr *PRepository) Delete(id uint) int {
	tag, err := pr.pool.Exec(context.Background(), DELETEPERSONA, id)

	if err != nil {
		return models.BADREQUEST
	}
	if tag.RowsAffected() == 0 {
		return models.NOTFOUND
	}

	return models.OKEY
}
