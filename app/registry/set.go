package registry

import (
	"github.com/financial_advisor/app/external/db/gorm"
	"github.com/financial_advisor/app/external/db/gorm/manticore"
	"github.com/financial_advisor/app/external/db/gorm/mysql"
	"github.com/financial_advisor/app/external/hasher"
	"github.com/financial_advisor/app/external/uuid"
	"github.com/financial_advisor/app/services/consumer"
	"github.com/financial_advisor/app/services/worker"
	"github.com/google/wire"
)

var singletonSet = wire.NewSet(
	gorm.GetMySQLIns,
	gorm.GetManticoreIns,
	worker.Get,
)

var repositorySet = wire.NewSet(
	mysql.NewNewsRepository,
	mysql.NewPublisherRepository,
	mysql.NewJobsRepository,
	manticore.NewNewsRepository,
)

var serviceSet = wire.NewSet(
	hasher.NewMD5,
	uuid.NewUUIDv7,
	consumer.NewManager,
)
