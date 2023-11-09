package datetime

import (
	"reflect"
	"testing"
	"time"
)

func TestDatetime_MarshalJSON(t *testing.T) {
	type fields struct {
		stringTime string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "succeed",
			fields: fields{
				stringTime: "2022-03-29T15:00:00Z",
			},
			want:    []byte("\"2022-03-29T15:00:00Z\""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, _ := time.Parse(RFC3339UTC, tt.fields.stringTime)
			s := DateTimeUTC(parsed)
			res, err := s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(tt.want, res) {
				t.Errorf("MarshalJSON() got = %[1]v(%[1]T), want %[2]v(%[2]T)", string(res), string(tt.want))
			}
		})
	}
}
