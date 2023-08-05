package cache

import (
	"encoding/binary"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/mmcdole/gofeed"
	"github.com/spaolacci/murmur3"
)

// ReadStatus is a set containing the hashes of the already read articles. We use a struct{} here
// because it takes up no space in memory. To hash the article, we use the URL of the feed and the title
// of the article.
type ReadStatus struct {
	set      map[uint32]struct{}
	filePath string
}

// New creates a new ReadStatus set.
func NewReadStatus(dir string) (*ReadStatus, error) {
	log.Println("Creating new read status")
	if dir == "" {
		defaultDir, err := getDefaultDir()
		if err != nil {
			return nil, err
		}

		dir = defaultDir
	}

	return &ReadStatus{
		filePath: filepath.Join(dir, "read_status"),
		set:      make(map[uint32]struct{}),
	}, nil
}

// Load reads the cache from disk
func (rs *ReadStatus) Load() error {
	log.Println("Loading read status from", rs.filePath)
	if _, err := os.Stat(rs.filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	data, err := os.ReadFile(rs.filePath)
	if err != nil {
		return err
	}

	rs.set, err = unmarshal(data)
	return err
}

// Save writes the cache to disk
func (rs ReadStatus) Save() error {
	data := marshal(rs.set)
	log.Println("Marshalling the data yielded a size of", len(data))

	// Try to write the data to the file
	if err := os.WriteFile(rs.filePath, data, 0600); err != nil {
		if err = os.MkdirAll(filepath.Dir(rs.filePath), 0755); err != nil {
			return err
		}

		if err = os.WriteFile(rs.filePath, data, 0600); err != nil {
			return err
		}
	}

	log.Println("Written succesffully")
	return nil
}

// MarkAsRead adds an article to the set.
func (rs *ReadStatus) MarkAsRead(item gofeed.Item) {
	rs.set[hashArticle(item)] = struct{}{}
}

// IsRead checks if an article is already in the set.
func (rs ReadStatus) IsRead(item gofeed.Item) bool {
	_, ok := rs.set[hashArticle(item)]
	return ok
}

// MarkAsUnread removes an article from the set.
func (rs *ReadStatus) MarkAsUnread(item gofeed.Item) {
	delete(rs.set, hashArticle(item))
}

// marshal converts the set to bytes.
func marshal(set map[uint32]struct{}) []byte {
	result := make([]byte, 0, len(set)*4)
	for k := range set {
		result = binary.LittleEndian.AppendUint32(result, k)
	}
	return result
}

// unmarshal converts bytes to a set.
func unmarshal(data []byte) (map[uint32]struct{}, error) {
	set := make(map[uint32]struct{})
	if len(data)%4 != 0 {
		return nil, errors.New("invalid data")
	}

	for i := 0; i < len(data); i += 4 {
		set[binary.LittleEndian.Uint32(data[i:i+4])] = struct{}{}
	}

	return set, nil
}

// hashArticle hashes the gofeed.Item to a uint32.
func hashArticle(item gofeed.Item) uint32 {
	h := murmur3.New32()
	h.Write([]byte(item.Title))
	h.Write([]byte(item.Link))
	h.Write([]byte(item.GUID))
	return h.Sum32()
}
