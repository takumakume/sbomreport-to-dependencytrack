package template

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		sbomReport map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *Template
	}{
		{
			name: "success",
			args: args{
				sbomReport: map[string]interface{}{
					"kind":       "SbomReport",
					"apiVersion": "aquasecurity.github.io/v1alpha1",
					"report": map[string]interface{}{
						"components": map[string]interface{}{},
					},
				},
			},
			want: &Template{
				values: map[string]interface{}{
					"sbomReport": map[string]interface{}{
						"kind":       "SbomReport",
						"apiVersion": "aquasecurity.github.io/v1alpha1",
						"report": map[string]interface{}{
							"components": map[string]interface{}{},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.sbomReport); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplate_Render(t *testing.T) {
	type fields struct {
		values map[string]interface{}
	}
	type args struct {
		tmpl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				values: map[string]interface{}{
					"sbomReport": map[string]interface{}{
						"kind":       "SbomReport",
						"apiVersion": "aquasecurity.github.io/v1alpha1",
						"report": map[string]interface{}{
							"components": map[string]interface{}{},
						},
					},
				},
			},
			args: args{
				tmpl: "kind:[[ .sbomReport.kind ]]",
			},
			want:    "kind:SbomReport",
			wantErr: false,
		},
		{
			name: "no key",
			fields: fields{
				values: map[string]interface{}{},
			},
			args: args{
				tmpl: "kind:[[ .sbomReport.kind ]]",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Template{
				values: tt.fields.values,
			}
			got, err := tr.Render(tt.args.tmpl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Template.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}
