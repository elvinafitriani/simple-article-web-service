package article

import (
	"article/entity"

	"github.com/gin-gonic/gin"
)

type UsecaseArticle interface {
	Create(*gin.Context) error
	GetAll(ctx *gin.Context) ([]entity.Article, error)
	GetByKeyword(ctx *gin.Context) ([]entity.Article, error)
	GetByAuthor(ctx *gin.Context) ([]entity.Article, error)
}

type RepositoryArticle interface {
	Create(entity.Article, *gin.Context) (*entity.Article, error)
	GetAll() ([]entity.Article, error)
	GetByKeyword(string) ([]entity.Article, error)
	GetByAuthor(string) ([]entity.Article, error)
}
