package database

import (
	"fmt"
	"strings"
	"text/template"
)

const (
	clauseKey = `Clause`
	sliceKey  = `Slice`
)

var (
	queryTemplateParser = template.New(`query`)
	sliceTemplateParser = template.New(`slice`)
)

type Filter map[string]string

type FilterParameter struct {
	IsSlice bool
	Value   string
	Item    string
	Join    string
}

type FilterConfig struct {
	Join       string
	Parameters map[string]*FilterParameter
}

func (f *Filter) Render(queryTemplate string, config *FilterConfig, args ...any) (string, []any, error) {

	// append slice to args
	start := len(args) + 1

	// let's go over filters we got
	conditions := make([]string, 0, len(*f)+10)
	for filter, value := range *f {
		if param, found := config.Parameters[filter]; found {
			if !param.IsSlice {
				conditions = append(conditions, fmt.Sprintf(param.Value, start))
				args = append(args, value)
				start += 1
			} else {
				// let's split slice first
				slice := strings.Split(value, `,`)
				if len(slice) > 0 {

					sliceConditions := make([]string, 0, len(slice))
					for _, item := range slice {
						sliceConditions = append(sliceConditions, fmt.Sprintf(param.Item, start))
						args = append(args, item)
						start += 1
					}
					parsedSliceTemplate, err := sliceTemplateParser.Parse(param.Value)
					if err != nil {
						return ``, nil, err
					}

					// render slice template
					var preparedSlice strings.Builder
					if err := parsedSliceTemplate.Execute(&preparedSlice, map[string]string{
						sliceKey: strings.Join(sliceConditions, param.Join),
					}); err != nil {
						return ``, nil, err
					}
					conditions = append(conditions, preparedSlice.String())

				}
			}
		}
	}

	// let's parse template
	parsedTemplate, err := queryTemplateParser.Parse(queryTemplate)
	if err != nil {
		return ``, nil, err
	}

	// render template
	var preparedQuery strings.Builder
	if err := parsedTemplate.Execute(&preparedQuery, map[string]string{
		clauseKey: strings.Join(conditions, config.Join),
	}); err != nil {
		return ``, nil, err
	}

	return preparedQuery.String(), args, nil

}
