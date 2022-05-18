package command

import (
	"reflect"
	"strings"
	"testing"
)

var (
	testString string = "this is a test"
)

func TestCreateData(t *testing.T) {
	type args struct {
		lst    dataCreatorListener
		name   string
		params []string
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
	}{
		{
			name: "Pass: Should return the data from stdout",
			args: args{
				&mockDataCreatorListener{},
				"echo",
				strings.Split(testString, " "),
			},
			wantData: []byte(testString),
		},
		{
			name: "Pass: No commands should return empty byteslice",
			args: args{
				&mockDataCreatorListener{},
				"",
				[]string{},
			},
			wantData: []byte{},
		},
		{
			name: "Pass: No listener should return empty data byteslice",
			args: args{
				nil,
				"echo",
				[]string{"hello"},
			},
			wantData: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateData(tt.args.lst, tt.args.name, tt.args.params...)

			if ml, ok := tt.args.lst.(*mockDataCreatorListener); ok {
				ml.data = ml.data[:len(tt.wantData)]
				if !reflect.DeepEqual(ml.data, tt.wantData) {
					t.Errorf("wanted '%s' but got '%s'\n", tt.wantData, ml.data)
				}
			}
		})
	}
}

type mockDataCreatorListener struct {
	data []byte
}

func (c *mockDataCreatorListener) OnDataCreation(data []byte) {
	c.data = data
}
