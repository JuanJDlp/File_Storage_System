package model

import "time"

type FilesDTO struct {
	FileName       string    `db:"filename" json:"filename"`
	Size           int64     `db:"size" json:"size"`
	Date_of_upload time.Time `db:"date_of_upload" json:"date_of_upload"`
}
