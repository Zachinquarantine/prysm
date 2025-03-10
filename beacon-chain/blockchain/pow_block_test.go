package blockchain

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/testing/require"
)

func Test_validTerminalPowBlock(t *testing.T) {
	tests := []struct {
		name              string
		currentDifficulty *uint256.Int
		parentDifficulty  *uint256.Int
		ttd               uint64
		want              bool
	}{
		{
			name:              "current > ttd, parent > ttd",
			currentDifficulty: uint256.NewInt(2),
			parentDifficulty:  uint256.NewInt(2),
			ttd:               1,
			want:              false,
		},
		{
			name:              "current < ttd, parent < ttd",
			currentDifficulty: uint256.NewInt(2),
			parentDifficulty:  uint256.NewInt(2),
			ttd:               3,
			want:              false,
		},
		{
			name:              "current == ttd, parent == ttd",
			currentDifficulty: uint256.NewInt(2),
			parentDifficulty:  uint256.NewInt(2),
			ttd:               2,
			want:              false,
		},
		{
			name:              "current > ttd, parent == ttd",
			currentDifficulty: uint256.NewInt(2),
			parentDifficulty:  uint256.NewInt(1),
			ttd:               1,
			want:              false,
		},
		{
			name:              "current == ttd, parent < ttd",
			currentDifficulty: uint256.NewInt(2),
			parentDifficulty:  uint256.NewInt(1),
			ttd:               2,
			want:              true,
		},
		{
			name:              "current > ttd, parent < ttd",
			currentDifficulty: uint256.NewInt(3),
			parentDifficulty:  uint256.NewInt(1),
			ttd:               2,
			want:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := params.BeaconConfig()
			cfg.TerminalTotalDifficulty = fmt.Sprint(tt.ttd)
			params.OverrideBeaconConfig(cfg)
			got, err := validTerminalPowBlock(tt.currentDifficulty, tt.parentDifficulty)
			require.NoError(t, err)
			if got != tt.want {
				t.Errorf("validTerminalPowBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validTerminalPowBlockSpecConfig(t *testing.T) {
	cfg := params.BeaconConfig()
	cfg.TerminalTotalDifficulty = "115792089237316195423570985008687907853269984665640564039457584007913129638912"
	params.OverrideBeaconConfig(cfg)

	i, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129638912", 10)
	current, of := uint256.FromBig(i)
	require.Equal(t, of, false)
	i, _ = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129638911", 10)
	parent, of := uint256.FromBig(i)
	require.Equal(t, of, false)

	got, err := validTerminalPowBlock(current, parent)
	require.NoError(t, err)
	require.Equal(t, true, got)
}
