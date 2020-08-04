package broadcaster

import (
	"sync"
)

// Subscriber is is registered to a broadcast agent and will receive messages
// through a channel
type Subscriber struct {
	Comm    chan *Broadcast
	ErrChan chan error

	workersAmmount uint32
	frequency      string
	handler        SubscriberHandler
	wg             sync.WaitGroup
	agentWg        *sync.WaitGroup
	agentWgBusy    *sync.WaitGroup
}

func newSubscriber(frequency string, handler SubscriberHandler, workersAmmount uint32, wgBusy, wg *sync.WaitGroup) *Subscriber {
	// ErrChan needs to be buffered to be non blocking in case it's not used
	errChanBufferSize := 1
	if workersAmmount > 0 {
		errChanBufferSize = int(workersAmmount)
	}

	return &Subscriber{
		Comm:    make(chan *Broadcast, workersAmmount),
		ErrChan: make(chan error, errChanBufferSize),

		workersAmmount: workersAmmount,
		frequency:      frequency,
		handler:        handler,
		agentWg:        wg,
		agentWgBusy:    wgBusy,
		wg:             sync.WaitGroup{},
	}
}

func (sub *Subscriber) stop() {
	close(sub.Comm)
	sub.wg.Wait()
	close(sub.ErrChan)
}

func (sub *Subscriber) startWorker() {
	defer sub.agentWg.Done()
	defer sub.wg.Done()

	for {
		broadcast := <-sub.Comm
		if broadcast == nil {
			return
		}

		sub.agentWgBusy.Add(1)
		if err := sub.handler(broadcast); err != nil {
			// for non blocking error buffering in the chan
			// avoids deadlocks if ErrChan is never read
			select {
			case sub.ErrChan <- err:
				break
			}
		}
		sub.agentWgBusy.Done()
	}
}

// Start will start listening to the primary Comm channel for message
// broadcasted on its agent
func (sub *Subscriber) Start() {
	if sub.workersAmmount == 0 {
		sub.agentWg.Add(1)
		sub.wg.Add(1)
		go sub.startWorker()

		return
	}

	for i := uint32(0); i < sub.workersAmmount; i++ {
		sub.agentWg.Add(1)
		sub.wg.Add(1)
		go sub.startWorker()
	}
}
