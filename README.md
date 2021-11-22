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
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			log.Printf("SYSTEM EVENT: client connected")
		},
		openhab.Debounce(openhab.OnConnect(), 1*time.Minute),
	)

	client.AddRule(
		openhab.RuleData{Name: "Disconnected from openHAB events"},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			log.Print("SYSTEM EVENT: client disconnected")
		},
		openhab.Debounce(openhab.OnDisconnect(), 10*time.Second),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving item command"},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			if ev, ok := e.(event.ItemReceivedCommand); ok {
				log.Printf("COMMAND EVENT: Back_Garden_Lighting_Switch received command %+v", ev.Command)
			}
		},
		openhab.OnItemReceivedCommand("Back_Garden_Lighting_Switch", nil),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving ON command"},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			log.Print("COMMAND EVENT: Back_Garden_Lighting_Switch switched ON")
		},
		openhab.OnItemReceivedCommand("Back_Garden_Lighting_Switch", openhab.SwitchON),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving OFF command"},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			log.Print("COMMAND EVENT: Back_Garden_Lighting_Switch switched OFF")
		},
		openhab.OnItemReceivedCommand("Back_Garden_Lighting_Switch", openhab.SwitchOFF),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving updated state"},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			if ev, ok := e.(event.ItemReceivedState); ok {
				log.Printf("STATE EVENT: Back_Garden_Lighting_Switch received state %+v", ev.State)
			}
		},
		openhab.OnItemReceivedState("Back_Garden_Lighting_Switch", nil),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving ON state"},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			log.Printf("STATE EVENT: Back_Garden_Lighting_Switch state is now ON")
		},
		openhab.OnItemReceivedState("Back_Garden_Lighting_Switch", openhab.SwitchON),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving OFF state"},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			log.Printf("STATE EVENT: Back_Garden_Lighting_Switch state is now OFF")
		},
		openhab.OnItemReceivedState("Back_Garden_Lighting_Switch", openhab.SwitchOFF),
	)

	client.AddRule(
		openhab.RuleData{Name: "Receiving state changed"},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
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
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
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


```

# Unit test your rules

To be able to run some unit tests I created a *mock* openHAB server, which can trigger events and can keep items in memory. This is work in progress but you can use it to test your rules.

## How to test a simple event

Imagine you have a function `calculateZoneTemperature` that takes an array of values coming from the `Context` and sends an average temperature to an output item. The context of the function will be as such:

```go

openhab.RuleData{
	Name: "Calculate average",
	Context: zoneContext{
		name:   "test-zone",
		config: ZoneConfiguration{Output: "ZoneTemperature", Sensors: []string{"temperature1", "temperature2"}},
	},
},
```

I'm not giving the code of `calculateAverage`, but simply imagine it sends the result to the item in output configuration, from the example the name of the item is `ZoneTemperature`.

Here's how you can test it with the openHAB mock server:

```go
func TestCalculateZoneTemperature(t *testing.T) {
	const temperatureItem1 = "temperature1"
	const temperatureItem2 = "temperature2"
	const averageTemperatureItem = "ZoneTemperature"

	// Create the openHAB mock server that will publish events from changes coming from the API calls
	server := openhabtest.NewServer(openhabtest.Config{Log: t, SendEventsFromAPI: true})
	defer server.Close()

	// setup all 3 items in the mock server
	server.SetItem(api.Item{
		Name:  averageTemperatureItem,
		Type:  "Number",
		State: "0.0",
	})
	server.SetItem(api.Item{
		Name:  temperatureItem1,
		Type:  "Number",
		State: "10.0",
	})
	server.SetItem(api.Item{
		Name:  temperatureItem2,
		Type:  "Number",
		State: "10.0",
	})

	// create a client that connects to our mock server
	client := openhab.NewClient(openhab.Config{URL: server.URL()})

	// standard rule to calculate the average
	client.AddRule(
		openhab.RuleData{
			Name: "Calculate average",
			Context: zoneContext{
				name:   "test-zone",
				config: ZoneConfiguration{Output: averageTemperatureItem, Sensors: []string{temperatureItem1, temperatureItem2}},
			},
		},
		calculateZoneTemperature, // this is the code to test
		openhab.OnItemReceivedState(temperatureItem1, nil),
		openhab.OnItemReceivedState(temperatureItem2, nil),
	)

	// testing rule to verify the calculation
	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		openhab.RuleData{
			Name: "Test rule",
		},
		func(client openhab.RuleClient, ruleData openhab.RuleData, e event.Event) {
			// test is finished after receiving this event
			defer wg.Done()

			ev, ok := e.(event.ItemReceivedCommand)
			if !ok {
				t.Fatalf("expected event to be of type ItemReceivedCommand")
			}
			assert.Equal(t, "10.5", ev.Command)
		},
		openhab.OnItemReceivedCommand(averageTemperatureItem, nil),
	)

	// start the client in the background so we can send events to it
	go func() {
		client.Start()
	}()

	// make sure the client is ready (it surely needs less than that)
	time.Sleep(10 * time.Millisecond)
	// we simulate openhab receiving an event: temperature2 item received a new value of 11 degrees, brrrr!
	server.Event(event.NewItemReceivedState(temperatureItem2, "Number", "11.0"))

	wg.Wait()
	client.Stop()
}
```
