package repository

import "github.com/mytrix-technology/mylibgo/datastore"

type queryConfig struct {
	filters []datastore.Filter
	limit int
	offset int
	sorts datastore.Sorts
}

type QueryOption func(qc *queryConfig)

func AddFilter(field string, comparison datastore.ComparisonString, value interface{}) QueryOption {
	return func(qc *queryConfig) {
		filter := datastore.Filter{
			Field:     field,
			Comparison: comparison,
			Value:      value,
		}

		qc.filters = append(qc.filters, filter)
	}
}

func WithFilters(filters ...datastore.Filter) QueryOption {
	return func(qc *queryConfig) {
		qc.filters = filters
	}
}

func WithLimit(limit int) QueryOption {
	return func(qc *queryConfig) {
		qc.limit = limit
	}
}

func WithOffset(offset int) QueryOption {
	return func(qc *queryConfig) {
		qc.offset = offset
	}
}

func WithSort(sorts ...datastore.Sort) QueryOption {
	return func(qc *queryConfig) {
		qc.sorts = sorts
	}
}
