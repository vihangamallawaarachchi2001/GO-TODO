package models

import "time"

type TODO struct {
	ID int
	Title string
	Description string
	CreatedAt time.Time
	Completed bool
	UpdatedAt time.Time
};