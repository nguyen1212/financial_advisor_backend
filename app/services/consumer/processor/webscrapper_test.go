package processor

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/usecases"
	"github.com/financial_advisor/app/usecases/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_NewWebScrapperProcessor(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	webScrapperUcMock := mock.NewMockWebScrapperUsecase(mockCtrl)
	fallbackUcMock := mock.NewMockFallbackScrapperUsecase(mockCtrl)
	mUc := map[entity.WebDomain]usecases.WebScrapperUsecase{
		entity.WebDomainVnExpress: webScrapperUcMock,
	}

	assert.Equal(
		t,
		WebScrapperProcessor{mUsecases: mUc, fallback: fallbackUcMock},
		NewWebScrapperProcessor(mUc, fallbackUcMock),
	)
}

func Test_WebScrapperProcessor_Execute(t *testing.T) {
	t.Run("failed - unmarshal msg", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			msg       = []byte(`invalid-json`)
			processor = NewWebScrapperProcessor(nil, nil)
			wantErr   = errors.New("invalid character 'i' looking for beginning of value")
		)

		err := processor.Execute(msg)

		assert.EqualError(t, wantErr, err.Error())
	})

	t.Run("failed - not found usecases for domain, fallback to default and fail", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			msg        = []byte(`{"domain":"unknown-domain","url":"http://example.com"}`)
			fallbackUc = mock.NewMockFallbackScrapperUsecase(mockCtrl)
			processor  = NewWebScrapperProcessor(nil, fallbackUc)
			wantErr    = errors.New("usecase for domain not found")
		)

		fallbackUc.EXPECT().Execute(context.Background(), entity.WebScrapperJob{
			Domain: "unknown-domain",
			URL:    "http://example.com",
		}).Return(wantErr)

		err := processor.Execute(msg)

		assert.EqualError(t, wantErr, err.Error())
	})

	t.Run("failed - usecase execute job failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			msg = []byte(`{"domain":"vnexpress.net","url":"http://example.com"}`)
			uc  = mock.NewMockWebScrapperUsecase(mockCtrl)
			job = entity.WebScrapperJob{
				Domain: entity.WebDomainVnExpress,
				URL:    "http://example.com",
			}
			processor = NewWebScrapperProcessor(map[entity.WebDomain]usecases.WebScrapperUsecase{
				entity.WebDomainVnExpress: uc,
			}, nil)
			wantErr = errors.New("failed to scrape web page")
		)

		uc.EXPECT().Execute(context.Background(), job).Return(wantErr)

		err := processor.Execute(msg)

		assert.EqualError(t, wantErr, err.Error())
	})

	t.Run("success - usecase execute job", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			msg = []byte(`{"domain":"vnexpress.net","url":"http://example.com"}`)
			uc  = mock.NewMockWebScrapperUsecase(mockCtrl)
			job = entity.WebScrapperJob{
				Domain: entity.WebDomainVnExpress,
				URL:    "http://example.com",
			}
			processor = NewWebScrapperProcessor(map[entity.WebDomain]usecases.WebScrapperUsecase{
				entity.WebDomainVnExpress: uc,
			}, nil)
		)

		uc.EXPECT().Execute(context.Background(), job).Return(nil)

		err := processor.Execute(msg)

		require.NoError(t, err)
	})
}
