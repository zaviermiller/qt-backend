package models

import (

)

type Section struct {
	QTModel
	Type string `json:"type"` 
	Position int `json:"position"`
	Test Test
	TestID uint
}