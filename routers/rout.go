package routers

import (
	handArticle "article/article/handlers"
	repoArticle "article/article/repository"
	useArticle "article/article/usecase"

	"article/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Routes struct {
	Db    *gorm.DB
	R     *gin.Engine
	Redis *redis.Client
}

func (r Routes) Routers() {
	middleware.Add(r.R, middleware.CORSMiddleware())
	v1 := r.R.Group("article")
	repositoryArticle := repoArticle.NewRepository(r.Db)
	usecaseArticle := useArticle.NewUsecase(repositoryArticle, r.Redis)
	handArticle.NewHandlers(usecaseArticle, v1)

}
