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
	item, err := client.Items().GetItem("Back_Garden_Lighting_Switch")
	if err != nil {
		log.Fatal(err)
	}

	client.AddRule(
		openhab.RuleConfig{
			Name: "Test rule",
		},
		func() {
			err := item.SendCommand(openhab.SwitchON)
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
		openhab.NewCron("*/10 * * ? * *"),
	)
	client.Start()
}

```
