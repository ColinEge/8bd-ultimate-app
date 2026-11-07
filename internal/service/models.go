package service

type Device struct {
	DeviceID    string
	ProductName string
}

type Keyboard struct {
	Model       uint16
	DeviceID    string
	ProductName string
	HasKeypad   bool
	Keys        map[string]Key
}

type Profile struct {
	Title string           `json:"title"`
	Keys  []ProfileBinding `json:"keys"`
}

type ProfileBinding struct {
	Key    uint8  `json:"key"`
	Action string `json:"action"`
}

type Key struct {
	ID           string `json:"id"`
	PhysicalKey  uint8  `json:"physicalKey"`
	KeyLabel     string `json:"keyLabel"`
	BindingState string `json:"bindingState"`
	ActionCode   string `json:"actionCode"`
	ActionLabel  string `json:"actionLabel"`
}

type Action struct {
	ActionCode  string `json:"actionCode"`
	ActionLabel string `json:"actionLabel"`
}

type ActionCategory struct {
	Name    string
	Actions []Action
}
