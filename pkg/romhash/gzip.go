package romhash

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func decodeGZip(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening gzip file %q: %w", path, err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("error reading gzip file %q: %w", path, err)
	}
	defer gzr.Close()

	ext := filepath.Ext(gzr.Header.Name)
	if decoder, ok := getDecoder(ext); ok {
		d, err := io.ReadAll(gzr)
		if err != nil {
			return nil, fmt.Errorf("error decoding gzip file %q: %w", path, err)
		}

		r := bytes.NewReader(d)

		return decoder(io.NopCloser(r), int64(r.Len()))
	}

	return nil, fmt.Errorf("no roms found in gzip file %q: %w", path, ErrNoRoms)
}
