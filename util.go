package main

import "context"

func findDevice(ctx context.Context, m *ObjectManager, name string) []*Device {
	paths := m.GetManagedObjects(ctx)

	devices := make([]*Device, 0)

	for path, ifaces := range paths {
		for iface, attrs := range ifaces {
			if iface != "org.bluez.Device1" {
				continue
			}

			if attrs["Name"] == name {
				devices = append(devices, NewDevice(m.conn, path))
			}
		}
	}

	return devices
}
