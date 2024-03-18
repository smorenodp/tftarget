package terraform

func diffMap(old, new map[string]interface{}) (added, removed, updated string) {
	aux := copyMap(new)
	for key, value := range old {
		switch value.(type) {
			case
		}
	}

}

func copyMap(m map[string]interface{}) map[string]interface{} {
	aux := make(map[string]interface{})
	for key, value := range m {
		aux[key] = value
	}
	return aux
}
