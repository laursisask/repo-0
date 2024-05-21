package secret_service

import (
	"fmt"

	dbus "github.com/godbus/dbus/v5"
)

const (
	itemInterface    = "org.freedesktop.Secret.Item"
	promptInterface  = "org.freedesktop.Secret.Prompt"
	serviceInterface = "org.freedesktop.Secret.Service"
	serviceName      = "org.freedesktop.secrets"
)

// Secret to be stored in the secret service.
type Secret struct {
	Session     dbus.ObjectPath
	Parameters  []byte
	Value       []byte
	ContentType string `dbus:"content_type"`
}

// DBusManager is a wrapper around the dbus connection and bus object.
type DBusManager struct {
	*dbus.Conn
	busObject dbus.BusObject
}

// Creates a new DBusManager for the secret service.
func NewDBusManager() (*DBusManager, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	return &DBusManager{
		conn,
		conn.Object(serviceName, "/org/freedesktop/secrets"),
	}, nil
}

// Closes the session with the secret service.
func (m *DBusManager) CloseSession(session dbus.BusObject) error {
	return session.Call("org.freedesktop.Secret.Session.Close", 0).Err
}

// Creates and saves the item/secret to be stored in the collection.
func (m *DBusManager) CreateItem(collection dbus.BusObject, label string, attributes map[string]string, secret Secret) error {
	properties := map[string]dbus.Variant{
		itemInterface + ".Label":      dbus.MakeVariant(label),
		itemInterface + ".Attributes": dbus.MakeVariant(attributes),
	}

	var item, prompt dbus.ObjectPath
	err := collection.Call("org.freedesktop.Secret.Collection.CreateItem", 0,
		properties, secret, true).Store(&item, &prompt)
	if err != nil {
		return err
	}

	return nil
}

// Deletes the item/secret from the collection.
func (s *DBusManager) Delete(itemPath dbus.ObjectPath) error {
	var prompt dbus.ObjectPath
	err := s.Object(serviceName, itemPath).Call(itemInterface+".Delete", 0).Store(&prompt)
	if err != nil {
		return err
	}

	_, _, err = s.handlePrompt(prompt)
	if err != nil {
		return err
	}

	return nil
}

// Handles the prompt signal from the secret service
func (s *DBusManager) handlePrompt(prompt dbus.ObjectPath) (bool, dbus.Variant, error) {
	if prompt != dbus.ObjectPath("/") {
		err := s.AddMatchSignal(dbus.WithMatchObjectPath(prompt),
			dbus.WithMatchInterface(promptInterface),
		)
		if err != nil {
			return false, dbus.MakeVariant(""), err
		}

		defer func(s *DBusManager, options ...dbus.MatchOption) {
			_ = s.RemoveMatchSignal(options...)
		}(s, dbus.WithMatchObjectPath(prompt), dbus.WithMatchInterface(promptInterface))

		promptSignal := make(chan *dbus.Signal, 1)
		s.Signal(promptSignal)

		err = s.Object(serviceName, prompt).Call(promptInterface+".Prompt", 0, "").Err
		if err != nil {
			return false, dbus.MakeVariant(""), err
		}

		signal := <-promptSignal
		switch signal.Name {
		case promptInterface + ".Completed":
			dismissed := signal.Body[0].(bool)
			result := signal.Body[1].(dbus.Variant)
			return dismissed, result, nil
		}
	}

	return false, dbus.MakeVariant(""), nil
}

// Opens a new session with the secret service.
func (m *DBusManager) OpenSession() (dbus.BusObject, error) {
	var disregard dbus.Variant
	var sessionPath dbus.ObjectPath
	err := m.busObject.Call(serviceInterface+".OpenSession", 0, "plain", dbus.MakeVariant("")).Store(&disregard, &sessionPath)
	if err != nil {
		return nil, err
	}

	return m.Object(serviceName, sessionPath), nil
}

// Ensures the collection is in an unlocked state before saving or fetching secrets.
func (m *DBusManager) Unlock(collection dbus.ObjectPath) error {
	var unlocked []dbus.ObjectPath
	var prompt dbus.ObjectPath
	err := m.busObject.Call(serviceInterface+".Unlock", 0, []dbus.ObjectPath{collection}).Store(&unlocked, &prompt)
	if err != nil {
		return err
	}

	_, v, err := m.handlePrompt(prompt)
	if err != nil {
		return err
	}

	collections := v.Value()
	switch c := collections.(type) {
	case []dbus.ObjectPath:
		unlocked = append(unlocked, c...)
	}

	if len(unlocked) != 1 || (collection != "/org/freedesktop/secrets/aliases/default" && unlocked[0] != collection) {
		return fmt.Errorf("failed to unlock correct collection '%v'", collection)
	}

	return nil
}

// Returns the default collection for the secret service.
func GetDefaultCollection(m *DBusManager) dbus.BusObject {
	return m.Object(serviceName, "/org/freedesktop/secrets/collection/login")
}

// Create a new secret object.
func NewSecret(session dbus.ObjectPath, secret string) Secret {
	return Secret{
		Session:     session,
		Parameters:  []byte{},
		Value:       []byte(secret),
		ContentType: "text/plain; charset=utf8",
	}
}
