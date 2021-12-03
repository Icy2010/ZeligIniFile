![image](https://github.com/Icy2010/ZeligIniFile/blob/main/zelig.iniFile.png)
# ZeligIniFile
#### 缩略的 Ini文件 操作
#### 作者 Icy 
#### Web http://zelig.cn
```golang

package main

import (
	"fmt"
	z "github.com/Icy2010/ZeligIniFile"
)

type TContacInfo struct {
	Name   string `ini:"name"`
	Web    string `ini:"web"`
	EMail  string `ini:"email"`
	WeChat string `ini:"wechat"`
	QQ     string `ini:"qq"`
}

func main() {
	zini, err := z.NewZeligIniFromMemory([]byte(`[options]`))
	if err == nil {
		Sec := zini.Section(`options`)
		Sec.SetString(`web`, `https://zelig.cn`)
		Sec.SetString(`name`, `icy`)
		Sec.SetString(`email`, `icy2010@hotmail.com`)
		Sec.SetString(`wechat`, `IcySoft`)
		Sec.SetString(`qq`, `2261206`)

		Val := Sec.String(`name`)
		fmt.Println(Val)

		Val = Sec.String(`nick`, `meow`)
		fmt.Println(Val)

		Val, err = zini.FindString(`options.web`)
		fmt.Printf("Web: %s\n", Val)

		Info := TContacInfo{}
		if Sec.Struct(&Info) == nil {
			fmt.Println(Info)
		}

		err = zini.SaveTo(`config.ini`)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

  
```
