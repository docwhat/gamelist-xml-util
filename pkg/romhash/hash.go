//nolint:gomnd
package romhash

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"strings"
)

var (
	ErrNoRoms     = fmt.Errorf("no valid roms found")
	ErrInvalidRom = fmt.Errorf("invalid rom")
	ErrUnknownExt = fmt.Errorf("unknown rom extension")
)

// var extra map[string]bool //nolint:gochecknoglobals

// func init() {
// 	extra = make(map[string]bool)
// }

// func AddExtra(e ...string) {
// 	for _, x := range e {
// 		extra[strings.ToLower(x)] = true
// 	}
// }

// func DelExtra(e ...string) {
// 	for _, x := range e {
// 		delete(extra, strings.ToLower(x))
// 	}
// }

// func HasExtra(e string) bool {
// 	return extra[e]
// }

// func ClearExtra() {
// 	extra = make(map[string]bool)
// }

type decoder func(io.ReadCloser, int64) (io.ReadCloser, error)

// Noop does nothong but return the passed in file.
func noop(f io.ReadCloser, s int64) (io.ReadCloser, error) {
	return f, nil
}

func decodeLNX(reader io.ReadCloser, s int64) (io.ReadCloser, error) {
	if s < 4 {
		return nil, ErrInvalidRom
	}

	tmp := make([]byte, 64)
	_, err := io.ReadFull(reader, tmp)

	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if errors.Is(err, io.ErrUnexpectedEOF) || !bytes.Equal(tmp[:4], []byte("LYNX")) {
		return newMultiReader(tmp, reader), nil
	}

	return reader, nil
}

type multireader struct {
	r  io.ReadCloser
	mr io.Reader
}

func (mr *multireader) Read(b []byte) (int, error) {
	//nolint:wrapcheck
	return mr.mr.Read(b)
}

func (mr *multireader) Close() error {
	//nolint:wrapcheck
	return mr.r.Close()
}

func newMultiReader(b []byte, r io.ReadCloser) io.ReadCloser {
	return &multireader{r, io.MultiReader(bytes.NewReader(b), r)}
}

func decodeA78(reader io.ReadCloser, s int64) (io.ReadCloser, error) {
	tmpBuffer := make([]byte, 128)

	_, err := io.ReadFull(reader, tmpBuffer)

	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, fmt.Errorf("error reading A78 rom: %w", err)
	}

	if errors.Is(err, io.ErrUnexpectedEOF) || !bytes.Equal(tmpBuffer[1:10], []byte("ATARI7800")) {
		return newMultiReader(tmpBuffer, reader), nil
	}

	return reader, nil
}

func deinterleave(interleavedBuffer []byte) []byte {
	interleavedBufferLength := len(interleavedBuffer)
	midpointIndex := interleavedBufferLength / 2
	deinterleavedBuffer := make([]byte, interleavedBufferLength)

	for interleavedIndex, dataByte := range interleavedBuffer {
		if interleavedIndex < midpointIndex {
			deinterleavedBuffer[interleavedIndex*2+1] = dataByte
		} else {
			deinterleavedBuffer[interleavedIndex*2-interleavedBufferLength] = dataByte
		}
	}

	return deinterleavedBuffer
}

func decodeSMD(f io.ReadCloser, s int64) (io.ReadCloser, error) {
	return decodeMD(f, s, ".smd")
}

func decodeMGD(f io.ReadCloser, s int64) (io.ReadCloser, error) {
	return decodeMD(f, s, ".mgd")
}

func decodeGEN(f io.ReadCloser, s int64) (io.ReadCloser, error) {
	return decodeMD(f, s, ".gen")
}

