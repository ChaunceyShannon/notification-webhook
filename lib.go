package main

import (
	"github.com/alecthomas/chroma"
	htmlChromaFormatter "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	
	"github.com/corpix/uarand"
	"github.com/imroc/req"
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unicode/utf8"
	"github.com/denisbrodbeck/striphtmltags"
	"github.com/miekg/dns"
	"github.com/dustin/go-humanize"
	"github.com/robertkrimen/otto"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
	"github.com/xtaci/smux"
	"github.com/Go-ini/ini"
	"encoding/json"
	"github.com/bregydoc/gtranslate"
	"golang.org/x/text/language"
	wwwhoisgo "github.com/likexian/whois"
	"encoding/csv"
	"hash/crc32"
	"net"
	"bufio"
	"github.com/rakyll/statik/fs"
	"github.com/wangtuanjie/ip17mon"
	"sync"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"context"
	"github.com/cloudflare/cloudflare-go"
	"crypto/sha256"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"github.com/abadojack/whatlanggo"
	"crypto/aes"
	"crypto/cipher"
	"github.com/elliotchance/orderedmap"
	"syreclabs.com/go/faker"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"html/template"
	"encoding/base64"
	"path"
	"regexp"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/lithammer/shortuuid"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/fsnotify/fsnotify"
	"github.com/wunderlist/ttlcache"
	"github.com/antchfx/jsonquery"
	"github.com/zcalusic/sysinfo"
	ua "github.com/mileusna/useragent"
	"github.com/mileusna/crontab"
	"github.com/hpcloud/tail"
	"gitlab.com/toby3d/telegraph"
	"github.com/go-redis/redis/v8"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
	pquernatotp "github.com/pquerna/otp/totp"
	"compress/zlib"
	"github.com/ulikunitz/xz"
	"github.com/fatih/color"
	"github.com/bwmarrin/snowflake"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/mozillazg/go-pinyin"
	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart"
	"runtime/debug"
	"github.com/matrix-org/gomatrix"
	"github.com/mmcdole/gofeed"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/h2non/filetype"
	"github.com/alexeyco/simpletable"
	"github.com/bndr/gotabulate"
	"github.com/streadway/amqp"
)

func md2html(md string) string {
	flags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_TOC

	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	htmlContent := string(blackfriday.MarkdownOptions([]byte(md), blackfriday.HtmlRenderer(flags, "", ""), blackfriday.Options{Extensions: extensions}))

	insideCode := false
	htmlCodeContent := ""
	lang := ""
	code := ""
	for _, line := range strSplitlines(htmlContent) {
		if strStartsWith(line, "<pre><code") {
			insideCode = true
			if !strIn("<pre><code>", line) {
				res := reFindAll("<pre><code class=\"language-(.+?)\">", line)
				if len(res) != 0 {
					line = strReplace(line, res[0][0], "")
					lang = res[0][1]
				}
			} else {
				line = strReplace(line, "<pre><code>", "")
			}
		}
		if strStartsWith(line, "</code></pre>") {
			insideCode = false
			codeHTML := getHightLightHTML(htmldecode(code), lang)
			htmlCodeContent += codeHTML
			code = ""
			lang = ""
		}
		if insideCode == false && !strStartsWith(line, "</code></pre>") {
			// res := reFindAll("\\[(.+?)\\]\\((.+?)\\)", line)
			// if len(res) != 0 && !strIn("!"+res[0][0], line) {
			// 	links = append(links, res[0])
			// 	line = strReplace(line, res[0][0], "<a href=\""+res[0][2]+"\">"+res[0][1]+"</a>")
			// }
			htmlCodeContent += line + "\n"
		}
		if insideCode == true {
			code += line + "\n"
		}
	}

	return htmlCodeContent
}

func getHightLightHTML(code string, codeType ...string) string {
	style := "emacs"

	var l chroma.Lexer
	var displayLineNumber bool
	if len(codeType) != 0 && codeType[0] != "" {
		l = lexers.Get(codeType[0])
		if l == nil {
			l = lexers.Fallback
		} else {
			displayLineNumber = true
		}
	} else {
		l = lexers.Fallback
	}

	l = chroma.Coalesce(l)

	f := htmlChromaFormatter.New(htmlChromaFormatter.Standalone(false), htmlChromaFormatter.WithLineNumbers(displayLineNumber), htmlChromaFormatter.TabWidth(2))
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}
	it, err := l.Tokenise(nil, code)
	panicerr(err)

	buf := getBuffer()
	err = f.Format(&buf, s, it)
	panicerr(err)
	return buf.String()
}
type httpHeader map[string]string
type httpParam map[string]interface{}

type httpResp struct {
	headers    map[string]string
	content    string
	statusCode int
	url        string
}

type httpConfig struct {
	timeout             int
	readBodySize        int
	doNotFollowRedirect bool
	httpProxy           string
	timeoutRetryTimes   int
}

