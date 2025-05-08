package database

import (
	"fmt"
	"strings"
	"text/template"
)

const (
	sliceTag = `Slice`
)

func PrepareSliceQuery(queryTemplate, sliceTemplate string, slice []any, args ...any) (string, []any, error) {

	// append slice to args
	start := len(args) + 1
	args = append(args, slice...)

	items := make([]string, 0, len(slice))
	for range slice {
		items = append(items, fmt.Sprintf(sliceTemplate, start))
		start += 1
	}

	parsedTemplate, err := template.New(`query`).Parse(queryTemplate)
	if err != nil {
		return ``, nil, err
	}

	var preparedQuery strings.Builder
	if err := parsedTemplate.Execute(&preparedQuery, map[string]string{
		sliceTag: strings.Join(items, `, `),
	}); err != nil {
		return ``, nil, err
	}

	return preparedQuery.String(), args, nil

}
