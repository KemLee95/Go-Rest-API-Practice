package collection_helper

import (
	"github.com/samber/lo"
)

func Unique[T comparable](items []T) []T {
	return lo.Uniq(items)
}
