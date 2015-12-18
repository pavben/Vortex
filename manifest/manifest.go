package manifest

import (
	"crypto/sha1"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ManifestEntity interface {
	Id() uint32
	Name() string
}

type Manifest struct {
	rootEntity ManifestEntity
	// Generated entityMap (index) keyed by ID
	entityMap map[uint32]ManifestEntity
}

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

type ManifestFile struct {
	id       uint32
	name     string
	fileSize uint64
	hashes   [][]byte
}

func (mf *ManifestFile) Id() uint32 {
	return mf.id
}

func (mf *ManifestFile) Name() string {
	return mf.name
}

func (mf *ManifestFile) Size() uint64 {
	return mf.fileSize
}

func (mf *ManifestFile) Hashes() [][]byte {
	return mf.hashes
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
		fileSize, hashes, err := getFileSizeAndHashes(currentPath)
		if err != nil {
			return nil, err
		}
		return indexAndReturn(entityMap, &ManifestFile{
			id:       takeId(nextId),
			name:     fileInfo.Name(),
			fileSize: fileSize,
			hashes:   hashes,
		}), nil
	}
}

func getFileSizeAndHashes(filePath string) (uint64, [][]byte, error) {
	const chunkSize uint32 = 4 * 1024 * 1024 // 4 MB
	f, err := os.Open(filePath)
	if err != nil {
		return 0, nil, err
	}
	defer f.Close()
	var bytesRead uint64 = 0
	var hashes [][]byte
	buf := make([]byte, chunkSize)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return 0, nil, err
		}
		if n <= 0 {
			break
		}
		bytesRead += uint64(n)
		hash := sha1.Sum(buf[:n])
		hashes = append(hashes, hash[:])
	}
	return bytesRead, hashes, nil
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
