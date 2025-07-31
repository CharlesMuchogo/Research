package dto

type GoogleLoginRequest struct {
	Token    string `json:"token"`
	DeviceId string `json:"device_id"`
	Country  string `json:"country"`
}
