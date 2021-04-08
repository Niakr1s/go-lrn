package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

const bucketName = "redirects"

func main() {
	path := flag.String("path", "", "path to redirect db")
	add := flag.String("add", "", "add redirect string in format 'url,redirectTo'")
	del := flag.String("del", "", "del redirect url")
	flag.Parse()

	if *path == "" {
		log.Fatalf("no path provided")
	}
	db, err := bolt.Open(*path, 0600, nil)
	if err != nil {
		log.Fatalf("error while opening bolt db: %v", err)
	}
	defer db.Close()

	if *del != "" {
		err := deleteRedirect(db, *del)
		if err != nil {
			log.Printf("couldn't delete redirect, reason: %v", err)
		}
		return
	}

	if *add != "" {
		splitted := strings.Split(*add, ",")
		if len(splitted) != 2 {
			log.Fatal("cannot add, invalid format")
		}

		err := addRedirect(db, splitted[0], splitted[1])
		if err != nil {
			log.Printf("couldn't add redirect, reason: %v", err)
		}
		return
	}

	printAllRedirects(db)
}

func printAllRedirects(db *bolt.DB) {
	db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			fmt.Printf("%s: %s\n", string(k), string(v))
			return nil
		})
	})
}

func deleteRedirect(db *bolt.DB, url string) error {
	return db.Update(func(t *bolt.Tx) error {
		bucket, err := t.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		err = bucket.Delete([]byte(url))

		if err != nil {
			return err
		}

		return err
	})
}

func addRedirect(db *bolt.DB, url, redirectTo string) error {
	return db.Update(func(t *bolt.Tx) error {
		bucket, err := t.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(url), []byte(redirectTo))

		if err != nil {
			return err
		}

		return err
	})
}
