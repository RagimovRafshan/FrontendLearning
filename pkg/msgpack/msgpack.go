package msgpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"
	"unsafe"
)

const (
	mpPositiveFixint uint8 = 0x00
	mpNegativeFixint uint8 = 0xe0
	mpFixArray       uint8 = 0x90
	mpFixStr         uint8 = 0xa0
	mpNil            uint8 = 0xc0
	mpFalse          uint8 = 0xc2
	mpTrue           uint8 = 0xc3
	mpFloat64        uint8 = 0xcb
	mpUint8          uint8 = 0xcc
	mpUint16         uint8 = 0xcd
	mpUint32         uint8 = 0xce
	mpUint64         uint8 = 0xcf
	mpInt8           uint8 = 0xd0
	mpInt16          uint8 = 0xd1
	mpInt32          uint8 = 0xd2
	mpInt64          uint8 = 0xd3
	mpStr16          uint8 = 0xda
	mpStr32          uint8 = 0xdb
	mpArray16        uint8 = 0xdc
	mpArray32        uint8 = 0xdd
	mpMap16          uint8 = 0xde
	mpMap32          uint8 = 0xdf
	mpBin8           uint8 = 0xc4
	mpBin16          uint8 = 0xc5
	mpBin32          uint8 = 0xc6
	mpExt8           uint8 = 0xc7
	mpExt16          uint8 = 0xc8
	mpExt32          uint8 = 0xc9
	mpFixExt1        uint8 = 0xd4
	mpFixExt2        uint8 = 0xd5
	mpFixExt4        uint8 = 0xd6
	mpFixExt8        uint8 = 0xd7
	mpFixExt16       uint8 = 0xd8
)

// Encoder кодирует данные в MessagePack
type Encoder struct {
	buf *bytes.Buffer
}

// NewEncoder создает новый encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{buf: &bytes.Buffer{}}
}

// Encode кодирует значение в MessagePack
func (e *Encoder) Encode(v interface{}) error {
	if err := e.encodeValue(reflect.ValueOf(v)); err != nil {
		return err
	}
	_, err := e.buf.WriteTo(nil)
	return err
}

