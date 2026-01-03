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

func Test_NewNewsDeleteUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)
	newsWithFullTextRepo := mock.NewMockNewsWithFullTextRepository(mockCtrl)

	assert.Equal(
		t,
		&newsDeleteUsecase{
			newsRepo:             newsRepo,
			newsWithFullTextRepo: newsWithFullTextRepo,
		},
		NewNewsDeleteUsecase(
			newsRepo,
			newsWithFullTextRepo,
		),
	)
}

func Test_NewsDeleteUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success with non-existing file", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsDeleteUsecase{newsRepo: newsRepo}

			news = entity.News{
				ID:          newsID,
				HashedURL:   []byte("nonexistent"),
				PublisherID: 999,
			}
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByID(newsID)))).Return(news, nil)
		newsRepo.EXPECT().Delete(ctx, newsID).Return(nil)

		err := uc.Execute(ctx, newsID)

		assert.NoError(t, err) // Should not fail for non-existent files
	})

	t.Run("success with empty file path", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsDeleteUsecase{newsRepo: newsRepo}

			news = entity.News{
				ID:          newsID,
				HashedURL:   []byte{},
				PublisherID: 0,
			}
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByID(newsID)))).Return(news, nil)
		newsRepo.EXPECT().Delete(ctx, newsID).Return(nil)

		err := uc.Execute(ctx, newsID)

		assert.NoError(t, err)
	})

	t.Run("news not found", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsDeleteUsecase{newsRepo: newsRepo}
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByID(newsID)))).Return(entity.News{}, appErrors.ErrNotFound)

		err := uc.Execute(ctx, newsID)

		assert.NoError(t, err) // Should return nil for not found
	})

	t.Run("get news error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsDeleteUsecase{newsRepo: newsRepo}

			getErr = errors.New("database error")
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByID(newsID)))).Return(entity.News{}, getErr)

		err := uc.Execute(ctx, newsID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get news by id")
	})

	t.Run("delete news error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsDeleteUsecase{newsRepo: newsRepo}

			news = entity.News{
				ID:          newsID,
				HashedURL:   []byte{},
				PublisherID: 0,
			}
			deleteErr = errors.New("database error")
		)

		newsRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByID(newsID)))).Return(news, nil)
		newsRepo.EXPECT().Delete(ctx, newsID).Return(deleteErr)

		err := uc.Execute(ctx, newsID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find news by date range")
	})
}
