package gopiper

import (
	"errors"
	"fmt"
	"html"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegisterFilter("preadd", preadd)
	RegisterFilter("postadd", postadd)
	RegisterFilter("replace", replace)
	RegisterFilter("split", split)
	RegisterFilter("join", join)
	RegisterFilter("trim", trim)
	RegisterFilter("trimspace", trimspace)
	RegisterFilter("substr", substr)
	RegisterFilter("intval", intval)
	RegisterFilter("floatval", floatval)
	RegisterFilter("hrefreplace", hrefreplace)
	RegisterFilter("wraphtml", wraphtml)
	RegisterFilter("tosbc", tosbc)
	RegisterFilter("unescape", unescape)
	RegisterFilter("escape", escape)
	RegisterFilter("sprintf", sprintf)
	RegisterFilter("sprintfmap", sprintfmap)
	RegisterFilter("unixtime", unixtime)
	RegisterFilter("unixmill", unixmill)
	RegisterFilter("paging", paging)
	RegisterFilter("quote", quote)
	RegisterFilter("unquote", unquote)
}

type FilterFunction func(src *reflect.Value, params *reflect.Value) (interface{}, error)

var filters = make(map[string]FilterFunction)

func RegisterFilter(name string, fn FilterFunction) {
	_, existing := filters[name]
	if existing {
		panic(fmt.Sprintf("Filter with name '%s' is already registered.", name))
	}
	filters[name] = fn
}

func ReplaceFilter(name string, fn FilterFunction) {
	_, existing := filters[name]
	if !existing {
		panic(fmt.Sprintf("Filter with name '%s' does not exist (therefore cannot be overridden).", name))
	}
	filters[name] = fn
}

var (
	filterExp     = regexp.MustCompile(`([a-zA-Z0-9\-_]+)(?:\(([\w\W]*?)\))?(\||$)`)
	hrefFilterExp = regexp.MustCompile(`href(\s*)=(\s*)([\w\W]+?)"`)
)

func applyFilter(name string, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	fn, existing := filters[name]
	if !existing {
		return nil, fmt.Errorf("Filter with name '%s' not found.", name)
	}
	return fn(src, params)
}

func callFilter(src interface{}, value string) (interface{}, error) {

	if src == nil || len(value) == 0 {
		return src, nil
	}

	vt := filterExp.FindAllStringSubmatch(value, -1)

	for _, v := range vt {
		if len(v) < 3 {
			continue
		}
		name := v[1]
		params := v[2]

		srcValue := reflect.ValueOf(src)
		paramValue := reflect.ValueOf(params)
		next, err := applyFilter(name, &srcValue, &paramValue)
		if err != nil {
			continue
		}
		src = next

	}

	return src, nil
}

func preadd(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return params.String() + src.String(), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = params.String() + vt[i]
		}
		return vt, nil
	}
	return params.String(), nil
}
func postadd(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return src.String() + params.String(), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = vt[i] + params.String()
		}
		return vt, nil
	}
	return params.String(), nil
}
func _substr(src string, params *reflect.Value) string {
	vt := strings.Split(params.String(), ",")
	switch len(vt) {
	case 1:
		start, _ := strconv.Atoi(vt[0])
		return src[start:]
	case 2:
		start, _ := strconv.Atoi(vt[0])
		end, _ := strconv.Atoi(vt[1])
		return src[start:end]
	}
	return src
}
func substr(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return _substr(src.String(), params), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = _substr(vt[i], params)
		}
		return vt, nil
	}
	return src.Interface(), nil
}
func _replace(src string, params *reflect.Value) string {
	vt := strings.Split(params.String(), ",")
	switch len(vt) {
	case 1:
		return strings.Replace(src, vt[0], "", -1)
	case 2:
		return strings.Replace(src, vt[0], vt[1], -1)
	case 3:
		n, _ := strconv.Atoi(vt[2])
		return strings.Replace(src, vt[0], vt[1], n)
	}
	return src
}
func replace(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return _replace(src.String(), params), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = _replace(vt[i], params)
		}
		return vt, nil
	}
	return src.Interface(), nil
}
func trim(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	if params == nil {
		return src.Interface(), errors.New("filter trim nil params")
	}
	switch src.Interface().(type) {
	case string:
		return strings.Trim(src.String(), params.String()), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = strings.Trim(vt[i], params.String())
		}
		return vt, nil
	}
	return src.Interface(), nil
}

