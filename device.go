package main

import (
	"context"
	"errors"

	dbus "github.com/godbus/dbus/v5"
)

type (
	Device struct {
		DBusObjectProxy

		Services map[string]*GattService
	}
)

var ErrServiceNotFound = errors.New("service not found")

func NewDevice(conn *dbus.Conn, path string) *Device {
	debug("NewDevice(%v, %v)", conn, path)

	return &Device{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, "org.bluez.Device1", path),
		Services:        make(map[string]*GattService),
	}
}

func (d *Device) attachService(s *GattService) error {
	uuid, err := s.UUID()
	if err != nil {
		return nil
	}

	d.Services[uuid] = s

	return nil
}

func (d *Device) Disconnect(ctx context.Context) error {
	debug("Device.Disconnect()")

	return d.CallWithContext(ctx, "org.bluez.Device1.Disconnect", 0).Store()
}

func (d *Device) Connect(ctx context.Context) error {
	debug("Device.Connect()")

	return d.CallWithContext(ctx, "org.bluez.Device1.Connect", 0).Store()
}

func (d *Device) DisconnectProfile(ctx context.Context, uuid string) error {
	debug("Device.DisconnectProfile()")

	return d.CallWithContext(ctx, "org.bluez.Device1.DisconnectProfile", 0).Store()
}

func (d *Device) ConnectProfile(ctx context.Context, uuid string) error {
	debug("Device.ConnectProfile()")

	return d.CallWithContext(ctx, "org.bluez.Device1.ConnectProfile", 0).Store()
}

func (d *Device) Pair(ctx context.Context) error {
	debug("Device.Pair()")

	return d.CallWithContext(ctx, "org.bluez.Device1.Pair", 0).Store()
}

func (d *Device) CancelPairing(ctx context.Context) error {
	debug("Device.CancelPairing()")

	return d.CallWithContext(ctx, "org.bluez.Device1.CancelPairing", 0).Store()
}

func (d *Device) Address() (string, error) {
	return d.GetStringProperty("Address")
}

func (d *Device) AddressType() (string, error) {
	return d.GetStringProperty("AddressType")
}

func (d *Device) Service(uuid string) (*GattService, error) {
	s, ok := d.Services[uuid]
	if !ok {
		return nil, ErrServiceNotFound
	}

	return s, nil
}

func (d *Device) Characteristic(uuid string) (*GattCharacteristic, error) {
	for _, s := range d.Services {
		if c, err := s.Characteristic(uuid); err == nil {
			return c, nil
		}
	}

	return nil, ErrCharacterisicNotFound
}

func (d *Device) Descriptor(uuid string) (*GattDescriptor, error) {
	for _, s := range d.Services {
		if d, err := s.Descriptor(uuid); err == nil {
			return d, nil
		}
	}

	return nil, ErrDescriptorNotFound
}

func (d *Device) Name() (string, error) {
	return d.GetStringProperty("Name")
}

func (d *Device) Icon() (string, error) {
	return d.GetStringProperty("Icon")
}

func (d *Device) Class() (uint32, error) {
	return d.GetUint32Property("Icon")
}

func (d *Device) Appearance() (uint32, error) {
	return d.GetUint32Property("Appearance")
}

func (d *Device) UUIDS() ([]string, error) {
	v, err := d.GetProperty("Appearance")
	if err != nil {
		return nil, err
	}

	return v.Value().([]string), nil
}

func (d *Device) Paried() (bool, error) {
	return d.GetBoolProperty("Paried")
}

func (d *Device) Connected() (bool, error) {
	return d.GetBoolProperty("Connected")
}

func (d *Device) Trusted() (bool, error) {
	return d.GetBoolProperty("Trusted")
}

func (d *Device) Blocked() (bool, error) {
	return d.GetBoolProperty("Blocked")
}

func (d *Device) Alias() (string, error) {
	return d.GetStringProperty("Alias")
}

func (d *Device) Adapter() (*Adapter, error) {
	v, err := d.GetProperty("Adapter")
	if err != nil {
		return nil, err
	}

	return NewAdapter(d.conn, v.Value().(string)), nil
}

func (d *Device) LegacyPairing() (bool, error) {
	return d.GetBoolProperty("LegacyPairing")
}

func (d *Device) Modalias() (string, error) {
	return d.GetStringProperty("Modalias")
}

func (d *Device) RSSI() (uint16, error) {
	return d.GetUint16Property("RSSI")
}

func (d *Device) TxPower() (uint16, error) {
	return d.GetUint16Property("TxPower")
}

func (d *Device) ManufacturerData() (map[uint16][]byte, error) {
	v, err := d.GetProperty("ManufacturerData")
	if err != nil {
		return nil, err
	}

	data := make(map[uint16][]byte)

	for k, v := range v.Value().(map[uint16]dbus.Variant) {
		data[k] = v.Value().([]byte)
	}

	return data, nil
}

func (d *Device) ServiceData() (map[string][]byte, error) {
	v, err := d.GetProperty("ServiceData")
	if err != nil {
		return nil, err
	}

	data := make(map[string][]byte)

	for k, v := range v.Value().(map[string]dbus.Variant) {
		data[k] = v.Value().([]byte)
	}

	return data, nil
}

func (d *Device) ServicesResolved() (bool, error) {
	return d.GetBoolProperty("ServicesResolved")
}

func (d *Device) AdvertisingFlags() ([]byte, error) {
	return d.GetByteSliceProperty("AdvertisingFlags")
}
