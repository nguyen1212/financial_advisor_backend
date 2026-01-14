package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/financial_advisor/app/domain/entity"
	"github.com/financial_advisor/app/domain/repository"
	"github.com/financial_advisor/app/domain/repository/specifications"
	appErrors "github.com/financial_advisor/app/errors"
	gormAdapter "github.com/financial_advisor/app/external/db/gorm"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type jobsRepository struct{ db *gorm.DB }

func NewJobsRepository(db *gormAdapter.MySQL) repository.JobsRepository {
	return &jobsRepository{db.DB()}
}

func (r *jobsRepository) Get(
	ctx context.Context,
	spec specifications.I,
) (entity.Job, error) {
	var (
		job entity.Job
		tx  = r.db.WithContext(ctx)
	)

	sql, err := spec.ToGet()
	if err != nil {
		return entity.Job{}, fmt.Errorf("convert to raw get sql: %w", err)
	}

	if err := tx.WithContext(ctx).Raw(sql).Take(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Job{}, appErrors.ErrNotFound
		}

		return entity.Job{}, err
	}

	return job, nil
}

func (r *jobsRepository) Create(
	ctx context.Context,
	job *entity.Job,
) error {
	if err := r.db.WithContext(ctx).Create(job).Error; err != nil {
		var mySQLErr *mysql.MySQLError

		if errors.As(err, &mySQLErr) && mySQLErr.Number == duplicateEntryMysqlCode {
			return appErrors.ErrConflicted
		}

		return err
	}

	return nil
}

func (r *jobsRepository) Update(
	ctx context.Context,
	job *entity.Job,
) error {
	if err := r.db.
		WithContext(ctx).
		Model(job).
		Updates(map[string]any{
			"status":  job.Status,
			"payload": job.Payload,
			"result":  job.ResultEnc,
		}).Error; err != nil {
		return err
	}

	return nil
}
