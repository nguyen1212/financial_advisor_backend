package usecases

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type vnExpressScrapper struct {
	newsRepo             repository.NewsRepository
	newsWithFullTextRepo repository.NewsWithFullTextRepository
}

type WebScrapperUsecase interface {
	Execute(context.Context, entity.WebScrapperJob) error
}

func NewVnExpressScrapperUsecase(
	newsRepo repository.NewsRepository,
	newsWithFullTextRepo repository.NewsWithFullTextRepository,
) WebScrapperUsecase {
	return vnExpressScrapper{
		newsRepo:             newsRepo,
		newsWithFullTextRepo: newsWithFullTextRepo,
	}
}

func (uc vnExpressScrapper) Execute(
	ctx context.Context,
	job entity.WebScrapperJob,
) error {
	var globalErr error

	defer func() {
		if globalErr != nil {
			if err := uc.errorHandler(ctx, job.NewsID); err != nil {
				logrus.WithField("news_id", job.NewsID).
					WithField("url", job.URL).
					Errorf("handle error from saving file: %v", err)

				globalErr = errors.Join(globalErr, err)
			}
		}
	}()

	var (
		content       = &strings.Builder{}
		author        string
		publishedDate time.Time
		title         string
		thumbnailURL  string
	)

	// preallocate content buffer to optimize memory usage
	content.Grow(1000)

	// instantiate a new collector object
	c := colly.NewCollector(
		colly.AllowedDomains(string(job.Domain)),
		colly.Async(true),
	)

	c.SetRequestTimeout(60 * time.Second)

	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 4 request at a time
		Parallelism: 1,
	})

	// set a global User Agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

	c.OnError(func(r *colly.Response, err error) {
		globalErr = err
	})

	// extract title
	c.OnHTML("head", func(e *colly.HTMLElement) {
		// itemProp := e.ChildAttr("meta", "itemprop")
		title = e.ChildAttr("meta[itemprop='headline']", "content")
		thumbnailURL = e.ChildAttr("meta[itemprop='thumbnailUrl']", "content")
	})

	// extract content
	c.OnHTML("article", func(e *colly.HTMLElement) {
		// initialize a new Product instance

		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			_, err := content.WriteString(el.Text + "\n")
			if err != nil {
				logrus.WithField("news_id", job.NewsID).
					WithField("url", job.URL).
					Errorf("write content string: %v", err)
			}

			// extract author
			_, exist := el.DOM.Next().Attr("id")
			if exist {
				author = el.Text
			}
		})
	})

	// extract published date
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		if e.Attr("itemprop") == "datePublished" {
			parsedDate, err := time.Parse(time.RFC3339, e.Attr("content"))
			if err != nil {
				logrus.WithField("news_id", job.NewsID).
					WithField("url", job.URL).
					WithField("date_string", e.Attr("content")).
					Errorf("parse published date: %v", err)
			}

			publishedDate = parsedDate
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if err := uc.saveFile(
			ctx,
			job.NewsID,
			title,
			content,
			author,
			publishedDate,
			thumbnailURL,
		); err != nil {
			logrus.WithField("news_id", job.NewsID).
				WithField("url", job.URL).
				Errorf("save extracted content to file: %v", err)

			globalErr = err
		}
	})

	// register all pages to scrape
	c.Visit(job.URL)

	// wait for Colly to visit all pages
	c.Wait()

	return globalErr
}

func (uc vnExpressScrapper) saveFile(
	ctx context.Context,
	newsID uint64,
	title string,
	content *strings.Builder,
	author string,
	publishedDate time.Time,
	thumbnailURL string,
) error {
	news, err := uc.newsRepo.Get(ctx, specifications.NewNewsByID(newsID))
	if err != nil {
		return fmt.Errorf("get news by id to save extracted content: %w", err)
	}

	newsWithFullText := entity.NewsWithFullText{
		NewsID:  strconv.FormatUint(news.ID, 10),
		Title:   news.Title,
		Content: content.String(),
	}

	if err = uc.newsWithFullTextRepo.Create(ctx, &newsWithFullText); err != nil {
		return fmt.Errorf("save news content to full text search index: %w", err)
	}

	news.Title = title
	news.Author = author
	news.PublishedAt = &publishedDate
	news.Status = entity.NewsStatusSynced
	news.Thumbnail = thumbnailURL
	news.FileSize = int64(len(newsWithFullText.Content))
	news.NewsWithFullTextID = newsWithFullText.ID

	if err = uc.newsRepo.Update(ctx, &news); err != nil {
		return fmt.Errorf("update news after saving extracted content: %w", err)
	}

	return nil
}

func (uc vnExpressScrapper) errorHandler(
	_ context.Context,
	newsID uint64,
) error {
	// NOTE: we use a new context.Background() here because the function should be called
	// regardless of the original context state as it serves as a cleanup function
	news, err := uc.newsRepo.Get(
		context.Background(),
		specifications.NewNewsByID(newsID),
	)
	if err != nil {
		return fmt.Errorf("get news by id to update error status: %w", err)
	}

	news.Status = entity.NewsStatusFailed

	if err = uc.newsRepo.Update(
		context.Background(),
		&news,
	); err != nil {
		return fmt.Errorf("update news status to failed: %w", err)
	}

	return nil
}
