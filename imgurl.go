package imgurl

import (
	. "github.com/nfnt/resize"
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
)

// Urilfy fetches the image referenced by the given url, scales it to the given sizes keeping the aspect ratio
// and transcods it to a base64 encoded data url.
func Urlify(url string, maxwidth,maxheight int) (ret string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	mt := resp.Header["Content-Type"]

	return transcode(resp.Body,mt[0],maxwidth,maxheight)
}

// transcode reads the given image, scales it to the given size keeping the aspect rati
// and transcods it to a bas64 encoded data url.
func transcode(source io.Reader,mt string,maxwidth,maxheight int) (ret string, err error) {
	var img image.Image
	switch mt {
		case "image/jpeg","image/jpg":
			img, err = jpeg.Decode(source)
		case "image/png":
			img, err = png.Decode(source)
		case "image/gif":
			img, err = gif.Decode(source)
		default:
			return "",errors.New(fmt.Sprintf("Unsupported content type: %s",mt))
	}

	scaled := Thumbnail(uint(maxwidth),uint(maxheight),img,Bilinear)

	buf:= new(bytes.Buffer)

	err = png.Encode(buf,scaled)
	if err != nil {
		return
	}

	ret = fmt.Sprintf("data:image/png;base64,%s",base64.StdEncoding.EncodeToString(buf.Bytes()))
	return
}
