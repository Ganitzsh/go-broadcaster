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
	workers := uint32(8)
	agent := broadcaster.NewBroadcastAgent()
	defer func() {
		fmt.Println("Closing agent...")
		now := time.Now()
		wg := agent.Close()
		wg.Wait()
		fmt.Printf("Closing the agent took %v\n", time.Since(now))
	}()

	fooFreq := "foo"
	subFoo1 := agent.Subscribe(fooFreq, fooHandlerOne, workers)
	subFoo2 := agent.Subscribe(fooFreq, fooHandlerTwo, workers)

	barFreq := "bar"
	subBar := agent.Subscribe(barFreq, barHandler, workers)

	subFoo1.Start()
	subFoo2.Start()
	subBar.Start()

	go logSubscriberHandlerError(subFoo1)
	go logSubscriberHandlerError(subFoo2)
	go logSubscriberHandlerError(subBar)

	now := time.Now()
	for i := 0; i < 100; i++ {
		agent.Broadcast(fooFreq, &fooPayloadOne{Value: "payload for handler 1"})
		agent.Broadcast(fooFreq, &fooPayloadTwo{Key: "payload", Value: 42})
		agent.Broadcast(barFreq, &barPayload{ID: "payload1", Message: "Hello from barHandler"})
	}

	agent.WaitForCompletion()
	fmt.Printf("The whole thing took %v to execute\n", time.Since(now))
}
