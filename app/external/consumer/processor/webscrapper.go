// Package processor holds implementations of various job processors
package processor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/usecases"
)

type webScrapperProcessor struct {
	mUsecases map[entity.WebDomain]usecases.WebScrapperUsecase
	fallback  usecases.FallbackScrapperUsecase
}

func NewWebScrapperProcessor(
	mUsecases map[entity.WebDomain]usecases.WebScrapperUsecase,
	fallback usecases.FallbackScrapperUsecase,
) consumer.Processor {
	return &webScrapperProcessor{
		mUsecases: mUsecases,
		fallback:  fallback,
	}
}

func (p *webScrapperProcessor) Execute(ctx context.Context, msg []byte) error {
	var job entity.WebScrapperJob

	err := json.Unmarshal(msg, &job)
	if err != nil {
		return err
	}

	// TODO: if there is a need to adjust timeout, consider making it configurable
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	uc, ok := p.mUsecases[job.Domain]
	if !ok {
		return p.fallback.Execute(timeoutCtx, job)
	}

	err = uc.Execute(
		timeoutCtx,
		job,
	)

	return err
}
