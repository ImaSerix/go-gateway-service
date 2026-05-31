package transformer

func DeepCopy(m map[string]any) map[string]any {

	// Очень странно, если template оказался всё таки nil в этом моменте, но всё возможно
	if m == nil {
		return nil
	}

	out := make(map[string]any)

	for k, v := range m {
		switch val := v.(type) {
		case map[string]any:
			out[k] = DeepCopy(val)
		default:
			out[k] = val
		}
	}

	return out
}
