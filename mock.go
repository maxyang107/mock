package mock

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/antlabs/mock/integer"
	"github.com/antlabs/mock/stringx"
	"github.com/antlabs/mock/timex"
)

func MockGen(resp any, opts ...Option) error {
	if resp == nil {
		return nil
	}

	opt := &Options{}
	opt.MinLen = 1
	for _, o := range opts {
		o(opt)
	}
	defaultOptions(opt)
	return mockData(reflect.ValueOf(resp), reflect.StructField{}, opt)
}
func defaultOptions(opt *Options) {
	if opt.Max == 0 {
		opt.Max = math.MaxInt32
	}

	if opt.MaxLen == 0 {
		opt.MaxLen = 10
	}

	if opt.FloatLen == "" {
		opt.FloatLen = "%.2f"
	}
}

func mockData(v reflect.Value, sf reflect.StructField, opt *Options) error {
	// 忽略
	if len(opt.IgnoreFields) > 0 && len(sf.Name) > 0 && opt.IgnoreFields[sf.Name] {
		return nil
	}

	switch v.Kind() {
	// 指针类型，需要先获取指针指向的值
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		return mockData(v.Elem(), sf, opt)
		// 结构体类型，需要遍历结构体的所有字段
	case reflect.Struct:
		// 如果是time.Time类型，直接返回当前时间
		if v.Type().String() == "time.Time" {
			tv := timex.TimeRFC3339(timex.WithMin(opt.Min), timex.WithMax(opt.Max))
			v.Set(reflect.ValueOf(tv))
			return nil
		}
		typ := v.Type()
		for i := 0; i < v.NumField(); i++ {
			sf := typ.Field(i)
			if sf.PkgPath != "" && !sf.Anonymous {
				continue
			}

			if err := mockData(v.Field(i), sf, opt); err != nil {
				return err
			}

		}

		// slice 或者 array 类型，需要遍历所有元素
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 && reflect.Array == v.Kind() {
			return nil
		}

		// 随机生成一个长度
		l := 0
		if _, ok := opt.MinMaxLenByField[sf.Name]; ok {
			l = integer.IntegerRangeInt(int(opt.MinMaxLenByField[sf.Name].MinLen), int(opt.MinMaxLenByField[sf.Name].MaxLen))
		} else {
			l = integer.IntegerRangeInt(int(opt.MinLen), int(opt.MaxLen))
		}

		// 如果是slice类型，那么就需要扩容
		if reflect.Slice == v.Kind() {
			v.Set(reflect.MakeSlice(v.Type(), l, l))
		}

		for i := 0; i < v.Len(); i++ {
			if err := mockData(v.Index(i), sf, opt); err != nil {
				return err
			}
		}

		// map类型，需要遍历map的所有key
	case reflect.Map:
		if v.Len() > 0 {
			return nil
		}

		// 随机生成map的长度
		minLen := opt.MinLen
		if minLen == 0 {
			minLen = 1
		}

		l := 0
		if _, ok := opt.MinMaxLenByField[sf.Name]; ok {
			l = integer.IntegerRangeInt(int(opt.MinMaxLenByField[sf.Name].MinLen), int(opt.MinMaxLenByField[sf.Name].MaxLen))
		} else {
			l = integer.IntegerRangeInt(int(minLen), int(opt.MaxLen))
		}

		// 创建一个map
		v.Set(reflect.MakeMapWithSize(v.Type(), l))
		// 遍历map的所有key

		for i := 0; i < l; i++ {
			// 创建一个key
			key := reflect.New(v.Type().Key()).Elem()
			// 创建一个value
			value := reflect.New(v.Type().Elem()).Elem()
			// 递归mock key
			if err := mockData(key, sf, opt); err != nil {
				return err
			}
			// 递归mock value
			if err := mockData(value, sf, opt); err != nil {
				return err
			}
			// 设置map的key和value
			v.SetMapIndex(key, value)
		}

		// 接口类型，需要先获取接口的值
	case reflect.Interface:
		// float32, float64 类型
	case reflect.Float32, reflect.Float64:
		if source, ok := opt.StringSource[sf.Name]; ok {
			i := rand.Intn(len(source))
			s := source[i]
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f := integer.Float64Range(float64(opt.Min), float64(opt.Max))
				f, _ = strconv.ParseFloat(fmt.Sprintf(opt.FloatLen, f), 64)
				v.SetFloat(f)
			}
			f, _ = strconv.ParseFloat(fmt.Sprintf(opt.FloatLen, f), 64)
			v.SetFloat(f)
		} else if source, ok := opt.MinMaxByField[sf.Name]; ok {
			f := integer.Float64Range(float64(source.MinLen), float64(source.MaxLen))
			f, _ = strconv.ParseFloat(fmt.Sprintf(opt.FloatLen, f), 64)
			v.SetFloat(f)
		} else {
			f := integer.Float64Range(float64(opt.Min), float64(opt.Max))
			f, _ = strconv.ParseFloat(fmt.Sprintf(opt.FloatLen, f), 64)
			v.SetFloat(f)
		}
		// int... 类型
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if source, ok := opt.StringSource[sf.Name]; ok {
			i := rand.Intn(len(source))
			s := source[i]
			u, _ := strconv.Atoi(s)
			v.SetInt(int64(u))
		} else if source, ok := opt.MinMaxByField[sf.Name]; ok {
			i := integer.IntegerRangeInt(int(source.MinLen), int(source.MaxLen))
			v.SetInt(int64(i))
		} else {
			i := integer.IntegerRangeInt(int(opt.Min), int(opt.Max))
			v.SetInt(int64(i))
		}
		// uint... 类型
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if source, ok := opt.StringSource[sf.Name]; ok {
			i := rand.Intn(len(source))
			s := source[i]
			u, _ := strconv.Atoi(s)
			v.SetUint(uint64(u))
		} else if source, ok := opt.MinMaxByField[sf.Name]; ok {
			u := integer.IntegerRangeUint(uint(source.MinLen), uint(source.MaxLen))
			v.SetUint(uint64(u))
		} else {
			u := integer.IntegerRangeUint(uint(opt.Min), uint(opt.Max))
			v.SetUint(uint64(u))
		}
		// string 类型
	case reflect.String:
		var err error
		var ok bool
		if source, ok := opt.StringSource[sf.Name]; ok {
			i := rand.Intn(len(source))
			v.SetString(source[i])
			return nil
		}
		// 猜测下数据类型
		if err, ok = guessStringType(v, sf, opt); err != nil {
			return err
		} else if !ok {
			return nil
		}

		s, err := stringx.StringRange(opt.MinLen, opt.MaxLen)
		if err != nil {
			return err
		}
		v.SetString(s)
	}
	return nil
}
