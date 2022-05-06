package models

import (
	"time"
)

type Model int64

const (
	TODO Model = 0
	RES  Model = 1
	NOTE Model = 2
)

type Todo struct {
	Info      string
	Timestamp time.Time
	Deadline  time.Time
}

type Resource struct {
	Info      string
	Timestamp time.Time
}

type Note struct {
	Info      string
	Timestamp time.Time
}
