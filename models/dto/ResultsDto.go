package dto

import "github.com/google/uuid"

type ResultDTO struct {
	UUID           uuid.UUID `json:"uuid"`
	Results        string    `json:"results"`
	PartnerResults string    `json:"partnerResults"`
	Status         string    `json:"status"`
	Reason         string    `json:"reason"`
}
