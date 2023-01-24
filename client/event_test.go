package client_test

import (
	"testing"

	"github.com/ContainerSolutions/go-utils"

	cobbler_testing "github.com/cobbler/cobblerclient/internal/testing"
)

func TestGetTaskStatus(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-task-status-req.xml", "get-task-status-res.xml")

	res, err := c.GetTaskStatus("2022-09-30_200403_Updating Signatures_8f2b3c1626fb4b158636059b31242ee6")
	utils.FailOnError(t, err)
	if res.Name == "" {
		t.Fatalf("Expected non emtpy string")
	}
}

func TestGetEvents(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-events-req.xml", "get-events-res.xml")

	err := c.GetEvents()
	utils.FailOnError(t, err)
}

func TestGetEventLog(t *testing.T) {
	c := cobbler_testing.CreateStubHTTPClient(t, "get-event-log-req.xml", "get-event-log-res.xml")

	err := c.GetEventLog()
	utils.FailOnError(t, err)
}
