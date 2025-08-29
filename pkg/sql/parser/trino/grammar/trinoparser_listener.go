// Code generated from TrinoParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar // TrinoParser
import "github.com/antlr4-go/antlr/v4"

// TrinoParserListener is a complete listener for a parse tree produced by TrinoParser.
type TrinoParserListener interface {
	antlr.ParseTreeListener

	// EnterParse is called when entering the parse production.
	EnterParse(c *ParseContext)

	// EnterStatements is called when entering the statements production.
	EnterStatements(c *StatementsContext)

	// EnterSingleStatement is called when entering the singleStatement production.
	EnterSingleStatement(c *SingleStatementContext)

	// EnterStandaloneExpression is called when entering the standaloneExpression production.
	EnterStandaloneExpression(c *StandaloneExpressionContext)

	// EnterStandalonePathSpecification is called when entering the standalonePathSpecification production.
	EnterStandalonePathSpecification(c *StandalonePathSpecificationContext)

	// EnterStandaloneType is called when entering the standaloneType production.
	EnterStandaloneType(c *StandaloneTypeContext)

	// EnterStandaloneRowPattern is called when entering the standaloneRowPattern production.
	EnterStandaloneRowPattern(c *StandaloneRowPatternContext)

	// EnterStandaloneFunctionSpecification is called when entering the standaloneFunctionSpecification production.
	EnterStandaloneFunctionSpecification(c *StandaloneFunctionSpecificationContext)

	// EnterStatementDefault is called when entering the statementDefault production.
	EnterStatementDefault(c *StatementDefaultContext)

	// EnterUse is called when entering the use production.
	EnterUse(c *UseContext)

	// EnterCreateCatalog is called when entering the createCatalog production.
	EnterCreateCatalog(c *CreateCatalogContext)

	// EnterDropCatalog is called when entering the dropCatalog production.
	EnterDropCatalog(c *DropCatalogContext)

	// EnterCreateSchema is called when entering the createSchema production.
	EnterCreateSchema(c *CreateSchemaContext)

	// EnterDropSchema is called when entering the dropSchema production.
	EnterDropSchema(c *DropSchemaContext)

	// EnterRenameSchema is called when entering the renameSchema production.
	EnterRenameSchema(c *RenameSchemaContext)

	// EnterSetSchemaAuthorization is called when entering the setSchemaAuthorization production.
	EnterSetSchemaAuthorization(c *SetSchemaAuthorizationContext)

	// EnterCreateTableAsSelect is called when entering the createTableAsSelect production.
	EnterCreateTableAsSelect(c *CreateTableAsSelectContext)

	// EnterCreateTable is called when entering the createTable production.
	EnterCreateTable(c *CreateTableContext)

	// EnterDropTable is called when entering the dropTable production.
	EnterDropTable(c *DropTableContext)

	// EnterInsertInto is called when entering the insertInto production.
	EnterInsertInto(c *InsertIntoContext)

	// EnterDelete is called when entering the delete production.
	EnterDelete(c *DeleteContext)

	// EnterTruncateTable is called when entering the truncateTable production.
	EnterTruncateTable(c *TruncateTableContext)

	// EnterCommentTable is called when entering the commentTable production.
	EnterCommentTable(c *CommentTableContext)

	// EnterCommentView is called when entering the commentView production.
	EnterCommentView(c *CommentViewContext)

	// EnterCommentColumn is called when entering the commentColumn production.
	EnterCommentColumn(c *CommentColumnContext)

	// EnterRenameTable is called when entering the renameTable production.
	EnterRenameTable(c *RenameTableContext)

	// EnterAddColumn is called when entering the addColumn production.
	EnterAddColumn(c *AddColumnContext)

	// EnterRenameColumn is called when entering the renameColumn production.
	EnterRenameColumn(c *RenameColumnContext)

	// EnterDropColumn is called when entering the dropColumn production.
	EnterDropColumn(c *DropColumnContext)

	// EnterSetColumnType is called when entering the setColumnType production.
	EnterSetColumnType(c *SetColumnTypeContext)

	// EnterSetTableAuthorization is called when entering the setTableAuthorization production.
	EnterSetTableAuthorization(c *SetTableAuthorizationContext)

	// EnterSetTableProperties is called when entering the setTableProperties production.
	EnterSetTableProperties(c *SetTablePropertiesContext)

	// EnterTableExecute is called when entering the tableExecute production.
	EnterTableExecute(c *TableExecuteContext)

	// EnterAnalyze is called when entering the analyze production.
	EnterAnalyze(c *AnalyzeContext)

	// EnterCreateMaterializedView is called when entering the createMaterializedView production.
	EnterCreateMaterializedView(c *CreateMaterializedViewContext)

	// EnterCreateView is called when entering the createView production.
	EnterCreateView(c *CreateViewContext)

	// EnterRefreshMaterializedView is called when entering the refreshMaterializedView production.
	EnterRefreshMaterializedView(c *RefreshMaterializedViewContext)

	// EnterDropMaterializedView is called when entering the dropMaterializedView production.
	EnterDropMaterializedView(c *DropMaterializedViewContext)

	// EnterRenameMaterializedView is called when entering the renameMaterializedView production.
	EnterRenameMaterializedView(c *RenameMaterializedViewContext)

	// EnterSetMaterializedViewProperties is called when entering the setMaterializedViewProperties production.
	EnterSetMaterializedViewProperties(c *SetMaterializedViewPropertiesContext)

	// EnterDropView is called when entering the dropView production.
	EnterDropView(c *DropViewContext)

	// EnterRenameView is called when entering the renameView production.
	EnterRenameView(c *RenameViewContext)

	// EnterSetViewAuthorization is called when entering the setViewAuthorization production.
	EnterSetViewAuthorization(c *SetViewAuthorizationContext)

	// EnterCall is called when entering the call production.
	EnterCall(c *CallContext)

	// EnterCreateFunction is called when entering the createFunction production.
	EnterCreateFunction(c *CreateFunctionContext)

	// EnterDropFunction is called when entering the dropFunction production.
	EnterDropFunction(c *DropFunctionContext)

	// EnterCreateRole is called when entering the createRole production.
	EnterCreateRole(c *CreateRoleContext)

	// EnterDropRole is called when entering the dropRole production.
	EnterDropRole(c *DropRoleContext)

	// EnterGrantRoles is called when entering the grantRoles production.
	EnterGrantRoles(c *GrantRolesContext)

	// EnterRevokeRoles is called when entering the revokeRoles production.
	EnterRevokeRoles(c *RevokeRolesContext)

	// EnterSetRole is called when entering the setRole production.
	EnterSetRole(c *SetRoleContext)

	// EnterGrant is called when entering the grant production.
	EnterGrant(c *GrantContext)

	// EnterDeny is called when entering the deny production.
	EnterDeny(c *DenyContext)

	// EnterRevoke is called when entering the revoke production.
	EnterRevoke(c *RevokeContext)

	// EnterShowGrants is called when entering the showGrants production.
	EnterShowGrants(c *ShowGrantsContext)

	// EnterExplain is called when entering the explain production.
	EnterExplain(c *ExplainContext)

	// EnterExplainAnalyze is called when entering the explainAnalyze production.
	EnterExplainAnalyze(c *ExplainAnalyzeContext)

	// EnterShowCreateTable is called when entering the showCreateTable production.
	EnterShowCreateTable(c *ShowCreateTableContext)

	// EnterShowCreateSchema is called when entering the showCreateSchema production.
	EnterShowCreateSchema(c *ShowCreateSchemaContext)

	// EnterShowCreateView is called when entering the showCreateView production.
	EnterShowCreateView(c *ShowCreateViewContext)

	// EnterShowCreateMaterializedView is called when entering the showCreateMaterializedView production.
	EnterShowCreateMaterializedView(c *ShowCreateMaterializedViewContext)

	// EnterShowTables is called when entering the showTables production.
	EnterShowTables(c *ShowTablesContext)

	// EnterShowSchemas is called when entering the showSchemas production.
	EnterShowSchemas(c *ShowSchemasContext)

	// EnterShowCatalogs is called when entering the showCatalogs production.
	EnterShowCatalogs(c *ShowCatalogsContext)

	// EnterShowColumns is called when entering the showColumns production.
	EnterShowColumns(c *ShowColumnsContext)

	// EnterShowStats is called when entering the showStats production.
	EnterShowStats(c *ShowStatsContext)

	// EnterShowStatsForQuery is called when entering the showStatsForQuery production.
	EnterShowStatsForQuery(c *ShowStatsForQueryContext)

	// EnterShowRoles is called when entering the showRoles production.
	EnterShowRoles(c *ShowRolesContext)

	// EnterShowRoleGrants is called when entering the showRoleGrants production.
	EnterShowRoleGrants(c *ShowRoleGrantsContext)

	// EnterShowFunctions is called when entering the showFunctions production.
	EnterShowFunctions(c *ShowFunctionsContext)

	// EnterShowSession is called when entering the showSession production.
	EnterShowSession(c *ShowSessionContext)

	// EnterSetSessionAuthorization is called when entering the setSessionAuthorization production.
	EnterSetSessionAuthorization(c *SetSessionAuthorizationContext)

	// EnterResetSessionAuthorization is called when entering the resetSessionAuthorization production.
	EnterResetSessionAuthorization(c *ResetSessionAuthorizationContext)

	// EnterSetSession is called when entering the setSession production.
	EnterSetSession(c *SetSessionContext)

	// EnterResetSession is called when entering the resetSession production.
	EnterResetSession(c *ResetSessionContext)

	// EnterStartTransaction is called when entering the startTransaction production.
	EnterStartTransaction(c *StartTransactionContext)

	// EnterCommit is called when entering the commit production.
	EnterCommit(c *CommitContext)

	// EnterRollback is called when entering the rollback production.
	EnterRollback(c *RollbackContext)

	// EnterPrepare is called when entering the prepare production.
	EnterPrepare(c *PrepareContext)

	// EnterDeallocate is called when entering the deallocate production.
	EnterDeallocate(c *DeallocateContext)

	// EnterExecute is called when entering the execute production.
	EnterExecute(c *ExecuteContext)

	// EnterExecuteImmediate is called when entering the executeImmediate production.
	EnterExecuteImmediate(c *ExecuteImmediateContext)

	// EnterDescribeInput is called when entering the describeInput production.
	EnterDescribeInput(c *DescribeInputContext)

	// EnterDescribeOutput is called when entering the describeOutput production.
	EnterDescribeOutput(c *DescribeOutputContext)

	// EnterSetPath is called when entering the setPath production.
	EnterSetPath(c *SetPathContext)

	// EnterSetTimeZone is called when entering the setTimeZone production.
	EnterSetTimeZone(c *SetTimeZoneContext)

	// EnterUpdate is called when entering the update production.
	EnterUpdate(c *UpdateContext)

	// EnterMerge is called when entering the merge production.
	EnterMerge(c *MergeContext)

	// EnterRootQuery is called when entering the rootQuery production.
	EnterRootQuery(c *RootQueryContext)

	// EnterWithFunction is called when entering the withFunction production.
	EnterWithFunction(c *WithFunctionContext)

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterWith is called when entering the with production.
	EnterWith(c *WithContext)

	// EnterTableElement is called when entering the tableElement production.
	EnterTableElement(c *TableElementContext)

	// EnterColumnDefinition is called when entering the columnDefinition production.
	EnterColumnDefinition(c *ColumnDefinitionContext)

	// EnterLikeClause is called when entering the likeClause production.
	EnterLikeClause(c *LikeClauseContext)

	// EnterProperties is called when entering the properties production.
	EnterProperties(c *PropertiesContext)

	// EnterPropertyAssignments is called when entering the propertyAssignments production.
	EnterPropertyAssignments(c *PropertyAssignmentsContext)

	// EnterProperty is called when entering the property production.
	EnterProperty(c *PropertyContext)

	// EnterDefaultPropertyValue is called when entering the defaultPropertyValue production.
	EnterDefaultPropertyValue(c *DefaultPropertyValueContext)

	// EnterNonDefaultPropertyValue is called when entering the nonDefaultPropertyValue production.
	EnterNonDefaultPropertyValue(c *NonDefaultPropertyValueContext)

	// EnterQueryNoWith is called when entering the queryNoWith production.
	EnterQueryNoWith(c *QueryNoWithContext)

	// EnterLimitRowCount is called when entering the limitRowCount production.
	EnterLimitRowCount(c *LimitRowCountContext)

	// EnterRowCount is called when entering the rowCount production.
	EnterRowCount(c *RowCountContext)

	// EnterQueryTermDefault is called when entering the queryTermDefault production.
	EnterQueryTermDefault(c *QueryTermDefaultContext)

	// EnterSetOperation is called when entering the setOperation production.
	EnterSetOperation(c *SetOperationContext)

	// EnterQueryPrimaryDefault is called when entering the queryPrimaryDefault production.
	EnterQueryPrimaryDefault(c *QueryPrimaryDefaultContext)

	// EnterTable is called when entering the table production.
	EnterTable(c *TableContext)

	// EnterInlineTable is called when entering the inlineTable production.
	EnterInlineTable(c *InlineTableContext)

	// EnterSubquery is called when entering the subquery production.
	EnterSubquery(c *SubqueryContext)

	// EnterSortItem is called when entering the sortItem production.
	EnterSortItem(c *SortItemContext)

	// EnterQuerySpecification is called when entering the querySpecification production.
	EnterQuerySpecification(c *QuerySpecificationContext)

	// EnterGroupBy is called when entering the groupBy production.
	EnterGroupBy(c *GroupByContext)

	// EnterSingleGroupingSet is called when entering the singleGroupingSet production.
	EnterSingleGroupingSet(c *SingleGroupingSetContext)

	// EnterRollup is called when entering the rollup production.
	EnterRollup(c *RollupContext)

	// EnterCube is called when entering the cube production.
	EnterCube(c *CubeContext)

	// EnterMultipleGroupingSets is called when entering the multipleGroupingSets production.
	EnterMultipleGroupingSets(c *MultipleGroupingSetsContext)

	// EnterGroupingSet is called when entering the groupingSet production.
	EnterGroupingSet(c *GroupingSetContext)

	// EnterWindowDefinition is called when entering the windowDefinition production.
	EnterWindowDefinition(c *WindowDefinitionContext)

	// EnterWindowSpecification is called when entering the windowSpecification production.
	EnterWindowSpecification(c *WindowSpecificationContext)

	// EnterNamedQuery is called when entering the namedQuery production.
	EnterNamedQuery(c *NamedQueryContext)

	// EnterSetQuantifier is called when entering the setQuantifier production.
	EnterSetQuantifier(c *SetQuantifierContext)

	// EnterSelectSingle is called when entering the selectSingle production.
	EnterSelectSingle(c *SelectSingleContext)

	// EnterSelectAll is called when entering the selectAll production.
	EnterSelectAll(c *SelectAllContext)

	// EnterRelationDefault is called when entering the relationDefault production.
	EnterRelationDefault(c *RelationDefaultContext)

	// EnterJoinRelation is called when entering the joinRelation production.
	EnterJoinRelation(c *JoinRelationContext)

	// EnterJoinType is called when entering the joinType production.
	EnterJoinType(c *JoinTypeContext)

	// EnterJoinCriteria is called when entering the joinCriteria production.
	EnterJoinCriteria(c *JoinCriteriaContext)

	// EnterSampledRelation is called when entering the sampledRelation production.
	EnterSampledRelation(c *SampledRelationContext)

	// EnterSampleType is called when entering the sampleType production.
	EnterSampleType(c *SampleTypeContext)

	// EnterTrimsSpecification is called when entering the trimsSpecification production.
	EnterTrimsSpecification(c *TrimsSpecificationContext)

	// EnterListAggOverflowBehavior is called when entering the listAggOverflowBehavior production.
	EnterListAggOverflowBehavior(c *ListAggOverflowBehaviorContext)

	// EnterListaggCountIndication is called when entering the listaggCountIndication production.
	EnterListaggCountIndication(c *ListaggCountIndicationContext)

	// EnterPatternRecognition is called when entering the patternRecognition production.
	EnterPatternRecognition(c *PatternRecognitionContext)

	// EnterMeasureDefinition is called when entering the measureDefinition production.
	EnterMeasureDefinition(c *MeasureDefinitionContext)

	// EnterRowsPerMatch is called when entering the rowsPerMatch production.
	EnterRowsPerMatch(c *RowsPerMatchContext)

	// EnterEmptyMatchHandling is called when entering the emptyMatchHandling production.
	EnterEmptyMatchHandling(c *EmptyMatchHandlingContext)

	// EnterSkipTo is called when entering the skipTo production.
	EnterSkipTo(c *SkipToContext)

	// EnterSubsetDefinition is called when entering the subsetDefinition production.
	EnterSubsetDefinition(c *SubsetDefinitionContext)

	// EnterVariableDefinition is called when entering the variableDefinition production.
	EnterVariableDefinition(c *VariableDefinitionContext)

	// EnterAliasedRelation is called when entering the aliasedRelation production.
	EnterAliasedRelation(c *AliasedRelationContext)

	// EnterColumnAliases is called when entering the columnAliases production.
	EnterColumnAliases(c *ColumnAliasesContext)

	// EnterTableName is called when entering the tableName production.
	EnterTableName(c *TableNameContext)

	// EnterSubqueryRelation is called when entering the subqueryRelation production.
	EnterSubqueryRelation(c *SubqueryRelationContext)

	// EnterUnnest is called when entering the unnest production.
	EnterUnnest(c *UnnestContext)

	// EnterLateral is called when entering the lateral production.
	EnterLateral(c *LateralContext)

	// EnterTableFunctionInvocation is called when entering the tableFunctionInvocation production.
	EnterTableFunctionInvocation(c *TableFunctionInvocationContext)

	// EnterParenthesizedRelation is called when entering the parenthesizedRelation production.
	EnterParenthesizedRelation(c *ParenthesizedRelationContext)

	// EnterTableFunctionCall is called when entering the tableFunctionCall production.
	EnterTableFunctionCall(c *TableFunctionCallContext)

	// EnterTableFunctionArgument is called when entering the tableFunctionArgument production.
	EnterTableFunctionArgument(c *TableFunctionArgumentContext)

	// EnterTableArgument is called when entering the tableArgument production.
	EnterTableArgument(c *TableArgumentContext)

	// EnterTableArgumentTable is called when entering the tableArgumentTable production.
	EnterTableArgumentTable(c *TableArgumentTableContext)

	// EnterTableArgumentQuery is called when entering the tableArgumentQuery production.
	EnterTableArgumentQuery(c *TableArgumentQueryContext)

	// EnterDescriptorArgument is called when entering the descriptorArgument production.
	EnterDescriptorArgument(c *DescriptorArgumentContext)

	// EnterDescriptorField is called when entering the descriptorField production.
	EnterDescriptorField(c *DescriptorFieldContext)

	// EnterCopartitionTables is called when entering the copartitionTables production.
	EnterCopartitionTables(c *CopartitionTablesContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterLogicalNot is called when entering the logicalNot production.
	EnterLogicalNot(c *LogicalNotContext)

	// EnterPredicated is called when entering the predicated production.
	EnterPredicated(c *PredicatedContext)

	// EnterOr is called when entering the or production.
	EnterOr(c *OrContext)

	// EnterAnd is called when entering the and production.
	EnterAnd(c *AndContext)

	// EnterComparison is called when entering the comparison production.
	EnterComparison(c *ComparisonContext)

	// EnterQuantifiedComparison is called when entering the quantifiedComparison production.
	EnterQuantifiedComparison(c *QuantifiedComparisonContext)

	// EnterBetween is called when entering the between production.
	EnterBetween(c *BetweenContext)

	// EnterInList is called when entering the inList production.
	EnterInList(c *InListContext)

	// EnterInSubquery is called when entering the inSubquery production.
	EnterInSubquery(c *InSubqueryContext)

	// EnterLike is called when entering the like production.
	EnterLike(c *LikeContext)

	// EnterNullPredicate is called when entering the nullPredicate production.
	EnterNullPredicate(c *NullPredicateContext)

	// EnterDistinctFrom is called when entering the distinctFrom production.
	EnterDistinctFrom(c *DistinctFromContext)

	// EnterValueExpressionDefault is called when entering the valueExpressionDefault production.
	EnterValueExpressionDefault(c *ValueExpressionDefaultContext)

	// EnterConcatenation is called when entering the concatenation production.
	EnterConcatenation(c *ConcatenationContext)

	// EnterArithmeticBinary is called when entering the arithmeticBinary production.
	EnterArithmeticBinary(c *ArithmeticBinaryContext)

	// EnterArithmeticUnary is called when entering the arithmeticUnary production.
	EnterArithmeticUnary(c *ArithmeticUnaryContext)

	// EnterAtTimeZone is called when entering the atTimeZone production.
	EnterAtTimeZone(c *AtTimeZoneContext)

	// EnterDereference is called when entering the dereference production.
	EnterDereference(c *DereferenceContext)

	// EnterTypeConstructor is called when entering the typeConstructor production.
	EnterTypeConstructor(c *TypeConstructorContext)

	// EnterJsonValue is called when entering the jsonValue production.
	EnterJsonValue(c *JsonValueContext)

	// EnterSpecialDateTimeFunction is called when entering the specialDateTimeFunction production.
	EnterSpecialDateTimeFunction(c *SpecialDateTimeFunctionContext)

	// EnterSubstring is called when entering the substring production.
	EnterSubstring(c *SubstringContext)

	// EnterCast is called when entering the cast production.
	EnterCast(c *CastContext)

	// EnterLambda is called when entering the lambda production.
	EnterLambda(c *LambdaContext)

	// EnterParenthesizedExpression is called when entering the parenthesizedExpression production.
	EnterParenthesizedExpression(c *ParenthesizedExpressionContext)

	// EnterTrim is called when entering the trim production.
	EnterTrim(c *TrimContext)

	// EnterParameter is called when entering the parameter production.
	EnterParameter(c *ParameterContext)

	// EnterNormalize is called when entering the normalize production.
	EnterNormalize(c *NormalizeContext)

	// EnterJsonObject is called when entering the jsonObject production.
	EnterJsonObject(c *JsonObjectContext)

	// EnterIntervalLiteral is called when entering the intervalLiteral production.
	EnterIntervalLiteral(c *IntervalLiteralContext)

	// EnterNumericLiteral is called when entering the numericLiteral production.
	EnterNumericLiteral(c *NumericLiteralContext)

	// EnterBooleanLiteral is called when entering the booleanLiteral production.
	EnterBooleanLiteral(c *BooleanLiteralContext)

	// EnterJsonArray is called when entering the jsonArray production.
	EnterJsonArray(c *JsonArrayContext)

	// EnterSimpleCase is called when entering the simpleCase production.
	EnterSimpleCase(c *SimpleCaseContext)

	// EnterColumnReference is called when entering the columnReference production.
	EnterColumnReference(c *ColumnReferenceContext)

	// EnterNullLiteral is called when entering the nullLiteral production.
	EnterNullLiteral(c *NullLiteralContext)

	// EnterRowConstructor is called when entering the rowConstructor production.
	EnterRowConstructor(c *RowConstructorContext)

	// EnterSubscript is called when entering the subscript production.
	EnterSubscript(c *SubscriptContext)

	// EnterJsonExists is called when entering the jsonExists production.
	EnterJsonExists(c *JsonExistsContext)

	// EnterCurrentPath is called when entering the currentPath production.
	EnterCurrentPath(c *CurrentPathContext)

	// EnterSubqueryExpression is called when entering the subqueryExpression production.
	EnterSubqueryExpression(c *SubqueryExpressionContext)

	// EnterBinaryLiteral is called when entering the binaryLiteral production.
	EnterBinaryLiteral(c *BinaryLiteralContext)

	// EnterCurrentUser is called when entering the currentUser production.
	EnterCurrentUser(c *CurrentUserContext)

	// EnterJsonQuery is called when entering the jsonQuery production.
	EnterJsonQuery(c *JsonQueryContext)

	// EnterMeasure is called when entering the measure production.
	EnterMeasure(c *MeasureContext)

	// EnterExtract is called when entering the extract production.
	EnterExtract(c *ExtractContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterArrayConstructor is called when entering the arrayConstructor production.
	EnterArrayConstructor(c *ArrayConstructorContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterCurrentSchema is called when entering the currentSchema production.
	EnterCurrentSchema(c *CurrentSchemaContext)

	// EnterExists is called when entering the exists production.
	EnterExists(c *ExistsContext)

	// EnterPosition is called when entering the position production.
	EnterPosition(c *PositionContext)

	// EnterListagg is called when entering the listagg production.
	EnterListagg(c *ListaggContext)

	// EnterSearchedCase is called when entering the searchedCase production.
	EnterSearchedCase(c *SearchedCaseContext)

	// EnterCurrentCatalog is called when entering the currentCatalog production.
	EnterCurrentCatalog(c *CurrentCatalogContext)

	// EnterGroupingOperation is called when entering the groupingOperation production.
	EnterGroupingOperation(c *GroupingOperationContext)

	// EnterJsonPathInvocation is called when entering the jsonPathInvocation production.
	EnterJsonPathInvocation(c *JsonPathInvocationContext)

	// EnterJsonValueExpression is called when entering the jsonValueExpression production.
	EnterJsonValueExpression(c *JsonValueExpressionContext)

	// EnterJsonRepresentation is called when entering the jsonRepresentation production.
	EnterJsonRepresentation(c *JsonRepresentationContext)

	// EnterJsonArgument is called when entering the jsonArgument production.
	EnterJsonArgument(c *JsonArgumentContext)

	// EnterJsonExistsErrorBehavior is called when entering the jsonExistsErrorBehavior production.
	EnterJsonExistsErrorBehavior(c *JsonExistsErrorBehaviorContext)

	// EnterJsonValueBehavior is called when entering the jsonValueBehavior production.
	EnterJsonValueBehavior(c *JsonValueBehaviorContext)

	// EnterJsonQueryWrapperBehavior is called when entering the jsonQueryWrapperBehavior production.
	EnterJsonQueryWrapperBehavior(c *JsonQueryWrapperBehaviorContext)

	// EnterJsonQueryBehavior is called when entering the jsonQueryBehavior production.
	EnterJsonQueryBehavior(c *JsonQueryBehaviorContext)

	// EnterJsonObjectMember is called when entering the jsonObjectMember production.
	EnterJsonObjectMember(c *JsonObjectMemberContext)

	// EnterProcessingMode is called when entering the processingMode production.
	EnterProcessingMode(c *ProcessingModeContext)

	// EnterNullTreatment is called when entering the nullTreatment production.
	EnterNullTreatment(c *NullTreatmentContext)

	// EnterBasicStringLiteral is called when entering the basicStringLiteral production.
	EnterBasicStringLiteral(c *BasicStringLiteralContext)

	// EnterUnicodeStringLiteral is called when entering the unicodeStringLiteral production.
	EnterUnicodeStringLiteral(c *UnicodeStringLiteralContext)

	// EnterTimeZoneInterval is called when entering the timeZoneInterval production.
	EnterTimeZoneInterval(c *TimeZoneIntervalContext)

	// EnterTimeZoneString is called when entering the timeZoneString production.
	EnterTimeZoneString(c *TimeZoneStringContext)

	// EnterComparisonOperator is called when entering the comparisonOperator production.
	EnterComparisonOperator(c *ComparisonOperatorContext)

	// EnterComparisonQuantifier is called when entering the comparisonQuantifier production.
	EnterComparisonQuantifier(c *ComparisonQuantifierContext)

	// EnterBooleanValue is called when entering the booleanValue production.
	EnterBooleanValue(c *BooleanValueContext)

	// EnterInterval is called when entering the interval production.
	EnterInterval(c *IntervalContext)

	// EnterIntervalField is called when entering the intervalField production.
	EnterIntervalField(c *IntervalFieldContext)

	// EnterNormalForm is called when entering the normalForm production.
	EnterNormalForm(c *NormalFormContext)

	// EnterRowType is called when entering the rowType production.
	EnterRowType(c *RowTypeContext)

	// EnterIntervalType is called when entering the intervalType production.
	EnterIntervalType(c *IntervalTypeContext)

	// EnterArrayType is called when entering the arrayType production.
	EnterArrayType(c *ArrayTypeContext)

	// EnterDoublePrecisionType is called when entering the doublePrecisionType production.
	EnterDoublePrecisionType(c *DoublePrecisionTypeContext)

	// EnterLegacyArrayType is called when entering the legacyArrayType production.
	EnterLegacyArrayType(c *LegacyArrayTypeContext)

	// EnterGenericType is called when entering the genericType production.
	EnterGenericType(c *GenericTypeContext)

	// EnterDateTimeType is called when entering the dateTimeType production.
	EnterDateTimeType(c *DateTimeTypeContext)

	// EnterLegacyMapType is called when entering the legacyMapType production.
	EnterLegacyMapType(c *LegacyMapTypeContext)

	// EnterRowField is called when entering the rowField production.
	EnterRowField(c *RowFieldContext)

	// EnterTypeParameter is called when entering the typeParameter production.
	EnterTypeParameter(c *TypeParameterContext)

	// EnterWhenClause is called when entering the whenClause production.
	EnterWhenClause(c *WhenClauseContext)

	// EnterFilter is called when entering the filter production.
	EnterFilter(c *FilterContext)

	// EnterMergeUpdate is called when entering the mergeUpdate production.
	EnterMergeUpdate(c *MergeUpdateContext)

	// EnterMergeDelete is called when entering the mergeDelete production.
	EnterMergeDelete(c *MergeDeleteContext)

	// EnterMergeInsert is called when entering the mergeInsert production.
	EnterMergeInsert(c *MergeInsertContext)

	// EnterOver is called when entering the over production.
	EnterOver(c *OverContext)

	// EnterWindowFrame is called when entering the windowFrame production.
	EnterWindowFrame(c *WindowFrameContext)

	// EnterFrameExtent is called when entering the frameExtent production.
	EnterFrameExtent(c *FrameExtentContext)

	// EnterUnboundedFrame is called when entering the unboundedFrame production.
	EnterUnboundedFrame(c *UnboundedFrameContext)

	// EnterCurrentRowBound is called when entering the currentRowBound production.
	EnterCurrentRowBound(c *CurrentRowBoundContext)

	// EnterBoundedFrame is called when entering the boundedFrame production.
	EnterBoundedFrame(c *BoundedFrameContext)

	// EnterQuantifiedPrimary is called when entering the quantifiedPrimary production.
	EnterQuantifiedPrimary(c *QuantifiedPrimaryContext)

	// EnterPatternConcatenation is called when entering the patternConcatenation production.
	EnterPatternConcatenation(c *PatternConcatenationContext)

	// EnterPatternAlternation is called when entering the patternAlternation production.
	EnterPatternAlternation(c *PatternAlternationContext)

	// EnterPatternVariable is called when entering the patternVariable production.
	EnterPatternVariable(c *PatternVariableContext)

	// EnterEmptyPattern is called when entering the emptyPattern production.
	EnterEmptyPattern(c *EmptyPatternContext)

	// EnterPatternPermutation is called when entering the patternPermutation production.
	EnterPatternPermutation(c *PatternPermutationContext)

	// EnterGroupedPattern is called when entering the groupedPattern production.
	EnterGroupedPattern(c *GroupedPatternContext)

	// EnterPartitionStartAnchor is called when entering the partitionStartAnchor production.
	EnterPartitionStartAnchor(c *PartitionStartAnchorContext)

	// EnterPartitionEndAnchor is called when entering the partitionEndAnchor production.
	EnterPartitionEndAnchor(c *PartitionEndAnchorContext)

	// EnterExcludedPattern is called when entering the excludedPattern production.
	EnterExcludedPattern(c *ExcludedPatternContext)

	// EnterZeroOrMoreQuantifier is called when entering the zeroOrMoreQuantifier production.
	EnterZeroOrMoreQuantifier(c *ZeroOrMoreQuantifierContext)

	// EnterOneOrMoreQuantifier is called when entering the oneOrMoreQuantifier production.
	EnterOneOrMoreQuantifier(c *OneOrMoreQuantifierContext)

	// EnterZeroOrOneQuantifier is called when entering the zeroOrOneQuantifier production.
	EnterZeroOrOneQuantifier(c *ZeroOrOneQuantifierContext)

	// EnterRangeQuantifier is called when entering the rangeQuantifier production.
	EnterRangeQuantifier(c *RangeQuantifierContext)

	// EnterUpdateAssignment is called when entering the updateAssignment production.
	EnterUpdateAssignment(c *UpdateAssignmentContext)

	// EnterExplainFormat is called when entering the explainFormat production.
	EnterExplainFormat(c *ExplainFormatContext)

	// EnterExplainType is called when entering the explainType production.
	EnterExplainType(c *ExplainTypeContext)

	// EnterIsolationLevel is called when entering the isolationLevel production.
	EnterIsolationLevel(c *IsolationLevelContext)

	// EnterTransactionAccessMode is called when entering the transactionAccessMode production.
	EnterTransactionAccessMode(c *TransactionAccessModeContext)

	// EnterReadUncommitted is called when entering the readUncommitted production.
	EnterReadUncommitted(c *ReadUncommittedContext)

	// EnterReadCommitted is called when entering the readCommitted production.
	EnterReadCommitted(c *ReadCommittedContext)

	// EnterRepeatableRead is called when entering the repeatableRead production.
	EnterRepeatableRead(c *RepeatableReadContext)

	// EnterSerializable is called when entering the serializable production.
	EnterSerializable(c *SerializableContext)

	// EnterPositionalArgument is called when entering the positionalArgument production.
	EnterPositionalArgument(c *PositionalArgumentContext)

	// EnterNamedArgument is called when entering the namedArgument production.
	EnterNamedArgument(c *NamedArgumentContext)

	// EnterQualifiedArgument is called when entering the qualifiedArgument production.
	EnterQualifiedArgument(c *QualifiedArgumentContext)

	// EnterUnqualifiedArgument is called when entering the unqualifiedArgument production.
	EnterUnqualifiedArgument(c *UnqualifiedArgumentContext)

	// EnterPathSpecification is called when entering the pathSpecification production.
	EnterPathSpecification(c *PathSpecificationContext)

	// EnterFunctionSpecification is called when entering the functionSpecification production.
	EnterFunctionSpecification(c *FunctionSpecificationContext)

	// EnterFunctionDeclaration is called when entering the functionDeclaration production.
	EnterFunctionDeclaration(c *FunctionDeclarationContext)

	// EnterParameterDeclaration is called when entering the parameterDeclaration production.
	EnterParameterDeclaration(c *ParameterDeclarationContext)

	// EnterReturnsClause is called when entering the returnsClause production.
	EnterReturnsClause(c *ReturnsClauseContext)

	// EnterLanguageCharacteristic is called when entering the languageCharacteristic production.
	EnterLanguageCharacteristic(c *LanguageCharacteristicContext)

	// EnterDeterministicCharacteristic is called when entering the deterministicCharacteristic production.
	EnterDeterministicCharacteristic(c *DeterministicCharacteristicContext)

	// EnterReturnsNullOnNullInputCharacteristic is called when entering the returnsNullOnNullInputCharacteristic production.
	EnterReturnsNullOnNullInputCharacteristic(c *ReturnsNullOnNullInputCharacteristicContext)

	// EnterCalledOnNullInputCharacteristic is called when entering the calledOnNullInputCharacteristic production.
	EnterCalledOnNullInputCharacteristic(c *CalledOnNullInputCharacteristicContext)

	// EnterSecurityCharacteristic is called when entering the securityCharacteristic production.
	EnterSecurityCharacteristic(c *SecurityCharacteristicContext)

	// EnterCommentCharacteristic is called when entering the commentCharacteristic production.
	EnterCommentCharacteristic(c *CommentCharacteristicContext)

	// EnterReturnStatement is called when entering the returnStatement production.
	EnterReturnStatement(c *ReturnStatementContext)

	// EnterAssignmentStatement is called when entering the assignmentStatement production.
	EnterAssignmentStatement(c *AssignmentStatementContext)

	// EnterSimpleCaseStatement is called when entering the simpleCaseStatement production.
	EnterSimpleCaseStatement(c *SimpleCaseStatementContext)

	// EnterSearchedCaseStatement is called when entering the searchedCaseStatement production.
	EnterSearchedCaseStatement(c *SearchedCaseStatementContext)

	// EnterIfStatement is called when entering the ifStatement production.
	EnterIfStatement(c *IfStatementContext)

	// EnterIterateStatement is called when entering the iterateStatement production.
	EnterIterateStatement(c *IterateStatementContext)

	// EnterLeaveStatement is called when entering the leaveStatement production.
	EnterLeaveStatement(c *LeaveStatementContext)

	// EnterCompoundStatement is called when entering the compoundStatement production.
	EnterCompoundStatement(c *CompoundStatementContext)

	// EnterLoopStatement is called when entering the loopStatement production.
	EnterLoopStatement(c *LoopStatementContext)

	// EnterWhileStatement is called when entering the whileStatement production.
	EnterWhileStatement(c *WhileStatementContext)

	// EnterRepeatStatement is called when entering the repeatStatement production.
	EnterRepeatStatement(c *RepeatStatementContext)

	// EnterCaseStatementWhenClause is called when entering the caseStatementWhenClause production.
	EnterCaseStatementWhenClause(c *CaseStatementWhenClauseContext)

	// EnterElseIfClause is called when entering the elseIfClause production.
	EnterElseIfClause(c *ElseIfClauseContext)

	// EnterElseClause is called when entering the elseClause production.
	EnterElseClause(c *ElseClauseContext)

	// EnterVariableDeclaration is called when entering the variableDeclaration production.
	EnterVariableDeclaration(c *VariableDeclarationContext)

	// EnterSqlStatementList is called when entering the sqlStatementList production.
	EnterSqlStatementList(c *SqlStatementListContext)

	// EnterPrivilege is called when entering the privilege production.
	EnterPrivilege(c *PrivilegeContext)

	// EnterQualifiedName is called when entering the qualifiedName production.
	EnterQualifiedName(c *QualifiedNameContext)

	// EnterQueryPeriod is called when entering the queryPeriod production.
	EnterQueryPeriod(c *QueryPeriodContext)

	// EnterRangeType is called when entering the rangeType production.
	EnterRangeType(c *RangeTypeContext)

	// EnterSpecifiedPrincipal is called when entering the specifiedPrincipal production.
	EnterSpecifiedPrincipal(c *SpecifiedPrincipalContext)

	// EnterCurrentUserGrantor is called when entering the currentUserGrantor production.
	EnterCurrentUserGrantor(c *CurrentUserGrantorContext)

	// EnterCurrentRoleGrantor is called when entering the currentRoleGrantor production.
	EnterCurrentRoleGrantor(c *CurrentRoleGrantorContext)

	// EnterUnspecifiedPrincipal is called when entering the unspecifiedPrincipal production.
	EnterUnspecifiedPrincipal(c *UnspecifiedPrincipalContext)

	// EnterUserPrincipal is called when entering the userPrincipal production.
	EnterUserPrincipal(c *UserPrincipalContext)

	// EnterRolePrincipal is called when entering the rolePrincipal production.
	EnterRolePrincipal(c *RolePrincipalContext)

	// EnterRoles is called when entering the roles production.
	EnterRoles(c *RolesContext)

	// EnterUnquotedIdentifier is called when entering the unquotedIdentifier production.
	EnterUnquotedIdentifier(c *UnquotedIdentifierContext)

	// EnterQuotedIdentifier is called when entering the quotedIdentifier production.
	EnterQuotedIdentifier(c *QuotedIdentifierContext)

	// EnterBackQuotedIdentifier is called when entering the backQuotedIdentifier production.
	EnterBackQuotedIdentifier(c *BackQuotedIdentifierContext)

	// EnterDigitIdentifier is called when entering the digitIdentifier production.
	EnterDigitIdentifier(c *DigitIdentifierContext)

	// EnterDecimalLiteral is called when entering the decimalLiteral production.
	EnterDecimalLiteral(c *DecimalLiteralContext)

	// EnterDoubleLiteral is called when entering the doubleLiteral production.
	EnterDoubleLiteral(c *DoubleLiteralContext)

	// EnterIntegerLiteral is called when entering the integerLiteral production.
	EnterIntegerLiteral(c *IntegerLiteralContext)

	// EnterIdentifierUser is called when entering the identifierUser production.
	EnterIdentifierUser(c *IdentifierUserContext)

	// EnterStringUser is called when entering the stringUser production.
	EnterStringUser(c *StringUserContext)

	// EnterNonReserved is called when entering the nonReserved production.
	EnterNonReserved(c *NonReservedContext)

	// ExitParse is called when exiting the parse production.
	ExitParse(c *ParseContext)

	// ExitStatements is called when exiting the statements production.
	ExitStatements(c *StatementsContext)

	// ExitSingleStatement is called when exiting the singleStatement production.
	ExitSingleStatement(c *SingleStatementContext)

	// ExitStandaloneExpression is called when exiting the standaloneExpression production.
	ExitStandaloneExpression(c *StandaloneExpressionContext)

	// ExitStandalonePathSpecification is called when exiting the standalonePathSpecification production.
	ExitStandalonePathSpecification(c *StandalonePathSpecificationContext)

	// ExitStandaloneType is called when exiting the standaloneType production.
	ExitStandaloneType(c *StandaloneTypeContext)

	// ExitStandaloneRowPattern is called when exiting the standaloneRowPattern production.
	ExitStandaloneRowPattern(c *StandaloneRowPatternContext)

	// ExitStandaloneFunctionSpecification is called when exiting the standaloneFunctionSpecification production.
	ExitStandaloneFunctionSpecification(c *StandaloneFunctionSpecificationContext)

	// ExitStatementDefault is called when exiting the statementDefault production.
	ExitStatementDefault(c *StatementDefaultContext)

	// ExitUse is called when exiting the use production.
	ExitUse(c *UseContext)

	// ExitCreateCatalog is called when exiting the createCatalog production.
	ExitCreateCatalog(c *CreateCatalogContext)

	// ExitDropCatalog is called when exiting the dropCatalog production.
	ExitDropCatalog(c *DropCatalogContext)

	// ExitCreateSchema is called when exiting the createSchema production.
	ExitCreateSchema(c *CreateSchemaContext)

	// ExitDropSchema is called when exiting the dropSchema production.
	ExitDropSchema(c *DropSchemaContext)

	// ExitRenameSchema is called when exiting the renameSchema production.
	ExitRenameSchema(c *RenameSchemaContext)

	// ExitSetSchemaAuthorization is called when exiting the setSchemaAuthorization production.
	ExitSetSchemaAuthorization(c *SetSchemaAuthorizationContext)

	// ExitCreateTableAsSelect is called when exiting the createTableAsSelect production.
	ExitCreateTableAsSelect(c *CreateTableAsSelectContext)

	// ExitCreateTable is called when exiting the createTable production.
	ExitCreateTable(c *CreateTableContext)

	// ExitDropTable is called when exiting the dropTable production.
	ExitDropTable(c *DropTableContext)

	// ExitInsertInto is called when exiting the insertInto production.
	ExitInsertInto(c *InsertIntoContext)

	// ExitDelete is called when exiting the delete production.
	ExitDelete(c *DeleteContext)

	// ExitTruncateTable is called when exiting the truncateTable production.
	ExitTruncateTable(c *TruncateTableContext)

	// ExitCommentTable is called when exiting the commentTable production.
	ExitCommentTable(c *CommentTableContext)

	// ExitCommentView is called when exiting the commentView production.
	ExitCommentView(c *CommentViewContext)

	// ExitCommentColumn is called when exiting the commentColumn production.
	ExitCommentColumn(c *CommentColumnContext)

	// ExitRenameTable is called when exiting the renameTable production.
	ExitRenameTable(c *RenameTableContext)

	// ExitAddColumn is called when exiting the addColumn production.
	ExitAddColumn(c *AddColumnContext)

	// ExitRenameColumn is called when exiting the renameColumn production.
	ExitRenameColumn(c *RenameColumnContext)

	// ExitDropColumn is called when exiting the dropColumn production.
	ExitDropColumn(c *DropColumnContext)

	// ExitSetColumnType is called when exiting the setColumnType production.
	ExitSetColumnType(c *SetColumnTypeContext)

	// ExitSetTableAuthorization is called when exiting the setTableAuthorization production.
	ExitSetTableAuthorization(c *SetTableAuthorizationContext)

	// ExitSetTableProperties is called when exiting the setTableProperties production.
	ExitSetTableProperties(c *SetTablePropertiesContext)

	// ExitTableExecute is called when exiting the tableExecute production.
	ExitTableExecute(c *TableExecuteContext)

	// ExitAnalyze is called when exiting the analyze production.
	ExitAnalyze(c *AnalyzeContext)

	// ExitCreateMaterializedView is called when exiting the createMaterializedView production.
	ExitCreateMaterializedView(c *CreateMaterializedViewContext)

	// ExitCreateView is called when exiting the createView production.
	ExitCreateView(c *CreateViewContext)

	// ExitRefreshMaterializedView is called when exiting the refreshMaterializedView production.
	ExitRefreshMaterializedView(c *RefreshMaterializedViewContext)

	// ExitDropMaterializedView is called when exiting the dropMaterializedView production.
	ExitDropMaterializedView(c *DropMaterializedViewContext)

	// ExitRenameMaterializedView is called when exiting the renameMaterializedView production.
	ExitRenameMaterializedView(c *RenameMaterializedViewContext)

	// ExitSetMaterializedViewProperties is called when exiting the setMaterializedViewProperties production.
	ExitSetMaterializedViewProperties(c *SetMaterializedViewPropertiesContext)

	// ExitDropView is called when exiting the dropView production.
	ExitDropView(c *DropViewContext)

	// ExitRenameView is called when exiting the renameView production.
	ExitRenameView(c *RenameViewContext)

	// ExitSetViewAuthorization is called when exiting the setViewAuthorization production.
	ExitSetViewAuthorization(c *SetViewAuthorizationContext)

	// ExitCall is called when exiting the call production.
	ExitCall(c *CallContext)

	// ExitCreateFunction is called when exiting the createFunction production.
	ExitCreateFunction(c *CreateFunctionContext)

	// ExitDropFunction is called when exiting the dropFunction production.
	ExitDropFunction(c *DropFunctionContext)

	// ExitCreateRole is called when exiting the createRole production.
	ExitCreateRole(c *CreateRoleContext)

	// ExitDropRole is called when exiting the dropRole production.
	ExitDropRole(c *DropRoleContext)

	// ExitGrantRoles is called when exiting the grantRoles production.
	ExitGrantRoles(c *GrantRolesContext)

	// ExitRevokeRoles is called when exiting the revokeRoles production.
	ExitRevokeRoles(c *RevokeRolesContext)

	// ExitSetRole is called when exiting the setRole production.
	ExitSetRole(c *SetRoleContext)

	// ExitGrant is called when exiting the grant production.
	ExitGrant(c *GrantContext)

	// ExitDeny is called when exiting the deny production.
	ExitDeny(c *DenyContext)

	// ExitRevoke is called when exiting the revoke production.
	ExitRevoke(c *RevokeContext)

	// ExitShowGrants is called when exiting the showGrants production.
	ExitShowGrants(c *ShowGrantsContext)

	// ExitExplain is called when exiting the explain production.
	ExitExplain(c *ExplainContext)

	// ExitExplainAnalyze is called when exiting the explainAnalyze production.
	ExitExplainAnalyze(c *ExplainAnalyzeContext)

	// ExitShowCreateTable is called when exiting the showCreateTable production.
	ExitShowCreateTable(c *ShowCreateTableContext)

	// ExitShowCreateSchema is called when exiting the showCreateSchema production.
	ExitShowCreateSchema(c *ShowCreateSchemaContext)

	// ExitShowCreateView is called when exiting the showCreateView production.
	ExitShowCreateView(c *ShowCreateViewContext)

	// ExitShowCreateMaterializedView is called when exiting the showCreateMaterializedView production.
	ExitShowCreateMaterializedView(c *ShowCreateMaterializedViewContext)

	// ExitShowTables is called when exiting the showTables production.
	ExitShowTables(c *ShowTablesContext)

	// ExitShowSchemas is called when exiting the showSchemas production.
	ExitShowSchemas(c *ShowSchemasContext)

	// ExitShowCatalogs is called when exiting the showCatalogs production.
	ExitShowCatalogs(c *ShowCatalogsContext)

	// ExitShowColumns is called when exiting the showColumns production.
	ExitShowColumns(c *ShowColumnsContext)

	// ExitShowStats is called when exiting the showStats production.
	ExitShowStats(c *ShowStatsContext)

	// ExitShowStatsForQuery is called when exiting the showStatsForQuery production.
	ExitShowStatsForQuery(c *ShowStatsForQueryContext)

	// ExitShowRoles is called when exiting the showRoles production.
	ExitShowRoles(c *ShowRolesContext)

	// ExitShowRoleGrants is called when exiting the showRoleGrants production.
	ExitShowRoleGrants(c *ShowRoleGrantsContext)

	// ExitShowFunctions is called when exiting the showFunctions production.
	ExitShowFunctions(c *ShowFunctionsContext)

	// ExitShowSession is called when exiting the showSession production.
	ExitShowSession(c *ShowSessionContext)

	// ExitSetSessionAuthorization is called when exiting the setSessionAuthorization production.
	ExitSetSessionAuthorization(c *SetSessionAuthorizationContext)

	// ExitResetSessionAuthorization is called when exiting the resetSessionAuthorization production.
	ExitResetSessionAuthorization(c *ResetSessionAuthorizationContext)

	// ExitSetSession is called when exiting the setSession production.
	ExitSetSession(c *SetSessionContext)

	// ExitResetSession is called when exiting the resetSession production.
	ExitResetSession(c *ResetSessionContext)

	// ExitStartTransaction is called when exiting the startTransaction production.
	ExitStartTransaction(c *StartTransactionContext)

	// ExitCommit is called when exiting the commit production.
	ExitCommit(c *CommitContext)

	// ExitRollback is called when exiting the rollback production.
	ExitRollback(c *RollbackContext)

	// ExitPrepare is called when exiting the prepare production.
	ExitPrepare(c *PrepareContext)

	// ExitDeallocate is called when exiting the deallocate production.
	ExitDeallocate(c *DeallocateContext)

	// ExitExecute is called when exiting the execute production.
	ExitExecute(c *ExecuteContext)

	// ExitExecuteImmediate is called when exiting the executeImmediate production.
	ExitExecuteImmediate(c *ExecuteImmediateContext)

	// ExitDescribeInput is called when exiting the describeInput production.
	ExitDescribeInput(c *DescribeInputContext)

	// ExitDescribeOutput is called when exiting the describeOutput production.
	ExitDescribeOutput(c *DescribeOutputContext)

	// ExitSetPath is called when exiting the setPath production.
	ExitSetPath(c *SetPathContext)

	// ExitSetTimeZone is called when exiting the setTimeZone production.
	ExitSetTimeZone(c *SetTimeZoneContext)

	// ExitUpdate is called when exiting the update production.
	ExitUpdate(c *UpdateContext)

	// ExitMerge is called when exiting the merge production.
	ExitMerge(c *MergeContext)

	// ExitRootQuery is called when exiting the rootQuery production.
	ExitRootQuery(c *RootQueryContext)

	// ExitWithFunction is called when exiting the withFunction production.
	ExitWithFunction(c *WithFunctionContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitWith is called when exiting the with production.
	ExitWith(c *WithContext)

	// ExitTableElement is called when exiting the tableElement production.
	ExitTableElement(c *TableElementContext)

	// ExitColumnDefinition is called when exiting the columnDefinition production.
	ExitColumnDefinition(c *ColumnDefinitionContext)

	// ExitLikeClause is called when exiting the likeClause production.
	ExitLikeClause(c *LikeClauseContext)

	// ExitProperties is called when exiting the properties production.
	ExitProperties(c *PropertiesContext)

	// ExitPropertyAssignments is called when exiting the propertyAssignments production.
	ExitPropertyAssignments(c *PropertyAssignmentsContext)

	// ExitProperty is called when exiting the property production.
	ExitProperty(c *PropertyContext)

	// ExitDefaultPropertyValue is called when exiting the defaultPropertyValue production.
	ExitDefaultPropertyValue(c *DefaultPropertyValueContext)

	// ExitNonDefaultPropertyValue is called when exiting the nonDefaultPropertyValue production.
	ExitNonDefaultPropertyValue(c *NonDefaultPropertyValueContext)

	// ExitQueryNoWith is called when exiting the queryNoWith production.
	ExitQueryNoWith(c *QueryNoWithContext)

	// ExitLimitRowCount is called when exiting the limitRowCount production.
	ExitLimitRowCount(c *LimitRowCountContext)

	// ExitRowCount is called when exiting the rowCount production.
	ExitRowCount(c *RowCountContext)

	// ExitQueryTermDefault is called when exiting the queryTermDefault production.
	ExitQueryTermDefault(c *QueryTermDefaultContext)

	// ExitSetOperation is called when exiting the setOperation production.
	ExitSetOperation(c *SetOperationContext)

	// ExitQueryPrimaryDefault is called when exiting the queryPrimaryDefault production.
	ExitQueryPrimaryDefault(c *QueryPrimaryDefaultContext)

	// ExitTable is called when exiting the table production.
	ExitTable(c *TableContext)

	// ExitInlineTable is called when exiting the inlineTable production.
	ExitInlineTable(c *InlineTableContext)

	// ExitSubquery is called when exiting the subquery production.
	ExitSubquery(c *SubqueryContext)

	// ExitSortItem is called when exiting the sortItem production.
	ExitSortItem(c *SortItemContext)

	// ExitQuerySpecification is called when exiting the querySpecification production.
	ExitQuerySpecification(c *QuerySpecificationContext)

	// ExitGroupBy is called when exiting the groupBy production.
	ExitGroupBy(c *GroupByContext)

	// ExitSingleGroupingSet is called when exiting the singleGroupingSet production.
	ExitSingleGroupingSet(c *SingleGroupingSetContext)

	// ExitRollup is called when exiting the rollup production.
	ExitRollup(c *RollupContext)

	// ExitCube is called when exiting the cube production.
	ExitCube(c *CubeContext)

	// ExitMultipleGroupingSets is called when exiting the multipleGroupingSets production.
	ExitMultipleGroupingSets(c *MultipleGroupingSetsContext)

	// ExitGroupingSet is called when exiting the groupingSet production.
	ExitGroupingSet(c *GroupingSetContext)

	// ExitWindowDefinition is called when exiting the windowDefinition production.
	ExitWindowDefinition(c *WindowDefinitionContext)

	// ExitWindowSpecification is called when exiting the windowSpecification production.
	ExitWindowSpecification(c *WindowSpecificationContext)

	// ExitNamedQuery is called when exiting the namedQuery production.
	ExitNamedQuery(c *NamedQueryContext)

	// ExitSetQuantifier is called when exiting the setQuantifier production.
	ExitSetQuantifier(c *SetQuantifierContext)

	// ExitSelectSingle is called when exiting the selectSingle production.
	ExitSelectSingle(c *SelectSingleContext)

	// ExitSelectAll is called when exiting the selectAll production.
	ExitSelectAll(c *SelectAllContext)

	// ExitRelationDefault is called when exiting the relationDefault production.
	ExitRelationDefault(c *RelationDefaultContext)

	// ExitJoinRelation is called when exiting the joinRelation production.
	ExitJoinRelation(c *JoinRelationContext)

	// ExitJoinType is called when exiting the joinType production.
	ExitJoinType(c *JoinTypeContext)

	// ExitJoinCriteria is called when exiting the joinCriteria production.
	ExitJoinCriteria(c *JoinCriteriaContext)

	// ExitSampledRelation is called when exiting the sampledRelation production.
	ExitSampledRelation(c *SampledRelationContext)

	// ExitSampleType is called when exiting the sampleType production.
	ExitSampleType(c *SampleTypeContext)

	// ExitTrimsSpecification is called when exiting the trimsSpecification production.
	ExitTrimsSpecification(c *TrimsSpecificationContext)

	// ExitListAggOverflowBehavior is called when exiting the listAggOverflowBehavior production.
	ExitListAggOverflowBehavior(c *ListAggOverflowBehaviorContext)

	// ExitListaggCountIndication is called when exiting the listaggCountIndication production.
	ExitListaggCountIndication(c *ListaggCountIndicationContext)

	// ExitPatternRecognition is called when exiting the patternRecognition production.
	ExitPatternRecognition(c *PatternRecognitionContext)

	// ExitMeasureDefinition is called when exiting the measureDefinition production.
	ExitMeasureDefinition(c *MeasureDefinitionContext)

	// ExitRowsPerMatch is called when exiting the rowsPerMatch production.
	ExitRowsPerMatch(c *RowsPerMatchContext)

	// ExitEmptyMatchHandling is called when exiting the emptyMatchHandling production.
	ExitEmptyMatchHandling(c *EmptyMatchHandlingContext)

	// ExitSkipTo is called when exiting the skipTo production.
	ExitSkipTo(c *SkipToContext)

	// ExitSubsetDefinition is called when exiting the subsetDefinition production.
	ExitSubsetDefinition(c *SubsetDefinitionContext)

	// ExitVariableDefinition is called when exiting the variableDefinition production.
	ExitVariableDefinition(c *VariableDefinitionContext)

	// ExitAliasedRelation is called when exiting the aliasedRelation production.
	ExitAliasedRelation(c *AliasedRelationContext)

	// ExitColumnAliases is called when exiting the columnAliases production.
	ExitColumnAliases(c *ColumnAliasesContext)

	// ExitTableName is called when exiting the tableName production.
	ExitTableName(c *TableNameContext)

	// ExitSubqueryRelation is called when exiting the subqueryRelation production.
	ExitSubqueryRelation(c *SubqueryRelationContext)

	// ExitUnnest is called when exiting the unnest production.
	ExitUnnest(c *UnnestContext)

	// ExitLateral is called when exiting the lateral production.
	ExitLateral(c *LateralContext)

	// ExitTableFunctionInvocation is called when exiting the tableFunctionInvocation production.
	ExitTableFunctionInvocation(c *TableFunctionInvocationContext)

	// ExitParenthesizedRelation is called when exiting the parenthesizedRelation production.
	ExitParenthesizedRelation(c *ParenthesizedRelationContext)

	// ExitTableFunctionCall is called when exiting the tableFunctionCall production.
	ExitTableFunctionCall(c *TableFunctionCallContext)

	// ExitTableFunctionArgument is called when exiting the tableFunctionArgument production.
	ExitTableFunctionArgument(c *TableFunctionArgumentContext)

	// ExitTableArgument is called when exiting the tableArgument production.
	ExitTableArgument(c *TableArgumentContext)

	// ExitTableArgumentTable is called when exiting the tableArgumentTable production.
	ExitTableArgumentTable(c *TableArgumentTableContext)

	// ExitTableArgumentQuery is called when exiting the tableArgumentQuery production.
	ExitTableArgumentQuery(c *TableArgumentQueryContext)

	// ExitDescriptorArgument is called when exiting the descriptorArgument production.
	ExitDescriptorArgument(c *DescriptorArgumentContext)

	// ExitDescriptorField is called when exiting the descriptorField production.
	ExitDescriptorField(c *DescriptorFieldContext)

	// ExitCopartitionTables is called when exiting the copartitionTables production.
	ExitCopartitionTables(c *CopartitionTablesContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitLogicalNot is called when exiting the logicalNot production.
	ExitLogicalNot(c *LogicalNotContext)

	// ExitPredicated is called when exiting the predicated production.
	ExitPredicated(c *PredicatedContext)

	// ExitOr is called when exiting the or production.
	ExitOr(c *OrContext)

	// ExitAnd is called when exiting the and production.
	ExitAnd(c *AndContext)

	// ExitComparison is called when exiting the comparison production.
	ExitComparison(c *ComparisonContext)

	// ExitQuantifiedComparison is called when exiting the quantifiedComparison production.
	ExitQuantifiedComparison(c *QuantifiedComparisonContext)

	// ExitBetween is called when exiting the between production.
	ExitBetween(c *BetweenContext)

	// ExitInList is called when exiting the inList production.
	ExitInList(c *InListContext)

	// ExitInSubquery is called when exiting the inSubquery production.
	ExitInSubquery(c *InSubqueryContext)

	// ExitLike is called when exiting the like production.
	ExitLike(c *LikeContext)

	// ExitNullPredicate is called when exiting the nullPredicate production.
	ExitNullPredicate(c *NullPredicateContext)

	// ExitDistinctFrom is called when exiting the distinctFrom production.
	ExitDistinctFrom(c *DistinctFromContext)

	// ExitValueExpressionDefault is called when exiting the valueExpressionDefault production.
	ExitValueExpressionDefault(c *ValueExpressionDefaultContext)

	// ExitConcatenation is called when exiting the concatenation production.
	ExitConcatenation(c *ConcatenationContext)

	// ExitArithmeticBinary is called when exiting the arithmeticBinary production.
	ExitArithmeticBinary(c *ArithmeticBinaryContext)

	// ExitArithmeticUnary is called when exiting the arithmeticUnary production.
	ExitArithmeticUnary(c *ArithmeticUnaryContext)

	// ExitAtTimeZone is called when exiting the atTimeZone production.
	ExitAtTimeZone(c *AtTimeZoneContext)

	// ExitDereference is called when exiting the dereference production.
	ExitDereference(c *DereferenceContext)

	// ExitTypeConstructor is called when exiting the typeConstructor production.
	ExitTypeConstructor(c *TypeConstructorContext)

	// ExitJsonValue is called when exiting the jsonValue production.
	ExitJsonValue(c *JsonValueContext)

	// ExitSpecialDateTimeFunction is called when exiting the specialDateTimeFunction production.
	ExitSpecialDateTimeFunction(c *SpecialDateTimeFunctionContext)

	// ExitSubstring is called when exiting the substring production.
	ExitSubstring(c *SubstringContext)

	// ExitCast is called when exiting the cast production.
	ExitCast(c *CastContext)

	// ExitLambda is called when exiting the lambda production.
	ExitLambda(c *LambdaContext)

	// ExitParenthesizedExpression is called when exiting the parenthesizedExpression production.
	ExitParenthesizedExpression(c *ParenthesizedExpressionContext)

	// ExitTrim is called when exiting the trim production.
	ExitTrim(c *TrimContext)

	// ExitParameter is called when exiting the parameter production.
	ExitParameter(c *ParameterContext)

	// ExitNormalize is called when exiting the normalize production.
	ExitNormalize(c *NormalizeContext)

	// ExitJsonObject is called when exiting the jsonObject production.
	ExitJsonObject(c *JsonObjectContext)

	// ExitIntervalLiteral is called when exiting the intervalLiteral production.
	ExitIntervalLiteral(c *IntervalLiteralContext)

	// ExitNumericLiteral is called when exiting the numericLiteral production.
	ExitNumericLiteral(c *NumericLiteralContext)

	// ExitBooleanLiteral is called when exiting the booleanLiteral production.
	ExitBooleanLiteral(c *BooleanLiteralContext)

	// ExitJsonArray is called when exiting the jsonArray production.
	ExitJsonArray(c *JsonArrayContext)

	// ExitSimpleCase is called when exiting the simpleCase production.
	ExitSimpleCase(c *SimpleCaseContext)

	// ExitColumnReference is called when exiting the columnReference production.
	ExitColumnReference(c *ColumnReferenceContext)

	// ExitNullLiteral is called when exiting the nullLiteral production.
	ExitNullLiteral(c *NullLiteralContext)

	// ExitRowConstructor is called when exiting the rowConstructor production.
	ExitRowConstructor(c *RowConstructorContext)

	// ExitSubscript is called when exiting the subscript production.
	ExitSubscript(c *SubscriptContext)

	// ExitJsonExists is called when exiting the jsonExists production.
	ExitJsonExists(c *JsonExistsContext)

	// ExitCurrentPath is called when exiting the currentPath production.
	ExitCurrentPath(c *CurrentPathContext)

	// ExitSubqueryExpression is called when exiting the subqueryExpression production.
	ExitSubqueryExpression(c *SubqueryExpressionContext)

	// ExitBinaryLiteral is called when exiting the binaryLiteral production.
	ExitBinaryLiteral(c *BinaryLiteralContext)

	// ExitCurrentUser is called when exiting the currentUser production.
	ExitCurrentUser(c *CurrentUserContext)

	// ExitJsonQuery is called when exiting the jsonQuery production.
	ExitJsonQuery(c *JsonQueryContext)

	// ExitMeasure is called when exiting the measure production.
	ExitMeasure(c *MeasureContext)

	// ExitExtract is called when exiting the extract production.
	ExitExtract(c *ExtractContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitArrayConstructor is called when exiting the arrayConstructor production.
	ExitArrayConstructor(c *ArrayConstructorContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitCurrentSchema is called when exiting the currentSchema production.
	ExitCurrentSchema(c *CurrentSchemaContext)

	// ExitExists is called when exiting the exists production.
	ExitExists(c *ExistsContext)

	// ExitPosition is called when exiting the position production.
	ExitPosition(c *PositionContext)

	// ExitListagg is called when exiting the listagg production.
	ExitListagg(c *ListaggContext)

	// ExitSearchedCase is called when exiting the searchedCase production.
	ExitSearchedCase(c *SearchedCaseContext)

	// ExitCurrentCatalog is called when exiting the currentCatalog production.
	ExitCurrentCatalog(c *CurrentCatalogContext)

	// ExitGroupingOperation is called when exiting the groupingOperation production.
	ExitGroupingOperation(c *GroupingOperationContext)

	// ExitJsonPathInvocation is called when exiting the jsonPathInvocation production.
	ExitJsonPathInvocation(c *JsonPathInvocationContext)

	// ExitJsonValueExpression is called when exiting the jsonValueExpression production.
	ExitJsonValueExpression(c *JsonValueExpressionContext)

	// ExitJsonRepresentation is called when exiting the jsonRepresentation production.
	ExitJsonRepresentation(c *JsonRepresentationContext)

	// ExitJsonArgument is called when exiting the jsonArgument production.
	ExitJsonArgument(c *JsonArgumentContext)

	// ExitJsonExistsErrorBehavior is called when exiting the jsonExistsErrorBehavior production.
	ExitJsonExistsErrorBehavior(c *JsonExistsErrorBehaviorContext)

	// ExitJsonValueBehavior is called when exiting the jsonValueBehavior production.
	ExitJsonValueBehavior(c *JsonValueBehaviorContext)

	// ExitJsonQueryWrapperBehavior is called when exiting the jsonQueryWrapperBehavior production.
	ExitJsonQueryWrapperBehavior(c *JsonQueryWrapperBehaviorContext)

	// ExitJsonQueryBehavior is called when exiting the jsonQueryBehavior production.
	ExitJsonQueryBehavior(c *JsonQueryBehaviorContext)

	// ExitJsonObjectMember is called when exiting the jsonObjectMember production.
	ExitJsonObjectMember(c *JsonObjectMemberContext)

	// ExitProcessingMode is called when exiting the processingMode production.
	ExitProcessingMode(c *ProcessingModeContext)

	// ExitNullTreatment is called when exiting the nullTreatment production.
	ExitNullTreatment(c *NullTreatmentContext)

	// ExitBasicStringLiteral is called when exiting the basicStringLiteral production.
	ExitBasicStringLiteral(c *BasicStringLiteralContext)

	// ExitUnicodeStringLiteral is called when exiting the unicodeStringLiteral production.
	ExitUnicodeStringLiteral(c *UnicodeStringLiteralContext)

	// ExitTimeZoneInterval is called when exiting the timeZoneInterval production.
	ExitTimeZoneInterval(c *TimeZoneIntervalContext)

	// ExitTimeZoneString is called when exiting the timeZoneString production.
	ExitTimeZoneString(c *TimeZoneStringContext)

	// ExitComparisonOperator is called when exiting the comparisonOperator production.
	ExitComparisonOperator(c *ComparisonOperatorContext)

	// ExitComparisonQuantifier is called when exiting the comparisonQuantifier production.
	ExitComparisonQuantifier(c *ComparisonQuantifierContext)

	// ExitBooleanValue is called when exiting the booleanValue production.
	ExitBooleanValue(c *BooleanValueContext)

	// ExitInterval is called when exiting the interval production.
	ExitInterval(c *IntervalContext)

	// ExitIntervalField is called when exiting the intervalField production.
	ExitIntervalField(c *IntervalFieldContext)

	// ExitNormalForm is called when exiting the normalForm production.
	ExitNormalForm(c *NormalFormContext)

	// ExitRowType is called when exiting the rowType production.
	ExitRowType(c *RowTypeContext)

	// ExitIntervalType is called when exiting the intervalType production.
	ExitIntervalType(c *IntervalTypeContext)

	// ExitArrayType is called when exiting the arrayType production.
	ExitArrayType(c *ArrayTypeContext)

	// ExitDoublePrecisionType is called when exiting the doublePrecisionType production.
	ExitDoublePrecisionType(c *DoublePrecisionTypeContext)

	// ExitLegacyArrayType is called when exiting the legacyArrayType production.
	ExitLegacyArrayType(c *LegacyArrayTypeContext)

	// ExitGenericType is called when exiting the genericType production.
	ExitGenericType(c *GenericTypeContext)

	// ExitDateTimeType is called when exiting the dateTimeType production.
	ExitDateTimeType(c *DateTimeTypeContext)

	// ExitLegacyMapType is called when exiting the legacyMapType production.
	ExitLegacyMapType(c *LegacyMapTypeContext)

	// ExitRowField is called when exiting the rowField production.
	ExitRowField(c *RowFieldContext)

	// ExitTypeParameter is called when exiting the typeParameter production.
	ExitTypeParameter(c *TypeParameterContext)

	// ExitWhenClause is called when exiting the whenClause production.
	ExitWhenClause(c *WhenClauseContext)

	// ExitFilter is called when exiting the filter production.
	ExitFilter(c *FilterContext)

	// ExitMergeUpdate is called when exiting the mergeUpdate production.
	ExitMergeUpdate(c *MergeUpdateContext)

	// ExitMergeDelete is called when exiting the mergeDelete production.
	ExitMergeDelete(c *MergeDeleteContext)

	// ExitMergeInsert is called when exiting the mergeInsert production.
	ExitMergeInsert(c *MergeInsertContext)

	// ExitOver is called when exiting the over production.
	ExitOver(c *OverContext)

	// ExitWindowFrame is called when exiting the windowFrame production.
	ExitWindowFrame(c *WindowFrameContext)

	// ExitFrameExtent is called when exiting the frameExtent production.
	ExitFrameExtent(c *FrameExtentContext)

	// ExitUnboundedFrame is called when exiting the unboundedFrame production.
	ExitUnboundedFrame(c *UnboundedFrameContext)

	// ExitCurrentRowBound is called when exiting the currentRowBound production.
	ExitCurrentRowBound(c *CurrentRowBoundContext)

	// ExitBoundedFrame is called when exiting the boundedFrame production.
	ExitBoundedFrame(c *BoundedFrameContext)

	// ExitQuantifiedPrimary is called when exiting the quantifiedPrimary production.
	ExitQuantifiedPrimary(c *QuantifiedPrimaryContext)

	// ExitPatternConcatenation is called when exiting the patternConcatenation production.
	ExitPatternConcatenation(c *PatternConcatenationContext)

	// ExitPatternAlternation is called when exiting the patternAlternation production.
	ExitPatternAlternation(c *PatternAlternationContext)

	// ExitPatternVariable is called when exiting the patternVariable production.
	ExitPatternVariable(c *PatternVariableContext)

	// ExitEmptyPattern is called when exiting the emptyPattern production.
	ExitEmptyPattern(c *EmptyPatternContext)

	// ExitPatternPermutation is called when exiting the patternPermutation production.
	ExitPatternPermutation(c *PatternPermutationContext)

	// ExitGroupedPattern is called when exiting the groupedPattern production.
	ExitGroupedPattern(c *GroupedPatternContext)

	// ExitPartitionStartAnchor is called when exiting the partitionStartAnchor production.
	ExitPartitionStartAnchor(c *PartitionStartAnchorContext)

	// ExitPartitionEndAnchor is called when exiting the partitionEndAnchor production.
	ExitPartitionEndAnchor(c *PartitionEndAnchorContext)

	// ExitExcludedPattern is called when exiting the excludedPattern production.
	ExitExcludedPattern(c *ExcludedPatternContext)

	// ExitZeroOrMoreQuantifier is called when exiting the zeroOrMoreQuantifier production.
	ExitZeroOrMoreQuantifier(c *ZeroOrMoreQuantifierContext)

	// ExitOneOrMoreQuantifier is called when exiting the oneOrMoreQuantifier production.
	ExitOneOrMoreQuantifier(c *OneOrMoreQuantifierContext)

	// ExitZeroOrOneQuantifier is called when exiting the zeroOrOneQuantifier production.
	ExitZeroOrOneQuantifier(c *ZeroOrOneQuantifierContext)

	// ExitRangeQuantifier is called when exiting the rangeQuantifier production.
	ExitRangeQuantifier(c *RangeQuantifierContext)

	// ExitUpdateAssignment is called when exiting the updateAssignment production.
	ExitUpdateAssignment(c *UpdateAssignmentContext)

	// ExitExplainFormat is called when exiting the explainFormat production.
	ExitExplainFormat(c *ExplainFormatContext)

	// ExitExplainType is called when exiting the explainType production.
	ExitExplainType(c *ExplainTypeContext)

	// ExitIsolationLevel is called when exiting the isolationLevel production.
	ExitIsolationLevel(c *IsolationLevelContext)

	// ExitTransactionAccessMode is called when exiting the transactionAccessMode production.
	ExitTransactionAccessMode(c *TransactionAccessModeContext)

	// ExitReadUncommitted is called when exiting the readUncommitted production.
	ExitReadUncommitted(c *ReadUncommittedContext)

	// ExitReadCommitted is called when exiting the readCommitted production.
	ExitReadCommitted(c *ReadCommittedContext)

	// ExitRepeatableRead is called when exiting the repeatableRead production.
	ExitRepeatableRead(c *RepeatableReadContext)

	// ExitSerializable is called when exiting the serializable production.
	ExitSerializable(c *SerializableContext)

	// ExitPositionalArgument is called when exiting the positionalArgument production.
	ExitPositionalArgument(c *PositionalArgumentContext)

	// ExitNamedArgument is called when exiting the namedArgument production.
	ExitNamedArgument(c *NamedArgumentContext)

	// ExitQualifiedArgument is called when exiting the qualifiedArgument production.
	ExitQualifiedArgument(c *QualifiedArgumentContext)

	// ExitUnqualifiedArgument is called when exiting the unqualifiedArgument production.
	ExitUnqualifiedArgument(c *UnqualifiedArgumentContext)

	// ExitPathSpecification is called when exiting the pathSpecification production.
	ExitPathSpecification(c *PathSpecificationContext)

	// ExitFunctionSpecification is called when exiting the functionSpecification production.
	ExitFunctionSpecification(c *FunctionSpecificationContext)

	// ExitFunctionDeclaration is called when exiting the functionDeclaration production.
	ExitFunctionDeclaration(c *FunctionDeclarationContext)

	// ExitParameterDeclaration is called when exiting the parameterDeclaration production.
	ExitParameterDeclaration(c *ParameterDeclarationContext)

	// ExitReturnsClause is called when exiting the returnsClause production.
	ExitReturnsClause(c *ReturnsClauseContext)

	// ExitLanguageCharacteristic is called when exiting the languageCharacteristic production.
	ExitLanguageCharacteristic(c *LanguageCharacteristicContext)

	// ExitDeterministicCharacteristic is called when exiting the deterministicCharacteristic production.
	ExitDeterministicCharacteristic(c *DeterministicCharacteristicContext)

	// ExitReturnsNullOnNullInputCharacteristic is called when exiting the returnsNullOnNullInputCharacteristic production.
	ExitReturnsNullOnNullInputCharacteristic(c *ReturnsNullOnNullInputCharacteristicContext)

	// ExitCalledOnNullInputCharacteristic is called when exiting the calledOnNullInputCharacteristic production.
	ExitCalledOnNullInputCharacteristic(c *CalledOnNullInputCharacteristicContext)

	// ExitSecurityCharacteristic is called when exiting the securityCharacteristic production.
	ExitSecurityCharacteristic(c *SecurityCharacteristicContext)

	// ExitCommentCharacteristic is called when exiting the commentCharacteristic production.
	ExitCommentCharacteristic(c *CommentCharacteristicContext)

	// ExitReturnStatement is called when exiting the returnStatement production.
	ExitReturnStatement(c *ReturnStatementContext)

	// ExitAssignmentStatement is called when exiting the assignmentStatement production.
	ExitAssignmentStatement(c *AssignmentStatementContext)

	// ExitSimpleCaseStatement is called when exiting the simpleCaseStatement production.
	ExitSimpleCaseStatement(c *SimpleCaseStatementContext)

	// ExitSearchedCaseStatement is called when exiting the searchedCaseStatement production.
	ExitSearchedCaseStatement(c *SearchedCaseStatementContext)

	// ExitIfStatement is called when exiting the ifStatement production.
	ExitIfStatement(c *IfStatementContext)

	// ExitIterateStatement is called when exiting the iterateStatement production.
	ExitIterateStatement(c *IterateStatementContext)

	// ExitLeaveStatement is called when exiting the leaveStatement production.
	ExitLeaveStatement(c *LeaveStatementContext)

	// ExitCompoundStatement is called when exiting the compoundStatement production.
	ExitCompoundStatement(c *CompoundStatementContext)

	// ExitLoopStatement is called when exiting the loopStatement production.
	ExitLoopStatement(c *LoopStatementContext)

	// ExitWhileStatement is called when exiting the whileStatement production.
	ExitWhileStatement(c *WhileStatementContext)

	// ExitRepeatStatement is called when exiting the repeatStatement production.
	ExitRepeatStatement(c *RepeatStatementContext)

	// ExitCaseStatementWhenClause is called when exiting the caseStatementWhenClause production.
	ExitCaseStatementWhenClause(c *CaseStatementWhenClauseContext)

	// ExitElseIfClause is called when exiting the elseIfClause production.
	ExitElseIfClause(c *ElseIfClauseContext)

	// ExitElseClause is called when exiting the elseClause production.
	ExitElseClause(c *ElseClauseContext)

	// ExitVariableDeclaration is called when exiting the variableDeclaration production.
	ExitVariableDeclaration(c *VariableDeclarationContext)

	// ExitSqlStatementList is called when exiting the sqlStatementList production.
	ExitSqlStatementList(c *SqlStatementListContext)

	// ExitPrivilege is called when exiting the privilege production.
	ExitPrivilege(c *PrivilegeContext)

	// ExitQualifiedName is called when exiting the qualifiedName production.
	ExitQualifiedName(c *QualifiedNameContext)

	// ExitQueryPeriod is called when exiting the queryPeriod production.
	ExitQueryPeriod(c *QueryPeriodContext)

	// ExitRangeType is called when exiting the rangeType production.
	ExitRangeType(c *RangeTypeContext)

	// ExitSpecifiedPrincipal is called when exiting the specifiedPrincipal production.
	ExitSpecifiedPrincipal(c *SpecifiedPrincipalContext)

	// ExitCurrentUserGrantor is called when exiting the currentUserGrantor production.
	ExitCurrentUserGrantor(c *CurrentUserGrantorContext)

	// ExitCurrentRoleGrantor is called when exiting the currentRoleGrantor production.
	ExitCurrentRoleGrantor(c *CurrentRoleGrantorContext)

	// ExitUnspecifiedPrincipal is called when exiting the unspecifiedPrincipal production.
	ExitUnspecifiedPrincipal(c *UnspecifiedPrincipalContext)

	// ExitUserPrincipal is called when exiting the userPrincipal production.
	ExitUserPrincipal(c *UserPrincipalContext)

	// ExitRolePrincipal is called when exiting the rolePrincipal production.
	ExitRolePrincipal(c *RolePrincipalContext)

	// ExitRoles is called when exiting the roles production.
	ExitRoles(c *RolesContext)

	// ExitUnquotedIdentifier is called when exiting the unquotedIdentifier production.
	ExitUnquotedIdentifier(c *UnquotedIdentifierContext)

	// ExitQuotedIdentifier is called when exiting the quotedIdentifier production.
	ExitQuotedIdentifier(c *QuotedIdentifierContext)

	// ExitBackQuotedIdentifier is called when exiting the backQuotedIdentifier production.
	ExitBackQuotedIdentifier(c *BackQuotedIdentifierContext)

	// ExitDigitIdentifier is called when exiting the digitIdentifier production.
	ExitDigitIdentifier(c *DigitIdentifierContext)

	// ExitDecimalLiteral is called when exiting the decimalLiteral production.
	ExitDecimalLiteral(c *DecimalLiteralContext)

	// ExitDoubleLiteral is called when exiting the doubleLiteral production.
	ExitDoubleLiteral(c *DoubleLiteralContext)

	// ExitIntegerLiteral is called when exiting the integerLiteral production.
	ExitIntegerLiteral(c *IntegerLiteralContext)

	// ExitIdentifierUser is called when exiting the identifierUser production.
	ExitIdentifierUser(c *IdentifierUserContext)

	// ExitStringUser is called when exiting the stringUser production.
	ExitStringUser(c *StringUserContext)

	// ExitNonReserved is called when exiting the nonReserved production.
	ExitNonReserved(c *NonReservedContext)
}
