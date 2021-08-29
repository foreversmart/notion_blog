package bdfanyi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	api       = "https://fanyi.baidu.com/v2transapi"
	referer   = "https://fanyi.baidu.com/"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"
)

const (
	GU    = "gu"    // 古吉拉特语
	ROM   = "rom"   // 罗马尼亚语
	GA    = "ga"    // 爱尔兰语
	GL    = "gl"    // 加利西亚语
	LB    = "lb"    // 卢森堡语
	TT    = "tt"    // 塔塔尔语
	TR    = "tr"    // 土耳其语
	TN    = "tn"    // 塞茨瓦纳语
	LT    = "lt"    // 立陶宛语
	TH    = "th"    // 泰语
	TE    = "te"    // 泰卢固语
	FIN   = "fin"   // 芬兰语
	TA    = "ta"    // 泰米尔语
	YO    = "yo"    // 约鲁巴语
	DAN   = "dan"   // 丹麦语
	DE    = "de"    // 德语
	QU    = "qu"    // 凯楚亚语
	EL    = "el"    // 希腊语
	EO    = "eo"    // 世界语
	EN    = "en"    // 英语
	ZH    = "zh"    // 中文
	ARA   = "ara"   // 阿拉伯语
	EU    = "eu"    // 巴斯克语
	ZU    = "zu"    // 祖鲁语
	RU    = "ru"    // 俄语
	EST   = "est"   // 爱沙尼亚语
	JPKA  = "jpka"  // 日语假名
	BE    = "be"    // 白俄罗斯语
	MS    = "may"   // 马来语
	BN    = "bn"    // 孟加拉语
	JP    = "jp"    // 日语
	BS    = "bs"    // 波斯尼亚语
	YUE   = "yue"   // 粤语
	OR    = "or"    // 奥利亚语
	VIE   = "vie"   // 越南语
	CA    = "ca"    // 加泰罗尼亚语
	CY    = "cy"    // 威尔士语
	CS    = "cs"    // 捷克语
	LV    = "lv"    // 拉脱维亚语
	FRA   = "fra"   // 法语
	PT    = "pt"    // 葡萄牙语
	SWE   = "swe"   // 瑞典语
	TL    = "tl"    // 菲律宾语
	PA    = "pa"    // 旁遮普语
	CHT   = "cht"   // 中文繁体
	KOR   = "kor"   // 韩语
	PL    = "pl"    // 波兰语
	HY    = "hy"    // 亚美尼亚语
	HR    = "hr"    // 克罗地亚语
	IU    = "iu"    // 因纽特语
	HU    = "hu"    // 匈牙利语
	HI    = "hi"    // 印地语
	BUL   = "bul"   // 保加利亚语
	HA    = "ha"    // 豪萨语
	UZ    = "uz"    // 乌兹别克语
	MI    = "mi"    // 毛利语
	MK    = "mk"    // 马其顿语
	UR    = "ur"    // 乌尔都语
	MT    = "mt"    // 马耳他语
	SLO   = "slo"   // 斯洛文尼亚语
	UK    = "uk"    // 乌克兰语
	MR    = "mr"    // 马拉提语
	AF    = "af"    // 南非语
	IS    = "is"    // 冰岛语
	IR    = "ir"    // 伊朗语
	AM    = "am"    // 阿姆哈拉语
	IT    = "it"    // 意大利语
	IW    = "iw"    // 希伯来语
	AS    = "as"    // 阿萨姆语
	SPA   = "spa"   // 西班牙语
	AZ    = "az"    // 阿塞拜疆语
	ID    = "id"    // 印尼语
	IG    = "ig"    // 伊博语
	NL    = "nl"    // 荷兰语
	PT_BR = "pt_BR" // 巴西语
	NO    = "no"    // 挪威语
	NE    = "ne"    // 尼泊尔语
	FA    = "fa"    // 波斯语
	WYW   = "wyw"   // 文言文
	KA    = "ka"    // 格鲁吉亚语
	KK    = "kk"    // 哈萨克语
	SR    = "sr"    // 塞尔维亚语
	SQ    = "sq"    // 阿尔巴尼亚语
	SW    = "sw"    // 斯瓦希里语
	KN    = "kn"    // 卡纳达语
	SK    = "sk"    // 斯洛伐克语
	SI    = "si"    // 僧加罗语
	KY    = "ky"    // 吉尔吉斯语
)

type Result struct {
	Errno       int    `json:"errno"`
	ErrMsg      string `json:"errmsg"`
	TransResult struct {
		Data []struct {
			Dst string `json:"dst"`
			Src string `json:"src"`
		} `json:"data"`
		From   string `json:"from"`
		To     string `json:"to"`
		Status int    `json:"status"`
		Type   int    `json:"type"`
	} `json:"trans_result"`
	LogID int `json:"logid"`
}

