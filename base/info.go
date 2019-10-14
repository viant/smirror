package base

//Info represents meta route information
type Info struct {
	Workflow          string `json:",omitempty"`
	Description       string `json:",omitempty"`
	ProjectURL        string `json:",omitempty"`
	LeadEngineer      string `json:",omitempty"`
	LeadEngineerEmail string `json:",omitempty"`
	SlackChannel      string `json:",omitempty"`
	SlackUser         string `json:",omitempty"`
}