func httpHead(uri string, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Head(uri, header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}

func httpPostFile(uri string, filePath string, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Post(uri, req.File(filePath), header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}

func httpPostRaw(uri string, body string, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Post(uri, body, header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}

func httpPostJSON(uri string, json interface{}, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Post(uri, req.BodyJSON(&json), header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}

// httpPost(url, httpHeader{}, httpParam{}) {
func httpPost(uri string, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Post(uri, header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}

// httpGet(url, httpHeader{}, httpParam{}) {
func httpGet(uri string, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}

	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Get(uri, header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}

func httpPutJSON(uri string, json interface{}, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Put(uri, req.BodyJSON(&json), header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}

// httpPost(url, httpHeader{}, httpParam{}) {
func httpPut(uri string, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Put(uri, header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}

func httpPutRaw(uri string, body string, args ...interface{}) httpResp {
	if !strStartsWith(uri, "http://") && !strStartsWith(uri, "https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case httpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case httpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case httpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Put(uri, body, header, param)
		if err != nil {
			// lg.trace(err)
			if strIn("context deadline exceeded", err.Error()) {
				timeouttimes += 1
				if timeouttimes >= timeoutRetryTimes {
					panicerr(err)
				} else {
					continue
				}
			} else {
				panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = strJoin(" ", v)
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if strIn("context deadline exceeded", err.Error()) {
					timeouttimes += 1
					if timeouttimes >= timeoutRetryTimes {
						panicerr(err)
					} else {
						continue
					}
				} else {
					if !strIn("unexpected EOF", err.Error()) {
						panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		content:    string(respBody),
		headers:    headers,
		statusCode: hresp.StatusCode,
		url:        hresp.Request.URL.String(),
	}
}
// 返回interface的类型
func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// 返回当前时间戳
// time.time()
func now() int64 {
	return time.Now().Unix()
}

// 计算字符串的 MD5 散列值
func md5sum(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

//  计算指定文件的 MD5 散列值
func md5File(path string) string {
	data, err := ioutil.ReadFile(path)
	panicerr(err)
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// 计算字符串的 sha1 散列值
func sha1sum(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

//  计算文件的 sha1 散列值
func sha1File(path string) string {
	data, err := ioutil.ReadFile(path)
	panicerr(err)
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func getBuffer() bytes.Buffer {
	var a bytes.Buffer
	return a
}

type execCmdStruct struct {
	buf string
	out chan string
}

func (m *execCmdStruct) Write(p []byte) (n int, err error) {
	for _, pp := range p {
		ps := string(pp)
		if ps == "\n" {
			try(func() {
				m.out <- m.buf
			})
			m.buf = ""
		} else {
			m.buf = m.buf + ps
		}
	}

	return len(p), nil
}

func tailCmdOutput(command string) chan string {
	w := execCmdStruct{out: make(chan string, 9999)}
	go func() {
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = &w
		cmd.Stderr = &w
		err := cmd.Start()
		if err != nil {
			close(w.out)
			panicerr(err)
		}
		err = cmd.Wait()
		close(w.out)
		if err != nil {
			panicerr(err)
		}
	}()
	return w.out
}

func getTimeDuration(t interface{}) time.Duration {
	var timeDuration time.Duration
	if typeof(t) == "float64" {
		tt := t.(float64) * 1000
		timeDuration = time.Duration(tt) * time.Millisecond
	}
	if typeof(t) == "int" || typeof(t) == "int8" || typeof(t) == "int16" || typeof(t) == "int32" || typeof(t) == "int64" {
		tt := toInt64(t)
		timeDuration = time.Duration(tt) * time.Second
	}
	return timeDuration
}

var num2char = "0123456789abcdefghijklmnopqrstuvwxyz"

func numToBHex(num, n int) string {
	var numStr string
	for num != 0 {
		yu := num % n
		numStr = string(num2char[yu]) + numStr
		num = num / n
	}
	return numStr
}

func bhex2Num(str string, n int) int {
	str = strings.ToLower(str)
	v := 0.0
	length := len(str)
	for i := 0; i < length; i++ {
		s := string(str[i])
		index := strings.Index(num2char, s)
		v += float64(index) * math.Pow(float64(n), float64(length-1-i)) // 倒序
	}
	return int(v)
}

func hasChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

func plural(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		result = strconv.Itoa(count) + " " + singular + " "
	} else {
		result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

func fmtTimeDuration(second int64) (result string) {
	years := math.Floor(float64(second) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := second % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = second % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = second % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = second % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = second % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = second % 60

	if years > 0 {
		result = plural(int(years), "year") + plural(int(months), "month") + plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if months > 0 {
		result = plural(int(months), "month") + plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if weeks > 0 {
		result = plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if days > 0 {
		result = plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if hours > 0 {
		result = plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if minutes > 0 {
		result = plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else {
		result = plural(int(seconds), "second")
	}

	result = strStrip(result)

	return
}

func getGoroutineID() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}

func cmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func randomStr(length int, charset ...string) string {
	var str string
	if len(charset) != 0 {
		str = charset[0]
	} else {
		str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = str[seededRand.Intn(len(str))]
	}
	return string(b)
}

func strFilter(str string, charts ...string) (res string) {
	strb := []byte(str)
	var chartsb []byte
	if len(charts) != 0 {
		chartsb = []byte(charts[0])
	} else {
		chartsb = []byte("1234567890_qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM.")
	}
	for _, c := range strb {
		if itemInArray(c, chartsb) {
			res = res + string(c)
		}
	}
	return
}

func zipDir(dir, zipFile string) {

	fz, err := os.Create(zipFile)
	if err != nil {
		log.Fatalf("Create zip file failed: %s\n", err.Error())
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fDest, err := w.Create(path[len(dir)+1:])
			if err != nil {
				log.Printf("Create failed: %s\n", err.Error())
				return nil
			}
			fSrc, err := os.Open(path)
			if err != nil {
				log.Printf("Open failed: %s\n", err.Error())
				return nil
			}
			defer fSrc.Close()
			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				log.Printf("Copy failed: %s\n", err.Error())
				return nil
			}
		}
		return nil
	})
}

func unzipDir(zipFile, dir string) {

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatalf("Open zip file failed: %s\n", err.Error())
	}
	defer r.Close()

	for _, f := range r.File {
		func() {
			path := dir + string(filepath.Separator) + f.Name
			os.MkdirAll(filepath.Dir(path), 0755)
			fDest, err := os.Create(path)
			if err != nil {
				log.Printf("Create failed: %s\n", err.Error())
				return
			}
			defer fDest.Close()

			fSrc, err := f.Open()
			if err != nil {
				log.Printf("Open failed: %s\n", err.Error())
				return
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				log.Printf("Copy failed: %s\n", err.Error())
				return
			}
		}()
	}
}

func getSelfDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func getTempFilePath() string {
	file, err := ioutil.TempFile("", "systemd-private-")
	panicerr(err)
	defer os.Remove(file.Name())

	return file.Name()
}

func getTempDirPath() string {
	dir, err := ioutil.TempDir("", "systemd-private-")
	panicerr(err)
	return dir
}

func setFileTime(fpath string, mtime string, atime ...string) {
	var at time.Time
	if len(atime) == 0 {
		fi, err := os.Stat(fpath)
		panicerr(err)
		//mtime := fi.ModTime()
		stat := fi.Sys().(*syscall.Stat_t)
		at = time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
		// ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	} else {
		at = time.Unix(strptime("%Y-%m-%d %H:%M:%S", atime[0]), 0)
	}

	err := os.Chtimes(fpath, at, time.Unix(strptime("%Y-%m-%d %H:%M:%S", mtime), 0))
	panicerr(err)
}

type fileTimeStruct struct {
	ctime int64
	mtime int64
	atime int64
}

func getFileTime(fpath string) *fileTimeStruct {
	fi, err := os.Stat(fpath)
	panicerr(err)
	mtime := fi.ModTime().Unix()
	stat := fi.Sys().(*syscall.Stat_t)
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)).Unix()
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec)).Unix()

	return &fileTimeStruct{
		mtime: mtime,
		ctime: ctime,
		atime: atime,
	}
}

func map2bin(m map[string]string) (res string) {
	btlen := make([]byte, 4)
	binary.LittleEndian.PutUint32(btlen, uint32(len(m)))
	res += string(btlen) // 元素个数

	for k, v := range m {
		binary.LittleEndian.PutUint32(btlen, uint32(len(k)))
		res += string(btlen) // key长度
		res += k             // key

		binary.LittleEndian.PutUint32(btlen, uint32(len(v)))
		res += string(btlen) // value长度
		res += v             // value
	}
	return
}

func bin2map(s string) (res map[string]string) {
	res = make(map[string]string)

	btlen := make([]byte, 4)

	ss := bytes.NewBuffer([]byte(s))

	ss.Read(btlen)
	tlen := int(binary.LittleEndian.Uint32(btlen)) // 元素个数

	for range rangeInt(tlen) {
		ss.Read(btlen)
		tlen := int(binary.LittleEndian.Uint32(btlen)) // key长度
		key := make([]byte, tlen)
		ss.Read(key) // 读取key

		ss.Read(btlen)
		tlen = int(binary.LittleEndian.Uint32(btlen)) // value长度
		value := make([]byte, tlen)
		ss.Read(value) // 读取key

		res[string(key)] = string(value)
	}
	return
}

func getSSLExpiryTime(domain string, port ...int) int64 {
	var p string
	if len(port) == 0 {
		p = "443"
	} else {
		p = str(port[0])
	}

	conn, err := tls.Dial("tcp", domain+":"+p, nil)
	if err != nil {
		panic("Server doesn't support SSL certificate err: " + err.Error())
	}

	err = conn.VerifyHostname(domain)
	if err != nil {
		panic("Hostname doesn't match with certificate: " + err.Error())
	}

	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter.Unix()
	return expiry
}

func nslookup(name string, querytype string, dnsService ...string) [][]string {
	c := dns.Client{
		Timeout: 5 * time.Second,
	}

	var server string
	if len(dnsService) == 0 {
		server = "8.8.8.8"
	} else {
		server = dnsService[0]
	}

	querytype = strLower(querytype)
	var qtype uint16
	if querytype == "ns" {
		qtype = dns.TypeNS
	} else if querytype == "a" {
		qtype = dns.TypeA
	} else if querytype == "txt" {
		qtype = dns.TypeTXT
	} else if querytype == "cname" {
		qtype = dns.TypeCNAME
	} else if querytype == "aaaa" {
		qtype = dns.TypeAAAA
	} else if querytype == "soa" {
		qtype = dns.TypeSOA
	} else if querytype == "mx" {
		qtype = dns.TypeMX
	}

	var err error
	var dst [][]string
	for i := 0; i < 3; i++ {
		m := dns.Msg{}
		if !strEndsWith(name, ".") {
			name = name + "."
		}
		m.SetQuestion(name, qtype)
		r, _, err := c.Exchange(&m, server+":53")
		if err != nil {
			time.Sleep(1 * time.Second * time.Duration(i+1))
			continue
		}

		for _, ans := range r.Answer {
			s := strings.Split(ans.String(), "\t")
			dst = append(dst, []string{s[0], s[3], s[4]})
		}
		err = nil
		break
	}

	panicerr(err)

	return dst
}

func strRemoveHtmlTag(in string) string {
	return striphtmltags.StripTags(in)
}

func strRemoveNonUTF8Character(text string) string {
	if !utf8.ValidString(text) {
		v := make([]rune, 0, len(text))
		for i, r := range text {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(text[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		text = string(v)
	}
	return text
}
func fmtsize(num uint64) string {
	return humanize.Bytes(num)
}
type javascriptVMStruct struct {
	vm *otto.Otto
}

func getJavascriptVM() *javascriptVMStruct {
	return &javascriptVMStruct{vm: otto.New()}
}

func (m *javascriptVMStruct) run(javascript string) *javascriptVMStruct {
	_, err := m.vm.Run(javascript)
	panicerr(err)
	return m
}

func (m *javascriptVMStruct) get(variableName string) string {
	value, err := m.vm.Get(variableName)
	panicerr(err)
	valueStr, err := value.ToString()
	panicerr(err)
	if valueStr == "undefined" {
		panicerr("变量" + variableName + "未定义")
	}
	return valueStr
}

func (m *javascriptVMStruct) set(variableName string, variableValue interface{}) {
	err := m.vm.Set(variableName, variableValue)
	panicerr(err)
}

func (m *javascriptVMStruct) isdefined(variableName string) bool {
	value, err := m.vm.Get(variableName)
	panicerr(err)
	valueStr, err := value.ToString()
	panicerr(err)
	return valueStr != "undefined"
}
func resizeImg(srcPath string, dstPath string, width int, height ...int) {
	file, err := os.Open(srcPath)
	panicerr(err)

	var img image.Image
	ftype := fileType(srcPath)
	if ftype == "jpg" {
		img, err = jpeg.Decode(file)
		panicerr(err)
	} else if ftype == "png" {
		img, err = png.Decode(file)
		panicerr(err)
	} else if ftype == "webp" {
		img, err = webp.Decode(file)
		panicerr(err)
	} else {
		err = errors.New("只支持读取jpg、png或者webp格式图片")
		panicerr(err)
	}
	file.Close()

	var h int
	if len(height) != 0 {
		h = height[0]
	} else {
		h = 0
	}

	resizeWidth := 0
	resizeHeight := 0
	if img.Bounds().Max.X > width {
		resizeWidth = width
	} else {
		resizeWidth = img.Bounds().Max.X
	}

	if h == 0 || img.Bounds().Max.Y > h {
		resizeHeight = h
	} else {
		resizeHeight = img.Bounds().Max.Y
	}

	m := resize.Resize(toUint(resizeWidth), toUint(resizeHeight), img, resize.Lanczos3)

	out, err := os.Create(dstPath)
	panicerr(err)
	defer out.Close()

	if strEndsWith(dstPath, ".jpg") || strEndsWith(dstPath, ".jpeg") {
		err = jpeg.Encode(out, m, nil)
	} else if strEndsWith(dstPath, ".png") {
		err = png.Encode(out, m)
		// } else if strEndsWith(dstPath, ".webp") {
		// 	err = webp.Encode(out, m)
	} else {
		// err = errors.New("只支持导出jpg、png或者webp，以文件扩展名来判定。")
		err = errors.New("只支持导出jpg或者png，以文件扩展名来判定。")
		panicerr(err)
	}

	panicerr(err)
}
type smuxConfig struct {
	KeepAliveInterval int
	KeepAliveTimeout  int
}

type smuxServerSideListener struct {
	listener *smux.Session
	isclose  bool
}

func smuxServerWrapper(conn io.ReadWriteCloser, cfg ...smuxConfig) *smuxServerSideListener {
	scfg := smux.DefaultConfig()
	if len(cfg) != 0 {
		if cfg[0].KeepAliveInterval != 0 {
			scfg.KeepAliveInterval = time.Duration(cfg[0].KeepAliveInterval) * time.Second
		}
		if cfg[0].KeepAliveTimeout != 0 {
			scfg.KeepAliveTimeout = time.Duration(cfg[0].KeepAliveTimeout) * time.Second
		}
	} else {
		scfg.KeepAliveInterval = time.Duration(2) * time.Second
		scfg.KeepAliveTimeout = time.Duration(7) * time.Second
	}

	listener, err := smux.Server(conn, scfg)
	panicerr(err)
	m := &smuxServerSideListener{listener: listener}

	go func() {
		for {
			m.isclose = m.listener.IsClosed()
			if m.isclose {
				lg.trace("listener is closed")
				break
			}
			sleep(0.1)
		}
	}()

	return m
}

type smuxServerSideConnection struct {
	stream  *smux.Stream
	isclose bool
}

func (m *smuxServerSideListener) accept() chan *smuxServerSideConnection {
	ch := make(chan *smuxServerSideConnection)

	go func() {
		err := try(func() {
			for {
				stream, err := m.listener.AcceptStream()
				panicerr(err)

				m := &smuxServerSideConnection{stream: stream}

				go func() {
					<-m.stream.GetDieCh()
					lg.trace("Stream is closed")
					m.isclose = true
				}()

				ch <- m
			}
		})
		lg.trace("smux接收新连接的时候报错:", err, "session为:", m.listener)
		close(ch)
	}()
	return ch
}

func (m *smuxServerSideConnection) Send(data map[string]string, timeout ...int) {
	m.send(data, timeout...)
}

func (m *smuxServerSideConnection) send(data map[string]string, timeout ...int) {
	if len(timeout) != 0 {
		m.stream.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}

	text := map2bin(data)

	btlen := make([]byte, 4)
	binary.LittleEndian.PutUint32(btlen, uint32(len(text)))
	_, err := m.stream.Write(btlen)
	panicerr(err)

	_, err = m.stream.Write([]byte(text))
	panicerr(err)

	m.stream.SetWriteDeadline(time.Time{})
}

func (m *smuxServerSideConnection) recv(timeout ...int) (data map[string]string) {
	if len(timeout) != 0 {
		m.stream.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}

	headertotallen := 4

	totalblen := ""
	buf := make([]byte, headertotallen)
	for {
		n, err := m.stream.Read(buf)
		if err != nil {
			if err.Error() == "timeout" {
				return nil
			}
			panicerr(err)
		}

		totalblen = totalblen + string(buf[:n])

		if len(totalblen) != toInt(headertotallen) {
			buf = make([]byte, toInt(headertotallen)-len(totalblen))
		} else {
			break
		}
	}

	totallen := binary.LittleEndian.Uint32([]byte(totalblen))

	totaldata := ""
	buf = make([]byte, totallen)
	for {
		n, err := m.stream.Read(buf)
		panicerr(err)

		totaldata = totaldata + string(buf[:n])

		if len(totaldata) != toInt(totallen) {
			buf = make([]byte, toInt(totallen)-len(totaldata))
		} else {
			break
		}
	}

	data = bin2map(totaldata)

	m.stream.SetReadDeadline(time.Time{})

	return
}

func (m *smuxServerSideConnection) close() {
	if !m.isclose {
		m.stream.Close()
	}
}

type smuxClientSideSession struct {
	session *smux.Session
	isclose bool
}

func smuxClientWrapper(conn io.ReadWriteCloser, cfg ...smuxConfig) *smuxClientSideSession {
	scfg := smux.DefaultConfig()
	if len(cfg) != 0 {
		if cfg[0].KeepAliveInterval != 0 {
			scfg.KeepAliveInterval = time.Duration(cfg[0].KeepAliveInterval) * time.Second
		}
		if cfg[0].KeepAliveTimeout != 0 {
			scfg.KeepAliveTimeout = time.Duration(cfg[0].KeepAliveTimeout) * time.Second
		}
	} else {
		scfg.KeepAliveInterval = time.Duration(2) * time.Second
		scfg.KeepAliveTimeout = time.Duration(7) * time.Second
	}

	session, err := smux.Client(conn, scfg)
	panicerr(err)

	m := &smuxClientSideSession{session: session}

	go func() {
		for {
			m.isclose = m.session.IsClosed()
			if m.isclose {
				break
			}
			sleep(0.1)
		}
	}()

	return m
}

type smuxClientSideConnection struct {
	stream  *smux.Stream
	isclose bool
}

func (m *smuxClientSideSession) connect() *smuxClientSideConnection {
	stream, err := m.session.OpenStream()
	panicerr(err)

	mm := &smuxClientSideConnection{stream: stream}

	go func() {
		<-mm.stream.GetDieCh()
		mm.isclose = true
	}()

	return mm
}

func (m *smuxClientSideSession) close() {
	if !m.isclose {
		m.isclose = true
		try(func() {
			m.session.Close()
		})
	}
}

func (m *smuxClientSideConnection) send(data map[string]string, timeout ...int) {
	if len(timeout) != 0 {
		m.stream.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}

	text := map2bin(data)

	btlen := make([]byte, 4)
	binary.LittleEndian.PutUint32(btlen, uint32(len(text)))
	_, err := m.stream.Write(btlen)
	panicerr(err)

	_, err = m.stream.Write([]byte(text))
	panicerr(err)

	m.stream.SetWriteDeadline(time.Time{})
}

func (m *smuxClientSideConnection) recv(timeout ...int) (data map[string]string) {
	if len(timeout) != 0 {
		m.stream.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}

	headertotallen := 4

	totalblen := ""
	buf := make([]byte, headertotallen)
	for {
		n, err := m.stream.Read(buf)
		if err != nil {
			if err.Error() == "timeout" {
				return nil
			}
			panicerr(err)
		}

		totalblen = totalblen + string(buf[:n])

		if len(totalblen) != toInt(headertotallen) {
			buf = make([]byte, toInt(headertotallen)-len(totalblen))
		} else {
			break
		}
	}

	totallen := binary.LittleEndian.Uint32([]byte(totalblen))

	totaldata := ""
	buf = make([]byte, totallen)
	for {
		n, err := m.stream.Read(buf)
		panicerr(err)

		totaldata = totaldata + string(buf[:n])

		if len(totaldata) != toInt(totallen) {
			buf = make([]byte, toInt(totallen)-len(totaldata))
		} else {
			break
		}
	}

	data = bin2map(totaldata)

	m.stream.SetReadDeadline(time.Time{})

	return
}

func (m *smuxClientSideConnection) close() {
	if !m.isclose {
		m.stream.Close()
	}
}
type iniStruct struct {
	cfg     *ini.File
	cfgPath string
}

func getIni(fpath ...string) *iniStruct {
	var cfg *ini.File
	var err error
	var cfgPath string
	if len(fpath) != 0 {
		cfgPath = fpath[0]
		if pathExists(cfgPath) {
			cfg, err = ini.Load(fpath[0])
			panicerr(err)
		} else {
			cfg = ini.Empty()
		}
	} else {
		cfgPath = ""
		cfg = ini.Empty()
	}
	return &iniStruct{cfg: cfg, cfgPath: cfgPath}
}

func (m *iniStruct) get(SectionKeyDefaultComment ...string) (res string) {
	if len(SectionKeyDefaultComment) == 2 {
		res = m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).String()
	} else if len(SectionKeyDefaultComment) == 3 {
		res = m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).String()
		if res == "" {
			res = SectionKeyDefaultComment[2]
			m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).SetValue(SectionKeyDefaultComment[2])
		}
	} else if len(SectionKeyDefaultComment) == 4 {
		res = m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).String()
		if res == "" {
			res = SectionKeyDefaultComment[2]
			m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).SetValue(SectionKeyDefaultComment[2])
		}
		m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).Comment = SectionKeyDefaultComment[3]
	}
	return
}

func (m *iniStruct) getInt(key ...string) int {
	return toInt(m.get(key...))
}

func (m *iniStruct) getInt64(key ...string) int64 {
	return toInt64(m.get(key...))
}

func (m *iniStruct) getFloat64(key ...string) float64 {
	return toFloat64(m.get(key...))
}

func (m *iniStruct) set(SectionKeyValueComment ...string) {
	if len(SectionKeyValueComment) == 3 {
		m.cfg.Section(SectionKeyValueComment[0]).Key(SectionKeyValueComment[1]).SetValue(SectionKeyValueComment[2])
	} else if len(SectionKeyValueComment) == 4 {
		m.cfg.Section(SectionKeyValueComment[0]).Key(SectionKeyValueComment[1]).SetValue(SectionKeyValueComment[2])
		m.cfg.Section(SectionKeyValueComment[0]).Key(SectionKeyValueComment[1]).Comment = SectionKeyValueComment[3]
	} else {
		panicerr("按顺序指定section, key, value(以及comment)")
	}
}

func (m *iniStruct) save(fpath ...string) (exist bool) {
	if len(fpath) != 0 {
		exist = pathExists(fpath[0])
		m.cfg.SaveTo(fpath[0])
	} else if m.cfgPath != "" {
		exist = pathExists(m.cfgPath)
		m.cfg.SaveTo(m.cfgPath)
	}
	return
}
// 有一定的次数限制, 具体多少未知, 但很少, 到限制就hang住
func gen2zh(text string) string {
	s, err := gtranslate.Translate(text, language.English, language.Chinese)
	panicerr(err)
	return s
}

type bresult struct {
	ErrorCode   string `json:"error_code"`
	ErrorMsg    string `json:"error_msg"`
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func ben2zh(text string) string {
	appID := "20200219000386130"
	secretKey := "FYNyhQRLtgbnzcYqqioj"
	apiURL := "http://api.fanyi.baidu.com/api/trans/vip/translate"

	bres := bresult{}

	for {
		salt := str(randint(32768, 65536))
		pre_sign := appID + text + salt + secretKey
		sign := md5sum(pre_sign)
		params := httpParam{
			"q":     text,
			"from":  "auto",
			"to":    "zh",
			"appid": appID,
			"salt":  salt,
			"sign":  sign,
		}

		err := json.Unmarshal([]byte(httpGet(apiURL, params).content), &bres)
		if err == nil {
			break
		} else {
			sleep(3)
			continue
		}
	}

	if bres.ErrorCode != "" && len(bres.TransResult) == 0 {
		panicerr("翻译出错, 状态码 " + bres.ErrorCode + ", 原因为 " + bres.ErrorMsg)
	}
	return bres.TransResult[0].Dst
}

func en2zh(text string) (res string) {
	res = ben2zh(text)
	return
}
func whois(s string, servers ...string) string {
	result, err := wwwhoisgo.Whois(s, servers...)
	panicerr(err)
	return result
}

var domainSuffixWhoisServerMap map[string]string

func getDomainWhoisInfo(domain string) string {
	darr := strSplit(domain, ".")
	suffix := darr[len(darr)-1]
	var dws string

	if domainSuffixWhoisServerMap == nil {
		domainSuffixWhoisServerMap = make(map[string]string)
	}

	if keyInMap(suffix, domainSuffixWhoisServerMap) {
		dws = domainSuffixWhoisServerMap[suffix]
	} else {
		//lg.trace("尝试查找顶级域" + suffix + "的whois服务器")
		s := reFindAll("whois:(.+)", whois(suffix, "whois.iana.org"))
		if len(s) == 0 {
			panicerr("找不到顶级域" + suffix + "的whois服务器")
		}
		dws = strStrip(s[0][1])

		domainSuffixWhoisServerMap[suffix] = dws
	}

	return whois(domain, dws)
}
// 检查是否是合法的年月日
func checkDate(month, day, year int) bool {
	if month < 1 || month > 12 || day < 1 || day > 31 || year < 1 || year > 32767 {
		return false
	}
	switch month {
	case 4, 6, 9, 11:
		if day > 30 {
			return false
		}
	case 2:
		// leap year
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			if day > 29 {
				return false
			}
		} else if day > 28 {
			return false
		}
	}

	return true
}

// 将字符串的首字母转换为大写
func ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

//  使一个字符串的第一个字符小写
func lcfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

//  将字符串中每个单词的首字母转换为大写
func ucwords(str string) string {
	return strings.Title(str)
}

//  返回字符串的子串
// func substr(str string, start uint, length int) string {
// 	if start < 0 || length < -1 {
// 		return str
// 	}
// 	switch {
// 	case length == -1:
// 		return str[start:]
// 	case length == 0:
// 		return ""
// 	}
// 	end := int(start) + length
// 	if end > len(str) {
// 		end = len(str)
// 	}
// 	return str[start:end]
// }

func substr(s string, start int, end int) string {
	start_str_idx := 0
	i := 0
	for j := range s {
		if i == start {
			start_str_idx = j
		}
		if i == end {
			return s[start_str_idx:j]
		}
		i++
	}
	return s[start_str_idx:]
}

// s = [i for i in 'abc']
// s.reverse()
// ''.join(s)
func strrev(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 处理url的参数字段
// f1=m&f2=n -> map[f1:m f2:n]
// f[a]=m&f[b]=n -> map[f:map[a:m b:n]]
// f[a][a]=m&f[a][b]=n -> map[f:map[a:map[a:m b:n]]]
// f[]=m&f[]=n -> map[f:[m n]]
// f[a][]=m&f[a][]=n -> map[f:map[a:[m n]]]
// f[][]=m&f[][]=n -> map[f:[map[]]] // Currently does not support nested slice.
// f=m&f[a]=n -> error // This is not the same as PHP.
// a .[[b=c -> map[a___[b:c]
//  将字符串解析成多个变量
func parseStr(encodedString string, result map[string]interface{}) error {
	// build nested map.
	var build func(map[string]interface{}, []string, interface{}) error

	build = func(result map[string]interface{}, keys []string, value interface{}) error {
		length := len(keys)
		// trim ',"
		key := strings.Trim(keys[0], "'\"")
		if length == 1 {
			result[key] = value
			return nil
		}

		// The end is slice. like f[], f[a][]
		if keys[1] == "" && length == 2 {
			// todo nested slice
			if key == "" {
				return nil
			}
			val, ok := result[key]
			if !ok {
				result[key] = []interface{}{value}
				return nil
			}
			children, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
			}
			result[key] = append(children, value)
			return nil
		}

		// The end is slice + map. like f[][a]
		if keys[1] == "" && length > 2 && keys[2] != "" {
			val, ok := result[key]
			if !ok {
				result[key] = []interface{}{}
				val = result[key]
			}
			children, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
			}
			if l := len(children); l > 0 {
				if child, ok := children[l-1].(map[string]interface{}); ok {
					if _, ok := child[keys[2]]; !ok {
						_ = build(child, keys[2:], value)
						return nil
					}
				}
			}
			child := map[string]interface{}{}
			_ = build(child, keys[2:], value)
			result[key] = append(children, child)

			return nil
		}

		// map. like f[a], f[a][b]
		val, ok := result[key]
		if !ok {
			result[key] = map[string]interface{}{}
			val = result[key]
		}
		children, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected type 'map[string]interface{}' for key '%s', but got '%T'", key, val)
		}

		return build(children, keys[1:], value)
	}

	// split encodedString.
	parts := strings.Split(encodedString, "&")
	for _, part := range parts {
		pos := strings.Index(part, "=")
		if pos <= 0 {
			continue
		}
		key, err := url.QueryUnescape(part[:pos])
		if err != nil {
			return err
		}
		for key[0] == ' ' {
			key = key[1:]
		}
		if key == "" || key[0] == '[' {
			continue
		}
		value, err := url.QueryUnescape(part[pos+1:])
		if err != nil {
			return err
		}

		// split into multiple keys
		var keys []string
		left := 0
		for i, k := range key {
			if k == '[' && left == 0 {
				left = i
			} else if k == ']' {
				if left > 0 {
					if len(keys) == 0 {
						keys = append(keys, key[:left])
					}
					keys = append(keys, key[left+1:i])
					left = 0
					if i+1 < len(key) && key[i+1] != '[' {
						break
					}
				}
			}
		}
		if len(keys) == 0 {
			keys = append(keys, key)
		}
		// first key
		first := ""
		for i, chr := range keys[0] {
			if chr == ' ' || chr == '.' || chr == '[' {
				first += "_"
			} else {
				first += string(chr)
			}
			if chr == '[' {
				first += keys[0][i+1:]
				break
			}
		}
		keys[0] = first

		// build nested map
		if err := build(result, keys, value); err != nil {
			return err
		}
	}

	return nil
}

// numberFormat(123456789.12933, 2, ".", ",") // 123,456,789.13
// numberFormat(123456789.12933, 5, "=", "-") // 123-456-789=12933
// numberFormat(123456789, 0, ".", ",") // 123,456,789
//  以千位分隔符方式格式化一个数字
func numberFormat(number float64, decimals uint, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	dec := int(decimals)
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(dec)+"F", number)
	prefix, suffix := "", ""
	if dec > 0 {
		prefix = str[:len(str)-(dec+1)]
		suffix = str[len(str)-dec:]
	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	// thousands sep num
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]
	}
	s := string(tmp)
	if dec > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}

// chunkSplit("1234567890", 3, ">") // 123>456>789>0>
//  将字符串分割成小块
func chunkSplit(body string, chunklen uint, end string) string {
	if end == "" {
		end = "\r\n"
	}
	runes, erunes := []rune(body), []rune(end)
	l := uint(len(runes))
	if l <= 1 || l < chunklen {
		return body + end
	}
	ns := make([]rune, 0, len(runes)+len(erunes))
	var i uint
	for i = 0; i < l; i += chunklen {
		if i+chunklen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunklen]...)
		}
		ns = append(ns, erunes...)
	}
	return string(ns)
}

// <?php
// $str = "An example of a long word is: Supercalifragulistic";
// echo wordwrap($str,15,"<br>\n");
// ?>
// 按照指定长度对字符串进行折行处理
func wordwrap(str string, width uint, br string, cut bool) string {
	strlen := len(str)
	brlen := len(br)
	linelen := int(width)

	if strlen == 0 {
		return ""
	}
	if brlen == 0 {
		panic("break string cannot be empty")
	}
	if linelen == 0 && cut {
		panic("can't force cut when width is zero")
	}

	current, laststart, lastspace := 0, 0, 0
	var ns []byte
	for current = 0; current < strlen; current++ {
		if str[current] == br[0] && current+brlen < strlen && str[current:current+brlen] == br {
			ns = append(ns, str[laststart:current+brlen]...)
			current += brlen - 1
			lastspace = current + 1
			laststart = lastspace
		} else if str[current] == ' ' {
			if current-laststart >= linelen {
				ns = append(ns, str[laststart:current]...)
				ns = append(ns, br[:]...)
				laststart = current + 1
			}
			lastspace = current
		} else if current-laststart >= linelen && cut && laststart >= lastspace {
			ns = append(ns, str[laststart:current]...)
			ns = append(ns, br[:]...)
			laststart = current
			lastspace = current
		} else if current-laststart >= linelen && laststart < lastspace {
			ns = append(ns, str[laststart:lastspace]...)
			ns = append(ns, br[:]...)
			lastspace++
			laststart = lastspace
		}
	}

	if laststart != current {
		ns = append(ns, str[laststart:current]...)
	}
	return string(ns)
}

// Strlen strlen()
//  获取字符串长度
func strlen(str string) int {
	return len(str)
}

// MbStrlen mb_strlen()
//  获取字符串的长度
func mbStrlen(str string) int {
	return utf8.RuneCountInString(str)
}

// <?php
// echo str_repeat("-=", 10);
// ?>
// 以上例程会输出：
// -=-=-=-=-=-=-=-=-=-=
// 重复一个字符串
func strRepeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// 返回 haystack 字符串从 needle 第一次出现的位置开始到 haystack 结尾的字符串。
// 查找字符串的首次出现
func strstr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}

// If the parameter length is 1, type is: map[string]string
// Strtr("baab", map[string]string{"ab": "01"}) will return "ba01"
// If the parameter length is 2, type is: string, string
// Strtr("baab", "ab", "01") will return "1001", a => 0; b => 1.
// 转换指定字符
func strtr(haystack string, params ...interface{}) string {
	ac := len(params)
	if ac == 1 {
		pairs := params[0].(map[string]string)
		length := len(pairs)
		if length == 0 {
			return haystack
		}
		oldnew := make([]string, length*2)
		for o, n := range pairs {
			if o == "" {
				return haystack
			}
			oldnew = append(oldnew, o, n)
		}
		return strings.NewReplacer(oldnew...).Replace(haystack)
	} else if ac == 2 {
		from := params[0].(string)
		to := params[1].(string)
		trlen, lt := len(from), len(to)
		if trlen > lt {
			trlen = lt
		}
		if trlen == 0 {
			return haystack
		}

		str := make([]uint8, len(haystack))
		var xlat [256]uint8
		var i int
		var j uint8
		if trlen == 1 {
			for i = 0; i < len(haystack); i++ {
				if haystack[i] == from[0] {
					str[i] = to[0]
				} else {
					str[i] = haystack[i]
				}
			}
			return string(str)
		}
		// trlen != 1
		for {
			xlat[j] = j
			if j++; j == 0 {
				break
			}
		}
		for i = 0; i < trlen; i++ {
			xlat[from[i]] = to[i]
		}
		for i = 0; i < len(haystack); i++ {
			str[i] = xlat[haystack[i]]
		}
		return string(str)
	}

	return haystack
}

//  随机打乱一个字符串
func strShuffle(str string) string {
	runes := []rune(str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := make([]rune, len(runes))
	for i, v := range r.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

//  使用一个字符串分割另一个字符串
func explode(delimiter, str string) []string {
	return strings.Split(str, delimiter)
}

// 返回指定的字符
func chr(ascii int) string {
	return string(ascii)
}

//  转换字符串第一个字节为 0-255 之间的值
func ord(char string) int {
	r, _ := utf8.DecodeRune([]byte(char))
	return int(r)
}

//  在字符串所有新行之前插入 HTML 换行标记
func nl2br(str string, isXhtml bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if isXhtml {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}
	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}
		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

//  使用反斜线引用字符串
func addslashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

///  反引用一个引用字符串. 返回一个去除转义反斜线后的字符串（\' 转换为 ' 等等）。双反斜线（\\）被转换为单个反斜线（\）。
func stripslashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// 返回 在下面这些特殊字符前加 反斜线(\) 转义后的字符串。 这些特殊字符包含：
// . \ + * ? [ ^ ] ( $ )
//  转义元字符集
func quotemeta(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

//  计算一个字符串的 crc32 多项式
func crc32sum(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// Levenshtein levenshtein()
// costIns: Defines the cost of insertion.
// costRep: Defines the cost of replacement.
// costDel: Defines the cost of deletion.
// 此函数返回两个字符串参数之间的编辑距离，如果其中一个字符串参数长度大于限制的255个字符时，返回-1。
//  计算两个字符串之间的编辑距离
func levenshtein(str1, str2 string, costIns, costRep, costDel int) int {
	var maxLen = 255
	l1 := len(str1)
	l2 := len(str2)
	if l1 == 0 {
		return l2 * costIns
	}
	if l2 == 0 {
		return l1 * costDel
	}
	if l1 > maxLen || l2 > maxLen {
		return -1
	}

	p1 := make([]int, l2+1)
	p2 := make([]int, l2+1)
	var c0, c1, c2 int
	var i1, i2 int
	for i2 := 0; i2 <= l2; i2++ {
		p1[i2] = i2 * costIns
	}
	for i1 = 0; i1 < l1; i1++ {
		p2[0] = p1[0] + costDel
		for i2 = 0; i2 < l2; i2++ {
			if str1[i1] == str2[i2] {
				c0 = p1[i2]
			} else {
				c0 = p1[i2] + costRep
			}
			c1 = p1[i2+1] + costDel
			if c1 < c0 {
				c0 = c1
			}
			c2 = p2[i2] + costIns
			if c2 < c0 {
				c0 = c2
			}
			p2[i2+1] = c0
		}
		tmp := p1
		p1 = p2
		p2 = tmp
	}
	c0 = p1[l2]

	return c0
}

// SimilarText similar_text()
// 返回在两个字符串中匹配字符的数目。
// 计算两个字符串的相似度
func SimilarText(first, second string, percent *float64) int {
	var similarText func(string, string, int, int) int
	similarText = func(str1, str2 string, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		// Find the longest segment of the same section in two strings
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for l := 0; (i+l < len1) && (j+l < len2) && (str1[i+l] == str2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(str1, str2, pos1, pos2)
			}
			if (pos1+max < len1) && (pos2+max < len2) {
				s1 := []byte(str1)
				s2 := []byte(str2)
				sum += similarText(string(s1[pos1+max:]), string(s2[pos2+max:]), len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}
	sim := similarText(first, second, l1, l2)
	if percent != nil {
		*percent = float64(sim*200) / float64(l1+l2)
	}
	return sim
}

// Soundex soundex()
// Calculate the soundex key of a string.
func soundex(str string) string {
	if str == "" {
		panic("str: cannot be an empty string")
	}
	table := [26]rune{
		// A, B, C, D
		'0', '1', '2', '3',
		// E, F, G
		'0', '1', '2',
		// H
		'0',
		// I, J, K, L, M, N
		'0', '2', '2', '4', '5', '5',
		// O, P, Q, R, S, T
		'0', '1', '2', '6', '2', '3',
		// U, V
		'0', '1',
		// W, X
		'0', '2',
		// Y, Z
		'0', '2',
	}
	last, code, small := -1, 0, 0
	sd := make([]rune, 4)
	// build soundex string
	for i := 0; i < len(str) && small < 4; i++ {
		// ToUpper
		char := str[i]
		if char < '\u007F' && 'a' <= char && char <= 'z' {
			code = int(char - 'a' + 'A')
		} else {
			code = int(char)
		}
		if code >= 'A' && code <= 'Z' {
			if small == 0 {
				sd[small] = rune(code)
				small++
				last = int(table[code-'A'])
			} else {
				code = int(table[code-'A'])
				if code != last {
					if code != 0 {
						sd[small] = rune(code)
						small++
					}
					last = code
				}
			}
		}
	}
	// pad with "0"
	for ; small < 4; small++ {
		sd[small] = '0'
	}
	return string(sd)
}

//////////// URL Functions ////////////

// Rawurlencode rawurlencode()
func rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// Rawurldecode rawurldecode()
//  按照 RFC 3986 对 URL 进行编码
func rawurldecode(str string) string {
	str, err := url.QueryUnescape(strings.Replace(str, "%20", "+", -1))
	panicerr(err)
	return str
}

// HTTPBuildQuery http_build_query()
//  生成 URL-encode 之后的请求字符串
func httpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

//////////// Array(Slice/Map) Functions ////////////

// ArrayFill array_fill()
// <?php
// $keys = array('foo', 5, 10, 'bar');
// $a = array_fill_keys($keys, 'banana');
// print_r($a);
// ?>
//
// 以上例程会输出：
//
// Array
// (
//     [foo] => banana
//     [5] => banana
//     [10] => banana
//     [bar] => banana
// )
// 使用指定的键和值填充数组
func arrayFill(startIndex int, num uint, value interface{}) map[int]interface{} {
	m := make(map[int]interface{})
	var i uint
	for i = 0; i < num; i++ {
		m[startIndex] = value
		startIndex++
	}
	return m
}

// ArrayFlip array_flip()
// 交换数组中的键和值. array_flip() 返回一个反转后的 array，例如 array 中的键名变成了值，而 array 中的值成了键名。
func arrayFlip(m map[interface{}]interface{}) map[interface{}]interface{} {
	n := make(map[interface{}]interface{})
	for i, v := range m {
		n[v] = i
	}
	return n
}

// ArrayKeys array_keys()
//  返回数组中部分的或所有的键名
func arrayKeys(elements map[interface{}]interface{}) []interface{} {
	i, keys := 0, make([]interface{}, len(elements))
	for key := range elements {
		keys[i] = key
		i++
	}
	return keys
}

// ArrayValues array_values()
//  返回数组中所有的值
func arrayValues(elements map[interface{}]interface{}) []interface{} {
	i, vals := 0, make([]interface{}, len(elements))
	for _, val := range elements {
		vals[i] = val
		i++
	}
	return vals
}

// ArrayMerge array_merge()
// 合并一个或多个数组
func arrayMerge(ss ...[]interface{}) []interface{} {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]interface{}, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

// ArrayChunk array_chunk()
//  将一个数组分割成多个
func arrayChunk(s []interface{}, size int) [][]interface{} {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]interface{}
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, s[i*size:end])
		i++
	}
	return n
}

// ArrayPad array_pad()
// array_pad() 返回 array 的一个拷贝，并用 value 将其填补到 size 指定的长度。如果 size 为正，则填补到数组的右侧，如果为负则从左侧开始填补。如果 size 的绝对值小于或等于 array 数组的长度则没有任何填补。有可能一次最多填补 1048576 个单元。
//  以指定长度将一个值填充进数组
func arrayPad(s []interface{}, size int, val interface{}) []interface{} {
	if size == 0 || (size > 0 && size < len(s)) || (size < 0 && size > -len(s)) {
		return s
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(s)
	tmp := make([]interface{}, n)
	for i := 0; i < n; i++ {
		tmp[i] = val
	}
	if size > 0 {
		return append(s, tmp...)
	}
	return append(tmp, s...)
}

// ArraySlice array_slice()
// array_slice() 返回根据 offset 和 length 参数所指定的 array 数组中的一段序列。
//  从数组中取出一段
func arraySlice(s []interface{}, offset, length uint) []interface{} {
	if offset > uint(len(s)) {
		panic("offset: the offset is less than the length of s")
	}
	end := offset + length
	if end < uint(len(s)) {
		return s[offset:end]
	}
	return s[offset:]
}

// ArrayRand array_rand()
//  从数组中随机取出一个或多个单元
func ArrayRand(elements []interface{}) []interface{} {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]interface{}, len(elements))
	for i, v := range r.Perm(len(elements)) {
		n[i] = elements[v]
	}
	return n
}

// ArrayColumn array_column()
// 返回数组中指定的一列
func arrayColumn(input map[string]map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns
}

// ArrayPush array_push()
// Push one or more elements onto the end of slice
//  将一个或多个单元压入数组的末尾（入栈）
func arrayPush(s *[]interface{}, elements ...interface{}) int {
	*s = append(*s, elements...)
	return len(*s)
}

// ArrayPop array_pop()
// Pop the element off the end of slice
// 弹出数组最后一个单元（出栈）
func ArrayPop(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	ep := len(*s) - 1
	e := (*s)[ep]
	*s = (*s)[:ep]
	return e
}

// ArrayUnshift array_unshift()
// Prepend one or more elements to the beginning of a slice
//  在数组开头插入一个或多个单元
func arrayUnshift(s *[]interface{}, elements ...interface{}) int {
	*s = append(elements, *s...)
	return len(*s)
}

// ArrayShift array_shift()
// Shift an element off the beginning of slice
//  将数组开头的单元移出数组
func arrayShift(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	f := (*s)[0]
	*s = (*s)[1:]
	return f
}

// ArrayKeyExists array_key_exists()
//  检查数组里是否有指定的键名或索引
func arrayKeyExists(key interface{}, m map[interface{}]interface{}) bool {
	_, ok := m[key]
	return ok
}

// ArrayCombine array_combine()
// 创建一个数组，用一个数组的值作为其键名，另一个数组的值作为其值
func ArrayCombine(s1, s2 []interface{}) map[interface{}]interface{} {
	if len(s1) != len(s2) {
		panic("the number of elements for each slice isn't equal")
	}
	m := make(map[interface{}]interface{}, len(s1))
	for i, v := range s1 {
		m[v] = s2[i]
	}
	return m
}

// ArrayReverse array_reverse()
//  返回单元顺序相反的数组
func arrayReverse(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// InArray in_array()
// haystack supported types: slice, array or map
// 检查数组中是否存在某个值
func inArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

//////////// Mathematical Functions ////////////

// Round round()
//  对浮点数进行四舍五入
func round(value float64) float64 {
	return math.Floor(value + 0.5)
}

// Floor floor()
//  舍去法取整. 返回不大于 value 的最接近的整数，舍去小数部分取整。
func floor(value float64) float64 {
	return math.Floor(value)
}

// Ceil ceil()
//  进一法取整. 返回不小于 value 的下一个整数，value 如果有小数部分则进一位。
func ceil(value float64) float64 {
	return math.Ceil(value)
}

// Pi pi()
//  得到圆周率值
func pi() float64 {
	return math.Pi
}

// Max max()
//  找出最大值
func max(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		max = math.Max(max, nums[i])
	}
	return max
}

// Min min()
//  找出最小值
func min(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	min := nums[0]
	for i := 1; i < len(nums); i++ {
		min = math.Min(min, nums[i])
	}
	return min
}

// Decbin decbin()
// 十进制转换为二进制
func decbin(number int64) string {
	return strconv.FormatInt(number, 2)
}

// Bindec bindec()
// 二进制转换为十进制
func bindec(str string) string {
	i, err := strconv.ParseInt(str, 2, 0)
	panicerr(err)
	return strconv.FormatInt(i, 10)
}

// Hex2bin hex2bin()
// 转换十六进制字符串为二进制字符串
func hex2bin(data string) string {
	i, err := strconv.ParseInt(data, 16, 0)
	panicerr(err)
	return strconv.FormatInt(i, 2)
}

// Bin2hex bin2hex()
// 函数把包含数据的二进制字符串转换为十六进制值
func bin2hex(str string) string {
	i, err := strconv.ParseInt(str, 2, 0)
	panicerr(err)
	return strconv.FormatInt(i, 16)
}

// Dechex dechex()
//  十进制转换为十六进制
func dechex(number int64) string {
	return strconv.FormatInt(number, 16)
}

// Hexdec hexdec()
//  十六进制转换为十进制
func hexdec(str string) int64 {
	i, err := strconv.ParseInt(str, 16, 0)
	panicerr(err)
	return i
}

// Decoct decoct()
//  十进制转换为八进制
func decoct(number int64) string {
	return strconv.FormatInt(number, 8)
}

// Octdec Octdec()
//  八进制转换为十进制
func octdec(str string) int64 {
	i, err := strconv.ParseInt(str, 8, 0)
	panicerr(err)
	return i
}

// BaseConvert base_convert()
//  在任意进制之间转换数字
func baseConvert(number string, frombase, tobase int) string {
	i, err := strconv.ParseInt(number, frombase, 0)
	panicerr(err)
	return strconv.FormatInt(i, tobase)
}

// IsNan is_nan()
//  判断是否为合法数值
func IsNan(val float64) bool {
	return math.IsNaN(val)
}

//////////// Directory/Filesystem Functions ////////////

// Stat stat()
//  给出文件的信息
func stat(filename string) os.FileInfo {
	f, err := os.Stat(filename)
	panicerr(err)
	return f
}

// Pathinfo pathinfo()
// -1: all; 1: dirname; 2: basename; 4: extension; 8: filename
// Usage:
// Pathinfo("/home/go/path/src/php2go/php2go.go", 1|2|4|8)
//  返回文件路径的信息
func pathinfo(path string) map[string]string {
	options := 1 | 2 | 4 | 8
	info := make(map[string]string)
	if (options & 1) == 1 {
		info["dirname"] = filepath.Dir(path)
	}
	if (options & 2) == 2 {
		info["basename"] = filepath.Base(path)
	}
	if ((options & 4) == 4) || ((options & 8) == 8) {
		basename := ""
		if (options & 2) == 2 {
			basename, _ = info["basename"]
		} else {
			basename = filepath.Base(path)
		}
		p := strings.LastIndex(basename, ".")
		filename, extension := "", ""
		if p > 0 {
			filename, extension = basename[:p], basename[p+1:]
		} else if p == -1 {
			filename = basename
		} else if p == 0 {
			extension = basename[p+1:]
		}
		if (options & 4) == 4 {
			info["extension"] = extension
		}
		if (options & 8) == 8 {
			info["filename"] = filename
		}
	}
	return info
}

// 取得文件大小
func fileSize(filename string) int64 {
	info, err := os.Stat(filename)
	panicerr(err)
	return info.Size()
}

// FilePutContents file_put_contents()
// 将一个字符串写入文件
func filePutContents(filename string, data string, mode os.FileMode) {
	err := ioutil.WriteFile(filename, []byte(data), mode)
	panicerr(err)
}

// FileGetContents file_get_contents()
// 将整个文件读入一个字符串
func fileGetContents(filename string) string {
	data, err := ioutil.ReadFile(filename)
	panicerr(err)
	return string(data)
}

// IsReadable is_readable()
//  判断给定文件名是否可读
func isReadable(filename string) bool {
	_, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return false
	}
	return true
}

// IsWriteable is_writeable()
// 判断给定的文件名是否可写
func isWriteable(filename string) bool {
	_, err := syscall.Open(filename, syscall.O_WRONLY, 0)
	if err != nil {
		return false
	}
	return true
}

// Touch touch()
// 设定文件的访问和修改时间
func touch(filename string) {
	fd, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	panicerr(err)
	fd.Close()
}

// Chmod chmod()
//  改变文件模式
func chmod(filename string, mode os.FileMode) bool {
	return os.Chmod(filename, mode) == nil
}

// Chown chown()
//  改变文件的所有者
func chown(filename string, uid, gid int) bool {
	return os.Chown(filename, uid, gid) == nil
}

// Fclose fclose()
// 关闭一个已打开的文件指针
func fclose(handle *os.File) {
	err := handle.Close()
	panicerr(err)
}

// Filemtime filemtime()
//  取得文件修改时间
func filemtime(filename string) int64 {
	fd, err := os.Open(filename)
	panicerr(err)
	defer fd.Close()
	fileinfo, err := fd.Stat()
	panicerr(err)
	return fileinfo.ModTime().Unix()
}

// Fgetcsv fgetcsv()
// 从文件指针中读入并解析 CSV 字段
func fgetcsv(handle *os.File, length int, delimiter rune) [][]string {
	reader := csv.NewReader(handle)
	reader.Comma = delimiter
	// TODO length limit
	str, err := reader.ReadAll()
	panicerr(err)
	return str
}

// Glob glob()
// 寻找与模式匹配的文件路径. glob() 函数依照 libc glob() 函数使用的规则寻找所有与 pattern 匹配的文件路径，类似于一般 shells 所用的规则一样。不进行缩写扩展或参数替代。
func glob(pattern string) []string {
	i, err := filepath.Glob(pattern)
	panicerr(err)
	return i
}

//////////// Variable handling Functions ////////////

// Empty empty()
// 检查一个变量是否为空
func empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

//////////// Program execution Functions ////////////

// Passthru passthru()
// returnVar, 0: succ; 1: fail
//  执行外部程序并且显示原始输出
func passthru(command string, returnVar *int) {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	// remove the " and ' on both sides
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		*returnVar = 1
		fmt.Println(err)
	} else {
		*returnVar = 0
	}
}

//////////// Network Functions ////////////

// Gethostbynamel gethostbynamel()
// Get a list of IPv4 addresses corresponding to a given Internet host name
//  获取互联网主机名对应的 IPv4 地址列表
// func gethostbynamel(hostname string) []string {
// 	ips, err := net.LookupIP(hostname)
// 	if err != nil {
// 		_, fn, line, _ := runtime.Caller(0)
// 		panic(filepath.Base(fn) + ":" + strconv.Itoa(line-2) + " >> " + err.Error() + " >> " + fmtDebugStack(string(debug.Stack())))
// 	}
// 	if ips != nil {
// 		var ipstrs []string
// 		for _, v := range ips {
// 			if v.To4() != nil {
// 				ipstrs = append(ipstrs, v.String())
// 			}
// 		}
// 		return ipstrs
// 	}
// 	return nil
// }

// Gethostbyaddr gethostbyaddr()
// Get the Internet host name corresponding to a given IP address
//  获取指定的IP地址对应的主机名
func gethostbyaddr(ipAddress string) string {
	names, err := net.LookupAddr(ipAddress)
	panicerr(err)
	if names != nil {
		return strings.TrimRight(names[0], ".")
	}
	return ""
}

// IP2long ip2long()
// IPv4
//  将 IPV4 的字符串互联网协议转换成长整型数字
func ip2long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

// Long2ip long2ip()
// IPv4
// 将长整型转化为字符串形式带点的互联网标准格式地址（IPV4）
func long2ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)
	return ip.String()
}

//////////// Misc. Functions ////////////

// Echo echo
//  输出一个或多个字符串
func echo(args ...interface{}) {
	fmt.Print(args...)
}

// Uniqid uniqid()
//  生成一个ID. 此函数不保证返回值的唯一性。
func uniqid(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%08x%05x", prefix, now.Unix(), now.UnixNano()%0x100000)
}

// Exit exit()
func exit(status int) {
	os.Exit(status)
}

// Die die()
func die(status int) {
	os.Exit(status)
}

// Putenv putenv()
// The setting, like "FOO=BAR"
func putenv(setting string) {
	s := strings.Split(setting, "=")
	if len(s) != 2 {
		panic("setting: invalid")
	}
	err := os.Setenv(s[0], s[1])
	panicerr(err)
}

// MemoryGetUsage memory_get_usage()
// return in bytes
// 返回分配给 golang 的内存量
func memoryGetUsage(realUsage bool) uint64 {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	return stat.Alloc
}

// VersionCompare version_compare()
// The possible operators are: <, lt, <=, le, >, gt, >=, ge, ==, =, eq, !=, <>, ne respectively.
// special version strings these are handled in the following order,
// (any string not found) < dev < alpha = a < beta = b < RC = rc < # < pl = p
// Usage:
// VersionCompare("1.2.3-alpha", "1.2.3RC7", '>=')
// VersionCompare("1.2.3-beta", "1.2.3pl", 'lt')
// VersionCompare("1.1_dev", "1.2any", 'eq')
//  对比两个「PHP 规范化」的版本数字字符串
func VersionCompare(version1, version2, operator string) bool {
	var vcompare func(string, string) int
	var canonicalize func(string) string
	var special func(string, string) int

	// version compare
	vcompare = func(origV1, origV2 string) int {
		if origV1 == "" || origV2 == "" {
			if origV1 == "" && origV2 == "" {
				return 0
			}
			if origV1 == "" {
				return -1
			}
			return 1
		}

		ver1, ver2, compare := "", "", 0
		if origV1[0] == '#' {
			ver1 = origV1
		} else {
			ver1 = canonicalize(origV1)
		}
		if origV2[0] == '#' {
			ver2 = origV2
		} else {
			ver2 = canonicalize(origV2)
		}
		n1, n2 := 0, 0
		for {
			p1, p2 := "", ""
			n1 = strings.IndexByte(ver1, '.')
			if n1 == -1 {
				p1, ver1 = ver1[:], ""
			} else {
				p1, ver1 = ver1[:n1], ver1[n1+1:]
			}
			n2 = strings.IndexByte(ver2, '.')
			if n2 == -1 {
				p2, ver2 = ver2, ""
			} else {
				p2, ver2 = ver2[:n2], ver2[n2+1:]
			}

			if (p1[0] >= '0' && p1[0] <= '9') && (p2[0] >= '0' && p2[0] <= '9') { // all is digit
				l1, _ := strconv.Atoi(p1)
				l2, _ := strconv.Atoi(p2)
				if l1 > l2 {
					compare = 1
				} else if l1 == l2 {
					compare = 0
				} else {
					compare = -1
				}
			} else if !(p1[0] >= '0' && p1[0] <= '9') && !(p2[0] >= '0' && p2[0] <= '9') { // all digit
				compare = special(p1, p2)
			} else { // part is digit
				if p1[0] >= '0' && p1[0] <= '9' { // is digit
					compare = special("#N#", p2)
				} else {
					compare = special(p1, "#N#")
				}
			}

			if compare != 0 || n1 == -1 || n2 == -1 {
				break
			}
		}

		if compare == 0 {
			if ver1 != "" {
				if ver1[0] >= '0' && ver1[0] <= '9' {
					compare = 1
				} else {
					compare = vcompare(ver1, "#N#")
				}
			} else if ver2 != "" {
				if ver2[0] >= '0' && ver2[0] <= '9' {
					compare = -1
				} else {
					compare = vcompare("#N#", ver2)
				}
			}
		}

		return compare
	}

	// canonicalize
	canonicalize = func(version string) string {
		ver := []byte(version)
		l := len(ver)
		if l == 0 {
			return ""
		}
		var buf = make([]byte, l*2)
		j := 0
		for i, v := range ver {
			next := uint8(0)
			if i+1 < l { // Have the next one
				next = ver[i+1]
			}
			if v == '-' || v == '_' || v == '+' { // replace '-', '_', '+' to '.'
				if j > 0 && buf[j-1] != '.' {
					buf[j] = '.'
					j++
				}
			} else if (next > 0) &&
				(!(next >= '0' && next <= '9') && (v >= '0' && v <= '9')) ||
				(!(v >= '0' && v <= '9') && (next >= '0' && next <= '9')) { // Insert '.' before and after a non-digit
				buf[j] = v
				j++
				if v != '.' && next != '.' {
					buf[j] = '.'
					j++
				}
				continue
			} else if !((v >= '0' && v <= '9') ||
				(v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z')) { // Non-letters and numbers
				if j > 0 && buf[j-1] != '.' {
					buf[j] = '.'
					j++
				}
			} else {
				buf[j] = v
				j++
			}
		}

		return string(buf[:j])
	}

	// compare special version forms
	special = func(form1, form2 string) int {
		found1, found2, len1, len2 := -1, -1, len(form1), len(form2)
		// (Any string not found) < dev < alpha = a < beta = b < RC = rc < # < pl = p
		forms := map[string]int{
			"dev":   0,
			"alpha": 1,
			"a":     1,
			"beta":  2,
			"b":     2,
			"RC":    3,
			"rc":    3,
			"#":     4,
			"pl":    5,
			"p":     5,
		}

		for name, order := range forms {
			if len1 < len(name) {
				continue
			}
			if strings.Compare(form1[:len(name)], name) == 0 {
				found1 = order
				break
			}
		}
		for name, order := range forms {
			if len2 < len(name) {
				continue
			}
			if strings.Compare(form2[:len(name)], name) == 0 {
				found2 = order
				break
			}
		}

		if found1 == found2 {
			return 0
		} else if found1 > found2 {
			return 1
		} else {
			return -1
		}
	}

	compare := vcompare(version1, version2)

	switch operator {
	case "<", "lt":
		return compare == -1
	case "<=", "le":
		return compare != 1
	case ">", "gt":
		return compare == 1
	case ">=", "ge":
		return compare != -1
	case "==", "=", "eq":
		return compare == 0
	case "!=", "<>", "ne":
		return compare != 0
	default:
		panic("operator: invalid")
	}
}

// ZipOpen zip_open()
//  打开ZIP存档文件
func zipOpen(filename string) *zip.ReadCloser {
	i, err := zip.OpenReader(filename)
	panicerr(err)
	return i
}

// Pack pack()
// 将数据打包成二进制字符串
func pack(order binary.ByteOrder, data interface{}) string {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	panicerr(err)

	return buf.String()
}

// Unpack unpack()
//  Unpack data from binary string
func unpack(order binary.ByteOrder, data string) interface{} {
	var result []byte
	r := bytes.NewReader([]byte(data))
	err := binary.Read(r, order, &result)
	panicerr(err)

	return result
}

// Ternary Ternary expression
// max := Ternary(a > b, a, b).(int)
func Ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
type fileStruct struct {
	path   string
	fd     *os.File
	mode   string
	reader *bufio.Reader
}

func (m *fileStruct) readlines() chan string {
	if m.reader == nil {
		m.reader = bufio.NewReader(m.fd)
	}

	lines := make(chan string)

	go func() {
		for {
			line, err := m.reader.ReadString('\n')
			if len(line) == 0 {
				if err != nil {
					close(lines)
					m.close()
					if err == io.EOF {
						return
					}
					_, fn, line, _ := runtime.Caller(0)
					panic(filepath.Base(fn) + ":" + strconv.Itoa(line-7) + " >> " + err.Error())
				}
			}
			line = strStrip(line, "\r\n")
			lines <- line
		}
	}()

	return lines
}

func (m *fileStruct) readline() string {
	b := make([]byte, 1)

	line := ""

	for {
		_, err := io.ReadAtLeast(m.fd, b, 1)
		if err != nil {
			if len(line) != 0 {
				return line
			}
			panicerr(err)
		}
		bs := string(b)
		if bs == "\n" {
			return line
		}
		line = line + bs
	}
}

func (m *fileStruct) close() {
	m.fd.Close()
}

func (m *fileStruct) write(str interface{}) *fileStruct {
	var buf []byte
	if typeof(str) == "string" {
		s := str.(string)
		buf = []byte(s)
	} else {
		s := str.([]byte)
		buf = s
	}
	m.fd.Write(buf)
	return m
}

func (m *fileStruct) read(num ...int) string {
	var bytes []byte
	var err error
	if len(num) == 0 {
		bytes, err = ioutil.ReadAll(m.fd)
		panicerr(err)
		m.close()
	} else {
		buffer := make([]byte, num[0])
		i, err := io.ReadAtLeast(m.fd, buffer, num[0])
		if err != nil {
			if !strIn("EOF", err.Error()) {
				panicerr(err)
			}
		}
		bytes = buffer[:i]
	}

	return string(bytes)
}

func (m *fileStruct) seek(num int64) {
	_, err := m.fd.Seek(num, 0)
	panicerr(err)
}

func open(args ...string) *fileStruct {
	path := args[0]
	var mode string
	if len(args) != 1 {
		mode = args[1]
	} else {
		mode = "r"
	}
	var fd *os.File
	var err error
	if string(mode[0]) == "r" {
		fd, err = os.Open(path)
	}
	if string(mode[0]) == "a" {
		fd, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if string(mode[0]) == "w" {
		fd, err = os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	}
	panicerr(err)
	return &fileStruct{
		path: path,
		fd:   fd,
		mode: mode,
	}
}

func getStdin() *fileStruct {
	return &fileStruct{fd: os.Stdin, mode: "r"}
}
var ip17modHadInit bool
var ip17DBFileModifyTime int64

type ipLocationInfo struct {
	Country string
	Region  string
	City    string
	Isp     string
}

func getIPLocation(ip string, dbpath ...string) *ipLocationInfo {
	if len(dbpath) == 0 {
		if ip17modHadInit == false {
			statikFS, err := fs.New()
			panicerr(err)

			ipdbfd, err := statikFS.Open("/qqwry.ipdb")
			panicerr(err)
			ipdbBytes, err := ioutil.ReadAll(ipdbfd)
			panicerr(err)
			ipdbfd.Close()

			ip17mon.InitWithIpdb(ipdbBytes)

			ip17modHadInit = true
		}
	} else {
		if ip17DBFileModifyTime == 0 || getFileTime(dbpath[0]).mtime != ip17DBFileModifyTime {
			ip17DBFileModifyTime = getFileTime(dbpath[0]).mtime

			ip17mon.InitWithIpdb([]byte(open(dbpath[0]).read()))
		}
	}

	loc, err := ip17mon.Find(ip)
	panicerr(err)

	return &ipLocationInfo{
		City:    loc.City,
		Region:  loc.Region,
		Country: loc.Country,
		Isp:     loc.Isp,
	}
}
type ttlCache struct {
	sync.RWMutex
	data    string
	expires *time.Time
}

func (item *ttlCache) touch(duration time.Duration) {
	item.Lock()
	expiration := time.Now().Add(duration)
	item.expires = &expiration
	item.Unlock()
}

func (item *ttlCache) expired() bool {
	var value bool
	item.RLock()
	if item.expires == nil {
		value = true
	} else {
		value = item.expires.Before(time.Now())
	}
	item.RUnlock()
	return value
}

type ttlCacheSCache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[string]*ttlCache
}

func (cache *ttlCacheSCache) Set(key string, data string) {
	cache.mutex.Lock()
	item := &ttlCache{data: data}
	item.touch(cache.ttl)
	cache.items[key] = item
	cache.mutex.Unlock()
}

func (cache *ttlCacheSCache) Get(key string) (data string, found bool) {
	cache.mutex.Lock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = ""
		found = false
	} else {
		item.touch(cache.ttl)
		data = item.data
		found = true
	}
	cache.mutex.Unlock()
	return
}

func (cache *ttlCacheSCache) Count() int {
	cache.mutex.RLock()
	count := len(cache.items)
	cache.mutex.RUnlock()
	return count
}

func (cache *ttlCacheSCache) cleanup() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

func (cache *ttlCacheSCache) startCleanupTimer() {
	ticker := time.Tick(3 * time.Second) // 3秒刷一次缓存
	go (func() {
		for {
			select {
			case <-ticker:
				cache.cleanup()
			}
		}
	})()
}

func NewCache(duration time.Duration) *ttlCacheSCache {
	cache := &ttlCacheSCache{
		ttl:   duration,
		items: map[string]*ttlCache{},
	}
	cache.startCleanupTimer()
	return cache
}

type ttlCacheStruct struct {
	cache *ttlCacheSCache
}

func getTTLCache(ttlsecond interface{}) *ttlCacheStruct {
	return &ttlCacheStruct{
		cache: NewCache(getTimeDuration(ttlsecond)),
	}
}

func (m *ttlCacheStruct) set(key string, value string) {
	m.cache.Set(key, value)
}

func (m *ttlCacheStruct) get(key string) string {
	value, exists := m.cache.Get(key)
	if exists != true {
		panicerr("Key " + key + " not found in cache")
	}
	return value
}

func (m *ttlCacheStruct) exists(key string) bool {
	_, exists := m.cache.Get(key)
	return exists
}

func (m *ttlCacheStruct) count() int {
	return m.cache.Count()
}
func gbkToUtf8(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewDecoder())
	d, err := ioutil.ReadAll(reader)
	panicerr(err)
	return string(d)
}

