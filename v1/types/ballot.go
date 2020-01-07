package types

import (
	"fmt"
	v1 "github.com/1xyz/paxossim/v1"
	"strings"
)

// A BallotNumber in paxos has the property that it has can
// lexicographically ordered. s.t they can be totally ordered
type BallotNumber struct {
	// a monotonically increasing integer
	Round int

	// Represents a the leader identifier with this ballot number
	LeaderID v1.Addr
}

func (bn BallotNumber) String() string {
	return fmt.Sprintf("(%d, %v)", bn.Round, bn.LeaderID)
}

// Compare two Ballot numbers lexicographically
// The result will be:
//   0 if bn1 == bn2,
//   -1 if bn1 < bn2, and
//   +1 if bn1 > bn2.
func Compare(bn1 *BallotNumber, bn2 *BallotNumber) int {
	res := returnOne(bn1.Round - bn2.Round)
	if res == 0 {
		id1 := int(bn1.LeaderID.ID())
		id2 := int(bn2.LeaderID.ID())
		return returnOne(id1 - id2)
	}
	return res
}

func returnOne(res int) int {
	if res == 0 {
		return res
	} else if res < 0 {
		return -1
	} else if res > 0 {
		return +1
	}
}
