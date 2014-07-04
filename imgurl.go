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
)

func Urlify(url string, maxwidth,maxheight int) (ret string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var img image.Image
	mt := resp.Header["Content-Type"]

	switch mt[0] {
		case "image/jpeg","image/jpg":
			img, err = jpeg.Decode(resp.Body)
		case "image/png":
			img, err = png.Decode(resp.Body)
		case "image/gif":
			img, err = gif.Decode(resp.Body)
		default:
			return "",errors.New(fmt.Sprintf("Unsupported content type: %s",mt[0]))
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
