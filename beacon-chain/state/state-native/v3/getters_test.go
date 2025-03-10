package v3

import (
	"sync"
	"testing"

	types "github.com/prysmaticlabs/eth2-types"
	fieldparams "github.com/prysmaticlabs/prysm/config/fieldparams"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/testing/assert"
	"github.com/prysmaticlabs/prysm/testing/require"
)

func TestBeaconState_SlotDataRace(t *testing.T) {
	headState, err := InitializeFromProto(&ethpb.BeaconStateBellatrix{Slot: 1})
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		require.NoError(t, headState.SetSlot(0))
		wg.Done()
	}()
	go func() {
		headState.Slot()
		wg.Done()
	}()

	wg.Wait()
}

func TestBeaconState_ValidatorByPubkey(t *testing.T) {
	keyCreator := func(input []byte) [fieldparams.BLSPubkeyLength]byte {
		nKey := [fieldparams.BLSPubkeyLength]byte{}
		copy(nKey[:1], input)
		return nKey
	}

	tests := []struct {
		name            string
		modifyFunc      func(b *BeaconState, k [fieldparams.BLSPubkeyLength]byte)
		exists          bool
		expectedIdx     types.ValidatorIndex
		largestIdxInSet types.ValidatorIndex
	}{
		{
			name: "retrieve validator",
			modifyFunc: func(b *BeaconState, key [fieldparams.BLSPubkeyLength]byte) {
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key[:]}))
			},
			exists:      true,
			expectedIdx: 0,
		},
		{
			name: "retrieve validator with multiple validators from the start",
			modifyFunc: func(b *BeaconState, key [fieldparams.BLSPubkeyLength]byte) {
				key1 := keyCreator([]byte{'C'})
				key2 := keyCreator([]byte{'D'})
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key[:]}))
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key1[:]}))
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key2[:]}))
			},
			exists:      true,
			expectedIdx: 0,
		},
		{
			name: "retrieve validator with multiple validators",
			modifyFunc: func(b *BeaconState, key [fieldparams.BLSPubkeyLength]byte) {
				key1 := keyCreator([]byte{'C'})
				key2 := keyCreator([]byte{'D'})
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key1[:]}))
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key2[:]}))
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key[:]}))
			},
			exists:      true,
			expectedIdx: 2,
		},
		{
			name: "retrieve validator with multiple validators from the start with shared state",
			modifyFunc: func(b *BeaconState, key [fieldparams.BLSPubkeyLength]byte) {
				key1 := keyCreator([]byte{'C'})
				key2 := keyCreator([]byte{'D'})
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key[:]}))
				_ = b.Copy()
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key1[:]}))
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key2[:]}))
			},
			exists:      true,
			expectedIdx: 0,
		},
		{
			name: "retrieve validator with multiple validators with shared state",
			modifyFunc: func(b *BeaconState, key [fieldparams.BLSPubkeyLength]byte) {
				key1 := keyCreator([]byte{'C'})
				key2 := keyCreator([]byte{'D'})
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key1[:]}))
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key2[:]}))
				n := b.Copy()
				// Append to another state
				assert.NoError(t, n.AppendValidator(&ethpb.Validator{PublicKey: key[:]}))

			},
			exists:      false,
			expectedIdx: 0,
		},
		{
			name: "retrieve validator with multiple validators with shared state at boundary",
			modifyFunc: func(b *BeaconState, key [fieldparams.BLSPubkeyLength]byte) {
				key1 := keyCreator([]byte{'C'})
				assert.NoError(t, b.AppendValidator(&ethpb.Validator{PublicKey: key1[:]}))
				n := b.Copy()
				// Append to another state
				assert.NoError(t, n.AppendValidator(&ethpb.Validator{PublicKey: key[:]}))

			},
			exists:      false,
			expectedIdx: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := InitializeFromProto(&ethpb.BeaconStateBellatrix{})
			require.NoError(t, err)
			nKey := keyCreator([]byte{'A'})
			tt.modifyFunc(s, nKey)
			idx, ok := s.ValidatorIndexByPubkey(nKey)
			assert.Equal(t, tt.exists, ok)
			assert.Equal(t, tt.expectedIdx, idx)
		})
	}
}
