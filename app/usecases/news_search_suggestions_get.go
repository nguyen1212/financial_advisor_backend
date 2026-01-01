package usecases

import (
	"context"
	"fmt"

	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
)

type NewsSearchSuggestionsGetUsecase interface {
	Execute(
		ctx context.Context,
		req dto.NewsSearchSuggestionsRequest,
	) ([]string, error)
}

type newsSearchSuggestionsGetUsecase struct {
	newsWithFullTextRepo repository.NewsWithFullTextRepository
}

func NewNewsSearchSuggestionsGetUsecase(
	newsWithFullTextRepo repository.NewsWithFullTextRepository,
) NewsSearchSuggestionsGetUsecase {
	return &newsSearchSuggestionsGetUsecase{
		newsWithFullTextRepo: newsWithFullTextRepo,
	}
}

func (uc *newsSearchSuggestionsGetUsecase) Execute(
	ctx context.Context,
	req dto.NewsSearchSuggestionsRequest,
) ([]string, error) {
	if len(req.Keywords) == 0 {
		return []string{}, nil
	}

	suggestions, err := uc.newsWithFullTextRepo.FindSearchSuggestions(
		ctx,
		specifications.NewNewsSearchSuggestions(
			req.Keywords,
			specifications.StrongFuzziness,
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get news search suggestions: %w", err)
	}

	return suggestions, nil
}
