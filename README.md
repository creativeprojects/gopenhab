[![Build](https://github.com/creativeprojects/gopenhab/actions/workflows/build.yml/badge.svg)](https://github.com/creativeprojects/gopenhab/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/creativeprojects/gopenhab/branch/main/graph/badge.svg?token=wyzDjPzIO3)](https://codecov.io/gh/creativeprojects/gopenhab)
[![Go Reference](https://pkg.go.dev/badge/github.com/creativeprojects/gopenhab.svg)](https://pkg.go.dev/github.com/creativeprojects/gopenhab)

# gopenHAB

Write your openHAB rules in Go. The power of openHAB rules with the simplicity of Go.

I had some existing code written in Go that I needed to connect to openHAB.
My first thought was to make a REST API to make this external system accessible from openHAB, but after fiddling with openHAB DSL rules and trying Jython scripts, I realized the best thing was to actually connect my system to the openHAB event bus and replicate a rule system, all in Go.

In theory, everything you can do in a DSL rule or in Jython should be available.

This is work in progress, but here's an example of what you can already do with it:

```go
package main

import (
	"log"
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/creativeprojects/gopenhab/openhab"
)

func main() {
	openhab.SetDebugLog(log.Default())
	client := openhab.NewClient(openhab.Config{
		URL: "http://localhost:8080",
	})

	client.AddRule(
		openhab.RuleData{Name: "Connected to openHAB events"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Printf("SYSTEM EVENT: client connected")
		},
		openhab.Debounce(openhab.OnConnect(), 1*time.Minute),
	)

	client.AddRule(
		openhab.RuleData{Name: "Disconnected from openHAB events"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Print("SYSTEM EVENT: client disconnected")
		},
		openhab.Debounce(openhab.OnDisconnect(), 10*time.Second),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving item command"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			if ev, ok := e.(event.ItemReceivedCommand); ok {
				log.Printf("COMMAND EVENT: Back_Garden_Lighting_Switch received command %+v", ev.Command)
			}
		},
		openhab.OnItemReceivedCommand("Back_Garden_Lighting_Switch", nil),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving ON command"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Print("COMMAND EVENT: Back_Garden_Lighting_Switch switched ON")
		},
		openhab.OnItemReceivedCommand("Back_Garden_Lighting_Switch", openhab.SwitchON),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving OFF command"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Print("COMMAND EVENT: Back_Garden_Lighting_Switch switched OFF")
		},
		openhab.OnItemReceivedCommand("Back_Garden_Lighting_Switch", openhab.SwitchOFF),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving updated state"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			if ev, ok := e.(event.ItemReceivedState); ok {
				log.Printf("STATE EVENT: Back_Garden_Lighting_Switch received state %+v", ev.State)
			}
		},
		openhab.OnItemReceivedState("Back_Garden_Lighting_Switch", nil),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving ON state"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Printf("STATE EVENT: Back_Garden_Lighting_Switch state is now ON")
		},
		openhab.OnItemReceivedState("Back_Garden_Lighting_Switch", openhab.SwitchON),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving OFF state"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Printf("STATE EVENT: Back_Garden_Lighting_Switch state is now OFF")
		},
		openhab.OnItemReceivedState("Back_Garden_Lighting_Switch", openhab.SwitchOFF),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving state changed"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			if ev, ok := e.(event.ItemStateChanged); ok {
				log.Printf("STATE CHANGED EVENT: Back_Garden_Lighting_Switch changed to state %+v", ev.NewState)
			}
		},
		openhab.OnItemStateChanged("Back_Garden_Lighting_Switch"),
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

			_, err = item.SendCommandWait(openhab.SwitchON, 2*time.Second)
			if err != nil {
				log.Printf("sending command: %s", err)
			}
			time.Sleep(4 * time.Second)

			_, err = item.SendCommandWait(openhab.SwitchOFF, 2*time.Second)
			if err != nil {
				log.Printf("sending command: %s", err)
			}
		},
		openhab.OnTimeCron("*/10 * * ? * *"),
	)
	client.Start()
}


```

# testing

To be able to run some unit tests I created a *mock* openHAB server, which can trigger events and can keep items in memory. This is work in progress but you can use it to test your rules.

I'll post more information about it, but here's the [godoc](https://pkg.go.dev/github.com/creativeprojects/gopenhab/openhabtest).
