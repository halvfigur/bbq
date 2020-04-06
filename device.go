package main

import (
	"context"

	"github.com/godbus/dbus"
)

type (
	Device struct {
		DBusObjectProxy
		Properties
		iface string
	}
)

func NewDevice(conn *dbus.Conn, path string) *Device {
	debug("NewDevice(%v, %v)", conn, path)

	return &Device{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, path),
		iface:           "org.bluez.Device1",
	}
}

func (d *Device) Disconnect(ctx context.Context) error {
	debug("Device.Disconnect()")

	_, err := d.call(ctx, "org.bluez.Device1.Disconnect", 0)

	return err
}

func (d *Device) Connect(ctx context.Context) error {
	debug("Device.Connect()")

	_, err := d.call(ctx, "org.bluez.Device1.Connect", 0)

	return err
}

func (d *Device) DisconnectProfile(ctx context.Context, uuid string) error {
	debug("Device.DisconnectProfile()")

	_, err := d.call(ctx, "org.bluez.Device1.DisconnectProfile", 0, uuid)

	return err
}

func (d *Device) ConnectProfile(ctx context.Context, uuid string) error {
	debug("Device.ConnectProfile()")

	_, err := d.call(ctx, "org.bluez.Device1.ConnectProfile", 0, uuid)

	return err
}

func (d *Device) Pair(ctx context.Context) error {
	debug("Device.Pair()")

	_, err := d.call(ctx, "org.bluez.Device1.Pair", 0)

	return err
}

func (d *Device) CancelPairing(ctx context.Context) error {
	debug("Device.CancelPairing()")

	_, err := d.call(ctx, "org.bluez.Device1.CancelPairing", 0)

	return err
}

func (d *Device) Address() (string, error) {
	return d.GetStringProperty(d.iface, "Address")
}

func (d *Device) AddressType() (string, error) {
	return d.GetStringProperty(d.iface, "AddressType")
}

func (d *Device) Name() (string, error) {
	return d.GetStringProperty(d.iface, "Name")
}

func (d *Device) Icon() (string, error) {
	return d.GetStringProperty(d.iface, "Icon")
}

func (d *Device) Class() (uint32, error) {
	return d.GetUint32Property(d.iface, "Icon")
}

func (d *Device) Appearance() (uint32, error) {
	return d.GetUint32Property(d.iface, "Appearance")
}

func (d *Device) UUIDS() ([]string, error) {
	v, err := d.GetProperty(d.iface, "Appearance")
	if err != nil {
		return nil, err
	}

	return v[0].([]string), nil
}

func (d *Device) Paried() (bool, error) {
	return d.GetBoolProperty(d.iface, "Paried")
}

func (d *Device) Connected() (bool, error) {
	return d.GetBoolProperty(d.iface, "Connected")
}

func (d *Device) Trusted() (bool, error) {
	return d.GetBoolProperty(d.iface, "Trusted")
}

func (d *Device) Blocked() (bool, error) {
	return d.GetBoolProperty(d.iface, "Blocked")
}

func (d *Device) Alias() (string, error) {
	return d.GetStringProperty(d.iface, "Alias")
}

func (d *Device) Adapter() (*Adapter, error) {
	v, err := d.GetProperty(d.iface, "Adapter")
	if err != nil {
		return nil, err
	}

	return NewAdapter(d.conn, v[0].(string)), nil
}

func (d *Device) LegacyPairing() (bool, error) {
	return d.GetBoolProperty(d.iface, "LegacyPairing")
}

func (d *Device) Modalias() (string, error) {
	return d.GetStringProperty(d.iface, "Modalias")
}

func (d *Device) RSSI() (uint16, error) {
	return d.GetUint16Property(d.iface, "RSSI")
}

func (d *Device) TxPower() (uint16, error) {
	return d.GetUint16Property(d.iface, "TxPower")
}

func (d *Device) ManufacturerData() (map[uint16][]byte, error) {
	v, err := d.GetProperty(d.iface, "ManufacturerData")
	if err != nil {
		return nil, err
	}

	data := make(map[uint16][]byte)

	for k, v := range v[0].(map[uint16]dbus.Variant) {
		data[k] = v.Value().([]byte)
	}

	return data, nil
}

func (d *Device) ServiceData() (map[string][]byte, error) {
	v, err := d.GetProperty(d.iface, "ServiceData")
	if err != nil {
		return nil, err
	}

	data := make(map[string][]byte)

	for k, v := range v[0].(map[string]dbus.Variant) {
		data[k] = v.Value().([]byte)
	}

	return data, nil
}

func (d *Device) ServicesResolved() (bool, error) {
	return d.GetBoolProperty(d.iface, "ServicesResolved")
}

func (d *Device) AdvertisingFlags() ([]byte, error) {
	v, err := d.GetProperty(d.iface, "AdvertisingFlags")
	if err != nil {
		return nil, err
	}

	return v[0].([]byte), nil
}
