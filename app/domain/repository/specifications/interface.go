// Package specifications defines repository specifications interface
package specifications

import "gorm.io/gorm"

// I interface supports builder pattern only
type I interface {
	// for any external packages that are not gorm related
	// we can provide a dry gorm DB instance to generate SQL query
	Query(*gorm.DB) *gorm.DB
}

type PagingI interface {
	Limit() int
	Offset() int
}
