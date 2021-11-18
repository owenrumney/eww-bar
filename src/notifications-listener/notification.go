package main

type Notification struct {
	AppName    string
	ReplacesID uint32

	AppIcon string
	Summary string
	Body    string
}

func decode(message []interface{}) (notification Notification) {

	return Notification{
		AppIcon: message[0].(string),
		Summary: message[3].(string),
		Body:    message[4].(string),
	}

}
