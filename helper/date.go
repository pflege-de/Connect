package helper

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack"
)

// Date is a Salesforce Date
type Date time.Time

const DateFormatSF = "2005-10-08T01:02:03Z"

// Implement Unmarshaler interface
func (t *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	tt, err := time.Parse(DateFormatSF, s)
	if err != nil {
		return err
	}
	*t = Date(tt)
	return nil
}

// Implement Marshaler interface
func (t *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*t))
}

func (t *Date) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.Encode(time.Time(*t))
}

func (t *Date) DecodeMsgpack(dec *msgpack.Decoder) error {
	var tm time.Time
	err := dec.Decode(&tm)
	if err != nil {
		return err
	}
	*t = Date(tm)
	return nil
}

func (t *Date) Format(s string) string {
	tt := time.Time(*t)
	return tt.Format(s)
}

func (t *Date) FormatSF() string {
	return t.Format(DateFormatSF)
}

func (t *Date) Copy() *Date {
	out := Date(time.Time(*t).Add(0 * time.Second))
	return &out
}
