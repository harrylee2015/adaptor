package common

type (
	ChannelLimiter struct {
		bufferChannel chan int
	}
)

func NewChannelLimiter(limit int) *ChannelLimiter {
	return &ChannelLimiter{bufferChannel: make(chan int, limit)}
}

func (l *ChannelLimiter) Allow() bool {
	select {
	case l.bufferChannel <- 1:
		return true
	default:
		return false
	}
}

func (l *ChannelLimiter) Release() bool {
	<-l.bufferChannel
	return true
}

func (l *ChannelLimiter) Close() {
	close(l.bufferChannel)
}
