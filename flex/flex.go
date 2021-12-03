package flex

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidJSONInput = errors.New("invalid JSON input")
	ErrEmptyDB          = errors.New("empty flex database")
)

// NewDB initializes and returns a new, empty DB instance
func NewDB() *DB {
	return &DB{
		Customers: make(Customers, 0),
	}
}

// GetFileOrStdinForReading returns os.Stdin, nil if fileName is "-",
// otherwise it will try to open the given fileName and return the file
// opened for reading.
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

// GetFileOrStdoutForWriting returns os.Stdout, nil if fileName is "-",
// otherwise it will try to create the given fileName and return the file
// opened for writing.
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

// EncodeDB encodes the given DB as JSON to the given writer
func EncodeDB(db *DB, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(db)
}

// DecodeDB tries to decode JSON input from the given reader
// into a new DB instance
func DecodeDB(reader io.Reader) (*DB, error) {
	db := &DB{}
	err := json.NewDecoder(reader).Decode(db)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if db.IsEmpty() {
		return nil, ErrEmptyDB
	}
	return db, nil
}
