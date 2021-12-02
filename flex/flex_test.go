package flex

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFileOrStdinForReadingExpectStdin(t *testing.T) {
	file, err := GetFileOrStdinForReading("-")
	assert.NoError(t, err)
	assert.Equal(t, os.Stdin, file)
}

func TestGetFileOrStdinForReadingExpectError(t *testing.T) {
	fileName := "/invalid/path/to/file.json"
	file, err := GetFileOrStdinForReading(fileName)
	assert.Nil(t, file)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			fmt.Sprintf("open %s: no such file or directory", fileName),
			err.Error(),
		)
	}
}

func TestGetFileOrStdinForReading(t *testing.T) {
	file, err := os.CreateTemp("", "flextime")
	if err != nil {
		t.Errorf("%v", err)
	}
	defer os.Remove(file.Name())
	file.Close()
	file2, err := GetFileOrStdinForReading(file.Name())
	assert.NoError(t, err)
	if assert.NotNil(t, file2) {
		assert.Equal(t, file.Name(), file2.Name())
		assert.IsType(
			t,
			(*os.File)(nil),
			file2,
		)
		file2.Close()
	}
}

func TestGetFileOrStdoutForWritingExpectStdout(t *testing.T) {
	file, err := GetFileOrStdoutForWriting("-")
	assert.NoError(t, err)
	assert.Equal(t, os.Stdout, file)
}

func TestGetFileOrStdoutForWritingExpectError(t *testing.T) {
	fileName := "/invalid/path/to/file.json"
	file, err := GetFileOrStdoutForWriting(fileName)
	assert.Nil(t, file)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			fmt.Sprintf("open %s: no such file or directory", fileName),
			err.Error(),
		)
	}
}

func TestGetFileOrStdoutForWriting(t *testing.T) {
	file, err := os.CreateTemp("", "flextime")
	if err != nil {
		t.Errorf("%v", err)
	}
	defer os.Remove(file.Name())
	file.Close()
	file2, err := GetFileOrStdoutForWriting(file.Name())
	assert.NoError(t, err)
	if assert.NotNil(t, file2) {
		assert.Equal(t, file.Name(), file2.Name())
		assert.IsType(
			t,
			(*os.File)(nil),
			file2,
		)
		file2.Close()
	}
}

func TestEncodeDB(t *testing.T) {
	today := time.Now()
	jsonDate, err := json.Marshal(today)
	if err != nil {
		t.Errorf("Unable to generate JSON date: %v", err)
	}
	entry := Entry{Date: today, Amount: 1 * time.Nanosecond}
	c1 := Customer{Name: "Customer1", Entries: Entries{&entry}}
	c2 := Customer{Name: "Customer2", Entries: Entries{&entry}}
	db := &DB{FileName: "flex.json", Customers: Customers{&c1, &c2}}

	expected := fmt.Sprintf(
		"{%q:[{%q:%q,%q:[{%q:%s,%q:1}]},{%q:%q,%q:[{%q:%s,%q:1}]}]}\n",
		"customers",
		"customer_name",
		"Customer1",
		"flex_entries",
		"date",
		string(jsonDate),
		"amount",
		"customer_name",
		"Customer2",
		"flex_entries",
		"date",
		string(jsonDate),
		"amount",
	)
	t.Logf("Expected JSON: %s", expected)

	builder := strings.Builder{}
	err = EncodeDB(db, &builder)
	assert.NoError(t, err)

	assert.Equal(
		t,
		expected,
		builder.String(),
	)
}

func TestDecodeDBWithInvalidJSON(t *testing.T) {
	json := `{ blah`
	reader := strings.NewReader(json)
	db, err := DecodeDB(reader)
	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestDecodeDBExpectErrEmptyDB(t *testing.T) {
	reader := strings.NewReader("{}")
	db, err := DecodeDB(reader)
	assert.Nil(t, db)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrEmptyDB)
	}
}

func TestDecodeDB(t *testing.T) {
	today := time.Now()
	jsonDate, err := json.Marshal(today)
	if err != nil {
		t.Errorf("Unable to generate JSON date: %v", err)
	}
	jsonInput := fmt.Sprintf(
		"{%q:[{%q:%q,%q:[{%q:%s,%q:1}]},{%q:%q,%q:[{%q:%s,%q:1}]}]}\n",
		"customers",
		"customer_name",
		"Customer1",
		"flex_entries",
		"date",
		string(jsonDate),
		"amount",
		"customer_name",
		"Customer2",
		"flex_entries",
		"date",
		string(jsonDate),
		"amount",
	)
	reader := strings.NewReader(jsonInput)
	db, err := DecodeDB(reader)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"Customer1",
		db.Customers[0].Name,
	)
	assert.True(t, today.Equal(db.Customers[0].Entries[0].Date))
	assert.Equal(
		t,
		time.Duration(1),
		db.Customers[0].Entries[0].Amount,
	)
	assert.Equal(
		t,
		"Customer2",
		db.Customers[1].Name,
	)
	assert.True(t, today.Equal(db.Customers[1].Entries[0].Date))
	assert.Equal(
		t,
		time.Duration(1),
		db.Customers[1].Entries[0].Amount,
	)
}
