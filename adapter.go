package main

import (
	"context"
	"errors"
	"time"

	"github.com/godbus/dbus"
)

type (
	Adapter struct {
		DBusObjectProxy
		Properties
		iface string
	}

	Device struct {
		DBusObjectProxy
		Properties
		iface string
	}

	ObjectManager struct {
		DBusObjectProxy
		iface string
	}
)

const (
	destOrgBluez = "org.bluez"
)

var ErrContext = errors.New("context done")

func NewAdapter(conn *dbus.Conn, path string) *Adapter {
	debug("NewAdapter(%v, %v)", conn, path)

	return &Adapter{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, path),
		iface:           "org.bluez.Adapter1",
	}
}

func (a *Adapter) SetProperty(ctx context.Context, key string, value interface{}) error {
	debug("Adapter.SetProperty(%v, %v)", key, value)

	return a.DBusObjectProxy.SetProperty(ctx, a.iface, key, value)
}

func (a *Adapter) SetDiscoveryFilter(ctx context.Context, filter map[string]interface{}) error {
	debug("Adapter.SetDiscoveryFilter(ctx, %v)", filter)

	_, err := a.call(ctx, "org.bluez.Adapter1.SetDiscoveryFilter",
		0, filter)

	return err
}

func (a *Adapter) StartDiscovery(ctx context.Context) ([]interface{}, error) {
	debug("Adapter.StartDiscovery(ctx)")

	return a.call(ctx, "org.bluez.Adapter1.StartDiscovery", 0)
}

func (a *Adapter) Address(ctx context.Context) (string, error) {
	return a.GetStringProperty(ctx, a.iface, "Address")
}

func (a *Adapter) AddressType(ctx context.Context) (string, error) {
	return a.GetStringProperty(ctx, a.iface, "AddressType")
}

func (a *Adapter) Name(ctx context.Context) (string, error) {
	return a.GetStringProperty(ctx, a.iface, "Name")
}

func (a *Adapter) Alias(ctx context.Context) (string, error) {
	return a.GetStringProperty(ctx, a.iface, "Alias")
}

func (a *Adapter) Class(ctx context.Context) (uint32, error) {
	return a.GetUint32Property(ctx, a.iface, "Class")
}

func (a *Adapter) Powered(ctx context.Context) (bool, error) {
	return a.GetBoolProperty(ctx, a.iface, "Powered")
}

func (a *Adapter) Discoverable(ctx context.Context) (bool, error) {
	return a.GetBoolProperty(ctx, a.iface, "Discoverable")
}

func (a *Adapter) Pairable(ctx context.Context) (bool, error) {
	return a.GetBoolProperty(ctx, a.iface, "Pairable")
}

func (a *Adapter) PairableTimeout(ctx context.Context) (time.Duration, error) {
	return a.GetDurationProperty(ctx, a.iface, "Pairable")
}

func (a *Adapter) DiscoverableTimeout(ctx context.Context) (time.Duration, error) {
	return a.GetDurationProperty(ctx, a.iface, "DiscoverableTimeout")
}

func (a *Adapter) Discovering(ctx context.Context) (bool, error) {
	return a.GetBoolProperty(ctx, a.iface, "Discovering")
}

func (a *Adapter) UUIDS(ctx context.Context) ([]string, error) {
	v, err := a.GetProperty(ctx, a.iface, "UUIDS")
	if err != nil {
		return nil, err
	}

	return v[0].([]string), nil
}
