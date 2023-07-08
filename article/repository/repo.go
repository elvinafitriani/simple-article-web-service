package repository

import (
	"article/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Database {
	return Database{
		Db: db,
	}
}

type Database struct {
	Db *gorm.DB
}

func (db Database) Create(article entity.Article, ctx *gin.Context) (*entity.Article, error) {
	if err := db.Db.Create(&article).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

func (db Database) GetAll() (article []entity.Article, err error) {
	err = db.Db.Order("created_at desc").Find(&article).Error

	if err != nil {
		return nil, err
	}

	return article, nil
}

func (db Database) GetByKeyword(key string) ([]entity.Article, error) {
	var article []entity.Article

	if err := db.Db.Order("created_at desc").First(&article, "body LIKE ? OR title LIKE ?", "%"+key+"%", "%"+key+"%").Error; err != nil {
		return nil, err
	}

	err := db.Db.Order("created_at desc").Find(&article, "body LIKE ? OR title LIKE ?", "%"+key+"%", "%"+key+"%").Error
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (db Database) GetByAuthor(author string) ([]entity.Article, error) {
	var article []entity.Article

	if err := db.Db.Order("created_at desc").First(&article, "author LIKE ?", "%"+author+"%").Error; err != nil {
		return nil, err
	}

	err := db.Db.Order("created_at desc").Find(&article, "author LIKE ?", "%"+author+"%").Error
	if err != nil {
		return nil, err
	}

	return article, nil
}
