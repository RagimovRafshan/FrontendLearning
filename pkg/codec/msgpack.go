package codec

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
)

// Encoder пул для кодирования
var encoderPool = sync.Pool{
	New: func() interface{} {
		return msgpack.NewEncoder(nil)
	},
}

// Decoder пул для декодирования
var decoderPool = sync.Pool{
	New: func() interface{} {
		return msgpack.NewDecoder(nil)
	},
}

// Encode кодирует структуру в бинарный формат MessagePack
func Encode(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := encoderPool.Get().(*msgpack.Encoder)
	defer encoderPool.Put(enc)

	enc.Reset(buf)
	// Оптимизация для высокой нагрузки
	enc.UseCompactInts(true)
	enc.UseCompactFloats(true)
	
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("encode error: %w", err)
	}
	
	return buf.Bytes(), nil
}

// Decode декодирует бинарные данные MessagePack в структуру
func Decode(data []byte, v interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data")
	}

	dec := decoderPool.Get().(*msgpack.Decoder)
	defer decoderPool.Put(dec)

	dec.Reset(bytes.NewReader(data))
	// Строгий режим для безопасности
	dec.SetCustomStructTag("msgpack")
	
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}
	
	return nil
}

// EncodeToBuffer кодирует напрямую в буфер (для оптимизации памяти)
func EncodeToBuffer(v interface{}, buf *bytes.Buffer) error {
	enc := encoderPool.Get().(*msgpack.Encoder)
	defer encoderPool.Put(enc)

	enc.Reset(buf)
	enc.UseCompactInts(true)
	enc.UseCompactFloats(true)
	
	return enc.Encode(v)
}

// DecodeFromReader декодирует из reader (для потоковой обработки)
func DecodeFromReader(r *bytes.Reader, v interface{}) error {
	dec := decoderPool.Get().(*msgpack.Decoder)
	defer decoderPool.Put(dec)

	dec.Reset(r)
	dec.SetCustomStructTag("msgpack")
	
	return dec.Decode(v)
}
