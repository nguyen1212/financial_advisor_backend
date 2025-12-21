// Package specifications defines repository specifications interface
package specifications

import "gorm.io/gorm"

// I interface supports builder pattern only
type I[T gorm.DB] interface {
	Query(*T) *T
}

type Paging interface {
	Limit() int
	Offset() int
}
