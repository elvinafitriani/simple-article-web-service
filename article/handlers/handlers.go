package handlers

import (
	"article/article"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHandlers(usecs article.UsecaseArticle, r *gin.RouterGroup) {
	eng := &usecase{
		use: usecs,
	}

	r.POST("", eng.Create)
	r.GET("", eng.GetAll)
	r.GET("/search/keyword/:keyword", eng.GetByKeyword)
	r.GET("/search/author/:author", eng.GetByAuthor)
}

type usecase struct {
	use article.UsecaseArticle
}

func (us usecase) Create(ctx *gin.Context) {
	err := us.use.Create(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "Data created successfully."})
}

func (us usecase) GetAll(ctx *gin.Context) {
	result, err := us.use.GetAll(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Articles": result, "Response": "Data retrieved successfully."})
}

func (us usecase) GetByKeyword(ctx *gin.Context) {
	result, err := us.use.GetByKeyword(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (us usecase) GetByAuthor(ctx *gin.Context) {
	result, err := us.use.GetByAuthor(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return	
	}

	ctx.JSON(http.StatusOK, result)
}
