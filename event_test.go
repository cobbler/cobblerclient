package cobblerclient

import (
	"reflect"
	"testing"
)

func TestGetTaskStatus(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-task-status")
	expectedResult := CobblerEvent{
		ID:        "2000-01-01_000000_(CLI) ACL Configuration_00000000000000000000000000000007",
		StateTime: 0.0,
		Name:      "(CLI) ACL Configuration",
		State:     "failed",
		ReadByWho: []string{},
	}

	result, err := c.GetTaskStatus("2000-01-01_000000_(CLI) ACL Configuration_00000000000000000000000000000007")
	FailOnError(t, err)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Result from 'get_task_status' did not match expected result.")
	}
}

func TestGetEvents(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-events")

	_, err := c.GetEvents("")
	FailOnError(t, err)
}

func TestGetEventLog(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-event-log")

	_, err := c.GetEventLog("2000-01-01_000000_(CLI) ACL Configuration_00000000000000000000000000000007")
	FailOnError(t, err)
}
