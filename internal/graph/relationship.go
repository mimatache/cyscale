package graph

import (
	guuid "github.com/google/uuid"
)

type Relationship struct {
	ID    string
	Label string
	From  string
	To    string
}

func newRelationship(from, to *Node, label string) Relationship {
	return Relationship{
		ID:    guuid.New().String(),
		Label: label,
		To:    to.id,
		From:  from.id,
	}
}
