package bn_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go-sandbox/consensus/spos"
	"github.com/ElrondNetwork/elrond-go-sandbox/consensus/spos/bn"
	"github.com/ElrondNetwork/elrond-go-sandbox/consensus/spos/mock"
	"github.com/ElrondNetwork/elrond-go-sandbox/crypto"
	"github.com/ElrondNetwork/elrond-go-sandbox/data"
	"github.com/ElrondNetwork/elrond-go-sandbox/data/block"
	"github.com/ElrondNetwork/elrond-go-sandbox/data/blockchain"
	"github.com/stretchr/testify/assert"
)

func initSubroundEndRound() bn.SubroundEndRound {
	blockChain := mock.BlockChainMock{}
	blockProcessorMock := mock.InitBlockProcessorMock()
	consensusState := mock.InitConsensusState()
	multiSignerMock := mock.InitMultiSignerMock()
	rounderMock := initRounderMock()
	syncTimerMock := mock.SyncTimerMock{}

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, _ := bn.NewSubroundEndRound(
		sr,
		&blockChain,
		blockProcessorMock,
		consensusState,
		multiSignerMock,
		rounderMock,
		syncTimerMock,
		broadcastBlock,
		extend,
	)

	return srEndRound
}

func TestSubroundEndRound_NewSubroundEndRoundNilSubroundShouldFail(t *testing.T) {
	t.Parallel()

	blockChain := mock.BlockChainMock{}
	blockProcessorMock := mock.InitBlockProcessorMock()
	consensusState := mock.InitConsensusState()
	multiSignerMock := mock.InitMultiSignerMock()
	rounderMock := initRounderMock()
	syncTimerMock := mock.SyncTimerMock{}

	srEndRound, err := bn.NewSubroundEndRound(
		nil,
		&blockChain,
		blockProcessorMock,
		consensusState,
		multiSignerMock,
		rounderMock,
		syncTimerMock,
		broadcastBlock,
		extend,
	)

	assert.Nil(t, srEndRound)
	assert.Equal(t, err, spos.ErrNilSubround)
}

func TestSubroundEndRound_NewSubroundEndRoundNilBlockChainShouldFail(t *testing.T) {
	t.Parallel()

	blockProcessorMock := mock.InitBlockProcessorMock()
	consensusState := mock.InitConsensusState()
	multiSignerMock := mock.InitMultiSignerMock()
	rounderMock := initRounderMock()
	syncTimerMock := mock.SyncTimerMock{}

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, err := bn.NewSubroundEndRound(
		sr,
		nil,
		blockProcessorMock,
		consensusState,
		multiSignerMock,
		rounderMock,
		syncTimerMock,
		broadcastBlock,
		extend,
	)

	assert.Nil(t, srEndRound)
	assert.Equal(t, err, spos.ErrNilBlockChain)
}

func TestSubroundEndRound_NewSubroundEndRoundNilBlockProcessorShouldFail(t *testing.T) {
	t.Parallel()

	blockChain := mock.BlockChainMock{}
	consensusState := mock.InitConsensusState()
	multiSignerMock := mock.InitMultiSignerMock()
	rounderMock := initRounderMock()
	syncTimerMock := mock.SyncTimerMock{}

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, err := bn.NewSubroundEndRound(
		sr,
		&blockChain,
		nil,
		consensusState,
		multiSignerMock,
		rounderMock,
		syncTimerMock,
		broadcastBlock,
		extend,
	)

	assert.Nil(t, srEndRound)
	assert.Equal(t, err, spos.ErrNilBlockProcessor)
}

func TestSubroundEndRound_NewSubroundEndRoundNilConsensusStateShouldFail(t *testing.T) {
	t.Parallel()

	blockChain := mock.BlockChainMock{}
	blockProcessorMock := mock.InitBlockProcessorMock()
	multiSignerMock := mock.InitMultiSignerMock()
	rounderMock := initRounderMock()
	syncTimerMock := mock.SyncTimerMock{}

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, err := bn.NewSubroundEndRound(
		sr,
		&blockChain,
		blockProcessorMock,
		nil,
		multiSignerMock,
		rounderMock,
		syncTimerMock,
		broadcastBlock,
		extend,
	)

	assert.Nil(t, srEndRound)
	assert.Equal(t, err, spos.ErrNilConsensusState)
}

