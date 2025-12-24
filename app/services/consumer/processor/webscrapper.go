// Package processor holds implementations of various job processors
package processor

import (
	"encoding/json"
	"errors"

	"github.com/financial_advisor/app/config"
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/usecases"
)

type WebScrapperProcessor struct {
	mUsecases map[entity.WebDomain]usecases.WebScrapperUsecase
}

func NewWebScrapperProcessor(
	mUsecases map[entity.WebDomain]usecases.WebScrapperUsecase,
) WebScrapperProcessor {
	return WebScrapperProcessor{
		mUsecases: mUsecases,
	}
}

func (p WebScrapperProcessor) Execute(msg []byte) error {
	var job entity.WebScrapperJob

	if err := json.Unmarshal(msg, &job); err != nil {
		return err
	}

	uc, ok := p.mUsecases[job.Domain]
	if !ok {
		return errors.New("usecase for domain not found")
	}

	return uc.Execute(
		config.Get().GlobalCtx(),
		job,
	)
}
