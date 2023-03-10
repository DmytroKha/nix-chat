package database

import (
	"github.com/DmytroKha/nix-chat/internal/domain"
	"gorm.io/gorm"
)

const UserTableName = "users"

type User struct {
	Id       int64 `gorm:"primary_key;auto_increment;not_null"`
	Uid      string
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
	Delete(uid string) error
	Find(uid string) (User, error)
	FindByName(name string) (User, error)
	FindAll() ([]domain.User, error)
}

func NewUserRepository(dbSession *gorm.DB) UserRepository {
	return &userRepository{
		sess: dbSession,
	}
}

func (user *User) GetId() string {
	return user.Uid
}

func (user *User) GetName() string {
	return user.Name
}

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

func (r userRepository) Delete(uid string) error {
	err := r.sess.Table(UserTableName).Where("uid = ?", uid).Delete(User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Find(uid string) (User, error) {
	var u User
	err := r.sess.Table(UserTableName).First(&u, "uid = ?", uid).Error
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

func (r userRepository) FindAll() ([]domain.User, error) {
	var usrs []User
	err := r.sess.Table(UserTableName).Find(&usrs).Error
	if err != nil {
		return nil, err
	}
	var users []domain.User
	for _, usr := range usrs {
		users = append(users, &usr)
	}

	return users, nil
}
