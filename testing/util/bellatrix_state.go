package util

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateutil"
	v3 "github.com/prysmaticlabs/prysm/beacon-chain/state/v3"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	enginev1 "github.com/prysmaticlabs/prysm/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// DeterministicGenesisStateBellatrix returns a genesis state in Bellatrix format made using the deterministic deposits.
func DeterministicGenesisStateBellatrix(t testing.TB, numValidators uint64) (state.BeaconState, []bls.SecretKey) {
	deposits, privKeys, err := DeterministicDepositsAndKeys(numValidators)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "failed to get %d deposits", numValidators))
	}
	eth1Data, err := DeterministicEth1Data(len(deposits))
	if err != nil {
		t.Fatal(errors.Wrapf(err, "failed to get eth1data for %d deposits", numValidators))
	}
	beaconState, err := genesisBeaconStateBellatrix(context.Background(), deposits, uint64(0), eth1Data)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "failed to get genesis beacon state of %d validators", numValidators))
	}
	resetCache()
	return beaconState, privKeys
}

// genesisBeaconStateBellatrix returns the genesis beacon state.
func genesisBeaconStateBellatrix(ctx context.Context, deposits []*ethpb.Deposit, genesisTime uint64, eth1Data *ethpb.Eth1Data) (state.BeaconState, error) {
	st, err := emptyGenesisStateBellatrix()
	if err != nil {
		return nil, err
	}

	// Process initial deposits.
	st, err = helpers.UpdateGenesisEth1Data(st, deposits, eth1Data)
	if err != nil {
		return nil, err
	}

	st, err = processPreGenesisDeposits(ctx, st, deposits)
	if err != nil {
		return nil, errors.Wrap(err, "could not process validator deposits")
	}

	return buildGenesisBeaconStateBellatrix(genesisTime, st, st.Eth1Data())
}

// emptyGenesisStateBellatrix returns an empty genesis state in Bellatrix format.
func emptyGenesisStateBellatrix() (state.BeaconState, error) {
	st := &ethpb.BeaconStateBellatrix{
		// Misc fields.
		Slot: 0,
		Fork: &ethpb.Fork{
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
			CurrentVersion:  params.BeaconConfig().AltairForkVersion,
			Epoch:           0,
		},
		// Validator registry fields.
		Validators:       []*ethpb.Validator{},
		Balances:         []uint64{},
		InactivityScores: []uint64{},

		JustificationBits:          []byte{0},
		HistoricalRoots:            [][]byte{},
		CurrentEpochParticipation:  []byte{},
		PreviousEpochParticipation: []byte{},

		// Eth1 data.
		Eth1Data:         &ethpb.Eth1Data{},
		Eth1DataVotes:    []*ethpb.Eth1Data{},
		Eth1DepositIndex: 0,

		LatestExecutionPayloadHeader: &ethpb.ExecutionPayloadHeader{},
	}
	return v3.InitializeFromProto(st)
}

