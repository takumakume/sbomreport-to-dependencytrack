package sbomreport

import (
	"encoding/json"
	"errors"
	"fmt"
)

type SbomReport struct {
	rawJSON []byte
	bom     []byte
	verb    string
}

func New(rawJSON []byte) (*SbomReport, error) {
	bom, verb, err := getBOMAndVerb(rawJSON)
	if err != nil {
		return nil, err
	}
	return &SbomReport{
		rawJSON: rawJSON,
		verb:    verb,
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

func (s *SbomReport) ISVerbUpdate() bool {
	return s.verb == "update"
}

func getBOMAndVerb(rawJSON []byte) ([]byte, string, error) {
	verb := "update"
	var data map[string]interface{}

	if err := json.Unmarshal(rawJSON, &data); err != nil {
		return nil, verb, err
	}

	obj := data

	v, ok := data["verb"].(string)
	if ok {
		verb = v
		if operatorObject, ok := data["operatorObject"].(map[string]interface{}); ok {
			obj = operatorObject
		} else {
			return nil, verb, errors.New("operatorObject is not found")
		}
	}

	kind, ok := obj["kind"].(string)
	if !ok || kind != "SbomReport" {
		return nil, verb, errors.New("kind is not SbomReport")
	}

	apiVersion, ok := obj["apiVersion"].(string)
	if !ok {
		return nil, verb, fmt.Errorf("apiVersion %q is not found", apiVersion)
	}

	report, ok := obj["report"].(map[string]interface{})
	if !ok {
		return nil, verb, errors.New("report is not found")
	}

	bom, ok := report["components"].(map[string]interface{})
	if !ok {
		return nil, verb, errors.New("bom is not found")
	}

	jsonBytes, err := json.Marshal(bom)
	if err != nil {
		return nil, verb, err
	}

	return jsonBytes, verb, nil
}

func (s *SbomReport) MetadataComponentBomRef() (string, error) {
	var obj map[string]interface{}

	if err := json.Unmarshal(s.bom, &obj); err != nil {
		return "", err
	}

	metadata, ok := obj["metadata"].(map[string]interface{})
	if !ok {
		return "nil", errors.New("bom metadata is not found")
	}

	component, ok := metadata["component"].(map[string]interface{})
	if !ok {
		return "", errors.New("bom metadata.component is not found")
	}

	bomRef, ok := component["bom-ref"].(string)
	if !ok || bomRef == "" {
		return "", errors.New("bom metadata.component.bom-ref is not found or empty")
	}

	return bomRef, nil
}
