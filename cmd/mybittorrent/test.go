package main
import (
	"reflect"
	"testing"
)

func Test_decodeBencode(t *testing.T) {
	type args struct {
		bencodedString string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Test 1", args{"5:hello"}, "hello", false},
		{"Test 2", args{"10:hello12345"}, "hello12345", false},
		{"Test 3", args{"i123e"}, 123, false},
		{"Test 4", args{"i-123e"}, -123, false},
		{"Test 5", args{"l4:spam4:eggse"}, []interface{}{"spam", "eggs"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err, _ := decodeBencode(tt.args.bencodedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeBencode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeBencode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_decodeInteger(t *testing.T) {
	type args struct {
		bencodedString string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Test 1", args{"i123e"}, 123, false},
		{"Test 2", args{"i-123e"}, -123, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err , _ := decodeInt(tt.args.bencodedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeInteger() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_decodeList(t *testing.T) {
	type args struct {
		bencodedString string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Test 1", args{"l4:spam4:eggse"}, []interface{}{"spam", "eggs"}, false},
		{"Test 2", args{"ll4:spame4:eggsi123ee"}, []interface{}{[]interface{}{"spam"}, "eggs", 123}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err, _ := decodeLists(tt.args.bencodedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_decodeString(t *testing.T) {
	type args struct {
		bencodedString string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Test 1", args{"5:hello"}, "hello", false},
		{"Test 2", args{"10:hello12345"}, "hello12345", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err , _:= decodeString(tt.args.bencodedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeString() got = %v, want %v", got, tt.want)
			}
		})
	}
}