package domain

type User interface {
	//
	GetUid() string
	GetName() string
	GetPhoto() string
	GetId() int64
}

//type UserRepository interface {
//	AddUser(user User)
//	RemoveUser(user User)
//	FindUserById(ID string) User
//	GetAllUsers() []User
//}
