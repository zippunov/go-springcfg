package config

// DeepMerge deeply merges two maps of type `map[string]interface{}`
// all fields of the source map override correspondent fields of the target map
// if both target and source maps contain child maps for the given path
// those maps will be also deeply merged.
// We assume that all child maps are also `map[string]interface{}`
func DeepMerge(target, source map[string]interface{}) map[string]interface{} {
	for k, val := range source {
		switch sm := val.(type) {
		case map[string]interface{}:
			switch tm := target[k].(type) {
			case map[string]interface{}:
				target[k] = DeepMerge(tm, sm)
			default:
				target[k] = val
			}
		default:
			target[k] = val
		}
	}
	return target
}

// DeepMergeAll deeply merges array of maps with latter overriding fields in the former map
func DeepMergeAll(maps ...map[string]interface{}) map[string]interface{} {
	target := map[string]interface{}{}
	if len(maps) > 0 {
		target = maps[0]
	}
	if len(maps) < 2 {
		return target
	}
	for _, source := range maps[1:] {
		target = DeepMerge(target, source)
	}
	return target
}
