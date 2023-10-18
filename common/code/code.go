package code

import (
	"math"
)

var CharacterSet = map[int32]int{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'A': 10,
	'B': 11,
	'C': 12,
	'D': 13,
	'E': 14,
	'F': 15,
	'G': 16,
	'H': 17,
	'J': 18,
	'K': 19,
	'L': 20,
	'M': 21,
	'N': 22,
	'P': 23,
	'Q': 24,
	'R': 25,
	'T': 26,
	'U': 27,
	'W': 28,
	'X': 29,
	'Y': 30,
}

type Code struct {
	Code             string
	NationCode       string
	FirstCode        string
	FirstDepartment  string
	SecondCode       string
	SecondDepartment string
	CheckCode        string
}

func NewCode(code string) *Code {
	if len(code) != 44 {
		return &Code{}
	}
	return &Code{
		Code:             code,
		NationCode:       code[:3],
		FirstCode:        code[3:21],
		FirstDepartment:  code[21:23],
		SecondCode:       code[23:41],
		SecondDepartment: code[41:43],
		CheckCode:        code[43:],
	}
}

func (c *Code) SumCheckCode(code string) string {
	var sum = 0
	for index, num := range code {
		var value = CharacterSet[num]
		//计算加权因子
		var weight = int(math.Pow(3, float64(index))) % 31
		sum += value * weight
	}
	var mod = sum % 31
	var sign = func() int {
		if mod == 0 {
			return 0
		}
		return 31 - mod
	}()
	var signChar int32

	for key, value := range CharacterSet {
		signChar = key
		if value == sign {
			break
		}
	}
	var signStr = string(signChar)
	return signStr
}

func (c *Code) isValidCheckCode(code string) bool {
	return c.SumCheckCode(code[:len(code)-1]) == code[len(code)-1:]
}

func (c *Code) IsValid() bool {
	ok1 := c.isValidCheckCode(c.FirstCode)

	ok2 := c.isValidCheckCode(c.SecondCode)

	ok3 := c.isValidCheckCode(c.Code)

	return ok1 && ok2 && ok3
}

func (c *Code) Generate() {

}
