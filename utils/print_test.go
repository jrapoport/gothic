package utils

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyString(t *testing.T) {
	t.Parallel()
	type Foo struct {
		Value  string
		Number int
	}
	f := Foo{
		Value:  "test",
		Number: 42,
	}
	s := PrettyString(f)
	assert.Equal(t, "{\n\t\"Value\": \"test\",\n\t\"Number\": 42\n}", s)
}

func TestPrettyPrint(t *testing.T) {
	t.Parallel()
	type Foo struct {
		Value  string
		Number int
	}
	f := Foo{
		Value:  "test",
		Number: 42,
	}
	assert.NotPanics(t, func() {
		PrettyPrint(f)
	})
}

func TestPrintGrid(t *testing.T) {
	t.Parallel()
	const (
		val   = "1234"
		grid1 = "1234\t1234\t1234\t\t\n"
		grid2 = "1234\t1234\t1234\t1234\t\n"
		grid3 = "1234\t1234\t1234\t1234\t\n1234\t1234\t1234\t1234\t\n"
		grid4 = "1234\t1234\t1234\t1234\t\n1234\t1234\t1234\t1234\t\n1234\t\t\t\t\n"
	)
	makeList := func(n int) []string {
		l := make([]string, n)
		for i := 0; i < n; i++ {
			l[i] = val
		}
		return l
	}
	tests := []struct {
		list []string
		grid string
	}{
		{makeList(3), grid1},
		{makeList(4), grid2},
		{makeList(8), grid3},
		{makeList(9), grid4},
	}
	for _, test := range tests {
		var b bytes.Buffer
		PrintGrid(&b, test.list, 4)
		assert.Equal(t, test.grid, b.String())
	}
	var b bytes.Buffer
	PrintGrid(&b, tests[0].list, 0)
	assert.Equal(t, "1234\t\n1234\t\n1234\t\n", b.String())
}

func TestWriteCSV(t *testing.T) {
	t.Parallel()
	const (
		val  = "1234"
		test = "value\n1234\n1234\n1234\n1234\n1234\n1234\n1234\n1234\n1234\n1" +
			"234\n1234\n1234\n1234\n1234\n1234\n1234\n1234\n1234\n1234\n1234\n"
	)
	makeList := func(n int) []string {
		l := make([]string, n)
		for i := 0; i < n; i++ {
			l[i] = val
		}
		return l
	}
	csv := filepath.Join(t.TempDir(), "test.csv")
	err := WriteCSV(csv, "value", makeList(20))
	assert.NoError(t, err)
	b, err := ioutil.ReadFile(csv)
	assert.NoError(t, err)
	assert.Equal(t, test, string(b))
}