func utf8ToGbk(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
	d, err := ioutil.ReadAll(reader)
	panicerr(err)
	return string(d)
}
type cloudflareRecord struct {
	ID                       string
	Data                     string
	Name                     string
	TTL                      int
	Type                     string
	Proxiable                bool
	Proxied                  bool
	OriginalCloudflareRecord cloudflare.DNSRecord
}

type cloudflareStruct struct {
	api   *cloudflare.API
	key   string
	email string
}

type cloudflareDomainInfoStruct struct {
	CreatedAt   time.Time
	Domain      string
	DomainID    string
	NameServers []string
	Status      string
}

type cloudflareDomainStruct struct {
	api        *cloudflare.API
	DomainID   string
	DomainName string
}

func getCloudflare(key string, email string) *cloudflareStruct {
	api, err := cloudflare.New(key, email)
	panicerr(err)
	return &cloudflareStruct{
		api:   api,
		key:   key,
		email: email,
	}
}

func (m *cloudflareStruct) add(domain string) cloudflare.Zone {
	zone, err := m.api.CreateZone(context.Background(), domain, false, cloudflare.Account{}, "full")
	panicerr(err)
	return zone
}

func (m *cloudflareStruct) list() (res []cloudflareDomainInfoStruct) {
	zones, err := m.api.ListZones(context.Background())
	panicerr(err)
	for _, zone := range zones {
		res = append(res, cloudflareDomainInfoStruct{
			CreatedAt:   zone.CreatedOn,
			Domain:      zone.Name,
			DomainID:    zone.ID,
			NameServers: zone.NameServers,
			Status:      zone.Status,
		})
	}
	return
}

func (m *cloudflareStruct) domain(domainName string) *cloudflareDomainStruct {
	id, err := m.api.ZoneIDByName(domainName)
	panicerr(err)
	return &cloudflareDomainStruct{
		api:        m.api,
		DomainID:   id,
		DomainName: domainName,
	}
}

func (m *cloudflareDomainStruct) list() (res []cloudflareRecord) {
	records, err := m.api.DNSRecords(context.Background(), m.DomainID, cloudflare.DNSRecord{})
	panicerr(err)

	var name string
	for _, record := range records {
		if record.Name == record.ZoneName {
			name = "@"
		} else {
			name = record.Name[:len(record.Name)-len(record.ZoneName)-1]
		}
		res = append(res, cloudflareRecord{
			Data:                     record.Content,
			Name:                     name,
			TTL:                      record.TTL,
			Type:                     strLower(record.Type),
			Proxiable:                record.Proxiable,
			Proxied:                  *record.Proxied,
			ID:                       record.ID,
			OriginalCloudflareRecord: record,
		})
	}

	return
}

func (m *cloudflareDomainStruct) delete(name string) {
	// 虽然之后有log的代码, 但是在这个函数里面关闭log
	// logLevel := lg.getLevel()
	// lg.setLevel("")
	//defer lg.setLevel(logLevel)

	// if name == "@" {
	// 	name = m.DomainName
	// }

	// lg.debug(name)
	for _, v := range m.list() {
		// lg.debug(v)
		if name == v.Name {
			err := m.api.DeleteDNSRecord(context.Background(), m.DomainID, v.ID)
			panicerr(err)
		}
	}
}

func (m *cloudflareDomainStruct) add(recordName string, recordType string, recordValue string, proxied ...bool) *cloudflare.DNSRecordResponse {
	if recordName == "@" {
		recordName = m.DomainName
	} else {
		recordName = recordName + "." + m.DomainName
	}
	var prox bool
	if len(proxied) == 0 {
		prox = false
	} else {
		prox = proxied[0]
	}
	var p uint16 = 10
	resp, err := m.api.CreateDNSRecord(context.Background(), m.DomainID, cloudflare.DNSRecord{
		Type:     strUpper(recordType),
		Name:     recordName,
		Content:  recordValue,
		TTL:      300,
		Priority: &p,
		Proxied:  &prox,
	})
	panicerr(err)
	return resp
}

func (m *cloudflareDomainStruct) setProxied(subdomain string, proxied bool) {
	for _, v := range m.list() {
		//lg.trace(v.Name, domain, proxied)
		if v.Name == subdomain {
			if v.Proxiable == false && proxied == true {
				panic(newerr("类型为" + v.Type + "的dns记录无法设置cdn代理加速"))
			} else if v.Proxied != proxied {
				v.OriginalCloudflareRecord.Proxied = &proxied
				m.api.UpdateDNSRecord(context.Background(), m.DomainID, v.ID, v.OriginalCloudflareRecord)
			}
		}
	}
}

func (m *cloudflareDomainStruct) update(recordName string, recordType string, recordValue string, proxied ...bool) {
	var prox bool
	if len(proxied) == 0 {
		prox = false
	} else {
		prox = proxied[0]
	}

	for _, v := range m.list() {
		//lg.trace(v.Name, recordName)
		if v.Name == recordName {
			m.delete(recordName)
			if len(proxied) == 0 {
				prox = v.Proxied
			}
		}
	}

	m.add(recordName, recordType, recordValue, prox)
}
type kcpRawServerSideListener struct {
	listener *kcp.Listener
}

func kcpRawListen(host string, port int, key string, salt string) *kcpRawServerSideListener {
	block, err := kcp.NewAESBlockCrypt(pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha256.New))
	panicerr(err)

	l, err := kcp.ListenWithOptions(host+":"+toString(port), block, 10, 3)
	panicerr(err)

	l.SetDSCP(46)
	l.SetReadBuffer(4194304)
	l.SetWriteBuffer(4194304)

	return &kcpRawServerSideListener{listener: l}
}

func (m *kcpRawServerSideListener) accept() chan *kcp.UDPSession {
	ch := make(chan *kcp.UDPSession)

	go func() {
		for {
			c, err := m.listener.AcceptKCP()
			if err != nil {
				if strIn("io: read/write on closed pipe", err.Error()) || strIn("use of closed network connection", err.Error()) {
					close(ch)
					break
				}
				panicerr(err)
			}

			c.SetNoDelay(0, 20, 2, 1)
			c.SetMtu(1400)
			c.SetWindowSize(1024, 1024)
			c.SetACKNoDelay(false)

			ch <- c
		}
	}()

	return ch
}

func kcpRawConnect(host string, port int, key string, salt string) *kcp.UDPSession {
	block, err := kcp.NewAESBlockCrypt(pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha256.New))
	panicerr(err)
	conn, err := kcp.DialWithOptions(host+":"+toString(port), block, 10, 3)
	panicerr(err)

	conn.SetMtu(1400)
	conn.SetWriteDelay(false)
	conn.SetNoDelay(0, 20, 2, 1)
	conn.SetWindowSize(128, 1024)
	conn.SetACKNoDelay(false)
	conn.SetDSCP(46)
	conn.SetReadBuffer(4194304)
	conn.SetWriteBuffer(4194304)

	return conn
}
type pexpectStruct struct {
	buf         string         // sendline之前的命令输出，每次sendline都会清空
	bufall      string         // 所有屏幕的显示内容，包括了输入
	in          io.WriteCloser // 不用管，会被sendline拿来输入
	cmd         *exec.Cmd      // 不用管，会被close()拿来控制子进程
	logToStdout bool           // 是否在屏幕打印出整个交互（适合做debug)
	isAlive     bool           // 子进程是否在运行
	exitCode    int            // 子进程的退出代码
	bufStdout   io.Writer
}

func (m *pexpectStruct) Write(p []byte) (n int, err error) {
	m.buf = m.buf + string(p)
	m.bufall = m.bufall + m.buf
	if m.logToStdout {
		os.Stdout.Write(p)
		//fmt.Fprint(m.bufStdout, string(p))
	}
	return len(p), nil
}

func (m *pexpectStruct) sendline(str string) {
	m.buf = ""
	m.bufall = m.bufall + str + "\n"
	if m.logToStdout {
		fmt.Print(str + "\n")
	}
	_, err := io.WriteString(m.in, str+"\n")
	panicerr(err)
}

func (m *pexpectStruct) close() {
	err := m.cmd.Process.Kill()
	panicerr(err)
}

func pexpect(command string) *pexpectStruct {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	// remove the " and ' on both sides
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}

	m := pexpectStruct{
		isAlive: true,
	}

	if !cmdExists(parts[0]) {
		panicerr("Command not exists")
	}

	go func() {
		cmd := exec.Command(parts[0], parts[1:]...)
		m.cmd = cmd

		cmd.Stdout = &m
		cmd.Stderr = &m

		stdin, err := cmd.StdinPipe()
		panicerr(err)
		defer stdin.Close()

		m.in = stdin

		err = cmd.Start()
		panicerr(err)
		err = cmd.Wait()
		if err != nil {
			e := err.(*exec.ExitError)
			m.exitCode = e.ExitCode()
		} else {
			m.exitCode = 0
		}
		m.isAlive = false
	}()

	return &m
}
func getLanguage(text string) string {
	return whatlanggo.Detect(text).Lang.String()
}
type aesStruct struct {
	key []byte
}

func getAES(key string) *aesStruct {
	return &aesStruct{key: []byte(key)}
}

func (a aesStruct) encrypt(plaintext string) string {
	block, err := aes.NewCipher(a.key)
	panicerr(err)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	panicerr(err)
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return string(ciphertext)
	//return hex.EncodeToString(ciphertext)
}