func buildGenesisBeaconStateBellatrix(genesisTime uint64, preState state.BeaconState, eth1Data *ethpb.Eth1Data) (state.BeaconState, error) {
	if eth1Data == nil {
		return nil, errors.New("no eth1data provided for genesis state")
	}

	randaoMixes := make([][]byte, params.BeaconConfig().EpochsPerHistoricalVector)
	for i := 0; i < len(randaoMixes); i++ {
		h := make([]byte, 32)
		copy(h, eth1Data.BlockHash)
		randaoMixes[i] = h
	}

	zeroHash := params.BeaconConfig().ZeroHash[:]

	activeIndexRoots := make([][]byte, params.BeaconConfig().EpochsPerHistoricalVector)
	for i := 0; i < len(activeIndexRoots); i++ {
		activeIndexRoots[i] = zeroHash
	}

	blockRoots := make([][]byte, params.BeaconConfig().SlotsPerHistoricalRoot)
	for i := 0; i < len(blockRoots); i++ {
		blockRoots[i] = zeroHash
	}

	stateRoots := make([][]byte, params.BeaconConfig().SlotsPerHistoricalRoot)
	for i := 0; i < len(stateRoots); i++ {
		stateRoots[i] = zeroHash
	}

	slashings := make([]uint64, params.BeaconConfig().EpochsPerSlashingsVector)

	genesisValidatorsRoot, err := stateutil.ValidatorRegistryRoot(preState.Validators())
	if err != nil {
		return nil, errors.Wrapf(err, "could not hash tree root genesis validators %v", err)
	}

	prevEpochParticipation, err := preState.PreviousEpochParticipation()
	if err != nil {
		return nil, err
	}
	currEpochParticipation, err := preState.CurrentEpochParticipation()
	if err != nil {
		return nil, err
	}
	scores, err := preState.InactivityScores()
	if err != nil {
		return nil, err
	}
	st := &ethpb.BeaconStateBellatrix{
		// Misc fields.
		Slot:                  0,
		GenesisTime:           genesisTime,
		GenesisValidatorsRoot: genesisValidatorsRoot[:],

		Fork: &ethpb.Fork{
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
			Epoch:           0,
		},

		// Validator registry fields.
		Validators:                 preState.Validators(),
		Balances:                   preState.Balances(),
		PreviousEpochParticipation: prevEpochParticipation,
		CurrentEpochParticipation:  currEpochParticipation,
		InactivityScores:           scores,

		// Randomness and committees.
		RandaoMixes: randaoMixes,

		// Finality.
		PreviousJustifiedCheckpoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  params.BeaconConfig().ZeroHash[:],
		},
		CurrentJustifiedCheckpoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  params.BeaconConfig().ZeroHash[:],
		},
		JustificationBits: []byte{0},
		FinalizedCheckpoint: &ethpb.Checkpoint{
			Epoch: 0,
			Root:  params.BeaconConfig().ZeroHash[:],
		},

		HistoricalRoots: [][]byte{},
		BlockRoots:      blockRoots,
		StateRoots:      stateRoots,
		Slashings:       slashings,

		// Eth1 data.
		Eth1Data:         eth1Data,
		Eth1DataVotes:    []*ethpb.Eth1Data{},
		Eth1DepositIndex: preState.Eth1DepositIndex(),
	}

	bodyRoot, err := (&ethpb.BeaconBlockBodyBellatrix{
		RandaoReveal: make([]byte, 96),
		Eth1Data: &ethpb.Eth1Data{
			DepositRoot: make([]byte, 32),
			BlockHash:   make([]byte, 32),
		},
		Graffiti: make([]byte, 32),
		SyncAggregate: &ethpb.SyncAggregate{
			SyncCommitteeBits:      make([]byte, len(bitfield.NewBitvector512())),
			SyncCommitteeSignature: make([]byte, 96),
		},
		ExecutionPayload: &enginev1.ExecutionPayload{
			ParentHash:    make([]byte, 32),
			FeeRecipient:  make([]byte, 20),
			StateRoot:     make([]byte, 32),
			ReceiptsRoot:  make([]byte, 32),
			LogsBloom:     make([]byte, 256),
			Random:        make([]byte, 32),
			BaseFeePerGas: make([]byte, 32),
			BlockHash:     make([]byte, 32),
		},
	}).HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not hash tree root empty block body")
	}

	st.LatestBlockHeader = &ethpb.BeaconBlockHeader{
		ParentRoot: zeroHash,
		StateRoot:  zeroHash,
		BodyRoot:   bodyRoot[:],
	}

	var pubKeys [][]byte
	for i := uint64(0); i < params.BeaconConfig().SyncCommitteeSize; i++ {
		pubKeys = append(pubKeys, bytesutil.PadTo([]byte{}, params.BeaconConfig().BLSPubkeyLength))
	}
	st.CurrentSyncCommittee = &ethpb.SyncCommittee{
		Pubkeys:         pubKeys,
		AggregatePubkey: bytesutil.PadTo([]byte{}, params.BeaconConfig().BLSPubkeyLength),
	}
	st.NextSyncCommittee = &ethpb.SyncCommittee{
		Pubkeys:         bytesutil.SafeCopy2dBytes(pubKeys),
		AggregatePubkey: bytesutil.PadTo([]byte{}, params.BeaconConfig().BLSPubkeyLength),
	}

	st.LatestExecutionPayloadHeader = &ethpb.ExecutionPayloadHeader{
		ParentHash:       make([]byte, 32),
		FeeRecipient:     make([]byte, 20),
		StateRoot:        make([]byte, 32),
		ReceiptRoot:      make([]byte, 32),
		LogsBloom:        make([]byte, 256),
		Random:           make([]byte, 32),
		BaseFeePerGas:    make([]byte, 32),
		BlockHash:        make([]byte, 32),
		TransactionsRoot: make([]byte, 32),
	}

	return v3.InitializeFromProto(st)
}
