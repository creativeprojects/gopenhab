[![Go Reference](https://pkg.go.dev/badge/github.com/creativeprojects/gopenhab.svg)](https://pkg.go.dev/github.com/creativeprojects/gopenhab)
[![Build](https://github.com/creativeprojects/gopenhab/actions/workflows/build.yml/badge.svg)](https://github.com/creativeprojects/gopenhab/actions/workflows/build.yml)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=creativeprojects_gopenhab&metric=coverage)](https://sonarcloud.io/summary/new_code?id=creativeprojects_gopenhab)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=creativeprojects_gopenhab&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=creativeprojects_gopenhab)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=creativeprojects_gopenhab&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=creativeprojects_gopenhab)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=creativeprojects_gopenhab&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=creativeprojects_gopenhab)

# gopenHAB

Write your openHAB rules in Go. The power of openHAB rules with the simplicity of Go.

I had some existing code written in Go that I needed to connect to openHAB.
My first thought was to make a REST API to make this external system accessible from openHAB, but after fiddling with openHAB DSL rules and trying Jython scripts, I realized the best thing was to actually connect my system to the openHAB event bus and replicate a rule system, all in Go.

In theory, everything you can do in a DSL rule or in Jython should be available.

Notes:
- I do fully use it on my home automation platform: I did move all my DSL rules to gopenhab.
- This is still work in progress: you might need some events that are not yet implemented, or some features that are not yet available.

Here's an example of what you can do with it right now:

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

Imagine the `calculateAverage` function simply sends the result to the item in output configuration, from the example the name of the item is `ZoneTemperature`.

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

# TODO

- Handle all state types. Handled for now are `String`, `Switch`, `Number`, `DateTime`.
- Add triggers for more events. All `item` events have triggers, and some `thing` events (but not all)
- Ability to update rules
- Handle more events on the openhab test server (typically, `things` are not supported yet)

# Limitations of the mock openHAB server for testing

- When the server receives a `command` event, it doesn't follow with any corresponding `state` event. You need to publish the `state` events manually if you need them in your tests.

# Compatibility

The library supports API version 3 to 6.

At some point in time, this library was tested (and running in production) with these versions of openHAB:
- 2.5
- 3.0
- 3.1
- 3.2
- 3.3
- 3.4
- 4.0
- 4.1

# Integration tests

I do have some integrations tests running against a real openHAB server (currently 4.1.2). The server is running an exact copy of my home configuration, with a script sending mock states to the server via MQTT.

For reliability testing, I also have a Toxiproxy between openHAB and gopenhab.

I might publish some of these tests at some point.
