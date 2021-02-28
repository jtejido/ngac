package obligations

import (
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func validateSchema(t *testing.T, schema, file string) {
	path, err := filepath.Abs(schema)
	if err != nil {
		t.Fatalf("cannot open schema: %v", err)
	}

	fp, err := filepath.Abs(file)
	if err != nil {
		t.Fatalf("cannot open file: %v", err)
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file:///" + path)
	documentLoader := gojsonschema.NewReferenceLoader("file:///" + fp)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Fatalf("cannot validate data: %v", err)
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			t.Errorf("- %s\n", desc)
		}
	}
}

func TestLabel(t *testing.T) {
	validateSchema(t, "../../api/obligations.json", "test.json")

	var b Obligation
	yamlFile, err := ioutil.ReadFile("test.json")
	if err != nil {
		t.Fatal("no yaml file found")
	}
	err = json.Unmarshal(yamlFile, &b)

	if err != nil {
		t.Fatalf("cannot unmarshal data: %v", err)
	}

	exp := "OA"
	if b.Rules[0].ResponsePattern.Actions[0].(*FunctionAction).Function.Args[3].Value != "OA" {
		t.Errorf("incorrect name: got %s, expected %s", b.Rules[0].ResponsePattern.Actions[0].(*FunctionAction).Function.Args[3].Value, exp)
	}
}
