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
	"go.uber.org/mock/gomock"
)

func Test_NewNewsDeleteUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)

	assert.Equal(t, &newsDeleteUsecase{newsRepo: newsRepo}, NewNewsDeleteUsecase(newsRepo))
}

func Test_NewsDeleteUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success with existing file", func(t *testing.T) {
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
				HashedURL:   []byte("test-hash"),
				PublisherID: 1,
				FilePath:    "",
			}
		)

		// Create the directory and file that StoragePath() would return
		storagePath := news.StoragePath()
		storageDir := news.StorageDir()
		os.MkdirAll(storageDir, 0755)
		os.WriteFile(storagePath, []byte("content"), 0644)
		defer func() {
			os.RemoveAll("scraped") // Clean up the entire test directory
		}()

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(news, nil)
		newsRepo.EXPECT().Delete(ctx, newsID).Return(nil)

		err := uc.Execute(ctx, newsID)

		assert.NoError(t, err)
		// File should be deleted
		_, err = os.Stat(storagePath)
		assert.True(t, os.IsNotExist(err))
	})

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

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(news, nil)
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

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(news, nil)
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

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(entity.News{}, appErrors.ErrNotFound)

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

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(entity.News{}, getErr)

		err := uc.Execute(ctx, newsID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get news by id")
	})

	t.Run("file removal error", func(t *testing.T) {
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
				HashedURL:   []byte("protected"),
				PublisherID: 999, // Will create a path that likely can't be accessed
			}
		)

		// Create a read-only directory to simulate permission error
		storagePath := news.StoragePath()
		storageDir := news.StorageDir()
		os.MkdirAll(storageDir, 0755)
		os.WriteFile(storagePath, []byte("content"), 0644)
		os.Chmod(storageDir, 0444) // Make directory read-only
		defer func() {
			os.Chmod(storageDir, 0755) // Restore permissions for cleanup
			os.RemoveAll("scraped")
		}()

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(news, nil)

		err := uc.Execute(ctx, newsID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "remove news file")
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

		newsRepo.EXPECT().Get(ctx, specifications.NewNewsByID(newsID)).Return(news, nil)
		newsRepo.EXPECT().Delete(ctx, newsID).Return(deleteErr)

		err := uc.Execute(ctx, newsID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find news by date range")
	})
}