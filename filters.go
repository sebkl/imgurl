package imgurl

import(
	. "github.com/sebkl/go-nude"
	"image"
)

// NudeFilter checks if the image may be a nude image and returns a boolean
// tag accordingly.
func NudeFilter(img image.Image) (image.Image,interface{}) {
	nude,_ := IsImageNude(img)
	return img,nude
}
