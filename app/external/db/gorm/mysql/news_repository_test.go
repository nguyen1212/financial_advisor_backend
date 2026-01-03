package mysql

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

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

func Test_NewNewsRepository(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		&newsRepository{},
		NewNewsRepository(&gormAdapter.MySQL{}),
	)
}

func Test_newsRepository_Count(t *testing.T) {
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
			r    = &newsRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)
		)

		spec.EXPECT().ToCount().Return("SELECT COUNT(*) FROM `news`", nil)

		rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(10)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `news`")).WillReturnRows(rows)

		count, err := r.Count(ctx, spec)

		require.NoError(t, err)
		assert.Equal(t, int64(10), count)
	})

	t.Run("failed - count news", func(t *testing.T) {
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
			r    = &newsRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)

			wantErr = errors.New("database error")
		)

		spec.EXPECT().ToCount().Return("SELECT COUNT(*) FROM `news`", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `news`")).WillReturnError(wantErr)

		_, err = r.Count(ctx, spec)

		assert.Equal(t, wantErr, err)
	})

	t.Run("failed - count news", func(t *testing.T) {
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
			r    = &newsRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)

			sqlErr  = errors.New("database error")
			wantErr = fmt.Errorf("convert to raw count sql: %w", sqlErr)
		)

		spec.EXPECT().ToCount().Return("", sqlErr)

		_, err = r.Count(ctx, spec)

		assert.Equal(t, wantErr, err)
	})
}

func Test_newsRepository_Find(t *testing.T) {
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
			r      = &newsRepository{gormDB}
			spec   = mock.NewMockI(mockCtrl)
			paging = mock.NewMockPagingI(mockCtrl)
		)

		spec.EXPECT().ToFind(paging).Return("SELECT * FROM `news` LIMIT 10", nil)

		rows := sqlmock.NewRows([]string{"id", "title", "author", "thumbnail", "url", "status", "category"}).
			AddRow(1, "Test News", "Author 1", "thumb1.jpg", "http://test1.com", 1, 1).
			AddRow(2, "Test News 2", "Author 2", "thumb2.jpg", "http://test2.com", 1, 2)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `news` LIMIT 10")).WillReturnRows(rows)

		newsList, err := r.Find(ctx, spec, paging)

		require.NoError(t, err)
		assert.Len(t, newsList, 2)

		// Validate first news item
		assert.Equal(t, uint64(1), newsList[0].ID)
		assert.Equal(t, "Test News", newsList[0].Title)
		assert.Equal(t, "Author 1", newsList[0].Author)
		assert.Equal(t, "thumb1.jpg", newsList[0].Thumbnail)
		assert.Equal(t, "http://test1.com", newsList[0].URL)

		// Validate second news item
		assert.Equal(t, uint64(2), newsList[1].ID)
		assert.Equal(t, "Test News 2", newsList[1].Title)
		assert.Equal(t, "Author 2", newsList[1].Author)
		assert.Equal(t, "thumb2.jpg", newsList[1].Thumbnail)
		assert.Equal(t, "http://test2.com", newsList[1].URL)
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
			r       = &newsRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			paging  = mock.NewMockPagingI(mockCtrl)
			wantErr = errors.New("database error")
		)

		spec.EXPECT().ToFind(paging).Return("SELECT * FROM `news` LIMIT 10", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `news` LIMIT 10")).WillReturnError(wantErr)

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
			r       = &newsRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			paging  = mock.NewMockPagingI(mockCtrl)
			wantErr = errors.New("spec error")
		)

		spec.EXPECT().ToFind(paging).Return("", wantErr)

		_, err = r.Find(ctx, spec, paging)

		assert.ErrorContains(t, err, "convert to raw find sql")
	})
}

func Test_newsRepository_Get(t *testing.T) {
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
			r    = &newsRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `news` WHERE id = 12", nil)

		rows := sqlmock.NewRows([]string{"id", "title", "author", "thumbnail", "url", "status", "category"}).
			AddRow(12, "Test News", "Author 1", "thumb1.jpg", "http://test1.com", 1, 1)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `news` WHERE id = 12")).WillReturnRows(rows)

		news, err := r.Get(ctx, spec)

		require.NoError(t, err)
		assert.Equal(t, uint64(12), news.ID)
		assert.Equal(t, "Test News", news.Title)
		assert.Equal(t, "Author 1", news.Author)
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
			r    = &newsRepository{gormDB}
			spec = mock.NewMockI(mockCtrl)
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `news` WHERE id = 12", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `news` WHERE id = 12")).WillReturnError(gorm.ErrRecordNotFound)

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
			r       = &newsRepository{gormDB}
			spec    = mock.NewMockI(mockCtrl)
			wantErr = errors.New("database error")
		)

		spec.EXPECT().ToGet().Return("SELECT * FROM `news` WHERE id = 12", nil)

		mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `news` WHERE id = 12")).WillReturnError(wantErr)

		_, err = r.Get(ctx, spec)

		assert.Equal(t, wantErr, err)
	})
}

