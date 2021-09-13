package common

type (
	ChannelLimiter struct {
		bufferChannel chan int
		index int
	}
)

func NewChannelLimiter(limit int) *ChannelLimiter {
	return &ChannelLimiter{bufferChannel: make(chan int, limit)}
}

func (l *ChannelLimiter) Allow() bool {
	select {
	case l.bufferChannel <- 1:
		 l.index++
		return true
	default:
		return false
	}
}

func (l *ChannelLimiter) GetIndex() int{
	return l.index
}

func (l *ChannelLimiter) Release() bool {
	<-l.bufferChannel
	l.index--
	return true
}

func (l *ChannelLimiter) Close() {
	close(l.bufferChannel)
}


