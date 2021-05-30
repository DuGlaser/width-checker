package model

type DeviceOption struct {
	Min      int
	Max      int
	Interval int
}

func (d *DeviceOption) Pattern() int {
	return (d.Max-d.Min)/d.Interval + 1
}
