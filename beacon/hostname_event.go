package beacon

// HostnameEvent contains the hostname data
type HostnameEvent struct {
	Hostname  string `json:"hostname"`
	UpdatedAt string `json:"updated_at"`
}

// NewHostnameEvent creates HostnameEvent
func NewHostnameEvent(
	hostname, updatedAt string,
) HostnameEvent {
	return HostnameEvent{
		Hostname:  hostname,
		UpdatedAt: updatedAt,
	}
}
