package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	goquSpec "github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewNewsFindUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)
	newsWithFullTextRepo := mock.NewMockNewsWithFullTextRepository(mockCtrl)

	assert.Equal(
		t,
		&newsFindUsecase{
			newsRepo:             newsRepo,
			newsWithFullTextRepo: newsWithFullTextRepo,
		},
		NewNewsFindUsecase(newsRepo, newsWithFullTextRepo),
	)
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

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsFindUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			newsEntities = []entity.News{
				{
					ID:                 1,
					Title:              "News 1",
					Category:           entity.NewsCategoryFinance,
					Status:             entity.NewsStatusAdded,
					NewsWithFullTextID: 101,
				},
				{
					ID:                 2,
					Title:              "News 2",
					Category:           entity.NewsCategoryMilitary,
					Status:             entity.NewsStatusSynced,
					NewsWithFullTextID: 102,
				},
			}

			newsWithFullTextEntities = []entity.NewsWithFullText{
				{
					ID:      101,
					Content: "Content 1",
				},
				{
					ID:      102,
					Content: "Content 2",
				},
			}
		)

		newsRepo.EXPECT().Count(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewsByDate(req.From, req.To, req.Status))),
		).Return(int64(2), nil)

		newsRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewsByDate(req.From, req.To, req.Status))),
			goquSpec.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(newsEntities, nil)

		newsWithFullTextRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsWithFullTextByFileIDs([]uint64{101, 102}, 256))),
			goquSpec.ToPaging(2, 1),
		).Return(newsWithFullTextEntities, nil)

		result, _, err := uc.Execute(ctx, req)

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

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsFindUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}
		)

		newsRepo.EXPECT().Count(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewsByDate(req.From, req.To, req.Status))),
		).Return(int64(0), nil)

		result, _, err := uc.Execute(ctx, req)

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

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsFindUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			repoErr = errors.New("database error")
		)

		newsRepo.EXPECT().Count(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewsByDate(req.From, req.To, req.Status))),
		).Return(int64(0), repoErr)

		_, _, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "count news by date range")
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

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsFindUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			newsEntities = []entity.News{
				{
					ID:                 1,
					Title:              "All News",
					NewsWithFullTextID: 103,
				},
			}

			newsWithFullTextEntities = []entity.NewsWithFullText{
				{
					ID:      103,
					Content: "All content",
				},
			}
		)

		newsRepo.EXPECT().Count(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewsByDate(req.From, req.To, req.Status))),
		).Return(int64(1), nil)

		newsRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewsByDate(req.From, req.To, req.Status))),
			goquSpec.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(newsEntities, nil)

		newsWithFullTextRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsWithFullTextByFileIDs([]uint64{103}, 256))),
			goquSpec.ToPaging(1, 1),
		).Return(newsWithFullTextEntities, nil)

		result, _, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Equal(t, "All News", result[0].Title)
	})
}