func trimspace(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return strings.TrimSpace(src.String()), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = strings.TrimSpace(vt[i])
		}
		return vt, nil
	}

	return src.Interface(), nil
}

func split(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	if params == nil {
		return src.Interface(), errors.New("filter split nil params")
	}
	switch src.Interface().(type) {
	case string:
		str := strings.TrimSpace(src.String())
		if len(str) == 0 {
			return []string{}, nil
		}
		return strings.Split(str, params.String()), nil
	case []string:
		vt, _ := src.Interface().([]string)
		rs := make([][]string, len(vt))
		for i := 0; i < len(vt); i++ {
			str := strings.TrimSpace(vt[i])
			if len(str) == 0 {
				rs[i] = []string{}
			} else {
				rs[i] = strings.Split(str, params.String())
			}
		}
		return rs, nil
	}

	return src.Interface(), nil
}

func join(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	if params == nil {
		return src.Interface(), errors.New("filter split nil params")
	}
	switch src.Interface().(type) {
	case []string:
		vt, _ := src.Interface().([]string)
		rs := make([]string, 0)
		for _, v := range vt {
			if len(v) > 0 {
				rs = append(rs, v)
			}
		}
		return strings.Join(rs, params.String()), nil
	}

	return src.Interface(), nil
}

func intval(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return strconv.Atoi(src.String())

	case []string:
		vt, _ := src.Interface().([]string)
		rs := make([]int, len(vt))
		for i := 0; i < len(vt); i++ {
			v, _ := strconv.Atoi(vt[i])
			rs[i] = v
		}
		return rs, nil
	}
	return 0, nil
}

func floatval(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return strconv.ParseFloat(src.String(), 64)

	case []string:
		vt, _ := src.Interface().([]string)
		rs := make([]float64, len(vt))
		for i := 0; i < len(vt); i++ {
			v, _ := strconv.ParseFloat(vt[i], 64)
			rs[i] = v
		}
		return rs, nil
	}
	return 0.0, nil
}

func hrefreplace(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return hrefFilterExp.ReplaceAllString(src.String(), params.String()), nil

	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = hrefFilterExp.ReplaceAllString(vt[i], params.String())
		}
		return vt, nil
	}
	return src.Interface(), nil
}

func regexpreplace(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return src.Interface(), nil
}

func _tosbc(src string, params *reflect.Value) string {
	var res string
	for _, t := range src {
		if t == 12288 {
			t = 32
		} else if t > 65280 && t < 65375 {
			t = t - 65248
		}
		res += string(t)
	}
	return res
}

func tosbc(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return _tosbc(src.String(), params), nil

	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = _tosbc(vt[i], params)
		}
		return vt, nil
	}
	return src.Interface(), nil
}

func unescape(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return html.UnescapeString(src.String()), nil

	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = html.UnescapeString(vt[i])
		}
		return vt, nil
	}
	return src.Interface(), nil
}

func escape(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	switch src.Interface().(type) {
	case string:
		return html.EscapeString(src.String()), nil

	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = html.EscapeString(vt[i])
		}
		return vt, nil
	}
	return src.Interface(), nil
}

func wraphtml(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	if params == nil {
		return src.Interface(), errors.New("filter wraphtml nil params")
	}

	switch src.Interface().(type) {
	case string:
		return fmt.Sprintf("<%s>%s</%s>", params.String(), src.String(), params.String()), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = fmt.Sprintf("<%s>%s</%s>", params.String(), vt[i], params.String())
		}
		return vt, nil
	}

	return src.Interface(), nil
}

func sprintf_multi_param(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	if params == nil {
		return src.Interface(), errors.New("filter split nil params ")
	}

	if src.Type().Kind() == reflect.Array || src.Type().Kind() == reflect.Slice {
		count := strings.Count(params.String(), "%")
		ret := make([]interface{}, 0)
		for i := 0; i < src.Len(); i++ {
			ret = append(ret, src.Index(i).Interface())
		}
		if len(ret) > count {
			return fmt.Sprintf(params.String(), ret[:count]...), nil
		}
		return fmt.Sprintf(params.String(), ret...), nil
	}

	return fmt.Sprintf(params.String(), src.Interface()), nil
}
func sprintf(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	if params == nil {
		return src.Interface(), errors.New("filter split nil params")
	}
	switch src.Interface().(type) {
	case string:
		return fmt.Sprintf(params.String(), src.String()), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			if len(vt[i]) <= 0 {
				continue
			}
			vt[i] = fmt.Sprintf(params.String(), vt[i])
		}
		return vt, nil
	}

	return fmt.Sprintf(params.String(), src.Interface()), nil
}
func sprintfmap(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	if params == nil {
		return src.Interface(), errors.New("filter split nil params")
	}
	msrc, ok := src.Interface().(map[string]interface{})
	if ok == false {
		return src.Interface(), errors.New("value is not map[string]interface{}")
	}
	vt := strings.Split(params.String(), ",")
	if len(vt) <= 1 {
		return src.Interface(), errors.New("params length must > 1")
	}
	pArray := []interface{}{}
	for _, x := range vt[1:] {
		if vm, ok := msrc[x]; ok {
			pArray = append(pArray, vm)
		}
	}
	return fmt.Sprintf(vt[0], pArray...), nil
}

