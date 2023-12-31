package mock

import (
	"reflect"
	"strings"

	"github.com/antlabs/mock/city"
	"github.com/antlabs/mock/country"
	"github.com/antlabs/mock/gid"
	"github.com/antlabs/mock/ipv4"
	"github.com/antlabs/mock/name"
	"github.com/antlabs/mock/timex"
	"github.com/antlabs/mock/urlx"
)

const (
	URL      = "url"
	UserName = "username"
	NickName = "nickname"
	ID       = "id"
	Time     = "time"
	fixEmail = "email"
	Country  = "country"
	Ipv4     = "ipv4"
	// 省
	Province = "province"
	// 市
	City = "city"
	// 区
	District = "district"
)

// 猜测string的实际类型
// 第一个返回值是错误, 第二个返回值决定是否继续猜测, false表示不继续猜测
func guessStringType(v reflect.Value, sf reflect.StructField, opt *Options) (err error, ok bool) {
	fieldName := sf.Name

	// TODO 优化成map查表函数
	if strings.Contains(strings.ToLower(fieldName), URL) {
		v.SetString(urlx.URL())
		return nil, false
	}

	// 如果字段名是Name，那么就随机生成一个名字
	if strings.Contains(strings.ToLower(fieldName), UserName) {
		// TODO 需要修改
		v.SetString(name.Name(name.WithChinese()))
		return nil, false
	}

	// 昵称
	if strings.Contains(strings.ToLower(fieldName), NickName) {
		v.SetString(name.Name(name.WithChinese()))
		return nil, false
	}

	// 如果字段名是ID，那么就生成一个uuid
	// 忽略大小写搜索id
	if strings.Contains(strings.ToLower(fieldName), ID) {
		v.SetString(gid.GID())
		return nil, false
	}

	// 如果字段名是Time，那么就随机生成一个时间
	if strings.Contains(strings.ToLower(fieldName), Time) {
		v.SetString(timex.TimeRFC3339String(timex.WithMin(opt.Min), timex.WithMax(opt.Max)))
		return nil, false
	}

	// 如果字段名是email，那么就随机生成一个email
	if strings.Contains(strings.ToLower(fieldName), fixEmail) {
		e, err := Email()
		if err != nil {
			return err, false
		}
		v.SetString(e)
		return nil, false
	}

	// 如果字段是country, 那么就随机生成一个国家
	if strings.Contains(strings.ToLower(fieldName), Country) {
		v.SetString(country.Country(opt.CountryChina))
		return nil, false
	}

	// 如果字段是country, 那么就随机生成一个国家
	if strings.Contains(strings.ToLower(fieldName), Ipv4) {
		v.SetString(ipv4.IPv4())
		return nil, false
	}

	// 省
	if strings.Contains(strings.ToLower(fieldName), Province) {
		province := city.Province()
		v.SetString(province)
		if opt.Province == "" {
			opt.Province = province
		}
		return nil, false
	}

	// 市
	if strings.Contains(strings.ToLower(fieldName), City) {
		city := city.City(city.WithProvinceName(opt.Province))
		v.SetString(city)
		if opt.City == "" {
			opt.City = city
		}
		return nil, false
	}

	// 区
	if strings.Contains(strings.ToLower(fieldName), District) {
		v.SetString(city.District(city.WithCityName(opt.City)))
		return nil, false
	}
	return nil, true
}