func TestSubroundEndRound_NewSubroundEndRoundNilMultisignerShouldFail(t *testing.T) {
	t.Parallel()

	blockChain := mock.BlockChainMock{}
	blockProcessorMock := mock.InitBlockProcessorMock()
	consensusState := mock.InitConsensusState()
	rounderMock := initRounderMock()
	syncTimerMock := mock.SyncTimerMock{}

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, err := bn.NewSubroundEndRound(
		sr,
		&blockChain,
		blockProcessorMock,
		consensusState,
		nil,
		rounderMock,
		syncTimerMock,
		broadcastBlock,
		extend,
	)

	assert.Nil(t, srEndRound)
	assert.Equal(t, err, spos.ErrNilMultiSigner)
}

func TestSubroundEndRound_NewSubroundEndRoundNilRounderShouldFail(t *testing.T) {
	t.Parallel()

	blockChain := mock.BlockChainMock{}
	blockProcessorMock := mock.InitBlockProcessorMock()
	consensusState := mock.InitConsensusState()
	multiSignerMock := mock.InitMultiSignerMock()
	syncTimerMock := mock.SyncTimerMock{}

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, err := bn.NewSubroundEndRound(
		sr,
		&blockChain,
		blockProcessorMock,
		consensusState,
		multiSignerMock,
		nil,
		syncTimerMock,
		broadcastBlock,
		extend,
	)

	assert.Nil(t, srEndRound)
	assert.Equal(t, err, spos.ErrNilRounder)
}

func TestSubroundEndRound_NewSubroundEndRoundNilSyncTimerShouldFail(t *testing.T) {
	t.Parallel()

	blockChain := mock.BlockChainMock{}
	blockProcessorMock := mock.InitBlockProcessorMock()
	consensusState := mock.InitConsensusState()
	multiSignerMock := mock.InitMultiSignerMock()
	rounderMock := initRounderMock()

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, err := bn.NewSubroundEndRound(
		sr,
		&blockChain,
		blockProcessorMock,
		consensusState,
		multiSignerMock,
		rounderMock,
		nil,
		broadcastBlock,
		extend,
	)

	assert.Nil(t, srEndRound)
	assert.Equal(t, err, spos.ErrNilSyncTimer)
}

func TestSubroundEndRound_NewSubroundEndRoundNilBroadcastBlockFunctionShouldFail(t *testing.T) {
	t.Parallel()

	blockChain := mock.BlockChainMock{}
	blockProcessorMock := mock.InitBlockProcessorMock()
	consensusState := mock.InitConsensusState()
	multiSignerMock := mock.InitMultiSignerMock()
	rounderMock := initRounderMock()
	syncTimerMock := mock.SyncTimerMock{}

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, err := bn.NewSubroundEndRound(
		sr,
		&blockChain,
		blockProcessorMock,
		consensusState,
		multiSignerMock,
		rounderMock,
		syncTimerMock,
		nil,
		extend,
	)

	assert.Nil(t, srEndRound)
	assert.Equal(t, err, spos.ErrNilBroadcastBlockFunction)
}

func TestSubroundEndRound_NewSubroundEndRoundShouldWork(t *testing.T) {
	t.Parallel()

	blockChain := mock.BlockChainMock{}
	blockProcessorMock := mock.InitBlockProcessorMock()
	consensusState := mock.InitConsensusState()
	multiSignerMock := mock.InitMultiSignerMock()
	rounderMock := initRounderMock()
	syncTimerMock := mock.SyncTimerMock{}

	ch := make(chan bool, 1)

	sr, _ := bn.NewSubround(
		int(bn.SrSignature),
		int(bn.SrEndRound),
		-1,
		int64(85*roundTimeDuration/100),
		int64(95*roundTimeDuration/100),
		"(END_ROUND)",
		ch,
	)

	srEndRound, err := bn.NewSubroundEndRound(
		sr,
		&blockChain,
		blockProcessorMock,
		consensusState,
		multiSignerMock,
		rounderMock,
		syncTimerMock,
		broadcastBlock,
		extend,
	)

	assert.NotNil(t, srEndRound)
	assert.Nil(t, err)
}

