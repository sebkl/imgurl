package imgurl


import (
	"testing"
	"regexp"
)

func TestBasic(t *testing.T) {
	resp, err := Urlify("https://raw.githubusercontent.com/sebkl/globejs/master/screenshots/sample_plain.png",100,100)
	if err != nil {
		t.Errorf("Image processing failed: %s",err)
	}

	if matched, err := regexp.MatchString("data:.*,.+",resp); !matched {
		t.Errorf("Image processing failed: %s,%s",resp,err)
	}
	t.Logf("%s",resp)
}



