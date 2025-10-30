// Code generated from TrinoParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar // TrinoParser
import "github.com/antlr4-go/antlr/v4"

// BaseTrinoParserListener is a complete listener for a parse tree produced by TrinoParser.
type BaseTrinoParserListener struct{}

var _ TrinoParserListener = &BaseTrinoParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseTrinoParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseTrinoParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseTrinoParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseTrinoParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseTrinoParserListener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseTrinoParserListener) ExitParse(ctx *ParseContext) {}

// EnterStatements is called when production statements is entered.
func (s *BaseTrinoParserListener) EnterStatements(ctx *StatementsContext) {}

// ExitStatements is called when production statements is exited.
func (s *BaseTrinoParserListener) ExitStatements(ctx *StatementsContext) {}

// EnterSingleStatement is called when production singleStatement is entered.
func (s *BaseTrinoParserListener) EnterSingleStatement(ctx *SingleStatementContext) {}

// ExitSingleStatement is called when production singleStatement is exited.
func (s *BaseTrinoParserListener) ExitSingleStatement(ctx *SingleStatementContext) {}

// EnterStandaloneExpression is called when production standaloneExpression is entered.
func (s *BaseTrinoParserListener) EnterStandaloneExpression(ctx *StandaloneExpressionContext) {}

// ExitStandaloneExpression is called when production standaloneExpression is exited.
func (s *BaseTrinoParserListener) ExitStandaloneExpression(ctx *StandaloneExpressionContext) {}

// EnterStandalonePathSpecification is called when production standalonePathSpecification is entered.
func (s *BaseTrinoParserListener) EnterStandalonePathSpecification(ctx *StandalonePathSpecificationContext) {
}

// ExitStandalonePathSpecification is called when production standalonePathSpecification is exited.
func (s *BaseTrinoParserListener) ExitStandalonePathSpecification(ctx *StandalonePathSpecificationContext) {
}

// EnterStandaloneType is called when production standaloneType is entered.
func (s *BaseTrinoParserListener) EnterStandaloneType(ctx *StandaloneTypeContext) {}

// ExitStandaloneType is called when production standaloneType is exited.
func (s *BaseTrinoParserListener) ExitStandaloneType(ctx *StandaloneTypeContext) {}

// EnterStandaloneRowPattern is called when production standaloneRowPattern is entered.
func (s *BaseTrinoParserListener) EnterStandaloneRowPattern(ctx *StandaloneRowPatternContext) {}

// ExitStandaloneRowPattern is called when production standaloneRowPattern is exited.
func (s *BaseTrinoParserListener) ExitStandaloneRowPattern(ctx *StandaloneRowPatternContext) {}

// EnterStandaloneFunctionSpecification is called when production standaloneFunctionSpecification is entered.
func (s *BaseTrinoParserListener) EnterStandaloneFunctionSpecification(ctx *StandaloneFunctionSpecificationContext) {
}

// ExitStandaloneFunctionSpecification is called when production standaloneFunctionSpecification is exited.
func (s *BaseTrinoParserListener) ExitStandaloneFunctionSpecification(ctx *StandaloneFunctionSpecificationContext) {
}

// EnterStatementDefault is called when production statementDefault is entered.
func (s *BaseTrinoParserListener) EnterStatementDefault(ctx *StatementDefaultContext) {}

// ExitStatementDefault is called when production statementDefault is exited.
func (s *BaseTrinoParserListener) ExitStatementDefault(ctx *StatementDefaultContext) {}

// EnterUse is called when production use is entered.
func (s *BaseTrinoParserListener) EnterUse(ctx *UseContext) {}

// ExitUse is called when production use is exited.
func (s *BaseTrinoParserListener) ExitUse(ctx *UseContext) {}

// EnterCreateCatalog is called when production createCatalog is entered.
func (s *BaseTrinoParserListener) EnterCreateCatalog(ctx *CreateCatalogContext) {}

// ExitCreateCatalog is called when production createCatalog is exited.
func (s *BaseTrinoParserListener) ExitCreateCatalog(ctx *CreateCatalogContext) {}

// EnterDropCatalog is called when production dropCatalog is entered.
func (s *BaseTrinoParserListener) EnterDropCatalog(ctx *DropCatalogContext) {}

// ExitDropCatalog is called when production dropCatalog is exited.
func (s *BaseTrinoParserListener) ExitDropCatalog(ctx *DropCatalogContext) {}

// EnterCreateSchema is called when production createSchema is entered.
func (s *BaseTrinoParserListener) EnterCreateSchema(ctx *CreateSchemaContext) {}

// ExitCreateSchema is called when production createSchema is exited.
func (s *BaseTrinoParserListener) ExitCreateSchema(ctx *CreateSchemaContext) {}

// EnterDropSchema is called when production dropSchema is entered.
func (s *BaseTrinoParserListener) EnterDropSchema(ctx *DropSchemaContext) {}

// ExitDropSchema is called when production dropSchema is exited.
func (s *BaseTrinoParserListener) ExitDropSchema(ctx *DropSchemaContext) {}

// EnterRenameSchema is called when production renameSchema is entered.
func (s *BaseTrinoParserListener) EnterRenameSchema(ctx *RenameSchemaContext) {}

// ExitRenameSchema is called when production renameSchema is exited.
func (s *BaseTrinoParserListener) ExitRenameSchema(ctx *RenameSchemaContext) {}

// EnterSetSchemaAuthorization is called when production setSchemaAuthorization is entered.
func (s *BaseTrinoParserListener) EnterSetSchemaAuthorization(ctx *SetSchemaAuthorizationContext) {}

// ExitSetSchemaAuthorization is called when production setSchemaAuthorization is exited.
func (s *BaseTrinoParserListener) ExitSetSchemaAuthorization(ctx *SetSchemaAuthorizationContext) {}

// EnterCreateTableAsSelect is called when production createTableAsSelect is entered.
func (s *BaseTrinoParserListener) EnterCreateTableAsSelect(ctx *CreateTableAsSelectContext) {}

// ExitCreateTableAsSelect is called when production createTableAsSelect is exited.
func (s *BaseTrinoParserListener) ExitCreateTableAsSelect(ctx *CreateTableAsSelectContext) {}

// EnterCreateTable is called when production createTable is entered.
func (s *BaseTrinoParserListener) EnterCreateTable(ctx *CreateTableContext) {}

// ExitCreateTable is called when production createTable is exited.
func (s *BaseTrinoParserListener) ExitCreateTable(ctx *CreateTableContext) {}

// EnterDropTable is called when production dropTable is entered.
func (s *BaseTrinoParserListener) EnterDropTable(ctx *DropTableContext) {}

// ExitDropTable is called when production dropTable is exited.
func (s *BaseTrinoParserListener) ExitDropTable(ctx *DropTableContext) {}

// EnterInsertInto is called when production insertInto is entered.
func (s *BaseTrinoParserListener) EnterInsertInto(ctx *InsertIntoContext) {}

// ExitInsertInto is called when production insertInto is exited.
func (s *BaseTrinoParserListener) ExitInsertInto(ctx *InsertIntoContext) {}

// EnterDelete is called when production delete is entered.
func (s *BaseTrinoParserListener) EnterDelete(ctx *DeleteContext) {}

// ExitDelete is called when production delete is exited.
func (s *BaseTrinoParserListener) ExitDelete(ctx *DeleteContext) {}

// EnterTruncateTable is called when production truncateTable is entered.
func (s *BaseTrinoParserListener) EnterTruncateTable(ctx *TruncateTableContext) {}

// ExitTruncateTable is called when production truncateTable is exited.
func (s *BaseTrinoParserListener) ExitTruncateTable(ctx *TruncateTableContext) {}

// EnterCommentTable is called when production commentTable is entered.
func (s *BaseTrinoParserListener) EnterCommentTable(ctx *CommentTableContext) {}

// ExitCommentTable is called when production commentTable is exited.
func (s *BaseTrinoParserListener) ExitCommentTable(ctx *CommentTableContext) {}

// EnterCommentView is called when production commentView is entered.
func (s *BaseTrinoParserListener) EnterCommentView(ctx *CommentViewContext) {}

// ExitCommentView is called when production commentView is exited.
func (s *BaseTrinoParserListener) ExitCommentView(ctx *CommentViewContext) {}

// EnterCommentColumn is called when production commentColumn is entered.
func (s *BaseTrinoParserListener) EnterCommentColumn(ctx *CommentColumnContext) {}

// ExitCommentColumn is called when production commentColumn is exited.
func (s *BaseTrinoParserListener) ExitCommentColumn(ctx *CommentColumnContext) {}

// EnterRenameTable is called when production renameTable is entered.
func (s *BaseTrinoParserListener) EnterRenameTable(ctx *RenameTableContext) {}

// ExitRenameTable is called when production renameTable is exited.
func (s *BaseTrinoParserListener) ExitRenameTable(ctx *RenameTableContext) {}

// EnterAddColumn is called when production addColumn is entered.
func (s *BaseTrinoParserListener) EnterAddColumn(ctx *AddColumnContext) {}

// ExitAddColumn is called when production addColumn is exited.
func (s *BaseTrinoParserListener) ExitAddColumn(ctx *AddColumnContext) {}

// EnterRenameColumn is called when production renameColumn is entered.
func (s *BaseTrinoParserListener) EnterRenameColumn(ctx *RenameColumnContext) {}

// ExitRenameColumn is called when production renameColumn is exited.
func (s *BaseTrinoParserListener) ExitRenameColumn(ctx *RenameColumnContext) {}

// EnterDropColumn is called when production dropColumn is entered.
func (s *BaseTrinoParserListener) EnterDropColumn(ctx *DropColumnContext) {}

// ExitDropColumn is called when production dropColumn is exited.
func (s *BaseTrinoParserListener) ExitDropColumn(ctx *DropColumnContext) {}

// EnterSetColumnType is called when production setColumnType is entered.
func (s *BaseTrinoParserListener) EnterSetColumnType(ctx *SetColumnTypeContext) {}

// ExitSetColumnType is called when production setColumnType is exited.
func (s *BaseTrinoParserListener) ExitSetColumnType(ctx *SetColumnTypeContext) {}

// EnterSetTableAuthorization is called when production setTableAuthorization is entered.
func (s *BaseTrinoParserListener) EnterSetTableAuthorization(ctx *SetTableAuthorizationContext) {}

// ExitSetTableAuthorization is called when production setTableAuthorization is exited.
func (s *BaseTrinoParserListener) ExitSetTableAuthorization(ctx *SetTableAuthorizationContext) {}

// EnterSetTableProperties is called when production setTableProperties is entered.
func (s *BaseTrinoParserListener) EnterSetTableProperties(ctx *SetTablePropertiesContext) {}

// ExitSetTableProperties is called when production setTableProperties is exited.
func (s *BaseTrinoParserListener) ExitSetTableProperties(ctx *SetTablePropertiesContext) {}

// EnterTableExecute is called when production tableExecute is entered.
func (s *BaseTrinoParserListener) EnterTableExecute(ctx *TableExecuteContext) {}

// ExitTableExecute is called when production tableExecute is exited.
func (s *BaseTrinoParserListener) ExitTableExecute(ctx *TableExecuteContext) {}

// EnterAnalyze is called when production analyze is entered.
func (s *BaseTrinoParserListener) EnterAnalyze(ctx *AnalyzeContext) {}

// ExitAnalyze is called when production analyze is exited.
func (s *BaseTrinoParserListener) ExitAnalyze(ctx *AnalyzeContext) {}

// EnterCreateMaterializedView is called when production createMaterializedView is entered.
func (s *BaseTrinoParserListener) EnterCreateMaterializedView(ctx *CreateMaterializedViewContext) {}

// ExitCreateMaterializedView is called when production createMaterializedView is exited.
func (s *BaseTrinoParserListener) ExitCreateMaterializedView(ctx *CreateMaterializedViewContext) {}

// EnterCreateView is called when production createView is entered.
func (s *BaseTrinoParserListener) EnterCreateView(ctx *CreateViewContext) {}

// ExitCreateView is called when production createView is exited.
func (s *BaseTrinoParserListener) ExitCreateView(ctx *CreateViewContext) {}

// EnterRefreshMaterializedView is called when production refreshMaterializedView is entered.
func (s *BaseTrinoParserListener) EnterRefreshMaterializedView(ctx *RefreshMaterializedViewContext) {}

// ExitRefreshMaterializedView is called when production refreshMaterializedView is exited.
func (s *BaseTrinoParserListener) ExitRefreshMaterializedView(ctx *RefreshMaterializedViewContext) {}

// EnterDropMaterializedView is called when production dropMaterializedView is entered.
func (s *BaseTrinoParserListener) EnterDropMaterializedView(ctx *DropMaterializedViewContext) {}

// ExitDropMaterializedView is called when production dropMaterializedView is exited.
func (s *BaseTrinoParserListener) ExitDropMaterializedView(ctx *DropMaterializedViewContext) {}

// EnterRenameMaterializedView is called when production renameMaterializedView is entered.
func (s *BaseTrinoParserListener) EnterRenameMaterializedView(ctx *RenameMaterializedViewContext) {}

// ExitRenameMaterializedView is called when production renameMaterializedView is exited.
func (s *BaseTrinoParserListener) ExitRenameMaterializedView(ctx *RenameMaterializedViewContext) {}

// EnterSetMaterializedViewProperties is called when production setMaterializedViewProperties is entered.
func (s *BaseTrinoParserListener) EnterSetMaterializedViewProperties(ctx *SetMaterializedViewPropertiesContext) {
}

// ExitSetMaterializedViewProperties is called when production setMaterializedViewProperties is exited.
func (s *BaseTrinoParserListener) ExitSetMaterializedViewProperties(ctx *SetMaterializedViewPropertiesContext) {
}

// EnterDropView is called when production dropView is entered.
func (s *BaseTrinoParserListener) EnterDropView(ctx *DropViewContext) {}

