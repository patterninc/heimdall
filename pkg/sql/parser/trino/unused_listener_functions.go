package trino

import (
	"log"
	"reflect"

	"github.com/patterninc/heimdall/pkg/sql/parser/trino/grammar"
)

// EnterParse is called when production parse is entered.
func (l *trinoListener) EnterParse(ctx *grammar.ParseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitParse is called when production parse is exited.
func (l *trinoListener) ExitParse(ctx *grammar.ParseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterStandaloneExpression is called when production standaloneExpression is entered.
func (l *trinoListener) EnterStandaloneExpression(ctx *grammar.StandaloneExpressionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitStandaloneExpression is called when production standaloneExpression is exited.
func (l *trinoListener) ExitStandaloneExpression(ctx *grammar.StandaloneExpressionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterStandalonePathSpecification is called when production standalonePathSpecification is entered.
func (l *trinoListener) EnterStandalonePathSpecification(ctx *grammar.StandalonePathSpecificationContext) {
}

// ExitStandalonePathSpecification is called when production standalonePathSpecification is exited.
func (l *trinoListener) ExitStandalonePathSpecification(ctx *grammar.StandalonePathSpecificationContext) {
}

// EnterStandaloneType is called when production standaloneType is entered.
func (l *trinoListener) EnterStandaloneType(ctx *grammar.StandaloneTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitStandaloneType is called when production standaloneType is exited.
func (l *trinoListener) ExitStandaloneType(ctx *grammar.StandaloneTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterStandaloneRowPattern is called when production standaloneRowPattern is entered.
func (l *trinoListener) EnterStandaloneRowPattern(ctx *grammar.StandaloneRowPatternContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitStandaloneRowPattern is called when production standaloneRowPattern is exited.
func (l *trinoListener) ExitStandaloneRowPattern(ctx *grammar.StandaloneRowPatternContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterStandaloneFunctionSpecification is called when production standaloneFunctionSpecification is entered.
func (l *trinoListener) EnterStandaloneFunctionSpecification(ctx *grammar.StandaloneFunctionSpecificationContext) {
}

// ExitStandaloneFunctionSpecification is called when production standaloneFunctionSpecification is exited.
func (l *trinoListener) ExitStandaloneFunctionSpecification(ctx *grammar.StandaloneFunctionSpecificationContext) {
}

// EnterUse is called when production use is entered.
func (l *trinoListener) EnterUse(ctx *grammar.UseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitUse is called when production use is exited.
func (l *trinoListener) ExitUse(ctx *grammar.UseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCreateCatalog is called when production createCatalog is entered.
func (l *trinoListener) EnterCreateCatalog(ctx *grammar.CreateCatalogContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCreateCatalog is called when production createCatalog is exited.
func (l *trinoListener) ExitCreateCatalog(ctx *grammar.CreateCatalogContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDropCatalog is called when production dropCatalog is entered.
func (l *trinoListener) EnterDropCatalog(ctx *grammar.DropCatalogContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDropCatalog is called when production dropCatalog is exited.
func (l *trinoListener) ExitDropCatalog(ctx *grammar.DropCatalogContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCreateSchema is called when production createSchema is entered.
func (l *trinoListener) EnterCreateSchema(ctx *grammar.CreateSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCreateSchema is called when production createSchema is exited.
func (l *trinoListener) ExitCreateSchema(ctx *grammar.CreateSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDropSchema is called when production dropSchema is entered.
func (l *trinoListener) EnterDropSchema(ctx *grammar.DropSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDropSchema is called when production dropSchema is exited.
func (l *trinoListener) ExitDropSchema(ctx *grammar.DropSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRenameSchema is called when production renameSchema is entered.
func (l *trinoListener) EnterRenameSchema(ctx *grammar.RenameSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRenameSchema is called when production renameSchema is exited.
func (l *trinoListener) ExitRenameSchema(ctx *grammar.RenameSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetSchemaAuthorization is called when production setSchemaAuthorization is entered.
func (l *trinoListener) EnterSetSchemaAuthorization(ctx *grammar.SetSchemaAuthorizationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetSchemaAuthorization is called when production setSchemaAuthorization is exited.
func (l *trinoListener) ExitSetSchemaAuthorization(ctx *grammar.SetSchemaAuthorizationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTruncateTable is called when production truncateTable is entered.
func (l *trinoListener) EnterTruncateTable(ctx *grammar.TruncateTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTruncateTable is called when production truncateTable is exited.
func (l *trinoListener) ExitTruncateTable(ctx *grammar.TruncateTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCommentTable is called when production commentTable is entered.
func (l *trinoListener) EnterCommentTable(ctx *grammar.CommentTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCommentTable is called when production commentTable is exited.
func (l *trinoListener) ExitCommentTable(ctx *grammar.CommentTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCommentView is called when production commentView is entered.
func (l *trinoListener) EnterCommentView(ctx *grammar.CommentViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCommentView is called when production commentView is exited.
func (l *trinoListener) ExitCommentView(ctx *grammar.CommentViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCommentColumn is called when production commentColumn is entered.
func (l *trinoListener) EnterCommentColumn(ctx *grammar.CommentColumnContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCommentColumn is called when production commentColumn is exited.
func (l *trinoListener) ExitCommentColumn(ctx *grammar.CommentColumnContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRenameColumn is called when production renameColumn is entered.
func (l *trinoListener) EnterRenameColumn(ctx *grammar.RenameColumnContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRenameColumn is called when production renameColumn is exited.
func (l *trinoListener) ExitRenameColumn(ctx *grammar.RenameColumnContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetColumnType is called when production setColumnType is entered.
func (l *trinoListener) EnterSetColumnType(ctx *grammar.SetColumnTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetColumnType is called when production setColumnType is exited.
func (l *trinoListener) ExitSetColumnType(ctx *grammar.SetColumnTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetTableAuthorization is called when production setTableAuthorization is entered.
func (l *trinoListener) EnterSetTableAuthorization(ctx *grammar.SetTableAuthorizationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetTableAuthorization is called when production setTableAuthorization is exited.
func (l *trinoListener) ExitSetTableAuthorization(ctx *grammar.SetTableAuthorizationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetTableProperties is called when production setTableProperties is entered.
func (l *trinoListener) EnterSetTableProperties(ctx *grammar.SetTablePropertiesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetTableProperties is called when production setTableProperties is exited.
func (l *trinoListener) ExitSetTableProperties(ctx *grammar.SetTablePropertiesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTableExecute is called when production tableExecute is entered.
func (l *trinoListener) EnterTableExecute(ctx *grammar.TableExecuteContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTableExecute is called when production tableExecute is exited.
func (l *trinoListener) ExitTableExecute(ctx *grammar.TableExecuteContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterAnalyze is called when production analyze is entered.
func (l *trinoListener) EnterAnalyze(ctx *grammar.AnalyzeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitAnalyze is called when production analyze is exited.
func (l *trinoListener) ExitAnalyze(ctx *grammar.AnalyzeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCreateMaterializedView is called when production createMaterializedView is entered.
func (l *trinoListener) EnterCreateMaterializedView(ctx *grammar.CreateMaterializedViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCreateMaterializedView is called when production createMaterializedView is exited.
func (l *trinoListener) ExitCreateMaterializedView(ctx *grammar.CreateMaterializedViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCreateView is called when production createView is entered.
func (l *trinoListener) EnterCreateView(ctx *grammar.CreateViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCreateView is called when production createView is exited.
func (l *trinoListener) ExitCreateView(ctx *grammar.CreateViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRefreshMaterializedView is called when production refreshMaterializedView is entered.
func (l *trinoListener) EnterRefreshMaterializedView(ctx *grammar.RefreshMaterializedViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRefreshMaterializedView is called when production refreshMaterializedView is exited.
func (l *trinoListener) ExitRefreshMaterializedView(ctx *grammar.RefreshMaterializedViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDropMaterializedView is called when production dropMaterializedView is entered.
func (l *trinoListener) EnterDropMaterializedView(ctx *grammar.DropMaterializedViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDropMaterializedView is called when production dropMaterializedView is exited.
func (l *trinoListener) ExitDropMaterializedView(ctx *grammar.DropMaterializedViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRenameMaterializedView is called when production renameMaterializedView is entered.
func (l *trinoListener) EnterRenameMaterializedView(ctx *grammar.RenameMaterializedViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRenameMaterializedView is called when production renameMaterializedView is exited.
func (l *trinoListener) ExitRenameMaterializedView(ctx *grammar.RenameMaterializedViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetMaterializedViewProperties is called when production setMaterializedViewProperties is entered.
func (l *trinoListener) EnterSetMaterializedViewProperties(ctx *grammar.SetMaterializedViewPropertiesContext) {
}

// ExitSetMaterializedViewProperties is called when production setMaterializedViewProperties is exited.
func (l *trinoListener) ExitSetMaterializedViewProperties(ctx *grammar.SetMaterializedViewPropertiesContext) {
}

// EnterDropView is called when production dropView is entered.
func (l *trinoListener) EnterDropView(ctx *grammar.DropViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDropView is called when production dropView is exited.
func (l *trinoListener) ExitDropView(ctx *grammar.DropViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRenameView is called when production renameView is entered.
func (l *trinoListener) EnterRenameView(ctx *grammar.RenameViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRenameView is called when production renameView is exited.
func (l *trinoListener) ExitRenameView(ctx *grammar.RenameViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetViewAuthorization is called when production setViewAuthorization is entered.
func (l *trinoListener) EnterSetViewAuthorization(ctx *grammar.SetViewAuthorizationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetViewAuthorization is called when production setViewAuthorization is exited.
func (l *trinoListener) ExitSetViewAuthorization(ctx *grammar.SetViewAuthorizationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCall is called when production call is entered.
func (l *trinoListener) EnterCall(ctx *grammar.CallContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCall is called when production call is exited.
func (l *trinoListener) ExitCall(ctx *grammar.CallContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCreateFunction is called when production createFunction is entered.
func (l *trinoListener) EnterCreateFunction(ctx *grammar.CreateFunctionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCreateFunction is called when production createFunction is exited.
func (l *trinoListener) ExitCreateFunction(ctx *grammar.CreateFunctionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDropFunction is called when production dropFunction is entered.
func (l *trinoListener) EnterDropFunction(ctx *grammar.DropFunctionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDropFunction is called when production dropFunction is exited.
func (l *trinoListener) ExitDropFunction(ctx *grammar.DropFunctionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCreateRole is called when production createRole is entered.
func (l *trinoListener) EnterCreateRole(ctx *grammar.CreateRoleContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCreateRole is called when production createRole is exited.
func (l *trinoListener) ExitCreateRole(ctx *grammar.CreateRoleContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDropRole is called when production dropRole is entered.
func (l *trinoListener) EnterDropRole(ctx *grammar.DropRoleContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDropRole is called when production dropRole is exited.
func (l *trinoListener) ExitDropRole(ctx *grammar.DropRoleContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterGrantRoles is called when production grantRoles is entered.
func (l *trinoListener) EnterGrantRoles(ctx *grammar.GrantRolesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitGrantRoles is called when production grantRoles is exited.
func (l *trinoListener) ExitGrantRoles(ctx *grammar.GrantRolesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRevokeRoles is called when production revokeRoles is entered.
func (l *trinoListener) EnterRevokeRoles(ctx *grammar.RevokeRolesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRevokeRoles is called when production revokeRoles is exited.
func (l *trinoListener) ExitRevokeRoles(ctx *grammar.RevokeRolesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetRole is called when production setRole is entered.
func (l *trinoListener) EnterSetRole(ctx *grammar.SetRoleContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetRole is called when production setRole is exited.
func (l *trinoListener) ExitSetRole(ctx *grammar.SetRoleContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterGrant is called when production grant is entered.
func (l *trinoListener) EnterGrant(ctx *grammar.GrantContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitGrant is called when production grant is exited.
func (l *trinoListener) ExitGrant(ctx *grammar.GrantContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDeny is called when production deny is entered.
func (l *trinoListener) EnterDeny(ctx *grammar.DenyContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDeny is called when production deny is exited.
func (l *trinoListener) ExitDeny(ctx *grammar.DenyContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRevoke is called when production revoke is entered.
func (l *trinoListener) EnterRevoke(ctx *grammar.RevokeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRevoke is called when production revoke is exited.
func (l *trinoListener) ExitRevoke(ctx *grammar.RevokeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowGrants is called when production showGrants is entered.
func (l *trinoListener) EnterShowGrants(ctx *grammar.ShowGrantsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowGrants is called when production showGrants is exited.
func (l *trinoListener) ExitShowGrants(ctx *grammar.ShowGrantsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExplain is called when production explain is entered.
func (l *trinoListener) EnterExplain(ctx *grammar.ExplainContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExplain is called when production explain is exited.
func (l *trinoListener) ExitExplain(ctx *grammar.ExplainContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExplainAnalyze is called when production explainAnalyze is entered.
func (l *trinoListener) EnterExplainAnalyze(ctx *grammar.ExplainAnalyzeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExplainAnalyze is called when production explainAnalyze is exited.
func (l *trinoListener) ExitExplainAnalyze(ctx *grammar.ExplainAnalyzeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowCreateTable is called when production showCreateTable is entered.
func (l *trinoListener) EnterShowCreateTable(ctx *grammar.ShowCreateTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowCreateTable is called when production showCreateTable is exited.
func (l *trinoListener) ExitShowCreateTable(ctx *grammar.ShowCreateTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowCreateSchema is called when production showCreateSchema is entered.
func (l *trinoListener) EnterShowCreateSchema(ctx *grammar.ShowCreateSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowCreateSchema is called when production showCreateSchema is exited.
func (l *trinoListener) ExitShowCreateSchema(ctx *grammar.ShowCreateSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowCreateView is called when production showCreateView is entered.
func (l *trinoListener) EnterShowCreateView(ctx *grammar.ShowCreateViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowCreateView is called when production showCreateView is exited.
func (l *trinoListener) ExitShowCreateView(ctx *grammar.ShowCreateViewContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowCreateMaterializedView is called when production showCreateMaterializedView is entered.
func (l *trinoListener) EnterShowCreateMaterializedView(ctx *grammar.ShowCreateMaterializedViewContext) {
}

// ExitShowCreateMaterializedView is called when production showCreateMaterializedView is exited.
func (l *trinoListener) ExitShowCreateMaterializedView(ctx *grammar.ShowCreateMaterializedViewContext) {
}

// EnterShowTables is called when production showTables is entered.
func (l *trinoListener) EnterShowTables(ctx *grammar.ShowTablesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowTables is called when production showTables is exited.
func (l *trinoListener) ExitShowTables(ctx *grammar.ShowTablesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowSchemas is called when production showSchemas is entered.
func (l *trinoListener) EnterShowSchemas(ctx *grammar.ShowSchemasContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowSchemas is called when production showSchemas is exited.
func (l *trinoListener) ExitShowSchemas(ctx *grammar.ShowSchemasContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowCatalogs is called when production showCatalogs is entered.
func (l *trinoListener) EnterShowCatalogs(ctx *grammar.ShowCatalogsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowCatalogs is called when production showCatalogs is exited.
func (l *trinoListener) ExitShowCatalogs(ctx *grammar.ShowCatalogsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowColumns is called when production showColumns is entered.
func (l *trinoListener) EnterShowColumns(ctx *grammar.ShowColumnsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowColumns is called when production showColumns is exited.
func (l *trinoListener) ExitShowColumns(ctx *grammar.ShowColumnsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowStats is called when production showStats is entered.
func (l *trinoListener) EnterShowStats(ctx *grammar.ShowStatsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowStats is called when production showStats is exited.
func (l *trinoListener) ExitShowStats(ctx *grammar.ShowStatsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowStatsForQuery is called when production showStatsForQuery is entered.
func (l *trinoListener) EnterShowStatsForQuery(ctx *grammar.ShowStatsForQueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowStatsForQuery is called when production showStatsForQuery is exited.
func (l *trinoListener) ExitShowStatsForQuery(ctx *grammar.ShowStatsForQueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowRoles is called when production showRoles is entered.
func (l *trinoListener) EnterShowRoles(ctx *grammar.ShowRolesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowRoles is called when production showRoles is exited.
func (l *trinoListener) ExitShowRoles(ctx *grammar.ShowRolesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowRoleGrants is called when production showRoleGrants is entered.
func (l *trinoListener) EnterShowRoleGrants(ctx *grammar.ShowRoleGrantsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowRoleGrants is called when production showRoleGrants is exited.
func (l *trinoListener) ExitShowRoleGrants(ctx *grammar.ShowRoleGrantsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowFunctions is called when production showFunctions is entered.
func (l *trinoListener) EnterShowFunctions(ctx *grammar.ShowFunctionsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowFunctions is called when production showFunctions is exited.
func (l *trinoListener) ExitShowFunctions(ctx *grammar.ShowFunctionsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterShowSession is called when production showSession is entered.
func (l *trinoListener) EnterShowSession(ctx *grammar.ShowSessionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitShowSession is called when production showSession is exited.
func (l *trinoListener) ExitShowSession(ctx *grammar.ShowSessionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetSessionAuthorization is called when production setSessionAuthorization is entered.
func (l *trinoListener) EnterSetSessionAuthorization(ctx *grammar.SetSessionAuthorizationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetSessionAuthorization is called when production setSessionAuthorization is exited.
func (l *trinoListener) ExitSetSessionAuthorization(ctx *grammar.SetSessionAuthorizationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterResetSessionAuthorization is called when production resetSessionAuthorization is entered.
func (l *trinoListener) EnterResetSessionAuthorization(ctx *grammar.ResetSessionAuthorizationContext) {
}

// ExitResetSessionAuthorization is called when production resetSessionAuthorization is exited.
func (l *trinoListener) ExitResetSessionAuthorization(ctx *grammar.ResetSessionAuthorizationContext) {
}

// EnterSetSession is called when production setSession is entered.
func (l *trinoListener) EnterSetSession(ctx *grammar.SetSessionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetSession is called when production setSession is exited.
func (l *trinoListener) ExitSetSession(ctx *grammar.SetSessionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterResetSession is called when production resetSession is entered.
func (l *trinoListener) EnterResetSession(ctx *grammar.ResetSessionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitResetSession is called when production resetSession is exited.
func (l *trinoListener) ExitResetSession(ctx *grammar.ResetSessionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterStartTransaction is called when production startTransaction is entered.
func (l *trinoListener) EnterStartTransaction(ctx *grammar.StartTransactionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitStartTransaction is called when production startTransaction is exited.
func (l *trinoListener) ExitStartTransaction(ctx *grammar.StartTransactionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCommit is called when production commit is entered.
func (l *trinoListener) EnterCommit(ctx *grammar.CommitContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCommit is called when production commit is exited.
func (l *trinoListener) ExitCommit(ctx *grammar.CommitContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRollback is called when production rollback is entered.
func (l *trinoListener) EnterRollback(ctx *grammar.RollbackContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRollback is called when production rollback is exited.
func (l *trinoListener) ExitRollback(ctx *grammar.RollbackContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPrepare is called when production prepare is entered.
func (l *trinoListener) EnterPrepare(ctx *grammar.PrepareContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPrepare is called when production prepare is exited.
func (l *trinoListener) ExitPrepare(ctx *grammar.PrepareContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDeallocate is called when production deallocate is entered.
func (l *trinoListener) EnterDeallocate(ctx *grammar.DeallocateContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDeallocate is called when production deallocate is exited.
func (l *trinoListener) ExitDeallocate(ctx *grammar.DeallocateContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExecute is called when production execute is entered.
func (l *trinoListener) EnterExecute(ctx *grammar.ExecuteContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExecute is called when production execute is exited.
func (l *trinoListener) ExitExecute(ctx *grammar.ExecuteContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExecuteImmediate is called when production executeImmediate is entered.
func (l *trinoListener) EnterExecuteImmediate(ctx *grammar.ExecuteImmediateContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExecuteImmediate is called when production executeImmediate is exited.
func (l *trinoListener) ExitExecuteImmediate(ctx *grammar.ExecuteImmediateContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDescribeInput is called when production describeInput is entered.
func (l *trinoListener) EnterDescribeInput(ctx *grammar.DescribeInputContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDescribeInput is called when production describeInput is exited.
func (l *trinoListener) ExitDescribeInput(ctx *grammar.DescribeInputContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDescribeOutput is called when production describeOutput is entered.
func (l *trinoListener) EnterDescribeOutput(ctx *grammar.DescribeOutputContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDescribeOutput is called when production describeOutput is exited.
func (l *trinoListener) ExitDescribeOutput(ctx *grammar.DescribeOutputContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetPath is called when production setPath is entered.
func (l *trinoListener) EnterSetPath(ctx *grammar.SetPathContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetPath is called when production setPath is exited.
func (l *trinoListener) ExitSetPath(ctx *grammar.SetPathContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetTimeZone is called when production setTimeZone is entered.
func (l *trinoListener) EnterSetTimeZone(ctx *grammar.SetTimeZoneContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetTimeZone is called when production setTimeZone is exited.
func (l *trinoListener) ExitSetTimeZone(ctx *grammar.SetTimeZoneContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterWithFunction is called when production withFunction is entered.
func (l *trinoListener) EnterWithFunction(ctx *grammar.WithFunctionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitWithFunction is called when production withFunction is exited.
func (l *trinoListener) ExitWithFunction(ctx *grammar.WithFunctionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// // EnterQuery is called when production query is entered.
// func (l *trinoListener) EnterQuery(ctx *grammar.QueryContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// // ExitQuery is called when production query is exited.
// func (l *trinoListener) ExitQuery(ctx *grammar.QueryContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// EnterTableElement is called when production tableElement is entered.
func (l *trinoListener) EnterTableElement(ctx *grammar.TableElementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTableElement is called when production tableElement is exited.
func (l *trinoListener) ExitTableElement(ctx *grammar.TableElementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLikeClause is called when production likeClause is entered.
func (l *trinoListener) EnterLikeClause(ctx *grammar.LikeClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLikeClause is called when production likeClause is exited.
func (l *trinoListener) ExitLikeClause(ctx *grammar.LikeClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterProperties is called when production properties is entered.
func (l *trinoListener) EnterProperties(ctx *grammar.PropertiesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitProperties is called when production properties is exited.
func (l *trinoListener) ExitProperties(ctx *grammar.PropertiesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPropertyAssignments is called when production propertyAssignments is entered.
func (l *trinoListener) EnterPropertyAssignments(ctx *grammar.PropertyAssignmentsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPropertyAssignments is called when production propertyAssignments is exited.
func (l *trinoListener) ExitPropertyAssignments(ctx *grammar.PropertyAssignmentsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterProperty is called when production property is entered.
func (l *trinoListener) EnterProperty(ctx *grammar.PropertyContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitProperty is called when production property is exited.
func (l *trinoListener) ExitProperty(ctx *grammar.PropertyContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDefaultPropertyValue is called when production defaultPropertyValue is entered.
func (l *trinoListener) EnterDefaultPropertyValue(ctx *grammar.DefaultPropertyValueContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDefaultPropertyValue is called when production defaultPropertyValue is exited.
func (l *trinoListener) ExitDefaultPropertyValue(ctx *grammar.DefaultPropertyValueContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterNonDefaultPropertyValue is called when production nonDefaultPropertyValue is entered.
func (l *trinoListener) EnterNonDefaultPropertyValue(ctx *grammar.NonDefaultPropertyValueContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitNonDefaultPropertyValue is called when production nonDefaultPropertyValue is exited.
func (l *trinoListener) ExitNonDefaultPropertyValue(ctx *grammar.NonDefaultPropertyValueContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// // EnterQueryNoWith is called when production queryNoWith is entered.
// func (l *trinoListener) EnterQueryNoWith(ctx *grammar.QueryNoWithContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// // ExitQueryNoWith is called when production queryNoWith is exited.
// func (l *trinoListener) ExitQueryNoWith(ctx *grammar.QueryNoWithContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// EnterLimitRowCount is called when production limitRowCount is entered.
func (l *trinoListener) EnterLimitRowCount(ctx *grammar.LimitRowCountContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLimitRowCount is called when production limitRowCount is exited.
func (l *trinoListener) ExitLimitRowCount(ctx *grammar.LimitRowCountContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRowCount is called when production rowCount is entered.
func (l *trinoListener) EnterRowCount(ctx *grammar.RowCountContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRowCount is called when production rowCount is exited.
func (l *trinoListener) ExitRowCount(ctx *grammar.RowCountContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// // EnterQueryTermDefault is called when production queryTermDefault is entered.
// func (l *trinoListener) EnterQueryTermDefault(ctx *grammar.QueryTermDefaultContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// // ExitQueryTermDefault is called when production queryTermDefault is exited.
// func (l *trinoListener) ExitQueryTermDefault(ctx *grammar.QueryTermDefaultContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// EnterSetOperation is called when production setOperation is entered.
func (l *trinoListener) EnterSetOperation(ctx *grammar.SetOperationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetOperation is called when production setOperation is exited.
func (l *trinoListener) ExitSetOperation(ctx *grammar.SetOperationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// // EnterQueryPrimaryDefault is called when production queryPrimaryDefault is entered.
// func (l *trinoListener) EnterQueryPrimaryDefault(ctx *grammar.QueryPrimaryDefaultContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// // ExitQueryPrimaryDefault is called when production queryPrimaryDefault is exited.
// func (l *trinoListener) ExitQueryPrimaryDefault(ctx *grammar.QueryPrimaryDefaultContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// EnterTable is called when production table is entered.
func (l *trinoListener) EnterTable(ctx *grammar.TableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTable is called when production table is exited.
func (l *trinoListener) ExitTable(ctx *grammar.TableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterInlineTable is called when production inlineTable is entered.
func (l *trinoListener) EnterInlineTable(ctx *grammar.InlineTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitInlineTable is called when production inlineTable is exited.
func (l *trinoListener) ExitInlineTable(ctx *grammar.InlineTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSubquery is called when production subquery is entered.
func (l *trinoListener) EnterSubquery(ctx *grammar.SubqueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSubquery is called when production subquery is exited.
func (l *trinoListener) ExitSubquery(ctx *grammar.SubqueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSortItem is called when production sortItem is entered.
func (l *trinoListener) EnterSortItem(ctx *grammar.SortItemContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSortItem is called when production sortItem is exited.
func (l *trinoListener) ExitSortItem(ctx *grammar.SortItemContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterQuerySpecification is called when production querySpecification is entered.
// func (l *trinoListener) EnterQuerySpecification(ctx *grammar.QuerySpecificationContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// // ExitQuerySpecification is called when production querySpecification is exited.
// func (l *trinoListener) ExitQuerySpecification(ctx *grammar.QuerySpecificationContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// EnterGroupBy is called when production groupBy is entered.
func (l *trinoListener) EnterGroupBy(ctx *grammar.GroupByContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitGroupBy is called when production groupBy is exited.
func (l *trinoListener) ExitGroupBy(ctx *grammar.GroupByContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSingleGroupingSet is called when production singleGroupingSet is entered.
func (l *trinoListener) EnterSingleGroupingSet(ctx *grammar.SingleGroupingSetContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSingleGroupingSet is called when production singleGroupingSet is exited.
func (l *trinoListener) ExitSingleGroupingSet(ctx *grammar.SingleGroupingSetContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRollup is called when production rollup is entered.
func (l *trinoListener) EnterRollup(ctx *grammar.RollupContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRollup is called when production rollup is exited.
func (l *trinoListener) ExitRollup(ctx *grammar.RollupContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCube is called when production cube is entered.
func (l *trinoListener) EnterCube(ctx *grammar.CubeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCube is called when production cube is exited.
func (l *trinoListener) ExitCube(ctx *grammar.CubeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterMultipleGroupingSets is called when production multipleGroupingSets is entered.
func (l *trinoListener) EnterMultipleGroupingSets(ctx *grammar.MultipleGroupingSetsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitMultipleGroupingSets is called when production multipleGroupingSets is exited.
func (l *trinoListener) ExitMultipleGroupingSets(ctx *grammar.MultipleGroupingSetsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterGroupingSet is called when production groupingSet is entered.
func (l *trinoListener) EnterGroupingSet(ctx *grammar.GroupingSetContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitGroupingSet is called when production groupingSet is exited.
func (l *trinoListener) ExitGroupingSet(ctx *grammar.GroupingSetContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterWindowDefinition is called when production windowDefinition is entered.
func (l *trinoListener) EnterWindowDefinition(ctx *grammar.WindowDefinitionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitWindowDefinition is called when production windowDefinition is exited.
func (l *trinoListener) ExitWindowDefinition(ctx *grammar.WindowDefinitionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterWindowSpecification is called when production windowSpecification is entered.
func (l *trinoListener) EnterWindowSpecification(ctx *grammar.WindowSpecificationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitWindowSpecification is called when production windowSpecification is exited.
func (l *trinoListener) ExitWindowSpecification(ctx *grammar.WindowSpecificationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSetQuantifier is called when production setQuantifier is entered.
func (l *trinoListener) EnterSetQuantifier(ctx *grammar.SetQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSetQuantifier is called when production setQuantifier is exited.
func (l *trinoListener) ExitSetQuantifier(ctx *grammar.SetQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSampleType is called when production sampleType is entered.
func (l *trinoListener) EnterSampleType(ctx *grammar.SampleTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSampleType is called when production sampleType is exited.
func (l *trinoListener) ExitSampleType(ctx *grammar.SampleTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTrimsSpecification is called when production trimsSpecification is entered.
func (l *trinoListener) EnterTrimsSpecification(ctx *grammar.TrimsSpecificationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTrimsSpecification is called when production trimsSpecification is exited.
func (l *trinoListener) ExitTrimsSpecification(ctx *grammar.TrimsSpecificationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterListAggOverflowBehavior is called when production listAggOverflowBehavior is entered.
func (l *trinoListener) EnterListAggOverflowBehavior(ctx *grammar.ListAggOverflowBehaviorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitListAggOverflowBehavior is called when production listAggOverflowBehavior is exited.
func (l *trinoListener) ExitListAggOverflowBehavior(ctx *grammar.ListAggOverflowBehaviorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterListaggCountIndication is called when production listaggCountIndication is entered.
func (l *trinoListener) EnterListaggCountIndication(ctx *grammar.ListaggCountIndicationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitListaggCountIndication is called when production listaggCountIndication is exited.
func (l *trinoListener) ExitListaggCountIndication(ctx *grammar.ListaggCountIndicationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterMeasureDefinition is called when production measureDefinition is entered.
func (l *trinoListener) EnterMeasureDefinition(ctx *grammar.MeasureDefinitionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitMeasureDefinition is called when production measureDefinition is exited.
func (l *trinoListener) ExitMeasureDefinition(ctx *grammar.MeasureDefinitionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRowsPerMatch is called when production rowsPerMatch is entered.
func (l *trinoListener) EnterRowsPerMatch(ctx *grammar.RowsPerMatchContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRowsPerMatch is called when production rowsPerMatch is exited.
func (l *trinoListener) ExitRowsPerMatch(ctx *grammar.RowsPerMatchContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterEmptyMatchHandling is called when production emptyMatchHandling is entered.
func (l *trinoListener) EnterEmptyMatchHandling(ctx *grammar.EmptyMatchHandlingContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitEmptyMatchHandling is called when production emptyMatchHandling is exited.
func (l *trinoListener) ExitEmptyMatchHandling(ctx *grammar.EmptyMatchHandlingContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSkipTo is called when production skipTo is entered.
func (l *trinoListener) EnterSkipTo(ctx *grammar.SkipToContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSkipTo is called when production skipTo is exited.
func (l *trinoListener) ExitSkipTo(ctx *grammar.SkipToContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSubsetDefinition is called when production subsetDefinition is entered.
func (l *trinoListener) EnterSubsetDefinition(ctx *grammar.SubsetDefinitionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSubsetDefinition is called when production subsetDefinition is exited.
func (l *trinoListener) ExitSubsetDefinition(ctx *grammar.SubsetDefinitionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterVariableDefinition is called when production variableDefinition is entered.
func (l *trinoListener) EnterVariableDefinition(ctx *grammar.VariableDefinitionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitVariableDefinition is called when production variableDefinition is exited.
func (l *trinoListener) ExitVariableDefinition(ctx *grammar.VariableDefinitionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// // EnterAliasedRelation is called when production aliasedRelation is entered.
// func (l *trinoListener) EnterAliasedRelation(ctx *grammar.AliasedRelationContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// // ExitAliasedRelation is called when production aliasedRelation is exited.
// func (l *trinoListener) ExitAliasedRelation(ctx *grammar.AliasedRelationContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// // EnterColumnAliases is called when production columnAliases is entered.
// func (l *trinoListener) EnterColumnAliases(ctx *grammar.ColumnAliasesContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// // ExitColumnAliases is called when production columnAliases is exited.
// func (l *trinoListener) ExitColumnAliases(ctx *grammar.ColumnAliasesContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// ExitTableName is called when production tableName is exited.
// func (l *trinoListener) ExitTableName(ctx *grammar.TableNameContext) {
// 	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
// }

// EnterSubqueryRelation is called when production subqueryRelation is entered.
func (l *trinoListener) EnterSubqueryRelation(ctx *grammar.SubqueryRelationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSubqueryRelation is called when production subqueryRelation is exited.
func (l *trinoListener) ExitSubqueryRelation(ctx *grammar.SubqueryRelationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterUnnest is called when production unnest is entered.
func (l *trinoListener) EnterUnnest(ctx *grammar.UnnestContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitUnnest is called when production unnest is exited.
func (l *trinoListener) ExitUnnest(ctx *grammar.UnnestContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLateral is called when production lateral is entered.
func (l *trinoListener) EnterLateral(ctx *grammar.LateralContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLateral is called when production lateral is exited.
func (l *trinoListener) ExitLateral(ctx *grammar.LateralContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTableFunctionInvocation is called when production tableFunctionInvocation is entered.
func (l *trinoListener) EnterTableFunctionInvocation(ctx *grammar.TableFunctionInvocationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTableFunctionInvocation is called when production tableFunctionInvocation is exited.
func (l *trinoListener) ExitTableFunctionInvocation(ctx *grammar.TableFunctionInvocationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterParenthesizedRelation is called when production parenthesizedRelation is entered.
func (l *trinoListener) EnterParenthesizedRelation(ctx *grammar.ParenthesizedRelationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitParenthesizedRelation is called when production parenthesizedRelation is exited.
func (l *trinoListener) ExitParenthesizedRelation(ctx *grammar.ParenthesizedRelationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTableFunctionCall is called when production tableFunctionCall is entered.
func (l *trinoListener) EnterTableFunctionCall(ctx *grammar.TableFunctionCallContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTableFunctionCall is called when production tableFunctionCall is exited.
func (l *trinoListener) ExitTableFunctionCall(ctx *grammar.TableFunctionCallContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTableFunctionArgument is called when production tableFunctionArgument is entered.
func (l *trinoListener) EnterTableFunctionArgument(ctx *grammar.TableFunctionArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTableFunctionArgument is called when production tableFunctionArgument is exited.
func (l *trinoListener) ExitTableFunctionArgument(ctx *grammar.TableFunctionArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTableArgument is called when production tableArgument is entered.
func (l *trinoListener) EnterTableArgument(ctx *grammar.TableArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTableArgument is called when production tableArgument is exited.
func (l *trinoListener) ExitTableArgument(ctx *grammar.TableArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTableArgumentTable is called when production tableArgumentTable is entered.
func (l *trinoListener) EnterTableArgumentTable(ctx *grammar.TableArgumentTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTableArgumentTable is called when production tableArgumentTable is exited.
func (l *trinoListener) ExitTableArgumentTable(ctx *grammar.TableArgumentTableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTableArgumentQuery is called when production tableArgumentQuery is entered.
func (l *trinoListener) EnterTableArgumentQuery(ctx *grammar.TableArgumentQueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTableArgumentQuery is called when production tableArgumentQuery is exited.
func (l *trinoListener) ExitTableArgumentQuery(ctx *grammar.TableArgumentQueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDescriptorArgument is called when production descriptorArgument is entered.
func (l *trinoListener) EnterDescriptorArgument(ctx *grammar.DescriptorArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDescriptorArgument is called when production descriptorArgument is exited.
func (l *trinoListener) ExitDescriptorArgument(ctx *grammar.DescriptorArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDescriptorField is called when production descriptorField is entered.
func (l *trinoListener) EnterDescriptorField(ctx *grammar.DescriptorFieldContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDescriptorField is called when production descriptorField is exited.
func (l *trinoListener) ExitDescriptorField(ctx *grammar.DescriptorFieldContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCopartitionTables is called when production copartitionTables is entered.
func (l *trinoListener) EnterCopartitionTables(ctx *grammar.CopartitionTablesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCopartitionTables is called when production copartitionTables is exited.
func (l *trinoListener) ExitCopartitionTables(ctx *grammar.CopartitionTablesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLogicalNot is called when production logicalNot is entered.
func (l *trinoListener) EnterLogicalNot(ctx *grammar.LogicalNotContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLogicalNot is called when production logicalNot is exited.
func (l *trinoListener) ExitLogicalNot(ctx *grammar.LogicalNotContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterInSubquery is called when production inSubquery is entered.
func (l *trinoListener) EnterInSubquery(ctx *grammar.InSubqueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitInSubquery is called when production inSubquery is exited.
func (l *trinoListener) ExitInSubquery(ctx *grammar.InSubqueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDistinctFrom is called when production distinctFrom is entered.
func (l *trinoListener) EnterDistinctFrom(ctx *grammar.DistinctFromContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDistinctFrom is called when production distinctFrom is exited.
func (l *trinoListener) ExitDistinctFrom(ctx *grammar.DistinctFromContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterConcatenation is called when production concatenation is entered.
func (l *trinoListener) EnterConcatenation(ctx *grammar.ConcatenationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitConcatenation is called when production concatenation is exited.
func (l *trinoListener) ExitConcatenation(ctx *grammar.ConcatenationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterAtTimeZone is called when production atTimeZone is entered.
func (l *trinoListener) EnterAtTimeZone(ctx *grammar.AtTimeZoneContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitAtTimeZone is called when production atTimeZone is exited.
func (l *trinoListener) ExitAtTimeZone(ctx *grammar.AtTimeZoneContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonValue is called when production jsonValue is entered.
func (l *trinoListener) EnterJsonValue(ctx *grammar.JsonValueContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonValue is called when production jsonValue is exited.
func (l *trinoListener) ExitJsonValue(ctx *grammar.JsonValueContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSpecialDateTimeFunction is called when production specialDateTimeFunction is entered.
func (l *trinoListener) EnterSpecialDateTimeFunction(ctx *grammar.SpecialDateTimeFunctionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSpecialDateTimeFunction is called when production specialDateTimeFunction is exited.
func (l *trinoListener) ExitSpecialDateTimeFunction(ctx *grammar.SpecialDateTimeFunctionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLambda is called when production lambda is entered.
func (l *trinoListener) EnterLambda(ctx *grammar.LambdaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLambda is called when production lambda is exited.
func (l *trinoListener) ExitLambda(ctx *grammar.LambdaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterParenthesizedExpression is called when production parenthesizedExpression is entered.
func (l *trinoListener) EnterParenthesizedExpression(ctx *grammar.ParenthesizedExpressionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitParenthesizedExpression is called when production parenthesizedExpression is exited.
func (l *trinoListener) ExitParenthesizedExpression(ctx *grammar.ParenthesizedExpressionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterParameter is called when production parameter is entered.
func (l *trinoListener) EnterParameter(ctx *grammar.ParameterContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitParameter is called when production parameter is exited.
func (l *trinoListener) ExitParameter(ctx *grammar.ParameterContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterNormalize is called when production normalize is entered.
func (l *trinoListener) EnterNormalize(ctx *grammar.NormalizeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitNormalize is called when production normalize is exited.
func (l *trinoListener) ExitNormalize(ctx *grammar.NormalizeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonObject is called when production jsonObject is entered.
func (l *trinoListener) EnterJsonObject(ctx *grammar.JsonObjectContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonObject is called when production jsonObject is exited.
func (l *trinoListener) ExitJsonObject(ctx *grammar.JsonObjectContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonArray is called when production jsonArray is entered.
func (l *trinoListener) EnterJsonArray(ctx *grammar.JsonArrayContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonArray is called when production jsonArray is exited.
func (l *trinoListener) ExitJsonArray(ctx *grammar.JsonArrayContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSimpleCase is called when production simpleCase is entered.
func (l *trinoListener) EnterSimpleCase(ctx *grammar.SimpleCaseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSimpleCase is called when production simpleCase is exited.
func (l *trinoListener) ExitSimpleCase(ctx *grammar.SimpleCaseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRowConstructor is called when production rowConstructor is entered.
func (l *trinoListener) EnterRowConstructor(ctx *grammar.RowConstructorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRowConstructor is called when production rowConstructor is exited.
func (l *trinoListener) ExitRowConstructor(ctx *grammar.RowConstructorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSubscript is called when production subscript is entered.
func (l *trinoListener) EnterSubscript(ctx *grammar.SubscriptContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSubscript is called when production subscript is exited.
func (l *trinoListener) ExitSubscript(ctx *grammar.SubscriptContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonExists is called when production jsonExists is entered.
func (l *trinoListener) EnterJsonExists(ctx *grammar.JsonExistsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonExists is called when production jsonExists is exited.
func (l *trinoListener) ExitJsonExists(ctx *grammar.JsonExistsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCurrentPath is called when production currentPath is entered.
func (l *trinoListener) EnterCurrentPath(ctx *grammar.CurrentPathContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCurrentPath is called when production currentPath is exited.
func (l *trinoListener) ExitCurrentPath(ctx *grammar.CurrentPathContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSubqueryExpression is called when production subqueryExpression is entered.
func (l *trinoListener) EnterSubqueryExpression(ctx *grammar.SubqueryExpressionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSubqueryExpression is called when production subqueryExpression is exited.
func (l *trinoListener) ExitSubqueryExpression(ctx *grammar.SubqueryExpressionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCurrentUser is called when production currentUser is entered.
func (l *trinoListener) EnterCurrentUser(ctx *grammar.CurrentUserContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCurrentUser is called when production currentUser is exited.
func (l *trinoListener) ExitCurrentUser(ctx *grammar.CurrentUserContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonQuery is called when production jsonQuery is entered.
func (l *trinoListener) EnterJsonQuery(ctx *grammar.JsonQueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonQuery is called when production jsonQuery is exited.
func (l *trinoListener) ExitJsonQuery(ctx *grammar.JsonQueryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterMeasure is called when production measure is entered.
func (l *trinoListener) EnterMeasure(ctx *grammar.MeasureContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitMeasure is called when production measure is exited.
func (l *trinoListener) ExitMeasure(ctx *grammar.MeasureContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExtract is called when production extract is entered.
func (l *trinoListener) EnterExtract(ctx *grammar.ExtractContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExtract is called when production extract is exited.
func (l *trinoListener) ExitExtract(ctx *grammar.ExtractContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterFunctionCall is called when production functionCall is entered.
func (l *trinoListener) EnterFunctionCall(ctx *grammar.FunctionCallContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitFunctionCall is called when production functionCall is exited.
func (l *trinoListener) ExitFunctionCall(ctx *grammar.FunctionCallContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCurrentSchema is called when production currentSchema is entered.
func (l *trinoListener) EnterCurrentSchema(ctx *grammar.CurrentSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCurrentSchema is called when production currentSchema is exited.
func (l *trinoListener) ExitCurrentSchema(ctx *grammar.CurrentSchemaContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExists is called when production exists is entered.
func (l *trinoListener) EnterExists(ctx *grammar.ExistsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExists is called when production exists is exited.
func (l *trinoListener) ExitExists(ctx *grammar.ExistsContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPosition is called when production position is entered.
func (l *trinoListener) EnterPosition(ctx *grammar.PositionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPosition is called when production position is exited.
func (l *trinoListener) ExitPosition(ctx *grammar.PositionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterListagg is called when production listagg is entered.
func (l *trinoListener) EnterListagg(ctx *grammar.ListaggContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitListagg is called when production listagg is exited.
func (l *trinoListener) ExitListagg(ctx *grammar.ListaggContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSearchedCase is called when production searchedCase is entered.
func (l *trinoListener) EnterSearchedCase(ctx *grammar.SearchedCaseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSearchedCase is called when production searchedCase is exited.
func (l *trinoListener) ExitSearchedCase(ctx *grammar.SearchedCaseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCurrentCatalog is called when production currentCatalog is entered.
func (l *trinoListener) EnterCurrentCatalog(ctx *grammar.CurrentCatalogContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCurrentCatalog is called when production currentCatalog is exited.
func (l *trinoListener) ExitCurrentCatalog(ctx *grammar.CurrentCatalogContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterGroupingOperation is called when production groupingOperation is entered.
func (l *trinoListener) EnterGroupingOperation(ctx *grammar.GroupingOperationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitGroupingOperation is called when production groupingOperation is exited.
func (l *trinoListener) ExitGroupingOperation(ctx *grammar.GroupingOperationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonPathInvocation is called when production jsonPathInvocation is entered.
func (l *trinoListener) EnterJsonPathInvocation(ctx *grammar.JsonPathInvocationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonPathInvocation is called when production jsonPathInvocation is exited.
func (l *trinoListener) ExitJsonPathInvocation(ctx *grammar.JsonPathInvocationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonValueExpression is called when production jsonValueExpression is entered.
func (l *trinoListener) EnterJsonValueExpression(ctx *grammar.JsonValueExpressionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonValueExpression is called when production jsonValueExpression is exited.
func (l *trinoListener) ExitJsonValueExpression(ctx *grammar.JsonValueExpressionContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonRepresentation is called when production jsonRepresentation is entered.
func (l *trinoListener) EnterJsonRepresentation(ctx *grammar.JsonRepresentationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonRepresentation is called when production jsonRepresentation is exited.
func (l *trinoListener) ExitJsonRepresentation(ctx *grammar.JsonRepresentationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonArgument is called when production jsonArgument is entered.
func (l *trinoListener) EnterJsonArgument(ctx *grammar.JsonArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonArgument is called when production jsonArgument is exited.
func (l *trinoListener) ExitJsonArgument(ctx *grammar.JsonArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonExistsErrorBehavior is called when production jsonExistsErrorBehavior is entered.
func (l *trinoListener) EnterJsonExistsErrorBehavior(ctx *grammar.JsonExistsErrorBehaviorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonExistsErrorBehavior is called when production jsonExistsErrorBehavior is exited.
func (l *trinoListener) ExitJsonExistsErrorBehavior(ctx *grammar.JsonExistsErrorBehaviorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonValueBehavior is called when production jsonValueBehavior is entered.
func (l *trinoListener) EnterJsonValueBehavior(ctx *grammar.JsonValueBehaviorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonValueBehavior is called when production jsonValueBehavior is exited.
func (l *trinoListener) ExitJsonValueBehavior(ctx *grammar.JsonValueBehaviorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonQueryWrapperBehavior is called when production jsonQueryWrapperBehavior is entered.
func (l *trinoListener) EnterJsonQueryWrapperBehavior(ctx *grammar.JsonQueryWrapperBehaviorContext) {
}

// ExitJsonQueryWrapperBehavior is called when production jsonQueryWrapperBehavior is exited.
func (l *trinoListener) ExitJsonQueryWrapperBehavior(ctx *grammar.JsonQueryWrapperBehaviorContext) {
}

// EnterJsonQueryBehavior is called when production jsonQueryBehavior is entered.
func (l *trinoListener) EnterJsonQueryBehavior(ctx *grammar.JsonQueryBehaviorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonQueryBehavior is called when production jsonQueryBehavior is exited.
func (l *trinoListener) ExitJsonQueryBehavior(ctx *grammar.JsonQueryBehaviorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterJsonObjectMember is called when production jsonObjectMember is entered.
func (l *trinoListener) EnterJsonObjectMember(ctx *grammar.JsonObjectMemberContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitJsonObjectMember is called when production jsonObjectMember is exited.
func (l *trinoListener) ExitJsonObjectMember(ctx *grammar.JsonObjectMemberContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterProcessingMode is called when production processingMode is entered.
func (l *trinoListener) EnterProcessingMode(ctx *grammar.ProcessingModeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitProcessingMode is called when production processingMode is exited.
func (l *trinoListener) ExitProcessingMode(ctx *grammar.ProcessingModeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterIntervalField is called when production intervalField is entered.
func (l *trinoListener) EnterIntervalField(ctx *grammar.IntervalFieldContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitIntervalField is called when production intervalField is exited.
func (l *trinoListener) ExitIntervalField(ctx *grammar.IntervalFieldContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterNormalForm is called when production normalForm is entered.
func (l *trinoListener) EnterNormalForm(ctx *grammar.NormalFormContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitNormalForm is called when production normalForm is exited.
func (l *trinoListener) ExitNormalForm(ctx *grammar.NormalFormContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDoublePrecisionType is called when production doublePrecisionType is entered.
func (l *trinoListener) EnterDoublePrecisionType(ctx *grammar.DoublePrecisionTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitDoublePrecisionType is called when production doublePrecisionType is exited.
func (l *trinoListener) ExitDoublePrecisionType(ctx *grammar.DoublePrecisionTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLegacyArrayType is called when production legacyArrayType is entered.
func (l *trinoListener) EnterLegacyArrayType(ctx *grammar.LegacyArrayTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLegacyArrayType is called when production legacyArrayType is exited.
func (l *trinoListener) ExitLegacyArrayType(ctx *grammar.LegacyArrayTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLegacyMapType is called when production legacyMapType is entered.
func (l *trinoListener) EnterLegacyMapType(ctx *grammar.LegacyMapTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLegacyMapType is called when production legacyMapType is exited.
func (l *trinoListener) ExitLegacyMapType(ctx *grammar.LegacyMapTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRowField is called when production rowField is entered.
func (l *trinoListener) EnterRowField(ctx *grammar.RowFieldContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRowField is called when production rowField is exited.
func (l *trinoListener) ExitRowField(ctx *grammar.RowFieldContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTypeParameter is called when production typeParameter is entered.
func (l *trinoListener) EnterTypeParameter(ctx *grammar.TypeParameterContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTypeParameter is called when production typeParameter is exited.
func (l *trinoListener) ExitTypeParameter(ctx *grammar.TypeParameterContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterWhenClause is called when production whenClause is entered.
func (l *trinoListener) EnterWhenClause(ctx *grammar.WhenClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitWhenClause is called when production whenClause is exited.
func (l *trinoListener) ExitWhenClause(ctx *grammar.WhenClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterFilter is called when production filter is entered.
func (l *trinoListener) EnterFilter(ctx *grammar.FilterContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitFilter is called when production filter is exited.
func (l *trinoListener) ExitFilter(ctx *grammar.FilterContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterMergeUpdate is called when production mergeUpdate is entered.
func (l *trinoListener) EnterMergeUpdate(ctx *grammar.MergeUpdateContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitMergeUpdate is called when production mergeUpdate is exited.
func (l *trinoListener) ExitMergeUpdate(ctx *grammar.MergeUpdateContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterMergeDelete is called when production mergeDelete is entered.
func (l *trinoListener) EnterMergeDelete(ctx *grammar.MergeDeleteContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitMergeDelete is called when production mergeDelete is exited.
func (l *trinoListener) ExitMergeDelete(ctx *grammar.MergeDeleteContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterOver is called when production over is entered.
func (l *trinoListener) EnterOver(ctx *grammar.OverContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitOver is called when production over is exited.
func (l *trinoListener) ExitOver(ctx *grammar.OverContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterWindowFrame is called when production windowFrame is entered.
func (l *trinoListener) EnterWindowFrame(ctx *grammar.WindowFrameContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitWindowFrame is called when production windowFrame is exited.
func (l *trinoListener) ExitWindowFrame(ctx *grammar.WindowFrameContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterFrameExtent is called when production frameExtent is entered.
func (l *trinoListener) EnterFrameExtent(ctx *grammar.FrameExtentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitFrameExtent is called when production frameExtent is exited.
func (l *trinoListener) ExitFrameExtent(ctx *grammar.FrameExtentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterUnboundedFrame is called when production unboundedFrame is entered.
func (l *trinoListener) EnterUnboundedFrame(ctx *grammar.UnboundedFrameContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitUnboundedFrame is called when production unboundedFrame is exited.
func (l *trinoListener) ExitUnboundedFrame(ctx *grammar.UnboundedFrameContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCurrentRowBound is called when production currentRowBound is entered.
func (l *trinoListener) EnterCurrentRowBound(ctx *grammar.CurrentRowBoundContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCurrentRowBound is called when production currentRowBound is exited.
func (l *trinoListener) ExitCurrentRowBound(ctx *grammar.CurrentRowBoundContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterBoundedFrame is called when production boundedFrame is entered.
func (l *trinoListener) EnterBoundedFrame(ctx *grammar.BoundedFrameContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitBoundedFrame is called when production boundedFrame is exited.
func (l *trinoListener) ExitBoundedFrame(ctx *grammar.BoundedFrameContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterQuantifiedPrimary is called when production quantifiedPrimary is entered.
func (l *trinoListener) EnterQuantifiedPrimary(ctx *grammar.QuantifiedPrimaryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitQuantifiedPrimary is called when production quantifiedPrimary is exited.
func (l *trinoListener) ExitQuantifiedPrimary(ctx *grammar.QuantifiedPrimaryContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPatternConcatenation is called when production patternConcatenation is entered.
func (l *trinoListener) EnterPatternConcatenation(ctx *grammar.PatternConcatenationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPatternConcatenation is called when production patternConcatenation is exited.
func (l *trinoListener) ExitPatternConcatenation(ctx *grammar.PatternConcatenationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPatternAlternation is called when production patternAlternation is entered.
func (l *trinoListener) EnterPatternAlternation(ctx *grammar.PatternAlternationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPatternAlternation is called when production patternAlternation is exited.
func (l *trinoListener) ExitPatternAlternation(ctx *grammar.PatternAlternationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPatternVariable is called when production patternVariable is entered.
func (l *trinoListener) EnterPatternVariable(ctx *grammar.PatternVariableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPatternVariable is called when production patternVariable is exited.
func (l *trinoListener) ExitPatternVariable(ctx *grammar.PatternVariableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterEmptyPattern is called when production emptyPattern is entered.
func (l *trinoListener) EnterEmptyPattern(ctx *grammar.EmptyPatternContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitEmptyPattern is called when production emptyPattern is exited.
func (l *trinoListener) ExitEmptyPattern(ctx *grammar.EmptyPatternContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPatternPermutation is called when production patternPermutation is entered.
func (l *trinoListener) EnterPatternPermutation(ctx *grammar.PatternPermutationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPatternPermutation is called when production patternPermutation is exited.
func (l *trinoListener) ExitPatternPermutation(ctx *grammar.PatternPermutationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterGroupedPattern is called when production groupedPattern is entered.
func (l *trinoListener) EnterGroupedPattern(ctx *grammar.GroupedPatternContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitGroupedPattern is called when production groupedPattern is exited.
func (l *trinoListener) ExitGroupedPattern(ctx *grammar.GroupedPatternContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPartitionStartAnchor is called when production partitionStartAnchor is entered.
func (l *trinoListener) EnterPartitionStartAnchor(ctx *grammar.PartitionStartAnchorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPartitionStartAnchor is called when production partitionStartAnchor is exited.
func (l *trinoListener) ExitPartitionStartAnchor(ctx *grammar.PartitionStartAnchorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPartitionEndAnchor is called when production partitionEndAnchor is entered.
func (l *trinoListener) EnterPartitionEndAnchor(ctx *grammar.PartitionEndAnchorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPartitionEndAnchor is called when production partitionEndAnchor is exited.
func (l *trinoListener) ExitPartitionEndAnchor(ctx *grammar.PartitionEndAnchorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExcludedPattern is called when production excludedPattern is entered.
func (l *trinoListener) EnterExcludedPattern(ctx *grammar.ExcludedPatternContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExcludedPattern is called when production excludedPattern is exited.
func (l *trinoListener) ExitExcludedPattern(ctx *grammar.ExcludedPatternContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterZeroOrMoreQuantifier is called when production zeroOrMoreQuantifier is entered.
func (l *trinoListener) EnterZeroOrMoreQuantifier(ctx *grammar.ZeroOrMoreQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitZeroOrMoreQuantifier is called when production zeroOrMoreQuantifier is exited.
func (l *trinoListener) ExitZeroOrMoreQuantifier(ctx *grammar.ZeroOrMoreQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterOneOrMoreQuantifier is called when production oneOrMoreQuantifier is entered.
func (l *trinoListener) EnterOneOrMoreQuantifier(ctx *grammar.OneOrMoreQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitOneOrMoreQuantifier is called when production oneOrMoreQuantifier is exited.
func (l *trinoListener) ExitOneOrMoreQuantifier(ctx *grammar.OneOrMoreQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterZeroOrOneQuantifier is called when production zeroOrOneQuantifier is entered.
func (l *trinoListener) EnterZeroOrOneQuantifier(ctx *grammar.ZeroOrOneQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitZeroOrOneQuantifier is called when production zeroOrOneQuantifier is exited.
func (l *trinoListener) ExitZeroOrOneQuantifier(ctx *grammar.ZeroOrOneQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRangeQuantifier is called when production rangeQuantifier is entered.
func (l *trinoListener) EnterRangeQuantifier(ctx *grammar.RangeQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRangeQuantifier is called when production rangeQuantifier is exited.
func (l *trinoListener) ExitRangeQuantifier(ctx *grammar.RangeQuantifierContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterUpdateAssignment is called when production updateAssignment is entered.
func (l *trinoListener) EnterUpdateAssignment(ctx *grammar.UpdateAssignmentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitUpdateAssignment is called when production updateAssignment is exited.
func (l *trinoListener) ExitUpdateAssignment(ctx *grammar.UpdateAssignmentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExplainFormat is called when production explainFormat is entered.
func (l *trinoListener) EnterExplainFormat(ctx *grammar.ExplainFormatContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExplainFormat is called when production explainFormat is exited.
func (l *trinoListener) ExitExplainFormat(ctx *grammar.ExplainFormatContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterExplainType is called when production explainType is entered.
func (l *trinoListener) EnterExplainType(ctx *grammar.ExplainTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitExplainType is called when production explainType is exited.
func (l *trinoListener) ExitExplainType(ctx *grammar.ExplainTypeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterIsolationLevel is called when production isolationLevel is entered.
func (l *trinoListener) EnterIsolationLevel(ctx *grammar.IsolationLevelContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitIsolationLevel is called when production isolationLevel is exited.
func (l *trinoListener) ExitIsolationLevel(ctx *grammar.IsolationLevelContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterTransactionAccessMode is called when production transactionAccessMode is entered.
func (l *trinoListener) EnterTransactionAccessMode(ctx *grammar.TransactionAccessModeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitTransactionAccessMode is called when production transactionAccessMode is exited.
func (l *trinoListener) ExitTransactionAccessMode(ctx *grammar.TransactionAccessModeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterReadUncommitted is called when production readUncommitted is entered.
func (l *trinoListener) EnterReadUncommitted(ctx *grammar.ReadUncommittedContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitReadUncommitted is called when production readUncommitted is exited.
func (l *trinoListener) ExitReadUncommitted(ctx *grammar.ReadUncommittedContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterReadCommitted is called when production readCommitted is entered.
func (l *trinoListener) EnterReadCommitted(ctx *grammar.ReadCommittedContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitReadCommitted is called when production readCommitted is exited.
func (l *trinoListener) ExitReadCommitted(ctx *grammar.ReadCommittedContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRepeatableRead is called when production repeatableRead is entered.
func (l *trinoListener) EnterRepeatableRead(ctx *grammar.RepeatableReadContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRepeatableRead is called when production repeatableRead is exited.
func (l *trinoListener) ExitRepeatableRead(ctx *grammar.RepeatableReadContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSerializable is called when production serializable is entered.
func (l *trinoListener) EnterSerializable(ctx *grammar.SerializableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSerializable is called when production serializable is exited.
func (l *trinoListener) ExitSerializable(ctx *grammar.SerializableContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPositionalArgument is called when production positionalArgument is entered.
func (l *trinoListener) EnterPositionalArgument(ctx *grammar.PositionalArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPositionalArgument is called when production positionalArgument is exited.
func (l *trinoListener) ExitPositionalArgument(ctx *grammar.PositionalArgumentContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPathSpecification is called when production pathSpecification is entered.
func (l *trinoListener) EnterPathSpecification(ctx *grammar.PathSpecificationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPathSpecification is called when production pathSpecification is exited.
func (l *trinoListener) ExitPathSpecification(ctx *grammar.PathSpecificationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterFunctionSpecification is called when production functionSpecification is entered.
func (l *trinoListener) EnterFunctionSpecification(ctx *grammar.FunctionSpecificationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitFunctionSpecification is called when production functionSpecification is exited.
func (l *trinoListener) ExitFunctionSpecification(ctx *grammar.FunctionSpecificationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterFunctionDeclaration is called when production functionDeclaration is entered.
func (l *trinoListener) EnterFunctionDeclaration(ctx *grammar.FunctionDeclarationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitFunctionDeclaration is called when production functionDeclaration is exited.
func (l *trinoListener) ExitFunctionDeclaration(ctx *grammar.FunctionDeclarationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterParameterDeclaration is called when production parameterDeclaration is entered.
func (l *trinoListener) EnterParameterDeclaration(ctx *grammar.ParameterDeclarationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitParameterDeclaration is called when production parameterDeclaration is exited.
func (l *trinoListener) ExitParameterDeclaration(ctx *grammar.ParameterDeclarationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterReturnsClause is called when production returnsClause is entered.
func (l *trinoListener) EnterReturnsClause(ctx *grammar.ReturnsClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitReturnsClause is called when production returnsClause is exited.
func (l *trinoListener) ExitReturnsClause(ctx *grammar.ReturnsClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLanguageCharacteristic is called when production languageCharacteristic is entered.
func (l *trinoListener) EnterLanguageCharacteristic(ctx *grammar.LanguageCharacteristicContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLanguageCharacteristic is called when production languageCharacteristic is exited.
func (l *trinoListener) ExitLanguageCharacteristic(ctx *grammar.LanguageCharacteristicContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterDeterministicCharacteristic is called when production deterministicCharacteristic is entered.
func (l *trinoListener) EnterDeterministicCharacteristic(ctx *grammar.DeterministicCharacteristicContext) {
}

// ExitDeterministicCharacteristic is called when production deterministicCharacteristic is exited.
func (l *trinoListener) ExitDeterministicCharacteristic(ctx *grammar.DeterministicCharacteristicContext) {
}

// EnterReturnsNullOnNullInputCharacteristic is called when production returnsNullOnNullInputCharacteristic is entered.
func (l *trinoListener) EnterReturnsNullOnNullInputCharacteristic(ctx *grammar.ReturnsNullOnNullInputCharacteristicContext) {
}

// ExitReturnsNullOnNullInputCharacteristic is called when production returnsNullOnNullInputCharacteristic is exited.
func (l *trinoListener) ExitReturnsNullOnNullInputCharacteristic(ctx *grammar.ReturnsNullOnNullInputCharacteristicContext) {
}

// EnterCalledOnNullInputCharacteristic is called when production calledOnNullInputCharacteristic is entered.
func (l *trinoListener) EnterCalledOnNullInputCharacteristic(ctx *grammar.CalledOnNullInputCharacteristicContext) {
}

// ExitCalledOnNullInputCharacteristic is called when production calledOnNullInputCharacteristic is exited.
func (l *trinoListener) ExitCalledOnNullInputCharacteristic(ctx *grammar.CalledOnNullInputCharacteristicContext) {
}

// EnterSecurityCharacteristic is called when production securityCharacteristic is entered.
func (l *trinoListener) EnterSecurityCharacteristic(ctx *grammar.SecurityCharacteristicContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSecurityCharacteristic is called when production securityCharacteristic is exited.
func (l *trinoListener) ExitSecurityCharacteristic(ctx *grammar.SecurityCharacteristicContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCommentCharacteristic is called when production commentCharacteristic is entered.
func (l *trinoListener) EnterCommentCharacteristic(ctx *grammar.CommentCharacteristicContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCommentCharacteristic is called when production commentCharacteristic is exited.
func (l *trinoListener) ExitCommentCharacteristic(ctx *grammar.CommentCharacteristicContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterReturnStatement is called when production returnStatement is entered.
func (l *trinoListener) EnterReturnStatement(ctx *grammar.ReturnStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitReturnStatement is called when production returnStatement is exited.
func (l *trinoListener) ExitReturnStatement(ctx *grammar.ReturnStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterAssignmentStatement is called when production assignmentStatement is entered.
func (l *trinoListener) EnterAssignmentStatement(ctx *grammar.AssignmentStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitAssignmentStatement is called when production assignmentStatement is exited.
func (l *trinoListener) ExitAssignmentStatement(ctx *grammar.AssignmentStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSimpleCaseStatement is called when production simpleCaseStatement is entered.
func (l *trinoListener) EnterSimpleCaseStatement(ctx *grammar.SimpleCaseStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSimpleCaseStatement is called when production simpleCaseStatement is exited.
func (l *trinoListener) ExitSimpleCaseStatement(ctx *grammar.SimpleCaseStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSearchedCaseStatement is called when production searchedCaseStatement is entered.
func (l *trinoListener) EnterSearchedCaseStatement(ctx *grammar.SearchedCaseStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSearchedCaseStatement is called when production searchedCaseStatement is exited.
func (l *trinoListener) ExitSearchedCaseStatement(ctx *grammar.SearchedCaseStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterIfStatement is called when production ifStatement is entered.
func (l *trinoListener) EnterIfStatement(ctx *grammar.IfStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitIfStatement is called when production ifStatement is exited.
func (l *trinoListener) ExitIfStatement(ctx *grammar.IfStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterIterateStatement is called when production iterateStatement is entered.
func (l *trinoListener) EnterIterateStatement(ctx *grammar.IterateStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitIterateStatement is called when production iterateStatement is exited.
func (l *trinoListener) ExitIterateStatement(ctx *grammar.IterateStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLeaveStatement is called when production leaveStatement is entered.
func (l *trinoListener) EnterLeaveStatement(ctx *grammar.LeaveStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLeaveStatement is called when production leaveStatement is exited.
func (l *trinoListener) ExitLeaveStatement(ctx *grammar.LeaveStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCompoundStatement is called when production compoundStatement is entered.
func (l *trinoListener) EnterCompoundStatement(ctx *grammar.CompoundStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCompoundStatement is called when production compoundStatement is exited.
func (l *trinoListener) ExitCompoundStatement(ctx *grammar.CompoundStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterLoopStatement is called when production loopStatement is entered.
func (l *trinoListener) EnterLoopStatement(ctx *grammar.LoopStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitLoopStatement is called when production loopStatement is exited.
func (l *trinoListener) ExitLoopStatement(ctx *grammar.LoopStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterWhileStatement is called when production whileStatement is entered.
func (l *trinoListener) EnterWhileStatement(ctx *grammar.WhileStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitWhileStatement is called when production whileStatement is exited.
func (l *trinoListener) ExitWhileStatement(ctx *grammar.WhileStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRepeatStatement is called when production repeatStatement is entered.
func (l *trinoListener) EnterRepeatStatement(ctx *grammar.RepeatStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRepeatStatement is called when production repeatStatement is exited.
func (l *trinoListener) ExitRepeatStatement(ctx *grammar.RepeatStatementContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCaseStatementWhenClause is called when production caseStatementWhenClause is entered.
func (l *trinoListener) EnterCaseStatementWhenClause(ctx *grammar.CaseStatementWhenClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCaseStatementWhenClause is called when production caseStatementWhenClause is exited.
func (l *trinoListener) ExitCaseStatementWhenClause(ctx *grammar.CaseStatementWhenClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterElseIfClause is called when production elseIfClause is entered.
func (l *trinoListener) EnterElseIfClause(ctx *grammar.ElseIfClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitElseIfClause is called when production elseIfClause is exited.
func (l *trinoListener) ExitElseIfClause(ctx *grammar.ElseIfClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterElseClause is called when production elseClause is entered.
func (l *trinoListener) EnterElseClause(ctx *grammar.ElseClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitElseClause is called when production elseClause is exited.
func (l *trinoListener) ExitElseClause(ctx *grammar.ElseClauseContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterVariableDeclaration is called when production variableDeclaration is entered.
func (l *trinoListener) EnterVariableDeclaration(ctx *grammar.VariableDeclarationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitVariableDeclaration is called when production variableDeclaration is exited.
func (l *trinoListener) ExitVariableDeclaration(ctx *grammar.VariableDeclarationContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSqlStatementList is called when production sqlStatementList is entered.
func (l *trinoListener) EnterSqlStatementList(ctx *grammar.SqlStatementListContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSqlStatementList is called when production sqlStatementList is exited.
func (l *trinoListener) ExitSqlStatementList(ctx *grammar.SqlStatementListContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterPrivilege is called when production privilege is entered.
func (l *trinoListener) EnterPrivilege(ctx *grammar.PrivilegeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitPrivilege is called when production privilege is exited.
func (l *trinoListener) ExitPrivilege(ctx *grammar.PrivilegeContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterQueryPeriod is called when production queryPeriod is entered.
func (l *trinoListener) EnterQueryPeriod(ctx *grammar.QueryPeriodContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitQueryPeriod is called when production queryPeriod is exited.
func (l *trinoListener) ExitQueryPeriod(ctx *grammar.QueryPeriodContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterSpecifiedPrincipal is called when production specifiedPrincipal is entered.
func (l *trinoListener) EnterSpecifiedPrincipal(ctx *grammar.SpecifiedPrincipalContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitSpecifiedPrincipal is called when production specifiedPrincipal is exited.
func (l *trinoListener) ExitSpecifiedPrincipal(ctx *grammar.SpecifiedPrincipalContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCurrentUserGrantor is called when production currentUserGrantor is entered.
func (l *trinoListener) EnterCurrentUserGrantor(ctx *grammar.CurrentUserGrantorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCurrentUserGrantor is called when production currentUserGrantor is exited.
func (l *trinoListener) ExitCurrentUserGrantor(ctx *grammar.CurrentUserGrantorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterCurrentRoleGrantor is called when production currentRoleGrantor is entered.
func (l *trinoListener) EnterCurrentRoleGrantor(ctx *grammar.CurrentRoleGrantorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitCurrentRoleGrantor is called when production currentRoleGrantor is exited.
func (l *trinoListener) ExitCurrentRoleGrantor(ctx *grammar.CurrentRoleGrantorContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterUnspecifiedPrincipal is called when production unspecifiedPrincipal is entered.
func (l *trinoListener) EnterUnspecifiedPrincipal(ctx *grammar.UnspecifiedPrincipalContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitUnspecifiedPrincipal is called when production unspecifiedPrincipal is exited.
func (l *trinoListener) ExitUnspecifiedPrincipal(ctx *grammar.UnspecifiedPrincipalContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterUserPrincipal is called when production userPrincipal is entered.
func (l *trinoListener) EnterUserPrincipal(ctx *grammar.UserPrincipalContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitUserPrincipal is called when production userPrincipal is exited.
func (l *trinoListener) ExitUserPrincipal(ctx *grammar.UserPrincipalContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRolePrincipal is called when production rolePrincipal is entered.
func (l *trinoListener) EnterRolePrincipal(ctx *grammar.RolePrincipalContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRolePrincipal is called when production rolePrincipal is exited.
func (l *trinoListener) ExitRolePrincipal(ctx *grammar.RolePrincipalContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterRoles is called when production roles is entered.
func (l *trinoListener) EnterRoles(ctx *grammar.RolesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitRoles is called when production roles is exited.
func (l *trinoListener) ExitRoles(ctx *grammar.RolesContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterIdentifierUser is called when production identifierUser is entered.
func (l *trinoListener) EnterIdentifierUser(ctx *grammar.IdentifierUserContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitIdentifierUser is called when production identifierUser is exited.
func (l *trinoListener) ExitIdentifierUser(ctx *grammar.IdentifierUserContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// EnterStringUser is called when production stringUser is entered.
func (l *trinoListener) EnterStringUser(ctx *grammar.StringUserContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}

// ExitStringUser is called when production stringUser is exited.
func (l *trinoListener) ExitStringUser(ctx *grammar.StringUserContext) {
	log.Println("Unknown context in query: ", l.query, "contextType: ", reflect.TypeOf(ctx))
}
