package gopiper

import (
	"reflect"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/admpub/regexp2"
	"github.com/webx-top/com"
)

func init() {
	// 验证器类型的过滤器统一用下划线开头（验证不通过时，返回ErrInvalidContent错误便于后续处理）
	RegisterFilter("_required", required, "非空", `_required`, ``)
	RegisterFilter("_email", email, "E-mail地址", `_email`, ``)
	RegisterFilter("_username", username, "用户名(字母/数字/汉字)", `_username`, ``)
	RegisterFilter("_singleline", singleline, "单行文本", `_singleline`, ``)
	RegisterFilter("_mutiline", mutiline, "多行文本", `_mutiline`, ``)
	RegisterFilter("_url", url, "URL", `_url`, ``)
	RegisterFilter("_chinese", chinese, "全是汉字", `_chinese`, ``)
	RegisterFilter("_haschinese", haschinese, "包含汉字", `_haschinese`, ``)
	RegisterFilter("_minsize", minsize, "最小长度", `_minsize(5)`, ``)
	RegisterFilter("_maxsize", maxsize, "最大长度", `_maxsize(5)`, ``)
	RegisterFilter("_size", size, "匹配长度", `_size(5)`, ``)
	RegisterFilter("_alpha", alpha, "字母", `_alpha`, ``)
	RegisterFilter("_alphanum", alphanum, "字母或数字", `_alphanum`, ``)
	RegisterFilter("_numeric", numeric, "纯数字", `_numeric`, ``)
	RegisterFilter("_match", match, "正则匹配", `_match([a-z]+)`, ``)
	RegisterFilter("_unmatch", unmatch, "正则不匹配", `_unmatch([a-z]+)`, ``)
	RegisterFilter("_match2", match2, "正则匹配(兼容Perl5和.NET)", `_match2([a-z]+)`, ``)
	RegisterFilter("_unmatch2", unmatch2, "正则不匹配(兼容Perl5和.NET)", `_unmatch2([a-z]+)`, ``)
}

func required(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		if len(vt) == 0 {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if len(v) == 0 {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func email(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		if !com.IsEmailRFC(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if !com.IsEmailRFC(v) {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil
	default:
		return vt, nil
	}
}

func username(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		if !com.IsUsername(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if !com.IsUsername(v) {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil
	default:
		return vt, nil
	}
}

func singleline(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		if !com.IsSingleLineText(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if !com.IsSingleLineText(v) {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil
	default:
		return vt, nil
	}
}

func mutiline(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		if !com.IsMultiLineText(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if !com.IsMultiLineText(v) {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil
	default:
		return vt, nil
	}
}

func url(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		if !com.IsURL(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if !com.IsURL(v) {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil
	default:
		return vt, nil
	}
}

func chinese(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		if !com.IsChinese(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if !com.IsChinese(v) {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil
	default:
		return vt, nil
	}
}

func haschinese(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		if !com.HasChinese(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if !com.HasChinese(v) {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil
	default:
		return vt, nil
	}
}

func minsize(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	minSize, _ := strconv.Atoi(params.String())
	switch vt := src.Interface().(type) {
	case string:
		if utf8.RuneCountInString(vt) < minSize {
			return vt, ErrInvalidContent
		}
		return vt, nil
	case []string:
		var rt []string
		for _, v := range vt {
			if utf8.RuneCountInString(v) < minSize {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil
	default:
		return vt, nil
	}
}

func maxsize(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	maxSize, _ := strconv.Atoi(params.String())
	switch vt := src.Interface().(type) {
	case string:
		if utf8.RuneCountInString(vt) > maxSize {
			return vt, ErrInvalidContent
		}
		return vt, nil

	case []string:
		var rt []string
		for _, v := range vt {
			if utf8.RuneCountInString(v) > maxSize {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func size(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	size, _ := strconv.Atoi(params.String())
	switch vt := src.Interface().(type) {
	case string:
		if utf8.RuneCountInString(vt) != size {
			return vt, ErrInvalidContent
		}
		return vt, nil

	case []string:
		var rt []string
		for _, v := range vt {
			if utf8.RuneCountInString(v) != size {
				continue
			}
			rt = append(rt, v)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func alpha(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		for _, v := range vt {
			if !com.IsAlpha(v) {
				return vt, ErrInvalidContent
			}
		}
		return vt, nil

	case []string:
		var rt []string
		for _, str := range vt {
			isAlpha := true
			for _, v := range str {
				if !com.IsAlpha(v) {
					isAlpha = false
				}
			}
			if !isAlpha {
				continue
			}
			rt = append(rt, str)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func alphanum(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		for _, v := range vt {
			if !com.IsAlphaNumeric(v) {
				return vt, ErrInvalidContent
			}
		}
		return vt, nil

	case []string:
		var rt []string
		for _, str := range vt {
			isAlphaNumeric := true
			for _, v := range str {
				if !com.IsAlphaNumeric(v) {
					isAlphaNumeric = false
				}
			}
			if !isAlphaNumeric {
				continue
			}
			rt = append(rt, str)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func numeric(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch vt := src.Interface().(type) {
	case string:
		for _, v := range vt {
			if !com.IsNumeric(v) {
				return vt, ErrInvalidContent
			}
		}
		return vt, nil

	case []string:
		var rt []string
		for _, str := range vt {
			isNumeric := true
			for _, v := range str {
				if !com.IsNumeric(v) {
					isNumeric = false
				}
			}
			if !isNumeric {
				continue
			}
			rt = append(rt, str)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func match(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	re, err := regexp.Compile(params.String())
	if err != nil {
		return src.Interface(), err
	}
	switch vt := src.Interface().(type) {
	case string:
		if !re.MatchString(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil

	case []string:
		var rt []string
		for _, str := range vt {
			if !re.MatchString(str) {
				continue
			}
			rt = append(rt, str)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func unmatch(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	re, err := regexp.Compile(params.String())
	if err != nil {
		return src.Interface(), err
	}
	switch vt := src.Interface().(type) {
	case string:
		if re.MatchString(vt) {
			return vt, ErrInvalidContent
		}
		return vt, nil

	case []string:
		var rt []string
		for _, str := range vt {
			if re.MatchString(str) {
				continue
			}
			rt = append(rt, str)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func match2(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	re, err := regexp2.Compile(params.String(), 0)
	if err != nil {
		return src.Interface(), err
	}
	switch vt := src.Interface().(type) {
	case string:
		if ok, _ := re.MatchString(vt); !ok {
			return vt, ErrInvalidContent
		}
		return vt, nil

	case []string:
		var rt []string
		for _, str := range vt {
			if ok, _ := re.MatchString(str); !ok {
				continue
			}
			rt = append(rt, str)
		}
		return vt, nil

	default:
		return vt, nil
	}
}

func unmatch2(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	re, err := regexp2.Compile(params.String(), 0)
	if err != nil {
		return src.Interface(), err
	}
	switch vt := src.Interface().(type) {
	case string:
		if ok, _ := re.MatchString(vt); ok {
			return vt, ErrInvalidContent
		}
		return vt, nil

	case []string:
		var rt []string
		for _, str := range vt {
			if ok, _ := re.MatchString(str); ok {
				continue
			}
			rt = append(rt, str)
		}
		return vt, nil

	default:
		return vt, nil
	}
}
