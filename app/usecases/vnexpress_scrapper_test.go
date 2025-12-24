package usecases

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewVnExpressScrapperUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)

	expected := vnExpressScrapper{
		newsRepo: newsRepo,
	}

	result := NewVnExpressScrapperUsecase(newsRepo)
	assert.Equal(t, expected, result)
}

func Test_VnExpressScrapper_ErrorHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		newsID        uint64
		setupMock     func(*mock.MockNewsRepository)
		expectedError string
		expectNoError bool
	}{
		{
			name:   "success - update news status to failed",
			newsID: 123,
			setupMock: func(repo *mock.MockNewsRepository) {
				news := entity.News{
					ID:     123,
					Status: entity.NewsStatusAdded,
				}
				updatedNews := news
				updatedNews.Status = entity.NewsStatusFailed

				repo.EXPECT().Get(context.Background(), specifications.NewNewsByID(uint64(123))).Return(news, nil)
				repo.EXPECT().Update(gomock.Any(), &updatedNews).Return(nil)
			},
			expectNoError: true,
		},
		{
			name:   "error - get news fails",
			newsID: 456,
			setupMock: func(repo *mock.MockNewsRepository) {
				repo.EXPECT().Get(context.Background(), specifications.NewNewsByID(uint64(456))).Return(entity.News{}, fmt.Errorf("news not found"))
			},
			expectedError: "get news by id to update error status: news not found",
		},
		{
			name:   "error - update news fails",
			newsID: 789,
			setupMock: func(repo *mock.MockNewsRepository) {
				news := entity.News{
					ID:     789,
					Status: entity.NewsStatusAdded,
				}
				updatedNews := news
				updatedNews.Status = entity.NewsStatusFailed

				repo.EXPECT().Get(context.Background(), specifications.NewNewsByID(uint64(789))).Return(news, nil)
				repo.EXPECT().Update(gomock.Any(), &updatedNews).Return(fmt.Errorf("database error"))
			},
			expectedError: "update news status to failed: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			newsRepo := mock.NewMockNewsRepository(mockCtrl)
			tt.setupMock(newsRepo)

			uc := vnExpressScrapper{
				newsRepo: newsRepo,
			}

			err := uc.errorHandler(context.Background(), tt.newsID)

			if tt.expectNoError {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func Test_VnExpressScrapper_SaveFile(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name          string
		newsID        uint64
		title         string
		content       string
		author        string
		publishedDate time.Time
		setupMock     func(*mock.MockNewsRepository)
		setupEnv      func()
		expectedError string
		expectNoError bool
	}{
		{
			name:          "success - save file and update news",
			newsID:        123,
			title:         "Test News Title",
			content:       "This is the news content",
			author:        "John Doe",
			publishedDate: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			setupMock: func(repo *mock.MockNewsRepository) {
				news := entity.News{
					ID:          123,
					PublisherID: 1,
					HashedURL:   []byte("test-hash"),
					Status:      entity.NewsStatusAdded,
				}

				// Expected updated news
				updatedNews := news
				updatedNews.Status = entity.NewsStatusSynced
				updatedNews.Title = "Test News Title"
				updatedNews.Author = "John Doe"
				publishedDate := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
				updatedNews.PublishedAt = &publishedDate
				updatedNews.FilePath = news.StoragePath()
				updatedNews.FileSize = int64(len("This is the news content"))

				repo.EXPECT().Get(gomock.Any(), specifications.NewNewsByID(uint64(123), "Publisher")).Return(news, nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updateNews *entity.News) error {
					// Verify the updated fields
					assert.Equal(t, entity.NewsStatusSynced, updateNews.Status)
					assert.Equal(t, "Test News Title", updateNews.Title)
					assert.Equal(t, "John Doe", updateNews.Author)
					assert.NotNil(t, updateNews.PublishedAt)
					assert.NotEmpty(t, updateNews.FilePath)
					assert.Greater(t, updateNews.FileSize, int64(0))
					return nil
				})
			},
			setupEnv: func() {
				// Set temp directory for storage
				os.Setenv("STORAGE_ROOT", tempDir)
			},
			expectNoError: true,
		},
		{
			name:   "error - get news fails",
			newsID: 456,
			setupMock: func(repo *mock.MockNewsRepository) {
				repo.EXPECT().Get(gomock.Any(), specifications.NewNewsByID(uint64(456), "Publisher")).Return(entity.News{}, fmt.Errorf("news not found"))
			},
			expectedError: "get news by id to save extracted content: news not found",
		},
		{
			name:          "error - update news fails",
			newsID:        789,
			title:         "Test Title",
			content:       "Test content",
			author:        "Test Author",
			publishedDate: time.Now(),
			setupMock: func(repo *mock.MockNewsRepository) {
				news := entity.News{
					ID:          789,
					PublisherID: 1,
					HashedURL:   []byte("test-hash"),
					Status:      entity.NewsStatusAdded,
				}

				repo.EXPECT().Get(gomock.Any(), specifications.NewNewsByID(uint64(789), "Publisher")).Return(news, nil)

				// First call will be the main update that fails
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(fmt.Errorf("database error")).Times(1)

				// Second call will be from the defer block to set status to failed
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updateNews *entity.News) error {
					assert.Equal(t, entity.NewsStatusFailed, updateNews.Status)
					return nil
				}).Times(1)
			},
			setupEnv: func() {
				os.Setenv("STORAGE_ROOT", tempDir)
			},
			expectedError: "update news after saving extracted content: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			newsRepo := mock.NewMockNewsRepository(mockCtrl)
			tt.setupMock(newsRepo)

			if tt.setupEnv != nil {
				tt.setupEnv()
			}

			uc := vnExpressScrapper{
				newsRepo: newsRepo,
			}

			content := &strings.Builder{}
			content.WriteString(tt.content)

			err := uc.saveFile(
				context.Background(),
				tt.newsID,
				tt.title,
				content,
				tt.author,
				tt.publishedDate,
			)

			if tt.expectNoError {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

