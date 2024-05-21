package keyring

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var testTempDir string

// Create a test directory for the config files used in subsequent unit tests.
func setup() string {
	var err error
	// MkDirTemp creates a directory and returns the path to it. The directory
	// name is prefixed with the given pattern and suffixed with a random 
	// string generated at initialization.
	testTempDir, err = os.MkdirTemp(os.TempDir(), "lku-test-")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return testTempDir
}

// Remove the test directory and all of its contents.
func teardown(d string) {
	os.RemoveAll(d)
}

// Run setup code before and teardown code after all tests have completed.
func TestMain(m *testing.M) {
	dir := setup()
	code := m.Run()
	teardown(dir)
	os.Exit(code)
}


func TestAttributesFunc(t *testing.T) {
	// Test the Attributes function
	attrs := Attributes("TestApp")
	if attrs["Application"] != "TestApp" {
		t.Errorf("Attributes function did not return expected value")
	}
}

func TestConvertSecretToBase64Func(t *testing.T) {
	// Test the ConvertSecretToBase64 function
	var testJsonPath = "test.json"
	var testJson = `{"username":"user", "password":"pass"}`
	var testB64Str = "eyJteVByZWNpb3VzIjogIk9uZSBSaW5nIHRvIHJ1bGUgdGhlbSBhbGwsIE9uZSByaW5nIHRvIGZpbmQgdGhlbTsgT25lIHJpbmcgdG8gYnJpbmcgdGhlbSBhbGwgYW5kIGluIHRoZSBkYXJrbmVzcyBiaW5kIHRoZW0uIn0="

	// Create a test JSON file
	jsonFile, err := os.Create(filepath.Join(testTempDir, testJsonPath))
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}
	defer jsonFile.Close()
	
	_, err = jsonFile.WriteString(testJson)
	if err != nil {
		t.Errorf("Error writing to test file: %v", err)
	}

	// Test the ConvertSecretToBase64 function with a base64 encoded string
	validB64Str, err := ConvertSecretToBase64(testB64Str)
	if err != nil {
		t.Errorf("ConvertSecretToBase64 function failed: %v", err)
	}
	if validB64Str != testB64Str {
		t.Errorf("ConvertSecretToBase64 function did not return expected value")
	}

	// Test the ConvertSecretToBase64 function with a path to a JSON file
	jsonFileStr, err := ConvertSecretToBase64(jsonFile.Name())
	if err != nil {
		t.Errorf("ConvertSecretToBase64 function failed: %v", err)
	}

	if jsonFileStr != base64.StdEncoding.EncodeToString([]byte(testJson)) {
		t.Errorf("ConvertSecretToBase64 function did not return expected value")
	}

	// Test the ConvertSecretToBase64 function with a valid JSON string
	jsonStr, err := ConvertSecretToBase64(testJson)
	if err != nil {
		t.Errorf("ConvertSecretToBase64 function failed: %v", err)
	}
	if jsonStr != base64.StdEncoding.EncodeToString([]byte(testJson)) {
		t.Errorf("ConvertSecretToBase64 function did not return expected value")
	}

	// Test the ConvertSecretToBase64 function with an invalid JSON string
	_, err = ConvertSecretToBase64("invalid json")
	if err == nil {
		t.Errorf("ConvertSecretToBase64 function failed to return an error")
	}
}