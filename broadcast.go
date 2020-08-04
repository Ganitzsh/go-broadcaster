package broadcaster

// Broadcast is the payload sent across the subscribers of an
// agent
type Broadcast struct {
	Frequency string
	Data      interface{}
}

// NewBroadcast will return a new instance of BNewBroadcast
func newBroadcast(freq string, data interface{}) *Broadcast {
	return &Broadcast{
		Frequency: freq,
		Data:      data,
	}
}
