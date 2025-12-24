package usecases

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/financial_advisor/app/domain/repository"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/sirupsen/logrus"
)

type NewsGetUsecase interface {
	Execute(
		ctx context.Context,
		newsID uint64,
	) (dto.News, error)
}

type newsGetUsecase struct {
	newsRepo repository.NewsRepository
}

func NewNewsGetUsecase(
	newsRepo repository.NewsRepository,
) NewsGetUsecase {
	return &newsGetUsecase{
		newsRepo: newsRepo,
	}
}

func (uc *newsGetUsecase) Execute(
	ctx context.Context,
	newsID uint64,
) (dto.News, error) {
	news, err := uc.newsRepo.Get(
		ctx,
		specifications.NewNewsByID(
			newsID,
		),
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return dto.News{}, appErrors.NewErrorBadRequest(
				appErrors.ErrorCodeNewsNotFound,
				"news not found",
			)
		}

		return dto.News{}, fmt.Errorf("find news by date range: %w", err)
	}

	if news.FilePath == "" {
		return dto.ToDtoNews(news), nil
	}

	file, err := os.ReadFile(news.FilePath)
	if err != nil {
		logrus.WithError(err).WithField("news_id", news.ID).Errorln("file does not exist")
	}

	news.Content = string(file)

	return dto.ToDtoNews(news), nil
}
