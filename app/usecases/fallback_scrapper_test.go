package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	appErrors "github.com/financial_advisor/app/errors"
	goquSpec "github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewFallbackScrapperUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)

	expected := fallbackScrapperUsecase{newsRepo: newsRepo}

	assert.Equal(t, expected, NewFallbackScrapperUsecase(newsRepo))
}

func Test_FallbackScrapperUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			job = entity.WebScrapperJob{
				Domain: entity.WebDomain("example.com"),
				URL:    "https://example.com/news/article-1",
				NewsID: 1,
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = fallbackScrapperUsecase{newsRepo: newsRepo}

			news = entity.News{
				ID:          1,
				Title:       "Test News",
				Author:      "Test Author",
				URL:         "https://example.com/news/article-1",
				Status:      entity.NewsStatusAdded,
				Category:    entity.NewsCategoryFinance,
				PublisherID: 1,
			}
		)

		// Mock getting news by ID
		newsRepo.EXPECT().Get(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByID(job.NewsID))),
		).Return(news, nil)

		// Mock updating news status to failed
		expectedNews := news
		expectedNews.Status = entity.NewsStatusFailed
		newsRepo.EXPECT().Update(ctx, &expectedNews).Return(nil)

		err := uc.Execute(ctx, job)

		assert.NoError(t, err)
	})

	t.Run("failed - news not found", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			job = entity.WebScrapperJob{
				Domain: entity.WebDomain("example.com"),
				URL:    "https://example.com/news/article-1",
				NewsID: 999,
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = fallbackScrapperUsecase{newsRepo: newsRepo}
		)

		// Mock news not found
		newsRepo.EXPECT().Get(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByID(job.NewsID))),
		).Return(entity.News{}, appErrors.ErrNotFound)

		err := uc.Execute(ctx, job)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get news by id")
	})

	t.Run("failed - database error on get", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			job = entity.WebScrapperJob{
				Domain: entity.WebDomain("example.com"),
				URL:    "https://example.com/news/article-1",
				NewsID: 1,
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = fallbackScrapperUsecase{newsRepo: newsRepo}

			dbErr = errors.New("database connection error")
		)

		// Mock database error on get
		newsRepo.EXPECT().Get(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByID(job.NewsID))),
		).Return(entity.News{}, dbErr)

		err := uc.Execute(ctx, job)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get news by id")
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("failed - update error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			job = entity.WebScrapperJob{
				Domain: entity.WebDomain("example.com"),
				URL:    "https://example.com/news/article-1",
				NewsID: 1,
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = fallbackScrapperUsecase{newsRepo: newsRepo}

			news = entity.News{
				ID:          1,
				Title:       "Test News",
				Author:      "Test Author",
				URL:         "https://example.com/news/article-1",
				Status:      entity.NewsStatusAdded,
				Category:    entity.NewsCategoryFinance,
				PublisherID: 1,
			}

			updateErr = errors.New("database update error")
		)

		// Mock successful get
		newsRepo.EXPECT().Get(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByID(job.NewsID))),
		).Return(news, nil)

		// Mock update failure
		expectedNews := news
		expectedNews.Status = entity.NewsStatusFailed
		newsRepo.EXPECT().Update(ctx, &expectedNews).Return(updateErr)

		err := uc.Execute(ctx, job)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update news status to failed")
		assert.ErrorIs(t, err, updateErr)
	})

	t.Run("success - different job parameters", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			job = entity.WebScrapperJob{
				Domain: entity.WebDomain("vnexpress.net"),
				URL:    "https://vnexpress.net/kinh-te/tai-chinh/article-123",
				NewsID: 42,
			}

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = fallbackScrapperUsecase{newsRepo: newsRepo}

			news = entity.News{
				ID:          42,
				Title:       "Vietnam Economic News",
				Author:      "VnExpress Author",
				URL:         "https://vnexpress.net/kinh-te/tai-chinh/article-123",
				Status:      entity.NewsStatusSynced,
				Category:    entity.NewsCategoryFinance,
				PublisherID: 2,
			}
		)

		// Mock getting news by ID
		newsRepo.EXPECT().Get(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByID(job.NewsID))),
		).Return(news, nil)

		// Mock updating news status to failed
		expectedNews := news
		expectedNews.Status = entity.NewsStatusFailed
		newsRepo.EXPECT().Update(ctx, &expectedNews).Return(nil)

		err := uc.Execute(ctx, job)

		assert.NoError(t, err)
	})

	t.Run("success - news with different statuses", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name          string
			initialStatus entity.NewsStatus
			finalStatus   entity.NewsStatus
		}{
			{
				name:          "synced to failed",
				initialStatus: entity.NewsStatusSynced,
				finalStatus:   entity.NewsStatusFailed,
			},
			{
				name:          "added to failed",
				initialStatus: entity.NewsStatusAdded,
				finalStatus:   entity.NewsStatusFailed,
			},
			{
				name:          "unknown to failed",
				initialStatus: entity.NewsStatusUnknown,
				finalStatus:   entity.NewsStatusFailed,
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				var (
					ctx = context.Background()
					job = entity.WebScrapperJob{
						Domain: entity.WebDomain("example.com"),
						URL:    "https://example.com/news/article-1",
						NewsID: 1,
					}

					newsRepo = mock.NewMockNewsRepository(mockCtrl)
					uc       = fallbackScrapperUsecase{newsRepo: newsRepo}

					news = entity.News{
						ID:          1,
						Title:       "Test News",
						Author:      "Test Author",
						URL:         "https://example.com/news/article-1",
						Status:      tc.initialStatus,
						Category:    entity.NewsCategoryFinance,
						PublisherID: 1,
					}
				)

				// Mock getting news by ID
				newsRepo.EXPECT().Get(
					ctx,
					CustomMatcher(specMatcher(goquSpec.NewNewsByID(job.NewsID))),
				).Return(news, nil)

				// Mock updating news status to failed
				expectedNews := news
				expectedNews.Status = tc.finalStatus
				newsRepo.EXPECT().Update(ctx, &expectedNews).Return(nil)

				err := uc.Execute(ctx, job)

				assert.NoError(t, err)
			})
		}
	})
}

