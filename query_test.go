package temporalgoquerybuilder

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testWorkflowName = "TestMe"
var testWorkflowNameBackup = "TestMeBackup"
var testStatus = "Running"

func TestTemporalQueryWithAnd(t *testing.T) {
	q := QueryBuilder{}
	q.Query(SearchAttrWorkflowType, '=', testWorkflowName, LogicalOpAND)
	q.And(SearchAttrExecutionStatus, '=', testStatus)
	// e.g. "WorkflowType='TestMe' AND ExecutionStatus='Running'"
	var expected string = fmt.Sprintf("%s = '%s' %s %s = '%s'", SearchAttrWorkflowType, testWorkflowName, LogicalOpAND, SearchAttrExecutionStatus, testStatus)
	assert.Equal(t, expected, q.Encode())
}

func TestTemporalQueryWithOr(t *testing.T) {
	q := QueryBuilder{}
	q.Query(SearchAttrWorkflowType, '=', testWorkflowName, LogicalOpAND)
	q.Or(SearchAttrExecutionStatus, '=', testStatus)
	// e.g. "WorkflowType='TestMe' OR ExecutionStatus='Running'"
	var expected string = fmt.Sprintf("%s = '%s' %s %s = '%s'", SearchAttrWorkflowType, testWorkflowName, LogicalOpOR, SearchAttrExecutionStatus, testStatus)
	assert.Equal(t, expected, q.Encode())
}

func TestTemporalQueryBetween(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Minute)
	q := QueryBuilder{}
	q.Between(SearchAttrStartTime, start, end, LogicalOpAND)
	// e.g. (StartTime BETWEEN '2024-12-16T20:47:35Z' AND '2024-12-16T20:52:35Z')
	var expected string = fmt.Sprintf("(%s %s '%s' %s '%s')", SearchAttrStartTime, LogicalOpBetween, start.Format(time.RFC3339), LogicalOpAND, end.Format(time.RFC3339))
	assert.Equal(t, expected, q.Encode())
}

func TestTemporalQueryIn(t *testing.T) {
	q := QueryBuilder{}
	q.In(SearchAttrWorkflowType, []string{testWorkflowName, testWorkflowNameBackup}, LogicalOpAND)
	// e.g. WorkflowType IN ('TestMe', 'TestMeBackup')
	var expected string = fmt.Sprintf("%s %s ('%s', '%s')", SearchAttrWorkflowType, LogicalOpIn, testWorkflowName, testWorkflowNameBackup)
	assert.Equal(t, expected, q.Encode())
}

func TestTemporalQueryStartsWith(t *testing.T) {
	testQ := "foo"
	q := QueryBuilder{}
	q.StartsWith(testQ, "")
	// e.g. STARTS_WITH foo
	var expected string = fmt.Sprintf("%s '%s'", LogicalOpStartsWith, testQ)
	assert.Equal(t, expected, q.Encode())
}