// ExitDropView is called when production dropView is exited.
func (s *BaseTrinoParserListener) ExitDropView(ctx *DropViewContext) {}

// EnterRenameView is called when production renameView is entered.
func (s *BaseTrinoParserListener) EnterRenameView(ctx *RenameViewContext) {}

// ExitRenameView is called when production renameView is exited.
func (s *BaseTrinoParserListener) ExitRenameView(ctx *RenameViewContext) {}

// EnterSetViewAuthorization is called when production setViewAuthorization is entered.
func (s *BaseTrinoParserListener) EnterSetViewAuthorization(ctx *SetViewAuthorizationContext) {}

// ExitSetViewAuthorization is called when production setViewAuthorization is exited.
func (s *BaseTrinoParserListener) ExitSetViewAuthorization(ctx *SetViewAuthorizationContext) {}

// EnterCall is called when production call is entered.
func (s *BaseTrinoParserListener) EnterCall(ctx *CallContext) {}

// ExitCall is called when production call is exited.
func (s *BaseTrinoParserListener) ExitCall(ctx *CallContext) {}

// EnterCreateFunction is called when production createFunction is entered.
func (s *BaseTrinoParserListener) EnterCreateFunction(ctx *CreateFunctionContext) {}

// ExitCreateFunction is called when production createFunction is exited.
func (s *BaseTrinoParserListener) ExitCreateFunction(ctx *CreateFunctionContext) {}

// EnterDropFunction is called when production dropFunction is entered.
func (s *BaseTrinoParserListener) EnterDropFunction(ctx *DropFunctionContext) {}

// ExitDropFunction is called when production dropFunction is exited.
func (s *BaseTrinoParserListener) ExitDropFunction(ctx *DropFunctionContext) {}

// EnterCreateRole is called when production createRole is entered.
func (s *BaseTrinoParserListener) EnterCreateRole(ctx *CreateRoleContext) {}

// ExitCreateRole is called when production createRole is exited.
func (s *BaseTrinoParserListener) ExitCreateRole(ctx *CreateRoleContext) {}

// EnterDropRole is called when production dropRole is entered.
func (s *BaseTrinoParserListener) EnterDropRole(ctx *DropRoleContext) {}

// ExitDropRole is called when production dropRole is exited.
func (s *BaseTrinoParserListener) ExitDropRole(ctx *DropRoleContext) {}

// EnterGrantRoles is called when production grantRoles is entered.
func (s *BaseTrinoParserListener) EnterGrantRoles(ctx *GrantRolesContext) {}

// ExitGrantRoles is called when production grantRoles is exited.
func (s *BaseTrinoParserListener) ExitGrantRoles(ctx *GrantRolesContext) {}

// EnterRevokeRoles is called when production revokeRoles is entered.
func (s *BaseTrinoParserListener) EnterRevokeRoles(ctx *RevokeRolesContext) {}

// ExitRevokeRoles is called when production revokeRoles is exited.
func (s *BaseTrinoParserListener) ExitRevokeRoles(ctx *RevokeRolesContext) {}

// EnterSetRole is called when production setRole is entered.
func (s *BaseTrinoParserListener) EnterSetRole(ctx *SetRoleContext) {}

// ExitSetRole is called when production setRole is exited.
func (s *BaseTrinoParserListener) ExitSetRole(ctx *SetRoleContext) {}

// EnterGrant is called when production grant is entered.
func (s *BaseTrinoParserListener) EnterGrant(ctx *GrantContext) {}

// ExitGrant is called when production grant is exited.
func (s *BaseTrinoParserListener) ExitGrant(ctx *GrantContext) {}

// EnterDeny is called when production deny is entered.
func (s *BaseTrinoParserListener) EnterDeny(ctx *DenyContext) {}

// ExitDeny is called when production deny is exited.
func (s *BaseTrinoParserListener) ExitDeny(ctx *DenyContext) {}

// EnterRevoke is called when production revoke is entered.
func (s *BaseTrinoParserListener) EnterRevoke(ctx *RevokeContext) {}

// ExitRevoke is called when production revoke is exited.
func (s *BaseTrinoParserListener) ExitRevoke(ctx *RevokeContext) {}

// EnterShowGrants is called when production showGrants is entered.
func (s *BaseTrinoParserListener) EnterShowGrants(ctx *ShowGrantsContext) {}

// ExitShowGrants is called when production showGrants is exited.
func (s *BaseTrinoParserListener) ExitShowGrants(ctx *ShowGrantsContext) {}

// EnterExplain is called when production explain is entered.
func (s *BaseTrinoParserListener) EnterExplain(ctx *ExplainContext) {}

// ExitExplain is called when production explain is exited.
func (s *BaseTrinoParserListener) ExitExplain(ctx *ExplainContext) {}

// EnterExplainAnalyze is called when production explainAnalyze is entered.
func (s *BaseTrinoParserListener) EnterExplainAnalyze(ctx *ExplainAnalyzeContext) {}

// ExitExplainAnalyze is called when production explainAnalyze is exited.
func (s *BaseTrinoParserListener) ExitExplainAnalyze(ctx *ExplainAnalyzeContext) {}

// EnterShowCreateTable is called when production showCreateTable is entered.
func (s *BaseTrinoParserListener) EnterShowCreateTable(ctx *ShowCreateTableContext) {}

// ExitShowCreateTable is called when production showCreateTable is exited.
func (s *BaseTrinoParserListener) ExitShowCreateTable(ctx *ShowCreateTableContext) {}

// EnterShowCreateSchema is called when production showCreateSchema is entered.
func (s *BaseTrinoParserListener) EnterShowCreateSchema(ctx *ShowCreateSchemaContext) {}

// ExitShowCreateSchema is called when production showCreateSchema is exited.
func (s *BaseTrinoParserListener) ExitShowCreateSchema(ctx *ShowCreateSchemaContext) {}

// EnterShowCreateView is called when production showCreateView is entered.
func (s *BaseTrinoParserListener) EnterShowCreateView(ctx *ShowCreateViewContext) {}

// ExitShowCreateView is called when production showCreateView is exited.
func (s *BaseTrinoParserListener) ExitShowCreateView(ctx *ShowCreateViewContext) {}

// EnterShowCreateMaterializedView is called when production showCreateMaterializedView is entered.
func (s *BaseTrinoParserListener) EnterShowCreateMaterializedView(ctx *ShowCreateMaterializedViewContext) {
}

// ExitShowCreateMaterializedView is called when production showCreateMaterializedView is exited.
func (s *BaseTrinoParserListener) ExitShowCreateMaterializedView(ctx *ShowCreateMaterializedViewContext) {
}

// EnterShowTables is called when production showTables is entered.
func (s *BaseTrinoParserListener) EnterShowTables(ctx *ShowTablesContext) {}

// ExitShowTables is called when production showTables is exited.
func (s *BaseTrinoParserListener) ExitShowTables(ctx *ShowTablesContext) {}

// EnterShowSchemas is called when production showSchemas is entered.
func (s *BaseTrinoParserListener) EnterShowSchemas(ctx *ShowSchemasContext) {}

// ExitShowSchemas is called when production showSchemas is exited.
func (s *BaseTrinoParserListener) ExitShowSchemas(ctx *ShowSchemasContext) {}

// EnterShowCatalogs is called when production showCatalogs is entered.
func (s *BaseTrinoParserListener) EnterShowCatalogs(ctx *ShowCatalogsContext) {}

// ExitShowCatalogs is called when production showCatalogs is exited.
func (s *BaseTrinoParserListener) ExitShowCatalogs(ctx *ShowCatalogsContext) {}

// EnterShowColumns is called when production showColumns is entered.
func (s *BaseTrinoParserListener) EnterShowColumns(ctx *ShowColumnsContext) {}

// ExitShowColumns is called when production showColumns is exited.
func (s *BaseTrinoParserListener) ExitShowColumns(ctx *ShowColumnsContext) {}

// EnterShowStats is called when production showStats is entered.
func (s *BaseTrinoParserListener) EnterShowStats(ctx *ShowStatsContext) {}

// ExitShowStats is called when production showStats is exited.
func (s *BaseTrinoParserListener) ExitShowStats(ctx *ShowStatsContext) {}

// EnterShowStatsForQuery is called when production showStatsForQuery is entered.
func (s *BaseTrinoParserListener) EnterShowStatsForQuery(ctx *ShowStatsForQueryContext) {}

// ExitShowStatsForQuery is called when production showStatsForQuery is exited.
func (s *BaseTrinoParserListener) ExitShowStatsForQuery(ctx *ShowStatsForQueryContext) {}

// EnterShowRoles is called when production showRoles is entered.
func (s *BaseTrinoParserListener) EnterShowRoles(ctx *ShowRolesContext) {}

// ExitShowRoles is called when production showRoles is exited.
func (s *BaseTrinoParserListener) ExitShowRoles(ctx *ShowRolesContext) {}

// EnterShowRoleGrants is called when production showRoleGrants is entered.
func (s *BaseTrinoParserListener) EnterShowRoleGrants(ctx *ShowRoleGrantsContext) {}

// ExitShowRoleGrants is called when production showRoleGrants is exited.
func (s *BaseTrinoParserListener) ExitShowRoleGrants(ctx *ShowRoleGrantsContext) {}

// EnterShowFunctions is called when production showFunctions is entered.
func (s *BaseTrinoParserListener) EnterShowFunctions(ctx *ShowFunctionsContext) {}

// ExitShowFunctions is called when production showFunctions is exited.
func (s *BaseTrinoParserListener) ExitShowFunctions(ctx *ShowFunctionsContext) {}

// EnterShowSession is called when production showSession is entered.
func (s *BaseTrinoParserListener) EnterShowSession(ctx *ShowSessionContext) {}

// ExitShowSession is called when production showSession is exited.
func (s *BaseTrinoParserListener) ExitShowSession(ctx *ShowSessionContext) {}

// EnterSetSessionAuthorization is called when production setSessionAuthorization is entered.
func (s *BaseTrinoParserListener) EnterSetSessionAuthorization(ctx *SetSessionAuthorizationContext) {}

// ExitSetSessionAuthorization is called when production setSessionAuthorization is exited.
func (s *BaseTrinoParserListener) ExitSetSessionAuthorization(ctx *SetSessionAuthorizationContext) {}

// EnterResetSessionAuthorization is called when production resetSessionAuthorization is entered.
func (s *BaseTrinoParserListener) EnterResetSessionAuthorization(ctx *ResetSessionAuthorizationContext) {
}

// ExitResetSessionAuthorization is called when production resetSessionAuthorization is exited.
func (s *BaseTrinoParserListener) ExitResetSessionAuthorization(ctx *ResetSessionAuthorizationContext) {
}

// EnterSetSession is called when production setSession is entered.
func (s *BaseTrinoParserListener) EnterSetSession(ctx *SetSessionContext) {}

// ExitSetSession is called when production setSession is exited.
func (s *BaseTrinoParserListener) ExitSetSession(ctx *SetSessionContext) {}

// EnterResetSession is called when production resetSession is entered.
func (s *BaseTrinoParserListener) EnterResetSession(ctx *ResetSessionContext) {}

// ExitResetSession is called when production resetSession is exited.
func (s *BaseTrinoParserListener) ExitResetSession(ctx *ResetSessionContext) {}

// EnterStartTransaction is called when production startTransaction is entered.
func (s *BaseTrinoParserListener) EnterStartTransaction(ctx *StartTransactionContext) {}

// ExitStartTransaction is called when production startTransaction is exited.
func (s *BaseTrinoParserListener) ExitStartTransaction(ctx *StartTransactionContext) {}

// EnterCommit is called when production commit is entered.
func (s *BaseTrinoParserListener) EnterCommit(ctx *CommitContext) {}

// ExitCommit is called when production commit is exited.
func (s *BaseTrinoParserListener) ExitCommit(ctx *CommitContext) {}

// EnterRollback is called when production rollback is entered.
func (s *BaseTrinoParserListener) EnterRollback(ctx *RollbackContext) {}

// ExitRollback is called when production rollback is exited.
func (s *BaseTrinoParserListener) ExitRollback(ctx *RollbackContext) {}

// EnterPrepare is called when production prepare is entered.
func (s *BaseTrinoParserListener) EnterPrepare(ctx *PrepareContext) {}

// ExitPrepare is called when production prepare is exited.
func (s *BaseTrinoParserListener) ExitPrepare(ctx *PrepareContext) {}

// EnterDeallocate is called when production deallocate is entered.
func (s *BaseTrinoParserListener) EnterDeallocate(ctx *DeallocateContext) {}

// ExitDeallocate is called when production deallocate is exited.
func (s *BaseTrinoParserListener) ExitDeallocate(ctx *DeallocateContext) {}

// EnterExecute is called when production execute is entered.
func (s *BaseTrinoParserListener) EnterExecute(ctx *ExecuteContext) {}

// ExitExecute is called when production execute is exited.
func (s *BaseTrinoParserListener) ExitExecute(ctx *ExecuteContext) {}

// EnterExecuteImmediate is called when production executeImmediate is entered.
func (s *BaseTrinoParserListener) EnterExecuteImmediate(ctx *ExecuteImmediateContext) {}

// ExitExecuteImmediate is called when production executeImmediate is exited.
func (s *BaseTrinoParserListener) ExitExecuteImmediate(ctx *ExecuteImmediateContext) {}

// EnterDescribeInput is called when production describeInput is entered.
func (s *BaseTrinoParserListener) EnterDescribeInput(ctx *DescribeInputContext) {}

// ExitDescribeInput is called when production describeInput is exited.
func (s *BaseTrinoParserListener) ExitDescribeInput(ctx *DescribeInputContext) {}

// EnterDescribeOutput is called when production describeOutput is entered.
func (s *BaseTrinoParserListener) EnterDescribeOutput(ctx *DescribeOutputContext) {}

