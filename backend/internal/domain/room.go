package domain

type Room interface {
	GetUid() string
	GetName() string
	GetPrivate() bool
}

type RoomRepository interface {
	Save(room Room) (Room, error)
	FindByName(name string) (Room, error)
	FindAll() ([]Room, error)
}
