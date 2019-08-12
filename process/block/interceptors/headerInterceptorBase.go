package interceptors

import (
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/block"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

// HeaderInterceptorBase is the "abstract class" extended in HeaderInterceptor and ShardHeaderInterceptor
type HeaderInterceptorBase struct {
	marshalizer         marshal.Marshalizer
	headerValidator     process.HeaderValidator
	multiSigVerifier    crypto.MultiSigVerifier
	hasher              hashing.Hasher
	shardCoordinator    sharding.Coordinator
	chronologyValidator process.ChronologyValidator
}

// NewHeaderInterceptorBase creates a new HeaderIncterceptorBase instance
func NewHeaderInterceptorBase(
	marshalizer marshal.Marshalizer,
	headerValidator process.HeaderValidator,
	multiSigVerifier crypto.MultiSigVerifier,
	hasher hashing.Hasher,
	shardCoordinator sharding.Coordinator,
	chronologyValidator process.ChronologyValidator,
) (*HeaderInterceptorBase, error) {
	if marshalizer == nil {
		return nil, process.ErrNilMarshalizer
	}
	if headerValidator == nil {
		return nil, process.ErrNilHeaderHandlerValidator
	}
	if multiSigVerifier == nil {
		return nil, process.ErrNilMultiSigVerifier
	}
	if hasher == nil {
		return nil, process.ErrNilHasher
	}
	if shardCoordinator == nil {
		return nil, process.ErrNilShardCoordinator
	}
	if chronologyValidator == nil {
		return nil, process.ErrNilChronologyValidator
	}

	hdrIntercept := &HeaderInterceptorBase{
		marshalizer:         marshalizer,
		headerValidator:     headerValidator,
		multiSigVerifier:    multiSigVerifier,
		hasher:              hasher,
		shardCoordinator:    shardCoordinator,
		chronologyValidator: chronologyValidator,
	}

	return hdrIntercept, nil
}

// ParseReceivedMessage will transform the received p2p.Message in an InterceptedHeader.
// If the header hash is present in storage it will output an error
func (hib *HeaderInterceptorBase) ParseReceivedMessage(message p2p.MessageP2P) (*block.InterceptedHeader, error) {
	if message == nil {
		return nil, process.ErrNilMessage
	}
	if message.Data() == nil {
		return nil, process.ErrNilDataToProcess
	}

	hdrIntercepted := block.NewInterceptedHeader(hib.multiSigVerifier, hib.chronologyValidator)
	err := hib.marshalizer.Unmarshal(hdrIntercepted, message.Data())
	if err != nil {
		return nil, err
	}

	hashWithSig := hib.hasher.Compute(string(message.Data()))
	hdrIntercepted.SetHash(hashWithSig)

	err = hdrIntercepted.IntegrityAndValidity(hib.shardCoordinator)
	if err != nil {
		return nil, err
	}

	err = hdrIntercepted.VerifySig()
	if err != nil {
		return nil, err
	}

	return hdrIntercepted, nil
}

// CheckHeaderForCurrentShard checks if the header is for current shard
func (hib *HeaderInterceptorBase) CheckHeaderForCurrentShard(interceptedHdr *block.InterceptedHeader) bool {
	isHeaderForCurrentShard := hib.shardCoordinator.SelfId() == interceptedHdr.GetHeader().ShardId
	isMetachainShardCoordinator := hib.shardCoordinator.SelfId() == sharding.MetachainShardId

	return isHeaderForCurrentShard || isMetachainShardCoordinator
}
