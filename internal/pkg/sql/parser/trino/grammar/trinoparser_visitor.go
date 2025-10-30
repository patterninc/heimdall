// Code generated from TrinoParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar // TrinoParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by TrinoParser.
type TrinoParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by TrinoParser#parse.
	VisitParse(ctx *ParseContext) interface{}

	// Visit a parse tree produced by TrinoParser#statements.
	VisitStatements(ctx *StatementsContext) interface{}

	// Visit a parse tree produced by TrinoParser#singleStatement.
	VisitSingleStatement(ctx *SingleStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#standaloneExpression.
	VisitStandaloneExpression(ctx *StandaloneExpressionContext) interface{}

	// Visit a parse tree produced by TrinoParser#standalonePathSpecification.
	VisitStandalonePathSpecification(ctx *StandalonePathSpecificationContext) interface{}

	// Visit a parse tree produced by TrinoParser#standaloneType.
	VisitStandaloneType(ctx *StandaloneTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#standaloneRowPattern.
	VisitStandaloneRowPattern(ctx *StandaloneRowPatternContext) interface{}

	// Visit a parse tree produced by TrinoParser#standaloneFunctionSpecification.
	VisitStandaloneFunctionSpecification(ctx *StandaloneFunctionSpecificationContext) interface{}

	// Visit a parse tree produced by TrinoParser#statementDefault.
	VisitStatementDefault(ctx *StatementDefaultContext) interface{}

	// Visit a parse tree produced by TrinoParser#use.
	VisitUse(ctx *UseContext) interface{}

	// Visit a parse tree produced by TrinoParser#createCatalog.
	VisitCreateCatalog(ctx *CreateCatalogContext) interface{}

	// Visit a parse tree produced by TrinoParser#dropCatalog.
	VisitDropCatalog(ctx *DropCatalogContext) interface{}

	// Visit a parse tree produced by TrinoParser#createSchema.
	VisitCreateSchema(ctx *CreateSchemaContext) interface{}

	// Visit a parse tree produced by TrinoParser#dropSchema.
	VisitDropSchema(ctx *DropSchemaContext) interface{}

	// Visit a parse tree produced by TrinoParser#renameSchema.
	VisitRenameSchema(ctx *RenameSchemaContext) interface{}

	// Visit a parse tree produced by TrinoParser#setSchemaAuthorization.
	VisitSetSchemaAuthorization(ctx *SetSchemaAuthorizationContext) interface{}

	// Visit a parse tree produced by TrinoParser#createTableAsSelect.
	VisitCreateTableAsSelect(ctx *CreateTableAsSelectContext) interface{}

	// Visit a parse tree produced by TrinoParser#createTable.
	VisitCreateTable(ctx *CreateTableContext) interface{}

	// Visit a parse tree produced by TrinoParser#dropTable.
	VisitDropTable(ctx *DropTableContext) interface{}

	// Visit a parse tree produced by TrinoParser#insertInto.
	VisitInsertInto(ctx *InsertIntoContext) interface{}

	// Visit a parse tree produced by TrinoParser#delete.
	VisitDelete(ctx *DeleteContext) interface{}

	// Visit a parse tree produced by TrinoParser#truncateTable.
	VisitTruncateTable(ctx *TruncateTableContext) interface{}

	// Visit a parse tree produced by TrinoParser#commentTable.
	VisitCommentTable(ctx *CommentTableContext) interface{}

	// Visit a parse tree produced by TrinoParser#commentView.
	VisitCommentView(ctx *CommentViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#commentColumn.
	VisitCommentColumn(ctx *CommentColumnContext) interface{}

	// Visit a parse tree produced by TrinoParser#renameTable.
	VisitRenameTable(ctx *RenameTableContext) interface{}

	// Visit a parse tree produced by TrinoParser#addColumn.
	VisitAddColumn(ctx *AddColumnContext) interface{}

	// Visit a parse tree produced by TrinoParser#renameColumn.
	VisitRenameColumn(ctx *RenameColumnContext) interface{}

	// Visit a parse tree produced by TrinoParser#dropColumn.
	VisitDropColumn(ctx *DropColumnContext) interface{}

	// Visit a parse tree produced by TrinoParser#setColumnType.
	VisitSetColumnType(ctx *SetColumnTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#setTableAuthorization.
	VisitSetTableAuthorization(ctx *SetTableAuthorizationContext) interface{}

	// Visit a parse tree produced by TrinoParser#setTableProperties.
	VisitSetTableProperties(ctx *SetTablePropertiesContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableExecute.
	VisitTableExecute(ctx *TableExecuteContext) interface{}

	// Visit a parse tree produced by TrinoParser#analyze.
	VisitAnalyze(ctx *AnalyzeContext) interface{}

	// Visit a parse tree produced by TrinoParser#createMaterializedView.
	VisitCreateMaterializedView(ctx *CreateMaterializedViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#createView.
	VisitCreateView(ctx *CreateViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#refreshMaterializedView.
	VisitRefreshMaterializedView(ctx *RefreshMaterializedViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#dropMaterializedView.
	VisitDropMaterializedView(ctx *DropMaterializedViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#renameMaterializedView.
	VisitRenameMaterializedView(ctx *RenameMaterializedViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#setMaterializedViewProperties.
	VisitSetMaterializedViewProperties(ctx *SetMaterializedViewPropertiesContext) interface{}

	// Visit a parse tree produced by TrinoParser#dropView.
	VisitDropView(ctx *DropViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#renameView.
	VisitRenameView(ctx *RenameViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#setViewAuthorization.
	VisitSetViewAuthorization(ctx *SetViewAuthorizationContext) interface{}

	// Visit a parse tree produced by TrinoParser#call.
	VisitCall(ctx *CallContext) interface{}

	// Visit a parse tree produced by TrinoParser#createFunction.
	VisitCreateFunction(ctx *CreateFunctionContext) interface{}

	// Visit a parse tree produced by TrinoParser#dropFunction.
	VisitDropFunction(ctx *DropFunctionContext) interface{}

	// Visit a parse tree produced by TrinoParser#createRole.
	VisitCreateRole(ctx *CreateRoleContext) interface{}

	// Visit a parse tree produced by TrinoParser#dropRole.
	VisitDropRole(ctx *DropRoleContext) interface{}

	// Visit a parse tree produced by TrinoParser#grantRoles.
	VisitGrantRoles(ctx *GrantRolesContext) interface{}

	// Visit a parse tree produced by TrinoParser#revokeRoles.
	VisitRevokeRoles(ctx *RevokeRolesContext) interface{}

	// Visit a parse tree produced by TrinoParser#setRole.
	VisitSetRole(ctx *SetRoleContext) interface{}

	// Visit a parse tree produced by TrinoParser#grant.
	VisitGrant(ctx *GrantContext) interface{}

	// Visit a parse tree produced by TrinoParser#deny.
	VisitDeny(ctx *DenyContext) interface{}

	// Visit a parse tree produced by TrinoParser#revoke.
	VisitRevoke(ctx *RevokeContext) interface{}

	// Visit a parse tree produced by TrinoParser#showGrants.
	VisitShowGrants(ctx *ShowGrantsContext) interface{}

	// Visit a parse tree produced by TrinoParser#explain.
	VisitExplain(ctx *ExplainContext) interface{}

	// Visit a parse tree produced by TrinoParser#explainAnalyze.
	VisitExplainAnalyze(ctx *ExplainAnalyzeContext) interface{}

	// Visit a parse tree produced by TrinoParser#showCreateTable.
	VisitShowCreateTable(ctx *ShowCreateTableContext) interface{}

	// Visit a parse tree produced by TrinoParser#showCreateSchema.
	VisitShowCreateSchema(ctx *ShowCreateSchemaContext) interface{}

	// Visit a parse tree produced by TrinoParser#showCreateView.
	VisitShowCreateView(ctx *ShowCreateViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#showCreateMaterializedView.
	VisitShowCreateMaterializedView(ctx *ShowCreateMaterializedViewContext) interface{}

	// Visit a parse tree produced by TrinoParser#showTables.
	VisitShowTables(ctx *ShowTablesContext) interface{}

	// Visit a parse tree produced by TrinoParser#showSchemas.
	VisitShowSchemas(ctx *ShowSchemasContext) interface{}

	// Visit a parse tree produced by TrinoParser#showCatalogs.
	VisitShowCatalogs(ctx *ShowCatalogsContext) interface{}

	// Visit a parse tree produced by TrinoParser#showColumns.
	VisitShowColumns(ctx *ShowColumnsContext) interface{}

	// Visit a parse tree produced by TrinoParser#showStats.
	VisitShowStats(ctx *ShowStatsContext) interface{}

	// Visit a parse tree produced by TrinoParser#showStatsForQuery.
	VisitShowStatsForQuery(ctx *ShowStatsForQueryContext) interface{}

	// Visit a parse tree produced by TrinoParser#showRoles.
	VisitShowRoles(ctx *ShowRolesContext) interface{}

	// Visit a parse tree produced by TrinoParser#showRoleGrants.
	VisitShowRoleGrants(ctx *ShowRoleGrantsContext) interface{}

	// Visit a parse tree produced by TrinoParser#showFunctions.
	VisitShowFunctions(ctx *ShowFunctionsContext) interface{}

	// Visit a parse tree produced by TrinoParser#showSession.
	VisitShowSession(ctx *ShowSessionContext) interface{}

	// Visit a parse tree produced by TrinoParser#setSessionAuthorization.
	VisitSetSessionAuthorization(ctx *SetSessionAuthorizationContext) interface{}

	// Visit a parse tree produced by TrinoParser#resetSessionAuthorization.
	VisitResetSessionAuthorization(ctx *ResetSessionAuthorizationContext) interface{}

	// Visit a parse tree produced by TrinoParser#setSession.
	VisitSetSession(ctx *SetSessionContext) interface{}

	// Visit a parse tree produced by TrinoParser#resetSession.
	VisitResetSession(ctx *ResetSessionContext) interface{}

	// Visit a parse tree produced by TrinoParser#startTransaction.
	VisitStartTransaction(ctx *StartTransactionContext) interface{}

	// Visit a parse tree produced by TrinoParser#commit.
	VisitCommit(ctx *CommitContext) interface{}

	// Visit a parse tree produced by TrinoParser#rollback.
	VisitRollback(ctx *RollbackContext) interface{}

	// Visit a parse tree produced by TrinoParser#prepare.
	VisitPrepare(ctx *PrepareContext) interface{}

	// Visit a parse tree produced by TrinoParser#deallocate.
	VisitDeallocate(ctx *DeallocateContext) interface{}

	// Visit a parse tree produced by TrinoParser#execute.
	VisitExecute(ctx *ExecuteContext) interface{}

	// Visit a parse tree produced by TrinoParser#executeImmediate.
	VisitExecuteImmediate(ctx *ExecuteImmediateContext) interface{}

	// Visit a parse tree produced by TrinoParser#describeInput.
	VisitDescribeInput(ctx *DescribeInputContext) interface{}

	// Visit a parse tree produced by TrinoParser#describeOutput.
	VisitDescribeOutput(ctx *DescribeOutputContext) interface{}

	// Visit a parse tree produced by TrinoParser#setPath.
	VisitSetPath(ctx *SetPathContext) interface{}

	// Visit a parse tree produced by TrinoParser#setTimeZone.
	VisitSetTimeZone(ctx *SetTimeZoneContext) interface{}

	// Visit a parse tree produced by TrinoParser#update.
	VisitUpdate(ctx *UpdateContext) interface{}

	// Visit a parse tree produced by TrinoParser#merge.
	VisitMerge(ctx *MergeContext) interface{}

	// Visit a parse tree produced by TrinoParser#rootQuery.
	VisitRootQuery(ctx *RootQueryContext) interface{}

	// Visit a parse tree produced by TrinoParser#withFunction.
	VisitWithFunction(ctx *WithFunctionContext) interface{}

	// Visit a parse tree produced by TrinoParser#query.
	VisitQuery(ctx *QueryContext) interface{}

	// Visit a parse tree produced by TrinoParser#with.
	VisitWith(ctx *WithContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableElement.
	VisitTableElement(ctx *TableElementContext) interface{}

	// Visit a parse tree produced by TrinoParser#columnDefinition.
	VisitColumnDefinition(ctx *ColumnDefinitionContext) interface{}

	// Visit a parse tree produced by TrinoParser#likeClause.
	VisitLikeClause(ctx *LikeClauseContext) interface{}

	// Visit a parse tree produced by TrinoParser#properties.
	VisitProperties(ctx *PropertiesContext) interface{}

	// Visit a parse tree produced by TrinoParser#propertyAssignments.
	VisitPropertyAssignments(ctx *PropertyAssignmentsContext) interface{}

	// Visit a parse tree produced by TrinoParser#property.
	VisitProperty(ctx *PropertyContext) interface{}

	// Visit a parse tree produced by TrinoParser#defaultPropertyValue.
	VisitDefaultPropertyValue(ctx *DefaultPropertyValueContext) interface{}

	// Visit a parse tree produced by TrinoParser#nonDefaultPropertyValue.
	VisitNonDefaultPropertyValue(ctx *NonDefaultPropertyValueContext) interface{}

	// Visit a parse tree produced by TrinoParser#queryNoWith.
	VisitQueryNoWith(ctx *QueryNoWithContext) interface{}

	// Visit a parse tree produced by TrinoParser#limitRowCount.
	VisitLimitRowCount(ctx *LimitRowCountContext) interface{}

	// Visit a parse tree produced by TrinoParser#rowCount.
	VisitRowCount(ctx *RowCountContext) interface{}

	// Visit a parse tree produced by TrinoParser#queryTermDefault.
	VisitQueryTermDefault(ctx *QueryTermDefaultContext) interface{}

	// Visit a parse tree produced by TrinoParser#setOperation.
	VisitSetOperation(ctx *SetOperationContext) interface{}

	// Visit a parse tree produced by TrinoParser#queryPrimaryDefault.
	VisitQueryPrimaryDefault(ctx *QueryPrimaryDefaultContext) interface{}

	// Visit a parse tree produced by TrinoParser#table.
	VisitTable(ctx *TableContext) interface{}

	// Visit a parse tree produced by TrinoParser#inlineTable.
	VisitInlineTable(ctx *InlineTableContext) interface{}

	// Visit a parse tree produced by TrinoParser#subquery.
	VisitSubquery(ctx *SubqueryContext) interface{}

	// Visit a parse tree produced by TrinoParser#sortItem.
	VisitSortItem(ctx *SortItemContext) interface{}

	// Visit a parse tree produced by TrinoParser#querySpecification.
	VisitQuerySpecification(ctx *QuerySpecificationContext) interface{}

	// Visit a parse tree produced by TrinoParser#groupBy.
	VisitGroupBy(ctx *GroupByContext) interface{}

	// Visit a parse tree produced by TrinoParser#singleGroupingSet.
	VisitSingleGroupingSet(ctx *SingleGroupingSetContext) interface{}

	// Visit a parse tree produced by TrinoParser#rollup.
	VisitRollup(ctx *RollupContext) interface{}

	// Visit a parse tree produced by TrinoParser#cube.
	VisitCube(ctx *CubeContext) interface{}

	// Visit a parse tree produced by TrinoParser#multipleGroupingSets.
	VisitMultipleGroupingSets(ctx *MultipleGroupingSetsContext) interface{}

	// Visit a parse tree produced by TrinoParser#groupingSet.
	VisitGroupingSet(ctx *GroupingSetContext) interface{}

	// Visit a parse tree produced by TrinoParser#windowDefinition.
	VisitWindowDefinition(ctx *WindowDefinitionContext) interface{}

	// Visit a parse tree produced by TrinoParser#windowSpecification.
	VisitWindowSpecification(ctx *WindowSpecificationContext) interface{}

	// Visit a parse tree produced by TrinoParser#namedQuery.
	VisitNamedQuery(ctx *NamedQueryContext) interface{}

	// Visit a parse tree produced by TrinoParser#setQuantifier.
	VisitSetQuantifier(ctx *SetQuantifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#selectSingle.
	VisitSelectSingle(ctx *SelectSingleContext) interface{}

	// Visit a parse tree produced by TrinoParser#selectAll.
	VisitSelectAll(ctx *SelectAllContext) interface{}

	// Visit a parse tree produced by TrinoParser#relationDefault.
	VisitRelationDefault(ctx *RelationDefaultContext) interface{}

	// Visit a parse tree produced by TrinoParser#joinRelation.
	VisitJoinRelation(ctx *JoinRelationContext) interface{}

	// Visit a parse tree produced by TrinoParser#joinType.
	VisitJoinType(ctx *JoinTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#joinCriteria.
	VisitJoinCriteria(ctx *JoinCriteriaContext) interface{}

	// Visit a parse tree produced by TrinoParser#sampledRelation.
	VisitSampledRelation(ctx *SampledRelationContext) interface{}

	// Visit a parse tree produced by TrinoParser#sampleType.
	VisitSampleType(ctx *SampleTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#trimsSpecification.
	VisitTrimsSpecification(ctx *TrimsSpecificationContext) interface{}

	// Visit a parse tree produced by TrinoParser#listAggOverflowBehavior.
	VisitListAggOverflowBehavior(ctx *ListAggOverflowBehaviorContext) interface{}

	// Visit a parse tree produced by TrinoParser#listaggCountIndication.
	VisitListaggCountIndication(ctx *ListaggCountIndicationContext) interface{}

	// Visit a parse tree produced by TrinoParser#patternRecognition.
	VisitPatternRecognition(ctx *PatternRecognitionContext) interface{}

	// Visit a parse tree produced by TrinoParser#measureDefinition.
	VisitMeasureDefinition(ctx *MeasureDefinitionContext) interface{}

	// Visit a parse tree produced by TrinoParser#rowsPerMatch.
	VisitRowsPerMatch(ctx *RowsPerMatchContext) interface{}

	// Visit a parse tree produced by TrinoParser#emptyMatchHandling.
	VisitEmptyMatchHandling(ctx *EmptyMatchHandlingContext) interface{}

	// Visit a parse tree produced by TrinoParser#skipTo.
	VisitSkipTo(ctx *SkipToContext) interface{}

	// Visit a parse tree produced by TrinoParser#subsetDefinition.
	VisitSubsetDefinition(ctx *SubsetDefinitionContext) interface{}

	// Visit a parse tree produced by TrinoParser#variableDefinition.
	VisitVariableDefinition(ctx *VariableDefinitionContext) interface{}

	// Visit a parse tree produced by TrinoParser#aliasedRelation.
	VisitAliasedRelation(ctx *AliasedRelationContext) interface{}

	// Visit a parse tree produced by TrinoParser#columnAliases.
	VisitColumnAliases(ctx *ColumnAliasesContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableName.
	VisitTableName(ctx *TableNameContext) interface{}

	// Visit a parse tree produced by TrinoParser#subqueryRelation.
	VisitSubqueryRelation(ctx *SubqueryRelationContext) interface{}

	// Visit a parse tree produced by TrinoParser#unnest.
	VisitUnnest(ctx *UnnestContext) interface{}

	// Visit a parse tree produced by TrinoParser#lateral.
	VisitLateral(ctx *LateralContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableFunctionInvocation.
	VisitTableFunctionInvocation(ctx *TableFunctionInvocationContext) interface{}

	// Visit a parse tree produced by TrinoParser#parenthesizedRelation.
	VisitParenthesizedRelation(ctx *ParenthesizedRelationContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableFunctionCall.
	VisitTableFunctionCall(ctx *TableFunctionCallContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableFunctionArgument.
	VisitTableFunctionArgument(ctx *TableFunctionArgumentContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableArgument.
	VisitTableArgument(ctx *TableArgumentContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableArgumentTable.
	VisitTableArgumentTable(ctx *TableArgumentTableContext) interface{}

	// Visit a parse tree produced by TrinoParser#tableArgumentQuery.
	VisitTableArgumentQuery(ctx *TableArgumentQueryContext) interface{}

	// Visit a parse tree produced by TrinoParser#descriptorArgument.
	VisitDescriptorArgument(ctx *DescriptorArgumentContext) interface{}

	// Visit a parse tree produced by TrinoParser#descriptorField.
	VisitDescriptorField(ctx *DescriptorFieldContext) interface{}

	// Visit a parse tree produced by TrinoParser#copartitionTables.
	VisitCopartitionTables(ctx *CopartitionTablesContext) interface{}

	// Visit a parse tree produced by TrinoParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by TrinoParser#logicalNot.
	VisitLogicalNot(ctx *LogicalNotContext) interface{}

	// Visit a parse tree produced by TrinoParser#predicated.
	VisitPredicated(ctx *PredicatedContext) interface{}

	// Visit a parse tree produced by TrinoParser#or.
	VisitOr(ctx *OrContext) interface{}

	// Visit a parse tree produced by TrinoParser#and.
	VisitAnd(ctx *AndContext) interface{}

	// Visit a parse tree produced by TrinoParser#comparison.
	VisitComparison(ctx *ComparisonContext) interface{}

	// Visit a parse tree produced by TrinoParser#quantifiedComparison.
	VisitQuantifiedComparison(ctx *QuantifiedComparisonContext) interface{}

	// Visit a parse tree produced by TrinoParser#between.
	VisitBetween(ctx *BetweenContext) interface{}

	// Visit a parse tree produced by TrinoParser#inList.
	VisitInList(ctx *InListContext) interface{}

	// Visit a parse tree produced by TrinoParser#inSubquery.
	VisitInSubquery(ctx *InSubqueryContext) interface{}

	// Visit a parse tree produced by TrinoParser#like.
	VisitLike(ctx *LikeContext) interface{}

	// Visit a parse tree produced by TrinoParser#nullPredicate.
	VisitNullPredicate(ctx *NullPredicateContext) interface{}

	// Visit a parse tree produced by TrinoParser#distinctFrom.
	VisitDistinctFrom(ctx *DistinctFromContext) interface{}

	// Visit a parse tree produced by TrinoParser#valueExpressionDefault.
	VisitValueExpressionDefault(ctx *ValueExpressionDefaultContext) interface{}

	// Visit a parse tree produced by TrinoParser#concatenation.
	VisitConcatenation(ctx *ConcatenationContext) interface{}

	// Visit a parse tree produced by TrinoParser#arithmeticBinary.
	VisitArithmeticBinary(ctx *ArithmeticBinaryContext) interface{}

	// Visit a parse tree produced by TrinoParser#arithmeticUnary.
	VisitArithmeticUnary(ctx *ArithmeticUnaryContext) interface{}

	// Visit a parse tree produced by TrinoParser#atTimeZone.
	VisitAtTimeZone(ctx *AtTimeZoneContext) interface{}

	// Visit a parse tree produced by TrinoParser#dereference.
	VisitDereference(ctx *DereferenceContext) interface{}

	// Visit a parse tree produced by TrinoParser#typeConstructor.
	VisitTypeConstructor(ctx *TypeConstructorContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonValue.
	VisitJsonValue(ctx *JsonValueContext) interface{}

	// Visit a parse tree produced by TrinoParser#specialDateTimeFunction.
	VisitSpecialDateTimeFunction(ctx *SpecialDateTimeFunctionContext) interface{}

	// Visit a parse tree produced by TrinoParser#substring.
	VisitSubstring(ctx *SubstringContext) interface{}

	// Visit a parse tree produced by TrinoParser#cast.
	VisitCast(ctx *CastContext) interface{}

	// Visit a parse tree produced by TrinoParser#lambda.
	VisitLambda(ctx *LambdaContext) interface{}

	// Visit a parse tree produced by TrinoParser#parenthesizedExpression.
	VisitParenthesizedExpression(ctx *ParenthesizedExpressionContext) interface{}

	// Visit a parse tree produced by TrinoParser#trim.
	VisitTrim(ctx *TrimContext) interface{}

	// Visit a parse tree produced by TrinoParser#parameter.
	VisitParameter(ctx *ParameterContext) interface{}

	// Visit a parse tree produced by TrinoParser#normalize.
	VisitNormalize(ctx *NormalizeContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonObject.
	VisitJsonObject(ctx *JsonObjectContext) interface{}

	// Visit a parse tree produced by TrinoParser#intervalLiteral.
	VisitIntervalLiteral(ctx *IntervalLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#numericLiteral.
	VisitNumericLiteral(ctx *NumericLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#booleanLiteral.
	VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonArray.
	VisitJsonArray(ctx *JsonArrayContext) interface{}

	// Visit a parse tree produced by TrinoParser#simpleCase.
	VisitSimpleCase(ctx *SimpleCaseContext) interface{}

	// Visit a parse tree produced by TrinoParser#columnReference.
	VisitColumnReference(ctx *ColumnReferenceContext) interface{}

	// Visit a parse tree produced by TrinoParser#nullLiteral.
	VisitNullLiteral(ctx *NullLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#rowConstructor.
	VisitRowConstructor(ctx *RowConstructorContext) interface{}

	// Visit a parse tree produced by TrinoParser#subscript.
	VisitSubscript(ctx *SubscriptContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonExists.
	VisitJsonExists(ctx *JsonExistsContext) interface{}

	// Visit a parse tree produced by TrinoParser#currentPath.
	VisitCurrentPath(ctx *CurrentPathContext) interface{}

	// Visit a parse tree produced by TrinoParser#subqueryExpression.
	VisitSubqueryExpression(ctx *SubqueryExpressionContext) interface{}

	// Visit a parse tree produced by TrinoParser#binaryLiteral.
	VisitBinaryLiteral(ctx *BinaryLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#currentUser.
	VisitCurrentUser(ctx *CurrentUserContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonQuery.
	VisitJsonQuery(ctx *JsonQueryContext) interface{}

	// Visit a parse tree produced by TrinoParser#measure.
	VisitMeasure(ctx *MeasureContext) interface{}

	// Visit a parse tree produced by TrinoParser#extract.
	VisitExtract(ctx *ExtractContext) interface{}

	// Visit a parse tree produced by TrinoParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#arrayConstructor.
	VisitArrayConstructor(ctx *ArrayConstructorContext) interface{}

	// Visit a parse tree produced by TrinoParser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by TrinoParser#currentSchema.
	VisitCurrentSchema(ctx *CurrentSchemaContext) interface{}

	// Visit a parse tree produced by TrinoParser#exists.
	VisitExists(ctx *ExistsContext) interface{}

	// Visit a parse tree produced by TrinoParser#position.
	VisitPosition(ctx *PositionContext) interface{}

	// Visit a parse tree produced by TrinoParser#listagg.
	VisitListagg(ctx *ListaggContext) interface{}

	// Visit a parse tree produced by TrinoParser#searchedCase.
	VisitSearchedCase(ctx *SearchedCaseContext) interface{}

	// Visit a parse tree produced by TrinoParser#currentCatalog.
	VisitCurrentCatalog(ctx *CurrentCatalogContext) interface{}

	// Visit a parse tree produced by TrinoParser#groupingOperation.
	VisitGroupingOperation(ctx *GroupingOperationContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonPathInvocation.
	VisitJsonPathInvocation(ctx *JsonPathInvocationContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonValueExpression.
	VisitJsonValueExpression(ctx *JsonValueExpressionContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonRepresentation.
	VisitJsonRepresentation(ctx *JsonRepresentationContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonArgument.
	VisitJsonArgument(ctx *JsonArgumentContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonExistsErrorBehavior.
	VisitJsonExistsErrorBehavior(ctx *JsonExistsErrorBehaviorContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonValueBehavior.
	VisitJsonValueBehavior(ctx *JsonValueBehaviorContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonQueryWrapperBehavior.
	VisitJsonQueryWrapperBehavior(ctx *JsonQueryWrapperBehaviorContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonQueryBehavior.
	VisitJsonQueryBehavior(ctx *JsonQueryBehaviorContext) interface{}

	// Visit a parse tree produced by TrinoParser#jsonObjectMember.
	VisitJsonObjectMember(ctx *JsonObjectMemberContext) interface{}

	// Visit a parse tree produced by TrinoParser#processingMode.
	VisitProcessingMode(ctx *ProcessingModeContext) interface{}

	// Visit a parse tree produced by TrinoParser#nullTreatment.
	VisitNullTreatment(ctx *NullTreatmentContext) interface{}

	// Visit a parse tree produced by TrinoParser#basicStringLiteral.
	VisitBasicStringLiteral(ctx *BasicStringLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#unicodeStringLiteral.
	VisitUnicodeStringLiteral(ctx *UnicodeStringLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#timeZoneInterval.
	VisitTimeZoneInterval(ctx *TimeZoneIntervalContext) interface{}

	// Visit a parse tree produced by TrinoParser#timeZoneString.
	VisitTimeZoneString(ctx *TimeZoneStringContext) interface{}

	// Visit a parse tree produced by TrinoParser#comparisonOperator.
	VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{}

	// Visit a parse tree produced by TrinoParser#comparisonQuantifier.
	VisitComparisonQuantifier(ctx *ComparisonQuantifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#booleanValue.
	VisitBooleanValue(ctx *BooleanValueContext) interface{}

	// Visit a parse tree produced by TrinoParser#interval.
	VisitInterval(ctx *IntervalContext) interface{}

	// Visit a parse tree produced by TrinoParser#intervalField.
	VisitIntervalField(ctx *IntervalFieldContext) interface{}

	// Visit a parse tree produced by TrinoParser#normalForm.
	VisitNormalForm(ctx *NormalFormContext) interface{}

	// Visit a parse tree produced by TrinoParser#rowType.
	VisitRowType(ctx *RowTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#intervalType.
	VisitIntervalType(ctx *IntervalTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#doublePrecisionType.
	VisitDoublePrecisionType(ctx *DoublePrecisionTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#legacyArrayType.
	VisitLegacyArrayType(ctx *LegacyArrayTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#genericType.
	VisitGenericType(ctx *GenericTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#dateTimeType.
	VisitDateTimeType(ctx *DateTimeTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#legacyMapType.
	VisitLegacyMapType(ctx *LegacyMapTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#rowField.
	VisitRowField(ctx *RowFieldContext) interface{}

	// Visit a parse tree produced by TrinoParser#typeParameter.
	VisitTypeParameter(ctx *TypeParameterContext) interface{}

	// Visit a parse tree produced by TrinoParser#whenClause.
	VisitWhenClause(ctx *WhenClauseContext) interface{}

	// Visit a parse tree produced by TrinoParser#filter.
	VisitFilter(ctx *FilterContext) interface{}

	// Visit a parse tree produced by TrinoParser#mergeUpdate.
	VisitMergeUpdate(ctx *MergeUpdateContext) interface{}

	// Visit a parse tree produced by TrinoParser#mergeDelete.
	VisitMergeDelete(ctx *MergeDeleteContext) interface{}

	// Visit a parse tree produced by TrinoParser#mergeInsert.
	VisitMergeInsert(ctx *MergeInsertContext) interface{}

	// Visit a parse tree produced by TrinoParser#over.
	VisitOver(ctx *OverContext) interface{}

	// Visit a parse tree produced by TrinoParser#windowFrame.
	VisitWindowFrame(ctx *WindowFrameContext) interface{}

	// Visit a parse tree produced by TrinoParser#frameExtent.
	VisitFrameExtent(ctx *FrameExtentContext) interface{}

	// Visit a parse tree produced by TrinoParser#unboundedFrame.
	VisitUnboundedFrame(ctx *UnboundedFrameContext) interface{}

	// Visit a parse tree produced by TrinoParser#currentRowBound.
	VisitCurrentRowBound(ctx *CurrentRowBoundContext) interface{}

	// Visit a parse tree produced by TrinoParser#boundedFrame.
	VisitBoundedFrame(ctx *BoundedFrameContext) interface{}

	// Visit a parse tree produced by TrinoParser#quantifiedPrimary.
	VisitQuantifiedPrimary(ctx *QuantifiedPrimaryContext) interface{}

	// Visit a parse tree produced by TrinoParser#patternConcatenation.
	VisitPatternConcatenation(ctx *PatternConcatenationContext) interface{}

	// Visit a parse tree produced by TrinoParser#patternAlternation.
	VisitPatternAlternation(ctx *PatternAlternationContext) interface{}

	// Visit a parse tree produced by TrinoParser#patternVariable.
	VisitPatternVariable(ctx *PatternVariableContext) interface{}

	// Visit a parse tree produced by TrinoParser#emptyPattern.
	VisitEmptyPattern(ctx *EmptyPatternContext) interface{}

	// Visit a parse tree produced by TrinoParser#patternPermutation.
	VisitPatternPermutation(ctx *PatternPermutationContext) interface{}

	// Visit a parse tree produced by TrinoParser#groupedPattern.
	VisitGroupedPattern(ctx *GroupedPatternContext) interface{}

	// Visit a parse tree produced by TrinoParser#partitionStartAnchor.
	VisitPartitionStartAnchor(ctx *PartitionStartAnchorContext) interface{}

	// Visit a parse tree produced by TrinoParser#partitionEndAnchor.
	VisitPartitionEndAnchor(ctx *PartitionEndAnchorContext) interface{}

	// Visit a parse tree produced by TrinoParser#excludedPattern.
	VisitExcludedPattern(ctx *ExcludedPatternContext) interface{}

	// Visit a parse tree produced by TrinoParser#zeroOrMoreQuantifier.
	VisitZeroOrMoreQuantifier(ctx *ZeroOrMoreQuantifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#oneOrMoreQuantifier.
	VisitOneOrMoreQuantifier(ctx *OneOrMoreQuantifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#zeroOrOneQuantifier.
	VisitZeroOrOneQuantifier(ctx *ZeroOrOneQuantifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#rangeQuantifier.
	VisitRangeQuantifier(ctx *RangeQuantifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#updateAssignment.
	VisitUpdateAssignment(ctx *UpdateAssignmentContext) interface{}

	// Visit a parse tree produced by TrinoParser#explainFormat.
	VisitExplainFormat(ctx *ExplainFormatContext) interface{}

	// Visit a parse tree produced by TrinoParser#explainType.
	VisitExplainType(ctx *ExplainTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#isolationLevel.
	VisitIsolationLevel(ctx *IsolationLevelContext) interface{}

	// Visit a parse tree produced by TrinoParser#transactionAccessMode.
	VisitTransactionAccessMode(ctx *TransactionAccessModeContext) interface{}

	// Visit a parse tree produced by TrinoParser#readUncommitted.
	VisitReadUncommitted(ctx *ReadUncommittedContext) interface{}

	// Visit a parse tree produced by TrinoParser#readCommitted.
	VisitReadCommitted(ctx *ReadCommittedContext) interface{}

	// Visit a parse tree produced by TrinoParser#repeatableRead.
	VisitRepeatableRead(ctx *RepeatableReadContext) interface{}

	// Visit a parse tree produced by TrinoParser#serializable.
	VisitSerializable(ctx *SerializableContext) interface{}

	// Visit a parse tree produced by TrinoParser#positionalArgument.
	VisitPositionalArgument(ctx *PositionalArgumentContext) interface{}

	// Visit a parse tree produced by TrinoParser#namedArgument.
	VisitNamedArgument(ctx *NamedArgumentContext) interface{}

	// Visit a parse tree produced by TrinoParser#qualifiedArgument.
	VisitQualifiedArgument(ctx *QualifiedArgumentContext) interface{}

	// Visit a parse tree produced by TrinoParser#unqualifiedArgument.
	VisitUnqualifiedArgument(ctx *UnqualifiedArgumentContext) interface{}

	// Visit a parse tree produced by TrinoParser#pathSpecification.
	VisitPathSpecification(ctx *PathSpecificationContext) interface{}

	// Visit a parse tree produced by TrinoParser#functionSpecification.
	VisitFunctionSpecification(ctx *FunctionSpecificationContext) interface{}

	// Visit a parse tree produced by TrinoParser#functionDeclaration.
	VisitFunctionDeclaration(ctx *FunctionDeclarationContext) interface{}

	// Visit a parse tree produced by TrinoParser#parameterDeclaration.
	VisitParameterDeclaration(ctx *ParameterDeclarationContext) interface{}

	// Visit a parse tree produced by TrinoParser#returnsClause.
	VisitReturnsClause(ctx *ReturnsClauseContext) interface{}

	// Visit a parse tree produced by TrinoParser#languageCharacteristic.
	VisitLanguageCharacteristic(ctx *LanguageCharacteristicContext) interface{}

	// Visit a parse tree produced by TrinoParser#deterministicCharacteristic.
	VisitDeterministicCharacteristic(ctx *DeterministicCharacteristicContext) interface{}

	// Visit a parse tree produced by TrinoParser#returnsNullOnNullInputCharacteristic.
	VisitReturnsNullOnNullInputCharacteristic(ctx *ReturnsNullOnNullInputCharacteristicContext) interface{}

	// Visit a parse tree produced by TrinoParser#calledOnNullInputCharacteristic.
	VisitCalledOnNullInputCharacteristic(ctx *CalledOnNullInputCharacteristicContext) interface{}

	// Visit a parse tree produced by TrinoParser#securityCharacteristic.
	VisitSecurityCharacteristic(ctx *SecurityCharacteristicContext) interface{}

	// Visit a parse tree produced by TrinoParser#commentCharacteristic.
	VisitCommentCharacteristic(ctx *CommentCharacteristicContext) interface{}

	// Visit a parse tree produced by TrinoParser#returnStatement.
	VisitReturnStatement(ctx *ReturnStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#assignmentStatement.
	VisitAssignmentStatement(ctx *AssignmentStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#simpleCaseStatement.
	VisitSimpleCaseStatement(ctx *SimpleCaseStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#searchedCaseStatement.
	VisitSearchedCaseStatement(ctx *SearchedCaseStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#ifStatement.
	VisitIfStatement(ctx *IfStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#iterateStatement.
	VisitIterateStatement(ctx *IterateStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#leaveStatement.
	VisitLeaveStatement(ctx *LeaveStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#compoundStatement.
	VisitCompoundStatement(ctx *CompoundStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#loopStatement.
	VisitLoopStatement(ctx *LoopStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#whileStatement.
	VisitWhileStatement(ctx *WhileStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#repeatStatement.
	VisitRepeatStatement(ctx *RepeatStatementContext) interface{}

	// Visit a parse tree produced by TrinoParser#caseStatementWhenClause.
	VisitCaseStatementWhenClause(ctx *CaseStatementWhenClauseContext) interface{}

	// Visit a parse tree produced by TrinoParser#elseIfClause.
	VisitElseIfClause(ctx *ElseIfClauseContext) interface{}

	// Visit a parse tree produced by TrinoParser#elseClause.
	VisitElseClause(ctx *ElseClauseContext) interface{}

	// Visit a parse tree produced by TrinoParser#variableDeclaration.
	VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{}

	// Visit a parse tree produced by TrinoParser#sqlStatementList.
	VisitSqlStatementList(ctx *SqlStatementListContext) interface{}

	// Visit a parse tree produced by TrinoParser#privilege.
	VisitPrivilege(ctx *PrivilegeContext) interface{}

	// Visit a parse tree produced by TrinoParser#qualifiedName.
	VisitQualifiedName(ctx *QualifiedNameContext) interface{}

	// Visit a parse tree produced by TrinoParser#queryPeriod.
	VisitQueryPeriod(ctx *QueryPeriodContext) interface{}

	// Visit a parse tree produced by TrinoParser#rangeType.
	VisitRangeType(ctx *RangeTypeContext) interface{}

	// Visit a parse tree produced by TrinoParser#specifiedPrincipal.
	VisitSpecifiedPrincipal(ctx *SpecifiedPrincipalContext) interface{}

	// Visit a parse tree produced by TrinoParser#currentUserGrantor.
	VisitCurrentUserGrantor(ctx *CurrentUserGrantorContext) interface{}

	// Visit a parse tree produced by TrinoParser#currentRoleGrantor.
	VisitCurrentRoleGrantor(ctx *CurrentRoleGrantorContext) interface{}

	// Visit a parse tree produced by TrinoParser#unspecifiedPrincipal.
	VisitUnspecifiedPrincipal(ctx *UnspecifiedPrincipalContext) interface{}

	// Visit a parse tree produced by TrinoParser#userPrincipal.
	VisitUserPrincipal(ctx *UserPrincipalContext) interface{}

	// Visit a parse tree produced by TrinoParser#rolePrincipal.
	VisitRolePrincipal(ctx *RolePrincipalContext) interface{}

	// Visit a parse tree produced by TrinoParser#roles.
	VisitRoles(ctx *RolesContext) interface{}

	// Visit a parse tree produced by TrinoParser#unquotedIdentifier.
	VisitUnquotedIdentifier(ctx *UnquotedIdentifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#quotedIdentifier.
	VisitQuotedIdentifier(ctx *QuotedIdentifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#backQuotedIdentifier.
	VisitBackQuotedIdentifier(ctx *BackQuotedIdentifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#digitIdentifier.
	VisitDigitIdentifier(ctx *DigitIdentifierContext) interface{}

	// Visit a parse tree produced by TrinoParser#decimalLiteral.
	VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#doubleLiteral.
	VisitDoubleLiteral(ctx *DoubleLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#integerLiteral.
	VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{}

	// Visit a parse tree produced by TrinoParser#identifierUser.
	VisitIdentifierUser(ctx *IdentifierUserContext) interface{}

	// Visit a parse tree produced by TrinoParser#stringUser.
	VisitStringUser(ctx *StringUserContext) interface{}

	// Visit a parse tree produced by TrinoParser#nonReserved.
	VisitNonReserved(ctx *NonReservedContext) interface{}
}