func (a aesStruct) decrypt(d string) string {
	ciphertext := []byte(d)
	//ciphertext, err := hex.DecodeString(d)
	block, err := aes.NewCipher(a.key)
	panicerr(err)
	if len(ciphertext) < aes.BlockSize {
		panic(newerr("ciphertext too short"))
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext)
}
type argparseIniStruct struct {
	cfg         *ini.File
	cfgPath     string
	args        *orderedmap.OrderedMap
	description string
}

func argparser(description string) *argparseIniStruct {
	var cfg *ini.File
	var err error
	var cfgPath string

	for idx, value := range os.Args {
		if value == "-c" {
			cfgPath = os.Args[idx+1]
		}
	}

	if cfgPath == "" {
		c := pathJoin(getSelfDir(), basename(os.Args[0])+".ini")
		if pathExists(c) {
			cfgPath = c
		} else {
			c = pathJoin(basename(os.Args[0]) + ".ini")
			if pathExists(c) {
				cfgPath = c
			}
		}
	}

	if len(cfgPath) != 0 {
		if pathExists(cfgPath) {
			cfg, err = ini.Load(cfgPath)
			panicerr(err)
		} else {
			cfg = ini.Empty()
		}
	} else {
		cfgPath = ""
		cfg = ini.Empty()
	}

	return &argparseIniStruct{cfg: cfg, cfgPath: cfgPath, args: orderedmap.NewOrderedMap(), description: description}
}

func (m *argparseIniStruct) get(section, key, defaultValue, comment string) (res string) {
	res = m.cfg.Section(section).Key(key).String()
	if res == "" {
		res = defaultValue
		m.cfg.Section(section).Key(key).SetValue(defaultValue)
	}
	m.cfg.Section(section).Key(key).Comment = comment
	if section != "" {
		m.args.Set(section+"."+key, comment)
	} else {
		m.args.Set(key, comment)
	}

	if section != "" {
		if os.Getenv(section+"."+key) != "" {
			res = os.Getenv(section + "." + key)
		}
	} else {
		if os.Getenv(key) != "" {
			res = os.Getenv(key)
		}
	}

	for idx, value := range os.Args {
		if section == "" {
			if "--"+key == value {
				res = os.Args[idx+1]
			}
		} else {
			if "--"+section+"."+key == value {
				res = os.Args[idx+1]
			}
		}

	}
	return
}

func (m *argparseIniStruct) getInt(section, key, defaultValue, comment string) int {
	return toInt(m.get(section, key, defaultValue, comment))
}

func (m *argparseIniStruct) getInt64(section, key, defaultValue, comment string) int64 {
	return toInt64(m.get(section, key, defaultValue, comment))
}

func (m *argparseIniStruct) getFloat64(section, key, defaultValue, comment string) float64 {
	return toFloat64(m.get(section, key, defaultValue, comment))
}

func (m *argparseIniStruct) getBool(section, key, defaultValue, comment string) bool {
	return toBool(m.get(section, key, defaultValue, comment))
}

func (m *argparseIniStruct) save(fpath ...string) (exist bool) {
	exist = true
	if len(fpath) != 0 {
		exist = pathExists(fpath[0])
		m.cfg.SaveTo(fpath[0])
	} else if m.cfgPath != "" {
		exist = pathExists(m.cfgPath)
		m.cfg.SaveTo(m.cfgPath)
	}
	return
}

func (m *argparseIniStruct) getHelpString() (h string) {
	maxLength := 0
	for _, k := range m.args.Keys() {
		if len(toString(k)) > maxLength {
			maxLength = len(toString(k))
		}
	}

	h = "\n    " + m.description + "\n\n"
	h = h + "    -" + fmt.Sprintf("%-"+toString(maxLength+5+1)+"v", "c") + " 配置文件\n"
	h = h + "    -" + fmt.Sprintf("%-"+toString(maxLength+5+1)+"v", "b") + " 在后台运行\n"
	for _, k := range m.args.Keys() {
		v, _ := m.args.Get(k)
		h = h + "    --" + fmt.Sprintf("%-"+toString(maxLength+5)+"v", toString(k)) + " " + toString(v) + "\n"
	}
	return
}

func (m *argparseIniStruct) parseArgs() *argparseIniStruct {
	if itemInArray("-h", os.Args) || itemInArray("--help", os.Args) {
		fmt.Println(m.getHelpString())
		exit(0)
	}
	if itemInArray("-b", os.Args) {
		if runtime.GOOS == "linux" {
			args := os.Args[1:]
			for i := 0; i < len(args); i++ {
				if args[i] == "-b" {
					args[i] = ""
					break
				}
			}
			cmd := exec.Command(os.Args[0], args...)
			cmd.Start()
			os.Exit(0)
		} else {
			fmt.Println("错误: -b 参数只在Linux环境测试过")
			exit(0)
		}
	}
	if !m.save() { // 配置文件存在或者未指定则返回true，如果不存在则返回true
		fmt.Println("初始配置文件不存在，已保存模板配置文件。")
		exit(0)
	}
	return m
}
func fakeName() string {
	return faker.Name().Name()
}
type sshStruct struct {
	client  *ssh.Client
	session *ssh.Session
	user    string
	pass    string
	host    string
	port    int
}

func sshConnect(user string, pass string, host string, port int) *sshStruct {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", host+":"+toString(port), sshConfig)
	panicerr(err)

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panicerr(err)
	}

	return &sshStruct{
		client:  client,
		session: session,
		user:    user,
		pass:    pass,
		port:    port,
		host:    host,
	}
}

func (m *sshStruct) exec(cmd string) (output string) {
	out, err := m.session.CombinedOutput(cmd)
	panicerr(err)
	return string(out)
}

func (m *sshStruct) pushFile(local string, remote string) {
	sftp, err := sftp.NewClient(m.client)
	panicerr(err)
	defer sftp.Close()

	fr, err := sftp.Create(remote)
	panicerr(err)

	fl, err := os.Open(local)
	panicerr(err)

	io.Copy(fr, fl)

	err = fr.Close()
	panicerr(err)

	err = fl.Close()
	panicerr(err)
}

func (m *sshStruct) pullFile(remote string, local string) {
	sftp, err := sftp.NewClient(m.client)
	panicerr(err)
	defer sftp.Close()

	fr, err := sftp.Open(remote)
	panicerr(err)

	fl, err := os.Create(local)
	panicerr(err)

	io.Copy(fl, fr)

	err = fr.Close()
	panicerr(err)

	err = fl.Close()
	panicerr(err)
}
var errNegativeNotAllowed = errors.New("unable to cast negative value")

// toTimeE casts an interface to a time.Time type.
func toTimeE(i interface{}) (tim time.Time, err error) {
	i = indirect(i)

	switch v := i.(type) {
	case time.Time:
		return v, nil
	case string:
		return stringToDate(v)
	case int:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case int32:
		return time.Unix(int64(v), 0), nil
	case uint:
		return time.Unix(int64(v), 0), nil
	case uint64:
		return time.Unix(int64(v), 0), nil
	case uint32:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
	}
}

// toDurationE casts an interface to a time.Duration type.
func toDurationE(i interface{}) (d time.Duration, err error) {
	i = indirect(i)

	switch s := i.(type) {
	case time.Duration:
		return s, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		d = time.Duration(toInt64(s))
		return
	case float32, float64:
		d = time.Duration(toFloat64(s))
		return
	case string:
		if strings.ContainsAny(s, "nsuµmh") {
			d, err = time.ParseDuration(s)
		} else {
			d, err = time.ParseDuration(s + "ns")
		}
		return
	default:
		err = fmt.Errorf("unable to cast %#v of type %T to Duration", i, i)
		return
	}
}

// toBoolE casts an interface to a bool type.
func toBoolE(i interface{}) (bool, error) {
	i = indirect(i)

	switch b := i.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int:
		if i.(int) != 0 {
			return true, nil
		}
		return false, nil
	case string:
		return strconv.ParseBool(i.(string))
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	}
}

// toFloat64E casts an interface to a float64 type.
func toFloat64E(i interface{}) (float64, error) {
	i = indirect(i)

	switch s := i.(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	case int:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int8:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case string:
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	}
}

// toFloat32E casts an interface to a float32 type.
func toFloat32E(i interface{}) (float32, error) {
	i = indirect(i)

	switch s := i.(type) {
	case float64:
		return float32(s), nil
	case float32:
		return s, nil
	case int:
		return float32(s), nil
	case int64:
		return float32(s), nil
	case int32:
		return float32(s), nil
	case int16:
		return float32(s), nil
	case int8:
		return float32(s), nil
	case uint:
		return float32(s), nil
	case uint64:
		return float32(s), nil
	case uint32:
		return float32(s), nil
	case uint16:
		return float32(s), nil
	case uint8:
		return float32(s), nil
	case string:
		v, err := strconv.ParseFloat(s, 32)
		if err == nil {
			return float32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to float32", i, i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float32", i, i)
	}
}

// toInt64E casts an interface to an int64 type.
func toInt64E(i interface{}) (int64, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return int64(s), nil
	case int64:
		return s, nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case uint:
		return int64(s), nil
	case uint64:
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint8:
		return int64(s), nil
	case float64:
		return int64(s), nil
	case float32:
		return int64(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	}
}

// toInt32E casts an interface to an int32 type.
func toInt32E(i interface{}) (int32, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return int32(s), nil
	case int64:
		return int32(s), nil
	case int32:
		return s, nil
	case int16:
		return int32(s), nil
	case int8:
		return int32(s), nil
	case uint:
		return int32(s), nil
	case uint64:
		return int32(s), nil
	case uint32:
		return int32(s), nil
	case uint16:
		return int32(s), nil
	case uint8:
		return int32(s), nil
	case float64:
		return int32(s), nil
	case float32:
		return int32(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return int32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int32", i, i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int32", i, i)
	}
}

// toInt16E casts an interface to an int16 type.
func toInt16E(i interface{}) (int16, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return int16(s), nil
	case int64:
		return int16(s), nil
	case int32:
		return int16(s), nil
	case int16:
		return s, nil
	case int8:
		return int16(s), nil
	case uint:
		return int16(s), nil
	case uint64:
		return int16(s), nil
	case uint32:
		return int16(s), nil
	case uint16:
		return int16(s), nil
	case uint8:
		return int16(s), nil
	case float64:
		return int16(s), nil
	case float32:
		return int16(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return int16(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int16", i, i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int16", i, i)
	}
}

// toInt8E casts an interface to an int8 type.
func toInt8E(i interface{}) (int8, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return int8(s), nil
	case int64:
		return int8(s), nil
	case int32:
		return int8(s), nil
	case int16:
		return int8(s), nil
	case int8:
		return s, nil
	case uint:
		return int8(s), nil
	case uint64:
		return int8(s), nil
	case uint32:
		return int8(s), nil
	case uint16:
		return int8(s), nil
	case uint8:
		return int8(s), nil
	case float64:
		return int8(s), nil
	case float32:
		return int8(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return int8(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int8", i, i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int8", i, i)
	}
}

// toIntE casts an interface to an int type.
func toIntE(i interface{}) (int, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return s, nil
	case int64:
		return int(s), nil
	case int32:
		return int(s), nil
	case int16:
		return int(s), nil
	case int8:
		return int(s), nil
	case uint:
		return int(s), nil
	case uint64:
		return int(s), nil
	case uint32:
		return int(s), nil
	case uint16:
		return int(s), nil
	case uint8:
		return int(s), nil
	case float64:
		return int(s), nil
	case float32:
		return int(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return int(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	}
}

// toUintE casts an interface to a uint type.
func toUintE(i interface{}) (uint, error) {
	i = indirect(i)

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseUint(s, 0, 0)
		if err == nil {
			return uint(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v to uint: %s", i, err)
	case int:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case uint:
		return s, nil
	case uint64:
		return uint(s), nil
	case uint32:
		return uint(s), nil
	case uint16:
		return uint(s), nil
	case uint8:
		return uint(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint", i, i)
	}
}

// toUint64E casts an interface to a uint64 type.
func toUint64E(i interface{}) (uint64, error) {
	i = indirect(i)

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseUint(s, 0, 64)
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v to uint64: %s", i, err)
	case int:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case uint:
		return uint64(s), nil
	case uint64:
		return s, nil
	case uint32:
		return uint64(s), nil
	case uint16:
		return uint64(s), nil
	case uint8:
		return uint64(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint64", i, i)
	}
}

// toUint32E casts an interface to a uint32 type.
func toUint32E(i interface{}) (uint32, error) {
	i = indirect(i)

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseUint(s, 0, 32)
		if err == nil {
			return uint32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v to uint32: %s", i, err)
	case int:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case uint:
		return uint32(s), nil
	case uint64:
		return uint32(s), nil
	case uint32:
		return s, nil
	case uint16:
		return uint32(s), nil
	case uint8:
		return uint32(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint32", i, i)
	}
}

// toUint16E casts an interface to a uint16 type.
func toUint16E(i interface{}) (uint16, error) {
	i = indirect(i)

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseUint(s, 0, 16)
		if err == nil {
			return uint16(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v to uint16: %s", i, err)
	case int:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case uint:
		return uint16(s), nil
	case uint64:
		return uint16(s), nil
	case uint32:
		return uint16(s), nil
	case uint16:
		return s, nil
	case uint8:
		return uint16(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint16", i, i)
	}
}

// toUint8E casts an interface to a uint type.
func toUint8E(i interface{}) (uint8, error) {
	i = indirect(i)

	switch s := i.(type) {
	case string:
		v, err := strconv.ParseUint(s, 0, 8)
		if err == nil {
			return uint8(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v to uint8: %s", i, err)
	case int:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case uint:
		return uint8(s), nil
	case uint64:
		return uint8(s), nil
	case uint32:
		return uint8(s), nil
	case uint16:
		return uint8(s), nil
	case uint8:
		return s, nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint8", i, i)
	}
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirecttoStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error,
func indirecttoStringerOrError(a interface{}) interface{} {
	if a == nil {
		return nil
	}

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// toStringE casts an interface to a string type.
func toStringE(i interface{}) (string, error) {
	i = indirecttoStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

// toStringMapStringE casts an interface to a map[string]string type.
func toStringMapStringE(i interface{}) (map[string]string, error) {
	var m = map[string]string{}

	switch v := i.(type) {
	case map[string]string:
		return v, nil
	case map[string]interface{}:
		for k, val := range v {
			m[toString(k)] = toString(val)
		}
		return m, nil
	case map[interface{}]string:
		for k, val := range v {
			m[toString(k)] = toString(val)
		}
		return m, nil
	case map[interface{}]interface{}:
		for k, val := range v {
			m[toString(k)] = toString(val)
		}
		return m, nil
	case string:
		err := jsonStringToObject(v, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]string", i, i)
	}
}

// toStringMapStringSliceE casts an interface to a map[string][]string type.
func toStringMapStringSliceE(i interface{}) (map[string][]string, error) {
	var m = map[string][]string{}

	switch v := i.(type) {
	case map[string][]string:
		return v, nil
	case map[string][]interface{}:
		for k, val := range v {
			m[toString(k)] = toStringSlice(val)
		}
		return m, nil
	case map[string]string:
		for k, val := range v {
			m[toString(k)] = []string{val}
		}
	case map[string]interface{}:
		for k, val := range v {
			switch vt := val.(type) {
			case []interface{}:
				m[toString(k)] = toStringSlice(vt)
			case []string:
				m[toString(k)] = vt
			default:
				m[toString(k)] = []string{toString(val)}
			}
		}
		return m, nil
	case map[interface{}][]string:
		for k, val := range v {
			m[toString(k)] = toStringSlice(val)
		}
		return m, nil
	case map[interface{}]string:
		for k, val := range v {
			m[toString(k)] = toStringSlice(val)
		}
		return m, nil
	case map[interface{}][]interface{}:
		for k, val := range v {
			m[toString(k)] = toStringSlice(val)
		}
		return m, nil
	case map[interface{}]interface{}:
		for k, val := range v {
			key, err := toStringE(k)
			if err != nil {
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			}
			value, err := toStringSliceE(val)
			if err != nil {
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			}
			m[key] = value
		}
	case string:
		err := jsonStringToObject(v, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
	}
	return m, nil
}

// toStringMapBoolE casts an interface to a map[string]bool type.
func toStringMapBoolE(i interface{}) (map[string]bool, error) {
	var m = map[string]bool{}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[toString(k)] = toBool(val)
		}
		return m, nil
	case map[string]interface{}:
		for k, val := range v {
			m[toString(k)] = toBool(val)
		}
		return m, nil
	case map[string]bool:
		return v, nil
	case string:
		err := jsonStringToObject(v, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]bool", i, i)
	}
}

// toStringMapE casts an interface to a map[string]interface{} type.
func toStringMapE(i interface{}) (map[string]interface{}, error) {
	var m = map[string]interface{}{}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[toString(k)] = val
		}
		return m, nil
	case map[string]interface{}:
		return v, nil
	case string:
		err := jsonStringToObject(v, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]interface{}", i, i)
	}
}

// toStringMapIntE casts an interface to a map[string]int{} type.
func toStringMapIntE(i interface{}) (map[string]int, error) {
	var m = map[string]int{}
	if i == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
	}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[toString(k)] = toInt(val)
		}
		return m, nil
	case map[string]interface{}:
		for k, val := range v {
			m[k] = toInt(val)
		}
		return m, nil
	case map[string]int:
		return v, nil
	case string:
		err := jsonStringToObject(v, &m)
		return m, err
	}

	if reflect.TypeOf(i).Kind() != reflect.Map {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
	}

	mVal := reflect.ValueOf(m)
	v := reflect.ValueOf(i)
	for _, keyVal := range v.MapKeys() {
		val, err := toIntE(v.MapIndex(keyVal).Interface())
		if err != nil {
			return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
		}
		mVal.SetMapIndex(keyVal, reflect.ValueOf(val))
	}
	return m, nil
}

// toStringMapInt64E casts an interface to a map[string]int64{} type.
func toStringMapInt64E(i interface{}) (map[string]int64, error) {
	var m = map[string]int64{}
	if i == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
	}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[toString(k)] = toInt64(val)
		}
		return m, nil
	case map[string]interface{}:
		for k, val := range v {
			m[k] = toInt64(val)
		}
		return m, nil
	case map[string]int64:
		return v, nil
	case string:
		err := jsonStringToObject(v, &m)
		return m, err
	}

	if reflect.TypeOf(i).Kind() != reflect.Map {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
	}
	mVal := reflect.ValueOf(m)
	v := reflect.ValueOf(i)
	for _, keyVal := range v.MapKeys() {
		val, err := toInt64E(v.MapIndex(keyVal).Interface())
		if err != nil {
			return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
		}
		mVal.SetMapIndex(keyVal, reflect.ValueOf(val))
	}
	return m, nil
}

// toSliceE casts an interface to a []interface{} type.
func toSliceE(i interface{}) ([]interface{}, error) {
	var s []interface{}

	switch v := i.(type) {
	case []interface{}:
		return append(s, v...), nil
	case []map[string]interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []interface{}", i, i)
	}
}

// toBoolSliceE casts an interface to a []bool type.
func toBoolSliceE(i interface{}) ([]bool, error) {
	if i == nil {
		return []bool{}, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
	}

	switch v := i.(type) {
	case []bool:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]bool, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := toBoolE(s.Index(j).Interface())
			if err != nil {
				return []bool{}, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []bool{}, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
	}
}

// toStringSliceE casts an interface to a []string type.
func toStringSliceE(i interface{}) ([]string, error) {
	var a []string

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			a = append(a, toString(u))
		}
		return a, nil
	case []string:
		return v, nil
	case []int8:
		for _, u := range v {
			a = append(a, toString(u))
		}
		return a, nil
	case []int:
		for _, u := range v {
			a = append(a, toString(u))
		}
		return a, nil
	case []int32:
		for _, u := range v {
			a = append(a, toString(u))
		}
		return a, nil
	case []int64:
		for _, u := range v {
			a = append(a, toString(u))
		}
		return a, nil
	case []float32:
		for _, u := range v {
			a = append(a, toString(u))
		}
		return a, nil
	case []float64:
		for _, u := range v {
			a = append(a, toString(u))
		}
		return a, nil
	case string:
		return strings.Fields(v), nil
	case []error:
		for _, err := range i.([]error) {
			a = append(a, err.Error())
		}
		return a, nil
	case interface{}:
		str, err := toStringE(v)
		if err != nil {
			return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
		}
		return []string{str}, nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
	}
}

