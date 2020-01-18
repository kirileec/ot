package ot

import (
	"errors"
)

var (
	ErrSelecttionUnmarshalFailed = errors.New("ot/selection: unmarshal failed")
)

// the selection range
type Selection struct {
	Ranges []Range `json:"ranges"`
}

func (s *Selection) Transform(op *Operation) *Selection {
	tr := make([]Range, len(s.Ranges))
	for i, r := range s.Ranges {
		tr[i] = *r.Transform(op)
	}
	return &Selection{tr}
}

func (s *Selection) SelectionMarshal() map[string]interface{} {
	mr := make([]map[string]interface{}, len(s.Ranges))
	for i, r := range s.Ranges {
		mr[i] = map[string]interface{}{"anchor": r.Anchor, "head": r.Head}
	}
	return map[string]interface{}{"ranges": mr}
}

func SelectionUnmarshal(data map[string]interface{}) (*Selection, error) {
	if data["ranges"] == nil {
		return nil, ErrSelecttionUnmarshalFailed
	}

	dr, ok := data["ranges"].([]interface{})
	if !ok {
		return nil, ErrSelecttionUnmarshalFailed
	}

	ranges := make([]Range, len(dr))

	for i, o := range dr {
		r, ok := o.(map[string]interface{})
		if !ok {
			return nil, ErrSelecttionUnmarshalFailed
		}
		rng, err := selectionUnmarshalRange(r)
		if err != nil {
			return nil, err
		}
		ranges[i] = *rng
	}

	return &Selection{ranges}, nil
}

func selectionUnmarshalRange(data map[string]interface{}) (*Range, error) {
	a, ok := selectionParseNumber(data["anchor"])
	if !ok {
		return nil, ErrSelecttionUnmarshalFailed
	}
	h, ok := selectionParseNumber(data["head"])
	if !ok {
		return nil, ErrSelecttionUnmarshalFailed
	}
	return &Range{a, h}, nil
}

func selectionParseNumber(n interface{}) (int, bool) {
	switch n.(type) {
	case int:
		return n.(int), true
	case float64:
		return int(n.(float64)), true
	}
	return 0, false
}

type Range struct {
	Anchor int `json:"anchor"`
	Head   int `json:"head"`
}

func (r *Range) Transform(op *Operation) *Range {
	return &Range{transformIndex(r.Anchor, op), transformIndex(r.Head, op)}
}

func transformIndex(i int, op *Operation) int {
	// start cursor at index 0
	j := 0

	for _, op := range op.Ops {
		// if cursor index is greater than i, the rest of the ops are irrelevant
		if j > i {
			break
		}
		if IsRetain(op) {
			// advance cursor
			j += op.N
		} else if IsInsert(op) {
			// insertion increments index. also advance cursor
			i += len(op.S)
			j += len(op.S)
		} else if IsDelete(op) {
			// deletion decrements index, but only up to current cursor
			i = max(j, i+op.N) // N is negative
		}
	}

	return i
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
