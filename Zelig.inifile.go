package Zelig_IniFile

/*
   缩略的 ini 操作,纯属为了自己的习惯上使用
   icy
   zelig.cn
*/
import (
	"encoding/json"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type TIniSection struct {
	section *ini.Section
}

func (this *TIniSection) ToJson() string {
	keys := this.section.Keys()
	if len(keys) > 0 {
		data := make(map[string]string, 0)
		for _, p := range keys {
			data[p.Name()] = p.Value()
		}

		if bytes, err := json.Marshal(data); err == nil {
			return string(bytes)
		}
	}

	return `{}`
}

func (this *TIniSection) Int(key string, defval ...int64) int64 {
	if value_int, err := this.section.Key(key).Int64(); err == nil {
		return value_int
	}

	if len(defval) > 0 {
		return defval[0]
	}

	return 0
}

func (this *TIniSection) Float(key string, defflo ...float64) float64 {
	if value_float, err := this.section.Key(key).Float64(); err == nil {
		return value_float
	}

	if len(defflo) > 0 {
		return defflo[0]
	}

	return 0
}

func (this *TIniSection) String(key string, defstr ...string) string {
	s := this.section.Key(key).String()
	if len(defstr) > 0 && s == "" {
		return defstr[0]
	}

	return s
}

func (this *TIniSection) Bool(key string, defbool ...bool) bool {
	if b, err := strconv.ParseBool(this.section.Key(key).String()); err == nil {
		return b
	}

	if len(defbool) > 0 {
		return defbool[0]
	}

	return false
}

func (this *TIniSection) setValue(key string, value interface{}) {
	if !this.section.HasKey(key) {
		if _, err := this.section.NewKey(key, ""); err != nil {
			return
		}
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Int64, reflect.Bool:
		this.section.Key(key).SetValue(strconv.FormatInt(value.(int64), 10))
	case reflect.Float64:
		this.section.Key(key).SetValue(strconv.FormatFloat(value.(float64), 'f', 2, 64))
	case reflect.String:
		this.section.Key(key).SetValue(value.(string))
	}

}

func (this *TIniSection) SetInt(key string, Value int64) {
	this.setValue(key, Value)
}

func (this *TIniSection) SetString(key, Value string) {
	this.setValue(key, Value)
}

func (this *TIniSection) SetFloat(key string, Value float64) {
	this.setValue(key, Value)
}

func (this *TIniSection) SetBool(key string, Value bool) {
	if Value {
		this.setValue(key, int64(1))
	} else {
		this.setValue(key, int64(0))
	}
}

func (this *TIniSection) DeleteKey(key string) {
	this.section.DeleteKey(key)
}

func (this *TIniSection) SetStruct(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Struct {
		size := val.NumField()
		if size == 0 {
			return fmt.Errorf(`这是一个空的结构体`)
		}

		t := reflect.TypeOf(value)
		for i := 0; i < size; i++ {
			v := val.Field(i)

			name := t.Field(i).Tag.Get(`ini`)
			if name == "" {
				name = t.Field(i).Name
			}

			switch v.Kind() {
			case reflect.Bool:
				this.setValue(name, func() int64 {
					if v.Bool() {
						return 1
					} else {
						return 0
					}
				}())
			case reflect.String:
				this.setValue(name, v.String())
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				this.setValue(name, v.Int())
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				this.setValue(name, int64(v.Uint()))
			case reflect.Float32, reflect.Float64:
				this.setValue(name, v.Float())
			default:
				continue
			}
		}
	} else {
		return fmt.Errorf(`不是一个有效的结构体`)
	}

	return nil
}

func (this *TIniSection) Struct(value interface{}) error {
	vtype := reflect.TypeOf(value).Elem() //获取type的真实类型
	val := reflect.ValueOf(value).Elem()

	if vtype.Kind() == reflect.Struct {
		size := vtype.NumField()
		if size == 0 {
			return fmt.Errorf(`这是一个空的结构体`)
		}
		keys := this.section.Keys()
		if len(keys) > 0 {
			doSet := func(p *ini.Key) {
				for i := 0; i < size; i++ {
					t := vtype.Field(i)
					s := t.Tag.Get(`ini`)
					next := false
					if len(s) > 0 {
						next = strings.EqualFold(s, p.Name())
					} else {
						next = strings.EqualFold(t.Name, p.Name())
					}

					if next {
						v := val.Field(i)
						switch v.Kind() {
						case reflect.Bool:
							v.SetBool(func() bool {
								if i, err := p.Int(); err == nil {
									return i > 0
								} else {
									return false
								}
							}())
						case reflect.String:
							v.SetString(p.String())
						case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
							{
								if i64, err := p.Int64(); err == nil {
									v.SetInt(i64)
								}
							}
						case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							{
								if ui64, err := p.Uint64(); err == nil {
									v.SetUint(ui64)
								}
							}
						case reflect.Float32, reflect.Float64:
							{
								if f64, err := p.Float64(); err == nil {
									v.SetFloat(f64)
								}
							}
						}
					}
				}
			}

			for _, p := range keys {
				doSet(p)
			}
		} else {
			return fmt.Errorf(`段内是控的.`)
		}
	}

	return nil
}

type TZeligIni struct {
	ini      *ini.File
	FileName string
}

func NewZeligIniFromMemory(Bytes []byte) (TZeligIni, error) {
	i, err := ini.Load(Bytes)
	return TZeligIni{
		ini:      i,
		FileName: "",
	}, err
}

func NewZeligIniFromFile(FileName string) (TZeligIni, error) {
	Bytes, err := os.ReadFile(FileName)
	if err != nil {
		Bytes = []byte(`[general]`)
	}

	zini, e := NewZeligIniFromMemory(Bytes)
	zini.FileName = FileName

	return zini, e
}

func (this *TZeligIni) Section(Name string) TIniSection {
	Sec, err := this.ini.GetSection(Name)
	if err != nil {
		Sec = this.ini.Section(Name)
	}
	return TIniSection{section: Sec}
}

func (this *TZeligIni) KeyStrings(section string) []string {
	return this.ini.Section(section).KeyStrings()
}

func (this *TZeligIni) HasKey(section, key string) bool {
	if Sec, err := this.ini.GetSection(section); err == nil {
		return Sec.HasKey(key)
	}

	return false
}

func (this *TZeligIni) HasValue(section, value string) bool {
	if Sec, err := this.ini.GetSection(section); err == nil {
		return Sec.HasValue(value)
	}
	return false
}

func (this *TZeligIni) HasSection(name string) bool {
	return this.ini.HasSection(name)
}

func (this *TZeligIni) DeleteSection(name string) {
	this.ini.DeleteSection(name)
}

func (this *TZeligIni) DeleteKey(section, key string) {
	if Sec, err := this.ini.GetSection(section); err == nil {
		Sec.DeleteKey(key)
	}
}

func (this *TZeligIni) SectionStrings() []string {
	return this.ini.SectionStrings()
}

func (this *TZeligIni) Sections() []TIniSection {
	list := make([]TIniSection, 0)
	Names := this.SectionStrings()
	if len(Names) > 0 {
		for _, v := range Names {
			list = append(list, this.Section(v))
		}
	}
	return list
}

func (this *TZeligIni) getSectionKey(path string) (*ini.Key, error) {
	items := strings.Split(path, `.`)
	if len(items) != 2 {
		return nil, fmt.Errorf(`段键不正确`)
	}

	if Sec, err := this.ini.GetSection(items[0]); err == nil {
		if !Sec.HasKey(items[1]) {
			return nil, fmt.Errorf(`键不正确`)
		}

		return Sec.Key(items[1]), nil
	} else {
		return nil, fmt.Errorf(`段不正确`)
	}
}

func (this *TZeligIni) FindString(path string) (string, error) {
	if key, err := this.getSectionKey(path); err == nil {
		return key.String(), err
	} else {
		return "", err
	}
}

func (this *TZeligIni) FindInt(path string) (int64, error) {
	if key, err := this.getSectionKey(path); err == nil {
		var i int64 = 0
		i, err = key.Int64()
		return i, err
	} else {
		return 0, err
	}
}

func (this *TZeligIni) FindFloat(path string) (float64, error) {
	if key, err := this.getSectionKey(path); err == nil {
		var f float64 = 0
		f, err = key.Float64()
		return f, err
	} else {
		return 0, err
	}
}

func (this *TZeligIni) FindBool(path string) (bool, error) {
	if key, err := this.getSectionKey(path); err == nil {
		var b bool = false
		b, err = key.Bool()
		return b, err
	} else {
		return false, err
	}
}

func (this *TZeligIni) SaveTo(filename string) error {
	return this.ini.SaveTo(filename)
}

func (this *TZeligIni) Save() error {
	return this.ini.SaveTo(this.FileName)
}
