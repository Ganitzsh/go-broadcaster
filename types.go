package broadcaster

// SubscriberHandler is a function called when a broadcast is received
// by a subscriber
type SubscriberHandler func(broadcast *Broadcast) error
