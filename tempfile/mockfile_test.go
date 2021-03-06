package tempfile_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/lanrat/extsort/tempfile"
)

func TestSingleMockFile(t *testing.T) {
	line := "The quick brown fox jumps over the lazy dog"
	tempWriter := tempfile.Mock(0)

	n, err := tempWriter.WriteString(line)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(line) {
		t.Fatalf("WriteString returned %d, expected %d", n, len(line))
	}
	s := tempWriter.Size()
	if s != 1 {
		t.Fatalf("tempWriter.Size returned %d, expected %d", s, 1)
	}

	tempReader, err := tempWriter.Save()
	if err != nil {
		t.Fatal(err)
	}
	s = tempReader.Size()
	if s != 1 {
		t.Fatalf("tempReader.Size returned %d, expected %d", s, 1)
	}
	str, err := tempReader.Read(0).ReadString('\n')
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if str != line {
		t.Fatalf("tempReader.ReadString returned %q expected %q", str, line)
	}
	err = tempReader.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTempMockRepeat(t *testing.T) {
	iterations := 10
	line := "The quick brown fox jumps over the lazy dog"
	tempWriter := tempfile.Mock(10)

	for i := 0; i < iterations; i++ {
		_, err := tempWriter.WriteString(fmt.Sprintf("%d: %s", i, line))
		if err != nil {
			t.Fatal(err)
		}
		s := tempWriter.Size()
		if s != i+1 {
			t.Fatalf("tempWriter.Size returned %d, expected %d", s, i+1)
		}
		_, err = tempWriter.Next()
		if err != nil {
			t.Fatal(err)
		}
	}

	tempReader, err := tempWriter.Save()
	if err != nil {
		t.Fatal(err)
	}

	s := tempReader.Size()
	if s != iterations+1 {
		t.Fatalf("tempReader.Size returned %d, expected %d", s, iterations)
	}

	for i := iterations - 1; i >= 0; i-- {
		str, err := tempReader.Read(i).ReadString('\n')
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("%d: %s", i, line)
		if str != expected {
			t.Fatalf("tempReader.ReadString %d returned %q expected %q", i, str, expected)
		}
	}
	err = tempReader.Close()
	if err != nil {
		t.Fatal(err)
	}
}
