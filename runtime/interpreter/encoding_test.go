/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2021 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package interpreter_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/fxamacker/atree"
	. "github.com/onflow/cadence/runtime/interpreter"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/tests/utils"
	. "github.com/onflow/cadence/runtime/tests/utils"
)

type encodeDecodeTest struct {
	value        Value
	storable     atree.Storable
	encoded      []byte
	invalid      bool
	decodedValue Value
	decodeOnly   bool
	deepEquality bool
	storage      Storage
}

var testOwner = common.BytesToAddress([]byte{0x42})

func testEncodeDecode(t *testing.T, test encodeDecodeTest) {

	if test.storage == nil {
		test.storage = NewInMemoryStorage()
	}

	var encoded []byte
	if (test.value != nil || test.storable != nil) && !test.decodeOnly {

		if test.value != nil {
			test.value.SetOwner(&testOwner)
			if test.storable == nil {
				test.storable = test.value.Storable(test.storage)
			}
		}

		var err error
		encoded, err = atree.Encode(test.storable, test.storage)
		require.NoError(t, err)

		if test.encoded != nil {
			AssertEqualWithDiff(t, test.encoded, encoded)
		}
	} else {
		encoded = test.encoded
	}

	decoder := DecMode.NewByteStreamDecoder(encoded)
	decoded, err := DecodeStorableV6(decoder, test.storage)

	if test.invalid {
		require.Error(t, err)
	} else {
		require.NoError(t, err)

		decodedValue, err := StoredValue(decoded, test.storage)
		require.NoError(t, err)

		expectedValue := test.value
		if test.decodedValue != nil {
			test.decodedValue.SetOwner(&testOwner)
			expectedValue = test.decodedValue
		}
		if test.deepEquality {
			assert.Equal(t, expectedValue, decodedValue)
		} else {
			AssertValuesEqual(t, expectedValue, decodedValue.(Value))
		}
	}
}

func TestEncodeDecodeNilValue(t *testing.T) {

	t.Parallel()

	testEncodeDecode(t,
		encodeDecodeTest{
			value: NilValue{},
			encoded: []byte{
				// null
				0xf6,
			},
		},
	)
}

func TestEncodeDecodeVoidValue(t *testing.T) {

	t.Parallel()

	testEncodeDecode(t,
		encodeDecodeTest{
			value: VoidValue{},
			encoded: []byte{
				// tag
				0xd8, CBORTagVoidValue,
				// null
				0xf6,
			},
		},
	)
}

func TestEncodeDecodeBool(t *testing.T) {

	t.Parallel()

	t.Run("false", func(t *testing.T) {

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: BoolValue(false),
				encoded: []byte{
					// false
					0xf4,
				},
			},
		)
	})

	t.Run("true", func(t *testing.T) {

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: BoolValue(true),
				encoded: []byte{
					// true
					0xf5,
				},
			},
		)
	})
}

func TestEncodeDecodeString(t *testing.T) {

	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		expected := NewStringValue("")

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: expected,
				encoded: []byte{
					//  UTF-8 string, 0 bytes follow
					0x60,
				},
			})
	})

	t.Run("non-empty", func(t *testing.T) {
		expected := NewStringValue("foo")

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: expected,
				encoded: []byte{
					// UTF-8 string, 3 bytes follow
					0x63,
					// f, o, o
					0x66, 0x6f, 0x6f,
				},
			},
		)
	})
}

func TestEncodeDecodeArray(t *testing.T) {

	// TODO: type
	// TODO: owner

	t.Parallel()

	t.Run("empty", func(t *testing.T) {

		storage := NewInMemoryStorage()

		expected := NewArrayValueUnownedNonCopying(
			ConstantSizedStaticType{
				Type: PrimitiveStaticTypeAnyStruct,
				Size: 0,
			},
			storage,
		)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				value:   expected,
				encoded: []byte{
					// tag
					0xd8, atree.CBORTagStorageID,

					// storage ID
					0x1,
				},
			})
	})

	t.Run("string and bool", func(t *testing.T) {

		storage := NewInMemoryStorage()

		expectedString := NewStringValue("test")

		expected := NewArrayValueUnownedNonCopying(
			VariableSizedStaticType{
				Type: PrimitiveStaticTypeAnyStruct,
			},
			storage,
			expectedString,
			BoolValue(true),
		)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				value:   expected,
				encoded: []byte{
					// tag
					0xd8, atree.CBORTagStorageID,

					// storage ID
					0x1,
				},
			},
		)
	})
}

