package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	appErrors "github.com/financial_advisor/app/errors"
	goquSpec "github.com/financial_advisor/app/external/db/goqu/specifications"
	hasherMock "github.com/financial_advisor/app/services/hasher/mock"
	"github.com/financial_advisor/app/services/queue"
	workerMock "github.com/financial_advisor/app/services/worker/mock"
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
	worker := workerMock.NewMockI(mockCtrl)
	hasher := hasherMock.NewMockI(mockCtrl)

	expected := &newsCreateUsecase{
		newsRepo:      newsRepo,
		publisherRepo: publisherRepo,
		worker:        worker,
		hasher:        hasher,
	}

	assert.Equal(t, expected, NewNewsCreateUsecase(newsRepo, publisherRepo, worker, hasher))
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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
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
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("example.com")))).Return(publisher, nil)

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
		expectedMessage := queue.Message{
			Type: queue.MessageTypeWebScrapper,
			Body: expectedJobBytes,
		}
		expectedMessageBytes, _ := json.Marshal(expectedMessage)
		worker.EXPECT().Run(expectedMessageBytes).Return(nil)

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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(1), nil)

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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
			countErr  = errors.New("database error")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), countErr)

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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		var badRequestErr appErrors.SystemError
		assert.True(t, errors.As(err, &badRequestErr))
		assert.Contains(t, badRequestErr.Message(), "invalid URL domain")
	})

	t.Run("success - create new publisher", func(t *testing.T) {
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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url")
		)

		// Mock validation steps
		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL))),
		).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("unknown.com"))),
		).Return(entity.Publisher{}, appErrors.ErrNotFound)

		// Mock publisher creation
		publisherRepo.EXPECT().Create(ctx, &entity.Publisher{
			Name:   "unknown.com",
			Domain: "unknown.com",
		}).DoAndReturn(func(ctx context.Context, publisher *entity.Publisher) error {
			publisher.ID = 2
			return nil
		})

		// Mock successful news create
		newsRepo.EXPECT().Create(ctx, &entity.News{
			URL:         req.URL,
			HashedURL:   hashedURL,
			Status:      entity.NewsStatusAdded,
			Category:    req.Category,
			PublisherID: 2, // New publisher ID
		}).DoAndReturn(func(ctx context.Context, news *entity.News) error {
			news.ID = 1
			return nil
		})

		// Mock successful job enqueue
		expectedJob := entity.WebScrapperJob{
			Domain: entity.WebDomain("unknown.com"),
			URL:    req.URL,
			NewsID: 1,
		}
		expectedJobBytes, _ := json.Marshal(expectedJob)
		expectedMessage := queue.Message{
			Type: queue.MessageTypeWebScrapper,
			Body: expectedJobBytes,
		}
		expectedMessageBytes, _ := json.Marshal(expectedMessage)
		worker.EXPECT().Run(expectedMessageBytes).Return(nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), result.ID)
		assert.Equal(t, req.URL, result.URL)
		assert.Equal(t, entity.NewsStatusAdded.String(), result.Status)
	})

	t.Run("error - publisher creation fails", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://newpublisher.com/news/article-1",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL    = []byte("hashed-url")
			publisherErr = errors.New("database error")
		)

		// Mock validation steps
		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(
			gomock.Any(),
			CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL))),
		).Return(int64(0), nil)

		publisherRepo.EXPECT().Get(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("newpublisher.com"))),
		).Return(entity.Publisher{}, appErrors.ErrNotFound)

		// Mock publisher creation failure
		publisherRepo.EXPECT().Create(ctx, &entity.Publisher{
			Name:   "newpublisher.com",
			Domain: "newpublisher.com",
		}).Return(publisherErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create publisher")
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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL    = []byte("hashed-url")
			publisherErr = errors.New("database error")
		)

		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("example.com")))).Return(entity.Publisher{}, publisherErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get publisher by id")
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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
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
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("example.com")))).Return(publisher, nil)
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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
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
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("example.com")))).Return(publisher, nil)
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
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
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
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("example.com")))).Return(publisher, nil)
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
		expectedMessage := queue.Message{
			Type: queue.MessageTypeWebScrapper,
			Body: expectedJobBytes,
		}
		expectedMessageBytes, _ := json.Marshal(expectedMessage)
		worker.EXPECT().Run(expectedMessageBytes).Return(queueErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "enqueue web scrapper job")
	})

	t.Run("success - with different news category", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/military/defense-news",
				Category: entity.NewsCategoryMilitary,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url-military")
			publisher = entity.Publisher{
				ID:     1,
				Name:   "Example Publisher",
				Domain: "example.com",
			}
		)

		// Mock validation steps
		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("example.com")))).Return(publisher, nil)

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
			NewsID: 1,
		}
		expectedJobBytes, _ := json.Marshal(expectedJob)
		expectedMessage := queue.Message{
			Type: queue.MessageTypeWebScrapper,
			Body: expectedJobBytes,
		}
		expectedMessageBytes, _ := json.Marshal(expectedMessage)
		worker.EXPECT().Run(expectedMessageBytes).Return(nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), result.ID)
		assert.Equal(t, req.URL, result.URL)
		assert.Equal(t, entity.NewsCategoryMilitary.String(), result.Category)
		assert.Equal(t, entity.NewsStatusAdded.String(), result.Status)
	})

	t.Run("success - URL with query parameters", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://example.com/news?id=123&utm_source=test",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-url-with-params")
			publisher = entity.Publisher{
				ID:     1,
				Name:   "Example Publisher",
				Domain: "example.com",
			}
		)

		// Mock validation steps
		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)
		publisherRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("example.com")))).Return(publisher, nil)

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
			NewsID: 1,
		}
		expectedJobBytes, _ := json.Marshal(expectedJob)
		expectedMessage := queue.Message{
			Type: queue.MessageTypeWebScrapper,
			Body: expectedJobBytes,
		}
		expectedMessageBytes, _ := json.Marshal(expectedMessage)
		worker.EXPECT().Run(expectedMessageBytes).Return(nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.URL, result.URL) // Should preserve full URL with params
		assert.Equal(t, entity.NewsStatusAdded.String(), result.Status)
	})

	t.Run("success - subdomain URL", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsCreateRequest{
				URL:      "https://news.subdomain.example.com/article",
				Category: entity.NewsCategoryFinance,
			}

			newsRepo      = mock.NewMockNewsRepository(mockCtrl)
			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			worker        = workerMock.NewMockI(mockCtrl)
			hasher        = hasherMock.NewMockI(mockCtrl)
			uc            = &newsCreateUsecase{
				newsRepo:      newsRepo,
				publisherRepo: publisherRepo,
				worker:        worker,
				hasher:        hasher,
			}

			hashedURL = []byte("hashed-subdomain-url")
		)

		// Mock validation steps
		hasher.EXPECT().Hash(req.URL).Return(hashedURL)
		newsRepo.EXPECT().Count(ctx, CustomMatcher(specMatcher(goquSpec.NewNewsByHashedURL(hashedURL)))).Return(int64(0), nil)
		// subdomain.example.com should resolve to example.com as the effective TLD+1
		publisherRepo.EXPECT().Get(ctx, CustomMatcher(specMatcher(goquSpec.NewPublisherByDomain("example.com")))).Return(entity.Publisher{}, appErrors.ErrNotFound)

		// Mock publisher creation for the main domain
		publisherRepo.EXPECT().Create(ctx, &entity.Publisher{
			Name:   "example.com",
			Domain: "example.com",
		}).DoAndReturn(func(ctx context.Context, publisher *entity.Publisher) error {
			publisher.ID = 3
			return nil
		})

		// Mock successful news create
		newsRepo.EXPECT().Create(ctx, &entity.News{
			URL:         req.URL,
			HashedURL:   hashedURL,
			Status:      entity.NewsStatusAdded,
			Category:    req.Category,
			PublisherID: 3,
		}).DoAndReturn(func(ctx context.Context, news *entity.News) error {
			news.ID = 1
			return nil
		})

		// Mock successful job enqueue
		expectedJob := entity.WebScrapperJob{
			Domain: entity.WebDomain("example.com"), // Should use main domain, not subdomain
			URL:    req.URL,
			NewsID: 1,
		}
		expectedJobBytes, _ := json.Marshal(expectedJob)
		expectedMessage := queue.Message{
			Type: queue.MessageTypeWebScrapper,
			Body: expectedJobBytes,
		}
		expectedMessageBytes, _ := json.Marshal(expectedMessage)
		worker.EXPECT().Run(expectedMessageBytes).Return(nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.URL, result.URL)
		assert.Equal(t, entity.NewsStatusAdded.String(), result.Status)
	})
}
