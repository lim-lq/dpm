package metadata

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type MapStr map[string]interface{}

type DpmTime struct {
	time.Time
}

func (t *DpmTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"2006-01-02 15:04:05"`)), nil
}

func (t *DpmTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	parsed, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.UTC)
	if err == nil {
		*t = DpmTime{parsed}
		return nil
	}

	parsed, err = time.Parse(time.RFC3339, strings.Trim(string(data), "\""))
	if err == nil {
		*t = DpmTime{parsed}
		return nil
	}

	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		*t = DpmTime{time.Unix(timestamp, 0)}
	}
	return fmt.Errorf("parse unknown time format: %s, %v", data, err)
}

func (t *DpmTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsonx.Time(t.Time).MarshalBSONValue()
}

func (t *DpmTime) UnmarshalBSONValue(tType bsontype.Type, raw []byte) error {
	switch tType {
	case bsontype.Timestamp:
		return bson.Unmarshal(raw, &t.Time)
	case bsontype.Double:
		rv := bson.RawValue{Type: bsontype.Double, Value: raw}
		if dt, ok := rv.DoubleOK(); ok {
			t.Time = time.Unix(int64(dt/1000), int64(uint64(dt)%1000*1000000))
			return nil
		}
		return nil
	case bsontype.DateTime:
		rv := bson.RawValue{Type: bsontype.DateTime, Value: raw}
		t.Time = rv.Time()
		return nil
	case bsontype.String:
		rawStr := bson.RawValue{Type: bsontype.String, Value: raw}
		strTime := strings.TrimSpace(strings.Trim(rawStr.String(), "\""))
		vTime, err := time.Parse(time.RFC3339Nano, strTime)
		if err == nil {
			t.Time = vTime
			return nil
		}
		vTime, err = time.Parse(time.RFC3339, strTime)
		if err == nil {
			t.Time = vTime
			return nil
		}
		return fmt.Errorf("cannot decode %v into a metadata.DpmTime, err - %v", bsontype.String, err)
	}
	return bson.Unmarshal(raw, t)
}