func TestEncodeDecodeDictionary(t *testing.T) {

	// TODO: type
	// TODO: owner
	// TODO: storage ID

	t.Parallel()

	t.Run("empty", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		expected := NewDictionaryValueUnownedNonCopying(
			DictionaryStaticType{
				KeyType:   PrimitiveStaticTypeString,
				ValueType: PrimitiveStaticTypeAnyStruct,
			},
			storage,
		)

		encodedValue := []byte{
			// tag
			0xd8, atree.CBORTagStorageID,

			// storage ID
			1,
		}

		encodedStorable := []byte{
			// tag
			0xd8, CBORTagDictionaryValue,
			// array, 3 items follow
			0x83,

			// dictionary type tag
			0xd8, CBORTagDictionaryStaticType,
			// array, 2 items follow
			0x82,
			// key type
			0xd8, CBORTagPrimitiveStaticType, byte(PrimitiveStaticTypeString),
			// value type
			0xd8, CBORTagPrimitiveStaticType, byte(PrimitiveStaticTypeAnyStruct),

			// cbor Array Value tag
			0xd8, atree.CBORTagStorageID,

			// storage ID.
			// 2 instead of 1, because dictionary has storage ID 1
			0x2,

			// array, 0 items follow
			0x80,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				value:   expected,
				encoded: encodedValue,
			},
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				storable: DictionaryStorable{
					Dictionary: expected,
				},
				encoded:      encodedStorable,
				decodedValue: expected,
			},
		)
	})

	t.Run("non-empty", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		key1 := BoolValue(true)
		value1 := BoolValue(false)

		key2 := NewStringValue("foo")
		value2 := NewStringValue("bar")

		expected := NewDictionaryValueUnownedNonCopying(
			DictionaryStaticType{
				KeyType:   PrimitiveStaticTypeAnyStruct,
				ValueType: PrimitiveStaticTypeAnyStruct,
			},
			storage,
			key1, value1,
			key2, value2,
		)

		encodedValue := []byte{
			// tag
			0xd8, atree.CBORTagStorageID,

			// storage ID.
			1,
		}

		encodedStorable := []byte{
			// tag
			0xd8, CBORTagDictionaryValue,
			// array, 3 items follow
			0x83,

			// dictionary type tag
			0xd8, CBORTagDictionaryStaticType,
			// array, 2 items follow
			0x82,
			// key type
			0xd8, CBORTagPrimitiveStaticType, byte(PrimitiveStaticTypeAnyStruct),
			// value type
			0xd8, CBORTagPrimitiveStaticType, byte(PrimitiveStaticTypeAnyStruct),

			// Keys

			// cbor Array Value tag
			0xd8, atree.CBORTagStorageID,

			// storage ID.
			// 2 instead of 1, because dictionary has storage ID 1
			0x2,

			// Values

			// array, 2 items follow
			0x82,

			// false
			0xf4,

			// UTF-8 string, length 3
			0x63,
			// b, a, r
			0x62, 0x61, 0x72,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				value:   expected,
				encoded: encodedValue,
			},
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				storable: DictionaryStorable{
					Dictionary: expected,
				},
				encoded:      encodedStorable,
				decodedValue: expected,
			},
		)
	})
}

func TestEncodeDecodeComposite(t *testing.T) {

	// TODO: owner
	// TODO: storage ID

	t.Parallel()

	t.Run("empty structure, string location, qualified identifier", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		expected := NewCompositeValue(
			storage,
			utils.TestLocation,
			"TestStruct",
			common.CompositeKindStructure,
			NewStringValueOrderedMap(),
			nil,
		)

		encodedValue := []byte{
			// tag
			0xd8, atree.CBORTagStorageID,

			// storage ID
			1,
		}

		encodedStorable := []byte{
			// tag
			0xd8, CBORTagCompositeValue,
			// array, 4 items follow
			0x84,

			// tag
			0xd8, CBORTagStringLocation,
			// UTF-8 string, length 4
			0x64,
			// t, e, s, t
			0x74, 0x65, 0x73, 0x74,

			// positive integer 1
			0x1,

			// array, 0 items follow
			0x80,

			// UTF-8 string, length 10
			0x6a,
			0x54, 0x65, 0x73, 0x74, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				storable: CompositeStorable{
					Composite: expected,
				},
				encoded:      encodedStorable,
				decodedValue: expected,
			},
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				value:   expected,
				encoded: encodedValue,
			},
		)
	})

	t.Run("non-empty resource, qualified identifier", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		stringValue := NewStringValue("test")

		members := NewStringValueOrderedMap()
		members.Set("string", stringValue)
		members.Set("true", BoolValue(true))

		expected := NewCompositeValue(
			storage,
			utils.TestLocation,
			"TestResource",
			common.CompositeKindResource,
			members,
			nil,
		)

		encodedValue := []byte{
			// tag
			0xd8, atree.CBORTagStorageID,

			// storage ID
			1,
		}

		encodedStorable := []byte{
			// tag
			0xd8, CBORTagCompositeValue,
			// array, 4 items follow
			0x84,

			// tag
			0xd8, CBORTagStringLocation,
			// UTF-8 string, length 4
			0x64,
			// t, e, s, t
			0x74, 0x65, 0x73, 0x74,

			// positive integer 2
			0x2,

			// array, 4 items follow
			0x84,
			// UTF-8 string, length 6
			0x66,
			// s, t, r, i, n, g
			0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
			// UTF-8 string, length 4
			0x64,
			// t, e, s, t
			0x74, 0x65, 0x73, 0x74,
			// UTF-8 string, length 4
			0x64,
			// t, r, u, e
			0x74, 0x72, 0x75, 0x65,
			// true
			0xf5,

			// UTF-8 string, length 12
			0x6c,
			0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				storable: CompositeStorable{
					Composite: expected,
				},
				encoded:      encodedStorable,
				decodedValue: expected,
			},
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				value:   expected,
				encoded: encodedValue,
			},
		)
	})

	t.Run("empty, address location", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		expected := NewCompositeValue(
			storage,
			common.AddressLocation{
				Address: common.BytesToAddress([]byte{0x1}),
				Name:    "TestStruct",
			},
			"TestStruct",
			common.CompositeKindStructure,
			NewStringValueOrderedMap(),
			nil,
		)

		encodedValue := []byte{
			// tag
			0xd8, atree.CBORTagStorageID,

			// storage ID
			1,
		}

		encodedStorable := []byte{
			// tag
			0xd8, CBORTagCompositeValue,
			// array, 4 items follow
			0x84,

			// tag
			0xd8, CBORTagAddressLocation,
			// array, 2 items follow
			0x82,
			// byte sequence, length 1
			0x41,
			// positive integer 1
			0x1,
			// UTF-8 string, length 10
			0x6a,
			0x54, 0x65, 0x73, 0x74, 0x53, 0x74, 0x72, 0x75,
			0x63, 0x74,

			// positive integer 1
			0x1,

			// array, 0 items follow
			0x80,

			// UTF-8 string, length 10
			0x6a,
			0x54, 0x65, 0x73, 0x74, 0x53, 0x74, 0x72, 0x75,
			0x63, 0x74,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				value:   expected,
				encoded: encodedValue,
			},
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				storage: storage,
				storable: CompositeStorable{
					Composite: expected,
				},
				encoded:      encodedStorable,
				decodedValue: expected,
			},
		)
	})
}

