package handler

import (
	"context"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	wallet "github.com/rahmatrdn/go-skeleton/proto/pb/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletServiceHandler struct {
	wallet.UnimplementedWalletServiceServer
	walletUsecase usecase.WalletUsecase
}

func NewWalletServiceHandler(
	walletUsecase usecase.WalletUsecase,
) *WalletServiceHandler {
	return &WalletServiceHandler{
		wallet.UnimplementedWalletServiceServer{},
		walletUsecase,
	}
}

func (s WalletServiceHandler) Create(ctx context.Context, req *wallet.CreateRequest) (*wallet.CreateResponse, error) {
	walletReq := &entity.WalletReq{
		UserName: req.UserName,
		Balance:  req.Balance,
		UserID:   1242,
	}

	data, err := s.walletUsecase.Create(ctx, walletReq)
	if err != nil {

		return nil, status.Error(codes.InvalidArgument, "Please supply valid name")
	}

	return &wallet.CreateResponse{
		Id:       data.ID,
		UserName: data.UserName,
		Balance:  data.Balance,
	}, nil
}

func (s WalletServiceHandler) Update(ctx context.Context, req *wallet.UpdateRequest) (*wallet.UpdateResponse, error) {
	return &wallet.UpdateResponse{
		Message: "test",
	}, nil
}
