package broadcaster_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Ganitzsh/go-broadcaster"
)

var handled bool

func handler(broadcast *broadcaster.Broadcast) error {
	handled = true
	return nil
}

func handlerWithError(broadcast *broadcaster.Broadcast) error {
	return errors.New("error")
}

func clean(agent *broadcaster.BroadcastAgent) {
	agent.WaitForCompletion()
	wg := agent.Close()
	wg.Wait()
}

func TestAgentSubscriberNoError(t *testing.T) {
	handled = false
	freq := "event1"
	agent := broadcaster.NewBroadcastAgent()
	defer clean(agent)

	subscriber := agent.Subscribe(freq, handler, 1)
	subscriber.Start()
	agent.Broadcast(freq, nil)

	select {
	case err := <-subscriber.ErrChan:
		if err != nil {
			t.Fatalf("no error expected, got %v", err)
		}
	case <-time.After(1 * time.Millisecond):
		break
	}

	if !handled {
		t.Fatal("handled switch not flipped")
	}
}

func TestAgentSubscriberWithError(t *testing.T) {
	freq := "event1"
	agent := broadcaster.NewBroadcastAgent()
	defer clean(agent)

	subscriber := agent.Subscribe(freq, handlerWithError, 1)
	subscriber.Start()
	agent.Broadcast(freq, nil)

	select {
	case <-subscriber.ErrChan:
		break
	case <-time.After(1 * time.Millisecond):
		t.Fatalf("expected an error")
	}
}

func TestAgentSubscriberZero(t *testing.T) {
	freq := "event1"
	agent := broadcaster.NewBroadcastAgent()
	defer clean(agent)

	subscriber := agent.Subscribe(freq, handlerWithError, 0)
	subscriber.Start()
	agent.Broadcast(freq, nil)

	<-subscriber.ErrChan
}
