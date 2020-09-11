package monitors

type DigitalMonitorMap map[string]DigitalMonitor

func (monitors DigitalMonitorMap) Subscribe(observer *DigitalEventObserver) {
	for _, monitor := range monitors {
		monitor.Subscribe(observer)
	}
}

func (monitors DigitalMonitorMap) Unsubscribe(observer *DigitalEventObserver) {
	for _, monitor := range monitors {
		monitor.Unsubscribe(observer)
	}
}
