package sbomreport

import (
	"encoding/json"
	"errors"
	"fmt"
)

type SbomReport struct {
	rawJSON []byte
	bom     []byte
}

func New(rawJSON []byte) (*SbomReport, error) {
	bom, err := getBOM(rawJSON)
	if err != nil {
		return nil, err
	}
	return &SbomReport{
		rawJSON: rawJSON,
		bom:     bom,
	}, nil
}

func (s *SbomReport) ToMap() (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(s.rawJSON, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *SbomReport) BOM() []byte {
	return s.bom
}

func getBOM(rawJSON []byte) ([]byte, error) {
	var data map[string]interface{}

	if err := json.Unmarshal(rawJSON, &data); err != nil {
		return nil, err
	}

	kind, ok := data["kind"].(string)
	if !ok || kind != "SbomReport" {
		return nil, errors.New("kind is not SbomReport")
	}

	apiVersion, ok := data["apiVersion"].(string)
	if !ok {
		return nil, fmt.Errorf("apiVersion %q is not found", apiVersion)
	}

	report, ok := data["report"].(map[string]interface{})
	if !ok {
		return nil, errors.New("report is not found")
	}

	bom, ok := report["components"].(map[string]interface{})
	if !ok {
		return nil, errors.New("bom is not found")
	}

	jsonBytes, err := json.Marshal(bom)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