// EncodeBytes кодирует значение и возвращает байты
func (e *Encoder) EncodeBytes(v interface{}) ([]byte, error) {
	e.buf.Reset()
	if err := e.encodeValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

func (e *Encoder) encodeValue(v reflect.Value) error {
	if !v.IsValid() {
		e.buf.WriteByte(mpNil)
		return nil
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			e.buf.WriteByte(mpTrue)
		} else {
			e.buf.WriteByte(mpFalse)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.encodeInt(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeUint(v.Uint())
	case reflect.Float32, reflect.Float64:
		e.buf.WriteByte(mpFloat64)
		binary.Write(e.buf, binary.BigEndian, v.Float())
	case reflect.String:
		return e.encodeString(v.String())
	case reflect.Slice, reflect.Array:
		return e.encodeSlice(v)
	case reflect.Map:
		return e.encodeMap(v)
	case reflect.Struct:
		return e.encodeStruct(v)
	case reflect.Ptr:
		if v.IsNil() {
			e.buf.WriteByte(mpNil)
			return nil
		}
		return e.encodeValue(v.Elem())
	case reflect.Interface:
		if v.IsNil() {
			e.buf.WriteByte(mpNil)
			return nil
		}
		return e.encodeValue(v.Elem())
	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}
	return nil
}

func (e *Encoder) encodeInt(i int64) error {
	if i >= 0 {
		if i <= 127 {
			e.buf.WriteByte(uint8(i))
		} else if i <= 255 {
			e.buf.WriteByte(mpUint8)
			e.buf.WriteByte(uint8(i))
		} else if i <= 65535 {
			e.buf.WriteByte(mpUint16)
			binary.Write(e.buf, binary.BigEndian, uint16(i))
		} else if i <= 4294967295 {
			e.buf.WriteByte(mpUint32)
			binary.Write(e.buf, binary.BigEndian, uint32(i))
		} else {
			e.buf.WriteByte(mpUint64)
			binary.Write(e.buf, binary.BigEndian, uint64(i))
		}
	} else {
		if i >= -32 {
			e.buf.WriteByte(uint8(i))
		} else if i >= -128 {
			e.buf.WriteByte(mpInt8)
			e.buf.WriteByte(uint8(i))
		} else if i >= -32768 {
			e.buf.WriteByte(mpInt16)
			binary.Write(e.buf, binary.BigEndian, int16(i))
		} else if i >= -2147483648 {
			e.buf.WriteByte(mpInt32)
			binary.Write(e.buf, binary.BigEndian, int32(i))
		} else {
			e.buf.WriteByte(mpInt64)
			binary.Write(e.buf, binary.BigEndian, i)
		}
	}
	return nil
}

func (e *Encoder) encodeUint(u uint64) error {
	if u <= 127 {
		e.buf.WriteByte(uint8(u))
	} else if u <= 255 {
		e.buf.WriteByte(mpUint8)
		e.buf.WriteByte(uint8(u))
	} else if u <= 65535 {
		e.buf.WriteByte(mpUint16)
		binary.Write(e.buf, binary.BigEndian, uint16(u))
	} else if u <= 4294967295 {
		e.buf.WriteByte(mpUint32)
		binary.Write(e.buf, binary.BigEndian, uint32(u))
	} else {
		e.buf.WriteByte(mpUint64)
		binary.Write(e.buf, binary.BigEndian, u)
	}
	return nil
}

func (e *Encoder) encodeString(s string) error {
	b := []byte(s)
	length := len(b)
	
	if length <= 31 {
		e.buf.WriteByte(mpFixStr | uint8(length))
	} else if length <= 65535 {
		e.buf.WriteByte(mpStr16)
		binary.Write(e.buf, binary.BigEndian, uint16(length))
	} else {
		e.buf.WriteByte(mpStr32)
		binary.Write(e.buf, binary.BigEndian, uint32(length))
	}
	
	e.buf.Write(b)
	return nil
}

func (e *Encoder) encodeSlice(v reflect.Value) error {
	length := v.Len()
	
	if length == 0 {
		e.buf.WriteByte(mpFixArray)
		return nil
	}
	
	if length <= 15 {
		e.buf.WriteByte(mpFixArray | uint8(length))
	} else if length <= 65535 {
		e.buf.WriteByte(mpArray16)
		binary.Write(e.buf, binary.BigEndian, uint16(length))
	} else {
		e.buf.WriteByte(mpArray32)
		binary.Write(e.buf, binary.BigEndian, uint32(length))
	}
	
	for i := 0; i < length; i++ {
		if err := e.encodeValue(v.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeMap(v reflect.Value) error {
	keys := v.MapKeys()
	length := len(keys)
	
	if length == 0 {
		e.buf.WriteByte(mpFixMap())
		return nil
	}
	
	if length <= 15 {
		e.buf.WriteByte(mpFixMap() | uint8(length))
	} else if length <= 65535 {
		e.buf.WriteByte(mpMap16)
		binary.Write(e.buf, binary.BigEndian, uint16(length))
	} else {
		e.buf.WriteByte(mpMap32)
		binary.Write(e.buf, binary.BigEndian, uint32(length))
	}
	
	for _, key := range keys {
		if err := e.encodeValue(key); err != nil {
			return err
		}
		if err := e.encodeValue(v.MapIndex(key)); err != nil {
			return err
		}
	}
	return nil
}

func mpFixMap() uint8 {
	return 0x80
}

func (e *Encoder) encodeStruct(v reflect.Value) error {
	t := v.Type()
	fields := make(map[string]reflect.Value)
	
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("msgpack")
		
		if tag == "-" {
			continue
		}
		
		if tag == "" {
			tag = field.Name
		}
		
		fields[tag] = v.Field(i)
	}
	
	length := len(fields)
	if length <= 15 {
		e.buf.WriteByte(mpFixMap() | uint8(length))
	} else if length <= 65535 {
		e.buf.WriteByte(mpMap16)
		binary.Write(e.buf, binary.BigEndian, uint16(length))
	} else {
		e.buf.WriteByte(mpMap32)
		binary.Write(e.buf, binary.BigEndian, uint32(length))
	}
	
	for key, value := range fields {
		if err := e.encodeString(key); err != nil {
			return err
		}
		if err := e.encodeValue(value); err != nil {
			return err
		}
	}
	return nil
}

// Decoder декодирует данные из MessagePack
type Decoder struct {
	data []byte
	pos  int
}

// NewDecoder создает новый decoder
func NewDecoder(data []byte) *Decoder {
	return &Decoder{data: data, pos: 0}
}

// Decode декодирует данные в указанное значение
func (d *Decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("decode into non-pointer or nil pointer")
	}
	
	return d.decodeValue(rv.Elem())
}

func (d *Decoder) decodeValue(v reflect.Value) error {
	if d.pos >= len(d.data) {
		return fmt.Errorf("unexpected end of data")
	}
	
	b := d.data[d.pos]
	
	switch {
	case b == mpNil:
		d.pos++
		v.Set(reflect.Zero(v.Type()))
		return nil
	case b == mpFalse:
		d.pos++
		if v.Kind() == reflect.Bool {
			v.SetBool(false)
		}
		return nil
	case b == mpTrue:
		d.pos++
		if v.Kind() == reflect.Bool {
			v.SetBool(true)
		}
		return nil
	case b >= mpPositiveFixint && b <= 0x7f:
		d.pos++
		return d.setInt(v, int64(b))
	case b >= mpNegativeFixint:
		d.pos++
		return d.setInt(v, int64(int8(b)))
	case b == mpUint8:
		d.pos++
		if d.pos+1 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := int64(d.data[d.pos])
		d.pos++
		return d.setInt(v, val)
	case b == mpUint16:
		d.pos++
		if d.pos+2 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := int64(binary.BigEndian.Uint16(d.data[d.pos : d.pos+2]))
		d.pos += 2
		return d.setInt(v, val)
	case b == mpUint32:
		d.pos++
		if d.pos+4 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := int64(binary.BigEndian.Uint32(d.data[d.pos : d.pos+4]))
		d.pos += 4
		return d.setInt(v, val)
	case b == mpUint64:
		d.pos++
		if d.pos+8 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := int64(binary.BigEndian.Uint64(d.data[d.pos : d.pos+8]))
		d.pos += 8
		return d.setInt(v, val)
	case b == mpInt8:
		d.pos++
		if d.pos+1 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := int64(int8(d.data[d.pos]))
		d.pos++
		return d.setInt(v, val)
	case b == mpInt16:
		d.pos++
		if d.pos+2 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := int64(int16(binary.BigEndian.Uint16(d.data[d.pos : d.pos+2])))
		d.pos += 2
		return d.setInt(v, val)
	case b == mpInt32:
		d.pos++
		if d.pos+4 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := int64(int32(binary.BigEndian.Uint32(d.data[d.pos : d.pos+4])))
		d.pos += 4
		return d.setInt(v, val)
	case b == mpInt64:
		d.pos++
		if d.pos+8 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := int64(binary.BigEndian.Uint64(d.data[d.pos : d.pos+8]))
		d.pos += 8
		return d.setInt(v, val)
	case b == mpFloat64:
		d.pos++
		if d.pos+8 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		val := binary.BigEndian.Uint64(d.data[d.pos : d.pos+8])
		d.pos += 8
		if v.Kind() == reflect.Float64 {
			v.SetFloat(math.Float64frombits(val))
		}
		return nil
	case b >= mpFixStr && b <= 0xbf:
		length := int(b & 0x1f)
		d.pos++
		return d.decodeString(v, length)
	case b == mpStr16:
		d.pos++
		if d.pos+2 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		length := int(binary.BigEndian.Uint16(d.data[d.pos : d.pos+2]))
		d.pos += 2
		return d.decodeString(v, length)
	case b == mpStr32:
		d.pos++
		if d.pos+4 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		length := int(binary.BigEndian.Uint32(d.data[d.pos : d.pos+4]))
		d.pos += 4
		return d.decodeString(v, length)
	case b >= mpFixArray && b <= 0x9f:
		length := int(b & 0x0f)
		d.pos++
		return d.decodeSlice(v, length)
	case b == mpArray16:
		d.pos++
		if d.pos+2 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		length := int(binary.BigEndian.Uint16(d.data[d.pos : d.pos+2]))
		d.pos += 2
		return d.decodeSlice(v, length)
	case b == mpArray32:
		d.pos++
		if d.pos+4 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		length := int(binary.BigEndian.Uint32(d.data[d.pos : d.pos+4]))
		d.pos += 4
		return d.decodeSlice(v, length)
	case b >= 0x80 && b <= 0x8f:
		length := int(b & 0x0f)
		d.pos++
		return d.decodeMap(v, length)
	case b == mpMap16:
		d.pos++
		if d.pos+2 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		length := int(binary.BigEndian.Uint16(d.data[d.pos : d.pos+2]))
		d.pos += 2
		return d.decodeMap(v, length)
	case b == mpMap32:
		d.pos++
		if d.pos+4 > len(d.data) {
			return fmt.Errorf("unexpected end of data")
		}
		length := int(binary.BigEndian.Uint32(d.data[d.pos : d.pos+4]))
		d.pos += 4
		return d.decodeMap(v, length)
	default:
		return fmt.Errorf("unsupported messagepack type: 0x%02x", b)
	}
}

func (d *Decoder) setInt(v reflect.Value, i int64) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(i))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(float64(i))
	default:
		return fmt.Errorf("cannot assign int to %s", v.Kind())
	}
	return nil
}

