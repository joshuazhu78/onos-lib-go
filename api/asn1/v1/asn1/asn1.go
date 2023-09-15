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

func (b *BitString) FromUint(value uint) {
	mask := uint((1 << b.GetLen()) - 1)
	v := value & mask
	bytes := int(math.Ceil(float64(b.GetLen()) / 8.0))
	for i := 0; i < bytes; i++ {
		B := v & 0xff
		b.Value[bytes-1-i] = byte(B)
		v = v >> 8
	}
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
