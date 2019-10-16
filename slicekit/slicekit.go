package slicekit

// IntInSlice detect if an int value in a slice
func IntInSlice(target int, sli []int) bool {
	for _, element := range sli {
		if element == target {
			return true
		}
	}
	return false
}

// StringInSlice detect if a string value in a string slice
func StringInSlice(target string, sli []string) bool {
	for _, element := range sli {
		if element == target {
			return true
		}
	}
	return false
}

// Int64InSlice detect if an int64 value in a slice
func Int64InSlice(target int64, sli []int64) bool {
	for _, element := range sli {
		if element == target {
			return true
		}
	}
	return false
}

// UniqueIntSlice Unique an int slice and optional filter zero value by optional parameter
func UniqueIntSlice(sli []int, filter ...interface{}) []int {
	sl := make([]int, 0, len(sli))
	var filterZeroValue bool
	if len(filter) == 1 {
		switch v := filter[0].(type) {
		case bool:
			filterZeroValue = v
		}
	}
	uniqueMap := make(map[int]struct{})
	for _, s := range sli {
		if filterZeroValue && s == 0  {
			continue
		}
		_, ok := uniqueMap[s]
		if ok {
			continue
		} else {
			uniqueMap[s] = struct{}{}
		}
		sl = append(sl, s)
	}
	return sl
}

// UniqueIntSlice Unique an int64 slice and optional filter zero value by optional parameter
func UniqueInt64Slice(sli []int64, filter ...interface{}) []int64 {
	sl := make([]int64, 0, len(sli))
	var filterZeroValue bool
	if len(filter) == 1 {
		switch v := filter[0].(type) {
		case bool:
			filterZeroValue = v
		}
	}
	uniqueMap := make(map[int64]struct{})
	for _, s := range sli {
		if filterZeroValue && s == 0 {
			continue
		}
		_, ok := uniqueMap[s]
		if ok {
			continue
		} else {
			uniqueMap[s] = struct{}{}
		}
		sl = append(sl, s)
	}
	return sl
}

// UniqueIntSlice Unique an string slice and optional filter zero value by optional parameter
func UniqueStringSlice(sli []string, filter ...interface{}) []string {
	sl := make([]string, 0, len(sli))
	var filterZeroValue bool
	if len(filter) == 1 {
		switch v := filter[0].(type) {
		case bool:
			filterZeroValue = v
		}
	}
	uniqueMap := make(map[string]struct{})
	for _, s := range sli {
		if filterZeroValue && len(s) == 0 {
			continue
		}
		_, ok := uniqueMap[s]
		if ok {
			continue
		} else {
			uniqueMap[s] = struct{}{}
		}
		sl = append(sl, s)
	}
	return sl
}