// toIntSliceE casts an interface to a []int type.
func toIntSliceE(i interface{}) ([]int, error) {
	if i == nil {
		return []int{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	}

	switch v := i.(type) {
	case []int:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := toIntE(s.Index(j).Interface())
			if err != nil {
				return []int{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	}
}

// toDurationSliceE casts an interface to a []time.Duration type.
func toDurationSliceE(i interface{}) ([]time.Duration, error) {
	if i == nil {
		return []time.Duration{}, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
	}

	switch v := i.(type) {
	case []time.Duration:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]time.Duration, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := toDurationE(s.Index(j).Interface())
			if err != nil {
				return []time.Duration{}, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []time.Duration{}, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
	}
}

// stringToDate attempts to parse a string into a time.Time type using a
// predefined list of formats.  If no suitable format is found, an error is
// returned.
func stringToDate(s string) (time.Time, error) {
	return parseDateWith(s, []string{
		time.RFC3339,
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
		"2006-01-02",
		"02 Jan 2006",
		"2006-01-02T15:04:05-0700", // RFC3339 without timezone hh:mm colon
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05Z07:00", // RFC3339 without T
		"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
		"2006-01-02 15:04:05",
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	})
}

func parseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.Parse(dateType, s); e == nil {
			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}

// jsonStringToObject attempts to unmarshall a string as JSON into
// the object passed as pointer.
func jsonStringToObject(s string, v interface{}) error {
	data := []byte(s)
	return json.Unmarshal(data, v)
}

func toBool(i interface{}) bool {
	v, err := toBoolE(i)
	panicerr(err)
	return v
}

// toTime casts an interface to a time.Time type.
func toTime(i interface{}) time.Time {
	v, err := toTimeE(i)
	panicerr(err)
	return v
}

// toDuration casts an interface to a time.Duration type.
func toDuration(i interface{}) time.Duration {
	v, err := toDurationE(i)
	panicerr(err)
	return v
}

// toFloat64 casts an interface to a float64 type.
func toFloat64(i interface{}) float64 {
	v, err := toFloat64E(i)
	panicerr(err)
	return v
}

// toFloat32 casts an interface to a float32 type.
func toFloat32(i interface{}) float32 {
	v, err := toFloat32E(i)
	panicerr(err)
	return v
}

// toInt64 casts an interface to an int64 type.
func toInt64(i interface{}) int64 {
	v, err := toInt64E(i)
	panicerr(err)
	return v
}

// toInt32 casts an interface to an int32 type.
func toInt32(i interface{}) int32 {
	v, err := toInt32E(i)
	panicerr(err)
	return v
}

// toInt16 casts an interface to an int16 type.
func toInt16(i interface{}) int16 {
	v, err := toInt16E(i)
	panicerr(err)
	return v
}

// toInt8 casts an interface to an int8 type.
func toInt8(i interface{}) int8 {
	v, err := toInt8E(i)
	panicerr(err)
	return v
}

// toInt casts an interface to an int type.
func toInt(i interface{}) int {
	v, err := toIntE(i)
	panicerr(err)
	return v
}

// toUint casts an interface to a uint type.
func toUint(i interface{}) uint {
	v, err := toUintE(i)
	panicerr(err)
	return v
}

// toUint64 casts an interface to a uint64 type.
func toUint64(i interface{}) uint64 {
	v, err := toUint64E(i)
	panicerr(err)
	return v
}

// toUint32 casts an interface to a uint32 type.
func toUint32(i interface{}) uint32 {
	v, err := toUint32E(i)
	panicerr(err)
	return v
}

// toUint16 casts an interface to a uint16 type.
func toUint16(i interface{}) uint16 {
	v, err := toUint16E(i)
	panicerr(err)
	return v
}

// toUint8 casts an interface to a uint8 type.
func toUint8(i interface{}) uint8 {
	v, err := toUint8E(i)
	panicerr(err)
	return v
}

// toString casts an interface to a string type.
func toString(i interface{}) string {
	v, err := toStringE(i)
	panicerr(err)
	return v
}

func str(i interface{}) string {
	v, err := toStringE(i)
	panicerr(err)
	return v
}

// toStringMapString casts an interface to a map[string]string type.
func toStringMapString(i interface{}) map[string]string {
	v, err := toStringMapStringE(i)
	panicerr(err)
	return v
}

// toStringMapStringSlice casts an interface to a map[string][]string type.
func toStringMapStringSlice(i interface{}) map[string][]string {
	v, err := toStringMapStringSliceE(i)
	panicerr(err)
	return v
}

// toStringMapBool casts an interface to a map[string]bool type.
func toStringMapBool(i interface{}) map[string]bool {
	v, err := toStringMapBoolE(i)
	panicerr(err)
	return v
}

// toStringMapInt casts an interface to a map[string]int type.
func toStringMapInt(i interface{}) map[string]int {
	v, err := toStringMapIntE(i)
	panicerr(err)
	return v
}

// toStringMapInt64 casts an interface to a map[string]int64 type.
func toStringMapInt64(i interface{}) map[string]int64 {
	v, err := toStringMapInt64E(i)
	panicerr(err)
	return v
}

// toStringMap casts an interface to a map[string]interface{} type.
func toStringMap(i interface{}) map[string]interface{} {
	v, err := toStringMapE(i)
	panicerr(err)
	return v
}

// toSlice casts an interface to a []interface{} type.
func toSlice(i interface{}) []interface{} {
	v, err := toSliceE(i)
	panicerr(err)
	return v
}

// toBoolSlice casts an interface to a []bool type.
func toBoolSlice(i interface{}) []bool {
	v, err := toBoolSliceE(i)
	panicerr(err)
	return v
}

// toStringSlice casts an interface to a []string type.
func toStringSlice(i interface{}) []string {
	v, err := toStringSliceE(i)
	panicerr(err)
	return v
}

// toIntSlice casts an interface to a []int type.
func toIntSlice(i interface{}) []int {
	v, err := toIntSliceE(i)
	panicerr(err)
	return v
}

// toDurationSlice casts an interface to a []time.Duration type.
func toDurationSlice(i interface{}) []time.Duration {
	v, err := toDurationSliceE(i)
	panicerr(err)
	return v
}
type prometheusStruct struct {
	url string
}

type prometheusOriginalResultStruct struct {
	Status string
	Data   struct {
		ResultType string
		Result     []struct {
			Metric map[string]interface{}
			Value  []interface{}
		}
	}
}

type prometheusResultStruct struct {
	Label map[string]string
	Value float64
}

func getPrometheus(url string) *prometheusStruct {
	return &prometheusStruct{url: url + "/api/v1/query"}
}

func (m *prometheusStruct) query(query string, time ...int64) (res []prometheusResultStruct) {
	var ttime int64
	if len(time) == 0 {
		ttime = now()
	} else {
		ttime = time[0]
	}
	resp := httpPost(m.url, httpParam{"query": query, "time": ttime})
	// fmt.Println(resp.content)
	if resp.statusCode != 200 {
		panicerr("查询Prometheus出错, 查询语句: " + query + ", 状态码: " + str(resp.statusCode) + ", 错误消息:" + resp.content)
	}
	// fmt.Println(resp.content)

	var por prometheusOriginalResultStruct
	err := json.Unmarshal([]byte(resp.content), &por)
	panicerr(err)
	// lg.debug(por)
	if por.Status != "success" {
		panicerr("查询Prometheus出错, 查询语句: " + query + ", Prometheus查询结果状态: " + str(por.Status))
	}

	var pr prometheusResultStruct
	for _, i := range por.Data.Result {
		pr.Label = make(map[string]string)
		for k, v := range i.Metric {
			pr.Label[k] = str(v)
		}
		pr.Value = toFloat64(i.Value[1])
		res = append(res, pr)
	}
	return
}
func getChinaCityRank(city string) string {
	citys := []string{"上海", "北京", "深圳", "广州", "重庆", "苏州", "成都", "武汉", "杭州", "天津", "南京", "长沙", "宁波", "无锡", "青岛", "郑州", "佛山", "泉州", "东莞", "济南", "合肥", "福州", "南通", "西安", "烟台", "常州", "徐州", "大连", "唐山", "温州", "昆明", "沈阳", "厦门", "长春", "扬州", "石家庄", "绍兴", "盐城", "潍坊", "南昌", "嘉兴", "哈尔滨", "台州", "泰州", "洛阳", "襄阳", "漳州", "临沂", "金华", "南宁", "济宁", "宜昌", "惠州", "榆林", "镇江", "贵阳", "太原", "邯郸", "保定", "岳阳", "淮安", "南阳", "淄博", "常德", "芜湖", "鄂尔多斯", "沧州", "廊坊", "乌鲁木齐", "遵义", "赣州", "珠海", "菏泽", "茂名", "衡阳", "江门", "连云港", "柳州", "湖州", "九江", "许昌", "中山", "宿迁", "周口", "湛江", "德州", "株洲", "威海", "包头", "呼和浩特", "东营", "滁州", "新乡", "绵阳", "安阳", "兰州", "阜阳", "商丘", "信阳", "驻马店"}
	for idx, v := range citys {
		if v == city {
			return toString(idx + 1)
		}
	}
	return ">100"
}
// 转换字符串到时间戳
// datetime.datetime.strptime()
func strptime(format, strtime string) int64 {
	format = strings.ReplaceAll(format, "%Y", "2006")
	format = strings.ReplaceAll(format, "%m", "01")
	format = strings.ReplaceAll(format, "%d", "02")
	format = strings.ReplaceAll(format, "%H", "15")
	format = strings.ReplaceAll(format, "%M", "04")
	format = strings.ReplaceAll(format, "%S", "05")
	t, err := time.Parse(format, strtime)
	panicerr(err)

	_, offset := time.Now().Zone()
	return t.Unix() - toInt64(offset)
}

// 转换时间戳到时间字符串
// datetime.datetime.strftime()
func strftime(format string, timestamp int64) string {
	format = strings.ReplaceAll(format, "%Y", "2006")
	format = strings.ReplaceAll(format, "%m", "01")
	format = strings.ReplaceAll(format, "%d", "02")
	format = strings.ReplaceAll(format, "%H", "15")
	format = strings.ReplaceAll(format, "%M", "04")
	format = strings.ReplaceAll(format, "%S", "05")
	return time.Unix(timestamp, 0).Format(format)
}

// 休眠，秒
// time.sleep(0.1)
func sleep(t interface{}) {
	time.Sleep(getTimeDuration(t))
}

// 查找子字符串在字符串中出现的位置
// stringobj.index(substr)
func strIndex(str, substr string) int {
	pos := strings.Index(str, substr)
	return pos
}

// stringobj.replace(old, new)
func strReplace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// stringobj.upper()
func strUpper(str string) string {
	return strings.ToUpper(str)
}

// stringobj.lower()
func strLower(str string) string {
	return strings.ToLower(str)
}

// stringobj.join()
//  将一个一维数组的值转化为字符串,返回一个字符串，其内容为由 glue 分割开的数组的值。
func strJoin(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

// 去除字符串首尾处的空白字符（或者其他字符）
func strStrip(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimSpace(str)
	}
	return strings.Trim(str, characterMask[0])
}

// 删除字符串开头的空白字符（或其他字符）
func strLStrip(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimLeftFunc(str, unicode.IsSpace)
	}
	return strings.TrimLeft(str, characterMask[0])
}

// 删除字符串末端的空白字符（或者其他字符）
func strRStrip(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimRightFunc(str, unicode.IsSpace)
	}
	return strings.TrimRight(str, characterMask[0])
}

type urlComponents struct {
	schema   string
	host     string
	port     string
	user     string
	pass     string
	path     string
	query    string
	fragment string
}

// ParseURL parse_url()
// Parse a URL and return its components
// -1: all; 1: scheme; 2: host; 4: port; 8: user; 16: pass; 32: path; 64: query; 128: fragment
//  解析 URL，返回其组成部分
func urlparse(str string) *urlComponents {
	u, err := url.Parse(str)
	panicerr(err)

	pass, _ := u.User.Password()

	var port string

	if u.Port() == "" {
		if u.Scheme == "https" {
			port = "443"
		}
		if u.Scheme == "http" {
			port = "80"
		}
	} else {
		port = u.Port()
	}

	return &urlComponents{
		schema:   u.Scheme,
		host:     u.Hostname(),
		port:     port,
		user:     u.User.Username(),
		pass:     pass,
		path:     u.Path,
		query:    u.RawQuery,
		fragment: u.Fragment,
	}
}

// URLEncode urlencode()
// 编码 URL 字符串
func urlEncode(str string) string {
	return url.QueryEscape(str)
}

// URLDecode urldecode()
//  解码已编码的 URL 字符串
func urlDecode(str string) string {
	str, err := url.QueryUnescape(str)
	panicerr(err)
	return str
}

// Base64Encode base64_encode()
// 使用 MIME base64 对数据进行编码
func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64_decode()
//  对使用 MIME base64 编码的数据进行解码
func base64Decode(str string) string {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	panicerr(err)
	return string(data)
}

// Abs abs()
// 绝对值
func abs(number float64) float64 {
	return math.Abs(number)
}

// Rand rand()
// Range: [0, 2147483647]
// random.randint()
func randint(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

// FileExists file_exists()
// os.path.exists()
//  检查文件或目录是否存在
func pathExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsFile is_file()
// os.path.isfile()
//  判断给定文件名是否为一个正常的文件
func pathIsFile(filename string) bool {
	fd, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	fm := fd.Mode()
	return fm.IsRegular()
}

// IsDir is_dir()
// os.path.isdir()
//  判断给定文件名是否是一个目录
func pathIsDir(filename string) bool {
	fd, err := os.Stat(filename)
	panicerr(err)
	fm := fd.Mode()
	return fm.IsDir()
}

// Unlink unlink()
// os.unlink()
//  删除文件或者目录, 目录或者文件不存在不报错
func unlink(filename string) {
	err := os.RemoveAll(filename)
	panicerr(err)
}

// Copy copy()
// shutil.copy()
//  拷贝文件
func copy(source, dest string) {
	fd1, err := os.Open(source)
	panicerr(err)
	defer fd1.Close()
	fd2, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	panicerr(err)
	defer fd2.Close()
	_, err = io.Copy(fd2, fd1)
	panicerr(err)
}

// Rename rename()
// os.rename()
// 重命名一个文件或目录
func rename(oldname, newname string) {
	err := os.Rename(oldname, newname)
	panicerr(err)
}

// Mkdir mkdir()
// os.mkdir()
//  递归新建目录, 目录存在不报错，如果路径存在又是一个文件，则报错
func mkdir(filename string) {
	err := os.MkdirAll(filename, 0755)
	panicerr(err)
}

// Getcwd getcwd()
// os.getcwd()
//  取得当前工作目录
func getcwd() string {
	dir, err := os.Getwd()
	panicerr(err)
	return dir
}

// Basename basename()
// os.path.basename()
//  返回路径中的文件名部分
func basename(path string) string {
	return filepath.Base(path)
}

func basedir(path string) string {
	str, err := filepath.Abs(filepath.Dir(path))
	panicerr(err)
	return str
}

// IsNumeric is_numeric()
// Numeric strings consist of optional sign, any number of digits, optional decimal part and optional exponential part.
// Thus +0123.45e6 is a valid numeric value.
// In PHP hexadecimal (e.g. 0xf4c3b00c) is not supported, but IsNumeric is supported.
// 检测变量是否为数字或数字字符串
func isdigit(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.TrimSpace(str)
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9, Point, Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}

func getOutput(command string, timeoutSecond ...interface{}) string {
	_, o := getStatusOutput(command, timeoutSecond...)
	return o
}

// subprocess.getstautsoutput()
// command.getstatusoutput()
func getStatusOutput(command string, timeoutSecond ...interface{}) (int, string) {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	// remove the " and ' on both sides
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}

	if !cmdExists(parts[0]) {
		panicerr("Command not exists")
	}

	var statuscode int
	var output string
	if len(timeoutSecond) != 0 {
		t := timeoutSecond[0]
		timeoutDuration := getTimeDuration(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		// cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

		out, err := cmd.CombinedOutput()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
		output = string(out)
	} else {
		// cmd := exec.Command("/bin/sh", "-c", command)
		cmd := exec.Command(parts[0], parts[1:]...)

		out, err := cmd.CombinedOutput()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
		output = string(out)
	}
	return statuscode, output
}

// System system()
// returnVar, 0: succ; 1: fail
// Returns the last line of the command output on success, and "" on failure.
func system(command string, timeoutSecond ...interface{}) int {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	// remove the " and ' on both sides
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}

	if !cmdExists(parts[0]) {
		panicerr("Command not exists")
	}

	var statuscode int
	if len(timeoutSecond) != 0 {
		t := timeoutSecond[0]
		timeoutDuration := getTimeDuration(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		// cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command) // 如果不是bash会kill不掉
		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
	} else {
		// cmd := exec.Command("/bin/bash", "-c", command)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
	}
	return statuscode
}

func getOutputWithShell(command string, timeoutSecond ...interface{}) string {
	_, o := getStatusOutputWithShell(command, timeoutSecond...)
	return o
}

// subprocess.getstautsoutput()
// command.getstatusoutput()
func getStatusOutputWithShell(command string, timeoutSecond ...interface{}) (int, string) {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	// remove the " and ' on both sides
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}

	if !cmdExists(parts[0]) {
		panicerr("Command not exists")
	}

	var shell string
	for _, s := range []string{"/bin/bash"} {
		if pathExists(s) {
			shell = s
			break
		}
	}

	if shell == "" {
		err := errors.New("Shell not found")
		panicerr(err)
	}

	var statuscode int
	var output string
	if len(timeoutSecond) != 0 {
		t := timeoutSecond[0]
		timeoutDuration := getTimeDuration(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		cmd := exec.CommandContext(ctx, shell, "-c", command)
		//cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

		out, err := cmd.CombinedOutput()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
		output = string(out)
	} else {
		cmd := exec.Command(shell, "-c", command)

		out, err := cmd.CombinedOutput()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
		output = string(out)
	}
	return statuscode, output
}

// System system()
// returnVar, 0: succ; 1: fail
// Returns the last line of the command output on success, and "" on failure.
func systemWithShell(command string, timeoutSecond ...interface{}) int {
	q := rune(0)
	parts := strings.FieldsFunc(command, func(r rune) bool {
		switch {
		case r == q:
			q = rune(0)
			return false
		case q != rune(0):
			return false
		case unicode.In(r, unicode.Quotation_Mark):
			q = r
			return false
		default:
			return unicode.IsSpace(r)
		}
	})
	// remove the " and ' on both sides
	for i, v := range parts {
		f, l := v[0], len(v)
		if l >= 2 && (f == '"' || f == '\'') {
			parts[i] = v[1 : l-1]
		}
	}

	if !cmdExists(parts[0]) {
		panicerr("Command not exists")
	}

	var shell string
	for _, s := range []string{"/bin/bash"} {
		if pathExists(s) {
			shell = s
			break
		}
	}

	if shell == "" {
		err := errors.New("Shell not found")
		panicerr(err)
	}

	var statuscode int
	if len(timeoutSecond) != 0 {
		t := timeoutSecond[0]
		timeoutDuration := getTimeDuration(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		cmd := exec.CommandContext(ctx, shell, "-c", command)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
	} else {
		cmd := exec.Command(shell, "-c", command)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			e := err.(*exec.ExitError)
			statuscode = e.ExitCode()
		} else {
			statuscode = 0
		}
	}
	return statuscode
}

// Gethostname gethostname()
// socket.gethostname()
//  获取主机名
func gethostname() string {
	name, err := os.Hostname()
	panicerr(err)
	return name
}

// Gethostbyname gethostbyname()
// socket.gethostbyname()
func gethostbyname(hostname string, dnsserver ...string) (res []string) {
	if len(dnsserver) == 0 {
		ips, err := net.LookupIP(hostname)
		panicerr(err)
		if ips != nil {
			for _, v := range ips {
				if v.To4() != nil || v.To16() != nil {
					res = append(res, v.String())
				}
			}
		}
		return
	} else {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(8000),
				}
				if !strIn(":", dnsserver[0]) {
					dnsserver[0] = dnsserver[0] + ":53"
				}
				return d.DialContext(ctx, "udp", dnsserver[0])
			},
		}
		ips, err := r.LookupHost(context.Background(), hostname)
		panicerr(err)
		for _, ip := range ips {
			if net.ParseIP(ip) != nil {
				res = append(res, ip)
			}
		}
	}
	return
}

func getcname(hostname string, dnsserver ...string) (res string) {
	hostname = strRStrip(hostname, ".")
	var err error
	if len(dnsserver) == 0 {
		res, err = net.LookupCNAME(hostname)
		panicerr(err)
		if strRStrip(res, ".") == hostname {
			res = ""
		}
		return
	} else {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(8000),
				}
				if !strIn(":", dnsserver[0]) {
					dnsserver[0] = dnsserver[0] + ":53"
				}
				return d.DialContext(ctx, "udp", dnsserver[0])
			},
		}
		res, err = r.LookupCNAME(context.Background(), hostname)
		panicerr(err)
		lg.trace(res, hostname)
		if strRStrip(res, ".") == hostname {
			res = ""
		}
		return
	}
}

// Getenv getenv()
// os.getenv()
func getenv(varname string) string {
	e, exist := os.LookupEnv(varname)
	if !exist {
		err := errors.New("环境变量不存在")
		panicerr(err)
	}
	return e
}

func dirname(path string) string {
	return filepath.Dir(path)
}

func uuid4() string {
	return uuid.New().String()
}

func shortuuid4() string {
	return shortuuid.New()
}

func walk(path string) chan string {
	pc := make(chan string)
	go func() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if path != "." && path != ".." {
				pc <- path
			}
			return nil
		})
		if err != nil {
			close(pc)
			_, fn, line, _ := runtime.Caller(0)
			panic(filepath.Base(fn) + ":" + strconv.Itoa(line-9) + " >> " + err.Error())
		}
		close(pc)
	}()

	return pc
}

func reFindAll(pattern string, text string) [][]string {
	r, err := regexp.Compile(pattern)
	panicerr(err)
	return r.FindAllStringSubmatch(text, -1)
}

func reReplace(pattern string, newstring string, text string) string {
	r, err := regexp.Compile(pattern)
	panicerr(err)
	return r.ReplaceAllString(text, newstring)
}

func strSplit(str string, sep ...string) []string {
	var a []string
	if len(sep) != 0 {
		for _, v := range strings.Split(str, sep[0]) {
			a = append(a, strStrip(v))
		}
	} else {
		for _, v := range strings.Split(str, " ") {
			if strStrip(v) != "" {
				a = append(a, strStrip(v))
			}
		}
	}

	return a
}

func strSplitlines(str string) []string {
	var a []string
	for _, v := range strings.Split(str, "\n") {
		a = append(a, strRStrip(v, "\r"))
	}

	return a
}

func print(data ...interface{}) int {
	scheme := pp.ColorScheme{
		Bool:            pp.Cyan | pp.Bold,
		Integer:         pp.Blue | pp.Bold,
		Float:           pp.Magenta | pp.Bold,
		String:          pp.Green,
		StringQuotation: pp.Green | pp.Bold,
		EscapedChar:     pp.Magenta | pp.Bold,
		FieldName:       pp.Yellow,
		PointerAdress:   pp.Blue | pp.Bold,
		Nil:             pp.Cyan | pp.Bold,
		Time:            pp.Blue | pp.Bold,
		StructName:      pp.Green | pp.Bold,
		ObjectLength:    pp.Blue,
	}

	pp.SetColorScheme(scheme)

	num, err := pp.Println(data...)
	panicerr(err)
	return num
}

func printf(format string, data ...interface{}) int {
	scheme := pp.ColorScheme{
		Bool:            pp.Cyan | pp.Bold,
		Integer:         pp.Blue | pp.Bold,
		Float:           pp.Magenta | pp.Bold,
		String:          pp.Green,
		StringQuotation: pp.Green | pp.Bold,
		EscapedChar:     pp.Magenta | pp.Bold,
		FieldName:       pp.Yellow,
		PointerAdress:   pp.Blue | pp.Bold,
		Nil:             pp.Cyan | pp.Bold,
		Time:            pp.Blue | pp.Bold,
		StructName:      pp.Green | pp.Bold,
		ObjectLength:    pp.Blue,
	}

	pp.SetColorScheme(scheme)

	num, err := pp.Printf(format, data...)
	panicerr(err)
	return num
}

func sprint(data ...interface{}) string {
	scheme := pp.ColorScheme{
		Bool:            pp.Cyan | pp.Bold,
		Integer:         pp.Blue | pp.Bold,
		Float:           pp.Magenta | pp.Bold,
		String:          pp.Green,
		StringQuotation: pp.Green | pp.Bold,
		EscapedChar:     pp.Magenta | pp.Bold,
		FieldName:       pp.Yellow,
		PointerAdress:   pp.Blue | pp.Bold,
		Nil:             pp.Cyan | pp.Bold,
		Time:            pp.Blue | pp.Bold,
		StructName:      pp.Green | pp.Bold,
		ObjectLength:    pp.Blue,
	}

	//pp.ColoringEnabled = false
	pp.SetColorScheme(scheme)

	return strStrip(pp.Sprintln(data...))
}

func strInArr(item string, arr []string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func itemInArray(item interface{}, array interface{}) bool {
	// 获取值的列表
	arr := reflect.ValueOf(array)

	// 手工判断值的类型
	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		panicerr("Invalid data type of param \"array\": Not an Array")
	}

	// 遍历值的列表
	for i := 0; i < arr.Len(); i++ {
		// 取出值列表的元素并转换为interface
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func keyInMap(Key interface{}, Map interface{}) bool {
	arr := reflect.ValueOf(Map)
	if arr.Kind() != reflect.Map {
		panicerr("Invalid data type of param \"Map\": Not a Map")
	}
	for _, v := range arr.MapKeys() {
		if v.Interface() == Key {
			return true
		}
	}
	return false
}

func randomChoice(array interface{}) interface{} {
	rand.Seed(time.Now().Unix())
	arr := reflect.ValueOf(array)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		panicerr("Invalid data type of param \"array\": Not an Array")
	}

	return arr.Index(rand.Intn(arr.Len())).Interface()
}

type jsonMap map[string]interface{}
type jsonArr []interface{}

func jsonDumps(v interface{}, pretty ...bool) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	if len(pretty) != 0 {
		encoder.SetIndent(" ", " ")
	}
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)

	panicerr(err)
	return strStrip(buffer.String())
}

func jsonLoads(str string) jsonMap {
	var dat jsonMap
	err := json.Unmarshal([]byte(str), &dat)
	panicerr(err)
	return dat
}

func rangeInt(num ...int) []int {
	if len(num) != 1 && len(num) != 2 {
		panicerr("需要1～2个参数")
	}
	var a []int
	var start int
	var end int
	if len(num) == 1 {
		start = 0
		end = num[0]
	} else {
		start = num[0]
		end = num[1]
	}
	for i := start; i < end; i++ {
		a = append(a, i)
	}
	return a
}

func strIn(substr string, str string) bool {
	if strIndex(str, substr) != -1 {
		return true
	}
	return false
}

type lockStruct struct {
	lock *sync.Mutex
}

func getLock() *lockStruct {
	var a sync.Mutex
	return &lockStruct{lock: &a}
}

func (m *lockStruct) acquire() {
	m.lock.Lock()
}

func (m *lockStruct) release() {
	m.lock.Unlock()
}

func strStartsWith(str string, substr string) (res bool) {
	if len(substr) <= len(str) && str[:len(substr)] == substr {
		res = true
	} else {
		res = false
	}
	return
}

func strEndsWith(str string, substr string) (res bool) {
	if len(substr) <= len(str) && str[len(str)-len(substr):] == substr {
		res = true
	} else {
		res = false
	}
	return
}

func listdir(path string) (res []string) {
	files, err := ioutil.ReadDir(path)
	panicerr(err)

	for _, f := range files {
		res = append(res, f.Name())
	}
	return
}

func pathJoin(args ...string) string {
	return path.Join(args...)
}

func strCount(s string, substr string) int {
	return strings.Count(s, substr)
}

// Realpath abspath()
//  返回规范化的绝对路径名
func abspath(path string) string {
	str, err := filepath.Abs(path)
	panicerr(err)
	return str
}

func pathIsSymlink(path string) bool {
	fi, err := os.Lstat(path)
	panicerr(err)
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	} else {
		return false
	}
}
type progressBarStruct struct {
	bar *progressbar.ProgressBar
}

func getProgressBar(title string, total int64, showBytes ...bool) *progressBarStruct {
	var showBytesOption progressbar.Option
	if len(showBytes) != 0 && showBytes[0] == true {
		showBytesOption = progressbar.OptionShowBytes(true)
	} else {
		showBytesOption = progressbar.OptionShowBytes(false)
	}
	bar := progressbar.NewOptions64(total,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		showBytesOption,
		progressbar.OptionSetDescription("[cyan]*[reset] "+title),

		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]=[reset]",
			SaucerHead:    "[cyan]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return &progressBarStruct{
		bar: bar,
	}
}

func (m *progressBarStruct) add(num int64) {
	m.bar.Add64(num)
}

func (m *progressBarStruct) set(num int64) {
	m.bar.Set64(num)
}

func (m *progressBarStruct) setTotal(total int64) {
	m.bar.ChangeMax64(total)
}

func (m *progressBarStruct) clear() {
	m.bar.Clear()
}
// UDP - Client

type udpClientSideConn struct {
	isclose bool
	conn    net.Conn
}

func udpConnect(host string, port int) udpClientSideConn {
	conn, err := net.Dial("udp", host+":"+toString(port))
	panicerr(err)
	return udpClientSideConn{conn: conn}
}

func (m *udpClientSideConn) send(str string) {
	_, err := fmt.Fprintf(m.conn, str)
	panicerr(err)
}

func (m *udpClientSideConn) close() {
	if !m.isclose {
		m.isclose = true
		m.conn.Close()
	}
}

func (m *udpClientSideConn) recv(buffersize int) string {
	p := make([]byte, buffersize)
	n, err := bufio.NewReader(m.conn).Read(p)
	panicerr(err)
	return string(p[:n])
}

// UDP - Server

type udpServerSideConn struct {
	isclose bool
	conn    *net.UDPConn
}

func udpListen(host string, port int) udpServerSideConn {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	}
	ser, err := net.ListenUDP("udp", &addr)
	panicerr(err)
	return udpServerSideConn{conn: ser}
}

