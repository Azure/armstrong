package jsonpointerpos

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/go-openapi/jsonpointer"
)

type JSONPointerPosition struct {
	Ptr jsonpointer.Pointer
	Position
}

type Position struct {
	Line   int
	Column int
}

func newJSONPtr(tks []string) *jsonpointer.Pointer {
	if len(tks) == 0 {
		return nil
	}
	encTks := make([]string, len(tks))
	for i, tk := range tks {
		encTks[i] = jsonpointer.Escape(tk)
	}
	ptr, _ := jsonpointer.New("/" + strings.Join(encTks, "/"))
	return &ptr
}

type tokenTree struct {
	tk       string
	offset   *int
	children map[string]*tokenTree
}

func (tree *tokenTree) add(ptr jsonpointer.Pointer) {
	tks := ptr.DecodedTokens()
	if len(tks) == 0 || (len(tks) == 1 && tks[0] == "") {
		return
	}
	if tree.children == nil {
		tree.children = map[string]*tokenTree{}
	}
	tk, remains := tks[0], tks[1:]
	subTree, ok := tree.children[tk]
	if !ok {
		subTree = &tokenTree{tk: tk}
		tree.children[tk] = subTree
	}
	remainPtr := newJSONPtr(remains)
	if remainPtr != nil {
		subTree.add(*remainPtr)
	}
}

// flattenOffset flattens the token tree to a map whose key is a json pointer and its value is the offset.
// For token tree nodes that have no offset (implies they doesn't exist in the json document), they are skipped.
func (tree *tokenTree) flattenOffset(parentTks []string) map[string]int {
	out := map[string]int{}

	var tks []string
	for _, tk := range parentTks {
		// This is to skip the root node of the tree when building the pointer
		if tk == "" {
			continue
		}
		tks = append(tks, tk)
	}
	tks = append(tks, tree.tk)

	for _, child := range tree.children {
		m := child.flattenOffset(tks)
		for k, v := range m {
			out[k] = v
		}
	}

	if tree.offset != nil {
		ptr := newJSONPtr(tks)
		out[ptr.String()] = *tree.offset
	}

	return out
}

func buildTokenTree(ptrs []jsonpointer.Pointer) tokenTree {
	root := tokenTree{}
	for _, ptr := range ptrs {
		root.add(ptr)
	}
	return root
}

func GetPositions(document string, ptrs []jsonpointer.Pointer) (map[string]JSONPointerPosition, error) {
	if len(ptrs) == 0 {
		return nil, nil
	}
	tree := buildTokenTree(ptrs)
	dec := json.NewDecoder(strings.NewReader(document))
	dec.UseNumber()

	if _, err := offsetValue(dec, &tree); err != nil {
		return nil, err
	}

	m := tree.flattenOffset(nil)
	nm := map[string]int{}
	// Only keep the specified pointers from the flattened offset map
	for _, ptr := range ptrs {
		if v, ok := m[ptr.String()]; ok {
			nm[ptr.String()] = v
		}
	}
	m = nm

	type offsetItem struct {
		ptr    string
		offset int
	}
	ol := []offsetItem{}
	for ptr, offset := range m {
		ol = append(ol, offsetItem{
			ptr:    ptr,
			offset: offset,
		})
	}
	sort.Slice(ol, func(i, j int) bool {
		return ol[i].offset < ol[j].offset
	})

	var sc scanner.Scanner
	sc.Init(strings.NewReader(document))

	out := map[string]JSONPointerPosition{}

	start := 0
	for _, ov := range ol {
		for i := start; i < ov.offset; i++ {
			sc.Next()
		}
		ptr, err := jsonpointer.New(ov.ptr)
		if err != nil {
			return nil, err
		}
		pos := sc.Pos()
		out[ptr.String()] = JSONPointerPosition{
			Ptr: ptr,
			Position: Position{
				Line:   pos.Line,
				Column: pos.Column,
			},
		}
		start = ov.offset
	}
	return out, nil
}

// offsetValue fill ins the offset(s) of the specified tree for a JSON value.
// Meanwhile, it returns the value length.
func offsetValue(dec *json.Decoder, tree *tokenTree) (int, error) {
	tk, err := dec.Token()
	if err != nil {
		return 0, err
	}
	var length int
	switch tk := tk.(type) {
	case json.Delim:
		switch tk {
		case '{':
			startOffset := int(dec.InputOffset())
			err = offsetObject(dec, tree.children)
			if err != nil {
				return 0, err
			}
			// Consumes the ending delim
			if _, err := dec.Token(); err != nil {
				return 0, err
			}
			endOffset := int(dec.InputOffset())
			length = endOffset - startOffset + 1
		case '[':
			startOffset := int(dec.InputOffset())
			err = offsetArray(dec, tree.children)
			if err != nil {
				return 0, err
			}
			// Consumes the ending delim
			if _, err := dec.Token(); err != nil {
				return 0, err
			}
			endOffset := int(dec.InputOffset())
			length = endOffset - startOffset + 1
		default:
			return 0, fmt.Errorf("unexpected delim token %#v", tk)
		}
	case bool:
		if tk {
			length = 4 // true
		} else {
			length = 5 // false
		}
	case json.Number:
		length = len(tk.String())
	case string:
		length = len(tk) + 2 // quotes
	case nil:
		length = 4 // null
	default:
		return 0, fmt.Errorf("invalid token %#v", tk)
	}
	return length, nil
}

func offsetObject(dec *json.Decoder, trees map[string]*tokenTree) error {
	var tree *tokenTree
	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			return err
		}
		switch tk := tk.(type) {
		case string:
			var ok bool
			tree, ok = trees[tk]
			if !ok {
				if err := drainValue(dec); err != nil {
					return err
				}
				continue
			}
			length, err := offsetValue(dec, tree)
			if err != nil {
				return err
			}
			offset := int(dec.InputOffset()) - length
			tree.offset = &offset
		default:
			return fmt.Errorf("invalid object key token %#v", tk)
		}
	}
	return nil
}

func offsetArray(dec *json.Decoder, trees map[string]*tokenTree) error {
	i := -1
	for dec.More() {
		i++
		idx := strconv.Itoa(i)
		tree, ok := trees[idx]
		if !ok {
			if err := drainValue(dec); err != nil {
				return err
			}
			continue
		}
		length, err := offsetValue(dec, tree)
		if err != nil {
			return err
		}
		offset := int(dec.InputOffset()) - length
		tree.offset = &offset
	}
	return nil
}

// drainValue drains a single value, including object and array.
func drainValue(dec *json.Decoder) error {
	tk, err := dec.Token()
	if err != nil {
		return err
	}

	switch tk := tk.(type) {
	case json.Delim:
		switch tk {
		case '{':
			if err := drainInContainer(dec); err != nil {
				return err
			}
		case '[':
			if err := drainInContainer(dec); err != nil {
				return err
			}
		}
	}
	return nil
}

// drainInContainer drains a json container (object/array) by assuming the beginning delimiter is consumed.
func drainInContainer(dec *json.Decoder) error {
	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			return err
		}
		switch tk := tk.(type) {
		case json.Delim:
			switch tk {
			case '{':
				if err := drainInContainer(dec); err != nil {
					return err
				}
			case '[':
				if err := drainInContainer(dec); err != nil {
					return err
				}
			}
		}
	}
	// Consumes the ending delim
	if _, err := dec.Token(); err != nil {
		return err
	}
	return nil
}