// ExitDescribeOutput is called when production describeOutput is exited.
func (s *BaseTrinoParserListener) ExitDescribeOutput(ctx *DescribeOutputContext) {}

// EnterSetPath is called when production setPath is entered.
func (s *BaseTrinoParserListener) EnterSetPath(ctx *SetPathContext) {}

// ExitSetPath is called when production setPath is exited.
func (s *BaseTrinoParserListener) ExitSetPath(ctx *SetPathContext) {}

// EnterSetTimeZone is called when production setTimeZone is entered.
func (s *BaseTrinoParserListener) EnterSetTimeZone(ctx *SetTimeZoneContext) {}

// ExitSetTimeZone is called when production setTimeZone is exited.
func (s *BaseTrinoParserListener) ExitSetTimeZone(ctx *SetTimeZoneContext) {}

// EnterUpdate is called when production update is entered.
func (s *BaseTrinoParserListener) EnterUpdate(ctx *UpdateContext) {}

// ExitUpdate is called when production update is exited.
func (s *BaseTrinoParserListener) ExitUpdate(ctx *UpdateContext) {}

// EnterMerge is called when production merge is entered.
func (s *BaseTrinoParserListener) EnterMerge(ctx *MergeContext) {}

// ExitMerge is called when production merge is exited.
func (s *BaseTrinoParserListener) ExitMerge(ctx *MergeContext) {}

// EnterRootQuery is called when production rootQuery is entered.
func (s *BaseTrinoParserListener) EnterRootQuery(ctx *RootQueryContext) {}

// ExitRootQuery is called when production rootQuery is exited.
func (s *BaseTrinoParserListener) ExitRootQuery(ctx *RootQueryContext) {}

// EnterWithFunction is called when production withFunction is entered.
func (s *BaseTrinoParserListener) EnterWithFunction(ctx *WithFunctionContext) {}

// ExitWithFunction is called when production withFunction is exited.
func (s *BaseTrinoParserListener) ExitWithFunction(ctx *WithFunctionContext) {}

// EnterQuery is called when production query is entered.
func (s *BaseTrinoParserListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BaseTrinoParserListener) ExitQuery(ctx *QueryContext) {}

// EnterWith is called when production with is entered.
func (s *BaseTrinoParserListener) EnterWith(ctx *WithContext) {}

// ExitWith is called when production with is exited.
func (s *BaseTrinoParserListener) ExitWith(ctx *WithContext) {}

// EnterTableElement is called when production tableElement is entered.
func (s *BaseTrinoParserListener) EnterTableElement(ctx *TableElementContext) {}

// ExitTableElement is called when production tableElement is exited.
func (s *BaseTrinoParserListener) ExitTableElement(ctx *TableElementContext) {}

// EnterColumnDefinition is called when production columnDefinition is entered.
func (s *BaseTrinoParserListener) EnterColumnDefinition(ctx *ColumnDefinitionContext) {}

// ExitColumnDefinition is called when production columnDefinition is exited.
func (s *BaseTrinoParserListener) ExitColumnDefinition(ctx *ColumnDefinitionContext) {}

// EnterLikeClause is called when production likeClause is entered.
func (s *BaseTrinoParserListener) EnterLikeClause(ctx *LikeClauseContext) {}

// ExitLikeClause is called when production likeClause is exited.
func (s *BaseTrinoParserListener) ExitLikeClause(ctx *LikeClauseContext) {}

// EnterProperties is called when production properties is entered.
func (s *BaseTrinoParserListener) EnterProperties(ctx *PropertiesContext) {}

// ExitProperties is called when production properties is exited.
func (s *BaseTrinoParserListener) ExitProperties(ctx *PropertiesContext) {}

// EnterPropertyAssignments is called when production propertyAssignments is entered.
func (s *BaseTrinoParserListener) EnterPropertyAssignments(ctx *PropertyAssignmentsContext) {}

// ExitPropertyAssignments is called when production propertyAssignments is exited.
func (s *BaseTrinoParserListener) ExitPropertyAssignments(ctx *PropertyAssignmentsContext) {}

// EnterProperty is called when production property is entered.
func (s *BaseTrinoParserListener) EnterProperty(ctx *PropertyContext) {}

// ExitProperty is called when production property is exited.
func (s *BaseTrinoParserListener) ExitProperty(ctx *PropertyContext) {}

// EnterDefaultPropertyValue is called when production defaultPropertyValue is entered.
func (s *BaseTrinoParserListener) EnterDefaultPropertyValue(ctx *DefaultPropertyValueContext) {}

// ExitDefaultPropertyValue is called when production defaultPropertyValue is exited.
func (s *BaseTrinoParserListener) ExitDefaultPropertyValue(ctx *DefaultPropertyValueContext) {}

// EnterNonDefaultPropertyValue is called when production nonDefaultPropertyValue is entered.
func (s *BaseTrinoParserListener) EnterNonDefaultPropertyValue(ctx *NonDefaultPropertyValueContext) {}

// ExitNonDefaultPropertyValue is called when production nonDefaultPropertyValue is exited.
func (s *BaseTrinoParserListener) ExitNonDefaultPropertyValue(ctx *NonDefaultPropertyValueContext) {}

// EnterQueryNoWith is called when production queryNoWith is entered.
func (s *BaseTrinoParserListener) EnterQueryNoWith(ctx *QueryNoWithContext) {}

// ExitQueryNoWith is called when production queryNoWith is exited.
func (s *BaseTrinoParserListener) ExitQueryNoWith(ctx *QueryNoWithContext) {}

// EnterLimitRowCount is called when production limitRowCount is entered.
func (s *BaseTrinoParserListener) EnterLimitRowCount(ctx *LimitRowCountContext) {}

// ExitLimitRowCount is called when production limitRowCount is exited.
func (s *BaseTrinoParserListener) ExitLimitRowCount(ctx *LimitRowCountContext) {}

// EnterRowCount is called when production rowCount is entered.
func (s *BaseTrinoParserListener) EnterRowCount(ctx *RowCountContext) {}

// ExitRowCount is called when production rowCount is exited.
func (s *BaseTrinoParserListener) ExitRowCount(ctx *RowCountContext) {}

// EnterQueryTermDefault is called when production queryTermDefault is entered.
func (s *BaseTrinoParserListener) EnterQueryTermDefault(ctx *QueryTermDefaultContext) {}

// ExitQueryTermDefault is called when production queryTermDefault is exited.
func (s *BaseTrinoParserListener) ExitQueryTermDefault(ctx *QueryTermDefaultContext) {}

// EnterSetOperation is called when production setOperation is entered.
func (s *BaseTrinoParserListener) EnterSetOperation(ctx *SetOperationContext) {}

// ExitSetOperation is called when production setOperation is exited.
func (s *BaseTrinoParserListener) ExitSetOperation(ctx *SetOperationContext) {}

// EnterQueryPrimaryDefault is called when production queryPrimaryDefault is entered.
func (s *BaseTrinoParserListener) EnterQueryPrimaryDefault(ctx *QueryPrimaryDefaultContext) {}

// ExitQueryPrimaryDefault is called when production queryPrimaryDefault is exited.
func (s *BaseTrinoParserListener) ExitQueryPrimaryDefault(ctx *QueryPrimaryDefaultContext) {}

// EnterTable is called when production table is entered.
func (s *BaseTrinoParserListener) EnterTable(ctx *TableContext) {}

// ExitTable is called when production table is exited.
func (s *BaseTrinoParserListener) ExitTable(ctx *TableContext) {}

// EnterInlineTable is called when production inlineTable is entered.
func (s *BaseTrinoParserListener) EnterInlineTable(ctx *InlineTableContext) {}

// ExitInlineTable is called when production inlineTable is exited.
func (s *BaseTrinoParserListener) ExitInlineTable(ctx *InlineTableContext) {}

// EnterSubquery is called when production subquery is entered.
func (s *BaseTrinoParserListener) EnterSubquery(ctx *SubqueryContext) {}

// ExitSubquery is called when production subquery is exited.
func (s *BaseTrinoParserListener) ExitSubquery(ctx *SubqueryContext) {}

// EnterSortItem is called when production sortItem is entered.
func (s *BaseTrinoParserListener) EnterSortItem(ctx *SortItemContext) {}

// ExitSortItem is called when production sortItem is exited.
func (s *BaseTrinoParserListener) ExitSortItem(ctx *SortItemContext) {}

// EnterQuerySpecification is called when production querySpecification is entered.
func (s *BaseTrinoParserListener) EnterQuerySpecification(ctx *QuerySpecificationContext) {}

// ExitQuerySpecification is called when production querySpecification is exited.
func (s *BaseTrinoParserListener) ExitQuerySpecification(ctx *QuerySpecificationContext) {}

// EnterGroupBy is called when production groupBy is entered.
func (s *BaseTrinoParserListener) EnterGroupBy(ctx *GroupByContext) {}

// ExitGroupBy is called when production groupBy is exited.
func (s *BaseTrinoParserListener) ExitGroupBy(ctx *GroupByContext) {}

// EnterSingleGroupingSet is called when production singleGroupingSet is entered.
func (s *BaseTrinoParserListener) EnterSingleGroupingSet(ctx *SingleGroupingSetContext) {}

// ExitSingleGroupingSet is called when production singleGroupingSet is exited.
func (s *BaseTrinoParserListener) ExitSingleGroupingSet(ctx *SingleGroupingSetContext) {}

// EnterRollup is called when production rollup is entered.
func (s *BaseTrinoParserListener) EnterRollup(ctx *RollupContext) {}

// ExitRollup is called when production rollup is exited.
func (s *BaseTrinoParserListener) ExitRollup(ctx *RollupContext) {}

// EnterCube is called when production cube is entered.
func (s *BaseTrinoParserListener) EnterCube(ctx *CubeContext) {}

// ExitCube is called when production cube is exited.
func (s *BaseTrinoParserListener) ExitCube(ctx *CubeContext) {}

// EnterMultipleGroupingSets is called when production multipleGroupingSets is entered.
func (s *BaseTrinoParserListener) EnterMultipleGroupingSets(ctx *MultipleGroupingSetsContext) {}

// ExitMultipleGroupingSets is called when production multipleGroupingSets is exited.
func (s *BaseTrinoParserListener) ExitMultipleGroupingSets(ctx *MultipleGroupingSetsContext) {}

// EnterGroupingSet is called when production groupingSet is entered.
func (s *BaseTrinoParserListener) EnterGroupingSet(ctx *GroupingSetContext) {}

// ExitGroupingSet is called when production groupingSet is exited.
func (s *BaseTrinoParserListener) ExitGroupingSet(ctx *GroupingSetContext) {}

// EnterWindowDefinition is called when production windowDefinition is entered.
func (s *BaseTrinoParserListener) EnterWindowDefinition(ctx *WindowDefinitionContext) {}

// ExitWindowDefinition is called when production windowDefinition is exited.
func (s *BaseTrinoParserListener) ExitWindowDefinition(ctx *WindowDefinitionContext) {}

// EnterWindowSpecification is called when production windowSpecification is entered.
func (s *BaseTrinoParserListener) EnterWindowSpecification(ctx *WindowSpecificationContext) {}

// ExitWindowSpecification is called when production windowSpecification is exited.
func (s *BaseTrinoParserListener) ExitWindowSpecification(ctx *WindowSpecificationContext) {}

// EnterNamedQuery is called when production namedQuery is entered.
func (s *BaseTrinoParserListener) EnterNamedQuery(ctx *NamedQueryContext) {}

// ExitNamedQuery is called when production namedQuery is exited.
func (s *BaseTrinoParserListener) ExitNamedQuery(ctx *NamedQueryContext) {}

// EnterSetQuantifier is called when production setQuantifier is entered.
func (s *BaseTrinoParserListener) EnterSetQuantifier(ctx *SetQuantifierContext) {}

// ExitSetQuantifier is called when production setQuantifier is exited.
func (s *BaseTrinoParserListener) ExitSetQuantifier(ctx *SetQuantifierContext) {}

// EnterSelectSingle is called when production selectSingle is entered.
func (s *BaseTrinoParserListener) EnterSelectSingle(ctx *SelectSingleContext) {}

// ExitSelectSingle is called when production selectSingle is exited.
func (s *BaseTrinoParserListener) ExitSelectSingle(ctx *SelectSingleContext) {}

// EnterSelectAll is called when production selectAll is entered.
func (s *BaseTrinoParserListener) EnterSelectAll(ctx *SelectAllContext) {}

// ExitSelectAll is called when production selectAll is exited.
func (s *BaseTrinoParserListener) ExitSelectAll(ctx *SelectAllContext) {}

// EnterRelationDefault is called when production relationDefault is entered.
func (s *BaseTrinoParserListener) EnterRelationDefault(ctx *RelationDefaultContext) {}

// ExitRelationDefault is called when production relationDefault is exited.
func (s *BaseTrinoParserListener) ExitRelationDefault(ctx *RelationDefaultContext) {}

// EnterJoinRelation is called when production joinRelation is entered.
func (s *BaseTrinoParserListener) EnterJoinRelation(ctx *JoinRelationContext) {}

// ExitJoinRelation is called when production joinRelation is exited.
func (s *BaseTrinoParserListener) ExitJoinRelation(ctx *JoinRelationContext) {}

// EnterJoinType is called when production joinType is entered.
func (s *BaseTrinoParserListener) EnterJoinType(ctx *JoinTypeContext) {}

// ExitJoinType is called when production joinType is exited.
func (s *BaseTrinoParserListener) ExitJoinType(ctx *JoinTypeContext) {}

// EnterJoinCriteria is called when production joinCriteria is entered.
func (s *BaseTrinoParserListener) EnterJoinCriteria(ctx *JoinCriteriaContext) {}

// ExitJoinCriteria is called when production joinCriteria is exited.
func (s *BaseTrinoParserListener) ExitJoinCriteria(ctx *JoinCriteriaContext) {}

