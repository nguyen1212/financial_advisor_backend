package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewNewsFindUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)

	assert.Equal(t, &newsFindUsecase{newsRepo: newsRepo}, NewNewsFindUsecase(newsRepo))
}

func Test_NewsFindUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success with results", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsFindRequest{
				From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsFindUsecase{newsRepo: newsRepo}

			newsEntities = []entity.News{
				{
					ID:       1,
					Title:    "News 1",
					Category: entity.NewsCategoryFinance,
					Status:   entity.NewsStatusAdded,
				},
				{
					ID:       2,
					Title:    "News 2",
					Category: entity.NewsCategoryMilitary,
					Status:   entity.NewsStatusSynced,
				},
			}
		)

		newsRepo.EXPECT().Find(
			ctx,
			specifications.NewNewsByDate(req.From, req.To),
		).Return(newsEntities, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Equal(t, "News 1", result[0].Title)
		assert.Equal(t, uint64(2), result[1].ID)
		assert.Equal(t, "News 2", result[1].Title)
	})

	t.Run("success with empty results", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsFindRequest{
				From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsFindUsecase{newsRepo: newsRepo}
		)

		newsRepo.EXPECT().Find(
			ctx,
			specifications.NewNewsByDate(req.From, req.To),
		).Return([]entity.News{}, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsFindRequest{
				From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsFindUsecase{newsRepo: newsRepo}

			repoErr = errors.New("database error")
		)

		newsRepo.EXPECT().Find(
			ctx,
			specifications.NewNewsByDate(req.From, req.To),
		).Return(nil, repoErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find news by date range")
	})

	t.Run("success with zero time values", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsFindRequest{
				From: time.Time{},
				To:   time.Time{},
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsFindUsecase{newsRepo: newsRepo}

			newsEntities = []entity.News{
				{
					ID:    1,
					Title: "All News",
				},
			}
		)

		newsRepo.EXPECT().Find(
			ctx,
			specifications.NewNewsByDate(req.From, req.To),
		).Return(newsEntities, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Equal(t, "All News", result[0].Title)
	})
}