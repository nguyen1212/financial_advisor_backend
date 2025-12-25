// Package processor holds implementations of various job processors
package processor

import (
	"encoding/json"

	"github.com/financial_advisor/app/config"
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/usecases"
)

type WebScrapperProcessor struct {
	mUsecases map[entity.WebDomain]usecases.WebScrapperUsecase
	fallback  usecases.FallbackScrapperUsecase
}

func NewWebScrapperProcessor(
	mUsecases map[entity.WebDomain]usecases.WebScrapperUsecase,
	fallback usecases.FallbackScrapperUsecase,
) WebScrapperProcessor {
	return WebScrapperProcessor{
		mUsecases: mUsecases,
		fallback:  fallback,
	}
}

func (p WebScrapperProcessor) Execute(msg []byte) error {
	var job entity.WebScrapperJob

	if err := json.Unmarshal(msg, &job); err != nil {
		return err
	}

	uc, ok := p.mUsecases[job.Domain]
	if !ok {
		return p.fallback.Execute(config.Get().GlobalCtx(), job)
	}

	return uc.Execute(
		config.Get().GlobalCtx(),
		job,
	)
}
