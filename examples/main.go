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
		openhab.Debounce(1*time.Minute, openhab.OnConnect()),
	)

	client.AddRule(
		openhab.RuleData{Name: "Disconnected from openHAB events"},
		func(client *openhab.Client, ruleData openhab.RuleData, e event.Event) {
			log.Print("SYSTEM EVENT: client disconnected")
		},
		openhab.Debounce(10*time.Second, openhab.OnDisconnect()),
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
			item, err := client.GetItem("Back_Garden_Lighting_Switch")
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
