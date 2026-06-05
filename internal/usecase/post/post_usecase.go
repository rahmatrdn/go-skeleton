package post_usecase

import (
	"context"

	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/post/entity"
)

type PostUsecase struct {
	postRepo mysql.IPostRepository
}

func NewPostUsecase(
	postRepo mysql.IPostRepository,
) *PostUsecase {
	return &PostUsecase{postRepo}
}

type IPostUsecase interface {
	GetAll(ctx context.Context) (res []*entity.PostResponse, err error)
	GetByID(ctx context.Context, postID int64) (*entity.PostResponse, error)
	Create(ctx context.Context, req entity.PostReq) (*entity.PostResponse, error)
	UpdateByID(ctx context.Context, req entity.PostReq) error
	DeleteByID(ctx context.Context, postID int64) error
}
