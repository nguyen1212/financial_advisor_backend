package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/goqu/specifications"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewPublisherGetUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	publisherRepo := mock.NewMockPublisherRepository(mockCtrl)

	assert.Equal(t, &publisherGetUsecase{publisherRepo: publisherRepo}, NewPublisherGetUsecase(publisherRepo))
}

func Test_PublisherGetUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx                = context.Background()
			publisherID uint64 = 1

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publisherGetUsecase{publisherRepo: publisherRepo}

			publisher = entity.Publisher{
				ID:          publisherID,
				Name:        "Test Publisher",
				Domain:      "example.com",
				Description: "A test publisher",
			}
		)

		publisherRepo.EXPECT().Get(
			ctx,
			CustomMatcher(SpecMatcher(specifications.NewPublisherByID(publisherID))),
		).Return(publisher, nil)

		result, err := uc.Execute(ctx, publisherID)

		assert.NoError(t, err)
		assert.Equal(t, publisherID, result.ID)
		assert.Equal(t, "Test Publisher", result.Name)
		assert.Equal(t, "example.com", result.Domain)
		assert.Equal(t, "A test publisher", result.Description)
	})

	t.Run("publisher not found", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx                = context.Background()
			publisherID uint64 = 1

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publisherGetUsecase{publisherRepo: publisherRepo}

			wantErr = appErrors.NewErrorNotFound(
				appErrors.ErrorCodePublisherNotFound,
				"publisher not found",
			)
		)

		publisherRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(specifications.NewPublisherByID(publisherID)))).Return(entity.Publisher{}, appErrors.ErrNotFound)

		_, err := uc.Execute(ctx, publisherID)

		assert.Equal(t, wantErr, err)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx                = context.Background()
			publisherID uint64 = 1

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publisherGetUsecase{publisherRepo: publisherRepo}

			repoErr = errors.New("database error")
		)

		publisherRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(specifications.NewPublisherByID(publisherID)))).Return(entity.Publisher{}, repoErr)

		_, err := uc.Execute(ctx, publisherID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "find publisher by date range")
	})

	t.Run("success with empty description", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx                = context.Background()
			publisherID uint64 = 2

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publisherGetUsecase{publisherRepo: publisherRepo}

			publisher = entity.Publisher{
				ID:          publisherID,
				Name:        "Simple Publisher",
				Domain:      "simple.com",
				Description: "",
			}
		)

		publisherRepo.EXPECT().Get(ctx, CustomMatcher(SpecMatcher(specifications.NewPublisherByID(publisherID)))).Return(publisher, nil)

		result, err := uc.Execute(ctx, publisherID)

		assert.NoError(t, err)
		assert.Equal(t, publisherID, result.ID)
		assert.Equal(t, "Simple Publisher", result.Name)
		assert.Equal(t, "simple.com", result.Domain)
		assert.Equal(t, "", result.Description)
	})
}