type Options struct {
	From  string
	To    string
	Proxy string
	Gtk   string
}

func NewOptions(from, to, tk, proxy string) *Options {
	return &Options{
		From:  from,
		To:    to,
		Proxy: proxy,
		Gtk:   tk,
	}
}

func Do(str string, options *Options) (result Result, err error) {
	if options == nil {
		err = fmt.Errorf("options is nil")
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				if options.Proxy == "" {
					return nil, nil
				}
				return url.ParseRequestURI(options.Proxy)
			},
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
	}

	var s string
	if s, err = sing(str, options.Gtk); err != nil {
		return
	}

	body := strings.NewReader(url.Values{
		"query":             []string{str},
		"from":              []string{options.From},
		"to":                []string{options.To},
		"transtype":         []string{"translang"},
		"simple_means_flag": []string{"3"},
		"sign":              []string{s},
		"token":             []string{"b648679daa321e869ce6e39d16536011"},
		"domain":            []string{"common"},
	}.Encode())

	var req *http.Request
	if req, err = http.NewRequest("POST", api+fmt.Sprintf("?from=%s&to=%s", options.From, options.To), body); err != nil {
		return
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", referer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", `PSTM=1599186903; BIDUPSID=31395A560050EB7379D0CD59058FBBF2; H_WISE_SIDS=154034_148078_153899_155177_155440_149356_150686_152055_150075_156287_150084_148867_155683_156095_154605_153715_153444_152410_154796_148303_146573_127969_154412_154175_152982_155906_150345_155318_146732_152571_131423_114551_153292_154189_155344_154952_144966_155813_154118_154801_154903_155238_152310_154301_155865_110085; BDUSS=1BwdktjeG1va243UFJyYUdOY3Zzb3NQYVBNeDhsRjVVYnBzZ3VHckxGS2hCNGhmRVFBQUFBJCQAAAAAAAAAAAEAAADEm3FzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKF6YF-hemBfY; BDUSS_BFESS=1BwdktjeG1va243UFJyYUdOY3Zzb3NQYVBNeDhsRjVVYnBzZ3VHckxGS2hCNGhmRVFBQUFBJCQAAAAAAAAAAAEAAADEm3FzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKF6YF-hemBfY; BDSFRCVID=uCKOJexroG3SoLrr7foXJTD47eKK0gOTDYLEOwXPsp3LGJLVN4vPEG0Ptf8gjM--J8jwogKK0gOTH6KF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF=tbkD_C-MfIvDqTrP-trf5DCShUFs-xJiB2Q-XPoO3KO6DpR4K6Jthpkp3-8Jt-QiWbRM2MbgylRp8P3y0bb2DUA1y4vpWj3qLgTxoUJ25DnJjlCzqfCWMR-ebPRiJPb9Qg-qahQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0HPonHjLMDTQ33f; BDSFRCVID_BFESS=uCKOJexroG3SoLrr7foXJTD47eKK0gOTDYLEOwXPsp3LGJLVN4vPEG0Ptf8gjM--J8jwogKK0gOTH6KF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF_BFESS=tbkD_C-MfIvDqTrP-trf5DCShUFs-xJiB2Q-XPoO3KO6DpR4K6Jthpkp3-8Jt-QiWbRM2MbgylRp8P3y0bb2DUA1y4vpWj3qLgTxoUJ25DnJjlCzqfCWMR-ebPRiJPb9Qg-qahQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0HPonHjLMDTQ33f; __yjs_duid=1_f3edbd8fa57282507f759741436a968c1618209277624; MCITY=-:; BAIDUID_BFESS=3BFED7E9104BF32F1EA249AFB5B73A79:FG=1; delPer=0; ZD_ENTRY=google; BDRCVFR[feWj1Vr5u3D]=I67x6TjHwwYf0; PSINO=3; H_PS_PSSID=34438_34442_31253_33848_34449_34092_34106_26350_34421_34323; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1630228270,1630229337; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1630229337; __yjs_st=2_NTgyZTk4OTY1NDUyZTQ5NWQ4NGUwOWU3MjY1YzNlZWQyYmMyYmZlMTMwOGMyOTEyNjBkZGQyM2MwYWE4OTE3OTNmMDVmNWEwZWFjM2E4NjFkZTU3OTkzZmQ0YWUzNDIzYjgyNGJiOTcxOTE1NWY3ZTBlNmE4MjU2ZDQwZTJkMzkxYTBiNGNkZjg1MDQ2YWE2ZGMzZGEyNjA1ZThjN2U3MDQ3YjJiOTc4NjI1MDkzZGI3NTc0NzgxNGNjOTZmYTRjZWQwNjI5YmExYWVhZGQyNDAzMTRkNTRjNjg5MjRkMTY3MTc1NTc5OGFlZTVmNDcwOTIwZjI0MDFmOGI0MGY4NF83XzIxODQ3NGZk; ab_sr=1.0.1_ZjdlZmNhMjNlNjk3OTVkMTcyMmYyNzljZjBiYmI4MmZlY2Q5ZTc0YWUwODRjZjJhYjk3MDQ1ZmIxZWQ5MTg2M2IxMTliZjM0NGU0NjI5NWYwYTZjZDRiODY5ZDZhYWU5MjVkNGU3YTliMGU3YzczY2E1NDU3ZjY0ODViMmFkZDk5ODRiZGFhZWNkODI2MzEwYzZmYzE4NGUxYTQ2NjhiYTVkZTcxZDA3ZGI0OTAyMWQyOTliYWMxNDg3N2JhNTg4`)
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return
	}

	var bytes []byte
	if bytes, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}

	if res.StatusCode >= 300 {
		err = fmt.Errorf(string(bytes))
		return
	}

	if err = json.Unmarshal(bytes, &result); err != nil {
		fmt.Println(string(bytes))
		return
	}

	if result.Errno != 0 {
		err = fmt.Errorf(result.ErrMsg)
		return
	}

	return
}