func TestEncodeDecodeIntValue(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewIntValueFromInt64(0),
				encoded: []byte{
					0xd8, CBORTagIntValue,
					// positive bignum
					0xc2,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewIntValueFromInt64(42),
				encoded: []byte{
					0xd8, CBORTagIntValue,
					// positive bignum
					0xc2,
					// byte string, length 1
					0x41,
					0x2a,
				},
			},
		)
	})

	t.Run("negative one", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewIntValueFromInt64(-1),
				encoded: []byte{
					0xd8, CBORTagIntValue,
					// negative bignum
					0xc3,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewIntValueFromInt64(-42),
				encoded: []byte{
					0xd8, CBORTagIntValue,
					// negative bignum
					0xc3,
					// byte string, length 1
					0x41,
					// `-42` in decimal is is `0x2a` in hex.
					// CBOR requires negative values to be encoded as `-1-n`, which is `-n - 1`,
					// which is `0x2a - 0x01`, which equals to `0x29`.
					0x29,
				},
			},
		)
	})

	t.Run("negative, large (> 64 bit)", func(t *testing.T) {
		setString, ok := new(big.Int).SetString("-18446744073709551617", 10)
		require.True(t, ok)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewIntValueFromBigInt(setString),
				encoded: []byte{
					0xd8, CBORTagIntValue,
					// negative bignum
					0xc3,
					// byte string, length 9
					0x49,
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
		)
	})

	t.Run("positive, large (> 64 bit)", func(t *testing.T) {
		bigInt, ok := new(big.Int).SetString("18446744073709551616", 10)
		require.True(t, ok)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewIntValueFromBigInt(bigInt),
				encoded: []byte{
					// tag
					0xd8, CBORTagIntValue,
					// positive bignum
					0xc2,
					// byte string, length 9
					0x49,
					0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				},
			},
		)
	})
}