func (m *udpServerSideConn) recvfrom(buffersize int, timeout ...int) (string, *net.UDPAddr) {
	if len(timeout) != 0 {
		m.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	p := make([]byte, buffersize)
	n, remoteaddr, err := m.conn.ReadFromUDP(p)
	panicerr(err)
	return string(p[:n]), remoteaddr
}

func (m *udpServerSideConn) sendto(data string, address *net.UDPAddr, timeout ...int) {
	if len(timeout) != 0 {
		m.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	_, err := m.conn.WriteToUDP([]byte(data), address)
	panicerr(err)
}

func (m *udpServerSideConn) close() {
	if !m.isclose {
		m.isclose = true
		m.conn.Close()
	}
}
type fsnotifyFileEventStruct struct {
	action string
	path   string
}

func inotify(path string) chan *fsnotifyFileEventStruct {
	ch := make(chan *fsnotifyFileEventStruct)
	watchList := make(map[string]struct{})

	watcher, err := fsnotify.NewWatcher()
	panicerr(err)

	go func() {
		// 1秒内有同样的事件在同样的路径, 则忽略
		cache := ttlcache.NewCache(time.Second)

		for {
			select {
			case ev := <-watcher.Events:
				_, exists := cache.Get(ev.String() + abspath(ev.Name))
				if !exists {
					cache.Set(ev.String()+abspath(ev.Name), "")
					var action string
					if ev.Op == fsnotify.Create {
						action = "create"
						if pathIsDir(ev.Name) {
							err := watcher.Add(ev.Name)
							panicerr(err)
							watchList[ev.Name] = struct{}{}
						}
					} else if ev.Op == fsnotify.Chmod {
						action = "chmod"
					} else if ev.Op == fsnotify.Remove {
						action = "remove"
						delete(watchList, ev.Name)
						go func() {
							// 如果action是delete, 那么就会从监控的列表当中移除
							// 如果10秒内再次出现, 就再次添加监听
							// 原因是vim编辑的时候并不是直接编辑, 而是先写入.swp文件, 在保存的时候删除源文件, 然后改名.swp文件, 会生成删除源文件的action
							for range rangeInt(100) {
								if pathExists(ev.Name) {
									err = watcher.Add(ev.Name)
									panicerr(err)
									watchList[ev.Name] = struct{}{}
									break
								}
								sleep(0.1)
							}
							if !keyInMap(ev.Name, watchList) {
								ch <- &fsnotifyFileEventStruct{
									action: action,
									path:   abspath(ev.Name),
								}
							}
							if len(watchList) == 0 {
								watcher.Close()
								close(ch)
							}

						}()
						continue
					} else if ev.Op == fsnotify.Write {
						action = "write"
					} else if ev.Op == fsnotify.Rename {
						action = "rename"
						delete(watchList, ev.Name)
					}

					ch <- &fsnotifyFileEventStruct{
						action: action,
						path:   abspath(ev.Name),
					}
				}
			case err := <-watcher.Errors:
				panicerr(err)
			}
		}
	}()

	err = watcher.Add(path)
	panicerr(err)
	watchList[path] = struct{}{}

	if pathIsDir(path) {
		for p := range walk(path) {
			if pathIsDir(p) {
				err := watcher.Add(p)
				panicerr(err)
				watchList[p] = struct{}{}
			}
		}
	}

	return ch
}
type xpathJsonStruct struct {
	doc *jsonquery.Node
}

func getXPathJson(jsonstr string) *xpathJsonStruct {
	doc, err := jsonquery.Parse(strings.NewReader(jsonstr))
	panicerr(err)
	return &xpathJsonStruct{
		doc: doc,
	}
}

func (m *xpathJsonStruct) first(expr string) (res *xpathJsonStruct) {
	return &xpathJsonStruct{
		doc: jsonquery.FindOne(m.doc, expr),
	}
}

func (m *xpathJsonStruct) find(expr string) (res []*xpathJsonStruct) {
	for _, doc := range jsonquery.Find(m.doc, expr) {
		res = append(res, &xpathJsonStruct{doc: doc})
	}
	return
}

func (m *xpathJsonStruct) text() string {
	return m.doc.InnerText()
}
// TCP - Server

type tcpServerSideConn struct {
	conn    net.Conn
	isclose bool
}

type tcpServerSideListener struct {
	listener net.Listener
	isclose  bool
}

func tcpListen(host string, port int) *tcpServerSideListener {
	l, err := net.Listen("tcp", host+":"+toString(port))
	panicerr(err)

	return &tcpServerSideListener{listener: l}
}

func (m *tcpServerSideListener) accept() chan *tcpServerSideConn {
	ch := make(chan *tcpServerSideConn)

	go func() {
		for {
			c, err := m.listener.Accept()
			if err != nil {
				if strIn("use of closed network connection", err.Error()) {
					break
				}
				panicerr(err)
			}
			ct := &tcpServerSideConn{conn: c, isclose: false}
			ch <- ct
		}
	}()

	return ch
}

func (m *tcpServerSideListener) close() {
	if !m.isclose {
		m.isclose = true
		m.listener.Close()
	}

}

func (m *tcpServerSideConn) close() {
	if !m.isclose {
		m.isclose = true
		m.conn.Close()
	}
}

func (m *tcpServerSideConn) send(str string) {
	_, err := m.conn.Write([]byte(str))
	panicerr(err)
}

func (m *tcpServerSideConn) recv(buffersize int) string {
	reply := make([]byte, buffersize)
	n, err := m.conn.Read(reply)
	panicerr(err)
	return string(reply[:n])
}

// TCP - Client

type tcpClientSideConn struct {
	conn    *net.TCPConn
	isclose bool
}

func tcpConnect(host string, port int, timeout ...int) *tcpClientSideConn {
	servAddr := host + ":" + toString(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	panicerr(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	panicerr(err)
	return &tcpClientSideConn{conn: conn, isclose: false}
}

func (m *tcpClientSideConn) send(str string, timeout ...int) {
	if len(timeout) != 0 {
		m.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	_, err := m.conn.Write([]byte(str))
	panicerr(err)
}

func (m *tcpClientSideConn) recv(buffersize int, timeout ...int) string {
	if len(timeout) != 0 {
		m.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	reply := make([]byte, buffersize)
	n, err := m.conn.Read(reply)
	panicerr(err)
	return string(reply[:n])
}

func (m *tcpClientSideConn) close() {
	if !m.isclose {
		m.isclose = true
		err := m.conn.Close()
		panicerr(err)
	}
}
func getSysInfo() sysinfo.SysInfo {
	var si sysinfo.SysInfo
	si.GetSysInfo()
	return si
}
func parseUserAgent(UserAgent string) ua.UserAgent {
	return ua.Parse(UserAgent)
}
type crontabStruct struct {
	c *crontab.Crontab
}

func getCrontab() *crontabStruct {
	return &crontabStruct{c: crontab.New()}
}

func (m *crontabStruct) add(schedule string, fn interface{}, args ...interface{}) {
	err := m.c.AddJob(schedule, fn, args...)
	panicerr(err)
}

func (m *crontabStruct) destory() {
	m.c.Shutdown()
}
func tailf(path string, startFromEndOfFile ...bool) chan *tail.Line {
	if !pathExists(path) {
		panicerr("Cannot open file \"" + path + "\": no such file or directory")
	}
	var t *tail.Tail
	var err error

	if len(startFromEndOfFile) == 0 {
		t, err = tail.TailFile(path, tail.Config{Follow: true, Poll: true, Logger: tail.DiscardingLogger})
	} else {
		if startFromEndOfFile[0] == false {
			t, err = tail.TailFile(path, tail.Config{Follow: true, Poll: true, Logger: tail.DiscardingLogger})
		} else {
			t, err = tail.TailFile(path, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}, Poll: true, Logger: tail.DiscardingLogger})
		}
	}

	panicerr(err)

	return t.Lines
}
type telegraphStruct struct {
	account *telegraph.Account
}

type telegraphPageInfo struct {
	author  string
	title   string
	url     string
	content string
}

func getTelegraph(AuthorName string) *telegraphStruct {
	account, err := telegraph.CreateAccount(telegraph.Account{
		ShortName:  randomStr(6, "abcdefghijklmn1234567890"),
		AuthorName: AuthorName,
	})
	panicerr(err)
	return &telegraphStruct{
		account: account,
	}
}

func (m *telegraphStruct) post(title string, content string) *telegraphPageInfo {
	tcontent, err := telegraph.ContentFormat(content)
	panicerr(err)

	page, err := m.account.CreatePage(telegraph.Page{
		URL:        "this-is-a-test-url",
		Title:      title,
		AuthorName: m.account.AuthorName,
		Content:    tcontent,
	}, true)
	panicerr(err)
	return &telegraphPageInfo{
		author:  page.AuthorName,
		title:   page.Title,
		url:     page.URL,
		content: content,
	}
}
type RedisStruct struct {
	ctx                    context.Context
	rdb                    *redis.Client
	networkErrorRetryTimes int
}

type redisConfig struct {
	networkErrorRetryTimes int
}

// 用来过滤报错的信息， 如果包含有如下的某一个， 就判断为是网络错误
var redisNetworkErrorStrings = []string{
	"timeout",
	"connection reset by peer",
	"connection refused",
}

func getRedis(host string, port int, password string, db int, cfg ...redisConfig) *RedisStruct {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + toString(port),
		Password: password,
		DB:       db,
	})
	r := &RedisStruct{ctx: context.Background(), rdb: rdb}
	r.ping()

	if len(cfg) != 0 {
		r.networkErrorRetryTimes = cfg[0].networkErrorRetryTimes
	}

	return r
}

func (m *RedisStruct) ping() string {
	pong, err := m.rdb.Ping(m.ctx).Result()
	panicerr(err)
	return pong
}

func (m *RedisStruct) del(key string) {
	errortimes := 0
	var err error
	for {
		err = m.rdb.Del(m.ctx, key).Err()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(redisNetworkErrorStrings, err.Error()) && errortimes < m.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}
}

func (m *RedisStruct) set(key string, value string, ttl ...interface{}) {
	var t time.Duration

	if len(ttl) != 0 {
		if typeof(ttl[0]) == "float64" {
			tt := ttl[0].(float64) * 1000
			t = time.Duration(tt) * time.Millisecond
		}
		if typeof(ttl[0]) == "int" {
			tt := ttl[0].(int)
			t = time.Duration(tt) * time.Second
		}
	}

	errortimes := 0
	var err error
	for {
		err = m.rdb.Set(m.ctx, key, value, t).Err()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(redisNetworkErrorStrings, err.Error()) && errortimes < m.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}
}

func (m *RedisStruct) get(key string) *string {
	errortimes := 0
	var val string
	var err error
	for {
		val, err = m.rdb.Get(m.ctx, key).Result()
		if err != nil && err != redis.Nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(redisNetworkErrorStrings, err.Error()) && errortimes < m.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	if err == redis.Nil {
		return nil
	}

	return &val
}

type RedisLockStruct struct {
	redis      *RedisStruct
	key        string
	timeoutsec int
}

var redisLockMutex sync.Mutex

func (m *RedisStruct) getLock(key string, timeoutsec int) *RedisLockStruct {
	return &RedisLockStruct{
		redis:      m,
		key:        key,
		timeoutsec: timeoutsec, // 锁的超时时间, 为了防止进程崩溃而没有释放锁, 不是获取锁的超时时间
	}
}

func (m *RedisLockStruct) acquire() {
	redisLockMutex.Lock()
	defer redisLockMutex.Unlock()

	for {
		b, err := m.redis.rdb.SetNX(m.redis.ctx, m.key, 1, getTimeDuration(m.timeoutsec)).Result()
		panicerr(err)
		if b {
			return
		} else {
			sleep(0.1)
		}
	}
}

func (m *RedisLockStruct) release() {
	_, err := m.redis.rdb.Del(m.redis.ctx, m.key).Result()
	panicerr(err)
}
type xpathStruct struct {
	doc *html.Node
}

func getXPath(htmlString string) *xpathStruct {
	doc, err := htmlquery.Parse(strings.NewReader(htmlString))
	panicerr(err)
	return &xpathStruct{
		doc: doc,
	}
}

func (m *xpathStruct) first(expr string) (res *xpathStruct) {
	return &xpathStruct{
		doc: htmlquery.FindOne(m.doc, expr),
	}
}

func (m *xpathStruct) find(expr string) (res []*xpathStruct) {
	for _, doc := range htmlquery.Find(m.doc, expr) {
		res = append(res, &xpathStruct{doc: doc})
	}
	return
}

func (m *xpathStruct) text() string {
	return htmlquery.InnerText(m.doc)
}

func (m *xpathStruct) getAttr(attr string) string {
	return htmlquery.SelectAttr(m.doc, attr)
}

func (m *xpathStruct) html() string {
	return htmlquery.OutputHTML(m.doc, true)
}
func getGin() *gin.Engine {
	ginEng := gin.New()
	logLevel := "trace"

	ginEng.Use(func(c *gin.Context) {
		if logLevel == "trace" {
			nowTime := time.Now()
			c.Next()

			ipaddr := c.ClientIP()
			var ipinfo *ipLocationInfo
			try(func() {
				ipinfo = getIPLocation(ipaddr)
			}).except(func(e error) {
				ipinfo = &ipLocationInfo{}
			})

			ipinfostr := ipinfo.Country
			if !itemInArray(ipinfo.Region, []string{"", "N/A", ipinfo.Country}) {
				ipinfostr = ipinfostr + " " + ipinfo.Region
			}
			if ipinfo.City != "" && ipinfo.City != "N/A" {
				ipinfostr = ipinfostr + " " + ipinfo.City
			}

			ua := parseUserAgent(c.Request.Header["User-Agent"][0])

			tm := toString(time.Since(nowTime))
			tmf := fmt.Sprintf("%.2f", toFloat64(reFindAll("[0-9\\.]+", tm)[0][0]))
			tms := reFindAll("[a-zµ]+", tm)[0][0]
			tm = tmf + tms

			logstr := toString(c.Writer.Status()) + " " + fmt.Sprintf("%8v", tm) + " " + c.Request.Method + " " + c.Request.URL.String() + " " + ipaddr + " " + ipinfostr + " " + ua.OS + " " + ua.OSVersion + " " + ua.Name + " " + ua.Version + " " + ua.Device
			lg.info(logstr)
		} else {
			c.Next()
		}
	})
	return ginEng
}
type databaseStruct struct {
	engin                  *gorose.Engin
	driver                 string
	dsn                    string
	isclose                bool
	networkErrorRetryTimes int // 网络错误重试次数
}

type databaseConfig struct {
	timeout                int
	charset                string
	networkErrorRetryTimes int // 网络错误重试次数
}

// 用来过滤报错的信息， 如果包含有如下的某一个， 就判断为是网络错误
var databaseNetworkErrorStrings = []string{
	"timeout",
	"invalid connection",
	"no such host",
	"connection refused",
	"bad connection",
}

func getMySQL(host string, port int, user string, password string, db string, cfg ...databaseConfig) *databaseStruct {
	var timeoutt int
	var chartsett string
	var networkErrorRetryTimess int

	if len(cfg) != 0 {
		if cfg[0].timeout != 0 {
			timeoutt = cfg[0].timeout
		} else {
			timeoutt = 10
		}
		if cfg[0].charset != "" {
			chartsett = cfg[0].charset
		} else {
			chartsett = "utf8mb4"
		}
		if cfg[0].networkErrorRetryTimes != 0 {
			networkErrorRetryTimess = cfg[0].networkErrorRetryTimes
		} else {
			networkErrorRetryTimess = 10
		}
	} else {
		timeoutt = 10
		chartsett = "utf8mb4"
		networkErrorRetryTimess = 0
	}

	m := &databaseStruct{}
	m.networkErrorRetryTimes = networkErrorRetryTimess
	m.driver = "mysql"
	m.dsn = user + ":" + password + "@tcp(" + host + ":" + toString(port) + ")/" + db + "?timeout=" + str(timeoutt) + "s&readTimeout=" + str(timeoutt) + "s&writeTimeout=" + str(timeoutt) + "s&charset=" + chartsett
	var config = &gorose.Config{
		Driver: m.driver,
		Dsn:    m.dsn,
	}

	errortimes := 0
	var err error
	var engin *gorose.Engin
	for {
		engin, err = gorose.Open(config)
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.engin = engin

	// 它会重连，如果连接坏了
	go func(m *databaseStruct) {
		for {
			sleep(3)
			try(func() {
				m.engin.Ping()
			})
			if m.isclose {
				break
			}
		}
	}(m)

	return m
}

func getSQLite(dbpath string) *databaseStruct {
	m := &databaseStruct{}
	m.driver = "sqlite3"
	m.dsn = dbpath
	var config = &gorose.Config{
		Driver: m.driver,
		Dsn:    m.dsn,
	}

	engin, err := gorose.Open(config)
	panicerr(err)

	m.engin = engin

	return m
}

func (m *databaseStruct) query(sql string, args ...interface{}) []gorose.Data {
	db := m.engin.NewOrm()
	res, err := db.Query(sql, args...)
	panicerr(err)

	for idx, d := range res {
		for k, v := range d {
			if v != nil && typeof(v) == "time.Time" {
				res[idx][k] = strftime("%Y-%m-%d %H:%M:%S", v.(time.Time).Unix())
			}
		}
	}
	return res
}

func (m *databaseStruct) close() {
	m.isclose = true
	m.engin.GetQueryDB().Close()
}

func (m *databaseStruct) execute(sql string) int64 {
	db := m.engin.NewOrm()
	res, err := db.Execute(sql)
	panicerr(err)
	return res
}

type databaseOrmStruct struct {
	orm    gorose.IOrm
	db     *databaseStruct
	driver string
	table  string
	lock   *lockStruct
	lockby int64
}

func (m *databaseStruct) table(tbname string) *databaseOrmStruct {
	orm := m.engin.NewOrm()
	return &databaseOrmStruct{
		orm:    orm.Table("`" + tbname + "`"),
		driver: m.driver,
		table:  tbname,
		db:     m,
		lock:   getLock(), // 为了保证线程安全，链式操作当中要上锁，返回数据解锁。为了保证在多个线程中复用同一个databaseStruct的时候报错。
		lockby: -1,
	}
}

func (m *databaseStruct) renameTable(oldTableName string, newTableNname string) {
	if m.driver == "mysql" {
		m.execute("RENAME TABLE `" + oldTableName + "` TO `" + newTableNname + "`;")
	} else if m.driver == "sqlite3" {
		m.execute("ALTER TABLE `" + oldTableName + "` RENAME TO `" + newTableNname + "`;")
	}
}

func (m *databaseOrmStruct) fields(items ...string) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	var i []string
	for _, v := range items {
		i = append(i, "`"+v+"`")
	}
	m.orm.Fields(i...)
	return m
}

func (m *databaseOrmStruct) where(key string, operator string, value interface{}) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.Where(key, operator, value)
	return m
}

func (m *databaseOrmStruct) whereIn(key string, value []interface{}) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.WhereIn(key, value)
	return m
}

func (m *databaseOrmStruct) whereNull(columnName string) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.WhereNull(columnName)
	return m
}

func (m *databaseOrmStruct) whereNotNull(columnName string) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.WhereNotNull(columnName)
	return m
}

func (m *databaseOrmStruct) orWhere(key string, operator string, value interface{}) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.OrWhere(key, operator, value)
	return m
}

func (m *databaseOrmStruct) orWhereIn(key string, value []interface{}) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.OrWhereIn(key, value)
	return m
}

func (m *databaseOrmStruct) orderby(item ...string) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.OrderBy(strJoin(" ", item))
	return m
}

func (m *databaseOrmStruct) limit(number int) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.Limit(number)
	return m
}

func (m *databaseOrmStruct) get() (res []gorose.Data) {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}

	errortimes := 0
	var err error
	for {
		res, err = m.orm.Get()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.orm = m.db.table(m.table).orm

	//print(m)

	m.lock.release()
	m.lockby = -1

	for idx, d := range res {
		for k, v := range d {
			if v != nil && typeof(v) == "time.Time" {
				res[idx][k] = strftime("%Y-%m-%d %H:%M:%S", v.(time.Time).Unix())
			}
		}
	}

	return
}

func (m *databaseOrmStruct) paginate(pagesize int, page int) []gorose.Data {
	offset := pagesize * (page - 1)
	return m.limit(pagesize).offset(offset).get()
}

func (m *databaseOrmStruct) first() (res gorose.Data) {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}

	errortimes := 0
	var err error
	for {
		res, err = m.orm.First()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.orm = m.db.table(m.table).orm

	m.lock.release()
	m.lockby = -1

	for k, v := range res {
		if v != nil && typeof(v) == "time.Time" {
			res[k] = strftime("%Y-%m-%d %H:%M:%S", v.(time.Time).Unix())
		}
	}

	return
}

func (m *databaseOrmStruct) find(id int) gorose.Data {
	return m.where("id", "=", id).first()
}

func (m *databaseOrmStruct) count() (res int64) {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}

	errortimes := 0
	var err error
	for {
		res, err = m.orm.Count()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.orm = m.db.table(m.table).orm

	m.lock.release()
	m.lockby = -1

	return
}

func (m *databaseOrmStruct) exists() (res bool) {
	if m.count() == 0 {
		res = false
	} else {
		res = true
	}

	return
}

func (m *databaseOrmStruct) chunk(limit int, callback func([]gorose.Data) error) {
	err := m.orm.Chunk(limit, callback)
	panicerr(err)
}

func (m *databaseOrmStruct) buildSQL() (string, []interface{}) {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	sql, param, err := m.orm.BuildSql()
	panicerr(err)

	m.lock.release()
	m.lockby = -1

	return sql, param
}

func (m *databaseOrmStruct) data(data interface{}) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.Data(data)
	return m
}

func (m *databaseOrmStruct) offset(offset int) *databaseOrmStruct {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	m.orm.Offset(offset)
	return m
}

func (m *databaseOrmStruct) insertGetID() (num int64) {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	errortimes := 0
	var err error
	for {
		num, err = m.orm.InsertGetId()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}

	m.orm = m.db.table(m.table).orm

	m.lock.release()
	m.lockby = -1

	return
}

func (m *databaseOrmStruct) insert() (num int64) {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}

	errortimes := 0
	var err error
	for {
		num, err = m.orm.Insert()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}
	m.orm = m.db.table(m.table).orm

	m.lock.release()
	m.lockby = -1

	return
}

func (m *databaseOrmStruct) update(data ...interface{}) (num int64) {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}
	errortimes := 0
	var err error
	for {
		num, err = m.orm.Update(data...)
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}
	m.orm = m.db.table(m.table).orm

	m.lock.release()
	m.lockby = -1

	return
}

func (m *databaseOrmStruct) delete() (num int64) {
	if m.lockby != getGoroutineID() {
		m.lock.acquire()
		m.lockby = getGoroutineID()
	}

	errortimes := 0
	var err error
	for {
		num, err = m.orm.Delete()
		if err != nil {
			if func(errfilter []string, errmsg string) bool {
				for _, err := range errfilter {
					if strIn(err, errmsg) {
						return true
					}
				}
				return false
			}(databaseNetworkErrorStrings, err.Error()) && errortimes < m.db.networkErrorRetryTimes {
				errortimes += 1
				sleep(3)
			} else {
				panicerr(err)
			}
		} else {
			break
		}
	}
	m.orm = m.db.table(m.table).orm

	m.lock.release()
	m.lockby = -1

	return
}

func (m *databaseStruct) tables() (res []string) {
	if m.driver == "mysql" {
		for _, v := range m.query("show tables;") {
			for _, i := range v {
				res = append(res, toString(i))
			}
		}
	} else if m.driver == "sqlite3" {
		for _, i := range m.query("SELECT `name` FROM sqlite_master WHERE type='table';") {
			res = append(res, toString(i["name"]))
		}
	}
	return
}

func (m *databaseStruct) createTable(tableName string, engineName ...string) *databaseOrmStruct {
	if !itemInArray(tableName, m.tables()) {
		if len(engineName) != 0 && m.driver != "mysql" {
			panicerr("SQLite不支持设定存储引擎")
		}
		if m.driver == "mysql" {
			if len(engineName) != 0 {
				m.execute("CREATE TABLE `" + tableName + "`(`id` BIGINT UNSIGNED AUTO_INCREMENT,PRIMARY KEY ( `id` ))ENGINE=" + engineName[0] + " DEFAULT CHARSET=utf8mb4;")
			} else {
				m.execute("CREATE TABLE `" + tableName + "`(`id` BIGINT UNSIGNED AUTO_INCREMENT,PRIMARY KEY ( `id` ))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
			}
		} else {
			m.execute("CREATE TABLE `" + tableName + "` (`id` INTEGER PRIMARY KEY AUTOINCREMENT)")
		}
	}
	return m.table(tableName)
}

func (m *databaseOrmStruct) dropTable() int64 {
	return m.db.execute("DROP TABLE `" + m.table + "`")
}

func (m *databaseOrmStruct) truncateTable() (status int64) {
	if m.driver == "mysql" {
		status = m.db.execute("TRUNCATE TABLE `" + m.table + "`")
	} else {
		status = m.db.execute("DELETE FROM TABLE `" + m.table + "`")
	}
	return
}

func (m *databaseOrmStruct) addColumn(columnName string, columnType string, defaultValue ...string) *databaseOrmStruct {
	if !keyInMap(columnName, m.columns()) {
		if !itemInArray(columnType, []string{"int", "float", "string", "text", "datetime", "blob"}) {
			err := errors.New("只支持以下数据类型：\"int\", \"float\", \"string\", \"text\", \"datetime\", \"blob\"(SQLite支持, MySql不支持)")
			panicerr(err)
		}
		if m.driver == "mysql" {
			if columnType == "int" {
				columnType = "bigint"
			} else if columnType == "float" {
				columnType = "double"
			} else if columnType == "string" {
				columnType = "VARCHAR(512)"
			} else if columnType == "text" {
				columnType = "LONGTEXT"
			} else if columnType == "datetime" {
				columnType = "DATETIME"
			} else if columnType == "blob" {
				columnType = "LONGBLOB"
				panicerr("MySQL暂不支持blob")
			}
		} else {
			if columnType == "int" {
				columnType = "INTEGER"
			} else if columnType == "float" {
				columnType = "REAL"
			} else if columnType == "string" {
				columnType = "VARCHAR"
			} else if columnType == "text" {
				columnType = "LONGTEXT"
			} else if columnType == "datetime" {
				columnType = "DATETIME"
			} else if columnType == "blob" {
				columnType = "BLOB"
			}
		}

		var sql string
		if len(defaultValue) == 0 {
			sql = "ALTER TABLE `" + m.table + "` ADD `" + columnName + "` " + columnType + ";"
		} else {
			sql = "ALTER TABLE `" + m.table + "` ADD `" + columnName + "` " + columnType + " DEFAULT \"" + defaultValue[0] + "\";"
		}

		m.db.execute(sql)

	}
	return m
}

func (m *databaseOrmStruct) dropColumn(columnName string) *databaseOrmStruct {
	if keyInMap(columnName, m.columns()) {
		if m.driver == "mysql" {
			m.db.execute("ALTER TABLE `" + m.table + "`  DROP " + columnName)
		} else {
			panic(newerr("SQLite does not support drop column"))
		}
	}
	return m
}

func (m *databaseOrmStruct) addIndex(columnName ...string) *databaseOrmStruct {
	if !m.indexExists(columnName...) {
		columns := "`" + strJoin("`,`", columnName) + "`"
		indexName := "idx_" + strJoin("_", columnName)
		if m.driver == "mysql" {
			m.db.execute("ALTER TABLE `" + m.table + "` ADD INDEX " + indexName + "(" + columns + ")")
		} else {
			m.db.execute("CREATE INDEX " + indexName + " on `" + m.table + "` (" + columns + ");")
		}
	}
	return m
}

func (m *databaseOrmStruct) indexExists(columnName ...string) (exists bool) {
	indexName := "idx_" + strJoin("_", columnName)
	if m.driver == "mysql" {
		for _, v := range m.db.query("SHOW INDEX FROM `" + m.table + "`") {
			if v["Key_name"] == indexName {
				exists = true
			}
		}
	} else if m.driver == "sqlite3" {
		if toInt(m.db.query("SELECT count(name) as `count` FROM sqlite_master WHERE type='index' AND name='" + indexName + "';")[0]["count"]) == 1 {
			exists = true
		}
	}
	return
}

func (m *databaseOrmStruct) dropIndex(columnName ...string) *databaseOrmStruct {
	indexName := "idx_" + strJoin("_", columnName)
	if m.driver == "mysql" {
		m.db.execute("ALTER TABLE `" + m.table + "` DROP INDEX " + indexName)
	} else {
		m.db.execute("DROP INDEX " + indexName)
	}
	return m
}

func (m *databaseOrmStruct) columns() (res map[string]string) {
	res = make(map[string]string)

	if m.driver == "mysql" {
		for _, i := range m.db.query("SHOW COLUMNS FROM `" + m.table + "`;") {
			// lg.debug(i)
			if strLower(toString(i["Type"])) == "bigint(20)" {
				res[toString(i["Field"])] = "int"
			} else if strLower(toString(i["Type"])) == "double" {
				res[toString(i["Field"])] = "float"
			} else if strLower(toString(i["Type"])) == "varchar(512)" {
				res[toString(i["Field"])] = "string"
			} else if strLower(toString(i["Type"])) == "longtext" {
				res[toString(i["Field"])] = "text"
			} else if strLower(toString(i["Type"])) == "datetime" {
				res[toString(i["Field"])] = "datetime"
			} else if strLower(toString(i["Type"])) == "longblob" {
				res[toString(i["Field"])] = "blob"
			}
		}
	} else {
		for _, i := range m.db.query("PRAGMA table_info(`" + m.table + "`);") {
			if strUpper(toString(i["type"])) == "INTEGER" {
				res[toString(i["name"])] = "int"
			} else if strUpper(toString(i["type"])) == "REAL" {
				res[toString(i["name"])] = "float"
			} else if strUpper(toString(i["type"])) == "VARCHAR" {
				res[toString(i["name"])] = "string"
			} else if strUpper(toString(i["type"])) == "LONGTEXT" {
				res[toString(i["name"])] = "text"
			} else if strUpper(toString(i["type"])) == "DATETIME" {
				res[toString(i["name"])] = "datetime"
			} else if strUpper(toString(i["type"])) == "BLOB" {
				res[toString(i["name"])] = "blob"
			}
		}
	}
	return
}
type totpStruct struct {
	key string
}

func getTotp(key string) *totpStruct {
	return &totpStruct{key: key}
}

func (m *totpStruct) validate(pass string) bool {
	return pquernatotp.Validate(pass, m.key)
}

func (m *totpStruct) password() string {
	p, err := pquernatotp.GenerateCode(m.key, time.Now())
	panicerr(err)
	return p
}
func lzmaCompress(text string) string {
	var buf bytes.Buffer

	defer buf.Reset()

	w, err := xz.NewWriter(&buf)
	panicerr(err)

	_, err = io.WriteString(w, text)
	panicerr(err)

	err = w.Close()
	panicerr(err)

	return buf.String()
}

func lzmaDecompress(text string) string {
	var buf bytes.Buffer
	buf.Write([]byte(text))

	defer buf.Reset()

	r, err := xz.NewReader(&buf)
	panicerr(err)

	dbuf := new(strings.Builder)
	_, err = io.Copy(dbuf, r)
	panicerr(err)
	return dbuf.String()
}

func zlibCompress(text string) string {
	var buf bytes.Buffer

	w := zlib.NewWriter(&buf)
	_, err := w.Write([]byte(text))
	panicerr(err)

	err = w.Close()
	panicerr(err)

	return buf.String()
}

func zlibDecompress(text string) string {
	var buf bytes.Buffer
	buf.Write([]byte(text))

	defer buf.Reset()

	r, err := zlib.NewReader(&buf)
	panicerr(err)

	dbuf := new(strings.Builder)
	_, err = io.Copy(dbuf, r)
	panicerr(err)

	return dbuf.String()
}
type logStruct struct {
	level                    []string
	levelString              string
	json                     bool
	color                    bool
	logDir                   string
	logFileName              string
	logFileSuffix            string
	fd                       *fileStruct
	displayOnTerminal        bool
	lock                     *lockStruct
	logfiles                 []string
	maxlogfiles              int
	logFileSizeInMB          int
	currentLogFileSizeInByte int
	currentLogFileNumber     int
}

func getLogger() *logStruct {
	return &logStruct{
		level:                    []string{"TRAC", "DEBU", "INFO", "WARN", "ERRO"},
		color:                    true,
		displayOnTerminal:        true,
		lock:                     getLock(),
		logFileSizeInMB:          0,
		currentLogFileSizeInByte: 0,
		currentLogFileNumber:     0,
	}
}

func (m *logStruct) setLevel(level string) {
	if level == "trace" {
		m.level = []string{"TRAC", "DEBU", "INFO", "WARN", "ERRO"}
	} else if level == "debug" {
		m.level = []string{"DEBU", "INFO", "WARN", "ERRO"}
	} else if level == "info" {
		m.level = []string{"INFO", "WARN", "ERRO"}
	} else if level == "warn" {
		m.level = []string{"WARN", "ERRO"}
	} else if level == "error" {
		m.level = []string{"ERRO"}
	} else if level == "" {
		m.level = []string{}
	}
	m.levelString = level
}

func (m *logStruct) getLevel() string {
	return m.levelString
}

func (m *logStruct) setLogFile(path string, maxLogFileCount int, logFileSizeInMB ...int) {
	m.logDir = basedir(path)
	if !pathExists(m.logDir) {
		mkdir(m.logDir)
	}

	f := strSplit(basename(path), ".")
	m.logFileName = strJoin(".", f[:len(f)-1])
	m.logFileSuffix = f[len(f)-1]

	var logpath string
	if len(logFileSizeInMB) != 0 {
		m.logFileSizeInMB = logFileSizeInMB[0]
		for {
			logpath = pathJoin(m.logDir, m.logFileName+"."+str(m.currentLogFileNumber)+"."+m.logFileSuffix)
			if pathExists(logpath) {
				m.currentLogFileNumber++
			} else {
				break
			}
		}
	} else {
		logpath = pathJoin(m.logDir, m.logFileName+"."+strftime("%Y-%m-%d", now())+"."+m.logFileSuffix)
	}
	m.fd = open(logpath, "a")
	m.logfiles = append(m.logfiles, logpath)

	m.maxlogfiles = maxLogFileCount
}

func (m *logStruct) error(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", now())
	level := "ERRO"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, fmt.Sprint(a))
	}
	msg := strJoin(" ", msgarr)

	_, file, no, _ := runtime.Caller(1)
	position := basename(file) + ":" + toString(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) warn(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", now())
	level := "WARN"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, fmt.Sprint(a))
	}
	msg := strJoin(" ", msgarr)

	_, file, no, _ := runtime.Caller(1)
	position := basename(file) + ":" + toString(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) info(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", now())
	level := "INFO"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, fmt.Sprint(a))
	}
	msg := strJoin(" ", msgarr)

	_, file, no, _ := runtime.Caller(1)
	position := basename(file) + ":" + toString(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) trace(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", now())
	level := "TRAC"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, fmt.Sprint(a))
	}
	msg := strJoin(" ", msgarr)

	_, file, no, _ := runtime.Caller(1)
	position := basename(file) + ":" + toString(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) debug(args ...interface{}) {
	t := strftime("%m-%d %H:%M:%S", now())
	level := "DEBU"

	var msgarr []string
	for _, a := range args {
		msgarr = append(msgarr, sprint(a))
	}
	msg := strJoin(" ", msgarr)

	_, file, no, _ := runtime.Caller(1)
	position := basename(file) + ":" + toString(no)

	m.show(t, level, msg, position)
}

