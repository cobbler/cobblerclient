package cobblerclient

import (
	"testing"

	"github.com/ContainerSolutions/go-utils"
)

func TestGetTaskStatus(t *testing.T) {
	c := createStubHTTPClient(t, "get-task-status-req.xml", "get-task-status-res.xml")

	res, err := c.GetTaskStatus("2022-09-30_200403_Updating Signatures_8f2b3c1626fb4b158636059b31242ee6")
	utils.FailOnError(t, err)
	if res.name == "" {
		t.Fatalf("Expected non emtpy string")
	}
}

func TestGetEvents(t *testing.T) {
	c := createStubHTTPClient(t, "get-events-req.xml", "get-events-res.xml")

	_, err := c.GetEvents("")
	utils.FailOnError(t, err)
}

func TestGetEventLog(t *testing.T) {
	c := createStubHTTPClient(t, "get-event-log-req.xml", "get-event-log-res.xml")

	_, err := c.GetEventLog("2022-09-30_145124_Sync_2cabdc4eddfa4731b45f145d7b625e29")
	utils.FailOnError(t, err)
}