func unixtime(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return time.Now().Unix(), nil
}

func unixmill(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return time.Now().UnixNano() / int64(time.Millisecond), nil
}

func paging(src *reflect.Value, params *reflect.Value) (interface{}, error) {
	if params == nil {
		return src.Interface(), errors.New("filter paging nil params")
	}
	srcType := src.Type().Kind()
	if srcType != reflect.Slice && srcType != reflect.Array && srcType != reflect.String {
		return src.Interface(), errors.New("value is not slice ,array or string")
	}
	vt := strings.Split(params.String(), ",")
	if len(vt) < 2 {
		return src.Interface(), errors.New("params length must > 1")
	}

	start, err := strconv.Atoi(vt[0])
	end, err := strconv.Atoi(vt[1])
	if err != nil {
		return src.Interface(), errors.New("params type error:need int." + err.Error())
	}

	offset := -1
	if len(vt) == 3 {
		offset, err = strconv.Atoi(vt[2])
		if err != nil {
			return src.Interface(), errors.New("params type error:need int." + err.Error())
		}
		if offset < 1 {
			return src.Interface(), errors.New("offset must > 0")
		}
	}

	var result []string
	switch src.Interface().(type) {
	case []interface{}:
		vt, _ := src.Interface().([]interface{})
		for i := start; i <= end; i++ {
			for j := 0; j < len(vt); j++ {
				if offset > 0 {
					result = append(result, sprintf_replace(vt[j].(string), []string{strconv.Itoa(i * offset), strconv.Itoa((i + 1) * offset)}))
				} else {
					result = append(result, sprintf_replace(vt[j].(string), []string{strconv.Itoa(i)}))
				}
			}

		}
		return result, nil

	case []string:
		vt, _ := src.Interface().([]string)
		for i := start; i <= end; i++ {
			for j := 0; j < len(vt); j++ {
				if offset > 0 {
					result = append(result, sprintf_replace(vt[i], []string{strconv.Itoa(i * offset), strconv.Itoa((i + 1) * offset)}))
				} else {
					result = append(result, sprintf_replace(vt[i], []string{strconv.Itoa(i)}))
				}
			}

		}
		return result, nil

	case string:
		msrc1, ok := src.Interface().(string)
		if ok == true {
			for i := start; i <= end; i++ {
				if offset > 0 {
					result = append(result, sprintf_replace(msrc1, []string{strconv.Itoa(i * offset), strconv.Itoa((i + 1) * offset)}))
				} else {
					result = append(result, sprintf_replace(msrc1, []string{strconv.Itoa(i)}))
				}
			}
			return result, nil
		}
	}
	return src.Interface(), errors.New("do nothing,src type not support!")
}

func sprintf_replace(src string, param []string) string {
	for i := range param {
		src = strings.Replace(src, "{"+strconv.Itoa(i)+"}", param[i], -1)
	}
	return src
}

func quote(src *reflect.Value, params *reflect.Value) (interface{}, error) {

	switch src.Interface().(type) {
	case string:
		return strconv.Quote(src.String()), nil
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i] = strconv.Quote(vt[i])
		}
		return vt, nil
	}

	return src.Interface(), nil
}

func unquote(src *reflect.Value, params *reflect.Value) (interface{}, error) {

	switch src.Interface().(type) {
	case string:
		return strconv.Unquote(`"` + src.String() + `"`)
	case []string:
		vt, _ := src.Interface().([]string)
		for i := 0; i < len(vt); i++ {
			vt[i], _ = strconv.Unquote(`"` + vt[i] + `"`)
		}
		return vt, nil
	}

	return src.Interface(), nil
}
