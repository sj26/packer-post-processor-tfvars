package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	awscommon "github.com/mitchellh/packer/builder/amazon/common"
	"github.com/mitchellh/packer/packer"
)

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packer.PostProcessor = new(PostProcessor)
}

func TestPostProcessor_Export(t *testing.T) {
	postProcessor := testPP(t)

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tempDir)

	previousDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	err = os.Chdir(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer os.Chdir(previousDir)

	_, _, err = postProcessor.PostProcess(testUi(), testArtifact())
	if err != nil {
		t.Error("Couldn't post-process artifact: ", err)
	}

	if _, err := os.Stat("packer.tfvars"); err != nil {
		t.Error("Expected packer.tfvars to be created")
	}

	contents, err := ioutil.ReadFile("packer.tfvars")
	if err != nil {
		t.Error(err)
	}

	var vars map[string]string
	err = json.Unmarshal(contents, &vars)
	if err != nil {
		t.Error("Couldn't unmarshal JSON: ", err)
	}

	if _, ok := vars["packer-east-test"]; !ok {
		t.Error("Expected a map key `packer-east-test`: ", vars)
	}

	if vars["packer-east-test"] != "foo" {
		t.Error("Expected map key `packer-east-test` to be `foo`: ", vars)
	}

	if _, ok := vars["packer-west-test"]; !ok {
		t.Error("Expected a map key `packer-west-test`: ", vars)
	}

	if vars["packer-west-test"] != "bar" {
		t.Error("Expected map key `packer-west-test` to be `bar`: ", vars)
	}
}

// Helpers

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"packer_build_name": "test",
	}
}

func testPP(t *testing.T) *PostProcessor {
	var p PostProcessor
	if err := p.Configure(testConfig()); err != nil {
		t.Fatalf("err: %s", err)
	}
	return &p
}

func testUi() *packer.BasicUi {
	return &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

func testArtifact() packer.Artifact {
	artifact := &awscommon.Artifact{
		BuilderIdValue: "mitchellh.amazonebs",
		Amis: map[string]string{
			"east": "foo",
			"west": "bar",
		},
	}

	return artifact
}
