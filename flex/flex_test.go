package flex

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenFileForReadingWithStdin(t *testing.T) {
	file, err := OpenFileForReading("-")
	assert.NoError(t, err)
	assert.Equal(t, os.Stdin, file)
}

func TestOpenFileForReadingExpectError(t *testing.T) {
	fileName := "/invalid/path/to/file.json"
	file, err := OpenFileForReading(fileName)
	assert.Nil(t, file)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			fmt.Sprintf("open %s: no such file or directory", fileName),
			err.Error(),
		)
	}
}

func TestOpenFileForReading(t *testing.T) {
	file, err := os.CreateTemp("", "flextime")
	if err != nil {
		t.Errorf("%v", err)
	}
	defer os.Remove(file.Name())
	file.Close()
	file2, err := OpenFileForReading(file.Name())
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

func TestOpenFileForWritingWithStdout(t *testing.T) {
	file, err := OpenFileForWriting("-")
	assert.NoError(t, err)
	assert.Equal(t, os.Stdout, file)
}

func TestOpenFileForWritingExpectError(t *testing.T) {
	fileName := "/invalid/path/to/file.json"
	file, err := OpenFileForWriting(fileName)
	assert.Nil(t, file)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			fmt.Sprintf("open %s: no such file or directory", fileName),
			err.Error(),
		)
	}
}

func TestOpenFileForWriting(t *testing.T) {
	file, err := os.CreateTemp("", "flextime")
	if err != nil {
		t.Errorf("%v", err)
	}
	defer os.Remove(file.Name())
	file.Close()
	file2, err := OpenFileForWriting(file.Name())
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

func TestFlexDBFromFileWithValidJSONThatIsNotFlexDB(t *testing.T) {
	file, err := os.CreateTemp("", "flextime")
	if err != nil {
		t.Errorf("%v", err)
	}
	jsonDB := `{"something": 1}`
	_, err = file.Write([]byte(jsonDB))
	if err != nil {
		t.Errorf("%v", err)
	}
	ret, err := file.Seek(0, 0)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Return from file.Seek(): %d", ret)
	defer os.Remove(file.Name())

	db, err := FlexDBFromFile(file)
	assert.Nil(t, db)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrEmptyDB)
	}

	if err = file.Close(); err != nil {
		t.Errorf("%v", err)
	}
}

func TestFlexDBFromFile(t *testing.T) {
	file, err := os.CreateTemp("", "flextime")
	if err != nil {
		t.Errorf("%v", err)
	}
	jsonDB := []byte(`{"filename":"flexdb.json","customers":[{"customer_name":"Customer1","flex_entries":[{"date":"2021-12-01T09:25:05.3421539+01:00","amount":1}]}]}`)
	n, err := file.Write(jsonDB)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%d bytes written to %q", n, file.Name())
	ret, err := file.Seek(0, 0)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Return from file.Seek(): %d", ret)
	defer os.Remove(file.Name())

	db, err := FlexDBFromFile(file)
	assert.NoError(t, err)
	if assert.NotNil(t, db) {
		assert.IsType(
			t,
			(*DB)(nil),
			db,
		)
		assert.Equal(
			t,
			"",
			db.FileName,
		)
		assert.Equal(
			t,
			"Customer1",
			db.Customers[0].Name,
		)
	}
}

func TestFlexDBFromFileStdInWithInvalidData(t *testing.T) {
	input := []byte(`{}`) // valid json, but not valid FlexDB
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = writer.Write(input)
	if err != nil {
		t.Error(err)
	}
	writer.Close()

	stdin := os.Stdin
	defer func() { os.Stdin = stdin }()
	os.Stdin = reader

	db, err := FlexDBFromFile(os.Stdin)
	reader.Close()

	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestFlexDBToFileStdOut(t *testing.T) {
	db := &DB{
		FileName:  "flex.json",
		Customers: make(Customers, 0),
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = writer

	err = FlexDBToFile(db, os.Stdout)
	writer.Close()
	assert.NoError(t, err)

	bytes, err := io.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Got %d bytes", len(bytes))
	t.Logf("%s", string(bytes))

	reader.Close()
}

func TestFlexDBToFileExpectError(t *testing.T) {
	file, err := os.CreateTemp("", "flexytime")
	if err != nil {
		t.Error(err)
	}
	file.Close()
	defer os.Remove(file.Name())
	err = FlexDBToFile(&DB{}, file)
	if assert.Error(t, err) {
		t.Logf("%s", err.Error())
	}
}

func TestFlexDBToFile(t *testing.T) {
	db := &DB{
		FileName:  "flexdb.json",
		Customers: make(Customers, 0),
	}
	err := db.SetFlexForCustomer("Customer1", time.Now(), 1*time.Nanosecond, true)
	if err != nil {
		t.Error(err)
	}
	//customer, err := db.getCustomer("customer1")
	//if err != nil {
	//	t.Error(err)
	//}
	//customer.Print(os.Stdout, " ", 1)
	file, err := os.CreateTemp("", "flextime")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())
	err = FlexDBToFile(db, file)
	if err != nil {
		t.Error(err)
	}
	file.Close()
}

func TestDecodeFlexDBWithInvalidJSON(t *testing.T) {
	json := `{ blah`
	reader := strings.NewReader(json)
	db, err := DecodeFlexDB(reader)
	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestEncodeFlexDB(t *testing.T) {
	today := time.Now()
	jsonDate, err := json.Marshal(today)
	if err != nil {
		t.Errorf("Unable to generate JSON date: %v", err)
	}
	fe := Entry{Date: today, Amount: 1 * time.Nanosecond}
	c1 := Customer{Name: "Customer1", FlexEntries: Entries{&fe}}
	c2 := Customer{Name: "Customer2", FlexEntries: Entries{&fe}}
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
	err = EncodeFlexDB(db, &builder)
	assert.NoError(t, err)

	assert.Equal(
		t,
		expected,
		builder.String(),
	)
}

func TestDecodeFlexDB(t *testing.T) {
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
	db, err := DecodeFlexDB(reader)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"Customer1",
		db.Customers[0].Name,
	)
	assert.True(t, today.Equal(db.Customers[0].FlexEntries[0].Date))
	assert.Equal(
		t,
		time.Duration(1),
		db.Customers[0].FlexEntries[0].Amount,
	)
	assert.Equal(
		t,
		"Customer2",
		db.Customers[1].Name,
	)
	assert.True(t, today.Equal(db.Customers[1].FlexEntries[0].Date))
	assert.Equal(
		t,
		time.Duration(1),
		db.Customers[1].FlexEntries[0].Amount,
	)
}
