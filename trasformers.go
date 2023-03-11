package flatr

import "fmt"

type Transformer func(Entry) (Entry, error)

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

func MaxDeep(h int) Transformer {
	return func(e Entry) (Entry, error) {
		if e.H >= h {
			e.Stop = true
		}
		return e, nil
	}
}
