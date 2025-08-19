package sql

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/xwb1989/sqlparser"
)

type accessType string

const (
	Select accessType = "select"
	Insert accessType = "insert"
	Update accessType = "update"
	Delete accessType = "delete"
)

type ResourceAccess struct {
	AccessType string `json:"access_type"`
	Schema     string `json:"schema"`
	Table      string `json:"table"`
}

// ParseSQLToResourceAccess parses a SQL query string and extracts resource access information
// for each statement found in the query. It supports SELECT, INSERT, DELETE, and DDL statements,
// and returns a slice of ResourceAccess pointers representing the accessed resources and their
// access types. Unsupported or unhandled statement types are skipped, except for unknown types,
// which result in an error. If parsing or processing any statement fails, an error is returned.
//
// Parameters:
//   - query: The SQL query string to parse.
//
// Returns:
//   - []*ResourceAccess: A slice of pointers to ResourceAccess structs describing the accessed resources.
//   - error: An error if parsing or processing fails, or nil on success.
func ParseSQLToResourceAccess(query string) ([]*ResourceAccess, error) {
	query = normalizeQuery(query)
	tokens := sqlparser.NewTokenizer(bytes.NewReader([]byte(query)))

	var result []*ResourceAccess
	for {
		stmt, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse SQL: %w", err)
		}
		switch stmt := stmt.(type) {
		case *sqlparser.Select:
			selectStatements, err := processSelectStatement(stmt)
			if err != nil {
				return nil, fmt.Errorf("failed to process select statement: %w", err)
			}
			result = append(result, selectStatements...)

		case *sqlparser.Insert:
			insertStatements, err := processInsertStatement(stmt)
			if err != nil {
				return nil, fmt.Errorf("failed to process insert statement: %w", err)
			}
			result = append(result, insertStatements...)
		case *sqlparser.Update:
			continue
		case *sqlparser.Delete:
			deleteResult, err := processDeleteStatement(stmt)
			if err != nil {
				return nil, fmt.Errorf("failed to process delete statement: %w", err)
			}
			result = append(result, deleteResult...)

		case *sqlparser.DBDDL: // DBDDL represents a CREATE, DROP database statement.
			continue
		case *sqlparser.DDL: // DDL represents a CREATE, ALTER, DROP, RENAME or TRUNCATE statement.
			result = append(result, &ResourceAccess{
				AccessType: string(stmt.Action),
				Schema:     stmt.NewName.Qualifier.String(),
				Table:      stmt.NewName.Name.String(),
			})
			ddlResult, err := processExpr(stmt.TableSpec, accessType(stmt.Action))
			if err != nil {
				return nil, fmt.Errorf("failed to process DDL statement: %w", err)
			}
			result = append(result, ddlResult...)

			continue
		case *sqlparser.Set: // Set represents a SET statement.
			continue
		case *sqlparser.Show: // Show represents a SHOW statement.
			continue
		case *sqlparser.Use: // Use represents a USE statement.
			continue
		default:
			return nil, fmt.Errorf("unsupported statement type: %T", stmt)
		}
	}
	return result, nil
}

func processDeleteStatement(stmt *sqlparser.Delete) ([]*ResourceAccess, error) {
	var result []*ResourceAccess
	accessTypeForTableExprs := Delete
	for _, tbl := range stmt.Targets {
		tableName := tbl.Name.String()
		schemaName := tbl.Qualifier.String()
		if schemaName == "" {
			continue
		}
		accessTypeForTableExprs = Select
		result = append(result, &ResourceAccess{
			AccessType: string(Delete),
			Schema:     schemaName,
			Table:      tableName,
		})
	}

	for _, tbl := range stmt.TableExprs {
		tblResult, err := processExpr(tbl, accessTypeForTableExprs)
		if err != nil {
			return nil, fmt.Errorf("failed to process delete statement: %w", err)
		}
		result = append(result, tblResult...)
	}

	if stmt.Where != nil {
		whereResult, err := processExpr(stmt.Where.Expr, Select)
		if err != nil {
			return nil, fmt.Errorf("failed to process delete where clause: %w", err)
		}
		result = append(result, whereResult...)
	}
	return result, nil
}

func processInsertStatement(stmt *sqlparser.Insert) ([]*ResourceAccess, error) {
	result, err := processExpr(stmt.Rows, Select)
	if err != nil {
		return nil, fmt.Errorf("failed to process insert statement: %w", err)
	}
	tableName := stmt.Table.Name.String()
	schemaName := stmt.Table.Qualifier.String()

	if schemaName == "" {
		return nil, nil
	}

	resourceAccess := &ResourceAccess{
		AccessType: "insert",
		Schema:     schemaName,
		Table:      tableName,
	}
	result = append(result, resourceAccess)
	return result, nil
}

