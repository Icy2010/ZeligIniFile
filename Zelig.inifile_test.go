package Zelig_IniFile

import (
	"testing"
)

type TTestData struct {
	StringVal string  `ini:"string_value"`
	IntVal    int64   `ini:"integer_value"`
	FloatVal  float64 `ini:"Float_value"`
}

func TestTZeligIni(t *testing.T) {
	ini := TZeligIni{}
	ini.ReadFromString(`[test]
string_value = 哈哈 ; 测试1
integer_value = 1 ;测试2 
Float_value = 2.2

[defalut]

hwqqwdqwd = 123
saddas = 123123 ; 嘿嘿`)

	data := TTestData{}
	ini.Struct(`test`, &data)
	t.Log(data)

	ini.SaveToFile(`test.ini`)

	ini.ClearSection()
	sec := ini.AddSection("options")
	sec.SetInt("integer_value", 20220401)
	sec.SetString("string_value", "这是一个测试咯")
	sec.SetComment("string_value", "看看咯")
	sec.SetFloat("Float_value", 102.23)

	ini.Struct(`options`, &data)
	t.Log(data)

	t.Log(sec.String("string_value", ""))

	initext := ""
	ini.SaveToString(&initext)
	t.Log("\n" + initext)
}
