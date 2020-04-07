package main

import "github.com/godbus/dbus"

type (
	GattService struct {
		DBusObjectProxy

		Characteristics map[string]*GattCharacteristic
	}
)

func NewGattService(conn *dbus.Conn, path string) *GattService {
	return &GattService{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, "org.bluez.GattService1", path),
		Characteristics: make(map[string]*GattCharacteristic),
	}
}

func (s *GattService) attachCharacteristic(c *GattCharacteristic) error {
	uuid, err := c.UUID()
	if err != nil {
		return nil
	}

	s.Characteristics[uuid] = c

	return nil
}

func (s *GattService) UUID() (string, error) {
	return s.GetStringProperty("UUID")
}

func (s *GattService) Primary() (bool, error) {
	return s.GetBoolProperty("Primary")
}

func (s *GattService) Device() (string, error) {
	path, err := s.GetObjectPathProperty("Device")
	if err != nil {
		return "", err
	}

	return string(path), nil
}

func (s *GattService) Includes() ([]string, error) {
	return s.GetStringSliceProperty("Includes")
}
