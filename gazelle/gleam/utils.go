package gleam

import "github.com/bazelbuild/bazel-gazelle/rule"

func mapper[T any, O any](things []T, mapper func(thing T) O) []O {
    result := make([]O, 0, len(things))
    for _, thing := range things {
        result = append(result, mapper(thing))
    }
    return result
}

func filter[T any](things []T, filter func(thing T) bool) []T {
	result := make([]T, 0, len(things))
	for _, thing := range things {
		if filter(thing) {
			result = append(result, thing)
		}
	}
	return result
}

func asSet[K comparable](things []K) map[K]bool {
	result := make(map[K]bool, len(things))
	for _, thing := range things {
		result[thing] = true
	}
	return result
}

func collect[K comparable, V any](things map[K]V) []K {
	result := make([]K, 0, len(things))
	for k := range things {
		result = append(result, k)
	}
	return result
}

func isGleamLibrary(r *rule.Rule) bool {
	return r.Kind() == "gleam_library" || r.Kind() == "gleam_erl_library" || r.Kind() == "gleam_binary"
}