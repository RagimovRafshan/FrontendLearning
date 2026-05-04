package msgpack

import (
"bytes"
"encoding/json"
"io"
"net/http"

"github.com/vmihailenco/msgpack/v5"
)

// Encode кодирует структуру в MessagePack
func Encode(v interface{}) ([]byte, error) {
return msgpack.Marshal(v)
}

// Decode декодирует MessagePack в структуру
func Decode(data []byte, v interface{}) error {
return msgpack.Unmarshal(data, v)
}

// DecodeRequest декодирует тело запроса (MessagePack или JSON fallback)
func DecodeRequest(r *http.Request, v interface{}) error {
contentType := r.Header.Get("Content-Type")

if contentType == "application/msgpack" || contentType == "application/x-msgpack" {
data, err := io.ReadAll(r.Body)
if err != nil {
return err
}
defer r.Body.Close()
return msgpack.Unmarshal(data, v)
}

// Fallback на JSON
decoder := json.NewDecoder(r.Body)
defer r.Body.Close()
return decoder.Decode(v)
}

// Respond отправляет ответ в формате MessagePack или JSON
func Respond(w http.ResponseWriter, status int, v interface{}) {
w.Header().Set("Content-Type", "application/msgpack")
w.WriteHeader(status)

data, err := msgpack.Marshal(v)
if err != nil {
// Fallback на JSON при ошибке
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(v)
return
}

w.Write(data)
}

// RespondJSON отправляет ответ в формате JSON
func RespondJSON(w http.ResponseWriter, status int, v interface{}) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)
json.NewEncoder(w).Encode(v)
}

// EncodeToBuffer кодирует в bytes.Buffer
func EncodeToBuffer(v interface{}) (*bytes.Buffer, error) {
data, err := msgpack.Marshal(v)
if err != nil {
return nil, err
}
return bytes.NewBuffer(data), nil
}

// Decoder создает новый декодер из io.Reader
func Decoder(r io.Reader) *msgpack.Decoder {
return msgpack.NewDecoder(r)
}

// Encoder создает новый энкодер для io.Writer
func Encoder(w io.Writer) *msgpack.Encoder {
return msgpack.NewEncoder(w)
}
