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
	type args struct {
		intpufile string
	}
	tests := []struct {
		name    string
		args    args
		want    *map[string]addParam
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddParam(tt.args.intpufile)
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
