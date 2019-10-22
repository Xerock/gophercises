package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var tasksBucket = []byte("tasks")
var db *bolt.DB

// Task represent a simple task
type Task struct {
	Key   int
	Value string
}

// Init the database connexion and create a bucket for the task if doesn't exist
func Init(dbpath string) error {
	// Open the database
	var err error
	db, err = bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	// Create the TasksBucket if it doesn't exist
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(tasksBucket)
		return err
	})

}

// AddTask to the tasks list
func AddTask(s string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		id, _ := b.NextSequence()
		err := b.Put(itob(int(id)), []byte(s))
		return err
	})
}

// DeleteTask removes the task corresponding to the id
func DeleteTask(id int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		return b.Delete(itob(id))
	})
}

// Tasks gives a list of all tasks
func Tasks() ([]Task, error) {
	var ret []Task

	return ret, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tasksBucket)
		b.ForEach(func(k, v []byte) error {
			{
				ret = append(ret, Task{
					Key:   btoi(k),
					Value: string(v),
				})
				return nil
			}
		})
		return nil
	})
}

// Close the database
func Close() {
	db.Close()
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	i := binary.BigEndian.Uint64(b)
	return int(i)
}
