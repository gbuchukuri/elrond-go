package requestHandlers

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-sandbox/core"
	"github.com/ElrondNetwork/elrond-go-sandbox/core/logger"
	"github.com/ElrondNetwork/elrond-go-sandbox/core/partitioning"
	"github.com/ElrondNetwork/elrond-go-sandbox/dataRetriever"
)

type ResolverRequestHandler struct {
	resolversFinder dataRetriever.ResolversFinder
	txRequestTopic  string
	mbRequestTopic  string
	hdrRequestTopic string
	isMetaChain     bool
	maxTxsToRequest int
}

var log = logger.DefaultLogger()

// NewShardResolverRequestHandler creates a requestHandler interface implementation with request functions
func NewShardResolverRequestHandler(
	finder dataRetriever.ResolversFinder,
	txRequestTopic string,
	mbRequestTopic string,
	hdrRequestTopic string,
	maxTxsToRequest int,
) (*ResolverRequestHandler, error) {
	if finder == nil {
		return nil, dataRetriever.ErrNilResolverFinder
	}
	if len(txRequestTopic) == 0 {
		return nil, dataRetriever.ErrEmptyTxRequestTopic
	}
	if len(mbRequestTopic) == 0 {
		return nil, dataRetriever.ErrEmptyMiniBlockRequestTopic
	}
	if len(hdrRequestTopic) == 0 {
		return nil, dataRetriever.ErrEmptyHeaderRequestTopic
	}
	if maxTxsToRequest < 1 {
		return nil, dataRetriever.ErrInvalidMaxTxRequest
	}

	rrh := &ResolverRequestHandler{
		resolversFinder: finder,
		txRequestTopic:  txRequestTopic,
		mbRequestTopic:  mbRequestTopic,
		hdrRequestTopic: hdrRequestTopic,
		isMetaChain:     false,
		maxTxsToRequest: maxTxsToRequest,
	}

	return rrh, nil
}

// NewMetaResolverRequestHandler creates a requestHandler interface implementation with request functions
func NewMetaResolverRequestHandler(
	finder dataRetriever.ResolversFinder,
	hdrRequestTopic string,
) (*ResolverRequestHandler, error) {
	if finder == nil {
		return nil, dataRetriever.ErrNilResolverFinder
	}
	if len(hdrRequestTopic) == 0 {
		return nil, dataRetriever.ErrEmptyHeaderRequestTopic
	}

	rrh := &ResolverRequestHandler{
		resolversFinder: finder,
		hdrRequestTopic: hdrRequestTopic,
		isMetaChain:     true,
	}

	return rrh, nil
}

// RequestTransaction method asks for transactions from the connected peers
func (rrh *ResolverRequestHandler) RequestTransaction(destShardID uint32, txHashes [][]byte) {
	log.Debug(fmt.Sprintf("Requesting %d transactions from shard %d from network...\n", len(txHashes), destShardID))
	resolver, err := rrh.resolversFinder.CrossShardResolver(rrh.txRequestTopic, destShardID)
	if err != nil {
		log.Error(fmt.Sprintf("missing resolver to %s topic to shard %d", rrh.txRequestTopic, destShardID))
		return
	}

	txResolver, ok := resolver.(HashSliceResolver)
	if !ok {
		log.Error("wrong assertion type when creating transaction resolver")
		return
	}

	go func() {
		dataSplit := &partitioning.DataSplit{}
		sliceBatches, err := dataSplit.SplitDataInChunks(txHashes, rrh.maxTxsToRequest)
		if err != nil {
			log.Error("error requesting transactions: " + err.Error())
			return
		}

		for _, batch := range sliceBatches {
			err = txResolver.RequestDataFromHashArray(batch)
			if err != nil {
				log.Debug("error requesting tx batch: " + err.Error())
			}
		}
	}()
}

// RequestMiniBlock method asks for miniblocks from the connected peers
func (rrh *ResolverRequestHandler) RequestMiniBlock(shardId uint32, miniblockHash []byte) {
	rrh.requestByHash(shardId, miniblockHash, rrh.mbRequestTopic)
}

// RequestHeader method asks for header from the connected peers
func (rrh *ResolverRequestHandler) RequestHeader(shardId uint32, hash []byte) {
	rrh.requestByHash(shardId, hash, rrh.hdrRequestTopic)
}

func (rrh *ResolverRequestHandler) requestByHash(destShardID uint32, hash []byte, baseTopic string) {
	log.Debug(fmt.Sprintf("Requesting %s from shard %d with hash %s from network\n", baseTopic, destShardID, core.ToB64(hash)))
	resolver, err := rrh.resolversFinder.CrossShardResolver(baseTopic, destShardID)
	if err != nil {
		log.Error(fmt.Sprintf("missing resolver to %s topic to shard %d", baseTopic, destShardID))
		return
	}

	err = resolver.RequestDataFromHash(hash)
	if err != nil {
		log.Debug(err.Error())
	}
}

// RequestHeaderByNonce method asks for transactions from the connected peers
func (rrh *ResolverRequestHandler) RequestHeaderByNonce(destShardID uint32, nonce uint64) {
	var err error
	var resolver dataRetriever.Resolver
	if rrh.isMetaChain {
		resolver, err = rrh.resolversFinder.CrossShardResolver(rrh.hdrRequestTopic, destShardID)
	} else {
		resolver, err = rrh.resolversFinder.MetaChainResolver(rrh.hdrRequestTopic)
	}

	if err != nil {
		log.Error(fmt.Sprintf("missing resolver to %s topic to shard %d", rrh.hdrRequestTopic, destShardID))
		return
	}

	headerResolver, ok := resolver.(dataRetriever.HeaderResolver)
	if !ok {
		log.Error(fmt.Sprintf("resolver is not a header resolver to %s topic to shard %d", rrh.hdrRequestTopic, destShardID))
		return
	}

	err = headerResolver.RequestDataFromNonce(nonce)
	if err != nil {
		log.Debug(err.Error())
	}
}