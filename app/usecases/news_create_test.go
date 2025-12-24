package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	hasherMock "github.com/financial_advisor/app/services/hasher/mock"
	"github.com/financial_advisor/app/services/queue"
	queueMock "github.com/financial_advisor/app/services/queue/mock"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewNewsCreateUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)
	publisherRepo := mock.NewMockPublisherRepository(mockCtrl)
	jobQueue := queueMock.NewMockI(mockCtrl)
	hasher := hasherMock.NewMockI(mockCtrl)

	expected := &newsCreateUsecase{
		newsRepo:      newsRepo,
		publisherRepo: publisherRepo,
		jobQueue:      jobQueue,
		hasher:        hasher,
	}

	assert.Equal(t, expected, NewNewsCreateUsecase(newsRepo, publisherRepo, jobQueue, hasher))
}

func Test_NewsCreateUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
			publisher = entity.Publisher{
				ID:     1,
				Name:   "Example Publisher",
				Domain: "example.com",
			}
		)

		// Mock validation steps
		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, specifications.NewPublisherByDomain("example.com")).Return(publisher, nil)

		// Mock successful create
		newsRepo.EXPECT().Create(ctx, &entity.News{
			URL:         req.URL,
			HashedURL:   hashedURL,
			Status:      entity.NewsStatusAdded,
			Category:    req.Category,
			PublisherID: publisher.ID,
		}).DoAndReturn(func(ctx context.Context, news *entity.News) error {
			news.ID = 1
			return nil
		})

		// Mock successful job enqueue
		expectedJob := entity.WebScrapperJob{
			Domain: entity.WebDomain(publisher.Domain),
			URL:    req.URL,
			NewsID: 1, // This will be set after Create returns
		}
		expectedJobBytes, _ := json.Marshal(expectedJob)
		jobQueue.EXPECT().Enqueue(queue.Message{
			Type: queue.MessageTypeWebScrapper,
			Body: expectedJobBytes,
		}).Return(nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), result.ID)
		assert.Equal(t, req.URL, result.URL)
		assert.Equal(t, entity.NewsStatusAdded.String(), result.Status)
	})

	t.Run("URL too long", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      string(make([]byte, 501)), // URL longer than 500 characters
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}
		)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		var badRequestErr appErrors.SystemError
		assert.True(t, errors.As(err, &badRequestErr))
		assert.Contains(t, badRequestErr.Message(), "url length exceeds the limit")
	})

	t.Run("URL already exists", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(1), nil)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		var conflictErr appErrors.SystemError
		assert.True(t, errors.As(err, &conflictErr))
		assert.Contains(t, conflictErr.Message(), "news with the same URL already exists")
	})

	t.Run("count error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
			countErr  = errors.New("database error")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(0), countErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "count url by hashed value")
	})

	t.Run("invalid URL format", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "invalid-url",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(0), nil)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		var badRequestErr appErrors.SystemError
		assert.True(t, errors.As(err, &badRequestErr))
		assert.Contains(t, badRequestErr.Message(), "invalid URL domain")
	})

	t.Run("publisher not found", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://unknown.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, specifications.NewPublisherByDomain("unknown.com")).Return(entity.Publisher{}, appErrors.ErrNotFound)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		var badRequestErr appErrors.SystemError
		assert.True(t, errors.As(err, &badRequestErr))
		assert.Contains(t, badRequestErr.Message(), "publisher not found")
	})

	t.Run("publisher repo error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL    = []byte("hashed-url")
			publisherErr = errors.New("database error")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, specifications.NewPublisherByDomain("example.com")).Return(entity.Publisher{}, publisherErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "count publisher by id")
	})

	t.Run("create news error - conflicted", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
			publisher = entity.Publisher{
				ID:     1,
				Name:   "Example Publisher",
				Domain: "example.com",
			}
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, specifications.NewPublisherByDomain("example.com")).Return(publisher, nil)
		newsRepo.EXPECT().Create(ctx, &entity.News{
			URL:         req.URL,
			HashedURL:   hashedURL,
			Status:      entity.NewsStatusAdded,
			Category:    req.Category,
			PublisherID: publisher.ID,
		}).Return(appErrors.ErrConflicted)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		var conflictErr appErrors.SystemError
		assert.True(t, errors.As(err, &conflictErr))
		assert.Contains(t, conflictErr.Message(), "news with the same URL already exists")
	})

	t.Run("create news error - generic", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
			publisher = entity.Publisher{
				ID:     1,
				Name:   "Example Publisher",
				Domain: "example.com",
			}
			createErr = errors.New("database error")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, specifications.NewPublisherByDomain("example.com")).Return(publisher, nil)
		newsRepo.EXPECT().Create(ctx, &entity.News{
			URL:         req.URL,
			HashedURL:   hashedURL,
			Status:      entity.NewsStatusAdded,
			Category:    req.Category,
			PublisherID: publisher.ID,
		}).Return(createErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create news")
	})

	t.Run("job queue error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			jobQueue      = queueMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				jobQueue:      jobQueue,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
			publisher = entity.Publisher{
				ID:     1,
				Name:   "Example Publisher",
				Domain: "example.com",
			}
			queueErr = errors.New("queue error")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, specifications.NewNewsByHashedURL(hashedURL)).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, specifications.NewPublisherByDomain("example.com")).Return(publisher, nil)
		newsRepo.EXPECT().Create(ctx, &entity.News{
			URL:         req.URL,
			HashedURL:   hashedURL,
			Status:      entity.NewsStatusAdded,
			Category:    req.Category,
			PublisherID: publisher.ID,
		}).DoAndReturn(func(ctx context.Context, news *entity.News) error {
			news.ID = 1
			return nil
		})
		expectedJob := entity.WebScrapperJob{
			Domain: entity.WebDomain(publisher.Domain),
			URL:    req.URL,
			NewsID: 1,
		}
		expectedJobBytes, _ := json.Marshal(expectedJob)
		jobQueue.EXPECT().Enqueue(queue.Message{
			Type: queue.MessageTypeWebScrapper,
			Body: expectedJobBytes,
		}).Return(queueErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "enqueue web scrapper job")
	})
}

