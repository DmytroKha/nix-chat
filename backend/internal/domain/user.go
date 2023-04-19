package domain

type User interface {
	GetName() string
	GetPhoto() string
	GetId() int64
}
