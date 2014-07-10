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
	resp, _, err := Urlify(TESTIMAGE,TESTSIZE,TESTSIZE)
	if err != nil {
		t.Errorf("Image processing failed: %s",err)
	}

	if !check(resp) {
		t.Errorf("Image processing failed: %s,%s",resp,err)
	}
	t.Logf("%s",resp)
}

func TestTranscodeService(t *testing.T) {
	ts := NewTranscodeService(3,5)

	for i := 0;i < 5; i++ {
		ts.in <- &Request{Url: TESTIMAGE,Maxwidth: TESTSIZE,Maxheight: TESTSIZE}
	}

	for i := 0;i < 5; i++ {
		resp := <-ts.out
		if !check(resp.Image) {
			t.Errorf("Image processing failed: %s",resp.Image)
		}
	}
}

func TestTrasncodeServiceOverload(t *testing.T) {
	ts := NewTranscodeService(2,2)
	ts.Push( &Request{Url: TESTIMAGE,Maxwidth: TESTSIZE,Maxheight: TESTSIZE} )
	ts.Push( &Request{Url: TESTIMAGE,Maxwidth: TESTSIZE,Maxheight: TESTSIZE} )
	if !ts.Full() {
		t.Errorf("Request queue overloaded.")
	}
}

func TestTranscodeServiceFilter(t *testing.T) {
	ts := NewTranscodeService(2,2)
	filters := []Filter{NudeFilter}
	ts.Push( &Request{Url: TESTIMAGE,Maxwidth: TESTSIZE,Maxheight: TESTSIZE, Filters: filters} )
	resp := ts.Get()
	if len(resp.Tags) <= 0 {
		t.Errorf("Filter did not return tag.")
	}

	if r,ok := resp.Tags[0].(bool); !ok || r  {
		t.Errorf("Filter did not detect correct.")
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
		img,_:= Decode(buf,mt[0],TESTSIZE,TESTSIZE)
		encode(img)
	}
}
