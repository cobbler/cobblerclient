package cobblerclient

import (
	"reflect"
	"testing"
)

func TestGetTaskStatus(t *testing.T) {
	c := createStubHTTPClientSingle(t, "get-task-status")
	expectedResult := CobblerEvent{
		id:        "2022-09-30_200403_Updating Signatures_8f2b3c1626fb4b158636059b31242ee6",
		statetime: 1664568243.5196018,
		name:      "Updating Signatures",
		state:     "complete",
		readByWho: []string{},
	}

	result, err := c.GetTaskStatus("2022-09-30_200403_Updating Signatures_8f2b3c1626fb4b158636059b31242ee6")
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

	_, err := c.GetEventLog("2022-09-30_145124_Sync_2cabdc4eddfa4731b45f145d7b625e29")
	FailOnError(t, err)
}
