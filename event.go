package cobblerclient

import (
	"fmt"
)

type CobblerEvent struct {
	id        string
	statetime float64
	name      string
	state     string
	readByWho []string
}

var EMPTYEVENT = CobblerEvent{
	id:        "",
	statetime: 0.0,
	name:      "empty",
	state:     "invalid",
	readByWho: nil,
}

// GetEvents retrieves all events from the Cobbler server
func (c *Client) GetEvents(forUser string) ([]*CobblerEvent, error) {
	var events []*CobblerEvent
	unmarshalledResult, err := c.Call("get_events", forUser)
	if err != nil {
		return nil, err
	}
	for key, event := range unmarshalledResult.(map[string]interface{}) {
		eventData := event.([]interface{})
		eventObj := &CobblerEvent{
			id:        key,
			statetime: eventData[0].(float64),
			name:      eventData[1].(string),
			state:     eventData[2].(string),
			readByWho: nil, // eventData[3].([]string)
		}
		// TODO: Add readByWho
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
	// FIXME: Server has the wrong format. Needs to be fixed there.
	// return result.(string), err
	fmt.Printf("%#v", unmarshalledResult)
	return EMPTYEVENT, err
}
