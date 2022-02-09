package membership

// QueuedChangeHandler is a wrapper around a ChangeHandler that buffers membership changes to a channel.
//
// (This is meant to be used by the Discord bot to apply its own membership checks.)
type QueuedChangeHandler struct {
	*internalQueuedChangeHandler
}

// GetGained() returns a slice of queued-up integers.
func (q *QueuedChangeHandler) GetGained() []int {
	return q.sliceFromChannel(q.gained)
}

// GetGained() returns a slice of queued-up integers.
func (q *QueuedChangeHandler) GetLost() []int {
	return q.sliceFromChannel(q.lost)
}

func (q *QueuedChangeHandler) sliceFromChannel(c chan int) []int {
	var acc []int
	for {
		select {
		case item := <-c:
			acc = append(acc, item)
		default:
			return acc
		}
	}
}

type internalQueuedChangeHandler struct {
	gained, lost chan int

	ChangeHandler
}

func (c *internalQueuedChangeHandler) GainedMembership(userMembershipID int) {
	c.gained <- userMembershipID
}
func (c *internalQueuedChangeHandler) LostMembership(userMembershipID int) {
	c.lost <- userMembershipID
}

// NewQueuedChangeHandler creates a QueuedChangeHandler that puts UserMembership ent IDs in two separate, buffered channels that can be exhausted when "safe" to do so. Use the second argument in HookMembershipChanges.
//
// TODO: notes on how to wrap this properly with a mutex
func NewQueuedChangeHandler(bufferSize int) (*QueuedChangeHandler, ChangeHandler) {
	var (
		gained = make(chan int, bufferSize)
		lost   = make(chan int, bufferSize)
	)
	handler := &QueuedChangeHandler{
		internalQueuedChangeHandler: &internalQueuedChangeHandler{
			gained: gained,
			lost:   lost,
		},
	}
	return handler, handler.internalQueuedChangeHandler
}
