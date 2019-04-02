package util

// 一些小算法

import (
	"fmt"
	"reflect"
	"sort"
)

type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].Value > p[j].Value // > : 大到小,<: 小到大
}

// map按value值升序排序
func SortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0

	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	// 保证排序的稳定性, 减少时间复杂度
	sort.Stable(p)
	return p
}

// 通过反射取得变量类型
func TypeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func InSliceMap(slice []interface{}) map[interface{}]bool {
	res := make(map[interface{}]bool)

	for _, s := range slice {
		k, ok := s.(string)
		if !ok {
			panic("s is not string")
		} else {
			res[k] = true
		}
	}

	return res
}

// 计算飞行时间,时:分
func CountFlyTime(seconds uint64) string {

	m := seconds / 60

	h := m / 60
	m = m % 60
	return fmt.Sprintf("%02d:%02d", h, m)

}
