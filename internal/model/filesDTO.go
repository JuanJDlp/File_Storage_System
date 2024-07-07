package model

import (
	"encoding/json"
	"time"
)

type FilesDTO struct {
	FileName       string    `db:"filename" json:"filename"`
	Size           int64     `db:"size" json:"size"`
	Date_Of_Upload time.Time `db:"date_of_upload" json:"-"`
}

func (f FilesDTO) MarshalJSON() ([]byte, error) {
	type Alias FilesDTO
	return json.Marshal(&struct {
		Alias
		TimeFormatted string `json:"date_of_upload"`
	}{
		Alias:         (Alias)(f),
		TimeFormatted: f.Date_Of_Upload.Format("01-02-2006"),
	})
}
