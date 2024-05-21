package keyring

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	ss "github.com/Keeper-Security/linux-keyring-utility/pkg/secret_service"
	dbus "github.com/godbus/dbus/v5"
)

// SecretProvider implements the Keyring interface.
type SecretProvider struct{}

// Stores a secret in the keyring. If the secret already exists, it will be overwritten.
func (s *SecretProvider) Set(appName string, secret string) error {
	manager, err := ss.NewDBusManager()
	if err != nil {
		fmt.Println("Error initializing secret service: ", err)
		return err
	}

	// Get the default collection for the user. Usually this is the "login" keyring.
	collection := ss.GetDefaultCollection(manager)

	b64Str, err := ConvertSecretToBase64(secret)
	if err != nil {
		fmt.Println("Error converting secret to base64: ", err)
		return err
	}

	session, err := manager.OpenSession()
	if err != nil {
		fmt.Println("Error opening session: ", err)
		return err
	}
	defer manager.CloseSession(session)

	secretObj := ss.NewSecret(session.Path(), b64Str)

	err = manager.Unlock(collection.Path())
	if err != nil {
		fmt.Println("Error unlocking collection: ", err)
		return err
	}

	err = manager.CreateItem(collection, fmt.Sprintf("Secret for %s", appName), Attributes(appName), secretObj)
	if err != nil {
		fmt.Println("Error creating item: ", err)
		return err
	}

	return nil
}

// Retrieves a secret from the keyring.
func (s *SecretProvider) Get(appName string) (string, error) {
	svc, err := ss.NewDBusManager()
	if err != nil {
		fmt.Println("Error initializing secret service: ", err)
		return "", err
	}

	// Get the default collection for the user. Usually this is the "login" keyring.
	collection := ss.GetDefaultCollection(svc)

	// Ensure the secret exists before opening a session and attempting to fetch it.
	var results []dbus.ObjectPath
	err = collection.Call("org.freedesktop.Secret.Collection.SearchItems", 0, Attributes(appName)).Store(&results)
	if err != nil {
		fmt.Println("Error searching items: ", err)
		return "", err
	}

	if len(results) == 0 {
		fmt.Println("No secret found for app: ", appName)
		return "", err
	}

	session, err := svc.OpenSession()
	if err != nil {
		fmt.Println("Error opening session: ", err)
		return "", err
	}
	defer svc.CloseSession(session)

	err = svc.Unlock(collection.Path())
	if err != nil {
		fmt.Println("Error unlocking collection: ", err)
		return "", err
	}

	var secret ss.Secret
	err = svc.Object("org.freedesktop.secrets", results[0]).Call("org.freedesktop.Secret.Item.GetSecret", 0, session.Path()).Store(&secret)
	if err != nil {
		fmt.Println("Error getting secret: ", err)
		return "", err
	}

	return string(secret.Value), nil
}

// Returns a map of attributes for the secret service.
func Attributes(appName string) map[string]string {
	return map[string]string{
		"Application": appName,
	}
}

// Parses the given KSM config string into a base64 encoded string. Accepts either:
//
// 1) A path to a JSON configuration file.
// 2) A JSON configuration string.
// 3) A Base64 encoded string.
func ConvertSecretToBase64(input string) (string, error) {
	// Check if the secret is already base64 encoded
	_, err := base64.StdEncoding.DecodeString(input)
	if err == nil {
		return input, nil
	}

	// If the secret is a path to a json file, read the file and encode the contents.
	_, err = os.Stat(input)
	if err == nil {
		// Read the file
		file, err := os.ReadFile(input)
		if err != nil {
			fmt.Println("Error reading file: ", err)
			return "", err
		}
		// Encode the file contents
		base64Str := base64.StdEncoding.EncodeToString(file)
		return base64Str, nil
	}

	// Check if the string is valid json, if so, encode the json string, otherwise output an error.
	isValid := json.Valid([]byte(input))
	if isValid {
		base64Str := base64.StdEncoding.EncodeToString([]byte(input))
		return base64Str, nil
	} else {
		fmt.Println("Invalid JSON string")
		return "", fmt.Errorf("invalid JSON string")
	}
}
