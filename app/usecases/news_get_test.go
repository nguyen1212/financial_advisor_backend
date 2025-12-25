package usecases

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_NewNewsGetUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)

	assert.Equal(t, &newsGetUsecase{newsRepo: newsRepo}, NewNewsGetUsecase(newsRepo))
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

			wantErr = appErrors.NewErrorBadRequest(
				appErrors.ErrorCodeNewsNotFound,
				"news not found",
			)
		)

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(entity.News{}, appErrors.ErrNotFound)

		_, err := uc.Execute(ctx, newsID)

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

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(entity.News{}, repoErr)

		_, err := uc.Execute(ctx, newsID)

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

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(news, nil)

		result, err := uc.Execute(ctx, newsID)

		assert.NoError(t, err)
		assert.Equal(t, newsID, result.ID)
		assert.Equal(t, "Test News", result.Title)
	})

	t.Run("success with file path", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		// Create a temporary file
		tmpFile := "/tmp/test_news_content.txt"
		content := "This is test news content"
		os.WriteFile(tmpFile, []byte(content), 0644)
		defer os.Remove(tmpFile)

		var (
			ctx           = context.Background()
			newsID uint64 = 1

			newsRepo = mock.NewMockNewsRepository(mockCtrl)
			uc       = &newsGetUsecase{newsRepo: newsRepo}

			news = entity.News{
				ID:       newsID,
				Title:    "Test News",
				FilePath: tmpFile,
			}
		)

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(news, nil)

		result, err := uc.Execute(ctx, newsID)

		assert.NoError(t, err)
		assert.Equal(t, newsID, result.ID)
		assert.Equal(t, "Test News", result.Title)
		assert.Equal(t, content, result.Content)
	})

	t.Run("file read error", func(t *testing.T) {
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
				FilePath: "/nonexistent/file.txt",
			}
		)

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(news, nil)

		_, err := uc.Execute(ctx, newsID)

		require.NoError(t, err)
	})
}
