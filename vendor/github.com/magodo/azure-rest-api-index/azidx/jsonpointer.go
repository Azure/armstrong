package azidx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/jsonpointer"
)

func JSONPointerOffset(p jsonpointer.Pointer, document string) (int64, error) {
	dec := json.NewDecoder(strings.NewReader(document))
	var offset int64
	for _, ttk := range p.DecodedTokens() {
		tk, err := dec.Token()
		if err != nil {
			return 0, err
		}
		switch tk := tk.(type) {
		case json.Delim:
			switch tk {
			case '{':
				offset, err = offsetSingleObject(dec, ttk)
				if err != nil {
					return 0, err
				}
			case '[':
				offset, err = offsetSingleArray(dec, ttk)
				if err != nil {
					return 0, err
				}
			default:
				return 0, fmt.Errorf("invalid token %#v", tk)
			}
		default:
			return 0, fmt.Errorf("invalid token %#v", tk)
		}
	}
	return offset, nil
}

func offsetSingleObject(dec *json.Decoder, decodedToken string) (int64, error) {
	for dec.More() {
		offset := dec.InputOffset()
		tk, err := dec.Token()
		if err != nil {
			return 0, err
		}
		switch tk := tk.(type) {
		case json.Delim:
			switch tk {
			case '{':
				if err := drainSingle(dec); err != nil {
					return 0, err
				}
			case '[':
				if err := drainSingle(dec); err != nil {
					return 0, err
				}
			}
		case string:
			if tk == decodedToken {
				return offset, nil
			}
		default:
			return 0, fmt.Errorf("invalid token %#v", tk)
		}
	}
	return 0, fmt.Errorf("token reference %q not found", decodedToken)
}

func offsetSingleArray(dec *json.Decoder, decodedToken string) (int64, error) {
	idx, err := strconv.Atoi(decodedToken)
	if err != nil {
		return 0, fmt.Errorf("token reference %q is not a number: %v", decodedToken, err)
	}
	var i int
	for i = 0; i < idx && dec.More(); i++ {
		tk, err := dec.Token()
		if err != nil {
			return 0, err
		}
		switch tk := tk.(type) {
		case json.Delim:
			switch tk {
			case '{':
				if err := drainSingle(dec); err != nil {
					return 0, err
				}
			case '[':
				if err := drainSingle(dec); err != nil {
					return 0, err
				}
			}
		}
	}
	if !dec.More() {
		return 0, fmt.Errorf("token reference %q not found", decodedToken)
	}
	return dec.InputOffset(), nil
}

// drainSingle drains a single level of object or array.
// The decoder has to guarantee the begining delim (i.e. '{' or '[') has been consumed.
func drainSingle(dec *json.Decoder) error {
	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			return err
		}
		switch tk := tk.(type) {
		case json.Delim:
			switch tk {
			case '{':
				if err := drainSingle(dec); err != nil {
					return err
				}
			case '[':
				if err := drainSingle(dec); err != nil {
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
