package go_yandex_checkout

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StrFloat64 struct {
	V float64
}

func (sf *StrFloat64) MarshalJSON() ([]byte, error) {
	quoted := strconv.Quote(fmt.Sprintf("%f", sf.V))
	return []byte(quoted), nil
}

func (sf *StrFloat64) UnmarshalJSON(bytes []byte) error {
	var v string
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		sf.V = 0
		return nil
	}

	sf.V = f
	return nil
}

type StrInt64 struct {
	V int64
}

func (si *StrInt64) MarshalJSON() ([]byte, error) {
	quoted := strconv.Quote(fmt.Sprintf("%d", si.V))
	return []byte(quoted), nil
}

func (si *StrInt64) UnmarshalJSON(bytes []byte) error {
	var v string
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}

	f, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		si.V = 0
		return nil
	}

	si.V = f
	return nil
}
