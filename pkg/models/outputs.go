package models

type PluginOutputs struct {
	Output struct {
		Parameters []OutputsParameters `json:"parameters"`
	} `json:"output"`
}

type OutputsParameters struct {
	GeneratedPath string `json:"generatedPath"`
}

func (p *PluginOutputs) SetOutputParameters(op []OutputsParameters) {
	p.Output.Parameters = op
}
