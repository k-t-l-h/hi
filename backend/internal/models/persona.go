package models

//todo: убрать теги
type PersonaRequest struct {
	ID      uint   `json:"id,omitempty" validate:""`
	Name    string `json:"name"`
	Age     uint   `json:"age,omitempty"`
	Address string `json:"address,omitempty"`
	Work    string `json:"work,omitempty"`
}

type PersonaResponse struct {
	ID      uint   `json:"id,required" validate:""`
	Name    string `json:"name,required"`
	Age     uint   `json:"age,omitempty"`
	Address string `json:"address,omitempty"`
	Work    string `json:"work,omitempty"`
}
