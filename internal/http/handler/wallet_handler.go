package handler

import (
	"net/http"

	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/entity"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/http/middleware"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/parser"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/presenter"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/usecase"

	fiber "github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	parser        parser.Parser
	presenter     presenter.Presenter
	walletUsecase usecase.WalletUsecase
}

func NewWalletHandler(
	parser parser.Parser,
	presenter presenter.Presenter,
	walletUsecase usecase.WalletUsecase,
) *WalletHandler {
	return &WalletHandler{parser, presenter, walletUsecase}
}

func (w *WalletHandler) Register(app fiber.Router) {
	getWalletByID := "/wallet/:id"

	app.Get(getWalletByID, middleware.VerifyJWTToken, w.GetByID)
	app.Get("/wallet", middleware.VerifyJWTToken, w.GetByUserID)
	app.Post("/wallet", middleware.VerifyJWTToken, w.Create)
	app.Put(getWalletByID, middleware.VerifyJWTToken, w.Update)
	app.Delete(getWalletByID, middleware.VerifyJWTToken, w.Delete)
}

// Get selected Wallet
// @Summary			Get Wallet by ID
// @Description		Get Wallet by ID
// @Tags			Wallet
// @Accept			json
// @Produce			json
// @Security 		Bearer
// @Param           id path int true "id of wallet"
// @Success			201 {object} entity.GeneralResponse{data=entity.WalletResponse} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/wallet/{id} [get]
func (w *WalletHandler) GetByID(c *fiber.Ctx) error {
	walletID, err := w.parser.ParserIntIDFromPathParams(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.walletUsecase.GetByID(c.Context(), walletID)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// Get Wallets by User ID
// @Summary			Get Wallets by User ID
// @Description		Get Wallets by User ID
// @Tags			Wallet
// @Accept			json
// @Produce			json
// @Security 		Bearer
// @Success			201 {object} entity.GeneralResponse{data=[]entity.WalletResponse} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/wallet [get]
func (w *WalletHandler) GetByUserID(c *fiber.Ctx) error {
	userID, err := w.parser.ParserUserID(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.walletUsecase.GetByUserID(c.Context(), userID)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// Create New Wallet
// @Summary			Create New Wallet
// @Description		Create New Wallet
// @Tags			Wallet
// @Accept			json
// @Produce			json
// @Security 		Bearer
// @Param			req body entity.WalletReq true "Payload Request Body"
// @Success			201 {object} entity.GeneralResponse{data=entity.WalletReq} "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/wallet [post]
func (w *WalletHandler) Create(c *fiber.Ctx) error {
	var req entity.WalletReq

	err := w.parser.ParserBodyRequestWithUserID(c, &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	data, err := w.walletUsecase.Create(c.Context(), &req)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, data, "Success", http.StatusOK)
}

// Update existing Wallet
// @Summary			Update existing Wallet by ID
// @Description		Update existing Wallet
// @Tags			Wallet
// @Accept			json
// @Produce			json
// @Security 		Bearer
// @Param           id path int true "id of wallet"
// @Param			req body entity.WalletReq true "Payload Request Body"
// @Success			201 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/wallet [put]
func (w *WalletHandler) Update(c *fiber.Ctx) error {
	var walletReq entity.WalletReq
	err := w.parser.ParserBodyWithIntIDPathParamsAndUserID(c, &walletReq)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	err = w.walletUsecase.UpdateByID(c.Context(), walletReq)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, "", "Success", http.StatusOK)
}

// Delete existing Wallet
// @Summary			Delete existing Wallet by ID
// @Description		Delete existing Wallet
// @Tags			Wallet
// @Accept			json
// @Produce			json
// @Security 		Bearer
// @Param           id path int true "id of wallet"
// @Success			201 {object} entity.GeneralResponse "Success"
// @Failure			401 {object} entity.CustomErrorResponse "Unauthorized"
// @Failure			422 {object} entity.CustomErrorResponse "Invalid Request Body"
// @Failure			500 {object} entity.CustomErrorResponse "Internal server Error"
// @Router			/v1/api/wallet/{id} [delete]
func (w *WalletHandler) Delete(c *fiber.Ctx) error {
	walletID, err := w.parser.ParserIntIDFromPathParams(c)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	err = w.walletUsecase.DeleteByID(c.Context(), walletID)
	if err != nil {
		return w.presenter.BuildError(c, err)
	}

	return w.presenter.BuildSuccess(c, nil, "Success", http.StatusOK)
}
