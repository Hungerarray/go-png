package png

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// constants
const (
	pngHeaderSize = 8
)

// errors
var (
	ErrInvalidPngFile = errors.New("not a valid png file")
	ErrCorruptedFile  = errors.New("corrupted file")
)

var pngHeader []byte = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

type Png struct {
	pngHeaderInfo
	chunks               []pngChunk
	idatCount, exifCount int
	exifdata             string
}

type pngChunk struct {
	length    uint32
	chunkType string
	data      []byte
	crc       uint32
}

type pngHeaderInfo struct {
	Width, Height                                                         uint32
	BitDepth, ColorType, CompressionMethod, FilterMethod, InterlaceMethod uint8
}

func readFullBytes(r io.Reader, p []byte) error {
	count, err := r.Read(p)
	if err != nil {
		return err
	}
	if count != len(p) {
		return ErrCorruptedFile
	}
	return nil
}

func nextChunk(r io.Reader) (*pngChunk, error) {
	var buf [8]byte
	err := readFullBytes(r, buf[:])
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(buf[:4])
	chunkType := string(buf[4:])

	data := make([]byte, length)
	err = readFullBytes(r, data)
	if err != nil {
		return nil, err
	}

	err = readFullBytes(r, buf[:4])
	if err != nil {
		return nil, err
	}

	crc := binary.BigEndian.Uint32(buf[:4])

	// todo: validate crc
	return &pngChunk{
		length:    length,
		chunkType: chunkType,
		data:      data,
		crc:       crc,
	}, nil
}

func validate(r io.Reader) bool {
	var buf [pngHeaderSize]byte
	count, err := r.Read(buf[:])
	if count != pngHeaderSize || err != nil {
		return false
	}

	for i, v := range pngHeader {
		if buf[i] != v {
			return false
		}
	}
	return true
}

func readHeaderInfo(chunk *pngChunk) pngHeaderInfo {
	if chunk.chunkType != "IHDR" || chunk.length != 13 {
		panic("passed a non valid header chunk")
	}

	return pngHeaderInfo{
		Width:             binary.BigEndian.Uint32(chunk.data[:4]),
		Height:            binary.BigEndian.Uint32(chunk.data[4:8]),
		BitDepth:          chunk.data[8],
		ColorType:         chunk.data[9],
		CompressionMethod: chunk.data[10],
		FilterMethod:      chunk.data[11],
		InterlaceMethod:   chunk.data[12],
	}
}

func NewPng(filepath string) (*Png, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	if ok := validate(file); !ok {
		return nil, ErrInvalidPngFile
	}

	chunks := make([]pngChunk, 1)
	chunk, err := nextChunk(file)
	if err != nil {
		return nil, err
	}
	chunks = append(chunks, *chunk)
	headerInfo := readHeaderInfo(chunk)

	idatCount := 0
	exifCount := 0

	var exif string
outer:
	for {
		chunk, err := nextChunk(file)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, *chunk)

		switch chunk.chunkType {
		case "IEND":
			break outer
		case "IDAT":
			idatCount += 1
		case "eXIf":
			exifCount += 1
			exif += string(chunk.data) + "\n"
		}
	}

	return &Png{
		pngHeaderInfo: headerInfo,
		chunks:        chunks,
		idatCount:     idatCount,
		exifCount:     exifCount,
		exifdata:      exif,
	}, nil
}

func (p *Png) PrintInfo() {
	fmt.Println("== PNG INFO ==")
	fmt.Printf("\tWidth: %d\n", p.Width)
	fmt.Printf("\tHeight: %d\n", p.Height)
	fmt.Printf("\tBit Depth: %d\n", p.BitDepth)
	fmt.Printf("\tColor Type: %d\n", p.ColorType)
	fmt.Printf("\tCompression Method: %d\n", p.CompressionMethod)
	fmt.Printf("\tFilter method: %d\n", p.FilterMethod)
	fmt.Printf("\tInterlaceMethod: %d\n", p.InterlaceMethod)
	fmt.Printf("\tIDAT count: %d\n", p.idatCount)
	fmt.Printf("\teXIf count: %d\n", p.exifCount)
	fmt.Printf("\teXIf data: %s\n", p.exifdata)
}
