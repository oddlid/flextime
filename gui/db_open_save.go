package main

import (
	"errors"
	"fmt"

	"github.com/oddlid/flextime/flex"
)

func openDB(fileName string) (*flex.DB, error) {
	if fileName == "" {
		db := flex.NewDB()
		db.FileName = "-"
		return db, nil
	}

	file, err := flex.GetFileOrStdinForReading(fileName)
	if err != nil {
		db := flex.NewDB()
		db.FileName = fileName
		return db, err
	}

	db, err := flex.DecodeDB(file)
	if err != nil {
		if errors.Is(err, flex.ErrEmptyDB) {
			db = flex.NewDB()
		} else {
			return nil, err
		}
	}
	db.FileName = fileName
	if fileName != "-" {
		err = file.Close()
	}

	return db, err
}

func saveDB(db *flex.DB) error {
	if db == nil {
		return fmt.Errorf("refusing to save nil DB")
	}
	file, err := flex.GetFileOrStdoutForWriting(db.FileName)
	if err != nil {
		return err
	}
	err = flex.EncodeDB(db, file)
	if err != nil {
		return err
	}
	if db.FileName != "-" {
		if err = file.Close(); err != nil {
			return err
		}
	}
	return nil
}
