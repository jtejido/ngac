package obligations

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestLabel(t *testing.T) {
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
		t.Fatalf("incorrect name: got %s, expected %s", b.Rules[0].ResponsePattern.Actions[0].(*FunctionAction).Function.Args[3].Value, exp)
	}
}