func TestEncodeDecodeInt8Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int8Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt8Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int8Value(-42),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt8Value,
					// negative integer 42
					0x38,
					0x29,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int8Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt8Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int8Value(math.MinInt8),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt8Value,
					// negative integer 0x7f
					0x38,
					0x7f,
				},
			},
		)
	})

	t.Run("<min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagInt8Value,
					// negative integer 0xf00
					0x38,
					0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int8Value(math.MaxInt8),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt8Value,
					// positive integer 0x7f00
					0x18,
					0x7f,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagInt8Value,
					// positive integer 0xff
					0x18,
					0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeInt16Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int16Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt16Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int16Value(-42),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt16Value,
					// negative integer 42
					0x38,
					0x29,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int16Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt16Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int16Value(math.MinInt16),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt16Value,
					// negative integer 0x7fff
					0x39,
					0x7f, 0xff,
				},
			},
		)
	})

	t.Run("<min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagInt16Value,
					// negative integer 0xffff
					0x39,
					0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int16Value(math.MaxInt16),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt16Value,
					// positive integer 0x7fff
					0x19,
					0x7f, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagInt16Value,
					// positive integer 0xffff
					0x19,
					0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeInt32Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int32Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt32Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int32Value(-42),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt32Value,
					// negative integer 42
					0x38,
					0x29,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int32Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt32Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int32Value(math.MinInt32),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt32Value,
					// negative integer 0x7fffffff
					0x3a,
					0x7f, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run("<min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagInt32Value,
					// negative integer 0xffffffff
					0x3a,
					0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int32Value(math.MaxInt32),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt32Value,
					// positive integer 0x7fffffff
					0x1a,
					0x7f, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagInt32Value,
					// positive integer 0xffffffff
					0x1a,
					0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeInt64Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int64Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt64Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int64Value(-42),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt64Value,
					// negative integer 42
					0x38,
					0x29,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int64Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt64Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int64Value(math.MinInt64),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt64Value,
					// negative integer: 0x7fffffffffffffff
					0x3b,
					0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run("<min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagInt64Value,
					// negative integer 0xffffffffffffffff
					0x3b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Int64Value(math.MaxInt64),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt64Value,
					// positive integer: 0x7fffffffffffffff
					0x1b,
					0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagInt64Value,
					// positive integer 0xffffffffffffffff
					0x1b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeInt128Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt128ValueFromInt64(0),
				encoded: []byte{
					0xd8, CBORTagInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt128ValueFromInt64(42),
				encoded: []byte{
					0xd8, CBORTagInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 1
					0x41,
					0x2a,
				},
			},
		)
	})

	t.Run("negative one", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt128ValueFromInt64(-1),
				encoded: []byte{
					0xd8, CBORTagInt128Value,
					// negative bignum
					0xc3,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt128ValueFromInt64(-42),
				encoded: []byte{
					0xd8, CBORTagInt128Value,
					// negative bignum
					0xc3,
					// byte string, length 1
					0x41,
					0x29,
				},
			},
		)
	})

	t.Run("min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt128ValueFromBigInt(sema.Int128TypeMinIntBig),
				encoded: []byte{
					0xd8, CBORTagInt128Value,
					// negative bignum
					0xc3,
					// byte string, length 16
					0x50,
					0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run("<min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagInt128Value,
					// negative bignum
					0xc3,
					// byte string, length 16
					0x50,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt128ValueFromBigInt(sema.Int128TypeMaxIntBig),
				encoded: []byte{
					0xd8, CBORTagInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 16
					0x50,
					0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 16
					0x50,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("RFC", func(t *testing.T) {
		rfcValue, ok := new(big.Int).SetString("18446744073709551616", 10)
		require.True(t, ok)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt128ValueFromBigInt(rfcValue),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 9
					0x49,
					0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				},
			},
		)
	})
}

func TestEncodeDecodeInt256Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt256ValueFromInt64(0),
				encoded: []byte{
					0xd8, CBORTagInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt256ValueFromInt64(42),
				encoded: []byte{
					0xd8, CBORTagInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 1
					0x41,
					0x2a,
				},
			},
		)
	})

	t.Run("negative one", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt256ValueFromInt64(-1),
				encoded: []byte{
					0xd8, CBORTagInt256Value,
					// negative bignum
					0xc3,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt256ValueFromInt64(-42),
				encoded: []byte{
					0xd8, CBORTagInt256Value,
					// negative bignum
					0xc3,
					// byte string, length 1
					0x41,
					0x29,
				},
			},
		)
	})

	t.Run("min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt256ValueFromBigInt(sema.Int256TypeMinIntBig),
				encoded: []byte{
					0xd8, CBORTagInt256Value,
					// negative bignum
					0xc3,
					// byte string, length 32
					0x58, 0x20,
					0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run("<min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagInt256Value,
					// negative bignum
					0xc3,
					// byte string, length 32
					0x58, 0x20,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt256ValueFromBigInt(sema.Int256TypeMaxIntBig),
				encoded: []byte{
					0xd8, CBORTagInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 32
					0x58, 0x20,
					0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 32
					0x58, 0x20,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("RFC", func(t *testing.T) {

		rfcValue, ok := new(big.Int).SetString("18446744073709551616", 10)
		require.True(t, ok)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewInt256ValueFromBigInt(rfcValue),
				encoded: []byte{
					// tag
					0xd8, CBORTagInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 9
					0x49,
					0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				},
			},
		)
	})
}

func TestEncodeDecodeUIntValue(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUIntValueFromUint64(0),
				encoded: []byte{
					0xd8, CBORTagUIntValue,
					// positive bignum
					0xc2,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagUIntValue,
					// negative bignum
					0xc3,
					// byte string, length 1
					0x41,
					0x2a,
				},
				invalid: true,
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUIntValueFromUint64(42),
				encoded: []byte{
					0xd8, CBORTagUIntValue,
					// positive bignum
					0xc2,
					// byte string, length 1
					0x41,
					0x2a,
				},
			},
		)
	})

	t.Run("RFC", func(t *testing.T) {

		rfcValue, ok := new(big.Int).SetString("18446744073709551616", 10)
		require.True(t, ok)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUIntValueFromBigInt(rfcValue),
				encoded: []byte{
					// tag
					0xd8, CBORTagUIntValue,
					// positive bignum
					0xc2,
					// byte string, length 9
					0x49,
					0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				},
			},
		)
	})
}

