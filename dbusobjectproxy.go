package main

import (
	"strings"
	"time"

	dbus "github.com/godbus/dbus/v5"
)

type (
	DBusObjectProxy struct {
		dbus.BusObject
		conn  *dbus.Conn
		iface string
	}

	Properties map[string]interface{}
)

func newDBusObjectProxy(conn *dbus.Conn, dest, iface, path string) DBusObjectProxy {
	return DBusObjectProxy{
		BusObject: conn.Object(dest, dbus.ObjectPath(path)),
		conn:      conn,
		iface:     iface,
	}
}

func (p DBusObjectProxy) propName(key string) string {
	return strings.Join([]string{p.iface, key}, ".")
}

func (p DBusObjectProxy) GetObjectPathProperty(key string) (dbus.ObjectPath, error) {
	v, err := p.BusObject.GetProperty(p.propName(key))
	if err != nil {
		return "", err
	}

	return v.Value().(dbus.ObjectPath), nil
}

func (p DBusObjectProxy) GetStringProperty(key string) (string, error) {
	v, err := p.BusObject.GetProperty(p.propName(key))
	if err != nil {
		return "", err
	}

	return v.Value().(string), nil
}

func (p DBusObjectProxy) GetStringSliceProperty(key string) ([]string, error) {
	v, err := p.BusObject.GetProperty(p.propName(key))
	if err != nil {
		return nil, err
	}

	return v.Value().([]string), nil
}

func (p DBusObjectProxy) GetBoolProperty(key string) (bool, error) {
	v, err := p.BusObject.GetProperty(p.propName(key))
	if err != nil {
		return false, err
	}

	return v.Value().(bool), nil
}

func (p DBusObjectProxy) GetDurationProperty(key string) (time.Duration, error) {
	v, err := p.BusObject.GetProperty(p.propName(key))
	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(v.Value().(uint32)), nil
}

func (p DBusObjectProxy) GetUint32Property(key string) (uint32, error) {
	v, err := p.BusObject.GetProperty(p.propName(key))
	if err != nil {
		return uint32(0), err
	}

	return v.Value().(uint32), nil
}

func (p DBusObjectProxy) GetUint16Property(key string) (uint16, error) {
	v, err := p.BusObject.GetProperty(p.propName(key))
	if err != nil {
		return uint16(0), err
	}

	return v.Value().(uint16), nil
}

func (p DBusObjectProxy) GetByteSliceProperty(key string) ([]byte, error) {
	v, err := p.BusObject.GetProperty(p.propName(key))
	if err != nil {
		return nil, err
	}

	return v.Value().([]byte), nil
}
