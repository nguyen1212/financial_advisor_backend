package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	goquSpec "github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewNewsSearchUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsRepo := mock.NewMockNewsRepository(mockCtrl)
	newsWithFullTextRepo := mock.NewMockNewsWithFullTextRepository(mockCtrl)

	expected := &newsSearchUsecase{
		newsRepo:             newsRepo,
		newsWithFullTextRepo: newsWithFullTextRepo,
	}

	assert.Equal(t, expected, NewNewsSearchUsecase(newsRepo, newsWithFullTextRepo))
}

func Test_NewsSearchUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success with few keywords (proximity search)", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchRequest{
				Keywords: []string{"finance", "market"},
				Paging: dto.PagingRequest{
					Size: 10,
					Page: 1,
				},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			newsWithFullTextEntities = []entity.NewsWithFullText{
				{
					ID:      11,
					NewsID:  "1",
					Content: "Financial <mark>market</mark> analysis shows growth",
				},
				{
					ID:      22,
					NewsID:  "2",
					Content: "Stock <mark>market</mark> update for <mark>finance</mark> sector",
				},
			}

			newsEntities = []entity.News{
				{
					ID:                 1,
					NewsWithFullTextID: 11,
					Title:              "Market Analysis",
					Category:           entity.NewsCategoryFinance,
					Status:             entity.NewsStatusSynced,
				},
				{
					ID:                 2,
					NewsWithFullTextID: 22,
					Title:              "Stock Update",
					Category:           entity.NewsCategoryFinance,
					Status:             entity.NewsStatusSynced,
				},
			}
		)

		// Mock full-text search with proximity
		newsWithFullTextRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsWithFullTextByKeywords(
				req.Keywords,
				256,
				goquSpec.FullTextSearchOpProximity,
			))),
			goquSpec.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(newsWithFullTextEntities, nil)

		// Mock news repository find by IDs
		newsRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByIDs([]uint64{1, 2}))),
			goquSpec.ToPaging(2, 1),
		).Return(newsEntities, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Equal(t, "Market Analysis", result[0].Title)
		assert.Equal(t, "Financial <mark>market</mark> analysis shows growth", result[0].Content)
		assert.Equal(t, uint64(2), result[1].ID)
		assert.Equal(t, "Stock Update", result[1].Title)
		assert.Equal(t, "Stock <mark>market</mark> update for <mark>finance</mark> sector", result[1].Content)
	})

	t.Run("success with many keywords (quorum search)", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchRequest{
				Keywords: []string{"word1", "word2", "word3", "word4", "word5", "word6", "word7", "word8", "word9", "word10", "word11"},
				Paging: dto.PagingRequest{
					Size: 5,
					Page: 1,
				},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			newsWithFullTextEntities = []entity.NewsWithFullText{
				{
					NewsID:  "1",
					Content: "Content with multiple matching words",
				},
			}

			newsEntities = []entity.News{
				{
					ID:       1,
					Title:    "Multi-keyword Match",
					Category: entity.NewsCategoryFinance,
					Status:   entity.NewsStatusSynced,
				},
			}
		)

		// Mock full-text search with quorum (>10 keywords)
		newsWithFullTextRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsWithFullTextByKeywords(
				req.Keywords,
				256,
				goquSpec.FullTextSearchOpQuorum,
			))),
			goquSpec.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(newsWithFullTextEntities, nil)

		// Mock news repository find by IDs
		newsRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByIDs([]uint64{1}))),
			goquSpec.ToPaging(1, 1),
		).Return(newsEntities, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Equal(t, "Multi-keyword Match", result[0].Title)
		assert.Equal(t, "Content with multiple matching words", result[0].Content)
	})

	t.Run("empty keywords", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchRequest{
				Keywords: []string{},
				Paging: dto.PagingRequest{
					Size: 10,
					Page: 1,
				},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}
		)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("whitespace only keywords", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchRequest{
				Keywords: []string{"  ", "\t", "\n"},
				Paging: dto.PagingRequest{
					Size: 10,
					Page: 1,
				},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}
		)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("full text search error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchRequest{
				Keywords: []string{"test", "search"},
				Paging: dto.PagingRequest{
					Size: 10,
					Page: 1,
				},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			fullTextErr = errors.New("full text search error")
		)

		// Mock full-text search error
		newsWithFullTextRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsWithFullTextByKeywords(
				req.Keywords,
				256,
				goquSpec.FullTextSearchOpProximity,
			))),
			goquSpec.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(nil, fullTextErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find news with full text by file IDs")
	})

	t.Run("news repository error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchRequest{
				Keywords: []string{"test"},
				Paging: dto.PagingRequest{
					Size: 10,
					Page: 1,
				},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			newsWithFullTextEntities = []entity.NewsWithFullText{
				{
					NewsID:  "1",
					Content: "Test content",
				},
			}

			newsRepoErr = errors.New("news repository error")
		)

		// Mock successful full-text search
		newsWithFullTextRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsWithFullTextByKeywords(
				req.Keywords,
				256,
				goquSpec.FullTextSearchOpProximity,
			))),
			goquSpec.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(newsWithFullTextEntities, nil)

		// Mock news repository error
		newsRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByIDs([]uint64{1}))),
			goquSpec.ToPaging(1, 1),
		).Return(nil, newsRepoErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find news by date range")
	})

	t.Run("no full text results", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchRequest{
				Keywords: []string{"nonexistent"},
				Paging: dto.PagingRequest{
					Size: 10,
					Page: 1,
				},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}
		)

		// Mock empty full-text search results
		newsWithFullTextRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsWithFullTextByKeywords(
				req.Keywords,
				256,
				goquSpec.FullTextSearchOpProximity,
			))),
			goquSpec.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return([]entity.NewsWithFullText{}, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("success with mismatched news content", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchRequest{
				Keywords: []string{"test"},
				Paging: dto.PagingRequest{
					Size: 10,
					Page: 1,
				},
			}

			newsRepo             = mock.NewMockNewsRepository(mockCtrl)
			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchUsecase{
				newsRepo:             newsRepo,
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			newsWithFullTextEntities = []entity.NewsWithFullText{
				{
					ID:      11,
					NewsID:  "1",
					Content: "Test content 1",
				},
				{
					ID:      22,
					NewsID:  "2",
					Content: "Test content 2",
				},
			}

			// News repository returns different IDs (simulating data inconsistency)
			newsEntities = []entity.News{
				{
					ID:                 1,
					NewsWithFullTextID: 11,
					Title:              "News 1",
					Category:           entity.NewsCategoryFinance,
					Status:             entity.NewsStatusSynced,
				},
				{
					ID:                 3,
					NewsWithFullTextID: 33, // Different ID than full-text result
					Title:              "News 3",
					Category:           entity.NewsCategoryMilitary,
					Status:             entity.NewsStatusSynced,
				},
			}
		)

		// Mock full-text search
		newsWithFullTextRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsWithFullTextByKeywords(
				req.Keywords,
				256,
				goquSpec.FullTextSearchOpProximity,
			))),
			goquSpec.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(newsWithFullTextEntities, nil)

		// Mock news repository find by IDs
		newsRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsByIDs([]uint64{1, 2}))),
			goquSpec.ToPaging(2, 1),
		).Return(newsEntities, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		// First news item has matching content
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Equal(t, "News 1", result[0].Title)
		assert.Equal(t, "Test content 1", result[0].Content)
		// Second news item doesn't have matching content (ID 3 vs 2)
		assert.Equal(t, uint64(3), result[1].ID)
		assert.Equal(t, "News 3", result[1].Title)
		assert.Equal(t, "", result[1].Content) // No matching content
	})
}

