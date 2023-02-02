package database

import (
	"gorm.io/gorm"
	"time"
)

const ImagesTableName = "images"

type Image struct {
	Id          int64      `db:"id,omitempty"`
	UserId      int64      `db:"user_id"`
	Name        string     `db:"name"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type ImageRepository interface {
	Save(img Image) (Image, error)
	Update(img Image) (Image, error)
	Find(id int64) (Image, error)
	Delete(id int64) error
}

type imageRepository struct {
	sess *gorm.DB
}

func NewImageRepository(dbSession *gorm.DB) ImageRepository {
	return imageRepository{
		sess: dbSession,
	}
}

func (r imageRepository) Save(i Image) (Image, error) {
	i.CreatedDate = time.Now()
	i.UpdatedDate = i.CreatedDate
	err := r.sess.Table(ImagesTableName).Create(&i).Error
	if err != nil {
		return Image{}, err
	}
	return i, nil

}

func (r imageRepository) Update(i Image) (Image, error) {
	i.UpdatedDate = time.Now()
	err := r.sess.Save(&i).Error
	if err != nil {
		return Image{}, err
	}
	return i, nil
}

func (r imageRepository) Find(id int64) (Image, error) {
	var i Image
	err := r.sess.Table(UserTableName).First(&i, "id = ?", id).Error
	if err != nil {
		return Image{}, err
	}
	return i, nil
}

func (r imageRepository) Delete(id int64) error {
	err := r.sess.Table(UserTableName).Where("id = ?", id).Update("deleted_date", map[string]interface{}{"deleted_date": time.Now()}).Error
	if err != nil {
		return err
	}
	return nil
}
