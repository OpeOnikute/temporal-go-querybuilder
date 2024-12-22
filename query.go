package temporalgoquerybuilder

import (
	"fmt"
	"strings"
	"time"
)

const (
	// Query delimiters
	DelimOpenParentheses  = '('
	DelimCloseParentheses = ')'

	// Logical Operators
	// STARTS_WITH
	LogicalOpAND        = "AND"
	LogicalOpOR         = "OR"
	LogicalOpBetween    = "BETWEEN"
	LogicalOpIn         = "IN"
	LogicalOpStartsWith = "STARTS_WITH"

	// Default Search Attributes - used in queries
	SearchAttrBatcherUser                = "BatcherUser"
	SearchAttrBinaryChecksums            = "BinaryChecksums"
	SearchAttrBuildIds                   = "BuildIds"
	SearchAttrCloseTime                  = "CloseTime"
	SearchAttrExecutionDuration          = "ExecutionDuration"
	SearchAttrExecutionStatus            = "ExecutionStatus"
	SearchAttrExecutionTime              = "ExecutionTime"
	SearchAttrHistoryLength              = "HistoryLength"
	SearchAttrHistorySizeBytes           = "HistorySizeBytes"
	SearchAttrRunId                      = "RunId"
	SearchAttrStartTime                  = "StartTime"
	SearchAttrStateTransitionCount       = "StateTransitionCount"
	SearchAttrTaskQueue                  = "TaskQueue"
	SearchAttrTemporalChangeVersion      = "TemporalChangeVersion"
	SearchAttrTemporalScheduledStartTime = "TemporalScheduledStartTime"
	SearchAttrTemporalScheduledById      = "TemporalScheduledById"
	SearchAttrTemporalSchedulePaused     = "TemporalSchedulePaused"
	SearchAttrWorkflowId                 = "WorkflowId"
	SearchAttrWorkflowType               = "WorkflowType"

	// Execution statuses
	ExecStatusRunning        = "Running"
	ExecStatusCompleted      = "Completed"
	ExecStatusFailed         = "Failed"
	ExecStatusCanceled       = "Canceled"
	ExecStatusTerminated     = "Terminated"
	ExecStatusContinuedAsNew = "ContinuedAsNew"
	ExecStatusTimedOut       = "TimedOut"
)

// QueryBuilder allows us to encode Temporal queries using a
// clearly-defined API that can be tested and avoids flaky string formatting.
type QueryBuilder map[string]string

// StartQuery is used to add the first entry to a query. The only reason it exists
// is to remove the need for a preceding logical operator (e.g. AND), which is not
// needed when adding the first entry
func (t QueryBuilder) StartQuery(searchAttr string, comparisonOp byte, value string) {
	// No need for a logical operator when starting a query.
	var logicalOp string
	t.Query(searchAttr, comparisonOp, value, logicalOp)
}

// Or is used to add a query preceded with an OR logical operator.
// If there is no preceding query, the `OR` won't be added, but it's better
// for readability to always start with StartQuery.
func (t QueryBuilder) Or(searchAttr string, comparisonOp byte, value string) {
	t.Query(searchAttr, comparisonOp, value, LogicalOpOR)
}

// And is used to add a query preceded with an AND logical operator.
// If there is no preceding query, the `OR` won't be added, but it's better
// for readability to always start with StartQuery.
func (t QueryBuilder) And(searchAttr string, comparisonOp byte, value string) {
	t.Query(searchAttr, comparisonOp, value, LogicalOpAND)
}

// Query is used to build a query entry.
func (t QueryBuilder) Query(searchAttr string, comparisonOp byte, value string, logicalOp string) {
	var buf strings.Builder
	if len(t) > 0 {
		buf.WriteString(fmt.Sprintf(" %s ", logicalOp))
	}
	buf.WriteString(searchAttr)
	buf.WriteByte(comparisonOp)
	buf.WriteString(fmt.Sprintf("'%s'", value))
	t[searchAttr] = buf.String()
}

// Between is used to create a Between ... AND query. It will add start and closing parentheses
// e.g. (StartTime BETWEEN '2024-12-16T20:47:35Z' AND '2024-12-16T20:52:35Z')
func (t QueryBuilder) Between(searchAttr string, start, end time.Time, logicalOp string) {
	var buf strings.Builder
	if len(t) > 0 {
		buf.WriteString(fmt.Sprintf(" %s ", logicalOp))
	}
	buf.WriteByte(DelimOpenParentheses)
	buf.WriteString(searchAttr)
	buf.WriteString(fmt.Sprintf(" %s ", LogicalOpBetween))
	buf.WriteString(fmt.Sprintf("'%s'", start.Format(time.RFC3339)))
	buf.WriteString(fmt.Sprintf(" %s ", LogicalOpAND))
	buf.WriteString(fmt.Sprintf("'%s'", end.Format(time.RFC3339)))
	buf.WriteByte(DelimCloseParentheses)
	t[searchAttr] = buf.String()
}

// In is used to create a 'SearchAttr in ('foo', 'bar')' type query.
func (t QueryBuilder) In(searchAttr string, values []string, logicalOp string) {
	var buf strings.Builder
	if len(t) > 0 {
		buf.WriteString(fmt.Sprintf(" %s ", logicalOp))
	}

	// Add single quotes to every entry
	vs := []string{}
	for _, v := range values {
		vs = append(vs, fmt.Sprintf("'%s'", v))
	}

	buf.WriteString(searchAttr)
	buf.WriteString(fmt.Sprintf(" %s ", LogicalOpIn))
	buf.WriteByte(DelimOpenParentheses)
	buf.WriteString(strings.Join(vs, ", "))
	buf.WriteByte(DelimCloseParentheses)
	t[searchAttr] = buf.String()
}

// StartsWith is used to create a 'STARTS_WITH foo' type query.
func (t QueryBuilder) StartsWith(value string, logicalOp string) {
	var buf strings.Builder
	if len(t) > 0 {
		buf.WriteString(fmt.Sprintf(" %s ", logicalOp))
	}

	buf.WriteString(fmt.Sprintf("%s ", LogicalOpStartsWith))
	buf.WriteString(fmt.Sprintf("'%s'", value))
	t[value] = buf.String()
}

// Encode will provide the final query by combining all current queries
func (t QueryBuilder) Encode() string {
	values := make([]string, 0, len(t))
	for k := range t {
		values = append(values, t[k])
	}
	return strings.Join(values, "")
}
