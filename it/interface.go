package it

// Beacon interface for IT service
type BeaconSender interface {
	Send(path string)
}
