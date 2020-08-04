# go-broadcaster

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>

Broadcaster-Subscriber via channels in Go for easy communication within your goroutines

## Getting started

```Go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Ganitzsh/go-broadcaster"
)

type fooPayloadOne struct {
	Value string
}

type fooPayloadTwo struct {
	Key   string
	Value int
}

type barPayload struct {
	ID      string
	Message string
}

var errIncorrectPayloadType = errors.New("incorrect payload type")

func fooHandlerOne(b *broadcaster.Broadcast) error {
	payload, ok := b.Data.(*fooPayloadOne)
	if !ok {
		return errIncorrectPayloadType
	}

	fmt.Printf("fooPayloadOne payload value is: %s\n", payload.Value)
	return nil
}

func fooHandlerTwo(b *broadcaster.Broadcast) error {
	payload, ok := b.Data.(*fooPayloadTwo)
	if !ok {
		return errIncorrectPayloadType
	}

	fmt.Printf("fooPayloadTwo payload value is: %s = %d\n", payload.Key, payload.Value)
	return nil
}

func barHandler(b *broadcaster.Broadcast) error {
	payload, ok := b.Data.(*barPayload)
	if !ok {
		return errIncorrectPayloadType
	}

	fmt.Printf("barPayload payload value is: id=%s message=%s\n", payload.ID, payload.Message)
	return nil
}

func logSubscriberHandlerError(sub *broadcaster.Subscriber) {
	for {
		err := <-sub.ErrChan
		if err == nil {
			return
		}
		if err != errIncorrectPayloadType {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}

func main() {
  // Set the amount of workers per subscribers
	workers := uint32(8)
	agent := broadcaster.NewBroadcastAgent()
	defer func() {
		fmt.Println("Closing agent...")
		now := time.Now()
		wg := agent.Close()
    // Close returns a waitgroup that you can use to wait for all the routines to
    // be stopped before exiting the program
		wg.Wait()
		fmt.Printf("Closing the agent took %v\n", time.Since(now))
	}()

	fooFreq := "foo"
  // Subscribe to the foo frequency and handle multiple jobs when receiving
  // broadcasts on this frequency
	subFoo1 := agent.Subscribe(fooFreq, fooHandlerOne, workers)
	subFoo2 := agent.Subscribe(fooFreq, fooHandlerTwo, workers)

	barFreq := "bar"
  // Same with bar..
	subBar := agent.Subscribe(barFreq, barHandler, workers)

  // Start you subscribers to let them handle the work
	subFoo1.Start()
	subFoo2.Start()
	subBar.Start()

  // An example of how you could handle errors coming from the handlers
	go logSubscriberHandlerError(subFoo1)
	go logSubscriberHandlerError(subFoo2)
	go logSubscriberHandlerError(subBar)

	now := time.Now()
  // Broadcasts a lot of messages
	for i := 0; i < 100; i++ {
		agent.Broadcast(fooFreq, &fooPayloadOne{Value: "payload for handler 1"})
		agent.Broadcast(fooFreq, &fooPayloadTwo{Key: "payload", Value: 42})
		agent.Broadcast(barFreq, &barPayload{ID: "payload1", Message: "Hello from barHandler"})
	}

  // Wait for all the workers to be done handling the jobs
  // note that it does not stop the routines see defer at the top
	agent.WaitForCompletion()
	fmt.Printf("The whole thing took %v to execute\n", time.Since(now))
}

```
