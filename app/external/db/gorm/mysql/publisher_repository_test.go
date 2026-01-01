package mysql

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications/mock"
	appErrors "github.com/financial_advisor/app/errors"
	gormAdapter "github.com/financial_advisor/app/external/db/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_NewPublisherRepository(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		&publisherRepository{},
		NewPublisherRepository(&gormAdapter.MySQL{}),
	)
}

func Test_publisherRepository_Count(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)

			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:       sqlDB,
				DriverName: "mysql",
				// Skip initializing with version to avoid querying the database for version
				// SELECT VERSION()
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)

			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx  = context.Background()
			r    = &publisherRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)
		)

		spec.EXPECT().ToCount().Return("SELECT COUNT(*) FROM `publishers`", nil)

		rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(5)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `publishers`")).WillReturnRows(rows)

		count, err := r.Count(ctx, spec)

		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("failed - count publishers", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)

			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:       sqlDB,
				DriverName: "mysql",
				// Skip initializing with version to avoid querying the database for version
				// SELECT VERSION()
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)

			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx  = context.Background()
			r    = &publisherRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)

			wantErr = errors.New("database error")
		)

		spec.EXPECT().ToCount().Return("SELECT COUNT(*) FROM `publishers`", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `publishers`")).WillReturnError(wantErr)

		_, err = r.Count(ctx, spec)

		assert.Equal(t, wantErr, err)
	})

	t.Run("failed - spec error", func(t *testing.T) {
		t.Parallel()

		sqlDB, _, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)

			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:       sqlDB,
				DriverName: "mysql",
				// Skip initializing with version to avoid querying the database for version
				// SELECT VERSION()
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)

			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx  = context.Background()
			r    = &publisherRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)

			sqlErr  = errors.New("database error")
			wantErr = fmt.Errorf("convert to raw count sql: %w", sqlErr)
		)

		spec.EXPECT().ToCount().Return("", sqlErr)

		_, err = r.Count(ctx, spec)

		assert.Equal(t, wantErr, err)
	})
}

func Test_publisherRepository_Find(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx    = context.Background()
			r      = &publisherRepository{gormDB}
			spec   = mock.NewMockI(mockCtrl)
			paging = mock.NewMockPagingI(mockCtrl)
		)

		spec.EXPECT().ToFind(paging).Return("SELECT * FROM `publishers` LIMIT 10", nil)

		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "name", "description", "domain", "created_at", "updated_at"}).
			AddRow(1, "VnExpress", "Leading news portal in Vietnam", "vnexpress.net", now, now).
			AddRow(2, "Tuoi Tre", "Youth newspaper", "tuoitre.vn", now, now)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `publishers` LIMIT 10")).WillReturnRows(rows)

		publisherList, err := r.Find(ctx, spec, paging)

		require.NoError(t, err)
		assert.Len(t, publisherList, 2)

		// Validate first publisher
		assert.Equal(t, uint64(1), publisherList[0].ID)
		assert.Equal(t, "VnExpress", publisherList[0].Name)
		assert.Equal(t, "Leading news portal in Vietnam", publisherList[0].Description)
		assert.Equal(t, "vnexpress.net", publisherList[0].Domain)

		// Validate second publisher
		assert.Equal(t, uint64(2), publisherList[1].ID)
		assert.Equal(t, "Tuoi Tre", publisherList[1].Name)
		assert.Equal(t, "Youth newspaper", publisherList[1].Description)
		assert.Equal(t, "tuoitre.vn", publisherList[1].Domain)
	})

	t.Run("failed - database error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx     = context.Background()
			r       = &publisherRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			paging  = mock.NewMockPagingI(mockCtrl)
			wantErr = errors.New("database error")
		)

		spec.EXPECT().ToFind(paging).Return("SELECT * FROM `publishers` LIMIT 10", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `publishers` LIMIT 10")).WillReturnError(wantErr)

		_, err = r.Find(ctx, spec, paging)

		assert.Equal(t, wantErr, err)
	})

	t.Run("failed - spec error", func(t *testing.T) {
		t.Parallel()

		sqlDB, _, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx     = context.Background()
			r       = &publisherRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			paging  = mock.NewMockPagingI(mockCtrl)
			wantErr = errors.New("spec error")
		)

		spec.EXPECT().ToFind(paging).Return("", wantErr)

		_, err = r.Find(ctx, spec, paging)

		assert.ErrorContains(t, err, "convert to raw find sql")
	})
}

