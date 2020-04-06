package main

import "github.com/godbus/dbus"

type (
	GattService struct {
		DBusObjectProxy
		iface string
	}
)

func NewGattService(conn *dbus.Conn, path string) *GattService {
	return &GattService{
		DBusObjectProxy: newDBusObjectProxy(conn, destOrgBluez, path),
		iface:           "org.bluez.GattService1",
	}
}

func (s *GattService) UUID() (string, error) {
	return s.GetStringProperty(s.iface, "UUID")
}

/*
   readonly s UUID = '00001801-0000-1000-8000-00805f9b34fb';
     readonly o Device = '/org/bluez/hci0/dev_E9_13_7F_70_2C_51';
     readonly b Primary = true;
     readonly ao Includes = []
*/
