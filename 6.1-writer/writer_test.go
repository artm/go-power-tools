package writer_test

import (
	"os"
	"testing"
	"writer"

	"github.com/google/go-cmp/cmp"
)

func TestWriteToFile(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/write_test.txt"
	want := []byte{1, 2, 3}
	err := writer.WriteToFile(path, want)
	if err != nil {
		t.Fatal(err)
	}
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	perm := stat.Mode().Perm()
	if perm != 0600 {
		t.Errorf("want file mode 0600, got 0%o", perm)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestWriteToFileClobbers(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/clobber_test.txt"
	err := os.WriteFile(path, []byte{4, 5, 6}, 0600)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{1, 2, 3}
	err = writer.WriteToFile(path, want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestPermsClosed(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/perms_test.txt"
	// Pre-create empty file with open perms
	err := os.WriteFile(path, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.WriteToFile(path, []byte{})
	if err != nil {
		t.Fatal(err)
	}
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	perm := stat.Mode().Perm()
	if perm != 0600 {
		t.Errorf("want file mode 0600, got 0%o", perm)
	}
}

func TestWriteZerosWrites(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/write_zeros_test.dat"
	err := writer.WriteZeros(path, 3)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]byte, 3)
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestWriteZerosClobbers(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/zero_clobber_test.dat"
	err := os.WriteFile(path, []byte{4, 5, 6}, 0600)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]byte, 3)
	err = writer.WriteZeros(path, 3)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestZeroerWrites(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/zeroer_writes_test.dat"
	args := []string{"-size", "3", path}
	cli, err := writer.NewZeroer(
		writer.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = cli.Write()
	if err != nil {
		t.Fatal(err)
	}
	want := make([]byte, 3)
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestZeroerRetries(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/zeroer_retries_test.dat"
	args := []string{"-size", "3", "-retries", "3", path}
	cli, err := writer.NewZeroer(
		writer.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = cli.Write()
	if err != nil {
		t.Fatal(err)
	}
	want := make([]byte, 3)
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func benchmarkWriteZeros(count int, b *testing.B) {
	path := b.TempDir() + "/write_zeros_test.dat"
	for i := 0; i < b.N; i++ {
		writer.WriteZeros(path, count)
	}
}

func BenchmarkWriteZeros5(b *testing.B) {
	benchmarkWriteZeros(5, b)
}

func BenchmarkWriteZeros5k(b *testing.B) {
	benchmarkWriteZeros(5000, b)
}

func BenchmarkWriteZeros5M(b *testing.B) {
	benchmarkWriteZeros(5000000, b)
}