func (d *Decoder) decodeString(v reflect.Value, length int) error {
	if d.pos+length > len(d.data) {
		return fmt.Errorf("unexpected end of data")
	}
	
	s := string(d.data[d.pos : d.pos+length])
	d.pos += length
	
	if v.Kind() == reflect.String {
		v.SetString(s)
	}
	return nil
}

func (d *Decoder) decodeSlice(v reflect.Value, length int) error {
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice, got %s", v.Kind())
	}
	
	slice := reflect.MakeSlice(v.Type(), length, length)
	
	for i := 0; i < length; i++ {
		if err := d.decodeValue(slice.Index(i)); err != nil {
			return err
		}
	}
	
	v.Set(slice)
	return nil
}

func (d *Decoder) decodeMap(v reflect.Value, length int) error {
	if v.Kind() != reflect.Map {
		return fmt.Errorf("expected map, got %s", v.Kind())
	}
	
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	
	for i := 0; i < length; i++ {
		key := reflect.New(v.Type().Key()).Elem()
		val := reflect.New(v.Type().Elem()).Elem()
		
		if err := d.decodeValue(key); err != nil {
			return err
		}
		if err := d.decodeValue(val); err != nil {
			return err
		}
		
		v.SetMapIndex(key, val)
	}
	return nil
}

// Marshal маршалиет значение в MessagePack
func Marshal(v interface{}) ([]byte, error) {
	enc := NewEncoder(nil)
	return enc.EncodeBytes(v)
}

// Unmarshal демаршалиет данные из MessagePack
func Unmarshal(data []byte, v interface{}) error {
	dec := NewDecoder(data)
	return dec.Decode(v)
}