func TestEncodeDecodeUInt8Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt8Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt8Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt8Value,
					// negative integer 42
					0x38,
					0x29,
				},
				invalid: true,
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt8Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt8Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt8Value(math.MaxUint8),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt8Value,
					// positive integer 0xff
					0x18,
					0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt8Value,
					// positive integer 0xffff
					0x19,
					0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeUInt16Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt16Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt16Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt16Value,
					// negative integer 42
					0x38,
					0x29,
				},
				invalid: true,
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt16Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt16Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt16Value(math.MaxUint16),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt16Value,
					// positive integer 0xffff
					0x19,
					0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt16Value,
					// positive integer 0xffffffff
					0x1a,
					0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeUInt32Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt32Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt32Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt32Value,
					// negative integer 42
					0x38,
					0x29,
				},
				invalid: true,
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt32Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt32Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt32Value(math.MaxUint32),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt32Value,
					// positive integer 0xffffffff
					0x1a,
					0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt32Value,
					// positive integer 0xffffffffffffffff
					0x1b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeUInt64Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt64Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt64Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt64Value,
					// negative integer 42
					0x38,
					0x29,
				},
				invalid: true,
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt64Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt64Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UInt64Value(math.MaxUint64),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt64Value,
					// positive integer 0xffffffffffffffff
					0x1b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})
}

func TestEncodeDecodeUInt128Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUInt128ValueFromUint64(0),
				encoded: []byte{
					0xd8, CBORTagUInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUInt128ValueFromUint64(42),
				encoded: []byte{
					0xd8, CBORTagUInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 1
					0x41,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUInt128ValueFromBigInt(sema.UInt128TypeMaxIntBig),
				encoded: []byte{
					0xd8, CBORTagUInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 16
					0x50,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagUInt128Value,
					// negative bignum
					0xc3,
					// byte string, length 1
					0x41,
					0x2a,
				},
				invalid: true,
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagUInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 17
					0x51,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("RFC", func(t *testing.T) {
		rfcValue, ok := new(big.Int).SetString("18446744073709551616", 10)
		require.True(t, ok)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUInt128ValueFromBigInt(rfcValue),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt128Value,
					// positive bignum
					0xc2,
					// byte string, length 9
					0x49,
					0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				},
			},
		)
	})
}

func TestEncodeDecodeUInt256Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUInt256ValueFromUint64(0),
				encoded: []byte{
					0xd8, CBORTagUInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 0
					0x40,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUInt256ValueFromUint64(42),
				encoded: []byte{
					0xd8, CBORTagUInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 1
					0x41,
					0x2a,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagUInt256Value,
					// negative bignum
					0xc3,
					// byte string, length 1
					0x41,
					0x2a,
				},
				invalid: true,
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					0xd8, CBORTagUInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 65
					0x58, 0x41,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("RFC", func(t *testing.T) {
		rfcValue, ok := new(big.Int).SetString("18446744073709551616", 10)
		require.True(t, ok)

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: NewUInt256ValueFromBigInt(rfcValue),
				encoded: []byte{
					// tag
					0xd8, CBORTagUInt256Value,
					// positive bignum
					0xc2,
					// byte string, length 9
					0x49,
					0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				},
			},
		)
	})
}

func TestEncodeDecodeWord8Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word8Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord8Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagWord8Value,
					// negative integer 42
					0x38,
					0x29,
				},
				invalid: true,
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word8Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord8Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word8Value(math.MaxUint8),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord8Value,
					// positive integer 0xff
					0x18,
					0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagWord8Value,
					// positive integer 0xffff
					0x19,
					0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeWord16Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word16Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord16Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word16Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord16Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word16Value(math.MaxUint16),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord16Value,
					// positive integer 0xffff
					0x19,
					0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagWord16Value,
					// positive integer 0xffffffff
					0x1a,
					0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeWord32Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word32Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord32Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word32Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord32Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word32Value(math.MaxUint32),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord32Value,
					// positive integer 0xffffffff
					0x1a,
					0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagWord32Value,
					// positive integer 0xffffffffffffffff
					0x1b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})
}

func TestEncodeDecodeWord64Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word64Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord64Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word64Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord64Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Word64Value(math.MaxUint64),
				encoded: []byte{
					// tag
					0xd8, CBORTagWord64Value,
					// positive integer 0xffffffffffffffff
					0x1b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})
}

func TestEncodeDecodeSomeValue(t *testing.T) {

	// TODO:
	t.Skip()

	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: &SomeValue{
					Value: NilValue{},
				},
				encoded: []byte{
					// tag
					0xd8, CBORTagSomeValue,
					// null
					0xf6,
				},
			},
		)
	})

	t.Run("string", func(t *testing.T) {
		expectedString := NewStringValue("test")

		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: &SomeValue{
					Value: expectedString,
				},
				encoded: []byte{
					// tag
					0xd8, CBORTagSomeValue,
					// UTF-8 string, length 4
					0x64,
					// t, e, s, t
					0x74, 0x65, 0x73, 0x74,
				},
			},
		)
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: &SomeValue{
					Value: BoolValue(true),
				},
				encoded: []byte{
					// tag
					0xd8, CBORTagSomeValue,
					// true
					0xf5,
				},
			},
		)
	})
}

