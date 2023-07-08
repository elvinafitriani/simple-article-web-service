package usecase

import (
	"article/article"
	"article/entity"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func NewUsecase(repoo article.RepositoryArticle, rediss *redis.Client) repository {
	return repository{
		repo:  repoo,
		redis: rediss,
	}
}

type repository struct {
	repo  article.RepositoryArticle
	redis *redis.Client
}

func (rep repository) Create(ctx *gin.Context) error {
	var article entity.Article

	if err := ctx.ShouldBindJSON(&article); err != nil {
		return err
	}

	result, err := rep.repo.Create(article, ctx)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("article:%d", result.ID)
	jsonVal, err := json.Marshal(result)
	if err != nil {
		return err
	}
	err = rep.redis.Set(ctx, cacheKey, jsonVal, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rep repository) GetAll(ctx *gin.Context) ([]entity.Article, error) {
	var result []entity.Article
	var err error

	keys, err := rep.redis.Keys(ctx, "article:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		articleData, err := rep.redis.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var article entity.Article
		err = json.Unmarshal([]byte(articleData), &article)
		if err != nil {
			return nil, err
		}

		result = append(result, article)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	if result == nil {
		result, err = rep.repo.GetAll()
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (rep repository) GetByKeyword(ctx *gin.Context) ([]entity.Article, error) {
	var result []entity.Article
	var err error
	var Keyword struct {
		Key string `uri:"keyword"`
	}
	if err := ctx.ShouldBindUri(&Keyword); err != nil {
		return nil, err
	}

	keys, err := rep.redis.Keys(ctx, "article:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		articleData, err := rep.redis.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var article entity.Article
		err = json.Unmarshal([]byte(articleData), &article)
		if err != nil {
			return nil, err
		}

		if strings.Contains(article.Title, Keyword.Key) || strings.Contains(article.Body, Keyword.Key) {
			result = append(result, article)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	if result == nil {
		result, err = rep.repo.GetByKeyword(Keyword.Key)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (rep repository) GetByAuthor(ctx *gin.Context) ([]entity.Article, error) {
	var result []entity.Article
	var err error
	var Author struct {
		Author string `uri:"author"`
	}
	if err := ctx.ShouldBindUri(&Author); err != nil {
		return nil, err
	}

	keys, err := rep.redis.Keys(ctx, "article:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		articleData, err := rep.redis.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var article entity.Article
		err = json.Unmarshal([]byte(articleData), &article)
		if err != nil {
			return nil, err
		}

		if strings.Contains(article.Author, Author.Author) {
			result = append(result, article)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	if result == nil {
		result, err = rep.repo.GetByAuthor(Author.Author)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
