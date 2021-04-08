package redirects

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const boltBucket = "redirects"

func boltRedirects(boltPath string) (Redirects, error) {
	db, err := bolt.Open(boltPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't open bolt db: %v", err)
	}
	defer db.Close()

	err = db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists([]byte(boltBucket))
		return err
	})

	if err != nil {
		return nil, err
	}

	redirects := map[string]string{}
	err = db.View(func(t *bolt.Tx) error {
		// bucket is not nil because we created it before
		bucket := t.Bucket([]byte(boltBucket))

		err = bucket.ForEach(func(k, v []byte) error {
			redirects[string(k)] = string(v)
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return redirects, nil
}
