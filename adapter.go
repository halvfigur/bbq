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

func (a *Adapter) SetDiscoveryFilter(filter map[string]interface{}) error {
	debug("Adapter.SetDiscoveryFilter(ctx, %v)", filter)

	_, err := a.call(context.Background(), "org.bluez.Adapter1.SetDiscoveryFilter",
		0, filter)

	return err
}

func (a *Adapter) StartDiscovery(ctx context.Context) ([]interface{}, error) {
	debug("Adapter.StartDiscovery(ctx)")

	return a.call(ctx, "org.bluez.Adapter1.StartDiscovery", 0)
}

func (a *Adapter) Address() (string, error) {
	return a.GetStringProperty(a.iface, "Address")
}

func (a *Adapter) AddressType() (string, error) {
	return a.GetStringProperty(a.iface, "AddressType")
}

func (a *Adapter) Name() (string, error) {
	return a.GetStringProperty(a.iface, "Name")
}

func (a *Adapter) Alias() (string, error) {
	return a.GetStringProperty(a.iface, "Alias")
}

func (a *Adapter) Class() (uint32, error) {
	return a.GetUint32Property(a.iface, "Class")
}

func (a *Adapter) SetPowered(powered bool) error {
	return a.SetProperty(a.iface, "Powered", powered)
}

func (a *Adapter) Powered() (bool, error) {
	return a.GetBoolProperty(a.iface, "Powered")
}

func (a *Adapter) Discoverable() (bool, error) {
	return a.GetBoolProperty(a.iface, "Discoverable")
}

func (a *Adapter) Pairable() (bool, error) {
	return a.GetBoolProperty(a.iface, "Pairable")
}

func (a *Adapter) PairableTimeout() (time.Duration, error) {
	return a.GetDurationProperty(a.iface, "Pairable")
}

func (a *Adapter) DiscoverableTimeout() (time.Duration, error) {
	return a.GetDurationProperty(a.iface, "DiscoverableTimeout")
}

func (a *Adapter) Discovering() (bool, error) {
	return a.GetBoolProperty(a.iface, "Discovering")
}

func (a *Adapter) UUIDS() ([]string, error) {
	v, err := a.GetProperty(a.iface, "UUIDS")
	if err != nil {
		return nil, err
	}

	return v[0].([]string), nil
}
