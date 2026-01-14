package specifications

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/luci/go-render/render"
	"go.uber.org/mock/gomock"
)

type matcherValue struct {
	fieldName     string
	actualValue   any
	expectedValue any
}

type customMatcher[V any] struct {
	match  func(v V) (bool, matcherValue)
	arg    any
	mValue matcherValue
}

func (cm *customMatcher[V]) Matches(arg any) bool {
	cm.arg = arg

	v, ok := arg.(V)
	if !ok {
		return false
	}

	isMatched, mValue := cm.match(v)

	cm.mValue = mValue

	return isMatched
}

func (cm *customMatcher[V]) String() string {
	actual := asCode(cm.mValue.actualValue)
	expected := asCode(cm.mValue.expectedValue)

	return fmt.Sprintf(
		"%v\nField: %s\nActual: %v\nExpected: %v",
		cm.arg,
		cm.mValue.fieldName,
		actual,
		expected,
	)
}

func CustomMatcher[V any](f func(v V) (bool, matcherValue)) gomock.Matcher {
	return &customMatcher[V]{
		match: f,
	}
}

func asCode(v any) string {
	s := render.Render(v)

	s = strings.ReplaceAll(s, "*", "&")
	s = strings.ReplaceAll(s, "[]&", "[]*")

	// remove extra parens wrapping types, e.g. (&map[string]int){"bar":1, "foo":0} -> &map[string]int{"bar":1, "foo":0}
	re := regexp.MustCompile(`\(([\w.*&\[\]]*)\){`)

	s = re.ReplaceAllString(s, "$1{")

	// format nils as runable code e.g. (&render.innerStruct)(nil), -> nil,
	re = regexp.MustCompile(`\([*&\w]*.\w*\)\(nil\),`)

	s = re.ReplaceAllString(s, "nil,")

	return s
}

func SpecMatcher(
	expectedSpec I,
) func(I) (bool, matcherValue) {
	return func(actualSpec I) (bool, matcherValue) {
		var (
			actualQuery, actualErr     = actualSpec.ToGet()
			expectedQuery, expectedErr = expectedSpec.ToGet()
		)

		if !errors.Is(actualErr, expectedErr) {
			return false, matcherValue{
				fieldName:     "ToSQL() error",
				actualValue:   actualErr,
				expectedValue: expectedErr,
			}
		}

		if actualQuery != expectedQuery {
			return false, matcherValue{
				fieldName:     "SQL query",
				actualValue:   actualQuery,
				expectedValue: expectedQuery,
			}
		}

		return true, matcherValue{}
	}
}
