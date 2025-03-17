package logging

func logParamsToZeroParams(keys map[ExtraKey]any) map[string]any {
	params := map[string]any{}

	for k, v := range keys {
		params[string(k)] = v
	}

	return params
}
