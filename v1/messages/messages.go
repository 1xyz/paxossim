package messages

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/types"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
)

type basicMessage struct {
	src v1.Addr
}

func (bm basicMessage) Src() v1.Addr {
	return bm.src
}

func (bm basicMessage) String() string {
	return fmt.Sprintf("source %v", bm.src)
}

// RequestMessage - Request from Client to all Replicas encapsulating a Command
type RequestMessage struct {
	basicMessage
	Command types.Command
}

func (rm RequestMessage) String() string {
	return fmt.Sprintf("%v %v", rm.basicMessage, rm.Command)
}
