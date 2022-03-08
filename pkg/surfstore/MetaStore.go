package surfstore

import (
	context "context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	// "errors"
	"sync"
)

type MetaStore struct {
	FileMetaMap    map[string]*FileMetaData
	BlockStoreAddr string
	UnimplementedMetaStoreServer
}

var metaLock sync.Mutex

func (m *MetaStore) GetFileInfoMap(ctx context.Context, _ *emptypb.Empty) (*FileInfoMap, error) {
	metaLock.Lock()
	defer metaLock.Unlock()
	return &FileInfoMap{
		FileInfoMap: m.FileMetaMap,
	}, nil
}

func (m *MetaStore) UpdateFile(ctx context.Context, fileMetaData *FileMetaData) (*Version, error) {
	metaLock.Lock()
	defer metaLock.Unlock()
	filename := fileMetaData.Filename
	if _, ok := m.FileMetaMap[filename]; ok {
		if fileMetaData.Version-m.FileMetaMap[filename].Version == 1 {
			m.FileMetaMap[filename] = fileMetaData
			return &Version{
				Version: fileMetaData.Version,
			}, nil
		} else {
			return &Version{
				Version: -1,
			}, nil //errors.New("Version number should be exactly one greater than the current!")
		}
	} else {
		m.FileMetaMap[filename] = fileMetaData
		return &Version{
			Version: fileMetaData.Version,
		}, nil
	}
}

func (m *MetaStore) GetBlockStoreAddr(ctx context.Context, _ *emptypb.Empty) (*BlockStoreAddr, error) {
	metaLock.Lock()
	defer metaLock.Unlock()
	return &BlockStoreAddr{
		Addr: m.BlockStoreAddr,
	}, nil
}

// This line guarantees all method for MetaStore are implemented
var _ MetaStoreInterface = new(MetaStore)

func NewMetaStore(blockStoreAddr string) *MetaStore {
	return &MetaStore{
		FileMetaMap:    map[string]*FileMetaData{},
		BlockStoreAddr: blockStoreAddr,
	}
}
