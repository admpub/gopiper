package gopiper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/admpub/gohttp"
	"github.com/axgle/mahonia"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
)

func TestRegexp(t *testing.T) {
	/*
		for i := 65281; i < 65375; i++ {
			fmt.Println(`[` + string(i) + `] => [` + _tosbc(string(i)) + `]`)
		}
		panic(`[` + string(12288) + `] => [` + _tosbc(string(12288)) + `]`)
	// */
	params := SplitParams(`\,b,c,d`)
	assert.Equal(t, []string{`,b`, `c`, `d`}, params)
	showjson(hrefFilterExp.FindAllStringSubmatch(`href="http://admpub.com"`, -1))
	match, err := hrefFilterExp2.FindStringMatch(`href="http://admpub.com"`)
	if err != nil {
		panic(err)
	}
	showjson(match.Slice())
	result, err := hrefFilterExp2.Replace(`<a href="http://admpub.com/index" ></a><a href="http://admpub.com/admin" ></a>`, "data-url='$2'", 0, -1)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `<a data-url='http://admpub.com/index' ></a><a data-url='http://admpub.com/admin' ></a>`, result)
	pipe := PipeItem{
		Type:     "string-array",
		Selector: "regexp2:^<([a-z]+)>([a-z]+)([\\d]+)</\\1>$",
	}
	if js, err := pipe.parseRegexp("<a>bcdefgh123</a>", true); err != nil {
		t.Fatal(err)
	} else {
		/*
			[
				"bcdefgh",
				"123"
			]
		*/
		fmt.Println(`=== [ TestRegexp ] ===================================\`)
		showjson(js)
		fmt.Println(`=== [ TestRegexp ] ===================================/`)
		fmt.Println()
	}
}

func TestSelector(t *testing.T) {
	js, err := simplejson.NewJson([]byte(`{"value": ["1","2",{"data": ["3", "2", "1"]}]}`))
	if err != nil {
		log.Println(err)
		return
	}
	if js, err = parseJsonSelector(js, "this.value[2].data[1]"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(`=== [ TestSelector ] ===================================\`)
		showjson(js)
		fmt.Println(`=== [ TestSelector ] ===================================/`)
		fmt.Println()
	}
}

func TestJsonPipe(t *testing.T) {
	req := gohttp.New()

	resp, _ := req.Get("http://s.m.taobao.com/search?&q=qq&atype=b&searchfrom=1&from=1&sst=1&m=api4h5").End()

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	pipe := PipeItem{}
	json.Unmarshal([]byte(`
		{
			"type": "array",
			"selector": "this.listItem",
			"subitem": [
				{
					"type": "map",
					"subitem": [
						{
							"name": "now",
							"filter": "unixtime"
						},
						{
							"name": "nowmill",
							"filter": "unixmill"
						},
						{
							"type": "text",
							"selector": "nick",
							"name": "nickname"
						},
						{
							"type": "text",
							"selector": "name",
							"name": "title"
						}
					]
				}
			]
		}
	`), &pipe)

	if val, err := pipe.PipeBytes(body, "json"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(`=== [ TestJsonPipe ] ===================================\`)
		showjson(val)
		fmt.Println(`=== [ TestJsonPipe ] ===================================/`)
		fmt.Println()
	}
}

func TestHtmlDouban(t *testing.T) {
	pb := []byte(`
{
	"type": "map",
	"selector": "",
	"subitem": [
		{
			"type": "string",
			"selector": "title",
			"name": "name",
			"filter": "trimspace|replace((豆瓣))|trim( )"
		},
		{
			"type": "string",
			"selector": "#content .gtleft a.bn-sharing//attr[data-type]",
			"name": "fenlei"
		},
		{
			"type": "string",
			"selector": "#content .gtleft a.bn-sharing//attr[data-pic]",
			"name": "thumbnail"
		},
		{
			"type": "string-array",
			"selector": "#info span.attrs a[rel=v\\:directedBy]",
			"name": "direct"
		},
		{
			"type": "string-array",
			"selector": "#info span a[rel=v\\:starring]",
			"name": "starring"
		},
		{
			"type": "string-array",
			"selector": "#info span[property=v\\:genre]",
			"name": "type"
		},
		{
			"type": "string-array",
			"selector": "#related-pic .related-pic-bd a:not(.related-pic-video) img//attr[src]",
			"name": "imgs",
			"filter": "join($)|replace(albumicon,photo)|split($)"
		},
		{
			"type": "string-array",
			"selector": "#info span[property=v\\:initialReleaseDate]",
			"name": "releasetime"
		},
		{
			"type": "string",
			"selector": "regexp:<span class=\"pl\">单集片长:</span> ([\\w\\W]+?)<br/>",
			"name": "longtime"
		},
		{
			"type": "string",
			"selector": "regexp:<span class=\"pl\">制片国家/地区:</span> ([\\w\\W]+?)<br/>",
			"name": "country",
			"filter": "split(/)|trimspace"
		},
		{
			"type": "string",
			"selector": "regexp:<span class=\"pl\">语言:</span> ([\\w\\W]+?)<br/>",
			"name": "language",
			"filter": "split(/)|trimspace"
		},
		{
			"type": "int",
			"selector": "regexp:<span class=\"pl\">集数:</span> (\\d+)<br/>",
			"name": "episode"
		},
		{
			"type": "string",
			"selector": "regexp:<span class=\"pl\">又名:</span> ([\\w\\W]+?)<br/>",
			"name": "alias",
			"filter": "split(/)|trimspace"
		},
		{
			"type": "string",
			"selector": "#link-report span.hidden, #link-report span[property=v\\:summary]|last",
			"name": "brief",
			"filter": "trimspace|split(\n)|trimspace|wraphtml(p)|join"
		},
		{
			"type": "float",
			"selector": "#interest_sectl .rating_num",
			"name": "score"
		},
		{
			"type": "string",
			"selector": "#content h1 span.year",
			"name": "year",
			"filter": "replace(()|replace())|intval"
		},
		{
			"type": "string",
			"selector": "#comments-section > .mod-hd h2 a",
			"name": "comment",
			"filter": "replace(全部)|replace(条)|trimspace|intval"
		}
	]
}
`)

	log.Println(callFilter("美团他|女神||", `preadd(AAAA)|split(|)|join(,)`))
	if val, err := test_piper("http://movie.douban.com/subject/25850640/", "html", pb); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(`=== [ TestHtmlDouban ] ===================================\`)
		showjson(val)
		fmt.Println(`=== [ TestHtmlDouban ] ===================================/`)
		fmt.Println()
	}

	if val, err := test_piper("http://movie.douban.com/subject/2035218/?from=tag_all", "html", pb); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(`=== [ TestHtmlDouban ] ===================================\`)
		showjson(val)
		fmt.Println(`=== [ TestHtmlDouban ] ===================================/`)
		fmt.Println()
	}

}

func test_piper(u string, tp string, pb []byte, headers ...string) (interface{}, error) {
	pipe := PipeItem{}
	err := json.Unmarshal(pb, &pipe)

	if err != nil {
		return nil, err
	}

	req := gohttp.New()
	req.Get(u)
	for idx := 0; idx < len(headers); idx += 2 {
		req.Set(headers[idx], headers[idx+1])
	}
	req.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	body, _, err := req.Bytes()
	if err != nil {
		return nil, err
	}

	return pipe.PipeBytes(body, tp)
}

func showjson(val interface{}) {
	bd, _ := json.MarshalIndent(val, "", "    ")
	fmt.Println(string(bd))
}

func TestBaidu(t *testing.T) {
	pb := []byte(`
		{
			"type": "jsonparse",
			"selector": "regexp:runtime\\.modsData\\.userData = (\\{(?:.+?)\\});",
			"subitem": [
				{
					"type": "json",
					"selector": "user"
				}
			]
		}
	`)
	val, err := test_piper("https://author.baidu.com/profile?context={%22app_id%22:%221567569757829059%22}&cmdType=&pagelets[]=root&reqID=0&ispeed=1", "text", pb, "Cookie", "BAIDUID=D0FB1501E11F72B20AEC00CED2C220D5:FG=1")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(`=== [ TestBaidu ] ===================================\`)
	showjson(val)
	fmt.Println(`=== [ TestBaidu ] ===================================/`)
	fmt.Println()
}

func TestHtmlJingJiang(t *testing.T) {
	req := gohttp.New()

	resp, _ := req.Get("http://www.jjwxc.net/bookbase_slave.php?submit=&booktype=&opt=&page=3&endstr=&orderstr=4").End()

	/*
		   <table class="cytable" width="986" cellspacing="0" cellpadding="0" border="0" align="center">
		        <tbody>
		            <tr>
		                <td class="sptd" width="128">作者</td>
		                <td class="sptd" width="290">作品</td>
		                <td class="sptd" width="194">类型</td>
		                <td class="sptd" width="50">风格</td>
		                <td class="sptd" width="48">进度</td>
		                <td class="sptd" width="63">字数</td>
		                <td class="sptd" width="73">作品积分</td>
		                <td class="sptd" width="138">发表时间</td>
		            </tr>
		            <tr onmouseover="this.bgColor = '#ffffff';" onmouseout="this.bgColor = '#eefaee';" bgcolor="#eefaee">
		                <td align="left"><a href="oneauthor.php?authorid=335780" target="_blank">雾矢翊</a></td>
		                <td align="left"><a href="onebook.php?novelid=2680387" target="_blank" rel="
		                                               与神兽的秀恩爱日常！<br />标签：灵异神怪 重生 甜文                                                            " class="tooltip">与天同兽</a>
		                </td>
		            	<td align="center">
													   原创-言情-架空历史-爱情
						</td>
		                <td align="center">轻松                                        </td>
		                <td align="center">
		                    <font color="red">已完成</font>                                        </td>
		                <td align="right">2784511</td>
		                <td align="right">3216333568</td>
		                <td align="center">2016-02-13 12:03:00</td>
		            </tr>
		            <tr onmouseover="this.bgColor = '#ffffff';" onmouseout="this.bgColor = '#eefaee';" bgcolor="#eefaee">
		                <td align="left"><a href="oneauthor.php?authorid=821775" target="_blank">开花不结果</a></td>
		                <td align="left"><a href="onebook.php?novelid=3283072" target="_blank" rel="
		                                               她是他的药。<br />标签：甜文 快穿 爽文                                                            " class="tooltip">大佬都爱我 [快穿]</a>
		                </td>
		                <td align="center">
		                                               原创-言情-近代现代-爱情                                      </td>
		            	<td align="center">轻松                                        </td>
		                <td align="center">
		                                               连载中                                        </td>
		                <td align="right">564346</td>
		                <td align="right">1039372096</td>
		                <td align="center">2017-07-28 17:17:49</td>
					</tr>
					其它的省略...
		        </tbody>
		    </table>
	*/

	defer resp.Body.Close()
	bodyGBK, _ := ioutil.ReadAll(resp.Body)

	bodyUTF8 := mahonia.NewDecoder("gb18030").ConvertString(string(bodyGBK))
	body := []byte(bodyUTF8)

	pipe := PipeItem{}
	err := json.Unmarshal([]byte(`
		{
			"type": "array",
			"selector": ".cytable tr:not(:nth-child(1))",
			"subitem": [
						{
							"type": "map",
							"subitem": [
								{
									"type": "text",
									"selector": "td:nth-child(1) a",
									"name":"author"
								},
								{
									"type": "text",
									"selector": "td:nth-child(2) a",
									"name": "name"
								},
								{
									"type": "attr[href]",
									"selector": "td:nth-child(2) a",
									"name": "source",
									"filter": "preadd(http://www.jjwxc.net/)"
								}
							]
						}
				]
		}
	`), &pipe)
	if err != nil {
		log.Println(err)
		return
	}

	v, _ := (pipe.PipeBytes(body, "html"))
	fmt.Println(`=== [ TestHtmlJingJiang ] ===================================\`)
	showjson(v)
	fmt.Println(`=== [ TestHtmlJingJiang ] ===================================/`)
	fmt.Println()
}

func TestBaiduLocal(t *testing.T) {
	req := gohttp.New()

	resp, _ := req.Get("http://www.baidu.com/s?wd=采集关键词&rn=50&tn=baidulocal").End()

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	pipe := PipeItem{}
	err := json.Unmarshal([]byte(`
		{
                "type" : "array",
                "selector": "table td ol table",
                "subitem": [
                    {
                        "type": "map",
                        "subitem": [
                            {
                                "name": "title",
                                "type" : "text",
                                "selector" : "td > a"
                            },
                            {
                                "name" : "url",
                                "type": "href",
                                "selector": "td a"
                            },
                            {
                                "name" : "desc",
                                "selector": "td > font|rm(font[color=\\#008000], font > a)",
                                "type" : "text"
                            }
                        ]
                    }
                ]
		}
	`), &pipe)
	if err != nil {
		log.Println(err)
		return
	}

	v, _ := (pipe.PipeBytes(body, "html"))
	fmt.Println(`=== [ TestBaiduLocal ] ===================================\`)
	showjson(v)
	fmt.Println(`=== [ TestBaiduLocal ] ===================================/`)
	fmt.Println()
}

func TestTTKBPaging(t *testing.T) {

	req := gohttp.New()

	resp, _ := req.Get("http://r.cnews.qq.com/getSubChannels").End()
	/*
	   {
	     "ret": 0,
	     "version": "156.8",
	     "channellist": [
	       {
	         "chlid": "daily_timeline",
	         "chlname": "快报",
	         "recommend": 0,
	         "order": 0,
	         "interfaceType": "timeline"
	       },
	       {
	         "chlid": "kb_video_news",
	         "chlname": "视频",
	         "recommend": 0,
	         "order": 0,
	         "rendtype": "video",
	         "group": "推荐"
	   	},
	   	其它的省略...
	     ],
	     "function": "oldAction"
	   }
	*/
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	pipe := PipeItem{}

	err := json.Unmarshal([]byte(`{
		"selector": "channellist",
		"type": "array",
		"name": "ROOT",
		"filter": "paging(1,2)",
		"subitem": [
			{
				"name": "child",
				"selector": "chlid",
				"type": "text",
				"filter": "sprintf(http://r.cnews.qq.com/getSubNewsChlidInterest?chlid=%s&page={0})"
			}
		]
	}`), &pipe)

	if err != nil {
		log.Println(err)
		return
	}

	v, _ := (pipe.PipeBytes(body, "json"))
	fmt.Println(`=== [ TestTTKBPaging ] ===================================\`)
	showjson(v)
	fmt.Println(`=== [ TestTTKBPaging ] ===================================/`)
	fmt.Println()
}
