package main

import (
	"context"
	"time"

	"github.com/godbus/dbus"
)

type (
	DBusObjectProxy struct {
		conn *dbus.Conn
		dobj dbus.BusObject
	}

	Properties map[string]interface{}
)

func newDBusObjectProxy(conn *dbus.Conn, dest, path string) DBusObjectProxy {
	return DBusObjectProxy{conn, conn.Object(dest, dbus.ObjectPath(path))}
}

func (p DBusObjectProxy) call(ctx context.Context, method string, flags dbus.Flags, args ...interface{}) ([]interface{}, error) {
	debug("DBusProxyObject.call(ctx, %v, %v, %v)", method, flags, args)

	ch := make(chan *dbus.Call, 1)

	p.dobj.Go(method, flags, ch, args...)
	select {
	case <-ctx.Done():
		return nil, ErrContext
	case c := <-ch:
		return c.Body, c.Err
	}
}

func (p DBusObjectProxy) SetProperty(iface, key string, value interface{}) error {
	_, err := p.call(context.Background(), "org.freedesktop.DBus.Properties.Set",
		0, iface, key, dbus.MakeVariant(value))

	return err
}

func (p DBusObjectProxy) GetProperty(iface, key string) ([]interface{}, error) {
	v, err := p.call(context.Background(), "org.freedesktop.DBus.Properties.Get", 0, iface, key)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (p DBusObjectProxy) GetStringProperty(iface, key string) (string, error) {
	v, err := p.call(context.Background(), "org.freedesktop.DBus.Properties.Get", 0, iface, key)
	if err != nil {
		return "", err
	}

	return v[0].(string), nil
}

func (p DBusObjectProxy) GetBoolProperty(iface, key string) (bool, error) {
	v, err := p.call(context.Background(), "org.freedesktop.DBus.Properties.Get", 0, iface, key)
	if err != nil {
		return false, err
	}

	return v[0].(bool), nil
}

func (p DBusObjectProxy) GetDurationProperty(iface, key string) (time.Duration, error) {
	v, err := p.call(context.Background(), "org.freedesktop.DBus.Properties.Get", 0, iface, key)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(v[0].(uint32)), nil
}

func (p DBusObjectProxy) GetUint32Property(iface, key string) (uint32, error) {
	v, err := p.call(context.Background(), "org.freedesktop.DBus.Properties.Get", 0, iface, key)
	if err != nil {
		return uint32(0), err
	}

	return v[0].(uint32), nil
}

func (p DBusObjectProxy) GetUint16Property(iface, key string) (uint16, error) {
	v, err := p.call(context.Background(), "org.freedesktop.DBus.Properties.Get", 0, iface, key)
	if err != nil {
		return uint16(0), err
	}

	return v[0].(uint16), nil
}