func TestSubroundEndRound_DoEndRoundJobErrAggregatingSigShouldFail(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	multiSignerMock := mock.InitMultiSignerMock()

	multiSignerMock.AggregateSigsMock = func(bitmap []byte) ([]byte, error) {
		return nil, crypto.ErrNilHasher
	}

	sr.SetMultiSigner(multiSignerMock)
	sr.ConsensusState().Header = &block.Header{}

	r := sr.DoEndRoundJob()
	assert.False(t, r)
}

func TestSubroundEndRound_DoEndRoundJobErrCommitBlockShouldFail(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	blProcMock := mock.InitBlockProcessorMock()

	blProcMock.CommitBlockCalled = func(
		blockChain data.ChainHandler,
		header data.HeaderHandler,
		body data.BodyHandler,
	) error {
		return blockchain.ErrHeaderUnitNil
	}

	sr.SetBlockProcessor(blProcMock)
	sr.ConsensusState().Header = &block.Header{}

	r := sr.DoEndRoundJob()
	assert.False(t, r)
}

func TestSubroundEndRound_DoEndRoundJobErrBroadcastBlockOK(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	sr.SetBroadcastBlock(func(data.BodyHandler, data.HeaderHandler) error {
		return spos.ErrNilBroadcastBlockFunction
	})

	sr.ConsensusState().Header = &block.Header{}

	r := sr.DoEndRoundJob()
	assert.True(t, r)
}

func TestSubroundEndRound_DoEndRoundJobAllOK(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	sr.ConsensusState().Header = &block.Header{}

	r := sr.DoEndRoundJob()
	assert.True(t, r)
}

func TestSubroundEndRound_DoEndRoundConsensusCheckShouldReturnFalseWhenRoundIsCanceled(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	sr.ConsensusState().RoundCanceled = true

	ok := sr.DoEndRoundConsensusCheck()
	assert.False(t, ok)
}

func TestSubroundEndRound_DoEndRoundConsensusCheckShouldReturnTrueWhenRoundIsFinished(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	sr.ConsensusState().SetStatus(bn.SrEndRound, spos.SsFinished)

	ok := sr.DoEndRoundConsensusCheck()
	assert.True(t, ok)
}

func TestSubroundEndRound_DoEndRoundConsensusCheckShouldReturnFalseWhenRoundIsNotFinished(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	ok := sr.DoEndRoundConsensusCheck()
	assert.False(t, ok)
}

func TestSubroundEndRound_CheckSignaturesValidityShouldErrNilSignature(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	err := sr.CheckSignaturesValidity([]byte(string(2)))
	assert.Equal(t, spos.ErrNilSignature, err)
}

func TestSubroundEndRound_CheckSignaturesValidityShouldErrIndexOutOfBounds(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	_, _ = sr.MultiSigner().Create(nil, 0)

	sr.ConsensusState().SetJobDone(sr.ConsensusState().ConsensusGroup()[0], bn.SrSignature, true)

	multiSignerMock := mock.InitMultiSignerMock()
	multiSignerMock.SignatureShareMock = func(index uint16) ([]byte, error) {
		return nil, crypto.ErrIndexOutOfBounds
	}

	sr.SetMultiSigner(multiSignerMock)

	err := sr.CheckSignaturesValidity([]byte(string(1)))
	assert.Equal(t, crypto.ErrIndexOutOfBounds, err)
}

func TestSubroundEndRound_CheckSignaturesValidityShouldErrInvalidSignatureShare(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	multiSignerMock := mock.InitMultiSignerMock()

	err := errors.New("invalid signature share")
	multiSignerMock.VerifySignatureShareMock = func(index uint16, sig []byte, bitmap []byte) error {
		return err
	}

	sr.SetMultiSigner(multiSignerMock)

	sr.ConsensusState().SetJobDone(sr.ConsensusState().ConsensusGroup()[0], bn.SrSignature, true)

	err2 := sr.CheckSignaturesValidity([]byte(string(1)))
	assert.Equal(t, err, err2)
}

func TestSubroundEndRound_CheckSignaturesValidityShouldRetunNil(t *testing.T) {
	t.Parallel()

	sr := *initSubroundEndRound()

	sr.ConsensusState().SetJobDone(sr.ConsensusState().ConsensusGroup()[0], bn.SrSignature, true)

	err := sr.CheckSignaturesValidity([]byte(string(1)))
	assert.Equal(t, nil, err)
}