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
	res := bn1.Round - bn2.Round
	if res == 0 {
		addr1 := string(bn1.LeaderID.ID())
		addr2 := string(bn2.LeaderID.ID())
		return strings.Compare(addr1, addr2)
	} else if res < 0 {
		return -1
	} else {
		return +1
	}
}
