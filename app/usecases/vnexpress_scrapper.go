package usecases

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type vnExpressScrapper struct {
	newsRepo repository.NewsRepository
}

type WebScrapperUsecase interface {
	Execute(context.Context, entity.WebScrapperJob) error
}

func NewVnExpressScrapperUsecase(
	newsRepo repository.NewsRepository,
) WebScrapperUsecase {
	return vnExpressScrapper{
		newsRepo: newsRepo,
	}
}

func (uc vnExpressScrapper) Execute(
	ctx context.Context,
	job entity.WebScrapperJob,
) error {
	// instantiate a new collector object
	c := colly.NewCollector(
		colly.AllowedDomains(string(job.Domain)),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 4 request at a time
		Parallelism: 1,
	})

	// set a global User Agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

	// set up the proxy
	// err := c.SetProxy("http://35.185.196.38:3128")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode > 299 {
			if err := uc.errorHandler(ctx, job.NewsID); err != nil {
				logrus.WithField("news_id", job.NewsID).
					WithField("url", job.URL).
					Errorf("handle error status code: %v", err)
			}
		}
	})

	var (
		content       = &strings.Builder{}
		author        string
		publishedDate time.Time
		title         string
		thumbnailURL  string
	)

	content.Grow(1000)

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

	// register all pages to scrape
	c.Visit(job.URL)

	// store the data to a CSV after extraction
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
		}
	})

	// wait for Colly to visit all pages
	c.Wait()

	return nil
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
	news, err := uc.newsRepo.Get(ctx, specifications.NewNewsByID(newsID, "Publisher"))
	if err != nil {
		return fmt.Errorf("get news by id to save extracted content: %w", err)
	}

	// handle error to update news status to failed
	defer func() {
		if err != nil {
			news.Status = entity.NewsStatusFailed
			if updateErr := uc.newsRepo.Update(ctx, &news); updateErr != nil {
				logrus.WithField("news_id", newsID).
					Errorf("update news status to failed after save file error: %v", updateErr)
			}
		}
	}()

	filePath := news.StoragePath()

	if err = os.MkdirAll(news.StorageDir(), 0755); err != nil {
		return fmt.Errorf("create directories for news storage: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file to save extracted content: %w", err)
	}

	defer file.Close()

	_, err = file.WriteString(content.String())
	if err != nil {
		return fmt.Errorf("write extracted content to file: %w", err)
	}

	fileStat, err := file.Stat()
	if err == nil {
		news.FileSize = fileStat.Size()
	} else {
		news.FileSize = int64(len(content.String()))
	}

	news.FilePath = filePath
	news.Author = author
	news.PublishedAt = &publishedDate
	news.Status = entity.NewsStatusSynced
	news.Title = title
	news.Thumbnail = thumbnailURL

	if err = uc.newsRepo.Update(ctx, &news); err != nil {
		return fmt.Errorf("update news after saving extracted content: %w", err)
	}

	return nil
}

func (uc vnExpressScrapper) errorHandler(ctx context.Context, newsID uint64) error {
	news, err := uc.newsRepo.Get(context.Background(), specifications.NewNewsByID(newsID))
	if err != nil {
		return fmt.Errorf("get news by id to update error status: %w", err)
	}

	news.Status = entity.NewsStatusFailed

	if err = uc.newsRepo.Update(ctx, &news); err != nil {
		return fmt.Errorf("update news status to failed: %w", err)
	}

	return nil
}