func (m *logStruct) show(t string, level string, msg string, position string) {
	if itemInArray(level, m.level) {
		var strMsg string
		if m.json {
			strMsg = jsonDumps(map[string]string{
				"time":    t,
				"level":   level,
				"message": msg,
			})
		} else {
			msg = strReplace(msg, "\n", "\n                      ")
			if m.color {
				if level == "ERRO" {
					strMsg = color.RedString(t + fmt.Sprintf(" %3v", getGoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				} else if level == "WARN" {
					strMsg = color.YellowString(t + fmt.Sprintf(" %3v", getGoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				} else if level == "INFO" {
					strMsg = color.HiBlueString(t + fmt.Sprintf(" %3v", getGoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				} else if level == "TRAC" {
					strMsg = color.MagentaString(t + fmt.Sprintf(" %3v", getGoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				} else if level == "DEBU" {
					strMsg = color.HiCyanString(t + fmt.Sprintf(" %3v", getGoroutineID()) + " [" + level + "] (" + position + ") " + msg)
				}
			} else {
				strMsg = t + "[" + level + "] (" + position + ") " + msg
			}
		}

		m.lock.acquire()
		if m.displayOnTerminal {
			fmt.Println(strMsg)
		}
		if m.fd != nil {
			if m.logFileSizeInMB == 0 {
				if m.fd.path != pathJoin(m.logDir, m.logFileName+"."+strftime("%Y-%m-%d", now())+"."+m.logFileSuffix) {
					m.fd.close()
					logpath := pathJoin(m.logDir, m.logFileName+"."+strftime("%Y-%m-%d", now())+"."+m.logFileSuffix)
					m.fd = open(logpath, "a")
					m.logfiles = append(m.logfiles, logpath)
					if len(m.logfiles) > m.maxlogfiles {
						unlink(m.logfiles[0])
						m.logfiles = m.logfiles[1:]
					}
				}
			} else {
				if m.currentLogFileSizeInByte > m.logFileSizeInMB*1024*1024 {
					m.currentLogFileSizeInByte = 0
					m.fd.close()
					var logpath string
					for {
						logpath = pathJoin(m.logDir, m.logFileName+"."+str(m.currentLogFileNumber)+"."+m.logFileSuffix)
						if pathExists(logpath) {
							m.currentLogFileNumber++
						} else {
							break
						}
					}
					m.fd = open(logpath, "a")
					m.logfiles = append(m.logfiles, logpath)
					if len(m.logfiles) > m.maxlogfiles {
						unlink(m.logfiles[0])
						m.logfiles = m.logfiles[1:]
					}
				}
			}
			m.fd.write(strMsg + "\n")
			m.currentLogFileSizeInByte = m.currentLogFileSizeInByte + len(strMsg) + 1
		}
		m.lock.release()
	}
}
func getSnowflakeID(nodeNumber ...int) int64 {
	var node *snowflake.Node
	var err error
	if len(nodeNumber) != 0 {
		node, err = snowflake.NewNode(toInt64(nodeNumber[0]))
	} else {
		node, err = snowflake.NewNode(1)
	}
	panicerr(err)
	return node.Generate().Int64()
}
//  将字符转换为 HTML 转义字符
func htmlencode(str string) string {
	return html.EscapeString(str)
}

//  Convert HTML entities to their corresponding characters
func htmldecode(str string) string {
	return html.UnescapeString(str)
}
type networkPacketStruct struct {
	data  string
	sport int
	dport int
	proto string // tcp, udp
	ipv   int    // 4, 6
	sip   string
	dip   string
	smac  string
	dmac  string
}

func doPacketSource(packetSource *gopacket.PacketSource, pkgchan chan *networkPacketStruct, pcapFileHandler ...*pcap.Handle) {
	for packet := range packetSource.Packets() {
		//print("Packet found: ", packet)
		transportLayer := packet.TransportLayer()
		if transportLayer != nil {
			//print("transportLayer found.")
			pkg := networkPacketStruct{}

			linuxSLLLayer := packet.Layer(layers.LayerTypeLinuxSLL)
			if linuxSLLLayer != nil {
				linuxSLLPacket, _ := linuxSLLLayer.(*layers.LinuxSLL)
				pkg.smac = fmt.Sprintf("%s", linuxSLLPacket.Addr)
			}

			//print(packet)
			ethLayer := packet.Layer(layers.LayerTypeEthernet)
			if ethLayer != nil {
				//print("eth layer found")
				ethernetPacket, _ := ethLayer.(*layers.Ethernet)
				pkg.smac = fmt.Sprintf("%s", ethernetPacket.SrcMAC)
				pkg.dmac = fmt.Sprintf("%s", ethernetPacket.DstMAC)
				//fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
			}

			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer != nil {
				//print("ip layer found")
				ip, ok := ipLayer.(*layers.IPv4)
				if ok {
					pkg.ipv = 4
					pkg.sip = fmt.Sprintf("%s", ip.SrcIP)
					pkg.dip = fmt.Sprintf("%s", ip.DstIP)
				} else {
					pkg.ipv = 6
					ip6, _ := ipLayer.(*layers.IPv6)
					pkg.sip = fmt.Sprintf("%s", ip6.SrcIP)
					pkg.dip = fmt.Sprintf("%s", ip6.DstIP)
				}
			}

			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				//print("tcp layer found")
				pkg.proto = "tcp"
				tcp, _ := tcpLayer.(*layers.TCP)
				pkg.sport = toInt(fmt.Sprintf("%d", tcp.SrcPort))
				pkg.dport = toInt(fmt.Sprintf("%d", tcp.DstPort))
			}

			udpLayer := packet.Layer(layers.LayerTypeUDP)
			if udpLayer != nil {
				//print("udp layer found")
				pkg.proto = "udp"
				udp, _ := udpLayer.(*layers.UDP)
				pkg.sport = toInt(fmt.Sprintf("%d", udp.SrcPort))
				pkg.dport = toInt(fmt.Sprintf("%d", udp.DstPort))
			}

			applicationLayer := packet.TransportLayer()
			if applicationLayer != nil {
				pkg.data = toString(applicationLayer.LayerPayload())
				//print("Data:", pkg.data)
			}

			// if strStartsWith(pkg.data, "GET /action") {
			// 	print(packet)
			// }

			if pkg.data != "" {
				pkgchan <- &pkg
			}
		}
	}
	if len(pcapFileHandler) != 0 {
		pcapFileHandler[0].Close()
	}
	close(pkgchan)
}

func sniffer(interfaceName string, filterString string, promisc ...bool) chan *networkPacketStruct {
	// 4096是读取每一个包的buffer, mtu一般为1500, 所以4096是超出了很多, 除非mtu超出了4096, 才读不全
	// promisc为设置网卡为混杂模式
	// timeout为0.3秒, 是kernel每0.3秒就会吐一次数据给pcap, 如果这个为30秒, 则收到数据包之后会继续等待其他数据包, 30秒再一起吐出来
	var handle *pcap.Handle
	var err error
	if len(promisc) == 0 {
		handle, err = pcap.OpenLive(interfaceName, 4096, false, getTimeDuration(0.3))
	} else {
		handle, err = pcap.OpenLive(interfaceName, 4096, promisc[0], getTimeDuration(0.3))
	}

	panicerr(err)
	//defer handle.Close()

	err = handle.SetBPFFilter(filterString)
	panicerr(err)

	pkgchan := make(chan *networkPacketStruct)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go doPacketSource(packetSource, pkgchan)

	return pkgchan
}

func readPcapFile(pcapFile string) chan *networkPacketStruct {
	handle, err := pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}

	pkgchan := make(chan *networkPacketStruct)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go doPacketSource(packetSource, pkgchan, handle)

	return pkgchan
}
type alidnsStruct struct {
	client          *alidns20150109.Client
	accessKeyID     string
	accessKeySecret string
}

type alidnsDomainInfoStruct struct {
	DomainName string
	DomainID   string
	CreateTime string
}

type alidnsDomainStruct struct {
	client     *alidns20150109.Client
	DomainName string
}

type alidnsRecord struct {
	ID   string
	Data string
	Name string
	TTL  int
	Type string
}

func getAlidns(accessKeyID string, accessKeySecret string) *alidnsStruct {
	config := &openapi.Config{
		AccessKeyId:     &accessKeyID,
		AccessKeySecret: &accessKeySecret,
	}

	config.Endpoint = tea.String("dns.aliyuncs.com")
	client, err := alidns20150109.NewClient(config)
	panicerr(err)

	return &alidnsStruct{
		client:          client,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
	}
}

func (m *alidnsStruct) total() (TotalCount int64) {
	describeDomainsRequest := &alidns20150109.DescribeDomainsRequest{}
	result, err := m.client.DescribeDomains(describeDomainsRequest)
	panicerr(err)
	TotalCount = *result.Body.TotalCount
	return
}

func (m *alidnsStruct) list(PageSize int64, PageNumber int64) (res []alidnsDomainInfoStruct) {
	describeDomainsRequest := &alidns20150109.DescribeDomainsRequest{
		PageSize:   &PageSize,
		PageNumber: &PageNumber,
	}
	result, err := m.client.DescribeDomains(describeDomainsRequest)
	panicerr(err)

	for _, d := range result.Body.Domains.Domain {
		res = append(res, alidnsDomainInfoStruct{
			DomainName: *d.DomainName,
			DomainID:   *d.DomainId,
			CreateTime: *d.CreateTime,
		})
	}
	return
}

func (m *alidnsStruct) domain(domainName string) *alidnsDomainStruct {
	return &alidnsDomainStruct{
		client:     m.client,
		DomainName: domainName,
	}
}

func (m *alidnsDomainStruct) list() (res []alidnsRecord) {
	result, err := m.client.DescribeDomainRecords(&alidns20150109.DescribeDomainRecordsRequest{
		DomainName: &m.DomainName,
	})
	panicerr(err)
	for _, r := range result.Body.DomainRecords.Record {
		res = append(res, alidnsRecord{
			ID:   *r.RecordId,
			Data: *r.Value,
			Name: *r.RR,
			TTL:  toInt(*r.TTL),
			Type: *r.Type,
		})
	}
	return
}

func (m *alidnsDomainStruct) add(recordName string, recordType string, recordValue string) (id string) {
	recordType = strUpper(recordType)
	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: &m.DomainName,
		RR:         &recordName,
		Type:       &recordType,
		Value:      &recordValue,
	}
	res, err := m.client.AddDomainRecord(addDomainRecordRequest)
	panicerr(err)
	return *res.Body.RecordId
}

func (m *alidnsDomainStruct) delete(name string, dtype string, value string) {
	for _, d := range m.list() {
		var nameres, dtyperes, valueres bool = false, false, false
		if name == "" || d.Name == name {
			nameres = true
		}
		if dtype == "" || strLower(d.Type) == strLower(dtype) {
			dtyperes = true
		}
		if value == "" || d.Data == value {
			valueres = true
		}
		if nameres && dtyperes && valueres {
			_, err := m.client.DeleteDomainRecord(&alidns20150109.DeleteDomainRecordRequest{
				RecordId: &d.ID,
			})
			panicerr(err)
		}
	}
}

func (m *alidnsDomainStruct) modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordName string, dstRecordType string, dstRecordValue string) {
	m.delete(recordName, srcRecordType, srcRecordValue)
	m.add(dstRecordName, dstRecordType, dstRecordValue)
}
func getPinYin(zh string) (ress []string) {
	a := pinyin.NewArgs()
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	p := pinyin.LazyPinyin(zh, a)

	filterChar := "1234567890qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM"

	var res []string
	for idx, v := range p {
		if len(v) != 1 {
			res = append(res, v)
		} else {
			if idx == 0 {
				res = append(res, v)
			} else {
				if len(p[idx-1]) == 1 {
					if strIn(p[idx-1], filterChar) {
						res[len(res)-1] = res[len(res)-1] + v
					} else {
						res = append(res, v)
					}
				} else {
					res = append(res, v)
				}
			}
		}
	}

	for _, v := range res {
		vv := strFilter(v, filterChar)
		if len(vv) != 0 {
			ress = append(ress, vv)
		}
	}
	return
}
func drawLineChartWithTimeSeries(timestampX []int64, dataY []float64, xtitle string, ytitle string, title string, fpath string) {
	if len(timestampX) != len(dataY) {
		panic(newerr("The number of elements on the X axis and Y axis must be the same"))
	}

	var dataX []time.Time
	for _, i := range timestampX {
		tm := time.Unix(i, 0)
		dataX = append(dataX, tm)
	}

	statikFS, err := fs.New()
	panicerr(err)

	simHeiFile, err := statikFS.Open("/SimHei.ttf")
	panicerr(err)
	fontBytes, err := ioutil.ReadAll(simHeiFile)
	panicerr(err)
	font, err := truetype.Parse(fontBytes)
	panicerr(err)
	simHeiFile.Close()

	graph := chart.Chart{
		Font:  font,
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		Height: 768,
		Width:  2000,
		XAxis: chart.XAxis{
			Name: xtitle,
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Name: ytitle,
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
					Show:        true,
				},
				XValues: dataX,
				YValues: dataY,
			},
		},
	}

	f, err := os.Create(fpath)
	panicerr(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func drawLineChartWithNumberSeries(dataX []float64, dataY []float64, xtitle string, ytitle string, title string, fpath string) {
	if len(dataX) != len(dataY) {
		panic(newerr("The number of elements on the X axis and Y axis must be the same"))
	}

	statikFS, err := fs.New()
	panicerr(err)

	simHeiFile, err := statikFS.Open("/SimHei.ttf")
	panicerr(err)
	fontBytes, err := ioutil.ReadAll(simHeiFile)
	panicerr(err)
	font, err := truetype.Parse(fontBytes)
	panicerr(err)
	simHeiFile.Close()

	graph := chart.Chart{
		Font:  font,
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		Height: 768,
		Width:  2000,
		XAxis: chart.XAxis{
			Name: xtitle,
			Style: chart.Style{
				Show: true,
			},
			NameStyle: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Name: ytitle,
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: dataX,
				YValues: dataY,
			},
		},
	}

	f, err := os.Create(fpath)
	panicerr(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func drawBarChartWithSeries(dataX []string, dataY []float64, ytitle string, title string, fpath string) {
	if len(dataX) != len(dataY) {
		panic(newerr("The number of elements on the X axis and Y axis must be the same"))
	}

	statikFS, err := fs.New()
	panicerr(err)

	simHeiFile, err := statikFS.Open("/SimHei.ttf")
	panicerr(err)
	fontBytes, err := ioutil.ReadAll(simHeiFile)
	panicerr(err)
	font, err := truetype.Parse(fontBytes)
	panicerr(err)
	simHeiFile.Close()

	var chartValue []chart.Value
	for i := 0; i < len(dataX); i++ {
		chartValue = append(chartValue, chart.Value{Value: dataY[i], Label: dataX[i]})
	}

	graph := chart.BarChart{
		Font:  font,
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		YAxis: chart.YAxis{
			Name: ytitle,
			NameStyle: chart.Style{
				Show: true,
			},
			Style: chart.Style{
				Show: true,
			},
		},
		Height: 768,
		Width:  len(dataX) * 120,
		Bars:   chartValue,
	}

	f, err := os.Create(fpath)
	panicerr(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func drawPieChartWithSeries(dataX []string, dataY []float64, title string, fpath string) {
	if len(dataX) != len(dataY) {
		panic(newerr("The number of elements on the X axis and Y axis must be the same"))
	}

	statikFS, err := fs.New()
	panicerr(err)

	simHeiFile, err := statikFS.Open("/SimHei.ttf")
	panicerr(err)
	fontBytes, err := ioutil.ReadAll(simHeiFile)
	panicerr(err)
	font, err := truetype.Parse(fontBytes)
	panicerr(err)
	simHeiFile.Close()

	var chartValue []chart.Value
	for i := 0; i < len(dataX); i++ {
		chartValue = append(chartValue, chart.Value{Value: dataY[i], Label: dataX[i]})
	}

	graph := chart.PieChart{
		Font:   font,
		Height: 2000,
		Width:  2000,
		Values: chartValue,
		Title:  title,
		TitleStyle: chart.Style{
			Show: true,
		},
	}

	f, err := os.Create(fpath)
	panicerr(err)
	defer f.Close()
	graph.Render(chart.PNG, f)
}
// SSL - Server
// 只实现了一个Listener， 其他的方法是tcp的方法

func sslListen(host string, port int, key string, crt string) *tcpServerSideListener {
	cert, err := tls.LoadX509KeyPair(crt, key)
	panicerr(err)

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp4", host+":"+toString(port), tlsCfg)
	panicerr(err)

	return &tcpServerSideListener{listener: listener}
}

// SSL - Client

type sslClientSideConn struct {
	conn *tls.Conn
}

func sslConnect(host string, port int) *sslClientSideConn {
	servAddr := host + ":" + toString(port)
	conn, err := tls.Dial("tcp", servAddr, &tls.Config{})
	panicerr(err)
	return &sslClientSideConn{conn: conn}
}

func (m *sslClientSideConn) send(str string) {
	_, err := m.conn.Write([]byte(str))
	panicerr(err)
}

func (m *sslClientSideConn) recv(buffersize int) string {
	reply := make([]byte, buffersize)
	n, err := m.conn.Read(reply)
	panicerr(err)
	return string(reply[:n])
}

func (m *sslClientSideConn) close() {
	err := m.conn.Close()
	panicerr(err)
}
type errorStruct struct {
	msg string
}

func fmtDebugStack(msg string, stack string) string {
	//lg.debug("msg:", msg)
	//lg.debug("stack:", stack)

	blackFileList := []string{
		"lib.go",
		"stack.go",
	}

	l := reFindAll("([\\-a-zA-Z0-9]+\\.go:[0-9]+)", stack)
	//lg.debug(l)
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	//lg.debug(l)

	var link []string
	for _, f := range l {
		ff := strings.Split(f[0], ":")[0]
		inside := func(a string, list []string) bool {
			for _, b := range list {
				if b == a {
					return true
				}
			}
			return false
		}(ff, blackFileList)
		if !inside {
			link = append(link, f[0])
		}
	}
	//lg.debug(link)

	var strr string
	if len(link) != 1 {
		strr = link[len(link)-2] + " >> " + msg + " >> " + "(" + strJoin(" => ", link[:len(link)-1]) + ")"
	} else {
		strr = link[0] + " >> " + msg
	}

	//lg.debug("strr:", strr)
	return strr
}

func newerr(msg interface{}) *errorStruct {
	switch t := msg.(type) {
	case string:
		return &errorStruct{
			msg: fmtDebugStack(t, string(debug.Stack())),
		}
	case error:
		return &errorStruct{
			msg: fmtDebugStack(t.Error(), string(debug.Stack())),
		}
	default:
		return &errorStruct{
			msg: fmtDebugStack(fmt.Sprintf("%s", t), string(debug.Stack())),
		}
	}
}

func panicerr(err interface{}) {
	switch t := err.(type) {
	case string:
		// lg.trace("1")
		panic(fmtDebugStack(t, string(debug.Stack())))
	case error:
		// lg.trace("2")
		panic(fmtDebugStack(t.Error(), string(debug.Stack())))
	case *errorStruct:
		// lg.trace(3)
		panic(t.msg)
	case nil:
		return
	default:
		panic(fmtDebugStack(fmt.Sprintf("%s", t), string(debug.Stack())))
	}
}
type progressCPUUsageStruct struct {
	pid  int64
	cpu  float64
	cmd  string
	name string
}

func getSystemProgressCPUUsage() (res map[int64]progressCPUUsageStruct) {
	pg := make(map[string]float64)
	res = make(map[int64]progressCPUUsageStruct)

	var line string
	for _, line = range strSplit(open("/proc/stat").read(), "\n") {
		if strIn("cpu ", line) {
			break
		}
	}

	var totalCPUSlice1 float64
	for _, i := range strSplit(line)[2:] {
		totalCPUSlice1 = totalCPUSlice1 + toFloat64(i)
	}

	for _, pid := range listdir("/proc") {
		if !isdigit(pid) {
			continue
		}
		try(func() {
			a := strSplit(open("/proc/" + pid + "/stat").read())
			totalProcessSlice1 := toFloat64(toInt(a[13]) + toInt(a[14]) + toInt(a[15]) + toInt(a[16]))
			pg[pid] = totalProcessSlice1
		})
	}

	sleep(1)

	for _, line = range strSplit(open("/proc/stat").read(), "\n") {
		if strIn("cpu ", line) {
			break
		}
	}

	var totalCPUSlice2 float64
	for _, i := range strSplit(line)[2:] {
		totalCPUSlice2 = totalCPUSlice2 + toFloat64(i)
	}

	for _, pid := range listdir("/proc") {
		if !isdigit(pid) {
			continue
		}
		try(func() {
			a := strSplit(open("/proc/" + pid + "/stat").read())
			totalProcessSlice2 := toFloat64(toInt(a[13]) + toInt(a[14]) + toInt(a[15]) + toInt(a[16]))
			_, found := pg[pid]
			if found {
				cpuusage := (totalProcessSlice2 - pg[pid]) / (totalCPUSlice2 - totalCPUSlice1) * 100 * toFloat64(runtime.NumCPU())
				res[toInt64(pid)] = progressCPUUsageStruct{
					pid:  toInt64(pid),
					cpu:  cpuusage,
					cmd:  strStrip(strReplace(open("/proc/"+pid+"/cmdline").read(), "\x00", " ")),
					name: strStrip(a[1], "()"),
				}
			}
		})
	}

	return
}
type exception struct {
	finallyFunc func()
	Error       error
}

func throw() {
	panic("_____rethrow")
}

func try(f func()) (e exception) {
	e = exception{nil, nil}
	defer func() {
		if err := recover(); err != nil {
			e.Error = errors.New(fmt.Sprintf("%s", err))
		}
	}()
	f()
	return
}

func (e exception) except(f func(err error)) error {
	if e.Error != nil {
		defer func() {
			if err := recover(); err != nil {
				if err == "_____rethrow" {
					err = e.Error
				}
				panic(err)
			}
		}()
		f(e.Error)
	}
	return e.Error
}
var KCP_PING uint32 = 4000000000
var KCP_PONG uint32 = 4000000001
var KCP_CLOSE uint32 = 4000000002
var KCP_HEARTBEAT_SECOND int = 20

func kcpRecvSendChanIsClosed(ch chan map[string]string) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

// KCP - Server

type kcpServerSideConn struct {
	conn          *kcp.UDPSession
	isclose       bool
	readtimeout   int
	writetimeout  int
	sendchan      chan map[string]string
	recvchan      chan map[string]string
	heartbeatTime int64
}

type kcpServerSideListener struct {
	listener *kcp.Listener
	isclose  bool
}

func kcpListen(host string, port int, key string, salt string) *kcpServerSideListener {
	lg.trace("kcplisten", 1)
	block, err := kcp.NewAESBlockCrypt(pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha256.New))
	panicerr(err)

	l, err := kcp.ListenWithOptions(host+":"+toString(port), block, 10, 3)
	panicerr(err)

	l.SetDSCP(46)
	l.SetReadBuffer(4194304)
	l.SetWriteBuffer(4194304)

	return &kcpServerSideListener{listener: l}
}

func (m *kcpServerSideListener) accept() chan *kcpServerSideConn {
	ch := make(chan *kcpServerSideConn)

	go func() {
		lg.trace("开启accept的goroutine")
		for {
			lg.trace("循环接收KCP连接")
			var mc *kcpServerSideConn
			if err := try(func() {
				c, err := m.listener.AcceptKCP()
				lg.trace("收到一个新连接")
				if err != nil {
					lg.trace("但是出错了")
					if strIn("io: read/write on closed pipe", err.Error()) || strIn("use of closed network connection", err.Error()) {
						close(ch)
						m.close()
					}
					panicerr(err)
				}
				lg.trace("没有错误发生")
				//lg.debug(c)

				c.SetNoDelay(0, 20, 2, 1)
				c.SetMtu(1400)
				c.SetWindowSize(1024, 1024)
				c.SetACKNoDelay(false)

				mc = &kcpServerSideConn{
					conn:          c,
					isclose:       false,
					readtimeout:   120,
					writetimeout:  120,
					heartbeatTime: now(),
					sendchan:      make(chan map[string]string),
					recvchan:      make(chan map[string]string, 10),
				}
			}).Error; err != nil {
				lg.trace("有错误发生")
				if m.isclose {
					lg.trace("kcpServerSideListener已关闭，错误是正常，退出")
					break
				} else {
					lg.trace("kcpServerSideListener未关闭，错误是异常：", err)
					sleep(1)
					continue
				}
			}
			lg.trace("accept的goroutine启动完毕")

			// sender
			go func(mc *kcpServerSideConn) {
				lg.trace("sender的groutine启动")
				try(func() {
					for {
						// 如果连接，关闭，关掉chan，退出
						if mc.isclose {
							if !kcpRecvSendChanIsClosed(mc.sendchan) {
								if !kcpRecvSendChanIsClosed(mc.sendchan) {
									close(mc.sendchan)
								}
							}
							lg.trace("因为kcpServerSideConn已关闭，sender的groutine退出")
							break
						}
						// 如果没有数据就等待1秒再查看，有数据就处理
						select {
						case data, ok := <-mc.sendchan:
							if ok {
								mc.conn.SetWriteDeadline(time.Now().Add(time.Duration(mc.writetimeout) * time.Second))

								btlen := make([]byte, 4)
								// 要发送的字符串
								text := map2bin(data)
								binary.LittleEndian.PutUint32(btlen, uint32(len(text)))

								_, err := mc.conn.Write(btlen)
								if err != nil {
									if mc.isclose {
										lg.trace("sender的groutine退出，因为kcpServerSideConn已关闭")
										break
									} else {
										panicerr(err)
									}
								}

								_, err = mc.conn.Write([]byte(text))
								if err != nil {
									if mc.isclose {
										lg.trace("sender的groutine退出，因为kcpServerSideConn已关闭")
										break
									} else {
										panicerr(err)
									}
								}
							}
						case <-time.After(getTimeDuration(1)):

						}
					}
				}).except(func(e error) {
					lg.trace("kcpServerSideConn sender groutine 已退出, 报错:", e)
				})
			}(mc)

			// receiver
			go func(mc *kcpServerSideConn) {
				lg.trace("receiver的goroutine启动")
				try(func() {
					for {
						mc.conn.SetReadDeadline(time.Now().Add(time.Duration(mc.readtimeout) * time.Second))

						// 读取总长度, 然后读取所有数据
						totalblen := make([]byte, 4096)
						n, err := mc.conn.Read(totalblen)
						if err != nil {
							if mc.isclose || strIn("io: read/write on closed pipe", err.Error()) {
								if mc.isclose {
									lg.trace("因为kcpServerSideConn关闭，receiver退出")
								} else {
									lg.trace("因为kcpServerSideConn关闭导致写入失败，receiver退出")
								}
								mc.isclose = true
								if !kcpRecvSendChanIsClosed(mc.recvchan) {
									close(mc.recvchan)
								}
								break
							} else {
								panicerr(err)
							}
						}
						if n != 4 {
							continue
						}

						totallen := binary.LittleEndian.Uint32(totalblen)

						// 如果是心跳，更新心跳时间，回复它，处理下一个包
						if totallen == KCP_PING {
							btlen := make([]byte, 4)
							binary.LittleEndian.PutUint32(btlen, KCP_PONG)
							_, err := mc.conn.Write(btlen)
							if err != nil {
								if mc.isclose {
									lg.trace("因为kcpServerSideConn关闭，receiver退出")
									break
								} else {
									panicerr(err)
								}
							}
							mc.heartbeatTime = now()
							continue
						}

						totaldata := ""
						buf := make([]byte, 4096)
						for len(totaldata) != toInt(totallen) {
							n, err := mc.conn.Read(buf)
							if err != nil {
								if mc.isclose {
									lg.trace("因为kcpServerSideConn关闭，receiver退出")
									if !kcpRecvSendChanIsClosed(mc.recvchan) {
										close(mc.recvchan)
									}
									break
								} else {
									panicerr(err)
								}
							}

							if n == 4 {
								buflen := binary.LittleEndian.Uint32(buf)
								// 如果是心跳，更新心跳时间，回复它，处理下一个包
								if buflen == KCP_PING {
									btlen := make([]byte, 4)
									binary.LittleEndian.PutUint32(btlen, KCP_PONG)
									_, err := mc.conn.Write(btlen)
									if err != nil {
										if mc.isclose {
											lg.trace("因为kcpServerSideConn关闭，receiver退出")
											break
										} else {
											panicerr(err)
										}
									}
									mc.heartbeatTime = now()
									continue
								}
							}

							totaldata = totaldata + string(buf[:n])
						}

						res := bin2map(totaldata)
						mc.recvchan <- res
					}
				}).except(func(e error) {
					lg.trace("kcpServerSideConn receiver groutine 已退出, 报错:", e)
				})
			}(mc)

			// heartbeat checker
			// 如果连接被主动关闭，退出，如果3次没有收到心跳，关闭连接
			go func(mc *kcpServerSideConn) {
				lg.trace("heartbeat checker的goroutine启动")
				try(func() {
					for {
						if mc.isclose {
							break
						}
						if now()-mc.heartbeatTime > toInt64(KCP_HEARTBEAT_SECOND)*3 {
							if !kcpRecvSendChanIsClosed(mc.recvchan) {
								close(mc.recvchan)
							}
							if !kcpRecvSendChanIsClosed(mc.sendchan) {
								close(mc.sendchan)
							}
							lg.trace("因为kcpServerSideConn关闭，heartbeat checker退出")
							mc.close()
							break
						}
						sleep(1)
					}
				}).except(func(e error) {
					lg.trace("kcpServerSideConn heartbeat checker groutine 已退出, 报错:", e)
				})
			}(mc)

			ch <- mc
		}
	}()

	return ch
}

func (m *kcpServerSideListener) close() {
	if !m.isclose {
		m.isclose = true
		m.listener.Close()
	}
}

func (m *kcpServerSideConn) close() {
	if !m.isclose {
		m.isclose = true
		try(func() {
			m.conn.Close()
		})
	}
	if !kcpRecvSendChanIsClosed(m.sendchan) {
		close(m.sendchan)
	}
	if !kcpRecvSendChanIsClosed(m.recvchan) {
		close(m.recvchan)
	}
}

func (m *kcpServerSideConn) send(data map[string]string) {
	if m.isclose {
		err := errors.New("连接已关闭，不可发送数据")
		panicerr(err)
	}
	m.sendchan <- data
}

func (m *kcpServerSideConn) recv(timeoutSecond ...int) (res map[string]string) {
	if m.isclose {
		err := errors.New("连接已关闭，不可收取数据")
		panicerr(err)
	}
	if len(timeoutSecond) == 0 {
		res = <-m.recvchan
	} else {
		for range rangeInt(timeoutSecond[0] * 10) {
			select {
			case data, ok := <-m.recvchan:
				if ok {
					res = data
					break
				}
			default:
				sleep(0.1)
			}
			if res != nil {
				break
			}
		}
	}

	if res == nil {
		if m.isclose {
			err := errors.New("连接已关闭，不可收取数据")
			panicerr(err)
		}
	}
	return
}

// KCP - Client

type kcpClientSideConn struct {
	conn          *kcp.UDPSession
	isclose       bool
	readtimeout   int
	writetimeout  int
	recvchan      chan map[string]string
	sendchan      chan map[string]string
	heartbeatTime int64
}

func kcpConnect(host string, port int, key string, salt string) *kcpClientSideConn {
	block, err := kcp.NewAESBlockCrypt(pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha256.New))
	panicerr(err)
	conn, err := kcp.DialWithOptions(host+":"+toString(port), block, 10, 3)
	panicerr(err)

	conn.SetMtu(1400)
	conn.SetWriteDelay(false)
	conn.SetNoDelay(0, 20, 2, 1)
	conn.SetWindowSize(128, 1024)
	conn.SetACKNoDelay(false)
	conn.SetDSCP(46)
	conn.SetReadBuffer(4194304)
	conn.SetWriteBuffer(4194304)

	m := &kcpClientSideConn{
		conn:          conn,
		isclose:       false,
		readtimeout:   120,
		writetimeout:  120,
		sendchan:      make(chan map[string]string),
		recvchan:      make(chan map[string]string, 10),
		heartbeatTime: now(),
	}

	// sender
	go func(m *kcpClientSideConn) {
		try(func() {
			for {
				// 如果连接，关闭，关掉chan，退出
				if m.isclose {
					if !kcpRecvSendChanIsClosed(m.sendchan) {
						close(m.sendchan)
					}
					break
				}
				// 如果没有数据就等待0.3秒再查看，有数据就处理
				select {
				case data, ok := <-m.sendchan:
					if ok {
						m.conn.SetWriteDeadline(time.Now().Add(time.Duration(m.writetimeout) * time.Second))
						// 要发送的字符串
						text := map2bin(data)
						// 4字节的二进制无符号int
						btlen := make([]byte, 4)

						binary.LittleEndian.PutUint32(btlen, uint32(len(text)))

						_, err := m.conn.Write(btlen)
						if err != nil {
							if m.isclose {
								break
							} else {
								panicerr(err)
							}
						}

						_, err = m.conn.Write([]byte(text))
						if err != nil {
							if m.isclose {
								break
							} else {
								panicerr(err)
							}
						}
					}
				case <-time.After(getTimeDuration(1)):
				}
			}
		}).except(func(e error) {
			lg.trace("kcpClientSideConn sender groutine 已退出, 报错:", e)
		})
	}(m)

	// receiver
	go func(m *kcpClientSideConn) {
		try(func() {
			for {
				m.conn.SetReadDeadline(time.Now().Add(time.Duration(m.readtimeout) * time.Second))

				// 读取总长度, 然后读取所有数据
				totalblen := make([]byte, 4096)
				n, err := m.conn.Read(totalblen)
				if err != nil {
					if m.isclose {
						if !kcpRecvSendChanIsClosed(m.recvchan) {
							close(m.recvchan)
						}
						break
					} else {
						panicerr(err)
					}
				}
				if n != 4 {
					continue
				}

				totallen := binary.LittleEndian.Uint32(totalblen)

				// 如果是心跳的回应, 继续下一个包
				if totallen == KCP_PONG {
					m.heartbeatTime = now()
					continue
				}

				totaldata := ""
				buf := make([]byte, 4096)
				for len(totaldata) != toInt(totallen) {
					n, err := m.conn.Read(buf)
					if err != nil {
						if m.isclose {
							if !kcpRecvSendChanIsClosed(m.recvchan) {
								close(m.recvchan)
							}
							break
						} else {
							panicerr(err)
						}
					}

					if n == 4 {
						buflen := binary.LittleEndian.Uint32(buf)
						// 如果是心跳的回应, 继续下一个包
						if buflen == KCP_PONG {
							m.heartbeatTime = now()
							continue
						}
					}

					totaldata = totaldata + string(buf[:n])
				}

				res := bin2map(totaldata)
				m.recvchan <- res
			}
		}).except(func(e error) {
			lg.trace("kcpClientSideConn receiver groutine 已退出, 报错:", e)
		})
	}(m)

	// heartbeat
	go func(m *kcpClientSideConn) {
		try(func() {
			for {
				btlen := make([]byte, 4)
				binary.LittleEndian.PutUint32(btlen, KCP_PING)
				_, err := m.conn.Write(btlen)
				if err != nil {
					if m.isclose {
						if !kcpRecvSendChanIsClosed(m.sendchan) {
							close(m.sendchan)
						}
						if !kcpRecvSendChanIsClosed(m.recvchan) {
							close(m.recvchan)
						}
						break
					} else {
						panicerr(err)
					}
				}
				if now()-m.heartbeatTime > toInt64(KCP_HEARTBEAT_SECOND)*3 {
					m.close()
				}
				sleep(KCP_HEARTBEAT_SECOND)
			}
		}).except(func(e error) {
			lg.trace("kcpClientSideConn heartbeat groutine 已退出, 报错:", e)
		})
	}(m)

	return m
}

func (m *kcpClientSideConn) send(data map[string]string) {
	if m.isclose {
		err := errors.New("连接已关闭，不可发送数据")
		panicerr(err)
	}
	m.sendchan <- data
}

func (m *kcpClientSideConn) recv(timeoutSecond ...int) (res map[string]string) {
	if m.isclose {
		err := errors.New("连接已关闭，不可收取数据")
		panicerr(err)
	}
	if len(timeoutSecond) == 0 {
		res = <-m.recvchan
	} else {
		for range rangeInt(timeoutSecond[0] * 10) {
			select {
			case data, ok := <-m.recvchan:
				if ok {
					res = data
					break
				}
			default:
				sleep(0.1)
			}
			if res != nil {
				break
			}
		}
	}

	if res == nil {
		if m.isclose {
			err := errors.New("连接已关闭，不可收取数据")
			panicerr(err)
		}
	}
	return
}

func (m *kcpClientSideConn) close() {
	if !m.isclose {
		m.isclose = true
		m.conn.Close()
	}

	if !kcpRecvSendChanIsClosed(m.sendchan) {
		close(m.sendchan)
	}
	if !kcpRecvSendChanIsClosed(m.recvchan) {
		close(m.recvchan)
	}
}
type matrixStruct struct {
	cli    *gomatrix.Client
	roomID string
}

func getMatrix(homeserverURL string) *matrixStruct {
	cli, err := gomatrix.NewClient(homeserverURL, "", "")
	panicerr(err)

	return &matrixStruct{cli: cli}
}

func (c *matrixStruct) login(username string, password string) *matrixStruct {
	resp, err := c.cli.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     username,
		Password: password,
	})
	panicerr(err)

	c.cli.SetCredentials(resp.UserID, resp.AccessToken)
	return c
}

func (c *matrixStruct) setToken(userID string, token string) *matrixStruct {
	c.cli.SetCredentials(userID, token)
	return c
}

func (c *matrixStruct) setRoomID(roomID string) *matrixStruct {
	c.roomID = roomID
	return c
}

func (c *matrixStruct) send(text string) {
	_, err := c.cli.SendText(c.roomID, text)
	panicerr(err)
}
type rssConfig struct {
	proxy        string
	retryOnError bool
	timeout      int
}

func getRSS(url string, config ...rssConfig) *gofeed.Feed {
	fp := gofeed.NewParser()

	var cfg httpConfig
	var retryOnError bool
	if len(config) != 0 {
		cfg.httpProxy = config[0].proxy
		retryOnError = config[0].retryOnError
		if config[0].timeout != 0 {
			cfg.timeout = config[0].timeout
		}
	}

	var feed *gofeed.Feed
	var err error
	if retryOnError {
		for {
			if err := try(func() {
				content := httpGet(url, cfg).content

				// lg.trace("获取到的内容是:", content)

				feed, err = fp.ParseString(content)
				panicerr(err)
			}).Error; err == nil {
				break
			} else {
				// lg.trace("获取RSS失败:", err)
				sleep(1)
				// lg.trace("重试")
			}
		}
	} else {
		content := httpGet(url, cfg).content

		// lg.trace("获取到的内容是:", content)

		feed, err = fp.ParseString(content)
		panicerr(err)
	}

	return feed
}
type statikFileStruct struct {
	path   string
	fd     http.File
	mode   string
	reader *bufio.Reader
}

func (m *statikFileStruct) readlines() chan string {
	if m.reader == nil {
		m.reader = bufio.NewReader(m.fd)
	}

	lines := make(chan string)

	go func() {
		for {
			line, err := m.reader.ReadString('\n')
			if len(line) == 0 {
				if err != nil {
					close(lines)
					m.close()
					if err == io.EOF {
						return
					}
					_, fn, line, _ := runtime.Caller(0)
					panic(filepath.Base(fn) + ":" + strconv.Itoa(line-7) + " >> " + err.Error())
				}
			}
			line = line[:len(line)-1]
			lines <- line
		}
	}()

	return lines
}

func (m *statikFileStruct) readline() string {
	b := make([]byte, 1)

	line := ""

	for {
		_, err := io.ReadAtLeast(m.fd, b, 1)
		if err != nil {
			if len(line) != 0 {
				return line
			}
			panicerr(err)
		}
		bs := string(b)
		if bs == "\n" {
			return line
		}
		line = line + bs
	}
}

func (m *statikFileStruct) close() {
	m.fd.Close()
}

func (m *statikFileStruct) read(num ...int) string {
	var bytes []byte
	var err error
	if len(num) == 0 {
		bytes, err = ioutil.ReadAll(m.fd)
		panicerr(err)
		m.close()
	} else {
		buffer := make([]byte, num[0])
		_, err := io.ReadAtLeast(m.fd, buffer, num[0])
		panicerr(err)
		bytes = buffer
	}

	return string(bytes)
}

func (m *statikFileStruct) seek(num int64) {
	_, err := m.fd.Seek(num, 0)
	panicerr(err)
}

func statikOpen(path string) *statikFileStruct {
	statikFS, err := fs.New()
	panicerr(err)

	if !strStartsWith(path, "/") {
		path = "/" + path
	}

	fd, err := statikFS.Open(path)
	panicerr(err)

	return &statikFileStruct{
		path: path,
		fd:   fd,
		mode: "r",
	}
}
type godaddyRecord struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
}

type godaddyStruct struct {
	header httpHeader
}

type godaddyDomainStruct struct {
	header     httpHeader
	domainName string
}

type godaddyDomainInfoStruct struct {
	CreatedAt              time.Time   `json:"createdAt"`
	Domain                 string      `json:"domain"`
	DomainID               int         `json:"domainId"`
	ExpirationProtected    bool        `json:"expirationProtected"`
	Expires                time.Time   `json:"expires"`
	ExposeWhois            bool        `json:"exposeWhois"`
	HoldRegistrar          bool        `json:"holdRegistrar"`
	Locked                 bool        `json:"locked"`
	NameServers            interface{} `json:"nameServers"`
	Privacy                bool        `json:"privacy"`
	RenewAuto              bool        `json:"renewAuto"`
	RenewDeadline          time.Time   `json:"renewDeadline"`
	Renewable              bool        `json:"renewable"`
	Status                 string      `json:"status"`
	TransferAwayEligibleAt time.Time   `json:"transferAwayEligibleAt"`
	TransferProtected      bool        `json:"transferProtected"`
}

func getGodaddy(key string, secret string) *godaddyStruct {
	return &godaddyStruct{
		header: httpHeader{
			"Authorization": "sso-key " + key + ":" + secret,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
	}
}

func (m *godaddyStruct) list() (res []godaddyDomainInfoStruct) {
	resp := httpGet("https://api.godaddy.com/v1/domains", m.header)
	if resp.statusCode != 200 {
		panicerr(resp.content)
	}
	err := json.Unmarshal([]byte(resp.content), &res)
	panicerr(err)
	return
}

func (m *godaddyStruct) domain(domainName string) *godaddyDomainStruct {
	return &godaddyDomainStruct{
		header:     m.header,
		domainName: domainName,
	}
}

func (m *godaddyDomainStruct) list() (res []godaddyRecord) {
	resp := httpGet("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", m.header)
	if resp.statusCode != 200 {
		panicerr(resp.content)
	}
	err := json.Unmarshal([]byte(resp.content), &res)
	panicerr(err)
	return
}

// 参数留空字符串则忽略这个项目
func (m *godaddyDomainStruct) delete(name string, dtype string, value string) {
	var records []godaddyRecord
	for _, v := range m.list() {
		if v.Name != name && !(v.Data == "Parked" && v.Type == "A") {
			if name != "" && dtype != "" && value != "" {
				if v.Name != name && v.Type != dtype && v.Data != value {
					records = append(records, v)
				}
			} else if name == "" && dtype != "" && value != "" {
				if v.Type != dtype && v.Data != value {
					records = append(records, v)
				}
			} else if name != "" && dtype == "" && value != "" {
				if v.Name != name && v.Data != value {
					records = append(records, v)
				}
			} else if name != "" && dtype != "" && value == "" {
				if v.Name != name && v.Type != dtype {
					records = append(records, v)
				}
			} else if name != "" && dtype == "" && value == "" {
				if v.Name != name {
					records = append(records, v)
				}
			} else if name == "" && dtype == "" && value != "" {
				if v.Data != value {
					records = append(records, v)
				}
			} else if name == "" && dtype != "" && value == "" {
				if v.Type != dtype {
					records = append(records, v)
				}
			}
		}
	}

	resp := httpPutJSON("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", records, m.header)
	if resp.statusCode != 200 {
		panicerr(resp.content)
	}
}

func (m *godaddyDomainStruct) modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordType string, dstRecordValue string) {
	var records []godaddyRecord
	for _, v := range m.list() {
		if !(v.Data == "Parked" && v.Type == "A") {
			records = append(records, v)
		}
	}

	for idx := range records {
		if records[idx].Name == recordName {
			if records[idx].Type == srcRecordType && records[idx].Data == srcRecordValue {
				records[idx].Type = dstRecordType
				records[idx].Data = dstRecordValue
			}
		}
	}
	resp := httpPutJSON("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", records, m.header)
	if resp.statusCode != 200 {
		panicerr(resp.content)
	}
}

func (m *godaddyDomainStruct) add(recordName string, recordType string, recordValue string) {
	var records []godaddyRecord
	for _, v := range m.list() {
		if !(v.Data == "Parked" && v.Type == "A") {
			records = append(records, v)
		}
	}

	records = append(records, godaddyRecord{
		Data: recordValue,
		Name: recordName,
		TTL:  600,
		Type: strUpper(recordType),
	})

	resp := httpPutJSON("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", records, m.header)
	if resp.statusCode != 200 {
		panicerr(resp.content)
	}
}

// 更新域名的某个记录, 如果不存在则新增
// 如果同一个名字, 多个A记录这样, 应该会出问题
// func (m *godaddyDomainStruct) update(recordName string, recordType string, recordValue string) {
// 	recordType = strUpper(recordType)
// 	var records []godaddyRecord
// 	for _, v := range m.list() {
// 		if !(v.Data == "Parked" && v.Type == "A") {
// 			if !(v.Name == recordName && v.Type == recordType) {
// 				records = append(records, v)
// 			}
// 		}
// 	}

// 	records = append(records, godaddyRecord{
// 		Data: recordValue,
// 		Name: recordName,
// 		TTL:  600,
// 		Type: recordType,
// 	})

// 	resp := httpPutJSON("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", records, m.header)
// 	if resp.statusCode != 200 {
// 		_, fn, line, _ := runtime.Caller(0)
// 		panic(filepath.Base(fn) + ":" + strconv.Itoa(line+1) + " >> " + resp.content + " >> " + fmtDebugStack(string(debug.Stack())))
// 	}
// }
type telegramStruct struct {
	tg     *tgbotapi.BotAPI
	chatid int64
}

func getTelegram(token string) *telegramStruct {
	for {
		//lg.trace("开始初始化")
		bot, err := tgbotapi.NewBotAPI(token)
		if err == nil {
			//lg.trace("初始化成功")
			return &telegramStruct{tg: bot}
		} else {
			//lg.trace("初始化出错:", err)
			sleep(3)
		}
	}
}

func (m *telegramStruct) setChatID(chatid int64) *telegramStruct {
	m.chatid = chatid
	return m
}

func (m *telegramStruct) sendFile(path string) tgbotapi.Message {
	var err error
	var msg tgbotapi.Message
	sleepCount := 10
	for {
		msg, err = m.tg.Send(tgbotapi.NewDocumentUpload(m.chatid, path))
		if err == nil {
			break
		}
		sleep(sleepCount)
		sleepCount = sleepCount * 2
	}
	return msg
}

func (m *telegramStruct) sendImage(path string) tgbotapi.Message {
	var err error
	var msg tgbotapi.Message
	sleepCount := 10
	for {
		msg, err = m.tg.Send(tgbotapi.NewPhotoUpload(m.chatid, path))
		if err == nil {
			break
		}
		sleep(sleepCount)
		sleepCount = sleepCount * 2
	}
	return msg
}

type tgMsgConfig struct {
	parseMode             string
	DisableWebPagePreview bool
	DisableRetryOnError   bool
}

func (m *telegramStruct) send(text string, cfg ...tgMsgConfig) tgbotapi.Message {
	var err error
	var msg tgbotapi.Message

	mm := tgbotapi.NewMessage(m.chatid, text)
	if len(cfg) != 0 {
		mm.ParseMode = cfg[0].parseMode
		mm.DisableWebPagePreview = cfg[0].DisableWebPagePreview
	}

	sleepCount := 10
	for {
		msg, err = m.tg.Send(mm)
		if err == nil {
			break
		} else {
			if len(cfg) != 0 && cfg[0].DisableRetryOnError {
				break
			}
			sleep(sleepCount)
			sleepCount = sleepCount * 2
		}
	}
	return msg
}
func fileType(fpath string) string {
	kind, err := filetype.Match([]byte(open(fpath).read()))

	panicerr(err)

	if kind == filetype.Unknown {
		return ""
	}
	return kind.Extension
}
type tableStruct struct {
	row          [][]interface{}
	header       []string
	maxCellWidth int
}

func getTable(header ...string) *tableStruct {
	return &tableStruct{header: header, maxCellWidth: 0}
}

func (m *tableStruct) setMaxCellWidth(width ...int) {
	if len(width) == 0 {
		m.maxCellWidth = 30
	} else {
		m.maxCellWidth = width[0]
	}
}

func (m *tableStruct) addRow(row ...interface{}) {
	if len(row) != len(m.header) {
		panicerr("添加的数据个数跟表头的个数对不上")
	}
	m.row = append(m.row, row)
}

func (m *tableStruct) render() string {
	if m.maxCellWidth == 0 {
		table := simpletable.New()

		for _, header := range m.header {
			table.Header.Cells = append(table.Header.Cells, &simpletable.Cell{Align: simpletable.AlignLeft, Text: header})
		}

		for _, row := range m.row {
			var cell []*simpletable.Cell
			for _, r := range row {
				cell = append(cell, &simpletable.Cell{
					Align: simpletable.AlignLeft,
					Text:  toString(r),
				})
			}

			table.Body.Cells = append(table.Body.Cells, cell)
		}

		table.SetStyle(simpletable.StyleCompactLite)
		return table.String()
	}

	tabulate := gotabulate.Create(m.row)
	tabulate.SetHeaders(m.header)
	tabulate.SetWrapStrings(true)
	tabulate.SetMaxCellSize(m.maxCellWidth)
	return tabulate.Render("grid")
}
type rabbitConnectionStruct struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func getRabbitMQ(rabbitMQURL string, queueName string) *rabbitConnectionStruct {
	conn, err := amqp.Dial(rabbitMQURL)
	panicerr(err)

	ch, err := conn.Channel()
	panicerr(err)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	return &rabbitConnectionStruct{
		connection: conn,
		channel:    ch,
		queue:      q,
	}
}

func (m *rabbitConnectionStruct) send(data map[string]string) {
	err := m.channel.Publish(
		"",           // exchange
		m.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(map2bin(data)),
		})
	panicerr(err)
}

func (m *rabbitConnectionStruct) recv() chan map[string]string {
	msgs, err := m.channel.Consume(
		m.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	panicerr(err)

	resch := make(chan map[string]string)
	go func() {
		for d := range msgs {
			resch <- bin2map(str(d.Body))
		}
	}()

	return resch
}
var lg *logStruct

func init() {
	lg = getLogger()
	os.Setenv("TZ", "Asia/Hong_Kong")
}