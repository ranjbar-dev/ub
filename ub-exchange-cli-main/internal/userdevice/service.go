package userdevice

import "strings"

const (
	DeviceWeb     = "WEB"
	DeviceAndroid = "ANDROID"
	DeviceIos     = "IOS"
)

func GetDeviceUsingUserAgent(userAgent string) string {
	device := DeviceWeb
	userAgent = strings.ToLower(userAgent)
	if strings.HasPrefix(userAgent, "ub") && strings.Contains(userAgent, strings.ToLower(DeviceAndroid)) {
		return DeviceAndroid
	}

	if strings.HasPrefix(userAgent, "ub") && strings.Contains(userAgent, strings.ToLower(DeviceIos)) {
		return DeviceIos
	}

	return device

}
