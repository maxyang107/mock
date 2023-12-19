# API数据mock包

说明，本mock包给予github.com/antlabs/mock 二开改造而来,如果对你有所帮助，请给我和原作者点个start吧，谢谢

## 一、安装使用
```
go get github.com/maxyang107/mock
```

### 二、快速开始

```
package main

import (
	"fmt"
	"mock"
)

type Resp struct {
	Age   uint32  `json:"age"`
	Score float64 `json:"score"`
	Rank  int16   `json:"rank"`
	Name  string  `json:"name"`
	Img   string  `json:"img"`
	Email string  `json:"email"`
	Min   float64 `json:"min"`
	Max   int     `json:"max"`
}

func main() {
	var (
		re     = Resp{}
		mockOp = []mock.Option{}
	)
	mockOp = append(mockOp, mock.WithContainsFieldSourceString("Age", []string{"12", "15", "26"})) //设置指定字段的值，随机从后面数组取
	mockOp = append(mockOp, mock.WithContainsFieldSourceString("Score", []string{"12.56", "15.32", "26.1"}))//设置指定字段的值，随机从后面数组取
	mockOp = append(mockOp, mock.WithContainsFieldSourceString("Rank", []string{"1", "13", "2"}))//设置指定字段的值，随机从后面数组取
	mockOp = append(mockOp, mock.WithMinMaxByField("Min", 3, 7), mock.WithFloatLen("%.3f"))//设置指定字段的值的取值范围，并且设置小数取之位数
	mockOp = append(mockOp, mock.WithMinMaxByField("Max", 3, 700))//设置指定字段的值的取值范围
	mockOp = append(mockOp, mock.WithContainsFieldSourceString("Img", []string{"http://www.baidu.com?adjae&daskjewa1321.jpe", "http://www.baidu.com?adjae&daskjewa1321.jpg", "http://www.baidu.com?adjae&daskjewa1321.png"}))
	mock.MockGen(&re, mockOp...)
	fmt.Println(re)
}


re输出：
{26 15.32 13 7ec2d http://www.baidu.com?adjae&daskjewa1321.png 12f1f5b6@yahoo.com 3.052 185}
```


## 三、WithXXX各种配置函数
### 3.1 配置指定字段的数据生成范围`WithMinMaxLenByField`
```go
type Test_MinMaxLenByField struct {
	S     string
	Slice []int
}

// 控制slice的生成长度范围
func TestMinMaxLenByField(t *testing.T) {
	e := Test_MinMaxLenByField{}
	mock.MockData(&e, mock.WithMinMaxLenByField("S", 10, 20), mock.WithMinMaxLenByField("Slice", 10, 20))
}
```

### 3.2 配置指定字段的数据源`WithContainsFieldSourceString`
指定HeadPic字段的，数据源。支持string，int，uint,float中的所有类型
```go
var a ReferenceType
image := []string{"image.xxx.com/1.headpic", "image.xxx.com/2.headpic", "image.xxx.com/3.headpic"}
err := mock.MockData(&a, mock.WithContainsFieldSourceString("HeadPic", image))
```
### 3.3 设置为英文
```go
mock.WithCountryEn()

```

### 3.4 设置数据最大长度`WithMaxLen`
```go
mock.WithMaxLen()
```

### 3.5 设置数据最大长度`WithMinLen`
```go
mock.WithMaxLen()
```

### 3.6 设置数值的最大值`WithMax`
```go
mock.WithMax()
```

### 3.7 设置数值的最大值`WithMin`
```go
mock.WithMin()
```
### 3.8 设置忽略的字段名
字段有时候是由protobuf或者thrift生成，不能直接修改tag，可以使用mock.WithIgnoreFields接口忽略
```go
// 设置忽略的字段名
mock.WithIgnoreFields([]string{"Country", "NickName"})
```

### 3.9 设置浮点数返回小数位数`WithFloatLen`
```go
mock.WithFloatLen()
```


### 4.0 设置int，uint，float三种类型值的取值范围`WithMinMaxByField`
```go
mock.WithMinMaxByField()
```