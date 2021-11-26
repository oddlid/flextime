package flex

import (
	"bufio"
	"os"

	"github.com/rs/zerolog/log"
)

func NewFlexDB(fileName string) *FlexDB {
	db, err := FlexDBFromFile(fileName)
	if err != nil {
		log.Error().Err(err).Send()
		log.Info().Str("filename", fileName).Msg("Returning new, empty FlexDB")
		return &FlexDB{
			FileName:  fileName,
			Customers: make(Customers, 0),
		}
	}
	// We don't want the filename saved in the JSON file, since that could lead to problems,
	// say if you rename the file, but don't edit the file and change the filename field.
	// Then it would be saved back to whatever was set in the file, and not the same name
	// as it was loaded from.
	db.FileName = fileName
	return db
}

func FlexDBFromFile(fileName string) (*FlexDB, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	db := &FlexDB{}
	err = db.Decode(reader)
	if err != nil {
		return nil, err
	}
	return db, nil
}
