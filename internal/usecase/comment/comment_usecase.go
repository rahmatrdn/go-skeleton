package comment_usecase

import (
	"context"

	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/comment/entity"
)

type CommentUsecase struct {
	commentRepo mysql.ICommentRepository
	postRepo    mysql.IPostRepository
}

func NewCommentUsecase(
	commentRepo mysql.ICommentRepository,
	postRepo mysql.IPostRepository,
) *CommentUsecase {
	return &CommentUsecase{commentRepo, postRepo}
}

type ICommentUsecase interface {
	GetByPostID(ctx context.Context, postID int64) (res []*entity.CommentResponse, err error)
	GetByID(ctx context.Context, commentID int64) (*entity.CommentResponse, error)
	Create(ctx context.Context, req entity.CommentReq) (*entity.CommentResponse, error)
	UpdateByID(ctx context.Context, req entity.CommentReq) error
	DeleteByID(ctx context.Context, commentID int64) error
}
