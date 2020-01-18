package ot

import (
	"errors"
)

var (
	ErrInvalidRevision = errors.New("ot/session: invalid revision")
)

type Session struct {
	Document   string
	Operations []*Operation
	Clients    map[string]*Client
}

func NewSession(document string) *Session {
	return &Session{
		Document:   document,
		Operations: []*Operation{},
		Clients:    map[string]*Client{},
	}
}

func (s *Session) AddClient(id string) {
	s.Clients[id] = &Client{Selection: Selection{[]Range{}}}
}

func (s *Session) RemoveClient(id string) {
	delete(s.Clients, id)
}

func (s *Session) SetName(id, name string) {
	c := s.Clients[id]
	if c != nil {
		c.Name = name
	}
}

func (s *Session) SetSelection(id string, sel *Selection) {
	c := s.Clients[id]
	if c != nil {
		c.Selection = *sel
	}
}

func (s *Session) AddOperation(revision int, op *Operation) (*Operation, error) {
	if revision < 0 || len(s.Operations) < revision {
		return nil, ErrInvalidRevision
	}
	// find concurrent operations client isn't yet aware of
	otherOps := s.Operations[revision:]

	// transform given operation against these operations
	for _, otherOp := range otherOps {
		op1, _, err := OperationTransform(op, otherOp)
		if err != nil {
			return nil, err
		}
		if op.Meta != nil {
			if m, ok := op.Meta.(*Selection); ok {
				op1.Meta = m.Transform(otherOp)
			}
		}

		op = op1
	}

	// apply transformed op on the doc
	doc, err := op.Apply(s.Document)
	if err != nil {
		return nil, err
	}

	s.Document = doc
	s.Operations = append(s.Operations, op)

	return op, nil
}
