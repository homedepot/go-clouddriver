package kubernetes

var DefaultStatus = ManifestStatus{
	Available: Available{
		State: true,
	},
	Failed: Failed{
		State: false,
	},
	Paused: Paused{
		State: false,
	},
	Stable: Stable{
		State: true,
	},
}

var NoneReported = ManifestStatus{
	Available: Available{
		State:   false,
		Message: "No availability reported",
	},
	Failed: Failed{
		State: false,
	},
	Paused: Paused{
		State: false,
	},
	Stable: Stable{
		State:   false,
		Message: "No status reported yet",
	},
}
