package flatr

import "fmt"

type Trasformer func(entry) (entry, error)

func UseFieldAsIndex(idKey string) Trasformer {
	return func(e entry) (entry, error) {
		elements, ok := e.v.([]any)
		if !ok {
			return e, fmt.Errorf("on %s: UseFieldAsIndex support only slices of objects", e.k)
		}
		transformed := make(map[string]any)
		for _, element := range elements {
			castedElement, ok := element.(map[string]any)
			keyValue := fmt.Sprintf("%v", castedElement[idKey])
			if !ok || keyValue == "" {
				return e, fmt.Errorf("on %s:  key value must be a non-empty string", e.k)
			}
			transformed[keyValue] = element
		}
		e.v = transformed
		return e, nil
	}
}