func decodeMD(reader io.ReadCloser, size int64, extension string) (io.ReadCloser, error) {
	// Strip off 512 byte header if present.
	if size%16384 == 512 {
		tmp := make([]byte, 512)
		if _, err := io.ReadFull(reader, tmp); err != nil {
			return nil, fmt.Errorf("error reading MD: %w", err)
		}

		size -= 512
	}

	if size%16384 != 0 {
		return nil, fmt.Errorf("invalid MD size: %w", ErrInvalidRom)
	}

	readBuffer, err := io.ReadAll(reader)
	reader.Close()

	if err != nil {
		return nil, fmt.Errorf("error opening MD for read: %w", err)
	}

	if bytes.Equal(readBuffer[256:260], []byte("SEGA")) {
		return io.NopCloser(bytes.NewReader(readBuffer)), nil
	}

	if bytes.Equal(readBuffer[8320:8328], []byte("SG EEI  ")) || bytes.Equal(readBuffer[8320:8328], []byte("SG EADIE")) {
		for i := 0; int64(i) < (size / int64(16384)); i++ {
			x := i * 16384
			copy(readBuffer[x:x+16384], deinterleave(readBuffer[x:x+16384]))
		}

		return io.NopCloser(bytes.NewReader(readBuffer)), nil
	}

	if bytes.Equal(readBuffer[128:135], []byte("EAGNSS ")) || bytes.Equal(readBuffer[128:135], []byte("EAMG RV")) {
		readBuffer = deinterleave(readBuffer)

		return io.NopCloser(bytes.NewReader(readBuffer)), nil
	}

	switch extension {
	case ".smd":
		for i := 0; int64(i) < (size / int64(16384)); i++ {
			x := i * 16384
			copy(readBuffer[x:x+16384], deinterleave(readBuffer[x:x+16384]))
		}

		return io.NopCloser(bytes.NewReader(readBuffer)), nil
	case ".mgd":
		readBuffer = deinterleave(readBuffer)

		return io.NopCloser(bytes.NewReader(readBuffer)), nil
	case ".gen":
		return io.NopCloser(bytes.NewReader(readBuffer)), nil
	}

	return nil, ErrUnknownExt
}

func decodeN64(f io.ReadCloser, s int64) (io.ReadCloser, error) {
	i := 0
	swapper := swapReader{f, make([]byte, 4), &i, noSwap}

	if s < 4 {
		return nil, ErrInvalidRom
	}

	_, err := io.ReadFull(swapper.reader, swapper.buffer)
	*swapper.length = 4

	if err != nil {
		return nil, fmt.Errorf("error reading N64 rom: %w", err)
	}

	switch {
	case swapper.buffer[0] == 0x80:
		swapper.s = zSwap
	case swapper.buffer[3] == 0x80:
		swapper.s = nSwap
	}

	return swapper, nil
}

func noSwap(b []byte) {}

func zSwap(buffer []byte) {
	bufferLength := len(buffer)
	for idx := 0; idx < bufferLength; idx += 4 {
		if bufferLength-idx < 4 {
			continue
		}

		buffer[idx+1], buffer[idx] = buffer[idx], buffer[idx+1]
		buffer[idx+3], buffer[idx+2] = buffer[idx+2], buffer[idx+3]
	}
}

func nSwap(buffer []byte) {
	bufferLength := len(buffer)
	for idx := 0; idx < bufferLength; idx += 4 {
		if bufferLength-idx < 4 {
			continue
		}

		buffer[idx+2], buffer[idx] = buffer[idx], buffer[idx+2]
		buffer[idx+3], buffer[idx+1] = buffer[idx+1], buffer[idx+3]
	}
}

type swapReader struct {
	reader io.ReadCloser
	buffer []byte
	length *int
	s      func([]byte)
}

func (r swapReader) Read(readBuffer []byte) (int, error) {
	readBufferLength := len(readBuffer)
	rl := readBufferLength - *r.length
	tmpBufLen := rl + 4 - 1 - (rl-1)%4

	copy(readBuffer, r.buffer[:*r.length])

	if rl <= 0 {
		*r.length -= readBufferLength
		copy(r.buffer, r.buffer[readBufferLength:])

		return readBufferLength, nil
	}

	numberOfBytes := *r.length
	tmpBuf := make([]byte, tmpBufLen)
	tmpBufReadCount, err := io.ReadFull(r.reader, tmpBuf)

	if tmpBufReadCount == 0 || err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return numberOfBytes, fmt.Errorf("error reading: %w", err)
	}

	copy(readBuffer[numberOfBytes:readBufferLength], tmpBuf[:tmpBufReadCount])
	numberOfBytes += tmpBufReadCount

	if readBufferLength <= numberOfBytes {
		r.s(readBuffer)
		copy(r.buffer, tmpBuf[tmpBufReadCount+readBufferLength-numberOfBytes:tmpBufReadCount])
		*r.length = numberOfBytes - readBufferLength

		return readBufferLength, nil
	}

	r.s(readBuffer[:numberOfBytes])
	*r.length = 0

	return numberOfBytes, nil
}

func (r swapReader) Close() error {
	//nolint:wrapcheck
	return r.reader.Close()
}

type nesReader struct {
	f      io.ReadCloser
	offset int64
	size   int64
	start  int64
	end    int64
}

