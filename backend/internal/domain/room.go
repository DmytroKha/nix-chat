package domain

type Room interface {
	GetId() int64
	GetName() string
	GetPrivate() bool
}
