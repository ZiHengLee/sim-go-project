package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type DingDingWriter struct {
	opt  *DingDingOption
	buf  chan *LogRecord
	stop bool
}

func NewDingDingWriter(opt *DingDingOption) (w *DingDingWriter, err error) {
	w = &DingDingWriter{
		opt: opt,
		buf: make(chan *LogRecord, LogBufferLength),
	}
	go w.run()
	return
}

func (w *DingDingWriter) Close() {
	close(w.buf)
	w.stop = true
}

func (w *DingDingWriter) LogWrite(rec *LogRecord) {
	if strings.Index(rec.Message, w.opt.Key) < 0 {
		return
	}
	select {
	case w.buf <- rec:
	default:
	}
}

func (w *DingDingWriter) run() {
	for !w.stop {
		select {
		case rec := <-w.buf:
			if rec == nil {
				break
			} else {
				w.handle(rec)
			}
		}
	}
}

func (w *DingDingWriter) handle(rec *LogRecord) {
	hostname, _ := os.Hostname()
	text := fmt.Sprintf("#### Time: %v \n #### Src: %v \n ### Hostname: %v\n #### Msg: %v\n", rec.Created.Format("2006-01-02 15:04:05"), rec.Source, hostname, rec.Message)
	err := DingDingSendMarkdown(w.opt.Url, w.opt.Key, text, w.opt.AtMobiles)
	if err != nil {
		fmt.Printf("dingding send msg:%v err:%v", text, err)
	}
}

// 机器人文档
// https://open-doc.dingtalk.com/docs/doc.htm?spm=a219a.7629140.0.0.karFPe&treeId=257&articleId=105735&docType=1

const (
	DINGDING_MSG_TYPE_TEXT     = "text"
	DINGDING_MSG_TYPE_LINK     = "link"
	DINGDING_MSG_TYPE_MARKDOWN = "markdown"
)

type DingDingMsgSt struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Link struct {
		Text       string `json:"text"`
		Title      string `json:"title"`
		PicURL     string `json:"picUrl"`
		MessageURL string `json:"messageUrl"`
	} `json:"link"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

func DingDingSendMsg(param *DingDingMsgSt, URL string) (err error) {
	buf, err := json.Marshal(param)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", URL, bytes.NewReader(buf))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var respData struct {
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return
	}
	if respData.ErrCode != 0 {
		err = fmt.Errorf("dingding resp: %s", body)
	}
	return
}

func DingDingSendText(URL, content string, atMobiles []string) (err error) {
	param := &DingDingMsgSt{
		MsgType: DINGDING_MSG_TYPE_TEXT,
	}
	param.Text.Content = content
	param.At.AtMobiles = atMobiles
	return DingDingSendMsg(param, URL)
}

func DingDingSendTextWithAtAll(URL, content string, atMobiles []string, isAtAll bool) (err error) {
	param := &DingDingMsgSt{
		MsgType: DINGDING_MSG_TYPE_TEXT,
	}
	param.Text.Content = content
	param.At.AtMobiles = atMobiles
	param.At.IsAtAll = isAtAll
	return DingDingSendMsg(param, URL)
}

func DingDingSendLink(URL, title, text, msgURL, picURL string, atMobiles []string) (err error) {
	param := &DingDingMsgSt{
		MsgType: DINGDING_MSG_TYPE_LINK,
	}
	param.Link.Title = title
	param.Link.Text = text
	param.Link.MessageURL = msgURL
	param.Link.PicURL = picURL
	param.At.AtMobiles = atMobiles
	return DingDingSendMsg(param, URL)
}

func DingDingSendMarkdown(URL, title, text string, atMobiles []string) (err error) {
	param := &DingDingMsgSt{
		MsgType: DINGDING_MSG_TYPE_MARKDOWN,
	}
	param.Markdown.Title = title
	param.Markdown.Text = text
	param.At.AtMobiles = atMobiles
	return DingDingSendMsg(param, URL)
}
