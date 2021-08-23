package client

//import (
//	"github.com/33cn/chain33/types"
//	"sync"
//)
//
//type Pool interface {
//
//	Get() (types.Chain33Client, error)
//
//	Close() error
//
//	Status() string
//}
//
//type pool struct {
//	// atomic, used to get connection random
//	index uint32
//
//	// atomic, the current physical connection of pool
//	current int32
//
//	// atomic, the using logic connection of pool
//	// logic connection = physical connection * MaxConcurrentStreams
//	ref int32
//
//	// pool options
//	opt Options
//
//	// all of created physical connections
//	conns []types.Chain33Client
//
//	// the server address is to create connection.
//	address string
//
//	// control the atomic var current's concurrent read write.
//	sync.RWMutex
//}