func Test_newsRepository_Create(t *testing.T) {
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
			ctx  = context.Background()
			r    = &newsRepository{gormDB}
			news = &entity.News{
				Title:     "Test News",
				Author:    "Test Author",
				Thumbnail: "",
				URL:       "https://vnexpress.net/the-gioi-bat-dau-don-nam-moi-2026-5000678.html",
				HashedURL: []byte("<binary>"),
			}
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `news` (`title`,`author`,`thumbnail`,`url`,`hashed_url`,`status`,`category`,`file_path`,`file_size`,`news_with_fulltext_id`,`publisher_id`,`published_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.ExpectCommit()

		err = r.Create(ctx, news)

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
			ctx  = context.Background()
			r    = &newsRepository{gormDB}
			news = &entity.News{
				Title:     "Test News",
				Author:    "Test Author",
				Thumbnail: "test.jpg",
				URL:       "http://test.com",
			}
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `news` (`title`,`author`,`thumbnail`,`url`,`hashed_url`,`status`,`category`,`file_path`,`file_size`,`news_with_fulltext_id`,`publisher_id`,`published_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		)).
			WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
		mockDB.ExpectRollback()

		err = r.Create(ctx, news)

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
			ctx  = context.Background()
			r    = &newsRepository{gormDB}
			news = &entity.News{
				Title:     "Test News",
				Author:    "Test Author",
				Thumbnail: "test.jpg",
				URL:       "http://test.com",
			}
			wantErr = errors.New("database error")
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"INSERT INTO `news` (`title`,`author`,`thumbnail`,`url`,`hashed_url`,`status`,`category`,`file_path`,`file_size`,`news_with_fulltext_id`,`publisher_id`,`published_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		)).
			WillReturnError(wantErr)
		mockDB.ExpectRollback()

		err = r.Create(ctx, news)

		assert.Equal(t, wantErr, err)
	})
}

func Test_newsRepository_Update(t *testing.T) {
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
			ctx  = context.Background()
			r    = &newsRepository{gormDB}
			news = &entity.News{
				ID:        1,
				Title:     "Updated News",
				Author:    "Updated Author",
				Thumbnail: "updated.jpg",
				URL:       "http://updated.com",
			}
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"UPDATE `news` SET `author`=?,`file_size`=?,`news_with_fulltext_id`=?,`published_at`=?,`status`=?,`thumbnail`=?,`title`=? WHERE `id` = ?",
		)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mockDB.ExpectCommit()

		err = r.Update(ctx, news)

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
			ctx  = context.Background()
			r    = &newsRepository{gormDB}
			news = &entity.News{
				ID:        1,
				Title:     "Updated News",
				Author:    "Updated Author",
				Thumbnail: "updated.jpg",
				URL:       "http://updated.com",
			}
			wantErr = errors.New("database error")
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(
			"UPDATE `news` SET `author`=?,`file_size`=?,`news_with_fulltext_id`=?,`published_at`=?,`status`=?,`thumbnail`=?,`title`=? WHERE `id` = ?",
		)).
			WillReturnError(wantErr)
		mockDB.ExpectRollback()

		err = r.Update(ctx, news)

		assert.Equal(t, wantErr, err)
	})
}

func Test_newsRepository_Delete(t *testing.T) {
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
			ctx    = context.Background()
			r      = &newsRepository{gormDB}
			newsID = uint64(1)
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `news` WHERE `news`.`id` = ?")).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mockDB.ExpectCommit()

		err = r.Delete(ctx, newsID)

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
			r         = &newsRepository{gormDB}
			newsID    = uint64(1)
			deleteErr = errors.New("database error")
		)

		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `news` WHERE `news`.`id` = ?")).
			WillReturnError(deleteErr)
		mockDB.ExpectRollback()

		err = r.Delete(ctx, newsID)

		assert.Equal(t, deleteErr, err)
	})
}
