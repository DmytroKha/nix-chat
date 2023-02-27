package database

import (
	"gorm.io/gorm"
)

const UserTableName = "users"

type User struct {
	Id       int64 `gorm:"primary_key;auto_increment;not_null"`
	Name     string
	Password string
	Image    Image
}

type Users struct {
	Items []User
	Total uint64
	Pages uint64
}

type userRepository struct {
	sess *gorm.DB
}

//go:generate mockery --dir . --name UserRepository --output ./mocks
type UserRepository interface {
	Save(user User) (User, error)
	Update(user User) (User, error)
	Delete(id int64) error
	Find(id int64) (User, error)
	FindByName(name string) (User, error)
	FindAll() ([]User, error)
}

func NewUserRepository(dbSession *gorm.DB) UserRepository {
	return &userRepository{
		sess: dbSession,
	}
}

//func (user *User) GetId() int64 {
//	return user.Id
//}

//func (user *User) GetName() string {
//	return user.Name
//}

func (r userRepository) Save(u User) (User, error) {
	err := r.sess.Table(UserTableName).Create(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r userRepository) Update(u User) (User, error) {
	err := r.sess.Save(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r userRepository) Delete(id int64) error {
	err := r.sess.Table(UserTableName).Where("id = ?", id).Delete(User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Find(id int64) (User, error) {
	var u User
	err := r.sess.Table(UserTableName).First(&u, "id = ?", id).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *userRepository) FindByName(name string) (User, error) {
	var u User
	err := r.sess.Table(UserTableName).First(&u, "name = ?", name).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r userRepository) FindAll() ([]User, error) {
	var users []User
	err := r.sess.Table(UserTableName).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
