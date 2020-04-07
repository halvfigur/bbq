package main

import "github.com/godbus/dbus"

type (
	GattCharacteristic struct {
		DBusObjectProxy

		Descriptors map[string]*GattDescriptor
	}
)

func NewGattCharacteristic(conn *dbus.Conn, path string) *GattCharacteristic {
	debug("NewGattCharacteristic(%v, %v)", conn, path)

	return &GattCharacteristic{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, "org.bluez.GattCharacteristic1", path),
		Descriptors:     make(map[string]*GattDescriptor),
	}
}

func (c *GattCharacteristic) attachDescriptor(d *GattDescriptor) error {
	uuid, err := d.UUID()
	if err != nil {
		return nil
	}

	c.Descriptors[uuid] = d

	return nil
}

func (c *GattCharacteristic) ReadValue() ([]byte, error) {
	debug("GattCharacteristic.ReadValue()")

	blob := make([]byte, 0, 1024)
	if err := c.Call("org.bluez.GattCharacteristic1.ReadValue", 0).Store(&blob); err != nil {
		return nil, err
	}

	return blob, nil
}

func (c *GattCharacteristic) WriteValue(data []byte) error {
	debug("GattCharacteristic.WriteValue()")

	return c.Call("org.bluez.GattCharacteristic1.WriteValue", 0, data).Store()
}

func (c *GattCharacteristic) StartNofity() error {
	debug("GattCharacteristic.StartNofity()")

	return c.Call("org.bluez.GattCharacteristic1.StartNofity", 0).Store()
}

func (c *GattCharacteristic) StopNofity() error {
	debug("GattCharacteristic.StopNofity()")

	return c.Call("org.bluez.GattCharacteristic1.StopNofity", 0).Store()
}

func (c *GattCharacteristic) UUID() (string, error) {
	return c.GetStringProperty("UUID")
}

func (c *GattCharacteristic) Service() (string, error) {
	path, err := c.GetObjectPathProperty("Service")
	if err != nil {
		return "", err
	}

	return string(path), nil
}

func (c *GattCharacteristic) Notifying() (bool, error) {
	return c.GetBoolProperty("Notifying")
}

func (c *GattCharacteristic) Flags() ([]string, error) {
	return c.GetStringSliceProperty("Flags")
}
