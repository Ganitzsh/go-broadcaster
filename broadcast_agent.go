package broadcaster

import (
	"sync"
)

// BroadcastAgent will register subscribers on different frequencies
// and allow broadcasting to share a message across multiple
// channels
type BroadcastAgent struct {
	mutex        sync.RWMutex
	subsRegister map[string][]*Subscriber
	wg           sync.WaitGroup
	wgBusy       sync.WaitGroup
}

// NewBroadcastAgent will return a new BroadcastAgent instance
func NewBroadcastAgent() *BroadcastAgent {
	return &BroadcastAgent{
		subsRegister: map[string][]*Subscriber{},
		wg:           sync.WaitGroup{},
		wgBusy:       sync.WaitGroup{},
	}
}

// Broadcast sends a message across all the subscribers on a given frequency
func (agent *BroadcastAgent) Broadcast(freq string, data interface{}) {
	broadcast := newBroadcast(freq, data)
	agent.mutex.RLock()

	for _, sub := range agent.subsRegister[freq] {
		sub.Comm <- broadcast
	}

	agent.mutex.RUnlock()
}

// Subscribe will register a new subscriber to the broadcast agent and will return
// a new pointer to a Subscriber
func (agent *BroadcastAgent) Subscribe(freq string, handler SubscriberHandler, workersAmmount uint32) *Subscriber {
	ret := newSubscriber(freq, handler, workersAmmount, &agent.wgBusy, &agent.wg)

	agent.mutex.Lock()
	if agent.subsRegister[freq] == nil {
		agent.subsRegister[freq] = []*Subscriber{}
	}

	agent.subsRegister[freq] = append(agent.subsRegister[freq], ret)
	agent.mutex.Unlock()

	return ret
}

// Close will stop all the subscribers and return a waitgroup that will
// be done when all the workers are done and closed
func (agent *BroadcastAgent) Close() *sync.WaitGroup {
	for _, subs := range agent.subsRegister {
		for _, sub := range subs {
			sub.stop()
		}
	}

	return &agent.wg
}

// WaitForCompletion waits for all the subscribers workers
// to be done handling their work, it does not stop  them
func (agent *BroadcastAgent) WaitForCompletion() {
	agent.wgBusy.Wait()
}
