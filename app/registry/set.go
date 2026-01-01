package registry

import (
	"github.com/financial_advisor/app/external/db/gorm"
	"github.com/financial_advisor/app/external/db/gorm/manticore"
	"github.com/financial_advisor/app/external/db/gorm/mysql"
	"github.com/financial_advisor/app/external/hasher"
	memoryqueue "github.com/financial_advisor/app/external/queue/memory-queue"
	"github.com/google/wire"
)

var singletonSet = wire.NewSet(
	gorm.GetMySQLIns,
	gorm.GetManticoreIns,
	memoryqueue.Get,
)

var repositorySet = wire.NewSet(
	mysql.NewNewsRepository,
	mysql.NewPublisherRepository,
	manticore.NewNewsRepository,
)

var serviceSet = wire.NewSet(
	hasher.NewMD5,
)
