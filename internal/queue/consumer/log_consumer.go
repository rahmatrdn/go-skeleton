package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	mongoRepo "github.com/rahmatrdn/go-skeleton/internal/repository/mongodb"
	moentity "github.com/rahmatrdn/go-skeleton/internal/repository/mongodb/entity"
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

	var executionTime string
	if params.LogFields["execution_time"] != "" {
		executionTime = params.LogFields["execution_time"]
	}

	err := l.logMongoRepo.Create(l.ctx, moentity.LogCollection{
		Status:        string(params.Status),
		FuncName:      params.FuncName,
		ErrorMessage:  params.ErrorMessage,
		Process:       params.Process,
		LogFields:     params.LogFields,
		Created:       time.Now().UTC().Add(7 * time.Hour),
		ExecutionTime: helper.ToInt(executionTime),
	})

	if err != nil {
		fmt.Println("FAILED CREATE LOG TO MONGODB")

		return err
	}

	fmt.Println("SYNC SUCCESS!")
	fmt.Println(params)

	return nil
}
