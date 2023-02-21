package database

import (
	"gorm.io/gorm"
	"time"
)

const ImagesTableName = "images"

type Image struct {
	Id          int64     `db:"id,omitempty"`
	UserId      int64     `db:"user_id"`
	Name        string    `db:"name"`
	CreatedDate time.Time `db:"created_date,omitempty"`
	UpdatedDate time.Time `db:"updated_date,omitempty"`
	Deleted     int64     `db:"deleted"`
}

type ImageRepository interface {
	Save(img Image) (Image, error)
	Update(img Image) (Image, error)
	Find(id int64) (Image, error)
	FindAll(userId int64) ([]Image, error)
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
	err := r.sess.Table(ImagesTableName).First(&i, "id = ?", id).Error
	if err != nil {
		return Image{}, err
	}
	return i, nil
}

func (r imageRepository) FindAll(userId int64) ([]Image, error) {
	var images []Image

	//err := r.sess.Table(ImagesTableName).Where("user_id = ?", userId).Find(&images).Error
	err := r.sess.Table(ImagesTableName).Where(&Image{UserId: userId, Deleted: 0}).Find(&images).Error
	if err != nil {
		return []Image{}, err
	}
	return images, nil
}

func (r imageRepository) Delete(id int64) error {
	var i Image
	err := r.sess.Table(ImagesTableName).First(&i, "id = ?", id).Error
	if err != nil {
		return err
	}

	i.UpdatedDate = time.Now()
	i.Deleted = 1
	err = r.sess.Save(&i).Error
	if err != nil {
		return err
	}

	return nil
}
