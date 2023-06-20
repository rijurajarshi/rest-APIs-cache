package models

type University struct {
	Ranking  int    `gorm:"primary_key" json:"ranking" `
	Title    string `json:"title" gorm:"type:varchar(100);column:title"`
	Location string `json:"location" gorm:"type:varchar(100);column:location"`
	// StudentFrequency uint32 `json:"number students"`
	// StaffRatio       string `json:"students staff ratio" gorm:"type:varchar(100);column:staff_ratio"`
	// Percent          string `json:"perc intl students" gorm:"type:varchar(100);column:percent"`
	// GenderRatio      string `json:"gender ratio" gorm:"type:varchar(100);column:gender_ratio"`
}
