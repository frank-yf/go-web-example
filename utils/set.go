package utils

import (
	"github.com/frank-yf/go-web-example/utils/json"
)

type StringSet map[string]struct{}

// NewSet 根据不定长字符串生成一个set结构的数据集合
func NewSet(arr ...string) *StringSet {
	return NewSetFromSlice(arr)
}

// NewSetFromSlice 根据数组生成一个set结构的数据集合
// 即使数组长度为空也会返回一个可用的空集合
func NewSetFromSlice(arr []string) *StringSet {
	arrLen := len(arr)
	set := make(StringSet, arrLen)
	for _, s := range arr {
		set[s] = struct{}{}
	}
	return &set
}

// UnmarshalJSON 将array结构的Json数据反序列化为set结构
func (s *StringSet) UnmarshalJSON(data []byte) (err error) {
	var arr []string
	err = json.Unmarshal(data, &arr)
	if err != nil {
		return
	}
	*s = *NewSetFromSlice(arr)
	return
}

// MarshalJSON 将set集合序列化为array结构的Json数据
func (s *StringSet) MarshalJSON() (bs []byte, err error) {
	var arr []string
	if s != nil {
		arr = make([]string, 0, len(*s))
		for k := range *s {
			arr = append(arr, k)
		}
	}
	return json.Marshal(arr)
}

// IsEmpty 是否为空
func (s StringSet) IsEmpty() bool {
	return len(s) == 0
}

// Has 是否存在元素
func (s StringSet) Has(str string) (ok bool) {
	if str == "" {
		return
	}
	_, ok = s[str]
	return
}

// Intersection 与传入字符串数组是否存在交集
func (s StringSet) Intersection(arr []string) bool {
	for _, str := range arr {
		if s.Has(str) {
			return true
		}
	}
	return false
}
