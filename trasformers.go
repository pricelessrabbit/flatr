package flatr

import "fmt"

// Transformer custom function called by Flatter to preprocess nodes
type Transformer func(Entry) (Entry, error)

// UseFieldAsIndex applied to a list of maps, uses list element field as the list index.
// the list is transformed in a map with the choosen field as key and the element as value
func UseFieldAsIndex(idKey string) Transformer {
	return func(e Entry) (Entry, error) {
		elements, ok := e.V.([]any)
		if !ok {
			return e, fmt.Errorf("on %s: UseFieldAsIndex support only slices of objects", e.K)
		}
		transformed := make(map[string]any)
		for _, element := range elements {
			castedElement, ok := element.(map[string]any)
			keyValue := fmt.Sprintf("%v", castedElement[idKey])
			if !ok || keyValue == "" {
				return e, fmt.Errorf("on %s:  key value must be a non-empty string", e.K)
			}
			transformed[keyValue] = element
		}
		e.V = transformed
		return e, nil
	}
}

// MaxDeep add a limit to the maximum tree-height of an entry to be processed
// if the height exceeds the limit the flattening stops and node is saved unflatted
func MaxDeep(h int) Transformer {
	return func(e Entry) (Entry, error) {
		if e.H >= h {
			e.Stop = true
		}
		return e, nil
	}
}
