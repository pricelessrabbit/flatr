package flat

import "fmt"

type Trasformer func(any) any

func UseFieldAsIndex(idKey string) Trasformer {
	return func(n any) any {
		elements, ok := n.([]any)
		if !ok {
			panic("UseFieldAsIndex support only slices of objects")
		}
		transformed := make(map[string]any)
		for _, element := range elements {
			castedElement, ok := element.(map[string]any)
			keyValue := fmt.Sprintf("%v", castedElement[idKey])
			if !ok || keyValue == "" {
				panic("key value must be a non-empty string")
			}
			transformed[keyValue] = element
		}
		return transformed
	}
}
