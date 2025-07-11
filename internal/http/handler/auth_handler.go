package handler

import (
	"net/http"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/http/auth"
	"github.com/rahmatrdn/go-skeleton/internal/http/middleware"
	"github.com/rahmatrdn/go-skeleton/internal/parser"
	"github.com/rahmatrdn/go-skeleton/internal/presenter/json"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"

	fiber "github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	parser      parser.Parser
	presenter   json.JsonPresenter
	userUsecase usecase.UserUsecase
}

func NewAuthHandler(
	parser parser.Parser,
	presenter json.JsonPresenter,
	userUsecase usecase.UserUsecase,
) *AuthHandler {
	return &AuthHandler{parser, presenter, userUsecase}
}

func (w *AuthHandler) Register(app fiber.Router) {
	app.Post("/auth/register", w.CreateAsGuest)
	app.Post("/auth/login", w.Login)
	app.Get("/auth/check-token", middleware.VerifyJWTToken, w.CheckToken)
	app.Get("/auth/refresh-token", middleware.VerifyJWTToken, w.RefreshToken)
}

// @Summary			Create User as Guest
// @Description		Create User for Guest
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			req body entity.CreateUserReq true "Payload Request Body"
// @Success			201 {object} entity.GeneralResponse{data=entity.CreateUserResponse} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Invalid Access Token"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Payload Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/auth/register [post]
func (w *AuthHandler) CreateAsGuest(c *fiber.Ctx) error {
	var req *entity.CreateUserReq

	err := w.parser.ParserBodyRequest(c, &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	login, err := w.userUsecase.CreateAsGuest(c.Context(), req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, login, "Success", http.StatusOK)
}

// @Summary			Login
// @Description		Login by using registered account
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			req body entity.LoginReq true "Payload Request Body"
// @Success			201 {object} entity.GeneralResponse{data=entity.LoginResponse} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Invalid Access Token"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Payload Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/auth/login [post]
func (w *AuthHandler) Login(c *fiber.Ctx) error {
	var req *entity.LoginReq

	err := w.parser.ParserBodyRequest(c, &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	login, err := w.userUsecase.VerifyByEmailAndPassword(c.Context(), req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, login, "Success", http.StatusOK)
}

// Check Validation Token
// @Summary			Check Token
// @Description		Check Token
// @Tags			Auth
// @Produce			json
// @Security 		Bearer
// @Success			201 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Invalid Access Token"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Payload Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/api/v1/auth/check-token [get]
func (w *AuthHandler) CheckToken(c *fiber.Ctx) error {
	return w.presenter.BuildSuccess(c, "Token is valid!", "Success", http.StatusOK)
}

func (w *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	newToken, err := auth.RefreshToken(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, newToken, "Success", http.StatusOK)
}
