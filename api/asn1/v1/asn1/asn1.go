// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package asn1

import (
	"math"

	"github.com/onosproject/onos-lib-go/pkg/errors"
)

// UpdateValue - replace the bytes value with values from a new []byte
// the size stays the same
func (m *BitString) UpdateValue(newBytes []byte) ([]byte, error) {
	if m == nil {
		return m.Value, errors.NewInvalid("null")
	}
	expectedLen := int(math.Ceil(float64(m.Len) / 8.0))
	if len(newBytes) != expectedLen {
		return m.Value, errors.NewInvalid("too many bytes %d. Expecting %d", len(newBytes), expectedLen)
	}
	m.Value = newBytes
	return m.Value, nil
}

// TruncateValue - truncates value of trailing bits in the BitString the size stays the same
// Assuming that BitString has a non-empty length
func (m *BitString) TruncateValue() ([]byte, error) {
	if m == nil {
		return m.Value, errors.NewInvalid("null")
	}
	if m.Len == 0 {
		return nil, errors.NewInvalid("Length should not be 0")
	}
	// Computing the number of bytes
	expectedBytesLen := int(math.Ceil(float64(m.Len) / 8.0))
	if len(m.Value) != expectedBytesLen {
		return m.Value, errors.NewInvalid("too many bytes %d. Expecting %d", len(m.Value), expectedBytesLen)
	}
	// Creating set of truncated bytes, with trailing zeroes
	// Since we've got there, value in expectedLen is correct
	truncBytes := make([]byte, expectedBytesLen)
	for i := 0; i < expectedBytesLen; i++ {
		truncBytes[i] = m.Value[i]
	}

	bitsFull := expectedBytesLen * 8
	trailingBits := uint32(bitsFull) - m.Len

	mask := ^((1 << trailingBits) - 1)
	truncBytes[len(truncBytes)-1] = truncBytes[len(truncBytes)-1] & byte(mask)
	//fmt.Printf("Last byte after truncation is %x\n", truncBytes[len(truncBytes)-1])
	m.Value = truncBytes
	return m.Value, nil
}

func (b *BitString) GetBit(pos int) (bool, error) {
	if pos >= int(b.GetLen()) || pos < 0 {
		return false, errors.NewInvalid("pos %d is out of range [0,%d)", pos, b.GetLen())
	}
	B := b.GetValue()[len(b.GetValue())-1-int(pos/8)]
	return ((B & (1 << (pos % 8))) != 0), nil
}

func (b *BitString) GetMaxBitOne() int {
	var max int
	for max = int(b.GetLen()) - 1; max >= 0; max-- {
		if bb, err := b.GetBit(max); err == nil && bb {
			break
		}
	}
	return max
}

func UintToBytes(value uint, len uint32) []byte {
	mask := uint((1 << len) - 1)
	v := value & mask
	numByptes := int(math.Ceil(float64(len) / 8.0))
	bytes := make([]byte, numByptes)
	for i := 0; i < numByptes; i++ {
		B := v & 0xff
		bytes[numByptes-1-i] = byte(B)
		v = v >> 8
	}
	return bytes
}

func (b *BitString) FromUint(value uint) {
	b.Value = UintToBytes(value, b.GetLen())
}

func (b *BitString) ToUint() uint {
	var v uint = 0
	bytes := int(math.Ceil(float64(b.GetLen()) / 8.0))
	for i := 0; i < bytes; i++ {
		v = v << 8
		v = v + uint(b.Value[i])
	}
	return v
}

func (b *BitString) AddUint(a uint) {
	v := b.ToUint()
	v = v + a
	b.FromUint(v)
}

func (b *BitString) SubBitString(start int, len uint32) *BitString {
	numBytes := int(math.Ceil(float64(len) / 8.0))
	v := make([]byte, numBytes)
	bb := 0
	BB := 0
	var bitsFirstByte uint32
	if len < uint32(numBytes*8) {
		bitsFirstByte = len - (uint32(numBytes-1) * 8)
		var B byte = 0
		for ; bb < int(bitsFirstByte); bb++ {
			bit, _ := b.GetBit(start + bb)
			B = B << 1
			if bit {
				B = (B + byte(0x01))
			}
		}
		v[BB] = B
		BB = BB + 1
	}
	var B byte = 0
	for ; bb < int(len); bb++ {
		bit, _ := b.GetBit(start + bb)
		B = B << 1
		if bit {
			B = (B + byte(0x01))
		}
		if (bb-int(bitsFirstByte))%8 == 7 {
			v[BB] = B
			B = 0
			BB = BB + 1
		}
	}
	return &BitString{
		Len:   len,
		Value: v,
	}
}

func NewBitString(value uint, len uint32) *BitString {
	return &BitString{
		Value: UintToBytes(value, len),
		Len:   len,
	}
}
