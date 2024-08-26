package cobblerclient

import (
	"reflect"
	"testing"
)

func TestVersion(t *testing.T) {
	c := createStubHTTPClientSingle(t, "version")

	res, err := c.Version()
	FailOnError(t, err)
	if res != 3.4 {
		t.Errorf("Wrong version returned.")
	}
}

func TestExtendedVersion(t *testing.T) {
	c := createStubHTTPClientSingle(t, "extended-version")
	expectedResult := ExtendedVersion{
		Gitdate:      "Mon Jun 13 16:13:33 2022 +0200",
		Gitstamp:     "0e20f01b",
		Builddate:    "Mon Jun 27 06:34:23 2022",
		Version:      "3.4.0",
		VersionTuple: []int{3, 4, 0},
	}

	result, err := c.ExtendedVersion()
	FailOnError(t, err)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Result from 'extended_version' did not match expected result.")
	}
}

func TestCobblerVersion_GreaterThan(t *testing.T) {
	type fields struct {
		Major int
		Minor int
		Patch int
	}
	type args struct {
		otherVersion *CobblerVersion
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "greater",
			fields: fields{
				Major: 3,
				Minor: 3,
				Patch: 4,
			},
			args: args{
				otherVersion: &CobblerVersion{
					Major: 3,
					Minor: 3,
					Patch: 3,
				},
			},
			want: true,
		},
		{
			name: "equal",
			fields: fields{
				Major: 3,
				Minor: 3,
				Patch: 3,
			},
			args: args{
				otherVersion: &CobblerVersion{
					Major: 3,
					Minor: 3,
					Patch: 3,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CobblerVersion{
				Major: tt.fields.Major,
				Minor: tt.fields.Minor,
				Patch: tt.fields.Patch,
			}
			if got := cv.GreaterThan(tt.args.otherVersion); got != tt.want {
				t.Errorf("GreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCobblerVersion_LessThan(t *testing.T) {
	type fields struct {
		Major int
		Minor int
		Patch int
	}
	type args struct {
		otherVersion *CobblerVersion
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "less",
			fields: fields{
				Major: 3,
				Minor: 3,
				Patch: 3,
			},
			args: args{
				otherVersion: &CobblerVersion{
					Major: 3,
					Minor: 3,
					Patch: 4,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CobblerVersion{
				Major: tt.fields.Major,
				Minor: tt.fields.Minor,
				Patch: tt.fields.Patch,
			}
			if got := cv.LessThan(tt.args.otherVersion); got != tt.want {
				t.Errorf("LessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCobblerVersion_Equal(t *testing.T) {
	type fields struct {
		Major int
		Minor int
		Patch int
	}
	type args struct {
		otherVersion *CobblerVersion
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "equal",
			fields: fields{
				Major: 3,
				Minor: 3,
				Patch: 3,
			},
			args: args{
				otherVersion: &CobblerVersion{
					Major: 3,
					Minor: 3,
					Patch: 3,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CobblerVersion{
				Major: tt.fields.Major,
				Minor: tt.fields.Minor,
				Patch: tt.fields.Patch,
			}
			if got := cv.Equal(tt.args.otherVersion); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCobblerVersion_NotEqual(t *testing.T) {
	type fields struct {
		Major int
		Minor int
		Patch int
	}
	type args struct {
		otherVersion *CobblerVersion
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "not equal",
			fields: fields{
				Major: 3,
				Minor: 3,
				Patch: 3,
			},
			args: args{
				otherVersion: &CobblerVersion{
					Major: 3,
					Minor: 3,
					Patch: 4,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CobblerVersion{
				Major: tt.fields.Major,
				Minor: tt.fields.Minor,
				Patch: tt.fields.Patch,
			}
			if got := cv.NotEqual(tt.args.otherVersion); got != tt.want {
				t.Errorf("NotEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
