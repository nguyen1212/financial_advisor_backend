package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/financial_advisor/app/services/hasher"
	"github.com/financial_advisor/app/services/queue"
	"github.com/financial_advisor/app/usecases/dto"
	"golang.org/x/net/publicsuffix"
)

type NewsCreateUsecase interface {
	Execute(
		ctx context.Context,
		req dto.NewsCreateRequest,
	) (dto.News, error)
}

type newsCreateUsecase struct {
	newsRepo      repository.NewsRepository
	publisherRepo repository.PublisherRepository

	jobQueue queue.I
	hasher   hasher.I
}

func NewNewsCreateUsecase(
	newsRepo repository.NewsRepository,
	publisherRepo repository.PublisherRepository,
	jobQueue queue.I,
	hasher hasher.I,
) NewsCreateUsecase {
	return &newsCreateUsecase{
		newsRepo:      newsRepo,
		jobQueue:      jobQueue,
		publisherRepo: publisherRepo,
		hasher:        hasher,
	}
}

func (uc *newsCreateUsecase) Execute(
	ctx context.Context,
	req dto.NewsCreateRequest,
) (dto.News, error) {
	publisher, hashedURL, err := uc.validate(ctx, req)
	if err != nil {
		return dto.News{}, err
	}

	news := entity.News{
		URL:         req.URL,
		HashedURL:   hashedURL,
		Status:      entity.NewsStatusAdded,
		Category:    req.Category,
		PublisherID: publisher.ID,
	}

	if err = uc.newsRepo.Create(
		ctx,
		&news,
	); err != nil {
		// there must be concurrent requests creating the same news
		if errors.Is(err, appErrors.ErrConflicted) {
			return dto.News{}, appErrors.NewErrorConflicted(
				appErrors.ErrorCodeConflicted,
				"news with the same URL already exists",
			)
		}

		return dto.News{}, fmt.Errorf("create news: %w", err)
	}

	webScrapperJob := entity.WebScrapperJob{
		Domain: entity.WebDomain(publisher.Domain),
		URL:    news.URL,
		NewsID: news.ID,
	}

	jobEncoded, err := json.Marshal(webScrapperJob)
	if err != nil {
		return dto.News{}, fmt.Errorf("marshal web scrapper job: %w", err)
	}

	if err := uc.jobQueue.Enqueue(queue.Message{
		Type: queue.MessageTypeWebScrapper,
		Body: jobEncoded,
	}); err != nil {
		return dto.News{}, fmt.Errorf("enqueue web scrapper job: %w", err)
	}

	return dto.ToDtoNews(news), nil
}

func (uc *newsCreateUsecase) validate(
	ctx context.Context,
	req dto.NewsCreateRequest,
) (entity.Publisher, []byte, error) {
	if len(req.URL) > 500 {
		return entity.Publisher{}, nil, appErrors.NewErrorBadRequest(
			appErrors.ErrorCodeURLTooLong,
			"url length exceeds the limit",
		)
	}

	hashedURL := uc.hasher.Hash(req.URL)

	countURL, err := uc.newsRepo.Count(
		ctx,
		specifications.NewNewsByHashedURL(hashedURL),
	)
	if err != nil {
		return entity.Publisher{}, nil, fmt.Errorf("count url by hashed value: %w", err)
	}

	// As the request may contain different information about the news
	// so we cannot return 200 but 409 to indicates there is existing site.
	if countURL > 0 {
		return entity.Publisher{}, nil, appErrors.NewErrorConflicted(
			appErrors.ErrorCodeConflicted,
			"news with the same URL already exists",
		)
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil || parsedURL == nil {
		return entity.Publisher{}, nil, appErrors.NewErrorBadRequest(
			appErrors.ErrorCodeURLInvalid,
			"invalid URL format",
		)
	}

	domain, err := publicsuffix.EffectiveTLDPlusOne(parsedURL.Hostname())
	if err != nil {
		return entity.Publisher{}, nil, appErrors.NewErrorBadRequest(
			appErrors.ErrorCodeURLInvalid,
			"invalid URL domain",
		)
	}

	publisher, err := uc.publisherRepo.Get(
		ctx,
		specifications.NewPublisherByDomain(domain),
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return entity.Publisher{}, nil, appErrors.NewErrorBadRequest(
				appErrors.ErrorCodePublisherNotFound,
				"publisher not found",
			)
		}

		return entity.Publisher{}, nil, fmt.Errorf("count publisher by id: %w", err)

	}

	return publisher, hashedURL, nil
}
