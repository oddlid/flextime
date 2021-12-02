package flex

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidJSONInput = errors.New("invalid JSON input")
	ErrEmptyDB          = errors.New("empty flex Database")
)

func NewDB() *DB {
	return &DB{
		Customers: make(Customers, 0),
	}
}

func GetFileOrStdinForReading(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdin, nil
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetFileOrStdoutForWriting(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdout, nil
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func EncodeDB(db *DB, w io.Writer) error {
	return json.NewEncoder(w).Encode(db)
}

func DecodeDB(r io.Reader) (*DB, error) {
	db := &DB{}
	err := json.NewDecoder(r).Decode(db)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if db.IsEmpty() {
		return nil, ErrEmptyDB
	}
	return db, nil
}
