package main

import "github.com/godbus/dbus"

type (
	GattDescriptor struct {
		DBusObjectProxy
	}
)

func NewGattDescriptor(conn *dbus.Conn, path string) *GattDescriptor {
	debug("NewGattDescriptor(%v, %v)", conn, path)

	return &GattDescriptor{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, "org.bluez.GattDescriptor1", path),
	}
}

func (d *GattDescriptor) ReadValue() ([]byte, error) {
	debug("GattDescriptor.ReadValue()")

	blob := make([]byte, 0, 1024)
	if err := d.Call("org.bluez.GattDescriptor1.ReadValue", 0).Store(&blob); err != nil {
		return nil, err
	}

	return blob, nil
}

func (d *GattDescriptor) WriteValue(data []byte) error {
	debug("GattDescriptor.WriteValue()")

	return d.Call("org.bluez.GattDescriptor1.WriteValue", 0, data).Store()
}

func (d *GattDescriptor) UUID() (string, error) {
	return d.GetStringProperty("UUID")
}

func (d *GattDescriptor) Characteristic() (string, error) {
	path, err := d.GetObjectPathProperty("Characteristic")
	if err != nil {
		return "", nil
	}

	return string(path), nil
}

func (d *GattDescriptor) Value() ([]byte, error) {
	return d.GetByteSliceProperty("Value")
}

func (d *GattDescriptor) Flags() ([]string, error) {
	return d.GetStringSliceProperty("Flags")
}
