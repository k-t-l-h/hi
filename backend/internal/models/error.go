package models

type Error struct {
	Message string `json:"message"`
}

type ErrorValidation struct {
	Message string `json:"message"`
	Errors  struct {
		additionalProp1 string `json:"additionalProp1"`
		additionalProp2 string `json:"additionalProp2"`
		additionalProp3 string `json:"additionalProp3"`
	} `json:"errors"`
}
