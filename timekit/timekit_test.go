package timekit

import (
	"reflect"
	"testing"
	"time"
)

func TestZeroTime(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		// TODO: Add test cases.
		{
			name: "测试utc的zero time",
			want: TimeStamp2UTCTime(CurrentTimeStamp()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZeroTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroTime() = %v, want %v", got, tt.want)
			}
		})
	}
}