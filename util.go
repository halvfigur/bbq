package main

import "fmt"

func findDevices(m *ObjectManager, name string) ([]*Device, error) {
	paths, err := m.GetManagedObjects()
	if err != nil {
		return nil, err
	}

	deviceMap := make(map[string]*Device)
	serviceMap := make(map[string]*GattService)
	characteristicsMap := make(map[string]*GattCharacteristic)
	descriptorsMap := make(map[string]*GattDescriptor)

	for path, ifaces := range paths {
		for iface, attrs := range ifaces {
			switch iface {
			case "org.bluez.Device1":
				if attrs["Name"] == name {
					deviceMap[path] = NewDevice(m.conn, path)
				}
			case "org.bluez.GattService1":
				serviceMap[path] = NewGattService(m.conn, path)
			case "org.bluez.GattCharacteristic1":
				characteristicsMap[path] = NewGattCharacteristic(m.conn, path)
			case "org.bluez.GattDescriptor1":
				descriptorsMap[path] = NewGattDescriptor(m.conn, path)
			}
		}
	}

	for uuid, d := range descriptorsMap {
		charUUID, err := d.Characteristic()
		if err != nil {
			return nil, err
		}

		c, ok := characteristicsMap[charUUID]
		if !ok {
			return nil, fmt.Errorf("descriptor %s depends on unknown characteristic %s", uuid, charUUID)
		}

		c.attachDescriptor(d)
	}

	for uuid, c := range characteristicsMap {
		servUUID, err := c.Service()
		if err != nil {
			return nil, err
		}

		s, ok := serviceMap[servUUID]
		if !ok {
			return nil, fmt.Errorf("characteristic %s depends on unknown service %s", uuid, servUUID)
		}

		s.attachCharacteristic(c)
	}

	for uuid, s := range serviceMap {
		devUUID, err := s.Device()
		if err != nil {
			return nil, err
		}

		d, ok := deviceMap[devUUID]
		if !ok {
			return nil, fmt.Errorf("service %s depends on unknown device %s", uuid, devUUID)
		}

		d.attachService(s)
	}

	devices := make([]*Device, 0, len(deviceMap))
	for _, d := range deviceMap {
		devices = append(devices, d)
	}

	return devices, nil
}
