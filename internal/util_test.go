package root

import (
	"reflect"
	"testing"
)

func TestGenItem(t *testing.T) {
	type args struct {
		inputFileName string
		cases         []string
	}
	tests := []struct {
		name    string
		args    args
		want    *APISpec
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenItem(tt.args.inputFileName, tt.args.cases)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenScenario(t *testing.T) {
	type args struct {
		apiSpec        *APISpec
		outputFileName string
		opts           []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenScenario(tt.args.apiSpec, tt.args.outputFileName, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("GenScenario() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetResponse(t *testing.T) {
	type args struct {
		url    string
		query  any
		method string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResponse(tt.args.url, tt.args.query, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddParam(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		want      *map[string]addParam
		wantErr   bool
	}{
		{
			name:      "正常系: クエリパラメーター付き",
			inputFile: "testdata/valid.csv",
			want: &map[string]addParam{
				"/v1/users": {Method: "GET", Query: map[string]interface{}{"userId": "123", "role": "admin"}, Body: ""},
				"/v1/items": {Method: "POST", Query: map[string]interface{}{}, Body: `{"name":"test"}`},
			},
			wantErr: false,
		},
		{
			name:      "正常系: 列が2つのみ (body なし)",
			inputFile: "testdata/two_columns.csv",
			want: &map[string]addParam{
				"/v1/users": {Method: "GET", Query: map[string]interface{}{}, Body: ""},
			},
			wantErr: false,
		},
		{
			name:      "正常系: 列が1つの不正行はスキップ",
			inputFile: "testdata/missing_columns.csv",
			want: &map[string]addParam{
				"/v1/users": {Method: "GET", Query: map[string]interface{}{}, Body: ""},
				"/v1/items": {Method: "POST", Query: map[string]interface{}{}, Body: `{"name":"test"}`},
			},
			wantErr: false,
		},
		{
			name:      "正常系: =のないクエリパラメーター",
			inputFile: "testdata/no_query_value.csv",
			want: &map[string]addParam{
				"/v1/users": {Method: "GET", Query: map[string]interface{}{"flag": ""}, Body: ""},
			},
			wantErr: false,
		},
		{
			name:      "異常系: ファイルが存在しない",
			inputFile: "testdata/not_exist.csv",
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddParam(tt.inputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