func TestEncodeDecodeFix64Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Fix64Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagFix64Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Fix64Value(-42),
				encoded: []byte{
					// tag
					0xd8, CBORTagFix64Value,
					// negative integer 42
					0x38,
					0x29,
				},
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Fix64Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagFix64Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Fix64Value(math.MinInt64),
				encoded: []byte{
					// tag
					0xd8, CBORTagFix64Value,
					// negative integer: 0x7fffffffffffffff
					0x3b,
					0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run("<min", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagFix64Value,
					// negative integer 0xffffffffffffffff
					0x3b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: Fix64Value(math.MaxInt64),
				encoded: []byte{
					// tag
					0xd8, CBORTagFix64Value,
					// positive integer: 0x7fffffffffffffff
					0x1b,
					0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})

	t.Run(">max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagFix64Value,
					// positive integer 0xffffffffffffffff
					0x1b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				invalid: true,
			},
		)
	})

}

func TestEncodeDecodeUFix64Value(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UFix64Value(0),
				encoded: []byte{
					// tag
					0xd8, CBORTagUFix64Value,
					// integer 0
					0x0,
				},
			},
		)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagUFix64Value,
					// negative integer 42
					0x38,
					0x29,
				},
				invalid: true,
			},
		)
	})

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UFix64Value(42),
				encoded: []byte{
					// tag
					0xd8, CBORTagUFix64Value,
					// positive integer 42
					0x18,
					0x2a,
				},
			},
		)
	})

	t.Run("max", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: UFix64Value(math.MaxUint64),
				encoded: []byte{
					// tag
					0xd8, CBORTagUFix64Value,
					// positive integer 0xffffffffffffffff
					0x1b,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
		)
	})
}

func TestEncodeDecodeAddressValue(t *testing.T) {

	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: AddressValue{},
				encoded: []byte{
					// tag
					0xd8, CBORTagAddressValue,
					// byte sequence, length 0
					0x40,
				},
			},
		)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: AddressValue(common.BytesToAddress([]byte{0x42})),
				encoded: []byte{
					// tag
					0xd8, CBORTagAddressValue,
					// byte sequence, length 1
					0x41,
					// address
					0x42,
				},
			},
		)
	})

	t.Run("with leading zeros", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: AddressValue(common.BytesToAddress([]byte{0x0, 0x42})),
				encoded: []byte{
					// tag
					0xd8, CBORTagAddressValue,
					// byte sequence, length 1
					0x41,
					// address
					0x42,
				},
			},
		)
	})

	t.Run("with zeros in-between and at and", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				value: AddressValue(common.BytesToAddress([]byte{0x0, 0x42, 0x0, 0x43, 0x0})),
				encoded: []byte{
					// tag
					0xd8, CBORTagAddressValue,
					// byte sequence, length 4
					0x44,
					// address
					0x42, 0x0, 0x43, 0x0,
				},
			},
		)
	})

	t.Run("too long", func(t *testing.T) {
		t.Parallel()

		testEncodeDecode(t,
			encodeDecodeTest{
				encoded: []byte{
					// tag
					0xd8, CBORTagAddressValue,
					// byte sequence, length 22
					0x56,
					// address
					0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x1,
				},
				invalid: true,
			},
		)
	})
}

var privatePathValue = PathValue{
	Domain:     common.PathDomainPrivate,
	Identifier: "foo",
}

var publicPathValue = PathValue{
	Domain:     common.PathDomainPublic,
	Identifier: "bar",
}

func TestEncodeDecodePathValue(t *testing.T) {

	t.Parallel()

	t.Run("private", func(t *testing.T) {

		t.Parallel()

		encoded := []byte{
			// tag
			0xd8, CBORTagPathValue,
			// array, 2 items follow
			0x82,
			// positive integer 2
			0x2,
			// UTF-8 string, 3 bytes follow
			0x63,
			// f, o, o
			0x66, 0x6f, 0x6f,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   privatePathValue,
				encoded: encoded,
			},
		)
	})

	t.Run("public", func(t *testing.T) {

		t.Parallel()

		encoded := []byte{
			// tag
			0xd8, CBORTagPathValue,
			// array, 2 items follow
			0x82,
			// positive integer 3
			0x3,
			// UTF-8 string, 3 bytes follow
			0x63,
			// b, a, r
			0x62, 0x61, 0x72,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   publicPathValue,
				encoded: encoded,
			},
		)
	})
}

