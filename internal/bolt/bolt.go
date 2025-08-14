package bolt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dyammarcano/clonr/internal/model"
	"github.com/segmentio/ksuid"
	"go.etcd.io/bbolt"
)

type Bolt struct {
	*bbolt.DB
}

func InitBolt(path, id string) (*Bolt, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	repoMetrics := filepath.Join(path, fmt.Sprintf("%s.bolt", id))

	db, err := bbolt.Open(repoMetrics, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Bolt{db}, nil
}

func (b *Bolt) Close() error {
	return b.DB.Close()
}

func (b *Bolt) SaveStats(stats *model.StatsData) error {
	err := b.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("stats"))
		if err != nil {
			return err
		}

		key := ksuid.New().String()
		return bucket.Put([]byte(key), stats.Bytes())
	})

	return err
}
