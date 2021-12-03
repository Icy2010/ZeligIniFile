![image](https://github.com/Icy2010/ZeligIniFile/blob/main/zelig.iniFile.png)
# ZeligIniFile
#### 缩略的 Ini文件 操作
#### 作者 Icy 
#### Web http://zelig.cn
```golang

  type TTestData struct {
    StringVal string  `ini:"string_value"`
    IntVal    int64   `ini:"integer_value"`
    FloatVal  float64 `ini:"Float_value"`
  }

	zini, err := NewZeligFromFile("config.cfg") //如果没有这个文件不存在的 或者错误的 默认会生成一个 [general]的段
	t.Error(err)
	sec := zini.Section(`general`)
	sec.SetString("test", "heihei") // 如果没有 新增一个  如果有 修改
	sec.SetFloat(`fee`, 100.12)
	t.Error(zini.Save()) // 保存到 创建时候输入的 如果不需要  SaveTo(`config.cfg`)

	t.Log(zini.FindString(`general.test`))
	t.Log(sec.Int(`heihei`, 100))
	t.Log(zini.FindFloat(`general.fee`))

	t.Log(sec.ToJson())
	t.Error(sec.SetStruct(TTestData{
		StringVal: "test",
		IntVal:    1000,
		FloatVal:  10.10,
	}))

	data := TTestData{}
	t.Error(sec.Struct(&data))
	t.Log(data)
	t.Error(zini.Save())
  
```
