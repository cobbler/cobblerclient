package client

type CobblerEvent struct {
	Statetime float64
	Name      string
	State     string
	ReadByWho []string
}

var EMPTYEVENT = CobblerEvent{
	Statetime: 0.0,
	Name:      "empty",
	State:     "invalid",
	ReadByWho: nil,
}

// GetEvents is ..
func (c *Client) GetEvents() error {
	_, err := c.Call("get_events", c.Token)
	return err
}

// GetEventLog is ...
func (c *Client) GetEventLog() error {
	_, err := c.Call("get_event_log", c.Token)
	return err
}

// GetTaskStatus takes the event ID from Cobbler and returns its the status.
func (c *Client) GetTaskStatus(eventId string) (CobblerEvent, error) {
	_, err := c.Call("get_task_status", eventId)
	if err != nil {
		return EMPTYEVENT, err
	} else {
		// FIXME: Server has the wrong format. Needs to be fixed there.
		// return result.(string), err
		return EMPTYEVENT, err
	}
}
