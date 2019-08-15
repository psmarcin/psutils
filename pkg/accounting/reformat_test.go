package accounting

import (
	"log"
	"testing"
	"time"
)

func Test_generateName(t *testing.T) {
	date, err := time.Parse("2006 January", "2018 September")
	if err != nil {
		log.Fatalf("Test err %+v", err)
	}

	type args struct {
		name string
		date time.Time
		types string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return correct file Name",
			args: args{
				name: "serwer vpn",
				date: date,
				types: "faktura",
			},
			want: "2018-09-serwer-vpn-faktura.pdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateName(tt.args.name, tt.args.types, tt.args.date); got != tt.want {
				t.Errorf("generateName() = %v, want %v", got, tt.want)
			}
		})
	}
}
