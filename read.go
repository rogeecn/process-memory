package pm

import (
	"encoding/binary"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/rodrigocfd/windigo/win"
	"math"
	"strings"
)

type HProcess struct {
	hProcess win.HPROCESS
}

func New(h win.HPROCESS) HProcess {
	return HProcess{h}
}

func (h HProcess) read(address, size uint32) ([]byte, error) {
	buffer := make([]byte, size)
	_, err := h.hProcess.ReadProcessMemory(uintptr(address), buffer)
	if err != nil {
		return nil, err
	}
	return buffer, err
}

// Read process memory and convert the returned data to byte
func (h HProcess) ReadByte(lpBaseAddress uint32) (byte, error) {
	data, err := h.read(lpBaseAddress, 4)
	if err != nil {
		return 0, err
	}

	return data[0], nil
}

// Read process memory and convert the returned data to uint32
func (h HProcess) ReadUint32(lpBaseAddress uint32) (uint32, error) {
	data, err := h.read(lpBaseAddress, 4)

	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(data), err
}

// Read process memory and convert the returned data to float
func (h HProcess) ReadFloat32(lpBaseAddress uint32) (float32, error) {
	data, err := h.read(lpBaseAddress, 8)

	if err != nil {
		return 0, err
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(data)), err
}

// Read process memory and convert the returned data to string
func (h HProcess) ReadString(lpBaseAddress uint32, size uint32) (ConvertString, error) {
	data, err := h.read(lpBaseAddress, size)
	if err != nil {
		return ConvertString{}, err
	}

	var bytes []byte
	for _, c := range data {
		if c == 0 {
			break
		}
		bytes = append(bytes, c)
	}

	return ConvertString{bytes}, err
}

type ConvertString struct {
	data []byte
}

func (c ConvertString) FromGBK() (string, error) {
	enc := mahonia.NewDecoder("GBK")
	return strings.Trim(enc.ConvertString(fmt.Sprintf("%s", c.data)), " "), nil
}
