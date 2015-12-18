package manifest

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type ManifestEntity interface {
	Id() uint32
	Name() string
	String() string
}

type Manifest struct {
	rootEntity ManifestEntity
	// Generated entityMap (index) keyed by ID
	entityMap map[uint32]ManifestEntity
}

/*
func (m Manifest) String() string {
	return (m.rootEntity.(*ManifestFolder).contents[1]).(*ManifestFolder).contents[0].String()
}
*/

type ManifestFolder struct {
	id       uint32
	name     string
	contents []ManifestEntity
}

func (mf *ManifestFolder) Id() uint32 {
	return mf.id
}

func (mf *ManifestFolder) Name() string {
	return mf.name
}

func (mf *ManifestFolder) String() string {
	return mf.Name()
}

type ManifestFile struct {
	id         uint32
	name       string
	byteLength uint64
	pieces     [][]byte
}

func (mf *ManifestFile) Id() uint32 {
	return mf.id
}

func (mf *ManifestFile) Name() string {
	return mf.name
}

func (mf *ManifestFile) String() string {
	return mf.Name()
}

func GenerateManifestFromPath(p string) (*Manifest, error) {
	fileInfo, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	var nextId uint32 = 0
	entityMap := make(map[uint32]ManifestEntity)
	rootEntity, err := generateManifestEntityTree(p, fileInfo, entityMap, &nextId)
	if err != nil {
		return nil, err
	}
	return &Manifest{
		rootEntity: rootEntity,
		entityMap:  entityMap,
	}, nil
}

func generateManifestEntityTree(currentPath string, fileInfo os.FileInfo, entityMap map[uint32]ManifestEntity, nextId *uint32) (ManifestEntity, error) {
	if fileInfo.IsDir() {
		var contents []ManifestEntity
		childrenFileInfos, err := ioutil.ReadDir(currentPath)
		if err != nil {
			return nil, err
		}
		for _, fileInfo := range childrenFileInfos {
			childEntity, err := generateManifestEntityTree(filepath.Join(currentPath, fileInfo.Name()), fileInfo, entityMap, nextId)
			if err != nil {
				return nil, err
			}
			contents = append(contents, childEntity)
		}
		return indexAndReturn(entityMap, &ManifestFolder{
			id:       takeId(nextId),
			name:     fileInfo.Name(),
			contents: contents,
		}), nil
	} else {
		return indexAndReturn(entityMap, &ManifestFile{
			id:         takeId(nextId),
			name:       fileInfo.Name(),
			byteLength: 0,
			pieces:     make([][]byte, 0),
		}), nil
	}
}

func indexAndReturn(entityMap map[uint32]ManifestEntity, mf ManifestEntity) ManifestEntity {
	entityMap[mf.Id()] = mf
	return mf
}

func takeId(nextId *uint32) uint32 {
	id := *nextId
	*nextId = *nextId + 1
	return id
}
