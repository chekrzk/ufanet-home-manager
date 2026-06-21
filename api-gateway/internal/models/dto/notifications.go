package dto

type RegisterDeviceRequest struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

func (r RegisterDeviceRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Token == "" {
		errs["token"] = "required"
	}
	if r.Platform == "" {
		errs["platform"] = "required"
	}
	return errs
}

type UnregisterDeviceRequest struct {
	Token string `json:"token"`
}

func (r UnregisterDeviceRequest) Validate() map[string]string {
	if r.Token == "" {
		return map[string]string{"token": "required"}
	}
	return nil
}
