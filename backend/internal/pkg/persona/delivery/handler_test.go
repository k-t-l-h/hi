package delivery

import (
	"RSOI/internal/models"
	"RSOI/internal/pkg/persona"
	"RSOI/internal/pkg/persona/mock"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPHandler_Create(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockIUsecase(ctl)

	type fields struct {
		personaUsecase persona.IUsecase
	}

	type args struct {
		r      *http.Request
		result http.Response
		status int
		statusReturn int
		expected models.PersonaRequest
		times    int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "simple create",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "/persons",
					strings.NewReader(fmt.Sprintf(`{"name": "%s" }`, "name"))),
				expected: models.PersonaRequest{Name: "name"},
				status:   http.StatusCreated,
				statusReturn: models.OKEY,
				times:    1,
			}},
		{
			name:   "json err",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "/persons",
					strings.NewReader(fmt.Sprint())),
				expected: models.PersonaRequest{Name: "name"},
				status:   http.StatusBadRequest,
				statusReturn: models.BADREQUEST,
				times:    0,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := &PHandler{
				personaUsecase: tt.fields.personaUsecase,
			}
			w := httptest.NewRecorder()

			mockUsecase.EXPECT().Create(&tt.args.expected).Return(uint(0), tt.args.statusReturn).Times(tt.args.times)

			h.Create(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
			}
		})
	}
}

func TestPHandler_Read(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockIUsecase(ctl)

	type fields struct {
		personaUsecase persona.IUsecase
	}
	type args struct {
		r        *http.Request
		result   http.Response
		status   int
		expected models.PersonaRequest
		times    int
		state    int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "simple read",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r:        httptest.NewRequest("GET", "/persons/0", nil),
				expected: models.PersonaRequest{ID: 0},
				status:   http.StatusOK,
				times:    1,
				state:    models.OKEY,
			}},
		{
			name:   "read not found",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r:        httptest.NewRequest("GET", "/persons/1", nil),
				expected: models.PersonaRequest{ID: 1},
				status:   http.StatusNotFound,
				times:    1,
				state:    models.NOTFOUND,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PHandler{
				personaUsecase: tt.fields.personaUsecase,
			}

			tt.args.r = mux.SetURLVars(tt.args.r, map[string]string{
				"personID": fmt.Sprint(tt.args.expected.ID),
			})

			w := httptest.NewRecorder()

			gomock.InOrder(
				mockUsecase.EXPECT().Read(tt.args.expected.ID).Return(&models.PersonaResponse{ID: tt.args.expected.ID}, tt.args.state))

			h.Read(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
			}
		})
	}
}

func TestPHandler_ReadAll(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockIUsecase(ctl)

	type fields struct {
		personaUsecase persona.IUsecase
	}
	type args struct {
		r        *http.Request
		result   http.Response
		status   int
		expected []*models.PersonaResponse
		state    int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "simple read all",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r:        httptest.NewRequest("GET", "/persons", nil),
				expected: []*models.PersonaResponse{{}, {}},
				status:   http.StatusOK,
				state:    models.OKEY,
			}},
		{
			name:   "read all no users",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r:        httptest.NewRequest("GET", "/persons", nil),
				expected: []*models.PersonaResponse{},
				status:   http.StatusNotFound,
				state:    models.NOTFOUND,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PHandler{
				personaUsecase: tt.fields.personaUsecase,
			}

			w := httptest.NewRecorder()

			gomock.InOrder(
				mockUsecase.EXPECT().ReadAll().Return(tt.args.expected, tt.args.state))

			h.ReadAll(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
			}
			log.Print(w.Result())
		})
	}
}

func TestPHandler_Update(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockIUsecase(ctl)

	type fields struct {
		personaUsecase persona.IUsecase
	}
	type args struct {
		r        *http.Request
		result   http.Response
		status   int
		expected models.PersonaRequest
		id       uint
		state    int
		times    int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "simple update",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("PATCH", "/persons/0",
					strings.NewReader(fmt.Sprintf(`{"name": "%s" }`, "name"))),
				id:       0,
				expected: models.PersonaRequest{Name: "name"},
				status:   http.StatusOK,
				state:    models.OKEY,
				times:    1,
			}},
		{
			name:   "update not found",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("PATCH", "/persons/5",
					strings.NewReader(fmt.Sprintf(`{"name": "%s" }`, "name"))),
				id:       5,
				expected: models.PersonaRequest{Name: "name"},
				status:   http.StatusNotFound,
				state:    models.NOTFOUND,
				times:    1,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PHandler{
				personaUsecase: tt.fields.personaUsecase,
			}

			w := httptest.NewRecorder()

			tt.args.r = mux.SetURLVars(tt.args.r, map[string]string{
				"personID": fmt.Sprint(tt.args.id),
			})

			gomock.InOrder(
				mockUsecase.EXPECT().Update(tt.args.id, &tt.args.expected).Return(tt.args.state).Times(tt.args.times))

			h.Update(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
			}

		})
	}
}

func TestPHandler_Delete(t *testing.T) {

	r1, _ := http.NewRequest("DELETE", "/persons/1", nil)
	r2, _ := http.NewRequest("DELETE", "/persons/100", nil)
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockIUsecase(ctl)

	type fields struct {
		personaUsecase persona.IUsecase
	}
	type args struct {
		r        *http.Request
		result   http.Response
		status   int
		state    int
		expected models.PersonaRequest
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "simple delete",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r:        r1,
				expected: models.PersonaRequest{ID: 1},
				status:   http.StatusOK,
				state:    models.OKEY,
			}},
		{
			name:   "delete not found",
			fields: fields{personaUsecase: mockUsecase},
			args: args{
				r:        r2,
				expected: models.PersonaRequest{ID: 100},
				status:   http.StatusNotFound,
				state:    models.NOTFOUND,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PHandler{
				personaUsecase: tt.fields.personaUsecase,
			}

			w := httptest.NewRecorder()

			tt.args.r = mux.SetURLVars(tt.args.r, map[string]string{
				"personID": fmt.Sprint(tt.args.expected.ID),
			})

			gomock.InOrder(
				mockUsecase.EXPECT().Delete(tt.args.expected.ID).Return(tt.args.state))

			h.Delete(w, tt.args.r)

			if tt.args.status != w.Code {
				log.Print(w.Code)
				t.Error(tt.name)
			}

		})
	}
}
