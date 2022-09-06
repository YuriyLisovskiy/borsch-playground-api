package models

type JobOutputRowDbModel struct {
	Model

	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Text  string `json:"text"`
	JobID string `json:"-"`
}

type JobDbModel struct {
	Model

	Code     string                `json:"code"`
	Outputs  []JobOutputRowDbModel `json:"-" gorm:"foreignKey:JobID"`
	ExitCode *int                  `json:"exit_code"`
}
