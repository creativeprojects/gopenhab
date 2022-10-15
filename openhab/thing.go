package openhab

import "github.com/creativeprojects/gopenhab/api"

type Thing struct {
	uid    string
	data   api.Thing
	client *Client
}

func newThing(client *Client, uid string) *Thing {
	return &Thing{
		uid:    uid,
		client: client,
	}
}
