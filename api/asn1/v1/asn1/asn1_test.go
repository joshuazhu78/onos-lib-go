// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package asn1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateValue(t *testing.T) {
	bs := BitString{
		Value: []byte{0x3f, 0xff, 0xfd},
		Len:   22,
	}
	t.Logf("%x", bs.GetValue())
	newValue, err := bs.UpdateValue([]byte{0xfd, 0xee, 0x3f})
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0xfd, 0xee, 0x3f}, newValue)
	t.Logf("%x", newValue)
	assert.EqualValues(t, []byte{0xfd, 0xee, 0x3f}, bs.GetValue())
	t.Logf("%x", bs.GetValue())

	// Testing case of an empty BitString value
	bs1 := BitString{
		//Value: []byte{0x3f, 0xff, 0xfd},
		Len: 31,
	}
	t.Logf("%x", bs1.GetValue())
	newValue1, err := bs1.UpdateValue([]byte{0xfd, 0xe4, 0xff, 0x1c})
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0xfd, 0xe4, 0xff, 0x1c}, newValue1)
	t.Logf("%x", newValue1)
	assert.EqualValues(t, []byte{0xfd, 0xe4, 0xff, 0x1c}, bs1.GetValue())
	t.Logf("%x", bs1.GetValue())

	bs2 := BitString{
		Value: make([]byte, 0),
		Len:   40,
	}
	t.Logf("%x", bs2.GetValue())
	newValue2, err := bs2.UpdateValue([]byte{0xbd, 0xe4, 0xaa, 0x1c, 0xd3})
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0xbd, 0xe4, 0xaa, 0x1c, 0xd3}, newValue2)
	t.Logf("%x", newValue2)
	assert.EqualValues(t, []byte{0xbd, 0xe4, 0xaa, 0x1c, 0xd3}, bs2.GetValue())
	t.Logf("%x", bs2.GetValue())
}

func Test_TruncateValue(t *testing.T) {
	bs := BitString{
		Value: []byte{0x3f, 0xff, 0xfd},
		Len:   22,
	}
	t.Logf("%x", bs.GetValue())
	newValue, err := bs.TruncateValue()
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0xfc}, newValue)
	t.Logf("%x", newValue)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0xfc}, bs.GetValue())
	t.Logf("%x", bs.GetValue())

	bs1 := BitString{
		Value: []byte{0x3f, 0xff, 0xff, 0xfd},
		Len:   28,
	}
	t.Logf("%x", bs1.GetValue())
	newValue1, err := bs1.TruncateValue()
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0xff, 0xf0}, newValue1)
	t.Logf("%x", newValue1)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0xff, 0xf0}, bs1.GetValue())
	t.Logf("%x", bs1.GetValue())

	bs2 := BitString{
		Value: []byte{0x3f, 0xff, 0xfd, 0xff},
		Len:   25,
	}
	t.Logf("%x", bs2.GetValue())
	newValue2, err := bs2.TruncateValue()
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0xfd, 0x80}, newValue2)
	t.Logf("%x", newValue2)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0xfd, 0x80}, bs2.GetValue())
	t.Logf("%x", bs2.GetValue())

	bs3 := BitString{
		Value: []byte{0x3f, 0xff, 0xfd, 0xff, 0x55},
		Len:   34,
	}
	t.Logf("%x", bs3.GetValue())
	newValue3, err := bs3.TruncateValue()
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0xfd, 0xff, 0x40}, newValue3)
	t.Logf("%x", newValue3)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0xfd, 0xff, 0x40}, bs3.GetValue())
	t.Logf("%x", bs3.GetValue())

	bs4 := BitString{
		Value: []byte{0x3f, 0xff, 0xfd},
		Len:   17,
	}
	t.Logf("%x", bs4.GetValue())
	newValue4, err := bs4.TruncateValue()
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0x80}, newValue4)
	t.Logf("%x", newValue4)
	assert.EqualValues(t, []byte{0x3f, 0xff, 0x80}, bs4.GetValue())
	t.Logf("%x", bs4.GetValue())
}

func Test_BitOps(t *testing.T) {
	r := &BitString{
		Value: []byte{2},
		Len:   4,
	}
	m := r.GetMaxBitOne()
	assert.Equal(t, m, 1)

	r = &BitString{
		Value: []byte{7},
		Len:   4,
	}
	m = r.GetMaxBitOne()
	assert.Equal(t, m, 2)

	r = &BitString{
		Value: []byte{8},
		Len:   4,
	}
	m = r.GetMaxBitOne()
	assert.Equal(t, m, 3)

	r = &BitString{
		Value: []byte{2, 8},
		Len:   10,
	}
	r.FromUint(10)
	assert.Equal(t, []byte{0, 10}, r.GetValue())
	r.Value[1] = 11
	assert.Equal(t, uint(11), r.ToUint())

	r = &BitString{
		Value: []byte{2, 8},
		Len:   10,
	}
	r.FromUint(1023)
	assert.Equal(t, []byte{3, 255}, r.GetValue())
	assert.Equal(t, uint(1023), r.ToUint())
	r.AddUint(1)
	assert.Equal(t, []byte{0, 0}, r.GetValue())
}

func Test_SubBitString(t *testing.T) {
	r := &BitString{
		Value: []byte{2},
		Len:   4,
	}
	s := r.SubBitString(1, 1)
	assert.True(t, reflect.DeepEqual(s, &BitString{Len: 1, Value: []byte{1}}))

	r = &BitString{
		Value: []byte{0x6, 0x03},
		Len:   11,
	}
	s = r.SubBitString(0, 2)
	assert.True(t, reflect.DeepEqual(s, &BitString{Len: 2, Value: []byte{3}}))
	s = r.SubBitString(9, 2)
	assert.True(t, reflect.DeepEqual(s, &BitString{Len: 2, Value: []byte{3}}))

	r = &BitString{
		Value: []byte{0x07, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Len:   139,
	}
	s = r.SubBitString(128, 11)
	assert.True(t, reflect.DeepEqual(s, &BitString{Len: 11, Value: []byte{0x07, 0xff}}))
}
