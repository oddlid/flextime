package flex

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidJSONInput = errors.New("Invalid JSON input")
	ErrEmptyDB          = errors.New("Empty flex Database")
)

func NewFlexDB() *FlexDB {
	return &FlexDB{
		Customers: make(Customers, 0),
	}
}

func OpenFileForReading(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdin, nil
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func OpenFileForWriting(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdout, nil
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func EncodeFlexDB(db *FlexDB, w io.Writer) error {
	return json.NewEncoder(w).Encode(db)
}

func DecodeFlexDB(r io.Reader) (*FlexDB, error) {
	db := &FlexDB{}
	err := json.NewDecoder(r).Decode(db)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	//bytes, err := io.ReadAll(r)
	//if err != nil {
	//	log.Error().Err(err).Msg("Error from io.ReadAll()")
	//	return nil, err
	//}
	//err = json.Unmarshal(bytes, db)
	//if err != nil {
	//	log.Error().Err(err).Msg("Error from json.Unmarshal()")
	//	return nil, err
	//}

	if db.IsEmpty() {
		return nil, ErrEmptyDB
	}
	return db, nil
}

func FlexDBToFile(db *FlexDB, file *os.File) error {
	//writer := bufio.NewWriter(file)
	//if err := EncodeFlexDB(db, writer); err != nil {
	//	return err
	//}
	//writer.Flush()
	//return nil
	return EncodeFlexDB(db, file)
}

// FlexDBFromFile will try to decode the JSON from the file into a FlexDB struct pointer.
// It's important to call file.Seek(0, 0) before passing it to this function, if you've
// written to the file after opening it.
func FlexDBFromFile(file *os.File) (*FlexDB, error) {
	reader := bufio.NewReader(file)
	db, err := DecodeFlexDB(reader)
	if err != nil {
		return nil, err
	}
	return db, nil
}
