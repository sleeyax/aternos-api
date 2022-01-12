package aternos_api

import (
	"encoding/json"
	"testing"
)

type DataContent struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

type Example struct {
	Id   int  `json:"id"`
	Data Data `json:"data"`
}

func TestData(t *testing.T) {
	raw := `{"id": 0, "data": {"foo": "hello world!", "bar": 123}}`

	var example Example

	if err := json.Unmarshal([]byte(raw), &example); err != nil {
		t.Fatal(err)
	}

	var dataContent DataContent
	if err := json.Unmarshal(example.Data.ContentBytes, &dataContent); err != nil {
		t.Fatal(err)
	}

	if dataContent.Foo != "hello world!" {
		t.Fail()
	}
	if dataContent.Bar != 123 {
		t.Fail()
	}
}

func TestData_MarshalJSON(t *testing.T) {
	data := &Data{
		Content: "hello world!",
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	str := string(serialized)

	if str != `"hello world!"` {
		t.FailNow()
	}

	data.Content = `{"foo": "bar"}`

	serialized, err = json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	str = string(serialized)

	if str != `{"foo":"bar"}` {
		t.FailNow()
	}
}