func Gtk() (tk string, err error) {
	var resp *http.Response
	client := http.Client{}

	req, _ := http.NewRequest("GET", referer, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", referer)

	if resp, err = client.Do(req); err != nil {
		return
	}

	if resp.StatusCode >= 300 {
		err = fmt.Errorf(fmt.Sprintf("http response code %d", resp.StatusCode))
		return
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	regex := regexp.MustCompile(`;window.gtk = '(\d+\.\d+?)';`)
	re := regex.FindAllStringSubmatch(string(body), -1)
	if len(re) < 1 && len(re[0]) < 2 {
		err = fmt.Errorf("get gtk error")
		return
	}

	tk = re[0][1]
	return
}

func sing(str, uid string) (string, error) {
	ctx := otto.New()
	_, err := ctx.Run(`function n(r, o) {
    for (var t = 0; t < o.length - 2; t += 3) {
        var a = o.charAt(t + 2);
        a = a >= "a" ? a.charCodeAt(0) - 87 : Number(a),
            a = "+" === o.charAt(t + 1) ? r >>> a : r << a,
            r = "+" === o.charAt(t) ? r + a & 4294967295 : r ^ a
    }
    return r
}

function sign(r, u) {
	var o = r.match(/[\uD800-\uDBFF][\uDC00-\uDFFF]/g);
	var t = r.length;
	t > 30 && (r = "" + r.substr(0, 10) + r.substr(Math.floor(t / 2) - 5, 10) + r.substr(-10, 10))

    for (var d = u.split("."), m = Number(d[0]) || 0, s = Number(d[1]) || 0, S = [], c = 0, v = 0; v < r.length; v++) {
        var A = r.charCodeAt(v);
        128 > A ? S[c++] = A : (2048 > A ? S[c++] = A >> 6 | 192 : (55296 === (64512 & A) && v + 1 < r.length && 56320 === (64512 & r.charCodeAt(v + 1)) ? (A = 65536 + ((1023 & A) << 10) + (1023 & r.charCodeAt(++v)),
            S[c++] = A >> 18 | 240,
            S[c++] = A >> 12 & 63 | 128) : S[c++] = A >> 12 | 224,
            S[c++] = A >> 6 & 63 | 128),
            S[c++] = 63 & A | 128)
    }
    for (var p = m, F = "" + String.fromCharCode(43) + String.fromCharCode(45) + String.fromCharCode(97) + ("" + String.fromCharCode(94) + String.fromCharCode(43) + String.fromCharCode(54)), D = "" + String.fromCharCode(43) + String.fromCharCode(45) + String.fromCharCode(51) + ("" + String.fromCharCode(94) + String.fromCharCode(43) + String.fromCharCode(98)) + ("" + String.fromCharCode(43) + String.fromCharCode(45) + String.fromCharCode(102)), b = 0; b < S.length; b++)
        p += S[b], p = n(p, F);
    return p = n(p, D), p ^= s, 0 > p && (p = (2147483647 & p) + 2147483648), p %= 1e6, p.toString() + "." + (p ^ m)
}`)
	if err != nil {
		return "", err
	}

	if value, err := ctx.Call(`sign`, nil, str, uid); err != nil {
		return "", err
	} else {
		return value.String(), nil
	}
}
