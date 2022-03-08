package surfstore

import (
	context "context"
	"errors"
	"sync"
)

type BlockStore struct {
	BlockMap map[string]*Block
	UnimplementedBlockStoreServer
}

var blockLock sync.Mutex

func (bs *BlockStore) GetBlock(ctx context.Context, blockHash *BlockHash) (*Block, error) {
	blockLock.Lock()
	defer blockLock.Unlock()
	if block, ok := bs.BlockMap[blockHash.Hash]; ok {
		return block, nil
	} else {
		return nil, errors.New("getblockfailed")
	}
}

func (bs *BlockStore) PutBlock(ctx context.Context, block *Block) (*Success, error) {
	blockLock.Lock()
	defer blockLock.Unlock()
	bs.BlockMap[GetBlockHashString(block.BlockData)] = block
	return &Success{}, nil
}

// Given a list of hashes “in”, returns a list containing the
// subset of in that are stored in the key-value store
func (bs *BlockStore) HasBlocks(ctx context.Context, blockHashesIn *BlockHashes) (*BlockHashes, error) {
	blockLock.Lock()
	defer blockLock.Unlock()
	var blockHashesOut *BlockHashes
	for _, key := range blockHashesIn.Hashes {
		if _, ok := bs.BlockMap[key]; ok {
			blockHashesOut.Hashes = append(blockHashesOut.Hashes, key)
		}
	}
	return blockHashesOut, nil
}

// This line guarantees all method for BlockStore are implemented
var _ BlockStoreInterface = new(BlockStore)

func NewBlockStore() *BlockStore {
	return &BlockStore{
		BlockMap: map[string]*Block{},
	}
}
