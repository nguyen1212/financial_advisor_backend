package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	"github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewPublishersFindUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	publisherRepo := mock.NewMockPublisherRepository(mockCtrl)

	assert.Equal(t, &publishersFindUsecase{publisherRepo: publisherRepo}, NewPublishersFindUsecase(publisherRepo))
}

func Test_PublishersFindUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success with results", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.PublishersFindRequest{
				Paging: dto.PagingRequest{
					Page: 1,
					Size: 10,
				},
			}

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publishersFindUsecase{publisherRepo: publisherRepo}

			publisherEntities = []entity.Publisher{
				{
					ID:          1,
					Name:        "Publisher 1",
					Domain:      "example.com",
					Description: "Description 1",
				},
				{
					ID:          2,
					Name:        "Publisher 2",
					Domain:      "another.com",
					Description: "Description 2",
				},
			}
		)

		publisherRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(specifications.NewPublishersByNone())),
			specifications.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(publisherEntities, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Equal(t, "Publisher 1", result[0].Name)
		assert.Equal(t, "example.com", result[0].Domain)
		assert.Equal(t, uint64(2), result[1].ID)
		assert.Equal(t, "Publisher 2", result[1].Name)
		assert.Equal(t, "another.com", result[1].Domain)
	})

	t.Run("success with empty results", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.PublishersFindRequest{
				Paging: dto.PagingRequest{
					Page: 1,
					Size: 10,
				},
			}

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publishersFindUsecase{publisherRepo: publisherRepo}
		)

		publisherRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(specifications.NewPublishersByNone())),
			specifications.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return([]entity.Publisher{}, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.PublishersFindRequest{
				Paging: dto.PagingRequest{
					Page: 1,
					Size: 10,
				},
			}

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publishersFindUsecase{publisherRepo: publisherRepo}

			repoErr = errors.New("database error")
		)

		publisherRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(specifications.NewPublishersByNone())),
			specifications.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(nil, repoErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find news by date range")
	})

	t.Run("success with different paging", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.PublishersFindRequest{
				Paging: dto.PagingRequest{
					Page: 2,
					Size: 5,
				},
			}

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publishersFindUsecase{publisherRepo: publisherRepo}

			publisherEntities = []entity.Publisher{
				{
					ID:     3,
					Name:   "Publisher 3",
					Domain: "third.com",
				},
			}
		)

		publisherRepo.EXPECT().Find(
			ctx,
			CustomMatcher(specMatcher(specifications.NewPublishersByNone())),
			specifications.ToPaging(req.Paging.Size, req.Paging.Page),
		).Return(publisherEntities, nil)

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, uint64(3), result[0].ID)
		assert.Equal(t, "Publisher 3", result[0].Name)
	})
}
