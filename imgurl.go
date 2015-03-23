package imgurl

import (
	"github.com/nfnt/resize"
	"net/http"
	"image/jpeg"
	"image/png"
	"image/gif"
	"image"
	"fmt"
	"errors"
	"encoding/base64"
	"bytes"
	"io"
	"log"
)

type Filter func(image.Image) (image.Image,interface{})

// UrilfyC fetches the image referenced by the given url, scales it to the given sizes keeping the aspect ratio
// and transcods it to a base64 encoded data url.
// Hereby the given http Client is used.
func UrlifyC(c *http.Client, url string, maxwidth,maxheight int,filters ...Filter) (ret string, tags []interface{}, err error) {
	resp, err := c.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	mt := resp.Header["Content-Type"]

	return UrlifyR(resp.Body,mt[0],maxwidth,maxheight,filters...)
}


// UrlifyR reads the image fromt he given reader, scales it to the given sizes keeping the apsect ratio
// and transcods it to a base74 encoded data url.
func UrlifyR(source io.Reader,mt string, maxwidth, maxheight int,filters ...Filter) (ret string, tags []interface{}, err error) {
	defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("Panic in Urlify: %s",r))
				log.Println(err)
			}
		}()
	img,err := Decode(source,mt,maxwidth,maxheight)
	if err != nil {
		return
	}

	tags = make([]interface{},len(filters))
	for i,f := range filters {
		img,tags[i] = f(img)
	}

	ret,err = encode(img)
	return
}

// Urilfy fetches the image referenced by the given url, scales it to the given sizes keeping the aspect ratio
// and transcods it to a base64 encoded data url.
func Urlify(url string, maxwidth,maxheight int,filters ...Filter) (ret string, tags []interface{}, err error) {
	return UrlifyC(http.DefaultClient,url,maxwidth,maxheight,filters...)
}

// Decode reads the given image and scales it to the given size keeping the aspect ratio
func Decode(source io.Reader,mt string,mwh ...int) (i image.Image, err error) {
	var img image.Image
	switch mt {
		case "image/jpeg","image/jpg":
			img, err = jpeg.Decode(source)
		case "image/png":
			img, err = png.Decode(source)
		case "image/gif":
			img, err = gif.Decode(source)
		default:
			return nil,errors.New(fmt.Sprintf("Unsupported content type: %s",mt))
	}
	if len(mwh) > 0{
		var mw,mh int
		mw = mwh[0]
		if len(mwh) > 1 {
			mh = mwh[1]
		}

		if (mw > 0 && mh > 0) {
			return resize.Thumbnail(uint(mw),uint(mh),img,resize.Bilinear), err
		} else {
			return img,err
		}
	} else {
		return img, err
	}
}

// encode ecnodes the image into a base64 png data url.
func encode(source image.Image) (ret string, err error) {
	buf:= new(bytes.Buffer)
	err = png.Encode(buf,source)
	if err != nil {
		return
	}

	ret = fmt.Sprintf("data:image/png;base64,%s",base64.StdEncoding.EncodeToString(buf.Bytes()))
	return
}
