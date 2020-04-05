package main

import (
	"context"

	"github.com/godbus/dbus"
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

func (d *Device) Address(ctx context.Context) (string, error) {
	return d.GetStringProperty(ctx, d.iface, "Address")
}

func (d *Device) AddressType(ctx context.Context) (string, error) {
	return d.GetStringProperty(ctx, d.iface, "AddressType")
}

func (d *Device) Name(ctx context.Context) (string, error) {
	return d.GetStringProperty(ctx, d.iface, "Name")
}

func (d *Device) Icon(ctx context.Context) (string, error) {
	return d.GetStringProperty(ctx, d.iface, "Icon")
}

func (d *Device) Class(ctx context.Context) (uint32, error) {
	return d.GetUint32Property(ctx, d.iface, "Icon")
}

func (d *Device) Appearance(ctx context.Context) (uint32, error) {
	return d.GetUint32Property(ctx, d.iface, "Appearance")
}

func (d *Device) UUIDS(ctx context.Context) ([]string, error) {
	v, err := d.GetProperty(ctx, d.iface, "Appearance")
	if err != nil {
		return nil, err
	}

	return v[0].([]string), nil
}

func (d *Device) Paried(ctx context.Context) (bool, error) {
	return d.GetBoolProperty(ctx, d.iface, "Paried")
}

func (d *Device) Connected(ctx context.Context) (bool, error) {
	return d.GetBoolProperty(ctx, d.iface, "Connected")
}

func (d *Device) Trusted(ctx context.Context) (bool, error) {
	return d.GetBoolProperty(ctx, d.iface, "Trusted")
}

func (d *Device) Blocked(ctx context.Context) (bool, error) {
	return d.GetBoolProperty(ctx, d.iface, "Blocked")
}

func (d *Device) Alias(ctx context.Context) (string, error) {
	return d.GetStringProperty(ctx, d.iface, "Alias")
}

func (d *Device) Adapter(ctx context.Context) (*Adapter, error) {
	v, err := d.GetProperty(ctx, d.iface, "Adapter")
	if err != nil {
		return nil, err
	}

	return NewAdapter(d.conn, v[0].(string)), nil
}

func (d *Device) LegacyPairing(ctx context.Context) (bool, error) {
	return d.GetBoolProperty(ctx, d.iface, "LegacyPairing")
}

func (d *Device) Modalias(ctx context.Context) (string, error) {
	return d.GetStringProperty(ctx, d.iface, "Modalias")
}

func (d *Device) RSSI(ctx context.Context) (uint16, error) {
	return d.GetUint16Property(ctx, d.iface, "RSSI")
}

func (d *Device) TxPower(ctx context.Context) (uint16, error) {
	return d.GetUint16Property(ctx, d.iface, "TxPower")
}

func (d *Device) ManufacturerData(ctx context.Context) (map[uint16][]byte, error) {
	v, err := d.GetProperty(ctx, d.iface, "ManufacturerData")
	if err != nil {
		return nil, err
	}

	data := make(map[uint16][]byte)

	for k, v := range v[0].(map[uint16]dbus.Variant) {
		data[k] = v.Value().([]byte)
	}

	return data, nil
}

func (d *Device) ServiceData(ctx context.Context) (map[string][]byte, error) {
	v, err := d.GetProperty(ctx, d.iface, "ServiceData")
	if err != nil {
		return nil, err
	}

	data := make(map[string][]byte)

	for k, v := range v[0].(map[string]dbus.Variant) {
		data[k] = v.Value().([]byte)
	}

	return data, nil
}

func (d *Device) ServicesResolved(ctx context.Context) (bool, error) {
	return d.GetBoolProperty(ctx, d.iface, "ServicesResolved")
}

func (d *Device) AdvertisingFlags(ctx context.Context) ([]byte, error) {
	v, err := d.GetProperty(ctx, d.iface, "AdvertisingFlags")
	if err != nil {
		return nil, err
	}

	return v[0].([]byte), nil
}
