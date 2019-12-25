package paxossim_test

import (
	"github.com/1xyz/paxossim"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	L1 = "Leader:1"
	L2 = "Leader:2"
	R1 = 0
	R2 = 10
)

func TestBallotNumber(t *testing.T) {
	bn1 := &paxossim.BallotNumber{
		Round:    R1,
		LeaderID: L1,
	}
	bn2 := &paxossim.BallotNumber{
		Round:    R2,
		LeaderID: L2,
	}
	assert.Equal(t, -1, bn1.CompareTo(bn2))
	assert.Equal(t, 1, bn2.CompareTo(bn1))

	bn3 := &paxossim.BallotNumber{
		Round:    R1,
		LeaderID: L2,
	}
	assert.Equal(t, -1, bn1.CompareTo(bn3))

	bn4 := &paxossim.BallotNumber{
		Round:    R1,
		LeaderID: L1,
	}
	assert.Equal(t, 0, bn1.CompareTo(bn4))
	assert.Equal(t, 0, bn1.CompareTo(bn1))
}