// EnterSampledRelation is called when production sampledRelation is entered.
func (s *BaseTrinoParserListener) EnterSampledRelation(ctx *SampledRelationContext) {}

// ExitSampledRelation is called when production sampledRelation is exited.
func (s *BaseTrinoParserListener) ExitSampledRelation(ctx *SampledRelationContext) {}

// EnterSampleType is called when production sampleType is entered.
func (s *BaseTrinoParserListener) EnterSampleType(ctx *SampleTypeContext) {}

// ExitSampleType is called when production sampleType is exited.
func (s *BaseTrinoParserListener) ExitSampleType(ctx *SampleTypeContext) {}

// EnterTrimsSpecification is called when production trimsSpecification is entered.
func (s *BaseTrinoParserListener) EnterTrimsSpecification(ctx *TrimsSpecificationContext) {}

// ExitTrimsSpecification is called when production trimsSpecification is exited.
func (s *BaseTrinoParserListener) ExitTrimsSpecification(ctx *TrimsSpecificationContext) {}

// EnterListAggOverflowBehavior is called when production listAggOverflowBehavior is entered.
func (s *BaseTrinoParserListener) EnterListAggOverflowBehavior(ctx *ListAggOverflowBehaviorContext) {}

// ExitListAggOverflowBehavior is called when production listAggOverflowBehavior is exited.
func (s *BaseTrinoParserListener) ExitListAggOverflowBehavior(ctx *ListAggOverflowBehaviorContext) {}

// EnterListaggCountIndication is called when production listaggCountIndication is entered.
func (s *BaseTrinoParserListener) EnterListaggCountIndication(ctx *ListaggCountIndicationContext) {}

// ExitListaggCountIndication is called when production listaggCountIndication is exited.
func (s *BaseTrinoParserListener) ExitListaggCountIndication(ctx *ListaggCountIndicationContext) {}

// EnterPatternRecognition is called when production patternRecognition is entered.
func (s *BaseTrinoParserListener) EnterPatternRecognition(ctx *PatternRecognitionContext) {}

// ExitPatternRecognition is called when production patternRecognition is exited.
func (s *BaseTrinoParserListener) ExitPatternRecognition(ctx *PatternRecognitionContext) {}

// EnterMeasureDefinition is called when production measureDefinition is entered.
func (s *BaseTrinoParserListener) EnterMeasureDefinition(ctx *MeasureDefinitionContext) {}

// ExitMeasureDefinition is called when production measureDefinition is exited.
func (s *BaseTrinoParserListener) ExitMeasureDefinition(ctx *MeasureDefinitionContext) {}

// EnterRowsPerMatch is called when production rowsPerMatch is entered.
func (s *BaseTrinoParserListener) EnterRowsPerMatch(ctx *RowsPerMatchContext) {}

// ExitRowsPerMatch is called when production rowsPerMatch is exited.
func (s *BaseTrinoParserListener) ExitRowsPerMatch(ctx *RowsPerMatchContext) {}

// EnterEmptyMatchHandling is called when production emptyMatchHandling is entered.
func (s *BaseTrinoParserListener) EnterEmptyMatchHandling(ctx *EmptyMatchHandlingContext) {}

// ExitEmptyMatchHandling is called when production emptyMatchHandling is exited.
func (s *BaseTrinoParserListener) ExitEmptyMatchHandling(ctx *EmptyMatchHandlingContext) {}

// EnterSkipTo is called when production skipTo is entered.
func (s *BaseTrinoParserListener) EnterSkipTo(ctx *SkipToContext) {}

// ExitSkipTo is called when production skipTo is exited.
func (s *BaseTrinoParserListener) ExitSkipTo(ctx *SkipToContext) {}

// EnterSubsetDefinition is called when production subsetDefinition is entered.
func (s *BaseTrinoParserListener) EnterSubsetDefinition(ctx *SubsetDefinitionContext) {}

// ExitSubsetDefinition is called when production subsetDefinition is exited.
func (s *BaseTrinoParserListener) ExitSubsetDefinition(ctx *SubsetDefinitionContext) {}

// EnterVariableDefinition is called when production variableDefinition is entered.
func (s *BaseTrinoParserListener) EnterVariableDefinition(ctx *VariableDefinitionContext) {}

// ExitVariableDefinition is called when production variableDefinition is exited.
func (s *BaseTrinoParserListener) ExitVariableDefinition(ctx *VariableDefinitionContext) {}

// EnterAliasedRelation is called when production aliasedRelation is entered.
func (s *BaseTrinoParserListener) EnterAliasedRelation(ctx *AliasedRelationContext) {}

// ExitAliasedRelation is called when production aliasedRelation is exited.
func (s *BaseTrinoParserListener) ExitAliasedRelation(ctx *AliasedRelationContext) {}

// EnterColumnAliases is called when production columnAliases is entered.
func (s *BaseTrinoParserListener) EnterColumnAliases(ctx *ColumnAliasesContext) {}

// ExitColumnAliases is called when production columnAliases is exited.
func (s *BaseTrinoParserListener) ExitColumnAliases(ctx *ColumnAliasesContext) {}

// EnterTableName is called when production tableName is entered.
func (s *BaseTrinoParserListener) EnterTableName(ctx *TableNameContext) {}

// ExitTableName is called when production tableName is exited.
func (s *BaseTrinoParserListener) ExitTableName(ctx *TableNameContext) {}

// EnterSubqueryRelation is called when production subqueryRelation is entered.
func (s *BaseTrinoParserListener) EnterSubqueryRelation(ctx *SubqueryRelationContext) {}

// ExitSubqueryRelation is called when production subqueryRelation is exited.
func (s *BaseTrinoParserListener) ExitSubqueryRelation(ctx *SubqueryRelationContext) {}

// EnterUnnest is called when production unnest is entered.
func (s *BaseTrinoParserListener) EnterUnnest(ctx *UnnestContext) {}

// ExitUnnest is called when production unnest is exited.
func (s *BaseTrinoParserListener) ExitUnnest(ctx *UnnestContext) {}

// EnterLateral is called when production lateral is entered.
func (s *BaseTrinoParserListener) EnterLateral(ctx *LateralContext) {}

// ExitLateral is called when production lateral is exited.
func (s *BaseTrinoParserListener) ExitLateral(ctx *LateralContext) {}

// EnterTableFunctionInvocation is called when production tableFunctionInvocation is entered.
func (s *BaseTrinoParserListener) EnterTableFunctionInvocation(ctx *TableFunctionInvocationContext) {}

// ExitTableFunctionInvocation is called when production tableFunctionInvocation is exited.
func (s *BaseTrinoParserListener) ExitTableFunctionInvocation(ctx *TableFunctionInvocationContext) {}

// EnterParenthesizedRelation is called when production parenthesizedRelation is entered.
func (s *BaseTrinoParserListener) EnterParenthesizedRelation(ctx *ParenthesizedRelationContext) {}

// ExitParenthesizedRelation is called when production parenthesizedRelation is exited.
func (s *BaseTrinoParserListener) ExitParenthesizedRelation(ctx *ParenthesizedRelationContext) {}

// EnterTableFunctionCall is called when production tableFunctionCall is entered.
func (s *BaseTrinoParserListener) EnterTableFunctionCall(ctx *TableFunctionCallContext) {}

// ExitTableFunctionCall is called when production tableFunctionCall is exited.
func (s *BaseTrinoParserListener) ExitTableFunctionCall(ctx *TableFunctionCallContext) {}

// EnterTableFunctionArgument is called when production tableFunctionArgument is entered.
func (s *BaseTrinoParserListener) EnterTableFunctionArgument(ctx *TableFunctionArgumentContext) {}

// ExitTableFunctionArgument is called when production tableFunctionArgument is exited.
func (s *BaseTrinoParserListener) ExitTableFunctionArgument(ctx *TableFunctionArgumentContext) {}

// EnterTableArgument is called when production tableArgument is entered.
func (s *BaseTrinoParserListener) EnterTableArgument(ctx *TableArgumentContext) {}

// ExitTableArgument is called when production tableArgument is exited.
func (s *BaseTrinoParserListener) ExitTableArgument(ctx *TableArgumentContext) {}

// EnterTableArgumentTable is called when production tableArgumentTable is entered.
func (s *BaseTrinoParserListener) EnterTableArgumentTable(ctx *TableArgumentTableContext) {}

// ExitTableArgumentTable is called when production tableArgumentTable is exited.
func (s *BaseTrinoParserListener) ExitTableArgumentTable(ctx *TableArgumentTableContext) {}

// EnterTableArgumentQuery is called when production tableArgumentQuery is entered.
func (s *BaseTrinoParserListener) EnterTableArgumentQuery(ctx *TableArgumentQueryContext) {}

// ExitTableArgumentQuery is called when production tableArgumentQuery is exited.
func (s *BaseTrinoParserListener) ExitTableArgumentQuery(ctx *TableArgumentQueryContext) {}

// EnterDescriptorArgument is called when production descriptorArgument is entered.
func (s *BaseTrinoParserListener) EnterDescriptorArgument(ctx *DescriptorArgumentContext) {}

// ExitDescriptorArgument is called when production descriptorArgument is exited.
func (s *BaseTrinoParserListener) ExitDescriptorArgument(ctx *DescriptorArgumentContext) {}

// EnterDescriptorField is called when production descriptorField is entered.
func (s *BaseTrinoParserListener) EnterDescriptorField(ctx *DescriptorFieldContext) {}

// ExitDescriptorField is called when production descriptorField is exited.
func (s *BaseTrinoParserListener) ExitDescriptorField(ctx *DescriptorFieldContext) {}

// EnterCopartitionTables is called when production copartitionTables is entered.
func (s *BaseTrinoParserListener) EnterCopartitionTables(ctx *CopartitionTablesContext) {}

