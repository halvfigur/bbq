package main

import (
	"context"
	"errors"
	"time"

	dbus "github.com/godbus/dbus/v5"
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
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, "org.bluez.Adapter1", path),
	}
}

func (a *Adapter) SetDiscoveryFilter(filter map[string]interface{}) error {
	debug("Adapter.SetDiscoveryFilter(ctx, %v)", filter)
	return a.Call("org.bluez.Adapter1.SetDiscoveryFilter", 0, filter).Store()
}

func (a *Adapter) StartDiscovery(ctx context.Context) ([]interface{}, error) {
	debug("Adapter.StartDiscovery(ctx)")

	v := make([]interface{}, 0)
	if err := a.CallWithContext(ctx, "org.bluez.Adapter1.StartDiscovery", 0).Store(v); err != nil {
		return nil, err
	}

	return v, nil
}

func (a *Adapter) Address() (string, error) {
	return a.GetStringProperty("Address")
}

func (a *Adapter) AddressType() (string, error) {
	return a.GetStringProperty("AddressType")
}

func (a *Adapter) Name() (string, error) {
	return a.GetStringProperty("Name")
}

func (a *Adapter) Alias() (string, error) {
	return a.GetStringProperty("Alias")
}

func (a *Adapter) Class() (uint32, error) {
	return a.GetUint32Property("Class")
}

func (a *Adapter) SetPowered(powered bool) error {
	return a.SetProperty("Powered", powered)
}

func (a *Adapter) Powered() (bool, error) {
	return a.GetBoolProperty("Powered")
}

func (a *Adapter) Discoverable() (bool, error) {
	return a.GetBoolProperty("Discoverable")
}

func (a *Adapter) Pairable() (bool, error) {
	return a.GetBoolProperty("Pairable")
}

func (a *Adapter) PairableTimeout() (time.Duration, error) {
	return a.GetDurationProperty("Pairable")
}

func (a *Adapter) DiscoverableTimeout() (time.Duration, error) {
	return a.GetDurationProperty("DiscoverableTimeout")
}

func (a *Adapter) Discovering() (bool, error) {
	return a.GetBoolProperty("Discovering")
}

func (a *Adapter) UUIDS() ([]string, error) {
	return a.GetStringSliceProperty("UUIDS")
}
