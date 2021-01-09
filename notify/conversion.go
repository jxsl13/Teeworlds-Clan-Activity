package notify

import "sort"

func setToSortedList(set map[string]bool) []string {
	strList := make([]string, 0, len(set))
	for k := range set {
		strList = append(strList, k)
	}

	sort.Strings(strList)
	return strList
}
