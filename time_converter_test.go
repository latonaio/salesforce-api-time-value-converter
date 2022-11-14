package salesforce_api_time_value_converter

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertToTime(t *testing.T) {
	type args struct {
		salesforceTime string
	}
	type testStr struct {
		name string
		args args
		want time.Time
	}
	tests := []testStr{
		func() testStr {
			now := time.Now()
			return testStr{
				name: "OK now time",
				args: args{
					salesforceTime: fmt.Sprintf(`\/Date(%d)\/`, now.UnixMilli()),
				},
				want: now,
			}
		}(),
		func() testStr {
			now := time.Now()
			return testStr{
				name: "OK now time",
				args: args{
					salesforceTime: fmt.Sprintf(`/Date(%d)/`, now.UnixMilli()),
				},
				want: now,
			}
		}(),
		func() testStr {
			now := time.Now()
			return testStr{
				name: "OK now time",
				args: args{
					salesforceTime: fmt.Sprintf(`/Date(%d+0000)/`, now.UnixMilli()),
				},
				want: now,
			}
		}(),
		func() testStr {
			return testStr{
				name: "OK now time",
				args: args{
					salesforceTime: `\/Date(1642757478000)\/`,
				},
				want: time.Date(2022, 1, 21, 9, 31, 18, 0, time.UTC),
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToTimeFormat(tt.args.salesforceTime); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want.UnixMilli(), got.UnixMilli(), "not same time")
			}
		})
	}
}