func processSelectStatement(stmt *sqlparser.Select) ([]*ResourceAccess, error) {
	var result []*ResourceAccess
	var err error
	var resourceAccesses []*ResourceAccess
	for _, tbl := range stmt.From {
		result, err = processExpr(tbl, Select)
		if err != nil {
			return nil, err
		}
		resourceAccesses = append(resourceAccesses, result...)

	}
	if stmt.Where != nil {
		result, err = processExpr(stmt.Where.Expr, Select)
		if err != nil {
			return nil, err
		}
		resourceAccesses = append(resourceAccesses, result...)
	}
	if stmt.Having != nil && stmt.Having.Expr != nil {
		result, err = processExpr(stmt.Having.Expr, Select)
		if err != nil {
			return nil, err
		}
		resourceAccesses = append(resourceAccesses, result...)
	}

	return resourceAccesses, nil
}

func processExpr(expr sqlparser.SQLNode, accessType accessType) ([]*ResourceAccess, error) {
	if expr == nil {
		return nil, nil
	}
	var result []*ResourceAccess
	switch node := expr.(type) {
	case *sqlparser.AliasedTableExpr:
		return processExpr(node.Expr, accessType)
	case *sqlparser.JoinTableExpr:
		// Recursively process left and right sides of the join
		leftResult, err := processExpr(node.LeftExpr, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process join left expression: %w", err)
		}
		result = append(result, leftResult...)
		rightResult, err := processExpr(node.RightExpr, Select) // in join right table is usually SELECT
		if err != nil {
			return nil, fmt.Errorf("failed to process join right expression: %w", err)
		}
		result = append(result, rightResult...)
	case *sqlparser.ComparisonExpr:
		// Recursively process left and right sides of the comparison
		leftResult, err := processExpr(node.Left, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process comparison left expression: %w", err)
		}
		result = append(result, leftResult...)
		rightResult, err := processExpr(node.Right, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process comparison right expression: %w", err)
		}
		result = append(result, rightResult...)
	case *sqlparser.ColName:
		// heimdall doesn't need column support for now
		return nil, nil
	case *sqlparser.Subquery:
		return processExpr(node.Select, accessType)
	case *sqlparser.Select:
		return processSelectStatement(node)
	case sqlparser.TableName:
		if node.Qualifier.IsEmpty() {
			return nil, fmt.Errorf("table name must be fully qualified")
		}
		tableName := node.Name.String()
		result = append(result, &ResourceAccess{
			AccessType: string(accessType),
			Schema:     node.Qualifier.String(),
			Table:      tableName,
		})

	case *sqlparser.ExistsExpr:
		return processExpr(node.Subquery, accessType)
	case sqlparser.Values:
		for _, valTuple := range node {
			tableResults, err := processExpr(valTuple, accessType)
			if err != nil {
				return nil, fmt.Errorf("failed to process values row: %w", err)
			}
			result = append(result, tableResults...)
		}
	case sqlparser.ValTuple:
		for _, val := range node {
			valResult, err := processExpr(val, accessType)
			if err != nil {
				return nil, fmt.Errorf("failed to process value tuple: %w", err)
			}
			result = append(result, valResult...)
		}
	case *sqlparser.SQLVal:
		// heimdall doesn't need SQLVal support for now
		return nil, nil
	case *sqlparser.AndExpr:
		leftResult, err := processExpr(node.Left, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process AND expression left side: %w", err)
		}
		result = append(result, leftResult...)
		rightResult, err := processExpr(node.Right, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process AND expression right side: %w", err)
		}
		result = append(result, rightResult...)
	case *sqlparser.OrExpr:
		leftResult, err := processExpr(node.Left, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process OR expression left side: %w", err)
		}
		result = append(result, leftResult...)
		rightResult, err := processExpr(node.Right, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process OR expression right side: %w", err)
		}
		result = append(result, rightResult...)
	case *sqlparser.NotExpr:
		notResult, err := processExpr(node.Expr, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process NOT expression: %w", err)
		}
		result = append(result, notResult...)
	case *sqlparser.IsExpr:
		isResult, err := processExpr(node.Expr, accessType)
		if err != nil {
			return nil, fmt.Errorf("failed to process IS expression: %w", err)
		}
		result = append(result, isResult...)
	case *sqlparser.ParenExpr:
		return processExpr(node.Expr, accessType)
	case *sqlparser.TableSpec:
		// TableSpec is used in DDL statements we don't need to handle it now
		return result, nil
	case sqlparser.TableExprs:
		for _, tbl := range node {
			tblResult, err := processExpr(tbl, accessType)
			if err != nil {
				return nil, fmt.Errorf("failed to process table expression: %w", err)
			}
			result = append(result, tblResult...)
		}

	default:
		return nil, fmt.Errorf("unsupported expr type: %T", node)
	}
	return result, nil
}

func normalizeQuery(query string) string {

	// sqlparser doesn't support TEMPORARY VIEW, so we replace it with VIEW
	return strings.ReplaceAll(query, "TEMPORARY VIEW", "VIEW")

}
