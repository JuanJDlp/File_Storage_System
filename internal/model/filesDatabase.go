package model

import "time"

type FileDatabase struct {
	Hash           string    `db:"hash"`
	FileName       string    `db:"filename"`
	Size           int64     `db:"size"`
	Date_of_upload time.Time `db:"date_of_upload"`
	Owner          string    `db:"owner"`
}
