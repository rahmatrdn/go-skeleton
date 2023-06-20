package consumer

import (
	"context"
	"fmt"

	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/entity"
	mongoRepo "gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/repository/mongodb"
	moentity "gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/repository/mongodb/entity"
)

type LogQueue struct {
	ctx          context.Context
	logMongoRepo mongoRepo.LogRepository
}

type LogConsumer interface {
	ProcessSyncLog(payload map[string]interface{}) error
}

func NewLogConsumer(
	ctx context.Context,
	logMongoRepo mongoRepo.LogRepository,
) LogConsumer {
	return &LogQueue{ctx, logMongoRepo}
}

func (l *LogQueue) ProcessSyncLog(payload map[string]interface{}) error {
	var params entity.Log
	params.LoadFromMap(payload)

	err := l.logMongoRepo.Create(l.ctx, moentity.LogCollection{
		Status:       string(params.Status),
		FuncName:     params.FuncName,
		ErrorMessage: params.ErrorMessage,
		Process:      params.Process,
		LogFields:    params.LogFields,
	})

	if err != nil {
		fmt.Println("FAILED CREATE LOG TO MONGODB")

		return err
	}

	fmt.Println("SYNC SUCCESS!")
	fmt.Println(params)

	return nil
}
