package repository

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"context"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications"
)

type JobsRepository interface {
	Get(context.Context, specifications.I) (entity.Job, error)
	Create(context.Context, *entity.Job) error
	Update(context.Context, *entity.Job) error
}
