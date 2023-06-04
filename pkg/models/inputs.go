package models

import (
	"strings"
)

type Labels []string

type Parameters struct {
	Labels       string `json:"labels"`
	ExcludeLabel string `json:"excludeLabel"`
	BlackHole    string `json:"blackHole"`
	Path         string `json:"path"`
	Number       string `json:"number"`
}

type PluginInputs struct {
	ApplicationSetName string `json:"applicationSetName"`
	Input              struct {
		Parameters Parameters `json:"parameters"`
	} `json:"input"`
}

// convert string contains labels to slice of labels
func (p *Parameters) GetLabels() Labels {
	var labels Labels

	labelStr := p.Labels
	trimmedStr := strings.Trim(labelStr, "[]") // Remove the square brackets from the string
	labels = strings.Split(trimmedStr, " ")

	return labels

}

// check if the PR labels contain the exclude label
func (p *Parameters) Deployable() bool {

	labels := p.GetLabels()

	for _, label := range labels {
		if label == p.ExcludeLabel {
			return false
		}
	}
	return true

}
