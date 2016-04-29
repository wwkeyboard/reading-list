package reading

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

// Piece is a thing to be read
type Piece struct {
	Name   string
	URL    string
	Source string
}

// Database holds the underlying data
type Database struct {
	filename string
	db       *bolt.DB
}

// NewDatabase configures and returns a new piece DB
func NewDatabase(filename string) (*Database, error) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Database{
		filename,
		db,
	}, nil
}

var defaultBucket = []byte("default-bucket")

// List returns a all of the pieces for that bucket
func (d *Database) List(bucket string) ([]Piece, error) {
	var pieces []Piece
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		if b == nil {
			return nil
		}

		var piece Piece
		b.ForEach(func(k, v []byte) error {
			err := json.Unmarshal(v, &piece)
			if err != nil {
				log.Printf("Couldn't Unmarshal %+v", v)
			} else {
				pieces = append(pieces, piece)
			}
			return nil
		})

		return nil
	})

	return pieces, err
}

// AddPiece saves the piece to the default bucket
func (d *Database) AddPiece(piece *Piece) error {
	payload, err := json.Marshal(piece)
	if err != nil {
		return err
	}

	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		err = b.Put([]byte(piece.Name), payload)
		return err
	})
	return err
}

// EnsureBucket returns true if the bucket exists
func (d *Database) EnsureBucket() error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(defaultBucket))
		return err
	})
	return err
}
