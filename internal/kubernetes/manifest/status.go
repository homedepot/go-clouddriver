package manifest

var DefaultStatus = Status{
	Available: StatusAvailable{
		State: true,
	},
	Failed: StatusFailed{
		State: false,
	},
	Paused: StatusPaused{
		State: false,
	},
	Stable: StatusStable{
		State: true,
	},
}

var NoneReported = Status{
	Available: StatusAvailable{
		State:   false,
		Message: "No availability reported",
	},
	Failed: StatusFailed{
		State: false,
	},
	Paused: StatusPaused{
		State: false,
	},
	Stable: StatusStable{
		State:   false,
		Message: "No status reported yet",
	},
}

type Status struct {
	Available StatusAvailable `json:"available"`
	Failed    StatusFailed    `json:"failed"`
	Paused    StatusPaused    `json:"paused"`
	Stable    StatusStable    `json:"stable"`
}

type StatusAvailable struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

type StatusFailed struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

type StatusPaused struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

type StatusStable struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}
