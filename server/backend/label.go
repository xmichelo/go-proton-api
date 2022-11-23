package backend

import (
	"github.com/ProtonMail/go-proton-api"
	"github.com/google/uuid"
)

type label struct {
	labelID    string
	parentID   string
	name       string
	labelType  proton.LabelType
	messageIDs map[string]struct{}
}

func newLabel(labelName, parentID string, labelType proton.LabelType) *label {
	return &label{
		labelID:    uuid.NewString(),
		parentID:   parentID,
		name:       labelName,
		labelType:  labelType,
		messageIDs: make(map[string]struct{}),
	}
}

func (label *label) toLabel(labels map[string]*label) proton.Label {
	var path []string

	for labelID := label.labelID; labelID != ""; labelID = labels[labelID].parentID {
		path = append([]string{labels[labelID].name}, path...)
	}

	return proton.Label{
		ID:   label.labelID,
		Name: label.name,
		Path: path,
		Type: label.labelType,
	}
}
