package imgurl


import (
	"testing"
	"regexp"
	"net/http"
	"io/ioutil"
	"bytes"
)

const (
	TESTIMAGE = "https://raw.githubusercontent.com/sebkl/globejs/master/screenshots/sample_plain.png"
	TESTSIZE = 100
)


func check(in string) (ret bool) {
	ret,_ = regexp.MatchString("data:.*,.+",in)
	return
}

func TestBasic(t *testing.T) {
	resp, err := Urlify(TESTIMAGE,TESTSIZE,TESTSIZE)
	if err != nil {
		t.Errorf("Image processing failed: %s",err)
	}

	if !check(resp) {
		t.Errorf("Image processing failed: %s,%s",resp,err)
	}
	t.Logf("%s",resp)
}

func TestTranscodeService(t *testing.T) {
	ts := NewTranscodeService(3)

	for i := 0;i < 5; i++ {
		ts.in <- &Request{url: TESTIMAGE,maxwidth: TESTSIZE,maxheight: TESTSIZE}
	}

	for i := 0;i < 5; i++ {
		resp := <-ts.out
		if !check(resp.image) {
			t.Errorf("Image processing failed: %s",resp.image)
		}
	}
}

func BenchmarkTranscode(b *testing.B) {
	resp, err := http.Get(TESTIMAGE)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	mt := resp.Header["Content-Type"]

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(body)
		transcode(buf,mt[0],TESTSIZE,TESTSIZE)
	}
}
