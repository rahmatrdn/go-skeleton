package handler

import (
	"fmt"
	"net/http"

	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/http/middleware"
	"github.com/rahmatrdn/go-skeleton/internal/parser"
	"github.com/rahmatrdn/go-skeleton/internal/presenter/json"
	comment_usecase "github.com/rahmatrdn/go-skeleton/internal/usecase/comment"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/comment/entity"

	fiber "github.com/gofiber/fiber/v3"
)

type CommentHandler struct {
	parser         parser.Parser
	presenter      json.JsonPresenter
	commentUsecase comment_usecase.ICommentUsecase
}

func NewCommentHandler(
	parser parser.Parser,
	presenter json.JsonPresenter,
	commentUsecase comment_usecase.ICommentUsecase,
) *CommentHandler {
	return &CommentHandler{parser, presenter, commentUsecase}
}

func (w *CommentHandler) Register(app fiber.Router) {
	app.Get("/posts/:post_id/comment", w.GetByPostID)
	app.Get("/comment/:id", w.GetByID)
	app.Post("/comment", middleware.VerifyJWTToken, w.Create)
	app.Put("/comment/:id", middleware.VerifyJWTToken, w.Update)
	app.Delete("/comment/:id", middleware.VerifyJWTToken, w.Delete)
}

// @Summary         Get Comment by Post ID
// @Description     Get all comment for a specific post
// @Tags            Comment
// @Accept          json
// @Produce         json
// @Param           post_id path int true "ID of the Post"
// @Success			200 {object} entity.GeneralResponse{data=[]entity.CommentResponse} "Success"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/posts/{post_id}/comment [get]
func (w *CommentHandler) GetByPostID(c fiber.Ctx) error {
	rawPostID := c.Params("post_id")
	if rawPostID == "" {
		return w.presenter.BuildError(c, fmt.Errorf("PATH PARAM post_id EMPTY"))
	}
	postID := helper.ToInt64(rawPostID)

	data, err := w.commentUsecase.GetByPostID(c.Context(), postID)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// @Summary         Get Comment by ID
// @Description     Get a comment by its ID
// @Tags            Comment
// @Accept          json
// @Produce         json
// @Param           id path int true "ID of the Comment"
// @Success			200 {object} entity.GeneralResponse{data=entity.CommentResponse} "Success"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/comment/{id} [get]
func (w *CommentHandler) GetByID(c fiber.Ctx) error {
	id, err := w.parser.ParserIntIDFromPathParams(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.commentUsecase.GetByID(c.Context(), id)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// @Summary         Create a new Comment
// @Description     Create a new comment on a post. The post must exist. Use parent_id to reply to another comment.
// @Tags            Comment
// @Accept          json
// @Produce         json
// @Security        Bearer
// @Param           req body entity.CommentReq true "Payload Request Body"
// @Success			201 {object} entity.GeneralResponse{data=entity.CommentResponse} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/comment [post]
func (w *CommentHandler) Create(c fiber.Ctx) error {
	var req entity.CommentReq

	err := w.parser.ParserBodyRequestWithUserID(c, &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.commentUsecase.Create(c.Context(), req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusCreated)
}

// @Summary         Update a Comment by ID
// @Description     Update the body of an existing comment
// @Tags            Comment
// @Accept          json
// @Produce         json
// @Security        Bearer
// @Param           id path int true "ID of the Comment"
// @Param           req body entity.CommentReq true "Payload Request Body"
// @Success			200 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/comment/{id} [put]
func (w *CommentHandler) Update(c fiber.Ctx) error {
	var req entity.CommentReq

	err := w.parser.ParserBodyWithIntIDPathParamsAndUserID(c, &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	err = w.commentUsecase.UpdateByID(c.Context(), req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, nil, "Success", http.StatusOK)
}

// @Summary         Delete Comment by ID
// @Description     Delete an existing comment by its ID
// @Tags            Comment
// @Accept          json
// @Produce         json
// @Security        Bearer
// @Param           id path int true "ID of the Comment"
// @Success			200 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/comment/{id} [delete]
func (w *CommentHandler) Delete(c fiber.Ctx) error {
	id, err := w.parser.ParserIntIDFromPathParams(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	err = w.commentUsecase.DeleteByID(c.Context(), id)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, nil, "Success", http.StatusOK)
}
