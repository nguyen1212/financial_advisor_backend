package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/repository/mock"
	goquSpec "github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewNewsSearchSuggestionsGetUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	newsWithFullTextRepo := mock.NewMockNewsWithFullTextRepository(mockCtrl)

	expected := &newsSearchSuggestionsGetUsecase{
		newsWithFullTextRepo: newsWithFullTextRepo,
	}

	assert.Equal(t, expected, NewNewsSearchSuggestionsGetUsecase(newsWithFullTextRepo))
}

func Test_NewsSearchSuggestionsGetUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success with keywords", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchSuggestionsRequest{
				Keywords: []string{"finan", "market"},
			}

			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchSuggestionsGetUsecase{
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			expectedSuggestions = []string{
				"finance",
				"financial",
				"financing",
				"market",
				"marketing",
				"marketplace",
			}
		)

		// Mock search suggestions call
		newsWithFullTextRepo.EXPECT().FindSearchSuggestions(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsSearchSuggestions(
				req.Keywords,
				goquSpec.StrongFuzziness,
			))),
			nil,
		).Return(expectedSuggestions, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedSuggestions, result)
		assert.Len(t, result, 6)
		assert.Contains(t, result, "finance")
		assert.Contains(t, result, "market")
	})

	t.Run("success with single keyword", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchSuggestionsRequest{
				Keywords: []string{"tech"},
			}

			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchSuggestionsGetUsecase{
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			expectedSuggestions = []string{
				"technology",
				"technical",
				"technique",
			}
		)

		// Mock search suggestions call
		newsWithFullTextRepo.EXPECT().FindSearchSuggestions(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsSearchSuggestions(
				req.Keywords,
				goquSpec.StrongFuzziness,
			))),
			nil,
		).Return(expectedSuggestions, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedSuggestions, result)
		assert.Len(t, result, 3)
	})

	t.Run("empty keywords", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchSuggestionsRequest{
				Keywords: []string{},
			}

			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchSuggestionsGetUsecase{
				newsWithFullTextRepo: newsWithFullTextRepo,
			}
		)

		// No mock expectations since method should return early

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("nil keywords", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchSuggestionsRequest{
				Keywords: nil,
			}

			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchSuggestionsGetUsecase{
				newsWithFullTextRepo: newsWithFullTextRepo,
			}
		)

		// No mock expectations since method should return early

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchSuggestionsRequest{
				Keywords: []string{"search", "term"},
			}

			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchSuggestionsGetUsecase{
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			repoErr = errors.New("full text search repository error")
		)

		// Mock repository error
		newsWithFullTextRepo.EXPECT().FindSearchSuggestions(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsSearchSuggestions(
				req.Keywords,
				goquSpec.StrongFuzziness,
			))),
			nil,
		).Return(nil, repoErr)

		result, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get news search suggestions")
	})

	t.Run("no suggestions found", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchSuggestionsRequest{
				Keywords: []string{"nonexistentterm"},
			}

			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchSuggestionsGetUsecase{
				newsWithFullTextRepo: newsWithFullTextRepo,
			}
		)

		// Mock empty suggestions result
		newsWithFullTextRepo.EXPECT().FindSearchSuggestions(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsSearchSuggestions(
				req.Keywords,
				goquSpec.StrongFuzziness,
			))),
			nil,
		).Return([]string{}, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("success with many keywords", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.NewsSearchSuggestionsRequest{
				Keywords: []string{"fin", "mark", "econ", "stock", "trade"},
			}

			newsWithFullTextRepo = mock.NewMockNewsWithFullTextRepository(mockCtrl)
			uc                   = &newsSearchSuggestionsGetUsecase{
				newsWithFullTextRepo: newsWithFullTextRepo,
			}

			expectedSuggestions = []string{
				"finance",
				"financial",
				"market",
				"economy",
				"economic",
				"stock",
				"stocks",
				"trading",
				"trader",
			}
		)

		// Mock search suggestions call
		newsWithFullTextRepo.EXPECT().FindSearchSuggestions(
			ctx,
			CustomMatcher(specMatcher(goquSpec.NewNewsSearchSuggestions(
				req.Keywords,
				goquSpec.StrongFuzziness,
			))),
			nil,
		).Return(expectedSuggestions, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedSuggestions, result)
		assert.Len(t, result, 9)
		assert.Contains(t, result, "finance")
		assert.Contains(t, result, "market")
		assert.Contains(t, result, "economy")
		assert.Contains(t, result, "stock")
		assert.Contains(t, result, "trading")
	})
}

