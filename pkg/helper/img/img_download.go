package img

import (
	"github.com/pkg/errors"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"
)

func DownPic(src, dest string) (string, error) {
	re, err := http.Get(src)
	if err != nil {
		return "", err
	}
	defer re.Body.Close()
	fix := ""
	if idx := strings.LastIndex(src, "."); idx != -1 {
		fix = strings.ToLower(src[idx+1:])
	}
	if fix == "" {
		return "", errors.Errorf("unknow pic type, pic path: %s", src)
	}
	thumbF, err := os.OpenFile(dest+"."+fix, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer thumbF.Close()
	if fix == "jpeg" || fix == "jpg" {
		img, err := jpeg.Decode(re.Body)
		if err != nil {
			return "", err
		}
		if err = jpeg.Encode(thumbF, img, &jpeg.Options{Quality: 40}); err != nil {
			return "", err
		}
	} else if fix == "png" {
		img, err := png.Decode(re.Body)
		if err != nil {
			return "", err
		}
		if err = png.Encode(thumbF, img); err != nil {
			return "", err
		}
	} else if fix == "gif" {
		img, err := gif.Decode(re.Body)
		if err != nil {
			return "", err
		}
		if err = gif.Encode(thumbF, img, nil); err != nil {
			return "", err
		}
	} else {
		return "", errors.New("不支持的格式")
	}
	return "." + fix, nil
}
