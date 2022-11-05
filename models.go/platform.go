package models

type Platform struct {
	Model
	Name               string `json:"name" db:"name"`
	Category           string `json:"category" db:"category"`
	Website            string `json:"website" db:"website"`
	Country            string `json:"country" db:"country"`
	ContactInformation string `json:"contactInformation" db:"contact_information"`
	BusinessModelNotes string `json:"businessModelNotes" db:"business_model_notes"`
	Source             string `json:"source" db:"source"`
	GeneralNotes       string `json:"generalNotes" db:"general_notes"`
	Comments           string `json:"comments" db:"comments"`
}
