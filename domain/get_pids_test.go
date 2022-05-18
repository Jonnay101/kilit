package domain

import (
	"reflect"
	"strings"
	"testing"
)

func Test_getWordByIndex(t *testing.T) {
	type args struct {
		str string
		idx int
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantPanic bool
	}{
		{
			name: "ZeroVal: provided string is empty so should return empty",
			args: args{
				"",
				1,
			},
			want:      "",
			wantPanic: false,
		},
		{
			name: "ZeroVal: provided idx is lower than zero",
			args: args{
				"this bit is fine",
				-1,
			},
			want:      "",
			wantPanic: false,
		},
		{
			name: "Pass: should return the word at the provided idx",
			args: args{
				"Hello world, great to be here!",
				1,
			},
			want:      "world,",
			wantPanic: false,
		},
		{
			name: "Panic: idx is out of range",
			args: args{
				"only three words",
				3,
			},
			want:      "",
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if panErr := recover(); (panErr != nil) != tt.wantPanic {
					t.Errorf("wanted panic error to be %v but got %v\n", tt.wantPanic, panErr)
				}
			}()

			if got := getWordByIndex(tt.args.str, tt.args.idx); got != tt.want {
				t.Errorf("getWordByIndex() = '%v' but wanted '%v'", got, tt.want)
			}
		})
	}
}

func Test_readPIDFromString(t *testing.T) {
	type args struct {
		str string
		lst pidListener
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Exit: COMMAND in string - should exit gracefully",
			args: args{
				"COMMAND",
				nil,
			},
			want: "",
		},
		{
			name: "Panic: word count in string is less than 2 but panic handled",
			args: args{
				"One",
				nil,
			},
			want: "",
		},
		{
			name: "Pass: method completes but no lister is passed to handle PID",
			args: args{
				"progName 12345 otherField",
				nil,
			},
			want: "",
		},
		{
			name: "Pass: should find the PID from the line of data",
			args: args{
				"progName 12345 otherField",
				&mockPIDListener{},
			},
			want: "12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readPIDFromString(tt.args.str, tt.args.lst)

			if ml, ok := tt.args.lst.(*mockPIDListener); ok {
				if ml.result[0] != tt.want {
					t.Errorf("wanted PID to be '%s' but got '%s'\n", tt.want, ml.result)
				}
			}
		})
	}
}

type mockPIDListener struct {
	result []string
}

func (l *mockPIDListener) OnPIDFound(s string) {
	l.result = append(l.result, s)
}

var (
	testLines = []string{
		"a COMMAND line",
		"prog1 123 thing",
		"prog2 456 thing",
		"oops", // panic should be handled
		"",
		"prog3 789 thing",
	}
	wantPIDs = []string{
		"123",
		"456",
		"789",
	}
)

func Test_readAllPIDs(t *testing.T) {
	type args struct {
		lines []string
		lst   pidListener
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Pass: should find all the PIDs",
			args: args{
				testLines,
				&mockPIDListener{},
			},
			want: wantPIDs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readAllPIDs(tt.args.lines, tt.args.lst)

			if ml, ok := tt.args.lst.(*mockPIDListener); ok {
				for idx, pid := range ml.result {
					if pid != tt.want[idx] {
						t.Errorf("wanted pid at index %d to be '%s' but got '%s'\n", idx, tt.want[idx], pid)
					}
				}
			}
		})
	}
}

func Test_splitDataByLine(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Pass: should split data into lines",
			args: args{
				[]byte("hi, this is a \nstring split over 3\nlines of text"),
			},
			want: []string{
				"hi, this is a ",
				"string split over 3",
				"lines of text",
			},
		},
		{
			name: "Pass: should split data into lines",
			args: args{
				[]byte("this has no lines"),
			},
			want: []string{
				"this has no lines",
			},
		},
		{
			name: "Pass: should split data into lines",
			args: args{
				[]byte("\n\n\n"),
			},
			want: []string{"", "", ""},
		},
		{
			name: "Pass: should split data into lines",
			args: args{
				nil,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitDataByLine(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitDataByLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllPIDsFromData(t *testing.T) {
	testData := strings.Join(testLines, "\n")
	type args struct {
		data []byte
		lst  pidListener
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Pass: should get all PIDs from data",
			args: args{
				[]byte(testData),
				&mockPIDListener{},
			},
			want: wantPIDs,
		},
		{
			name: "Pass: should get all PIDs from data",
			args: args{
				nil,
				&mockPIDListener{},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetAllPIDsFromData(tt.args.data, tt.args.lst)

			if ml, ok := tt.args.lst.(*mockPIDListener); ok {
				for idx, pid := range ml.result {
					if pid != tt.want[idx] {
						t.Errorf("wanted pid at index %d to be '%s' but got '%s'\n", idx, tt.want[idx], pid)
					}
				}
			}
		})
	}
}

var (
	testData = []byte("hello, this is the data!")
)

type mockDataSource struct{}

func (s *mockDataSource) Data() []byte {
	return testData
}

var testString = "this is a test"
