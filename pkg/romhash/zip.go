package romhash

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
)

type zipReader struct {
	path *zip.ReadCloser
	rom  io.ReadCloser
}

func (r zipReader) Read(data []byte) (int, error) {
	readCount, err := r.rom.Read(data)
	if err != nil {
		return readCount, fmt.Errorf("error reading zip file: %w", err)
	}

	return readCount, nil
}

func (r zipReader) Close() error {
	r.rom.Close()

	if err := r.path.Close(); err != nil {
		return fmt.Errorf("error closing zip file: %w", err)
	}

	return nil
}

func decodeZip(path string) (io.ReadCloser, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read zip file %q: %w", path, err)
	}

	for _, zipFile := range reader.File {
		ext := filepath.Ext(zipFile.FileHeader.Name)
		if decoder, ok := getDecoder(ext); ok {
			romFile, err := zipFile.Open()
			if err != nil {
				continue
			}

			rs := zipFile.FileHeader.UncompressedSize64

			rom, err := decoder(romFile, int64(rs))
			if err != nil {
				continue
			}

			return zipReader{reader, rom}, nil
		}
	}

	return nil, fmt.Errorf("no valid roms found in %q: %w", path, ErrNoRoms)
}
