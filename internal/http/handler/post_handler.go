package handler

import (
	"net/http"

	"github.com/rahmatrdn/go-skeleton/internal/http/middleware"
	"github.com/rahmatrdn/go-skeleton/internal/parser"
	"github.com/rahmatrdn/go-skeleton/internal/presenter/json"
	post_usecase "github.com/rahmatrdn/go-skeleton/internal/usecase/post"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/post/entity"

	fiber "github.com/gofiber/fiber/v3"
)

type PostHandler struct {
	parser      parser.Parser
	presenter   json.JsonPresenter
	postUsecase post_usecase.IPostUsecase
}

func NewPostHandler(
	parser parser.Parser,
	presenter json.JsonPresenter,
	postUsecase post_usecase.IPostUsecase,
) *PostHandler {
	return &PostHandler{parser, presenter, postUsecase}
}

func (w *PostHandler) Register(app fiber.Router) {
	app.Get("/posts", w.GetAll)
	app.Get("/posts/:id", w.GetByID)
	app.Post("/posts", middleware.VerifyJWTToken, w.Create)
	app.Put("/posts/:id", middleware.VerifyJWTToken, w.Update)
	app.Delete("/posts/:id", middleware.VerifyJWTToken, w.Delete)
}

// @Summary         Get all Posts
// @Description     Get a list of all posts
// @Tags            Post
// @Accept          json
// @Produce         json
// @Success			200 {object} entity.GeneralResponse{data=[]entity.PostResponse} "Success"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/posts [get]
func (w *PostHandler) GetAll(c fiber.Ctx) error {
	data, err := w.postUsecase.GetAll(c.Context())
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// @Summary         Get Post by ID
// @Description     Get a Post by its ID
// @Tags            Post
// @Accept          json
// @Produce         json
// @Param           id path int true "ID of the Post"
// @Success			200 {object} entity.GeneralResponse{data=entity.PostResponse} "Success"
// @Failure			404 {object} entity.CustomErrorResponse "Not Found"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/posts/{id} [get]
func (w *PostHandler) GetByID(c fiber.Ctx) error {
	id, err := w.parser.ParserIntIDFromPathParams(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.postUsecase.GetByID(c.Context(), id)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// @Summary         Create a new Post
// @Description     Create a new Post. Slug is auto-generated from the title.
// @Tags            Post
// @Accept          json
// @Produce         json
// @Security        Bearer
// @Param           req body entity.PostReq true "Payload Request Body"
// @Success			201 {object} entity.GeneralResponse{data=entity.PostResponse} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/posts [post]
func (w *PostHandler) Create(c fiber.Ctx) error {
	var req entity.PostReq

	err := w.parser.ParserBodyRequestWithUserID(c, &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.postUsecase.Create(c.Context(), req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusCreated)
}

// @Summary         Update an existing Post by ID
// @Description     Update an existing Post. If the title changes, slug is regenerated automatically.
// @Tags            Post
// @Accept          json
// @Produce         json
// @Security        Bearer
// @Param           id path int true "ID of the Post"
// @Param           req body entity.PostReq true "Payload Request Body"
// @Success			200 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/posts/{id} [put]
func (w *PostHandler) Update(c fiber.Ctx) error {
	var req entity.PostReq

	err := w.parser.ParserBodyWithIntIDPathParamsAndUserID(c, &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	err = w.postUsecase.UpdateByID(c.Context(), req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, nil, "Success", http.StatusOK)
}

// @Summary         Delete Post by ID
// @Description     Delete an existing Post by its ID
// @Tags            Post
// @Accept          json
// @Produce         json
// @Security        Bearer
// @Param           id path int true "ID of the Post"
// @Success			200 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/posts/{id} [delete]
func (w *PostHandler) Delete(c fiber.Ctx) error {
	id, err := w.parser.ParserIntIDFromPathParams(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	err = w.postUsecase.DeleteByID(c.Context(), id)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, nil, "Success", http.StatusOK)
}
