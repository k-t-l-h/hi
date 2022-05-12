package delivery

import (
	"RSOI/internal/models"
	"RSOI/internal/pkg/persona"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
)

type PHandler struct {
	personaUsecase persona.IUsecase
}

func NewPHandler(personaUsecase persona.IUsecase) *PHandler {
	return &PHandler{personaUsecase: personaUsecase}
}

func (h *PHandler) Create(w http.ResponseWriter, r *http.Request) {

	person := &models.PersonaRequest{}
	err := easyjson.UnmarshalFromReader(r.Body, person)
	if err != nil {
		answer := models.ErrorValidation{
			Message: "Incorrect json",
		}
		jsn, _ := json.Marshal(answer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsn)
		return
	}

	id, code := h.personaUsecase.Create(person)

	switch code {
	case models.OKEY:
		w.Header().Set("Location",
			fmt.Sprintf("https://persona-service.herokuapp.com/persons/%d", id))
		w.WriteHeader(http.StatusCreated)
	case models.NOTFOUND:
		answer := models.Error{
			Message: "Person not found",
		}
		jsn, _ := json.Marshal(answer)
		w.Write(jsn)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *PHandler) Read(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	ids := v["personID"]

	id, err := strconv.Atoi(ids)
	if err != nil {
		answer := models.ErrorValidation{
			Message: "Incorrect id",
		}
		jsn, _ := json.Marshal(answer)


		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsn)
		return
	}

	p, code := h.personaUsecase.Read(uint(id))

	switch code {
	case models.OKEY:
		easyjson.MarshalToHTTPResponseWriter(p, w)
	case models.NOTFOUND:
		answer := models.Error{
			Message: "Person not found",
		}
		jsn, _ := json.Marshal(answer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsn)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (h *PHandler) ReadAll(w http.ResponseWriter, r *http.Request) {

	ps, code := h.personaUsecase.ReadAll()
	switch code {
	case models.OKEY:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ps)
		//w.WriteHeader(http.StatusOK)
	case models.NOTFOUND:
		w.Header().Set("Content-Type", "application/json")
		answer := models.Error{
			Message: "Person not found",
		}
		jsn, _ := json.Marshal(answer)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsn)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		answer := models.ErrorValidation{
			Message: "Incorrect data",
		}
		js, _ := json.Marshal(answer)
		w.Write(js)
	}
}

func (h *PHandler) Update(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	ids := v["personID"]
	id, _ := strconv.Atoi(ids)

	persona := &models.PersonaRequest{}
	easyjson.UnmarshalFromReader(r.Body, persona)

	code := h.personaUsecase.Update(uint(id), persona)

	switch code {
	case models.OKEY:
		w.WriteHeader(http.StatusOK)
	case models.NOTFOUND:
		answer := models.Error{
			Message: "Person not found",
		}
		jsn, _ := json.Marshal(answer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsn)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (h *PHandler) Delete(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	ids := v["personID"]
	id, err := strconv.Atoi(ids)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := h.personaUsecase.Delete(uint(id))

	switch code {
	case models.OKEY:
		w.WriteHeader(http.StatusOK)
	case models.NOTFOUND:
		answer := models.Error{
			Message: "Person not found",
		}

		jsn, _ := json.Marshal(answer)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsn)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}