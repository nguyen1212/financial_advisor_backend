package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	appErrors "github.com/financial_advisor/app/errors"
	goquSpec "github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewNewsGetUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)
	newsWithFullTextRepo := mock.NewMockNewsWithFullTextRepository(mockCtrl)

	assert.Equal(
		t,
		&newsGetUsecase{
			newsRepo:             newsRepo,
			newsWithFullTextRepo: newsWithFullTextRepo,
		},
		NewNewsGetUsecase(newsRepo, newsWithFullTextRepo),
	)
}

func Test_NewsGetUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("news not found", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsGetUsecase{newsRepo: newsRepo}

			wantErr = appErrors.NewErrorNotFound(
				appErrors.ErrorCodeNewsNotFound,
				"news not found",
			)
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(goquSpec.NewNewsByID(newsID)))).Return(entity.News{}, appErrors.ErrNotFound)

		_, err := uc.Execute(ctx, newsID, dto.NewsGetRequest{})

		assert.Equal(t, wantErr, err)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsGetUsecase{newsRepo: newsRepo}

			repoErr = errors.New("database error")
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(goquSpec.NewNewsByID(newsID)))).Return(entity.News{}, repoErr)

		_, err := uc.Execute(ctx, newsID, dto.NewsGetRequest{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find news by date range")
	})

	t.Run("success without file path", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsGetUsecase{newsRepo: newsRepo}

			news = entity.News{
				ID:       newsID,
				Title:    "Test News",
				FilePath: "",
			}
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(goquSpec.NewNewsByID(newsID)))).Return(news, nil)

		result, err := uc.Execute(ctx, newsID, dto.NewsGetRequest{})

		assert.NoError(t, err)
		assert.Equal(t, newsID, result.ID)
		assert.Equal(t, "Test News", result.Title)
	})

	t.Run("success with FileID > 0", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1
			req           = dto.NewsGetRequest{
				HighlightKeywords: []string{"keyword1", "keyword2"},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsGetUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			news = entity.News{
				ID:                 newsID,
				Title:              "Test News with Content",
				NewsWithFullTextID: 123,
				FileSize:           1024,
				Category:           entity.NewsCategoryFinance,
				Status:             entity.NewsStatusSynced,
			}

			newsWithFullText = entity.NewsWithFullText{
				ID:      123,
				Content: "This is the full content with <mark>keyword1</mark> and <mark>keyword2</mark>",
			}
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(goquSpec.NewNewsByID(newsID)))).Return(news, nil)
		newsWithFullTextRepo.EXPECT().Get(
			ctx,
			CustomMatcher(SpecMatcher(goquSpec.NewNewsWithFullTextByFileID(news.NewsWithFullTextID, news.FileSize, req.HighlightKeywords))),
		).Return(newsWithFullText, nil)

		result, err := uc.Execute(ctx, newsID, req)

		assert.NoError(t, err)
		assert.Equal(t, newsID, result.ID)
		assert.Equal(t, "Test News with Content", result.Title)
		assert.Equal(t, "This is the full content with <mark>keyword1</mark> and <mark>keyword2</mark>", result.Content)
	})

	t.Run("success with FileID > 0 without highlight keywords", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1
			req           = dto.NewsGetRequest{
				HighlightKeywords: []string{},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsGetUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			news = entity.News{
				ID:                 newsID,
				Title:              "Test News Plain Content",
				NewsWithFullTextID: 456,
				FileSize:           2048,
				Category:           entity.NewsCategoryMilitary,
				Status:             entity.NewsStatusSynced,
			}

			newsWithFullText = entity.NewsWithFullText{
				ID:      456,
				Content: "Plain content without highlighting",
			}
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(goquSpec.NewNewsByID(newsID)))).Return(news, nil)
		newsWithFullTextRepo.EXPECT().Get(
			ctx,
			CustomMatcher(SpecMatcher(goquSpec.NewNewsWithFullTextByFileID(news.NewsWithFullTextID, news.FileSize, req.HighlightKeywords))),
		).Return(newsWithFullText, nil)

		result, err := uc.Execute(ctx, newsID, req)

		assert.NoError(t, err)
		assert.Equal(t, newsID, result.ID)
		assert.Equal(t, "Test News Plain Content", result.Title)
		assert.Equal(t, "Plain content without highlighting", result.Content)
	})

	t.Run("error when getting full text content", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1
			req           = dto.NewsGetRequest{
				HighlightKeywords: []string{"test"},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsGetUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			news = entity.News{
				ID:                 newsID,
				Title:              "Test News",
				NewsWithFullTextID: 789,
				FileSize:           512,
			}

			fullTextErr = errors.New("full text repository error")
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(goquSpec.NewNewsByID(newsID)))).Return(news, nil)
		newsWithFullTextRepo.EXPECT().Get(
			ctx,
			CustomMatcher(SpecMatcher(goquSpec.NewNewsWithFullTextByFileID(news.NewsWithFullTextID, news.FileSize, req.HighlightKeywords))),
		).Return(entity.NewsWithFullText{}, fullTextErr)

		_, err := uc.Execute(ctx, newsID, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find news with content by document id")
	})

	t.Run("success with zero FileSize", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1
			req           = dto.NewsGetRequest{
				HighlightKeywords: []string{"keyword"},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsGetUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			news = entity.News{
				ID:                 newsID,
				Title:              "Test News Zero FileSize",
				NewsWithFullTextID: 999,
				FileSize:           0, // Zero file size
			}

			newsWithFullText = entity.NewsWithFullText{
				ID:      999,
				Content: "Content with zero file size",
			}
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(goquSpec.NewNewsByID(newsID)))).Return(news, nil)
		newsWithFullTextRepo.EXPECT().Get(
			ctx,
			CustomMatcher(SpecMatcher(goquSpec.NewNewsWithFullTextByFileID(news.NewsWithFullTextID, news.FileSize, req.HighlightKeywords))),
		).Return(newsWithFullText, nil)

		result, err := uc.Execute(ctx, newsID, req)

		assert.NoError(t, err)
		assert.Equal(t, newsID, result.ID)
		assert.Equal(t, "Test News Zero FileSize", result.Title)
		assert.Equal(t, "Content with zero file size", result.Content)
	})
}