// ExitCopartitionTables is called when production copartitionTables is exited.
func (s *BaseTrinoParserListener) ExitCopartitionTables(ctx *CopartitionTablesContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseTrinoParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseTrinoParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterLogicalNot is called when production logicalNot is entered.
func (s *BaseTrinoParserListener) EnterLogicalNot(ctx *LogicalNotContext) {}

// ExitLogicalNot is called when production logicalNot is exited.
func (s *BaseTrinoParserListener) ExitLogicalNot(ctx *LogicalNotContext) {}

// EnterPredicated is called when production predicated is entered.
func (s *BaseTrinoParserListener) EnterPredicated(ctx *PredicatedContext) {}

// ExitPredicated is called when production predicated is exited.
func (s *BaseTrinoParserListener) ExitPredicated(ctx *PredicatedContext) {}

// EnterOr is called when production or is entered.
func (s *BaseTrinoParserListener) EnterOr(ctx *OrContext) {}

// ExitOr is called when production or is exited.
func (s *BaseTrinoParserListener) ExitOr(ctx *OrContext) {}

// EnterAnd is called when production and is entered.
func (s *BaseTrinoParserListener) EnterAnd(ctx *AndContext) {}

// ExitAnd is called when production and is exited.
func (s *BaseTrinoParserListener) ExitAnd(ctx *AndContext) {}

// EnterComparison is called when production comparison is entered.
func (s *BaseTrinoParserListener) EnterComparison(ctx *ComparisonContext) {}

// ExitComparison is called when production comparison is exited.
func (s *BaseTrinoParserListener) ExitComparison(ctx *ComparisonContext) {}

// EnterQuantifiedComparison is called when production quantifiedComparison is entered.
func (s *BaseTrinoParserListener) EnterQuantifiedComparison(ctx *QuantifiedComparisonContext) {}

// ExitQuantifiedComparison is called when production quantifiedComparison is exited.
func (s *BaseTrinoParserListener) ExitQuantifiedComparison(ctx *QuantifiedComparisonContext) {}

// EnterBetween is called when production between is entered.
func (s *BaseTrinoParserListener) EnterBetween(ctx *BetweenContext) {}

// ExitBetween is called when production between is exited.
func (s *BaseTrinoParserListener) ExitBetween(ctx *BetweenContext) {}

// EnterInList is called when production inList is entered.
func (s *BaseTrinoParserListener) EnterInList(ctx *InListContext) {}

// ExitInList is called when production inList is exited.
func (s *BaseTrinoParserListener) ExitInList(ctx *InListContext) {}

// EnterInSubquery is called when production inSubquery is entered.
func (s *BaseTrinoParserListener) EnterInSubquery(ctx *InSubqueryContext) {}

// ExitInSubquery is called when production inSubquery is exited.
func (s *BaseTrinoParserListener) ExitInSubquery(ctx *InSubqueryContext) {}

// EnterLike is called when production like is entered.
func (s *BaseTrinoParserListener) EnterLike(ctx *LikeContext) {}

// ExitLike is called when production like is exited.
func (s *BaseTrinoParserListener) ExitLike(ctx *LikeContext) {}

// EnterNullPredicate is called when production nullPredicate is entered.
func (s *BaseTrinoParserListener) EnterNullPredicate(ctx *NullPredicateContext) {}

// ExitNullPredicate is called when production nullPredicate is exited.
func (s *BaseTrinoParserListener) ExitNullPredicate(ctx *NullPredicateContext) {}

// EnterDistinctFrom is called when production distinctFrom is entered.
func (s *BaseTrinoParserListener) EnterDistinctFrom(ctx *DistinctFromContext) {}

// ExitDistinctFrom is called when production distinctFrom is exited.
func (s *BaseTrinoParserListener) ExitDistinctFrom(ctx *DistinctFromContext) {}

// EnterValueExpressionDefault is called when production valueExpressionDefault is entered.
func (s *BaseTrinoParserListener) EnterValueExpressionDefault(ctx *ValueExpressionDefaultContext) {}

// ExitValueExpressionDefault is called when production valueExpressionDefault is exited.
func (s *BaseTrinoParserListener) ExitValueExpressionDefault(ctx *ValueExpressionDefaultContext) {}

// EnterConcatenation is called when production concatenation is entered.
func (s *BaseTrinoParserListener) EnterConcatenation(ctx *ConcatenationContext) {}

// ExitConcatenation is called when production concatenation is exited.
func (s *BaseTrinoParserListener) ExitConcatenation(ctx *ConcatenationContext) {}

// EnterArithmeticBinary is called when production arithmeticBinary is entered.
func (s *BaseTrinoParserListener) EnterArithmeticBinary(ctx *ArithmeticBinaryContext) {}

// ExitArithmeticBinary is called when production arithmeticBinary is exited.
func (s *BaseTrinoParserListener) ExitArithmeticBinary(ctx *ArithmeticBinaryContext) {}

// EnterArithmeticUnary is called when production arithmeticUnary is entered.
func (s *BaseTrinoParserListener) EnterArithmeticUnary(ctx *ArithmeticUnaryContext) {}

// ExitArithmeticUnary is called when production arithmeticUnary is exited.
func (s *BaseTrinoParserListener) ExitArithmeticUnary(ctx *ArithmeticUnaryContext) {}

// EnterAtTimeZone is called when production atTimeZone is entered.
func (s *BaseTrinoParserListener) EnterAtTimeZone(ctx *AtTimeZoneContext) {}

// ExitAtTimeZone is called when production atTimeZone is exited.
func (s *BaseTrinoParserListener) ExitAtTimeZone(ctx *AtTimeZoneContext) {}

// EnterDereference is called when production dereference is entered.
func (s *BaseTrinoParserListener) EnterDereference(ctx *DereferenceContext) {}

// ExitDereference is called when production dereference is exited.
func (s *BaseTrinoParserListener) ExitDereference(ctx *DereferenceContext) {}

// EnterTypeConstructor is called when production typeConstructor is entered.
func (s *BaseTrinoParserListener) EnterTypeConstructor(ctx *TypeConstructorContext) {}

// ExitTypeConstructor is called when production typeConstructor is exited.
func (s *BaseTrinoParserListener) ExitTypeConstructor(ctx *TypeConstructorContext) {}

// EnterJsonValue is called when production jsonValue is entered.
func (s *BaseTrinoParserListener) EnterJsonValue(ctx *JsonValueContext) {}

// ExitJsonValue is called when production jsonValue is exited.
func (s *BaseTrinoParserListener) ExitJsonValue(ctx *JsonValueContext) {}

// EnterSpecialDateTimeFunction is called when production specialDateTimeFunction is entered.
func (s *BaseTrinoParserListener) EnterSpecialDateTimeFunction(ctx *SpecialDateTimeFunctionContext) {}

// ExitSpecialDateTimeFunction is called when production specialDateTimeFunction is exited.
func (s *BaseTrinoParserListener) ExitSpecialDateTimeFunction(ctx *SpecialDateTimeFunctionContext) {}

// EnterSubstring is called when production substring is entered.
func (s *BaseTrinoParserListener) EnterSubstring(ctx *SubstringContext) {}

// ExitSubstring is called when production substring is exited.
func (s *BaseTrinoParserListener) ExitSubstring(ctx *SubstringContext) {}

// EnterCast is called when production cast is entered.
func (s *BaseTrinoParserListener) EnterCast(ctx *CastContext) {}

// ExitCast is called when production cast is exited.
func (s *BaseTrinoParserListener) ExitCast(ctx *CastContext) {}

// EnterLambda is called when production lambda is entered.
func (s *BaseTrinoParserListener) EnterLambda(ctx *LambdaContext) {}

// ExitLambda is called when production lambda is exited.
func (s *BaseTrinoParserListener) ExitLambda(ctx *LambdaContext) {}

// EnterParenthesizedExpression is called when production parenthesizedExpression is entered.
func (s *BaseTrinoParserListener) EnterParenthesizedExpression(ctx *ParenthesizedExpressionContext) {}

// ExitParenthesizedExpression is called when production parenthesizedExpression is exited.
func (s *BaseTrinoParserListener) ExitParenthesizedExpression(ctx *ParenthesizedExpressionContext) {}

// EnterTrim is called when production trim is entered.
func (s *BaseTrinoParserListener) EnterTrim(ctx *TrimContext) {}

// ExitTrim is called when production trim is exited.
func (s *BaseTrinoParserListener) ExitTrim(ctx *TrimContext) {}

// EnterParameter is called when production parameter is entered.
func (s *BaseTrinoParserListener) EnterParameter(ctx *ParameterContext) {}

// ExitParameter is called when production parameter is exited.
func (s *BaseTrinoParserListener) ExitParameter(ctx *ParameterContext) {}

// EnterNormalize is called when production normalize is entered.
func (s *BaseTrinoParserListener) EnterNormalize(ctx *NormalizeContext) {}

// ExitNormalize is called when production normalize is exited.
func (s *BaseTrinoParserListener) ExitNormalize(ctx *NormalizeContext) {}

// EnterJsonObject is called when production jsonObject is entered.
func (s *BaseTrinoParserListener) EnterJsonObject(ctx *JsonObjectContext) {}

// ExitJsonObject is called when production jsonObject is exited.
func (s *BaseTrinoParserListener) ExitJsonObject(ctx *JsonObjectContext) {}

// EnterIntervalLiteral is called when production intervalLiteral is entered.
func (s *BaseTrinoParserListener) EnterIntervalLiteral(ctx *IntervalLiteralContext) {}

// ExitIntervalLiteral is called when production intervalLiteral is exited.
func (s *BaseTrinoParserListener) ExitIntervalLiteral(ctx *IntervalLiteralContext) {}

// EnterNumericLiteral is called when production numericLiteral is entered.
func (s *BaseTrinoParserListener) EnterNumericLiteral(ctx *NumericLiteralContext) {}

// ExitNumericLiteral is called when production numericLiteral is exited.
func (s *BaseTrinoParserListener) ExitNumericLiteral(ctx *NumericLiteralContext) {}

// EnterBooleanLiteral is called when production booleanLiteral is entered.
func (s *BaseTrinoParserListener) EnterBooleanLiteral(ctx *BooleanLiteralContext) {}

// ExitBooleanLiteral is called when production booleanLiteral is exited.
func (s *BaseTrinoParserListener) ExitBooleanLiteral(ctx *BooleanLiteralContext) {}

// EnterJsonArray is called when production jsonArray is entered.
func (s *BaseTrinoParserListener) EnterJsonArray(ctx *JsonArrayContext) {}

// ExitJsonArray is called when production jsonArray is exited.
func (s *BaseTrinoParserListener) ExitJsonArray(ctx *JsonArrayContext) {}

// EnterSimpleCase is called when production simpleCase is entered.
func (s *BaseTrinoParserListener) EnterSimpleCase(ctx *SimpleCaseContext) {}

// ExitSimpleCase is called when production simpleCase is exited.
func (s *BaseTrinoParserListener) ExitSimpleCase(ctx *SimpleCaseContext) {}

// EnterColumnReference is called when production columnReference is entered.
func (s *BaseTrinoParserListener) EnterColumnReference(ctx *ColumnReferenceContext) {}

// ExitColumnReference is called when production columnReference is exited.
func (s *BaseTrinoParserListener) ExitColumnReference(ctx *ColumnReferenceContext) {}

// EnterNullLiteral is called when production nullLiteral is entered.
func (s *BaseTrinoParserListener) EnterNullLiteral(ctx *NullLiteralContext) {}

// ExitNullLiteral is called when production nullLiteral is exited.
func (s *BaseTrinoParserListener) ExitNullLiteral(ctx *NullLiteralContext) {}

// EnterRowConstructor is called when production rowConstructor is entered.
func (s *BaseTrinoParserListener) EnterRowConstructor(ctx *RowConstructorContext) {}

// ExitRowConstructor is called when production rowConstructor is exited.
func (s *BaseTrinoParserListener) ExitRowConstructor(ctx *RowConstructorContext) {}

// EnterSubscript is called when production subscript is entered.
func (s *BaseTrinoParserListener) EnterSubscript(ctx *SubscriptContext) {}

// ExitSubscript is called when production subscript is exited.
func (s *BaseTrinoParserListener) ExitSubscript(ctx *SubscriptContext) {}

// EnterJsonExists is called when production jsonExists is entered.
func (s *BaseTrinoParserListener) EnterJsonExists(ctx *JsonExistsContext) {}

// ExitJsonExists is called when production jsonExists is exited.
func (s *BaseTrinoParserListener) ExitJsonExists(ctx *JsonExistsContext) {}

// EnterCurrentPath is called when production currentPath is entered.
func (s *BaseTrinoParserListener) EnterCurrentPath(ctx *CurrentPathContext) {}

// ExitCurrentPath is called when production currentPath is exited.
func (s *BaseTrinoParserListener) ExitCurrentPath(ctx *CurrentPathContext) {}

// EnterSubqueryExpression is called when production subqueryExpression is entered.
func (s *BaseTrinoParserListener) EnterSubqueryExpression(ctx *SubqueryExpressionContext) {}

// ExitSubqueryExpression is called when production subqueryExpression is exited.
func (s *BaseTrinoParserListener) ExitSubqueryExpression(ctx *SubqueryExpressionContext) {}

// EnterBinaryLiteral is called when production binaryLiteral is entered.
func (s *BaseTrinoParserListener) EnterBinaryLiteral(ctx *BinaryLiteralContext) {}

// ExitBinaryLiteral is called when production binaryLiteral is exited.
func (s *BaseTrinoParserListener) ExitBinaryLiteral(ctx *BinaryLiteralContext) {}

// EnterCurrentUser is called when production currentUser is entered.
func (s *BaseTrinoParserListener) EnterCurrentUser(ctx *CurrentUserContext) {}

// ExitCurrentUser is called when production currentUser is exited.
func (s *BaseTrinoParserListener) ExitCurrentUser(ctx *CurrentUserContext) {}

// EnterJsonQuery is called when production jsonQuery is entered.
func (s *BaseTrinoParserListener) EnterJsonQuery(ctx *JsonQueryContext) {}

// ExitJsonQuery is called when production jsonQuery is exited.
func (s *BaseTrinoParserListener) ExitJsonQuery(ctx *JsonQueryContext) {}

// EnterMeasure is called when production measure is entered.
func (s *BaseTrinoParserListener) EnterMeasure(ctx *MeasureContext) {}

// ExitMeasure is called when production measure is exited.
func (s *BaseTrinoParserListener) ExitMeasure(ctx *MeasureContext) {}

// EnterExtract is called when production extract is entered.
func (s *BaseTrinoParserListener) EnterExtract(ctx *ExtractContext) {}

// ExitExtract is called when production extract is exited.
func (s *BaseTrinoParserListener) ExitExtract(ctx *ExtractContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseTrinoParserListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseTrinoParserListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterArrayConstructor is called when production arrayConstructor is entered.
func (s *BaseTrinoParserListener) EnterArrayConstructor(ctx *ArrayConstructorContext) {}

// ExitArrayConstructor is called when production arrayConstructor is exited.
func (s *BaseTrinoParserListener) ExitArrayConstructor(ctx *ArrayConstructorContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseTrinoParserListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseTrinoParserListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterCurrentSchema is called when production currentSchema is entered.
func (s *BaseTrinoParserListener) EnterCurrentSchema(ctx *CurrentSchemaContext) {}

// ExitCurrentSchema is called when production currentSchema is exited.
func (s *BaseTrinoParserListener) ExitCurrentSchema(ctx *CurrentSchemaContext) {}

// EnterExists is called when production exists is entered.
func (s *BaseTrinoParserListener) EnterExists(ctx *ExistsContext) {}

// ExitExists is called when production exists is exited.
func (s *BaseTrinoParserListener) ExitExists(ctx *ExistsContext) {}

// EnterPosition is called when production position is entered.
func (s *BaseTrinoParserListener) EnterPosition(ctx *PositionContext) {}

// ExitPosition is called when production position is exited.
func (s *BaseTrinoParserListener) ExitPosition(ctx *PositionContext) {}

// EnterListagg is called when production listagg is entered.
func (s *BaseTrinoParserListener) EnterListagg(ctx *ListaggContext) {}

// ExitListagg is called when production listagg is exited.
func (s *BaseTrinoParserListener) ExitListagg(ctx *ListaggContext) {}

// EnterSearchedCase is called when production searchedCase is entered.
func (s *BaseTrinoParserListener) EnterSearchedCase(ctx *SearchedCaseContext) {}

// ExitSearchedCase is called when production searchedCase is exited.
func (s *BaseTrinoParserListener) ExitSearchedCase(ctx *SearchedCaseContext) {}

// EnterCurrentCatalog is called when production currentCatalog is entered.
func (s *BaseTrinoParserListener) EnterCurrentCatalog(ctx *CurrentCatalogContext) {}

// ExitCurrentCatalog is called when production currentCatalog is exited.
func (s *BaseTrinoParserListener) ExitCurrentCatalog(ctx *CurrentCatalogContext) {}

// EnterGroupingOperation is called when production groupingOperation is entered.
func (s *BaseTrinoParserListener) EnterGroupingOperation(ctx *GroupingOperationContext) {}

// ExitGroupingOperation is called when production groupingOperation is exited.
func (s *BaseTrinoParserListener) ExitGroupingOperation(ctx *GroupingOperationContext) {}

// EnterJsonPathInvocation is called when production jsonPathInvocation is entered.
func (s *BaseTrinoParserListener) EnterJsonPathInvocation(ctx *JsonPathInvocationContext) {}

// ExitJsonPathInvocation is called when production jsonPathInvocation is exited.
func (s *BaseTrinoParserListener) ExitJsonPathInvocation(ctx *JsonPathInvocationContext) {}

// EnterJsonValueExpression is called when production jsonValueExpression is entered.
func (s *BaseTrinoParserListener) EnterJsonValueExpression(ctx *JsonValueExpressionContext) {}

// ExitJsonValueExpression is called when production jsonValueExpression is exited.
func (s *BaseTrinoParserListener) ExitJsonValueExpression(ctx *JsonValueExpressionContext) {}

// EnterJsonRepresentation is called when production jsonRepresentation is entered.
func (s *BaseTrinoParserListener) EnterJsonRepresentation(ctx *JsonRepresentationContext) {}

// ExitJsonRepresentation is called when production jsonRepresentation is exited.
func (s *BaseTrinoParserListener) ExitJsonRepresentation(ctx *JsonRepresentationContext) {}

// EnterJsonArgument is called when production jsonArgument is entered.
func (s *BaseTrinoParserListener) EnterJsonArgument(ctx *JsonArgumentContext) {}

// ExitJsonArgument is called when production jsonArgument is exited.
func (s *BaseTrinoParserListener) ExitJsonArgument(ctx *JsonArgumentContext) {}

// EnterJsonExistsErrorBehavior is called when production jsonExistsErrorBehavior is entered.
func (s *BaseTrinoParserListener) EnterJsonExistsErrorBehavior(ctx *JsonExistsErrorBehaviorContext) {}

// ExitJsonExistsErrorBehavior is called when production jsonExistsErrorBehavior is exited.
func (s *BaseTrinoParserListener) ExitJsonExistsErrorBehavior(ctx *JsonExistsErrorBehaviorContext) {}

// EnterJsonValueBehavior is called when production jsonValueBehavior is entered.
func (s *BaseTrinoParserListener) EnterJsonValueBehavior(ctx *JsonValueBehaviorContext) {}

// ExitJsonValueBehavior is called when production jsonValueBehavior is exited.
func (s *BaseTrinoParserListener) ExitJsonValueBehavior(ctx *JsonValueBehaviorContext) {}

// EnterJsonQueryWrapperBehavior is called when production jsonQueryWrapperBehavior is entered.
func (s *BaseTrinoParserListener) EnterJsonQueryWrapperBehavior(ctx *JsonQueryWrapperBehaviorContext) {
}

// ExitJsonQueryWrapperBehavior is called when production jsonQueryWrapperBehavior is exited.
func (s *BaseTrinoParserListener) ExitJsonQueryWrapperBehavior(ctx *JsonQueryWrapperBehaviorContext) {
}

// EnterJsonQueryBehavior is called when production jsonQueryBehavior is entered.
func (s *BaseTrinoParserListener) EnterJsonQueryBehavior(ctx *JsonQueryBehaviorContext) {}

// ExitJsonQueryBehavior is called when production jsonQueryBehavior is exited.
func (s *BaseTrinoParserListener) ExitJsonQueryBehavior(ctx *JsonQueryBehaviorContext) {}

// EnterJsonObjectMember is called when production jsonObjectMember is entered.
func (s *BaseTrinoParserListener) EnterJsonObjectMember(ctx *JsonObjectMemberContext) {}

// ExitJsonObjectMember is called when production jsonObjectMember is exited.
func (s *BaseTrinoParserListener) ExitJsonObjectMember(ctx *JsonObjectMemberContext) {}

// EnterProcessingMode is called when production processingMode is entered.
func (s *BaseTrinoParserListener) EnterProcessingMode(ctx *ProcessingModeContext) {}

// ExitProcessingMode is called when production processingMode is exited.
func (s *BaseTrinoParserListener) ExitProcessingMode(ctx *ProcessingModeContext) {}

// EnterNullTreatment is called when production nullTreatment is entered.
func (s *BaseTrinoParserListener) EnterNullTreatment(ctx *NullTreatmentContext) {}

// ExitNullTreatment is called when production nullTreatment is exited.
func (s *BaseTrinoParserListener) ExitNullTreatment(ctx *NullTreatmentContext) {}

// EnterBasicStringLiteral is called when production basicStringLiteral is entered.
func (s *BaseTrinoParserListener) EnterBasicStringLiteral(ctx *BasicStringLiteralContext) {}

// ExitBasicStringLiteral is called when production basicStringLiteral is exited.
func (s *BaseTrinoParserListener) ExitBasicStringLiteral(ctx *BasicStringLiteralContext) {}

// EnterUnicodeStringLiteral is called when production unicodeStringLiteral is entered.
func (s *BaseTrinoParserListener) EnterUnicodeStringLiteral(ctx *UnicodeStringLiteralContext) {}

// ExitUnicodeStringLiteral is called when production unicodeStringLiteral is exited.
func (s *BaseTrinoParserListener) ExitUnicodeStringLiteral(ctx *UnicodeStringLiteralContext) {}

// EnterTimeZoneInterval is called when production timeZoneInterval is entered.
func (s *BaseTrinoParserListener) EnterTimeZoneInterval(ctx *TimeZoneIntervalContext) {}

// ExitTimeZoneInterval is called when production timeZoneInterval is exited.
func (s *BaseTrinoParserListener) ExitTimeZoneInterval(ctx *TimeZoneIntervalContext) {}

// EnterTimeZoneString is called when production timeZoneString is entered.
func (s *BaseTrinoParserListener) EnterTimeZoneString(ctx *TimeZoneStringContext) {}

// ExitTimeZoneString is called when production timeZoneString is exited.
func (s *BaseTrinoParserListener) ExitTimeZoneString(ctx *TimeZoneStringContext) {}

// EnterComparisonOperator is called when production comparisonOperator is entered.
func (s *BaseTrinoParserListener) EnterComparisonOperator(ctx *ComparisonOperatorContext) {}

// ExitComparisonOperator is called when production comparisonOperator is exited.
func (s *BaseTrinoParserListener) ExitComparisonOperator(ctx *ComparisonOperatorContext) {}

// EnterComparisonQuantifier is called when production comparisonQuantifier is entered.
func (s *BaseTrinoParserListener) EnterComparisonQuantifier(ctx *ComparisonQuantifierContext) {}

// ExitComparisonQuantifier is called when production comparisonQuantifier is exited.
func (s *BaseTrinoParserListener) ExitComparisonQuantifier(ctx *ComparisonQuantifierContext) {}

// EnterBooleanValue is called when production booleanValue is entered.
func (s *BaseTrinoParserListener) EnterBooleanValue(ctx *BooleanValueContext) {}

// ExitBooleanValue is called when production booleanValue is exited.
func (s *BaseTrinoParserListener) ExitBooleanValue(ctx *BooleanValueContext) {}

// EnterInterval is called when production interval is entered.
func (s *BaseTrinoParserListener) EnterInterval(ctx *IntervalContext) {}

// ExitInterval is called when production interval is exited.
func (s *BaseTrinoParserListener) ExitInterval(ctx *IntervalContext) {}

// EnterIntervalField is called when production intervalField is entered.
func (s *BaseTrinoParserListener) EnterIntervalField(ctx *IntervalFieldContext) {}

// ExitIntervalField is called when production intervalField is exited.
func (s *BaseTrinoParserListener) ExitIntervalField(ctx *IntervalFieldContext) {}

// EnterNormalForm is called when production normalForm is entered.
func (s *BaseTrinoParserListener) EnterNormalForm(ctx *NormalFormContext) {}

// ExitNormalForm is called when production normalForm is exited.
func (s *BaseTrinoParserListener) ExitNormalForm(ctx *NormalFormContext) {}

// EnterRowType is called when production rowType is entered.
func (s *BaseTrinoParserListener) EnterRowType(ctx *RowTypeContext) {}

// ExitRowType is called when production rowType is exited.
func (s *BaseTrinoParserListener) ExitRowType(ctx *RowTypeContext) {}

// EnterIntervalType is called when production intervalType is entered.
func (s *BaseTrinoParserListener) EnterIntervalType(ctx *IntervalTypeContext) {}

// ExitIntervalType is called when production intervalType is exited.
func (s *BaseTrinoParserListener) ExitIntervalType(ctx *IntervalTypeContext) {}

// EnterArrayType is called when production arrayType is entered.
func (s *BaseTrinoParserListener) EnterArrayType(ctx *ArrayTypeContext) {}

// ExitArrayType is called when production arrayType is exited.
func (s *BaseTrinoParserListener) ExitArrayType(ctx *ArrayTypeContext) {}

// EnterDoublePrecisionType is called when production doublePrecisionType is entered.
func (s *BaseTrinoParserListener) EnterDoublePrecisionType(ctx *DoublePrecisionTypeContext) {}

// ExitDoublePrecisionType is called when production doublePrecisionType is exited.
func (s *BaseTrinoParserListener) ExitDoublePrecisionType(ctx *DoublePrecisionTypeContext) {}

// EnterLegacyArrayType is called when production legacyArrayType is entered.
func (s *BaseTrinoParserListener) EnterLegacyArrayType(ctx *LegacyArrayTypeContext) {}

// ExitLegacyArrayType is called when production legacyArrayType is exited.
func (s *BaseTrinoParserListener) ExitLegacyArrayType(ctx *LegacyArrayTypeContext) {}

// EnterGenericType is called when production genericType is entered.
func (s *BaseTrinoParserListener) EnterGenericType(ctx *GenericTypeContext) {}

// ExitGenericType is called when production genericType is exited.
func (s *BaseTrinoParserListener) ExitGenericType(ctx *GenericTypeContext) {}

// EnterDateTimeType is called when production dateTimeType is entered.
func (s *BaseTrinoParserListener) EnterDateTimeType(ctx *DateTimeTypeContext) {}

// ExitDateTimeType is called when production dateTimeType is exited.
func (s *BaseTrinoParserListener) ExitDateTimeType(ctx *DateTimeTypeContext) {}

// EnterLegacyMapType is called when production legacyMapType is entered.
func (s *BaseTrinoParserListener) EnterLegacyMapType(ctx *LegacyMapTypeContext) {}

// ExitLegacyMapType is called when production legacyMapType is exited.
func (s *BaseTrinoParserListener) ExitLegacyMapType(ctx *LegacyMapTypeContext) {}

// EnterRowField is called when production rowField is entered.
func (s *BaseTrinoParserListener) EnterRowField(ctx *RowFieldContext) {}

// ExitRowField is called when production rowField is exited.
func (s *BaseTrinoParserListener) ExitRowField(ctx *RowFieldContext) {}

// EnterTypeParameter is called when production typeParameter is entered.
func (s *BaseTrinoParserListener) EnterTypeParameter(ctx *TypeParameterContext) {}

// ExitTypeParameter is called when production typeParameter is exited.
func (s *BaseTrinoParserListener) ExitTypeParameter(ctx *TypeParameterContext) {}

// EnterWhenClause is called when production whenClause is entered.
func (s *BaseTrinoParserListener) EnterWhenClause(ctx *WhenClauseContext) {}

// ExitWhenClause is called when production whenClause is exited.
func (s *BaseTrinoParserListener) ExitWhenClause(ctx *WhenClauseContext) {}

// EnterFilter is called when production filter is entered.
func (s *BaseTrinoParserListener) EnterFilter(ctx *FilterContext) {}

// ExitFilter is called when production filter is exited.
func (s *BaseTrinoParserListener) ExitFilter(ctx *FilterContext) {}

// EnterMergeUpdate is called when production mergeUpdate is entered.
func (s *BaseTrinoParserListener) EnterMergeUpdate(ctx *MergeUpdateContext) {}

// ExitMergeUpdate is called when production mergeUpdate is exited.
func (s *BaseTrinoParserListener) ExitMergeUpdate(ctx *MergeUpdateContext) {}

// EnterMergeDelete is called when production mergeDelete is entered.
func (s *BaseTrinoParserListener) EnterMergeDelete(ctx *MergeDeleteContext) {}

// ExitMergeDelete is called when production mergeDelete is exited.
func (s *BaseTrinoParserListener) ExitMergeDelete(ctx *MergeDeleteContext) {}

// EnterMergeInsert is called when production mergeInsert is entered.
func (s *BaseTrinoParserListener) EnterMergeInsert(ctx *MergeInsertContext) {}

// ExitMergeInsert is called when production mergeInsert is exited.
func (s *BaseTrinoParserListener) ExitMergeInsert(ctx *MergeInsertContext) {}

// EnterOver is called when production over is entered.
func (s *BaseTrinoParserListener) EnterOver(ctx *OverContext) {}

// ExitOver is called when production over is exited.
func (s *BaseTrinoParserListener) ExitOver(ctx *OverContext) {}

// EnterWindowFrame is called when production windowFrame is entered.
func (s *BaseTrinoParserListener) EnterWindowFrame(ctx *WindowFrameContext) {}

// ExitWindowFrame is called when production windowFrame is exited.
func (s *BaseTrinoParserListener) ExitWindowFrame(ctx *WindowFrameContext) {}

// EnterFrameExtent is called when production frameExtent is entered.
func (s *BaseTrinoParserListener) EnterFrameExtent(ctx *FrameExtentContext) {}

// ExitFrameExtent is called when production frameExtent is exited.
func (s *BaseTrinoParserListener) ExitFrameExtent(ctx *FrameExtentContext) {}

// EnterUnboundedFrame is called when production unboundedFrame is entered.
func (s *BaseTrinoParserListener) EnterUnboundedFrame(ctx *UnboundedFrameContext) {}

// ExitUnboundedFrame is called when production unboundedFrame is exited.
func (s *BaseTrinoParserListener) ExitUnboundedFrame(ctx *UnboundedFrameContext) {}

// EnterCurrentRowBound is called when production currentRowBound is entered.
func (s *BaseTrinoParserListener) EnterCurrentRowBound(ctx *CurrentRowBoundContext) {}

// ExitCurrentRowBound is called when production currentRowBound is exited.
func (s *BaseTrinoParserListener) ExitCurrentRowBound(ctx *CurrentRowBoundContext) {}

// EnterBoundedFrame is called when production boundedFrame is entered.
func (s *BaseTrinoParserListener) EnterBoundedFrame(ctx *BoundedFrameContext) {}

// ExitBoundedFrame is called when production boundedFrame is exited.
func (s *BaseTrinoParserListener) ExitBoundedFrame(ctx *BoundedFrameContext) {}

// EnterQuantifiedPrimary is called when production quantifiedPrimary is entered.
func (s *BaseTrinoParserListener) EnterQuantifiedPrimary(ctx *QuantifiedPrimaryContext) {}

// ExitQuantifiedPrimary is called when production quantifiedPrimary is exited.
func (s *BaseTrinoParserListener) ExitQuantifiedPrimary(ctx *QuantifiedPrimaryContext) {}

// EnterPatternConcatenation is called when production patternConcatenation is entered.
func (s *BaseTrinoParserListener) EnterPatternConcatenation(ctx *PatternConcatenationContext) {}

// ExitPatternConcatenation is called when production patternConcatenation is exited.
func (s *BaseTrinoParserListener) ExitPatternConcatenation(ctx *PatternConcatenationContext) {}

// EnterPatternAlternation is called when production patternAlternation is entered.
func (s *BaseTrinoParserListener) EnterPatternAlternation(ctx *PatternAlternationContext) {}

// ExitPatternAlternation is called when production patternAlternation is exited.
func (s *BaseTrinoParserListener) ExitPatternAlternation(ctx *PatternAlternationContext) {}

// EnterPatternVariable is called when production patternVariable is entered.
func (s *BaseTrinoParserListener) EnterPatternVariable(ctx *PatternVariableContext) {}

// ExitPatternVariable is called when production patternVariable is exited.
func (s *BaseTrinoParserListener) ExitPatternVariable(ctx *PatternVariableContext) {}

// EnterEmptyPattern is called when production emptyPattern is entered.
func (s *BaseTrinoParserListener) EnterEmptyPattern(ctx *EmptyPatternContext) {}

// ExitEmptyPattern is called when production emptyPattern is exited.
func (s *BaseTrinoParserListener) ExitEmptyPattern(ctx *EmptyPatternContext) {}

// EnterPatternPermutation is called when production patternPermutation is entered.
func (s *BaseTrinoParserListener) EnterPatternPermutation(ctx *PatternPermutationContext) {}

// ExitPatternPermutation is called when production patternPermutation is exited.
func (s *BaseTrinoParserListener) ExitPatternPermutation(ctx *PatternPermutationContext) {}

// EnterGroupedPattern is called when production groupedPattern is entered.
func (s *BaseTrinoParserListener) EnterGroupedPattern(ctx *GroupedPatternContext) {}

// ExitGroupedPattern is called when production groupedPattern is exited.
func (s *BaseTrinoParserListener) ExitGroupedPattern(ctx *GroupedPatternContext) {}

// EnterPartitionStartAnchor is called when production partitionStartAnchor is entered.
func (s *BaseTrinoParserListener) EnterPartitionStartAnchor(ctx *PartitionStartAnchorContext) {}

// ExitPartitionStartAnchor is called when production partitionStartAnchor is exited.
func (s *BaseTrinoParserListener) ExitPartitionStartAnchor(ctx *PartitionStartAnchorContext) {}

// EnterPartitionEndAnchor is called when production partitionEndAnchor is entered.
func (s *BaseTrinoParserListener) EnterPartitionEndAnchor(ctx *PartitionEndAnchorContext) {}

// ExitPartitionEndAnchor is called when production partitionEndAnchor is exited.
func (s *BaseTrinoParserListener) ExitPartitionEndAnchor(ctx *PartitionEndAnchorContext) {}

// EnterExcludedPattern is called when production excludedPattern is entered.
func (s *BaseTrinoParserListener) EnterExcludedPattern(ctx *ExcludedPatternContext) {}

// ExitExcludedPattern is called when production excludedPattern is exited.
func (s *BaseTrinoParserListener) ExitExcludedPattern(ctx *ExcludedPatternContext) {}

// EnterZeroOrMoreQuantifier is called when production zeroOrMoreQuantifier is entered.
func (s *BaseTrinoParserListener) EnterZeroOrMoreQuantifier(ctx *ZeroOrMoreQuantifierContext) {}

// ExitZeroOrMoreQuantifier is called when production zeroOrMoreQuantifier is exited.
func (s *BaseTrinoParserListener) ExitZeroOrMoreQuantifier(ctx *ZeroOrMoreQuantifierContext) {}

// EnterOneOrMoreQuantifier is called when production oneOrMoreQuantifier is entered.
func (s *BaseTrinoParserListener) EnterOneOrMoreQuantifier(ctx *OneOrMoreQuantifierContext) {}

// ExitOneOrMoreQuantifier is called when production oneOrMoreQuantifier is exited.
func (s *BaseTrinoParserListener) ExitOneOrMoreQuantifier(ctx *OneOrMoreQuantifierContext) {}

// EnterZeroOrOneQuantifier is called when production zeroOrOneQuantifier is entered.
func (s *BaseTrinoParserListener) EnterZeroOrOneQuantifier(ctx *ZeroOrOneQuantifierContext) {}

// ExitZeroOrOneQuantifier is called when production zeroOrOneQuantifier is exited.
func (s *BaseTrinoParserListener) ExitZeroOrOneQuantifier(ctx *ZeroOrOneQuantifierContext) {}

// EnterRangeQuantifier is called when production rangeQuantifier is entered.
func (s *BaseTrinoParserListener) EnterRangeQuantifier(ctx *RangeQuantifierContext) {}

// ExitRangeQuantifier is called when production rangeQuantifier is exited.
func (s *BaseTrinoParserListener) ExitRangeQuantifier(ctx *RangeQuantifierContext) {}

// EnterUpdateAssignment is called when production updateAssignment is entered.
func (s *BaseTrinoParserListener) EnterUpdateAssignment(ctx *UpdateAssignmentContext) {}

// ExitUpdateAssignment is called when production updateAssignment is exited.
func (s *BaseTrinoParserListener) ExitUpdateAssignment(ctx *UpdateAssignmentContext) {}

// EnterExplainFormat is called when production explainFormat is entered.
func (s *BaseTrinoParserListener) EnterExplainFormat(ctx *ExplainFormatContext) {}

// ExitExplainFormat is called when production explainFormat is exited.
func (s *BaseTrinoParserListener) ExitExplainFormat(ctx *ExplainFormatContext) {}

// EnterExplainType is called when production explainType is entered.
func (s *BaseTrinoParserListener) EnterExplainType(ctx *ExplainTypeContext) {}

// ExitExplainType is called when production explainType is exited.
func (s *BaseTrinoParserListener) ExitExplainType(ctx *ExplainTypeContext) {}

// EnterIsolationLevel is called when production isolationLevel is entered.
func (s *BaseTrinoParserListener) EnterIsolationLevel(ctx *IsolationLevelContext) {}

// ExitIsolationLevel is called when production isolationLevel is exited.
func (s *BaseTrinoParserListener) ExitIsolationLevel(ctx *IsolationLevelContext) {}

// EnterTransactionAccessMode is called when production transactionAccessMode is entered.
func (s *BaseTrinoParserListener) EnterTransactionAccessMode(ctx *TransactionAccessModeContext) {}

// ExitTransactionAccessMode is called when production transactionAccessMode is exited.
func (s *BaseTrinoParserListener) ExitTransactionAccessMode(ctx *TransactionAccessModeContext) {}

// EnterReadUncommitted is called when production readUncommitted is entered.
func (s *BaseTrinoParserListener) EnterReadUncommitted(ctx *ReadUncommittedContext) {}

// ExitReadUncommitted is called when production readUncommitted is exited.
func (s *BaseTrinoParserListener) ExitReadUncommitted(ctx *ReadUncommittedContext) {}

// EnterReadCommitted is called when production readCommitted is entered.
func (s *BaseTrinoParserListener) EnterReadCommitted(ctx *ReadCommittedContext) {}

// ExitReadCommitted is called when production readCommitted is exited.
func (s *BaseTrinoParserListener) ExitReadCommitted(ctx *ReadCommittedContext) {}

// EnterRepeatableRead is called when production repeatableRead is entered.
func (s *BaseTrinoParserListener) EnterRepeatableRead(ctx *RepeatableReadContext) {}

// ExitRepeatableRead is called when production repeatableRead is exited.
func (s *BaseTrinoParserListener) ExitRepeatableRead(ctx *RepeatableReadContext) {}

// EnterSerializable is called when production serializable is entered.
func (s *BaseTrinoParserListener) EnterSerializable(ctx *SerializableContext) {}

// ExitSerializable is called when production serializable is exited.
func (s *BaseTrinoParserListener) ExitSerializable(ctx *SerializableContext) {}

// EnterPositionalArgument is called when production positionalArgument is entered.
func (s *BaseTrinoParserListener) EnterPositionalArgument(ctx *PositionalArgumentContext) {}

// ExitPositionalArgument is called when production positionalArgument is exited.
func (s *BaseTrinoParserListener) ExitPositionalArgument(ctx *PositionalArgumentContext) {}

// EnterNamedArgument is called when production namedArgument is entered.
func (s *BaseTrinoParserListener) EnterNamedArgument(ctx *NamedArgumentContext) {}

// ExitNamedArgument is called when production namedArgument is exited.
func (s *BaseTrinoParserListener) ExitNamedArgument(ctx *NamedArgumentContext) {}

// EnterQualifiedArgument is called when production qualifiedArgument is entered.
func (s *BaseTrinoParserListener) EnterQualifiedArgument(ctx *QualifiedArgumentContext) {}

// ExitQualifiedArgument is called when production qualifiedArgument is exited.
func (s *BaseTrinoParserListener) ExitQualifiedArgument(ctx *QualifiedArgumentContext) {}

// EnterUnqualifiedArgument is called when production unqualifiedArgument is entered.
func (s *BaseTrinoParserListener) EnterUnqualifiedArgument(ctx *UnqualifiedArgumentContext) {}

// ExitUnqualifiedArgument is called when production unqualifiedArgument is exited.
func (s *BaseTrinoParserListener) ExitUnqualifiedArgument(ctx *UnqualifiedArgumentContext) {}

// EnterPathSpecification is called when production pathSpecification is entered.
func (s *BaseTrinoParserListener) EnterPathSpecification(ctx *PathSpecificationContext) {}

// ExitPathSpecification is called when production pathSpecification is exited.
func (s *BaseTrinoParserListener) ExitPathSpecification(ctx *PathSpecificationContext) {}

// EnterFunctionSpecification is called when production functionSpecification is entered.
func (s *BaseTrinoParserListener) EnterFunctionSpecification(ctx *FunctionSpecificationContext) {}

// ExitFunctionSpecification is called when production functionSpecification is exited.
func (s *BaseTrinoParserListener) ExitFunctionSpecification(ctx *FunctionSpecificationContext) {}

// EnterFunctionDeclaration is called when production functionDeclaration is entered.
func (s *BaseTrinoParserListener) EnterFunctionDeclaration(ctx *FunctionDeclarationContext) {}

// ExitFunctionDeclaration is called when production functionDeclaration is exited.
func (s *BaseTrinoParserListener) ExitFunctionDeclaration(ctx *FunctionDeclarationContext) {}

// EnterParameterDeclaration is called when production parameterDeclaration is entered.
func (s *BaseTrinoParserListener) EnterParameterDeclaration(ctx *ParameterDeclarationContext) {}

// ExitParameterDeclaration is called when production parameterDeclaration is exited.
func (s *BaseTrinoParserListener) ExitParameterDeclaration(ctx *ParameterDeclarationContext) {}

// EnterReturnsClause is called when production returnsClause is entered.
func (s *BaseTrinoParserListener) EnterReturnsClause(ctx *ReturnsClauseContext) {}

// ExitReturnsClause is called when production returnsClause is exited.
func (s *BaseTrinoParserListener) ExitReturnsClause(ctx *ReturnsClauseContext) {}

// EnterLanguageCharacteristic is called when production languageCharacteristic is entered.
func (s *BaseTrinoParserListener) EnterLanguageCharacteristic(ctx *LanguageCharacteristicContext) {}

// ExitLanguageCharacteristic is called when production languageCharacteristic is exited.
func (s *BaseTrinoParserListener) ExitLanguageCharacteristic(ctx *LanguageCharacteristicContext) {}

// EnterDeterministicCharacteristic is called when production deterministicCharacteristic is entered.
func (s *BaseTrinoParserListener) EnterDeterministicCharacteristic(ctx *DeterministicCharacteristicContext) {
}

// ExitDeterministicCharacteristic is called when production deterministicCharacteristic is exited.
func (s *BaseTrinoParserListener) ExitDeterministicCharacteristic(ctx *DeterministicCharacteristicContext) {
}

// EnterReturnsNullOnNullInputCharacteristic is called when production returnsNullOnNullInputCharacteristic is entered.
func (s *BaseTrinoParserListener) EnterReturnsNullOnNullInputCharacteristic(ctx *ReturnsNullOnNullInputCharacteristicContext) {
}

// ExitReturnsNullOnNullInputCharacteristic is called when production returnsNullOnNullInputCharacteristic is exited.
func (s *BaseTrinoParserListener) ExitReturnsNullOnNullInputCharacteristic(ctx *ReturnsNullOnNullInputCharacteristicContext) {
}

// EnterCalledOnNullInputCharacteristic is called when production calledOnNullInputCharacteristic is entered.
func (s *BaseTrinoParserListener) EnterCalledOnNullInputCharacteristic(ctx *CalledOnNullInputCharacteristicContext) {
}

// ExitCalledOnNullInputCharacteristic is called when production calledOnNullInputCharacteristic is exited.
func (s *BaseTrinoParserListener) ExitCalledOnNullInputCharacteristic(ctx *CalledOnNullInputCharacteristicContext) {
}

// EnterSecurityCharacteristic is called when production securityCharacteristic is entered.
func (s *BaseTrinoParserListener) EnterSecurityCharacteristic(ctx *SecurityCharacteristicContext) {}

// ExitSecurityCharacteristic is called when production securityCharacteristic is exited.
func (s *BaseTrinoParserListener) ExitSecurityCharacteristic(ctx *SecurityCharacteristicContext) {}

// EnterCommentCharacteristic is called when production commentCharacteristic is entered.
func (s *BaseTrinoParserListener) EnterCommentCharacteristic(ctx *CommentCharacteristicContext) {}

// ExitCommentCharacteristic is called when production commentCharacteristic is exited.
func (s *BaseTrinoParserListener) ExitCommentCharacteristic(ctx *CommentCharacteristicContext) {}

// EnterReturnStatement is called when production returnStatement is entered.
func (s *BaseTrinoParserListener) EnterReturnStatement(ctx *ReturnStatementContext) {}

// ExitReturnStatement is called when production returnStatement is exited.
func (s *BaseTrinoParserListener) ExitReturnStatement(ctx *ReturnStatementContext) {}

// EnterAssignmentStatement is called when production assignmentStatement is entered.
func (s *BaseTrinoParserListener) EnterAssignmentStatement(ctx *AssignmentStatementContext) {}

// ExitAssignmentStatement is called when production assignmentStatement is exited.
func (s *BaseTrinoParserListener) ExitAssignmentStatement(ctx *AssignmentStatementContext) {}

// EnterSimpleCaseStatement is called when production simpleCaseStatement is entered.
func (s *BaseTrinoParserListener) EnterSimpleCaseStatement(ctx *SimpleCaseStatementContext) {}

// ExitSimpleCaseStatement is called when production simpleCaseStatement is exited.
func (s *BaseTrinoParserListener) ExitSimpleCaseStatement(ctx *SimpleCaseStatementContext) {}

// EnterSearchedCaseStatement is called when production searchedCaseStatement is entered.
func (s *BaseTrinoParserListener) EnterSearchedCaseStatement(ctx *SearchedCaseStatementContext) {}

// ExitSearchedCaseStatement is called when production searchedCaseStatement is exited.
func (s *BaseTrinoParserListener) ExitSearchedCaseStatement(ctx *SearchedCaseStatementContext) {}

// EnterIfStatement is called when production ifStatement is entered.
func (s *BaseTrinoParserListener) EnterIfStatement(ctx *IfStatementContext) {}

// ExitIfStatement is called when production ifStatement is exited.
func (s *BaseTrinoParserListener) ExitIfStatement(ctx *IfStatementContext) {}

// EnterIterateStatement is called when production iterateStatement is entered.
func (s *BaseTrinoParserListener) EnterIterateStatement(ctx *IterateStatementContext) {}

// ExitIterateStatement is called when production iterateStatement is exited.
func (s *BaseTrinoParserListener) ExitIterateStatement(ctx *IterateStatementContext) {}

// EnterLeaveStatement is called when production leaveStatement is entered.
func (s *BaseTrinoParserListener) EnterLeaveStatement(ctx *LeaveStatementContext) {}

// ExitLeaveStatement is called when production leaveStatement is exited.
func (s *BaseTrinoParserListener) ExitLeaveStatement(ctx *LeaveStatementContext) {}

// EnterCompoundStatement is called when production compoundStatement is entered.
func (s *BaseTrinoParserListener) EnterCompoundStatement(ctx *CompoundStatementContext) {}

// ExitCompoundStatement is called when production compoundStatement is exited.
func (s *BaseTrinoParserListener) ExitCompoundStatement(ctx *CompoundStatementContext) {}

// EnterLoopStatement is called when production loopStatement is entered.
func (s *BaseTrinoParserListener) EnterLoopStatement(ctx *LoopStatementContext) {}

// ExitLoopStatement is called when production loopStatement is exited.
func (s *BaseTrinoParserListener) ExitLoopStatement(ctx *LoopStatementContext) {}

// EnterWhileStatement is called when production whileStatement is entered.
func (s *BaseTrinoParserListener) EnterWhileStatement(ctx *WhileStatementContext) {}

// ExitWhileStatement is called when production whileStatement is exited.
func (s *BaseTrinoParserListener) ExitWhileStatement(ctx *WhileStatementContext) {}

// EnterRepeatStatement is called when production repeatStatement is entered.
func (s *BaseTrinoParserListener) EnterRepeatStatement(ctx *RepeatStatementContext) {}

// ExitRepeatStatement is called when production repeatStatement is exited.
func (s *BaseTrinoParserListener) ExitRepeatStatement(ctx *RepeatStatementContext) {}

// EnterCaseStatementWhenClause is called when production caseStatementWhenClause is entered.
func (s *BaseTrinoParserListener) EnterCaseStatementWhenClause(ctx *CaseStatementWhenClauseContext) {}

// ExitCaseStatementWhenClause is called when production caseStatementWhenClause is exited.
func (s *BaseTrinoParserListener) ExitCaseStatementWhenClause(ctx *CaseStatementWhenClauseContext) {}

// EnterElseIfClause is called when production elseIfClause is entered.
func (s *BaseTrinoParserListener) EnterElseIfClause(ctx *ElseIfClauseContext) {}

// ExitElseIfClause is called when production elseIfClause is exited.
func (s *BaseTrinoParserListener) ExitElseIfClause(ctx *ElseIfClauseContext) {}

// EnterElseClause is called when production elseClause is entered.
func (s *BaseTrinoParserListener) EnterElseClause(ctx *ElseClauseContext) {}

// ExitElseClause is called when production elseClause is exited.
func (s *BaseTrinoParserListener) ExitElseClause(ctx *ElseClauseContext) {}

// EnterVariableDeclaration is called when production variableDeclaration is entered.
func (s *BaseTrinoParserListener) EnterVariableDeclaration(ctx *VariableDeclarationContext) {}

// ExitVariableDeclaration is called when production variableDeclaration is exited.
func (s *BaseTrinoParserListener) ExitVariableDeclaration(ctx *VariableDeclarationContext) {}

// EnterSqlStatementList is called when production sqlStatementList is entered.
func (s *BaseTrinoParserListener) EnterSqlStatementList(ctx *SqlStatementListContext) {}

// ExitSqlStatementList is called when production sqlStatementList is exited.
func (s *BaseTrinoParserListener) ExitSqlStatementList(ctx *SqlStatementListContext) {}

// EnterPrivilege is called when production privilege is entered.
func (s *BaseTrinoParserListener) EnterPrivilege(ctx *PrivilegeContext) {}

// ExitPrivilege is called when production privilege is exited.
func (s *BaseTrinoParserListener) ExitPrivilege(ctx *PrivilegeContext) {}

// EnterQualifiedName is called when production qualifiedName is entered.
func (s *BaseTrinoParserListener) EnterQualifiedName(ctx *QualifiedNameContext) {}

// ExitQualifiedName is called when production qualifiedName is exited.
func (s *BaseTrinoParserListener) ExitQualifiedName(ctx *QualifiedNameContext) {}

// EnterQueryPeriod is called when production queryPeriod is entered.
func (s *BaseTrinoParserListener) EnterQueryPeriod(ctx *QueryPeriodContext) {}

// ExitQueryPeriod is called when production queryPeriod is exited.
func (s *BaseTrinoParserListener) ExitQueryPeriod(ctx *QueryPeriodContext) {}

// EnterRangeType is called when production rangeType is entered.
func (s *BaseTrinoParserListener) EnterRangeType(ctx *RangeTypeContext) {}

// ExitRangeType is called when production rangeType is exited.
func (s *BaseTrinoParserListener) ExitRangeType(ctx *RangeTypeContext) {}

// EnterSpecifiedPrincipal is called when production specifiedPrincipal is entered.
func (s *BaseTrinoParserListener) EnterSpecifiedPrincipal(ctx *SpecifiedPrincipalContext) {}

// ExitSpecifiedPrincipal is called when production specifiedPrincipal is exited.
func (s *BaseTrinoParserListener) ExitSpecifiedPrincipal(ctx *SpecifiedPrincipalContext) {}

// EnterCurrentUserGrantor is called when production currentUserGrantor is entered.
func (s *BaseTrinoParserListener) EnterCurrentUserGrantor(ctx *CurrentUserGrantorContext) {}

// ExitCurrentUserGrantor is called when production currentUserGrantor is exited.
func (s *BaseTrinoParserListener) ExitCurrentUserGrantor(ctx *CurrentUserGrantorContext) {}

// EnterCurrentRoleGrantor is called when production currentRoleGrantor is entered.
func (s *BaseTrinoParserListener) EnterCurrentRoleGrantor(ctx *CurrentRoleGrantorContext) {}

// ExitCurrentRoleGrantor is called when production currentRoleGrantor is exited.
func (s *BaseTrinoParserListener) ExitCurrentRoleGrantor(ctx *CurrentRoleGrantorContext) {}

// EnterUnspecifiedPrincipal is called when production unspecifiedPrincipal is entered.
func (s *BaseTrinoParserListener) EnterUnspecifiedPrincipal(ctx *UnspecifiedPrincipalContext) {}

// ExitUnspecifiedPrincipal is called when production unspecifiedPrincipal is exited.
func (s *BaseTrinoParserListener) ExitUnspecifiedPrincipal(ctx *UnspecifiedPrincipalContext) {}

// EnterUserPrincipal is called when production userPrincipal is entered.
func (s *BaseTrinoParserListener) EnterUserPrincipal(ctx *UserPrincipalContext) {}

// ExitUserPrincipal is called when production userPrincipal is exited.
func (s *BaseTrinoParserListener) ExitUserPrincipal(ctx *UserPrincipalContext) {}

// EnterRolePrincipal is called when production rolePrincipal is entered.
func (s *BaseTrinoParserListener) EnterRolePrincipal(ctx *RolePrincipalContext) {}

// ExitRolePrincipal is called when production rolePrincipal is exited.
func (s *BaseTrinoParserListener) ExitRolePrincipal(ctx *RolePrincipalContext) {}

// EnterRoles is called when production roles is entered.
func (s *BaseTrinoParserListener) EnterRoles(ctx *RolesContext) {}

// ExitRoles is called when production roles is exited.
func (s *BaseTrinoParserListener) ExitRoles(ctx *RolesContext) {}

// EnterUnquotedIdentifier is called when production unquotedIdentifier is entered.
func (s *BaseTrinoParserListener) EnterUnquotedIdentifier(ctx *UnquotedIdentifierContext) {}

// ExitUnquotedIdentifier is called when production unquotedIdentifier is exited.
func (s *BaseTrinoParserListener) ExitUnquotedIdentifier(ctx *UnquotedIdentifierContext) {}

// EnterQuotedIdentifier is called when production quotedIdentifier is entered.
func (s *BaseTrinoParserListener) EnterQuotedIdentifier(ctx *QuotedIdentifierContext) {}

// ExitQuotedIdentifier is called when production quotedIdentifier is exited.
func (s *BaseTrinoParserListener) ExitQuotedIdentifier(ctx *QuotedIdentifierContext) {}

// EnterBackQuotedIdentifier is called when production backQuotedIdentifier is entered.
func (s *BaseTrinoParserListener) EnterBackQuotedIdentifier(ctx *BackQuotedIdentifierContext) {}

// ExitBackQuotedIdentifier is called when production backQuotedIdentifier is exited.
func (s *BaseTrinoParserListener) ExitBackQuotedIdentifier(ctx *BackQuotedIdentifierContext) {}

// EnterDigitIdentifier is called when production digitIdentifier is entered.
func (s *BaseTrinoParserListener) EnterDigitIdentifier(ctx *DigitIdentifierContext) {}

// ExitDigitIdentifier is called when production digitIdentifier is exited.
func (s *BaseTrinoParserListener) ExitDigitIdentifier(ctx *DigitIdentifierContext) {}

// EnterDecimalLiteral is called when production decimalLiteral is entered.
func (s *BaseTrinoParserListener) EnterDecimalLiteral(ctx *DecimalLiteralContext) {}

// ExitDecimalLiteral is called when production decimalLiteral is exited.
func (s *BaseTrinoParserListener) ExitDecimalLiteral(ctx *DecimalLiteralContext) {}

// EnterDoubleLiteral is called when production doubleLiteral is entered.
func (s *BaseTrinoParserListener) EnterDoubleLiteral(ctx *DoubleLiteralContext) {}

// ExitDoubleLiteral is called when production doubleLiteral is exited.
func (s *BaseTrinoParserListener) ExitDoubleLiteral(ctx *DoubleLiteralContext) {}

// EnterIntegerLiteral is called when production integerLiteral is entered.
func (s *BaseTrinoParserListener) EnterIntegerLiteral(ctx *IntegerLiteralContext) {}

// ExitIntegerLiteral is called when production integerLiteral is exited.
func (s *BaseTrinoParserListener) ExitIntegerLiteral(ctx *IntegerLiteralContext) {}

// EnterIdentifierUser is called when production identifierUser is entered.
func (s *BaseTrinoParserListener) EnterIdentifierUser(ctx *IdentifierUserContext) {}

// ExitIdentifierUser is called when production identifierUser is exited.
func (s *BaseTrinoParserListener) ExitIdentifierUser(ctx *IdentifierUserContext) {}

// EnterStringUser is called when production stringUser is entered.
func (s *BaseTrinoParserListener) EnterStringUser(ctx *StringUserContext) {}

// ExitStringUser is called when production stringUser is exited.
func (s *BaseTrinoParserListener) ExitStringUser(ctx *StringUserContext) {}

// EnterNonReserved is called when production nonReserved is entered.
func (s *BaseTrinoParserListener) EnterNonReserved(ctx *NonReservedContext) {}

// ExitNonReserved is called when production nonReserved is exited.
func (s *BaseTrinoParserListener) ExitNonReserved(ctx *NonReservedContext) {}