func (r *nesReader) Read(p []byte) (int, error) {
	readCount, err := r.f.Read(p)
	if err != nil {
		return readCount, fmt.Errorf("error reading rom: %w", err)
	}

	if r.offset+int64(readCount) > r.end {
		readCount = int(r.end - r.offset)
	}

	r.offset += int64(readCount)

	return readCount, nil
}

func (r *nesReader) Close() error {
	//nolint:wrapcheck
	return r.f.Close()
}

func decodeNES(reader io.ReadCloser, s int64) (io.ReadCloser, error) {
	header := make([]byte, 16)

	readCount, err := io.ReadFull(reader, header)
	if err != nil {
		return nil, fmt.Errorf("error reading header from NES rom: %w", err)
	}

	if readCount < 16 {
		return nil, ErrInvalidRom
	}

	if !bytes.Equal(header[:3], []byte("NES")) {
		return newMultiReader(header, reader), nil
	}

	prgSize := int64(header[4])
	chrSize := int64(header[5])

	if header[7]&12 == 8 {
		romSize := int64(header[9])
		chrSize = romSize&0x0F<<8 + chrSize
		prgSize = romSize&0xF0<<4 + prgSize
	}

	prg := 16 * 1024 * prgSize
	chr := 8 * 1024 * chrSize
	hasTrainer := header[6]&4 == 4
	offset := int64(16)

	if hasTrainer {
		tmp := make([]byte, 512)

		n, err := io.ReadFull(reader, tmp)
		if err != nil {
			return nil, fmt.Errorf("error reading trainer from NES rom: %w", err)
		}

		offset += int64(n)
	}

	return &nesReader{reader, offset, prg + chr, offset, offset + prg + chr}, nil
}

func decodeSNES(reader io.ReadCloser, s int64) (io.ReadCloser, error) {
	if s%1024 == 512 {
		tmp := make([]byte, 512)

		_, err := io.ReadFull(reader, tmp)
		if err != nil {
			return nil, fmt.Errorf("error reading header from SNES rom: %w", err)
		}
	}

	return reader, nil
}

func getDecoder(ext string) (decoder, bool) {
	ext = strings.ToLower(ext)
	switch ext {
	case ".bin", ".a26", ".a52", ".rom", ".cue", ".gdi", ".gb", ".gba", ".gbc", ".32x", ".gg",
		".pce", ".sms", ".col", ".ngp", ".ngc", ".sg", ".int", ".vb", ".vec", ".gam", ".j64",
		".jag", ".mgw", ".nds", ".fds", ".ctg", ".sgx", ".tgx", ".ws", ".wsc", ".iso":
		return noop, true
	case ".a78":
		return decodeA78, true
	case ".lnx", ".lyx":
		return decodeLNX, true
	case ".smd":
		return decodeSMD, true
	case ".mgd":
		return decodeMGD, true
	case ".gen", ".md":
		return decodeGEN, true
	case ".n64", ".v64", ".z64":
		return decodeN64, true
	case ".nes":
		return decodeNES, true
	case ".smc", ".sfc", ".fig", ".swc":
		return decodeSNES, true
	default:
		// return noop, extra[ext]
		return noop, false
	}
}

// KnownExt returns true if the ext is recognized.
func KnownExt(ext string) bool {
	ext = strings.ToLower(ext)

	if ext == ".zip" || ext == ".gz" {
		return true
	}

	if ext == ".7z" && has7z() {
		return true
	}

	_, ok := getDecoder(ext)

	return ok
}

// decode takes a path and returns a reader for the inner rom data.
func decode(filePath string) (io.ReadCloser, error) {
	ext := strings.ToLower(path.Ext(filePath))

	if ext == ".zip" {
		return decodeZip(filePath)
	}

	if ext == ".gz" {
		return decodeGZip(filePath)
	}

	if ext == ".7z" && has7z() {
		return decode7Zip(filePath)
	}

	decode, _ := getDecoder(ext)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filePath, err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %q: %w", filePath, err)
	}

	return decode(file, fileInfo.Size())
}

// Hash returns the hash of a rom given a path to the file and hash function.
func Hash(filePath string, hasher hash.Hash, buf []byte) (string, error) {
	reader, err := decode(filePath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	for {
		readCount, err := reader.Read(buf)

		if err != nil && !errors.Is(err, io.EOF) {
			return "", fmt.Errorf("failed to read file %q: %w", filePath, err)
		}

		if readCount == 0 {
			break
		}

		hasher.Write(buf[:readCount])
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
