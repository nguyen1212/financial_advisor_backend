package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/mock"
	appErrors "github.com/financial_advisor/app/errors"
	"github.com/financial_advisor/app/external/db/gorm/specifications"
	"github.com/financial_advisor/app/usecases/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewPublisherCreateUsecase(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	publisherRepo := mock.NewMockPublisherRepository(mockCtrl)

	assert.Equal(t, &publisherCreateUsecase{publisherRepo: publisherRepo}, NewPublisherCreateUsecase(publisherRepo))
}

func Test_PublisherCreateUsecase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.PublisherCreateRequest{
				Name:        "Test Publisher",
				Domain:      "example.com",
				Description: "A test publisher",
			}

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publisherCreateUsecase{publisherRepo: publisherRepo}
		)

		// Mock count returns 0 (no existing publisher)
		publisherRepo.EXPECT().Count(
			ctx,
			specifications.NewPublisherByDomain("example.com"),
		).Return(int64(0), nil)

		// Mock successful create
		publisherRepo.EXPECT().Create(
			ctx,
			&entity.Publisher{
				Name:        req.Name,
				Domain:      "example.com", // Normalized domain
				Description: req.Description,
			},
		).DoAndReturn(func(ctx context.Context, publisher *entity.Publisher) error {
			publisher.ID = 1 // Simulate database assigning ID
			return nil
		})

		result, err := uc.Execute(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), result.ID)
		assert.Equal(t, "Test Publisher", result.Name)
		assert.Equal(t, "example.com", result.Domain)
		assert.Equal(t, "A test publisher", result.Description)
	})

	t.Run("publisher already exists", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.PublisherCreateRequest{
				Name:        "Test Publisher",
				Domain:      "example.com",
				Description: "A test publisher",
			}
			forbiddenErr appErrors.SystemError

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publisherCreateUsecase{publisherRepo: publisherRepo}
		)

		// Mock count returns 1 (existing publisher)
		publisherRepo.EXPECT().Count(
			ctx,
			specifications.NewPublisherByDomain("example.com"),
		).Return(int64(1), nil)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.True(t, errors.As(err, &forbiddenErr))
		assert.Equal(t, appErrors.ErrorTypeConflicted, forbiddenErr.Type())
	})

	t.Run("count error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.PublisherCreateRequest{
				Name:        "Test Publisher",
				Domain:      "example.com",
				Description: "A test publisher",
			}

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publisherCreateUsecase{publisherRepo: publisherRepo}

			countErr = errors.New("database error")
		)

		publisherRepo.EXPECT().Count(
			ctx,
			specifications.NewPublisherByDomain("example.com"),
		).Return(int64(0), countErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "count existing publisher by domain")
	})

	t.Run("create error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx = context.Background()
			req = dto.PublisherCreateRequest{
				Name:        "Test Publisher",
				Domain:      "example.com",
				Description: "A test publisher",
			}

			publisherRepo = mock.NewMockPublisherRepository(mockCtrl)
			uc            = &publisherCreateUsecase{publisherRepo: publisherRepo}

			createErr = errors.New("database error")
		)

		publisherRepo.EXPECT().Count(
			ctx,
			specifications.NewPublisherByDomain("example.com"),
		).Return(int64(0), nil)

		publisherRepo.EXPECT().Create(
			ctx,
			&entity.Publisher{
				Name:        req.Name,
				Domain:      "example.com", // Normalized domain
				Description: req.Description,
			},
		).Return(createErr)

		_, err := uc.Execute(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create new publisher")
	})
}