func Test_publisherRepository_Get(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx  = context.Background()
			r    = &publisherRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `publishers` WHERE id = 1", nil)

		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "name", "description", "domain", "created_at", "updated_at"}).
			AddRow(1, "VnExpress", "Leading news portal in Vietnam", "vnexpress.net", now, now)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `publishers` WHERE id = 1")).WillReturnRows(rows)

		publisher, err := r.Get(ctx, spec)

		require.NoError(t, err)
		assert.Equal(t, uint64(1), publisher.ID)
		assert.Equal(t, "VnExpress", publisher.Name)
		assert.Equal(t, "Leading news portal in Vietnam", publisher.Description)
		assert.Equal(t, "vnexpress.net", publisher.Domain)
	})

	t.Run("failed - record not found", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx  = context.Background()
			r    = &publisherRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `publishers` WHERE id = 999", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `publishers` WHERE id = 999")).WillReturnError(gorm.ErrRecordNotFound)

		_, err = r.Get(ctx, spec)

		assert.ErrorIs(t, err, appErrors.ErrNotFound)
	})

	t.Run("failed - database error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx     = context.Background()
			r       = &publisherRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			wantErr = errors.New("database error")
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `publishers` WHERE id = 1", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `publishers` WHERE id = 1")).WillReturnError(wantErr)

		_, err = r.Get(ctx, spec)

		assert.Equal(t, wantErr, err)
	})

	t.Run("failed - spec error", func(t *testing.T) {
		t.Parallel()

		sqlDB, _, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			ctx     = context.Background()
			r       = &publisherRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			wantErr = errors.New("spec error")
		)

		spec.EXPECT().ToGet().Return("", wantErr)

		_, err = r.Get(ctx, spec)

		assert.ErrorContains(t, err, "convert to raw get sql")
	})
}

func Test_publisherRepository_Create(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		var (
			ctx       = context.Background()
			r         = &publisherRepository{gormDB}
			publisher = &entity.Publisher{
				Name:        "VnExpress",
				Description: "Leading news portal in Vietnam",
				Domain:      "vnexpress.net",
			}
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `publishers` (`name`,`description`,`domain`) VALUES (?,?,?)",
		)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.ExpectCommit()

		err = r.Create(ctx, publisher)

		require.NoError(t, err)
	})

	t.Run("failed - database error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mockDB, err := sqlmock.New()
		if err != nil {
			t.Errorf("failed to open sqlmock database: %s", err)
			return
		}
		defer sqlDB.Close()

		gormDB, err := gorm.Open(
			mysql.New(mysql.Config{
				Conn:                      sqlDB,
				DriverName:                "mysql",
				SkipInitializeWithVersion: true,
			}),
			&gorm.Config{},
		)
		if err != nil {
			t.Errorf("failed to open gorm db: %s", err)
			return
		}

		var (
			ctx       = context.Background()
			r         = &publisherRepository{gormDB}
			publisher = &entity.Publisher{
				Name:        "VnExpress",
				Description: "Leading news portal in Vietnam",
				Domain:      "vnexpress.net",
			}
			wantErr = errors.New("database error")
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `publishers` (`name`,`description`,`domain`) VALUES (?,?,?)",
		)).
			WillReturnError(wantErr)
		mockDB.ExpectRollback()

		err = r.Create(ctx, publisher)

		assert.Equal(t, wantErr, err)
	})
}