func TestEncodeDecodeCapabilityValue(t *testing.T) {

	t.Parallel()

	t.Run("private path, untyped capability, new format", func(t *testing.T) {

		t.Parallel()

		value := CapabilityValue{
			Address: NewAddressValueFromBytes([]byte{0x2}),
			Path:    privatePathValue,
		}

		encoded := []byte{
			// tag
			0xd8, CBORTagCapabilityValue,
			// array, 3 items follow
			0x83,
			// tag for address
			0xd8, CBORTagAddressValue,
			// byte sequence, length 1
			0x41,
			// address
			0x02,
			// tag for address
			0xd8, CBORTagPathValue,
			// array, 2 items follow
			0x82,
			// positive integer 2
			0x2,
			// UTF-8 string, length 3
			0x63,
			// f, o, o
			0x66, 0x6f, 0x6f,
			// nil
			0xf6,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("private path, typed capability", func(t *testing.T) {

		t.Parallel()

		value := CapabilityValue{
			Address:    NewAddressValueFromBytes([]byte{0x2}),
			Path:       privatePathValue,
			BorrowType: PrimitiveStaticTypeBool,
		}

		encoded := []byte{
			// tag
			0xd8, CBORTagCapabilityValue,
			// array, 3 items follow
			0x83,
			// tag for address
			0xd8, CBORTagAddressValue,
			// byte sequence, length 1
			0x41,
			// address
			0x02,
			// tag for address
			0xd8, CBORTagPathValue,
			// aray, 2 items follow
			0x82,
			// positive integer 2
			0x2,
			// UTF-8 string, length 3
			0x63,
			// f, o, o
			0x66, 0x6f, 0x6f,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			// bool
			0x6,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("public path, untyped capability, new format", func(t *testing.T) {

		t.Parallel()

		value := CapabilityValue{
			Address: NewAddressValueFromBytes([]byte{0x3}),
			Path:    publicPathValue,
		}

		encoded := []byte{
			// tag
			0xd8, CBORTagCapabilityValue,
			// array, 3 items follow
			0x83,
			// tag for address
			0xd8, CBORTagAddressValue,
			// byte sequence, length 1
			0x41,
			// address
			0x03,
			// tag for address
			0xd8, CBORTagPathValue,
			// array, 2 items follow
			0x82,
			// positive integer 3
			0x3,
			// UTF-8 string, length 3
			0x63,
			// b, a, r
			0x62, 0x61, 0x72,
			// nil
			0xf6,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)

	})

	t.Run("public path, typed capability", func(t *testing.T) {

		t.Parallel()

		value := CapabilityValue{
			Address:    NewAddressValueFromBytes([]byte{0x3}),
			Path:       publicPathValue,
			BorrowType: PrimitiveStaticTypeBool,
		}

		encoded := []byte{
			// tag
			0xd8, CBORTagCapabilityValue,
			// array, 3 items follow
			0x83,
			// tag for address
			0xd8, CBORTagAddressValue,
			// byte sequence, length 1
			0x41,
			// address
			0x03,
			// tag for address
			0xd8, CBORTagPathValue,
			// array, 2 items follow
			0x82,
			// positive integer 3
			0x3,
			// UTF-8 string, length 3
			0x63,
			// b, a, r
			0x62, 0x61, 0x72,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			// bool
			0x6,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	// For testing backward compatibility for native composite types
	t.Run("public path, public account typed capability", func(t *testing.T) {

		t.Parallel()

		capabilityValue := CapabilityValue{
			Address:    NewAddressValueFromBytes([]byte{0x3}),
			Path:       publicPathValue,
			BorrowType: PrimitiveStaticTypePublicAccount,
		}

		encoded := []byte{
			// tag
			0xd8, CBORTagCapabilityValue,
			// array, 3 items follow
			0x83,
			// tag for address
			0xd8, CBORTagAddressValue,
			// byte sequence, length 1
			0x41,
			// address
			0x03,
			// tag for address
			0xd8, CBORTagPathValue,
			// array, 2 items follow
			0x82,
			// positive integer 3
			0x3,
			// UTF-8 string, length 3
			0x63,
			// b, a, r
			0x62, 0x61, 0x72,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			// positive integer to follow
			0x18,
			// public account (tag)
			0x5b,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   capabilityValue,
				encoded: encoded,
			},
		)
	})
}

func TestEncodeDecodeLinkValue(t *testing.T) {

	t.Parallel()

	expectedLinkEncodingPrefix := []byte{
		// tag
		0xd8, CBORTagLinkValue,
		// array, 2 items follow
		0x82,
		0xd8, CBORTagPathValue,
		// array, 2 items follow
		0x82,
		// positive integer 3
		0x3,
		// UTF-8 string, length 3
		0x63,
		// b, a, r
		0x62, 0x61, 0x72,
	}

	t.Run("primitive, Bool", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type:       ConvertSemaToPrimitiveStaticType(sema.BoolType),
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x6,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("optional, primitive, bool", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: OptionalStaticType{
				Type: PrimitiveStaticTypeBool,
			},
		}

		encodedType := []byte{
			// tag
			0xd8, CBORTagOptionalStaticType,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x6,
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			encodedType...,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("composite, struct, qualified identifier", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: CompositeStaticType{
				Location:            utils.TestLocation,
				QualifiedIdentifier: "SimpleStruct",
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagCompositeStaticType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, CBORTagStringLocation,
			// UTF-8 string, length 4
			0x64,
			// t, e, s, t
			0x74, 0x65, 0x73, 0x74,
			// UTF-8 string, length 12
			0x6c,
			// SimpleStruct
			0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("interface, struct, qualified identifier", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: InterfaceStaticType{
				Location:            utils.TestLocation,
				QualifiedIdentifier: "SimpleInterface",
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagInterfaceStaticType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, CBORTagStringLocation,
			// UTF-8 string, length 4
			0x64,
			// t, e, s, t
			0x74, 0x65, 0x73, 0x74,
			// UTF-8 string, length 22
			0x6F,
			// SimpleInterface
			0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("variable-sized, bool", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: VariableSizedStaticType{
				Type: PrimitiveStaticTypeBool,
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagVariableSizedStaticType,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x6,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("constant-sized, bool", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: ConstantSizedStaticType{
				Type: PrimitiveStaticTypeBool,
				Size: 42,
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagConstantSizedStaticType,
			// array, 2 items follow
			0x82,
			// positive integer 42
			0x18, 0x2A,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x6,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("reference type, authorized, bool", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: ReferenceStaticType{
				Authorized: true,
				Type:       PrimitiveStaticTypeBool,
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagReferenceStaticType,
			// array, 2 items follow
			0x82,
			// true
			0xf5,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x6,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("reference type, unauthorized, bool", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: ReferenceStaticType{
				Authorized: false,
				Type:       PrimitiveStaticTypeBool,
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagReferenceStaticType,
			// array, 2 items follow
			0x82,
			// false
			0xf4,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x6,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("dictionary, bool, string", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: DictionaryStaticType{
				KeyType:   PrimitiveStaticTypeBool,
				ValueType: PrimitiveStaticTypeString,
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagDictionaryStaticType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x6,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x8,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("restricted", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: &RestrictedStaticType{
				Type: CompositeStaticType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "S",
				},
				Restrictions: []InterfaceStaticType{
					{
						Location:            utils.TestLocation,
						QualifiedIdentifier: "I1",
					},
					{
						Location:            utils.TestLocation,
						QualifiedIdentifier: "I2",
					},
				},
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagRestrictedStaticType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, CBORTagCompositeStaticType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, CBORTagStringLocation,
			// UTF-8 string, length 4
			0x64,
			// t, e, s, t
			0x74, 0x65, 0x73, 0x74,
			// UTF-8 string, length 1
			0x61,
			// S
			0x53,
			// array, length 2
			0x82,
			// tag
			0xd8, CBORTagInterfaceStaticType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, CBORTagStringLocation,
			// UTF-8 string, length 4
			0x64,
			// t, e, s, t
			0x74, 0x65, 0x73, 0x74,
			// UTF-8 string, length 2
			0x62,
			// I1
			0x49, 0x31,
			// tag
			0xd8, CBORTagInterfaceStaticType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, CBORTagStringLocation,
			// UTF-8 string, length 4
			0x64,
			// t, e, s, t
			0x74, 0x65, 0x73, 0x74,
			// UTF-8 string, length 2
			0x62,
			// I2
			0x49, 0x32,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("capability, none", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type:       CapabilityStaticType{},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagCapabilityStaticType,
			// null
			0xf6,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("capability, primitive, bool", func(t *testing.T) {

		t.Parallel()

		value := LinkValue{
			TargetPath: publicPathValue,
			Type: CapabilityStaticType{
				BorrowType: PrimitiveStaticTypeBool,
			},
		}

		//nolint:gocritic
		encoded := append(
			expectedLinkEncodingPrefix[:],
			// tag
			0xd8, CBORTagCapabilityStaticType,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			0x6,
		)

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})
}

func TestEncodeDecodeTypeValue(t *testing.T) {

	t.Parallel()

	t.Run("primitive, Bool", func(t *testing.T) {

		t.Parallel()

		value := TypeValue{
			Type: ConvertSemaToPrimitiveStaticType(sema.BoolType),
		}

		encoded := []byte{
			// tag
			0xd8, CBORTagTypeValue,
			// array, 1 items follow
			0x81,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			// positive integer 0
			0x6,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("primitive, Int", func(t *testing.T) {

		t.Parallel()

		value := TypeValue{
			Type: ConvertSemaToPrimitiveStaticType(sema.IntType),
		}

		encoded := []byte{
			// tag
			0xd8, CBORTagTypeValue,
			// array, 1 items follow
			0x81,
			// tag
			0xd8, CBORTagPrimitiveStaticType,
			// positive integer 36
			0x18, 0x24,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
			},
		)
	})

	t.Run("without static type", func(t *testing.T) {

		t.Parallel()

		value := TypeValue{
			Type: nil,
		}

		encoded := []byte{
			// tag
			0xd8, CBORTagTypeValue,
			// array, 1 items follow
			0x81,
			// nil
			0xf6,
		}

		testEncodeDecode(t,
			encodeDecodeTest{
				value:   value,
				encoded: encoded,
				// type values without a static type are not semantically equal,
				// so check deep equality
				deepEquality: true,
			},
		)
	})
}
