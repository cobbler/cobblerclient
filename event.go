package cobblerclient

type CobblerEvent struct {
	ID        string
	StateTime float64
	Name      string
	State     string
	ReadByWho []string
}

var EMPTYEVENT = CobblerEvent{
	ID:        "",
	StateTime: 0.0,
	Name:      "empty",
	State:     "invalid",
	ReadByWho: nil,
}

func unmarshalEvent(eventId string, data interface{}) *CobblerEvent {
	eventData := data.([]interface{})
	readByWho := make([]string, 0)
	for _, user := range eventData[3].([]interface{}) {
		readByWho = append(readByWho, user.(string))
	}
	return &CobblerEvent{
		ID:        eventId,
		StateTime: eventData[0].(float64),
		Name:      eventData[1].(string),
		State:     eventData[2].(string),
		ReadByWho: readByWho,
	}
}

// GetEvents retrieves all events from the Cobbler server
func (c *Client) GetEvents(forUser string) ([]*CobblerEvent, error) {
	var events []*CobblerEvent
	unmarshalledResult, err := c.Call("get_events", forUser)
	if err != nil {
		return nil, err
	}
	for key, event := range unmarshalledResult.(map[string]interface{}) {
		eventObj := unmarshalEvent(key, event)
		events = append(events, eventObj)
	}
	return events, err
}

// GetEventLog retrieves the logged messages for a given event id.
func (c *Client) GetEventLog(eventId string) (string, error) {
	res, err := c.Call("get_event_log", eventId)
	return returnString(res, err)
}

// GetTaskStatus takes the event ID from Cobbler and returns its status.
func (c *Client) GetTaskStatus(eventId string) (CobblerEvent, error) {
	unmarshalledResult, err := c.Call("get_task_status", eventId)
	if err != nil {
		return EMPTYEVENT, err
	}
	eventObj := unmarshalEvent(eventId, unmarshalledResult)
	return *eventObj, err
}
