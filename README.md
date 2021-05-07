# gopenHAB

Write your openHAB rules in Go.

This is very much work in progress, but you can this working:

```go
package main

import (
	"log"
	"time"

	"github.com/creativeprojects/gopenhab/openhab"
)

func main() {
	openhab.SetLogger(log.Default())
	client := openhab.NewClient(openhab.Config{
		URL: "http://localhost:8080",
	})

	client.AddRule(
		openhab.RuleData{Name: "Connected to openHAB events"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Printf("EVENT: client connected")
		},
		openhab.OnConnect(),
	)

	client.AddRule(
		openhab.RuleData{Name: "Disconnected from openHAB events"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Print("EVENT: client disconnected")
		},
		openhab.OnDisconnect(),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving item command"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			if ev, ok := e.(*event.ItemReceivedCommand); ok {
				log.Printf("EVENT: Back_Garden_Lighting_Switch received command %+v", ev.Command)
			}
		},
		openhab.OnItemReceivedCommand("Back_Garden_Lighting_Switch", nil),
	)

	client.AddRule(
		openhab.RuleData{
			Name: "Test rule",
		},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			item, err := client.Items().GetItem("Back_Garden_Lighting_Switch")
			if err != nil {
				log.Print(err)
			}
			err = item.SendCommand(openhab.SwitchON)
			if err != nil {
				log.Printf("sending command: %s", err)
			}

			state, err := item.State()
			if err != nil {
				log.Printf("getting state: %s", err)
			}

			log.Printf("switched %s for 4 seconds...", state)
			time.Sleep(4 * time.Second)

			err = item.SendCommand(openhab.SwitchOFF)
			if err != nil {
				log.Printf("sending command: %s", err)
			}

			state, err = item.State()
			if err != nil {
				log.Printf("getting state: %s", err)
			}

			log.Printf("switched %s", state)
		},
		openhab.OnTimeCron("*/10 * * ? * *"),
	)
	client.Start()
}

```
