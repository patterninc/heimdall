// Code generated from TrinoParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar // TrinoParser
import "github.com/antlr4-go/antlr/v4"

type BaseTrinoParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseTrinoParserVisitor) VisitParse(ctx *ParseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStatements(ctx *StatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSingleStatement(ctx *SingleStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStandaloneExpression(ctx *StandaloneExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStandalonePathSpecification(ctx *StandalonePathSpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStandaloneType(ctx *StandaloneTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStandaloneRowPattern(ctx *StandaloneRowPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStandaloneFunctionSpecification(ctx *StandaloneFunctionSpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStatementDefault(ctx *StatementDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUse(ctx *UseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCreateCatalog(ctx *CreateCatalogContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDropCatalog(ctx *DropCatalogContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCreateSchema(ctx *CreateSchemaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDropSchema(ctx *DropSchemaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRenameSchema(ctx *RenameSchemaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetSchemaAuthorization(ctx *SetSchemaAuthorizationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCreateTableAsSelect(ctx *CreateTableAsSelectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCreateTable(ctx *CreateTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDropTable(ctx *DropTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitInsertInto(ctx *InsertIntoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDelete(ctx *DeleteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTruncateTable(ctx *TruncateTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCommentTable(ctx *CommentTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCommentView(ctx *CommentViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCommentColumn(ctx *CommentColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRenameTable(ctx *RenameTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitAddColumn(ctx *AddColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRenameColumn(ctx *RenameColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDropColumn(ctx *DropColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetColumnType(ctx *SetColumnTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetTableAuthorization(ctx *SetTableAuthorizationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetTableProperties(ctx *SetTablePropertiesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableExecute(ctx *TableExecuteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitAnalyze(ctx *AnalyzeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCreateMaterializedView(ctx *CreateMaterializedViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCreateView(ctx *CreateViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRefreshMaterializedView(ctx *RefreshMaterializedViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDropMaterializedView(ctx *DropMaterializedViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRenameMaterializedView(ctx *RenameMaterializedViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetMaterializedViewProperties(ctx *SetMaterializedViewPropertiesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDropView(ctx *DropViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRenameView(ctx *RenameViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetViewAuthorization(ctx *SetViewAuthorizationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCall(ctx *CallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCreateFunction(ctx *CreateFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDropFunction(ctx *DropFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCreateRole(ctx *CreateRoleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDropRole(ctx *DropRoleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitGrantRoles(ctx *GrantRolesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRevokeRoles(ctx *RevokeRolesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetRole(ctx *SetRoleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitGrant(ctx *GrantContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDeny(ctx *DenyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRevoke(ctx *RevokeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowGrants(ctx *ShowGrantsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExplain(ctx *ExplainContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExplainAnalyze(ctx *ExplainAnalyzeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowCreateTable(ctx *ShowCreateTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowCreateSchema(ctx *ShowCreateSchemaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowCreateView(ctx *ShowCreateViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowCreateMaterializedView(ctx *ShowCreateMaterializedViewContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowTables(ctx *ShowTablesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowSchemas(ctx *ShowSchemasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowCatalogs(ctx *ShowCatalogsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowColumns(ctx *ShowColumnsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowStats(ctx *ShowStatsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowStatsForQuery(ctx *ShowStatsForQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowRoles(ctx *ShowRolesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowRoleGrants(ctx *ShowRoleGrantsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowFunctions(ctx *ShowFunctionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitShowSession(ctx *ShowSessionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetSessionAuthorization(ctx *SetSessionAuthorizationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitResetSessionAuthorization(ctx *ResetSessionAuthorizationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetSession(ctx *SetSessionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitResetSession(ctx *ResetSessionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStartTransaction(ctx *StartTransactionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCommit(ctx *CommitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRollback(ctx *RollbackContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPrepare(ctx *PrepareContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDeallocate(ctx *DeallocateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExecute(ctx *ExecuteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExecuteImmediate(ctx *ExecuteImmediateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDescribeInput(ctx *DescribeInputContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDescribeOutput(ctx *DescribeOutputContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetPath(ctx *SetPathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetTimeZone(ctx *SetTimeZoneContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUpdate(ctx *UpdateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitMerge(ctx *MergeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRootQuery(ctx *RootQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitWithFunction(ctx *WithFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitWith(ctx *WithContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableElement(ctx *TableElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitColumnDefinition(ctx *ColumnDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLikeClause(ctx *LikeClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitProperties(ctx *PropertiesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPropertyAssignments(ctx *PropertyAssignmentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitProperty(ctx *PropertyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDefaultPropertyValue(ctx *DefaultPropertyValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNonDefaultPropertyValue(ctx *NonDefaultPropertyValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQueryNoWith(ctx *QueryNoWithContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLimitRowCount(ctx *LimitRowCountContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRowCount(ctx *RowCountContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQueryTermDefault(ctx *QueryTermDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetOperation(ctx *SetOperationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQueryPrimaryDefault(ctx *QueryPrimaryDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTable(ctx *TableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitInlineTable(ctx *InlineTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSubquery(ctx *SubqueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSortItem(ctx *SortItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQuerySpecification(ctx *QuerySpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitGroupBy(ctx *GroupByContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSingleGroupingSet(ctx *SingleGroupingSetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRollup(ctx *RollupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCube(ctx *CubeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitMultipleGroupingSets(ctx *MultipleGroupingSetsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitGroupingSet(ctx *GroupingSetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitWindowDefinition(ctx *WindowDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitWindowSpecification(ctx *WindowSpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNamedQuery(ctx *NamedQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSetQuantifier(ctx *SetQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSelectSingle(ctx *SelectSingleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSelectAll(ctx *SelectAllContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRelationDefault(ctx *RelationDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJoinRelation(ctx *JoinRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJoinType(ctx *JoinTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJoinCriteria(ctx *JoinCriteriaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSampledRelation(ctx *SampledRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSampleType(ctx *SampleTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTrimsSpecification(ctx *TrimsSpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitListAggOverflowBehavior(ctx *ListAggOverflowBehaviorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitListaggCountIndication(ctx *ListaggCountIndicationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPatternRecognition(ctx *PatternRecognitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitMeasureDefinition(ctx *MeasureDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRowsPerMatch(ctx *RowsPerMatchContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitEmptyMatchHandling(ctx *EmptyMatchHandlingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSkipTo(ctx *SkipToContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSubsetDefinition(ctx *SubsetDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitVariableDefinition(ctx *VariableDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitAliasedRelation(ctx *AliasedRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitColumnAliases(ctx *ColumnAliasesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableName(ctx *TableNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSubqueryRelation(ctx *SubqueryRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUnnest(ctx *UnnestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLateral(ctx *LateralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableFunctionInvocation(ctx *TableFunctionInvocationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitParenthesizedRelation(ctx *ParenthesizedRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableFunctionCall(ctx *TableFunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableFunctionArgument(ctx *TableFunctionArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableArgument(ctx *TableArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableArgumentTable(ctx *TableArgumentTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTableArgumentQuery(ctx *TableArgumentQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDescriptorArgument(ctx *DescriptorArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDescriptorField(ctx *DescriptorFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCopartitionTables(ctx *CopartitionTablesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLogicalNot(ctx *LogicalNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPredicated(ctx *PredicatedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitOr(ctx *OrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitAnd(ctx *AndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitComparison(ctx *ComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQuantifiedComparison(ctx *QuantifiedComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitBetween(ctx *BetweenContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitInList(ctx *InListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitInSubquery(ctx *InSubqueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLike(ctx *LikeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNullPredicate(ctx *NullPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDistinctFrom(ctx *DistinctFromContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitValueExpressionDefault(ctx *ValueExpressionDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitConcatenation(ctx *ConcatenationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitArithmeticBinary(ctx *ArithmeticBinaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitArithmeticUnary(ctx *ArithmeticUnaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitAtTimeZone(ctx *AtTimeZoneContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDereference(ctx *DereferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTypeConstructor(ctx *TypeConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonValue(ctx *JsonValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSpecialDateTimeFunction(ctx *SpecialDateTimeFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSubstring(ctx *SubstringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCast(ctx *CastContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLambda(ctx *LambdaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitParenthesizedExpression(ctx *ParenthesizedExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTrim(ctx *TrimContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitParameter(ctx *ParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNormalize(ctx *NormalizeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonObject(ctx *JsonObjectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitIntervalLiteral(ctx *IntervalLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNumericLiteral(ctx *NumericLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonArray(ctx *JsonArrayContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSimpleCase(ctx *SimpleCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitColumnReference(ctx *ColumnReferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNullLiteral(ctx *NullLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRowConstructor(ctx *RowConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSubscript(ctx *SubscriptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonExists(ctx *JsonExistsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCurrentPath(ctx *CurrentPathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSubqueryExpression(ctx *SubqueryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitBinaryLiteral(ctx *BinaryLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCurrentUser(ctx *CurrentUserContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonQuery(ctx *JsonQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitMeasure(ctx *MeasureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExtract(ctx *ExtractContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitArrayConstructor(ctx *ArrayConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCurrentSchema(ctx *CurrentSchemaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExists(ctx *ExistsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPosition(ctx *PositionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitListagg(ctx *ListaggContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSearchedCase(ctx *SearchedCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCurrentCatalog(ctx *CurrentCatalogContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitGroupingOperation(ctx *GroupingOperationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonPathInvocation(ctx *JsonPathInvocationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonValueExpression(ctx *JsonValueExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonRepresentation(ctx *JsonRepresentationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonArgument(ctx *JsonArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonExistsErrorBehavior(ctx *JsonExistsErrorBehaviorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonValueBehavior(ctx *JsonValueBehaviorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonQueryWrapperBehavior(ctx *JsonQueryWrapperBehaviorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonQueryBehavior(ctx *JsonQueryBehaviorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitJsonObjectMember(ctx *JsonObjectMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitProcessingMode(ctx *ProcessingModeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNullTreatment(ctx *NullTreatmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitBasicStringLiteral(ctx *BasicStringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUnicodeStringLiteral(ctx *UnicodeStringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTimeZoneInterval(ctx *TimeZoneIntervalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTimeZoneString(ctx *TimeZoneStringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitComparisonQuantifier(ctx *ComparisonQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitBooleanValue(ctx *BooleanValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitInterval(ctx *IntervalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitIntervalField(ctx *IntervalFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNormalForm(ctx *NormalFormContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRowType(ctx *RowTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitIntervalType(ctx *IntervalTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDoublePrecisionType(ctx *DoublePrecisionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLegacyArrayType(ctx *LegacyArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitGenericType(ctx *GenericTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDateTimeType(ctx *DateTimeTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLegacyMapType(ctx *LegacyMapTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRowField(ctx *RowFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTypeParameter(ctx *TypeParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitWhenClause(ctx *WhenClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitFilter(ctx *FilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitMergeUpdate(ctx *MergeUpdateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitMergeDelete(ctx *MergeDeleteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitMergeInsert(ctx *MergeInsertContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitOver(ctx *OverContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitWindowFrame(ctx *WindowFrameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitFrameExtent(ctx *FrameExtentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUnboundedFrame(ctx *UnboundedFrameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCurrentRowBound(ctx *CurrentRowBoundContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitBoundedFrame(ctx *BoundedFrameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQuantifiedPrimary(ctx *QuantifiedPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPatternConcatenation(ctx *PatternConcatenationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPatternAlternation(ctx *PatternAlternationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPatternVariable(ctx *PatternVariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitEmptyPattern(ctx *EmptyPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPatternPermutation(ctx *PatternPermutationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitGroupedPattern(ctx *GroupedPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPartitionStartAnchor(ctx *PartitionStartAnchorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPartitionEndAnchor(ctx *PartitionEndAnchorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExcludedPattern(ctx *ExcludedPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitZeroOrMoreQuantifier(ctx *ZeroOrMoreQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitOneOrMoreQuantifier(ctx *OneOrMoreQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitZeroOrOneQuantifier(ctx *ZeroOrOneQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRangeQuantifier(ctx *RangeQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUpdateAssignment(ctx *UpdateAssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExplainFormat(ctx *ExplainFormatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitExplainType(ctx *ExplainTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitIsolationLevel(ctx *IsolationLevelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitTransactionAccessMode(ctx *TransactionAccessModeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitReadUncommitted(ctx *ReadUncommittedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitReadCommitted(ctx *ReadCommittedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRepeatableRead(ctx *RepeatableReadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSerializable(ctx *SerializableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPositionalArgument(ctx *PositionalArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNamedArgument(ctx *NamedArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQualifiedArgument(ctx *QualifiedArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUnqualifiedArgument(ctx *UnqualifiedArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPathSpecification(ctx *PathSpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitFunctionSpecification(ctx *FunctionSpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitFunctionDeclaration(ctx *FunctionDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitParameterDeclaration(ctx *ParameterDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitReturnsClause(ctx *ReturnsClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLanguageCharacteristic(ctx *LanguageCharacteristicContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDeterministicCharacteristic(ctx *DeterministicCharacteristicContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitReturnsNullOnNullInputCharacteristic(ctx *ReturnsNullOnNullInputCharacteristicContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCalledOnNullInputCharacteristic(ctx *CalledOnNullInputCharacteristicContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSecurityCharacteristic(ctx *SecurityCharacteristicContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCommentCharacteristic(ctx *CommentCharacteristicContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitReturnStatement(ctx *ReturnStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitAssignmentStatement(ctx *AssignmentStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSimpleCaseStatement(ctx *SimpleCaseStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSearchedCaseStatement(ctx *SearchedCaseStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitIfStatement(ctx *IfStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitIterateStatement(ctx *IterateStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLeaveStatement(ctx *LeaveStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCompoundStatement(ctx *CompoundStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitLoopStatement(ctx *LoopStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitWhileStatement(ctx *WhileStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRepeatStatement(ctx *RepeatStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCaseStatementWhenClause(ctx *CaseStatementWhenClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitElseIfClause(ctx *ElseIfClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitElseClause(ctx *ElseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSqlStatementList(ctx *SqlStatementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitPrivilege(ctx *PrivilegeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQualifiedName(ctx *QualifiedNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQueryPeriod(ctx *QueryPeriodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRangeType(ctx *RangeTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitSpecifiedPrincipal(ctx *SpecifiedPrincipalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCurrentUserGrantor(ctx *CurrentUserGrantorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitCurrentRoleGrantor(ctx *CurrentRoleGrantorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUnspecifiedPrincipal(ctx *UnspecifiedPrincipalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUserPrincipal(ctx *UserPrincipalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRolePrincipal(ctx *RolePrincipalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitRoles(ctx *RolesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitUnquotedIdentifier(ctx *UnquotedIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitQuotedIdentifier(ctx *QuotedIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitBackQuotedIdentifier(ctx *BackQuotedIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDigitIdentifier(ctx *DigitIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitDoubleLiteral(ctx *DoubleLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitIdentifierUser(ctx *IdentifierUserContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitStringUser(ctx *StringUserContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTrinoParserVisitor) VisitNonReserved(ctx *NonReservedContext) interface{} {
	return v.VisitChildren(ctx)
}
