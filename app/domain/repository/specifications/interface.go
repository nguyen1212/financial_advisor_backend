// Package specifications defines repository specifications interface
package specifications

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

type I interface {
	// string: SQL query with/without placeholders
	// []any: values for placeholders
	// error: error if any
	ToCount() (string, error)
	ToFind(paging PagingI) (string, error)
	ToGet() (string, error)
}

type PagingI interface {
	Limit() int
	Offset() int
}
