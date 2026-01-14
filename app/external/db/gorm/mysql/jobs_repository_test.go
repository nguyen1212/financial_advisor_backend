package mysql

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository/specifications/mock"
	appErrors "github.com/financial_advisor/app/errors"
	gormAdapter "github.com/financial_advisor/app/external/db/gorm"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_NewJobsRepository(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		&jobsRepository{},
		NewJobsRepository(&gormAdapter.MySQL{}),
	)
}

func Test_jobsRepository_Get(t *testing.T) {
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
			r    = &jobsRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `jobs` WHERE uuid = 'test-uuid-123'", nil)

		now := time.Now()
		rows := sqlmock.NewRows([]string{"uuid", "payload", "result", "status", "type", "created_at", "updated_at"}).
			AddRow("test-uuid-123", []byte(`{"domain":"vnexpress.net","url":"test.com"}`), []byte(`{"error":""}`), 2, 1, now, now)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE uuid = 'test-uuid-123'")).WillReturnRows(rows)

		job, err := r.Get(ctx, spec)

		require.NoError(t, err)
		assert.Equal(t, "test-uuid-123", job.UUID)
		assert.Equal(t, entity.JobStatusCompleted, job.Status)
		assert.Equal(t, entity.JobTypeWebScrapper, job.Type)
		assert.NotNil(t, job.Payload)
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
			r    = &jobsRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `jobs` WHERE uuid = 'non-existent'", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE uuid = 'non-existent'")).WillReturnError(gorm.ErrRecordNotFound)

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
			r       = &jobsRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			wantErr = errors.New("database error")
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `jobs` WHERE uuid = 'test-uuid'", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE uuid = 'test-uuid'")).WillReturnError(wantErr)

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
			r       = &jobsRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			wantErr = errors.New("spec error")
		)

		spec.EXPECT().ToGet().Return("", wantErr)

		_, err = r.Get(ctx, spec)

		assert.ErrorContains(t, err, "convert to raw get sql")
		assert.ErrorContains(t, err, "spec error")
	})
}

func Test_jobsRepository_Create(t *testing.T) {
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
			ctx = context.Background()
			r   = &jobsRepository{gormDB}
			job = &entity.Job{
				UUID:      "new-job-uuid",
				Payload:   []byte(`{"test":"data"}`),
				ResultEnc: []byte(`{}`),
				Status:    entity.JobStatusNew,
				Type:      entity.JobTypeWebScrapper,
			}
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `jobs` (`uuid`,`payload`,`result`,`status`,`type`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)",
		)).WithArgs(
			"new-job-uuid",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.ExpectCommit()

		err = r.Create(ctx, job)

		require.NoError(t, err)
	})

	t.Run("failed - duplicate key", func(t *testing.T) {
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
			ctx = context.Background()
			r   = &jobsRepository{gormDB}
			job = &entity.Job{
				UUID:      "duplicate-uuid",
				Payload:   []byte(`{"test":"data"}`),
				ResultEnc: []byte(`{}`),
				Status:    entity.JobStatusNew,
				Type:      entity.JobTypeWebScrapper,
			}
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `jobs` (`uuid`,`payload`,`result`,`status`,`type`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)",
		)).WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
		mockDB.ExpectRollback()

		err = r.Create(ctx, job)

		assert.ErrorIs(t, err, appErrors.ErrConflicted)
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
			ctx = context.Background()
			r   = &jobsRepository{gormDB}
			job = &entity.Job{
				UUID:      "test-uuid",
				Payload:   []byte(`{"test":"data"}`),
				ResultEnc: []byte(`{}`),
				Status:    entity.JobStatusNew,
				Type:      entity.JobTypeWebScrapper,
			}
			wantErr = errors.New("database error")
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `jobs` (`uuid`,`payload`,`result`,`status`,`type`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)",
		)).WillReturnError(wantErr)
		mockDB.ExpectRollback()

		err = r.Create(ctx, job)

		assert.Equal(t, wantErr, err)
	})
}

func Test_jobsRepository_Update(t *testing.T) {
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
			ctx = context.Background()
			r   = &jobsRepository{gormDB}
			job = &entity.Job{
				UUID:      "update-uuid",
				Payload:   []byte(`{"updated":"data"}`),
				ResultEnc: []byte(`{"result":"success"}`),
				Status:    entity.JobStatusCompleted,
				Type:      entity.JobTypeWebScrapper,
			}
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"UPDATE `jobs` SET `payload`=?,`result`=?,`status`=?,`updated_at`=? WHERE `uuid` = ?",
		)).WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			"update-uuid",
		).WillReturnResult(sqlmock.NewResult(0, 1))
		mockDB.ExpectCommit()

		err = r.Update(ctx, job)

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
			ctx = context.Background()
			r   = &jobsRepository{gormDB}
			job = &entity.Job{
				UUID:      "error-uuid",
				Payload:   []byte(`{"test":"data"}`),
				ResultEnc: []byte(`{}`),
				Status:    entity.JobStatusFailed,
				Type:      entity.JobTypeWebScrapper,
			}
			wantErr = errors.New("database error")
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"UPDATE `jobs` SET `payload`=?,`result`=?,`status`=?,`updated_at`=? WHERE `uuid` = ?",
		)).WillReturnError(wantErr)
		mockDB.ExpectRollback()

		err = r.Update(ctx, job)

		assert.Equal(t, wantErr, err)
	})
}