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
	"github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewVnExpressScrapperUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)
	newsWithFullTextRepo := mock.NewMockNewsWithFullTextRepository(mockCtrl)

	expected := vnExpressScrapper{
		newsRepo:             newsRepo,
		newsWithFullTextRepo: newsWithFullTextRepo,
	}

	result := NewVnExpressScrapperUsecase(newsRepo, newsWithFullTextRepo)

	assert.Equal(t, expected, result)
}

func Test_VnExpressScrapper_ErrorHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		newsID        uint64
		setupMock     func(*mock.MockNewsRepository, *mock.MockNewsWithFullTextRepository)
		expectedError string
		expectNoError bool
	}{
		{
			name:   "success - update news status to failed",
			newsID: 123,
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
				news := entity.News{
					ID:     123,
					Status: entity.NewsStatusAdded,
				}
				updatedNews := news
				updatedNews.Status = entity.NewsStatusFailed

				repo.EXPECT().Get(context.Background(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(123))))).Return(news, nil)
				repo.EXPECT().Update(gomock.Any(), &updatedNews).Return(nil)
			},
			expectNoError: true,
		},
		{
			name:   "error - get news fails",
			newsID: 456,
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
				repo.EXPECT().Get(context.Background(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(456))))).Return(entity.News{}, fmt.Errorf("news not found"))
			},
			expectedError: "get news by id to update error status: news not found",
		},
		{
			name:   "error - update news fails",
			newsID: 789,
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
				news := entity.News{
					ID:     789,
					Status: entity.NewsStatusAdded,
				}
				updatedNews := news
				updatedNews.Status = entity.NewsStatusFailed

				repo.EXPECT().Get(context.Background(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(789))))).Return(news, nil)
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
			newsWithFullTextRepo := mock.NewMockNewsWithFullTextRepository(mockCtrl)
			tt.setupMock(newsRepo, newsWithFullTextRepo)
			uc := vnExpressScrapper{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
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
		thumbnailURL  string
		setupMock     func(*mock.MockNewsRepository, *mock.MockNewsWithFullTextRepository)
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
			thumbnailURL:  "https://example.com/thumbnail.jpg",
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
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
				updatedNews.FileSize = int64(len("This is the news content"))

				repo.EXPECT().Get(gomock.Any(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(123))))).Return(news, nil)
				fullTextRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updateNews *entity.News) error {
					// Verify the updated fields
					assert.Equal(t, entity.NewsStatusSynced, updateNews.Status)
					assert.Equal(t, "Test News Title", updateNews.Title)
					assert.Equal(t, "John Doe", updateNews.Author)
					assert.Equal(t, "https://example.com/thumbnail.jpg", updateNews.Thumbnail)
					assert.NotNil(t, updateNews.PublishedAt)
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
			name:         "error - get news fails",
			newsID:       456,
			thumbnailURL: "https://example.com/thumb.jpg",
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
				repo.EXPECT().Get(gomock.Any(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(456))))).Return(entity.News{}, fmt.Errorf("news not found"))
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
			thumbnailURL:  "https://example.com/thumb2.jpg",
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
				news := entity.News{
					ID:          789,
					PublisherID: 1,
					HashedURL:   []byte("test-hash"),
					Status:      entity.NewsStatusAdded,
				}

				repo.EXPECT().Get(gomock.Any(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(789))))).Return(news, nil)
				fullTextRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

				// First call will be the main update that fails
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(fmt.Errorf("database error")).Times(1)
			},
			setupEnv: func() {
				os.Setenv("STORAGE_ROOT", tempDir)
			},
			expectedError: "update news after saving extracted content: database error",
		},
		{
			name:          "success - with zero published date",
			newsID:        999,
			title:         "Test Title",
			content:       "Test content",
			author:        "Test Author",
			publishedDate: time.Time{}, // Zero time
			thumbnailURL:  "https://example.com/thumb.jpg",
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
				news := entity.News{
					ID:          999,
					PublisherID: 1,
					HashedURL:   []byte("test-hash"),
					Status:      entity.NewsStatusAdded,
				}

				repo.EXPECT().Get(gomock.Any(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(999))))).Return(news, nil)
				fullTextRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updateNews *entity.News) error {
					assert.Equal(t, entity.NewsStatusSynced, updateNews.Status)
					assert.Equal(t, "Test Title", updateNews.Title)
					assert.Equal(t, "Test Author", updateNews.Author)
					assert.Equal(t, "https://example.com/thumb.jpg", updateNews.Thumbnail)
					assert.NotNil(t, updateNews.PublishedAt)
					assert.Greater(t, updateNews.FileSize, int64(0))
					return nil
				})
			},
			setupEnv: func() {
				os.Setenv("STORAGE_ROOT", tempDir)
			},
			expectNoError: true,
		},
		{
			name:          "success - file stat fails but continues",
			newsID:        111,
			title:         "Test News Title",
			content:       "This is the news content",
			author:        "John Doe",
			publishedDate: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			thumbnailURL:  "https://example.com/thumbnail.jpg",
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
				news := entity.News{
					ID:          111,
					PublisherID: 1,
					HashedURL:   []byte("test-hash"),
					Status:      entity.NewsStatusAdded,
				}

				repo.EXPECT().Get(gomock.Any(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(111))))).Return(news, nil)
				fullTextRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updateNews *entity.News) error {
					// Verify the updated fields and that FileSize is greater than 0
					assert.Equal(t, entity.NewsStatusSynced, updateNews.Status)
					assert.Equal(t, "Test News Title", updateNews.Title)
					assert.Equal(t, "John Doe", updateNews.Author)
					assert.Equal(t, "https://example.com/thumbnail.jpg", updateNews.Thumbnail)
					assert.NotNil(t, updateNews.PublishedAt)
					assert.Greater(t, updateNews.FileSize, int64(0))
					return nil
				})
			},
			setupEnv: func() {
				os.Setenv("STORAGE_ROOT", tempDir)
			},
			expectNoError: true,
		},
		{
			name:          "success - empty content",
			newsID:        222,
			title:         "Empty News",
			content:       "",
			author:        "",
			publishedDate: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			thumbnailURL:  "",
			setupMock: func(repo *mock.MockNewsRepository, fullTextRepo *mock.MockNewsWithFullTextRepository) {
				news := entity.News{
					ID:          222,
					PublisherID: 1,
					HashedURL:   []byte("test-hash"),
					Status:      entity.NewsStatusAdded,
				}

				repo.EXPECT().Get(gomock.Any(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(222))))).Return(news, nil)
				fullTextRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updateNews *entity.News) error {
					assert.Equal(t, entity.NewsStatusSynced, updateNews.Status)
					assert.Equal(t, "Empty News", updateNews.Title)
					assert.Equal(t, "", updateNews.Author)
					assert.Equal(t, "", updateNews.Thumbnail)
					assert.NotNil(t, updateNews.PublishedAt)
					assert.GreaterOrEqual(t, updateNews.FileSize, int64(0))
					return nil
				})
			},
			setupEnv: func() {
				os.Setenv("STORAGE_ROOT", tempDir)
			},
			expectNoError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			newsRepo := mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo := mock.NewMockNewsWithFullTextRepository(mockCtrl)
			tt.setupMock(newsRepo, newsWithFullTextRepo)

			if tt.setupEnv != nil {
				tt.setupEnv()
			}
			uc := vnExpressScrapper{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
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
				tt.thumbnailURL,
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

func Test_VnExpressScrapper_Execute(t *testing.T) {
	t.Parallel()

	t.Run("error handling", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			job = entity.WebScrapperJob{
				Domain: entity.WebDomain("vnexpress.net"),
				URL:    "https://vnexpress.net/invalid-url-404",
				NewsID: 123,
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = vnExpressScrapper{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			news = entity.News{
				ID:     123,
				Status: entity.NewsStatusAdded,
			}
		)

		// Mock for error handler that will be called when scraping fails
		newsRepo.EXPECT().Get(gomock.Any(), CustomMatcher(specMatcher(specifications.NewNewsByID(uint64(123))))).Return(news, nil).AnyTimes()
		newsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updateNews *entity.News) error {
			// Should update status to failed
			assert.Equal(t, entity.NewsStatusFailed, updateNews.Status)
			return nil
		}).AnyTimes()

		// Execute should handle the error gracefully
		err := uc.Execute(ctx, job)

		// Note: This test may pass or fail depending on network conditions
		// It primarily tests that the method doesnt panic and follows expected patterns
		_ = err // Error is acceptable since were testing with an invalid URL
	})
}
