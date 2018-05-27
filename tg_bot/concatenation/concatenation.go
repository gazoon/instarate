package concatenation

import (
	"image"
	"image/color"
	"net/http"
	"strings"

	"bytes"
	"github.com/disintegration/imaging"
	"github.com/i/paralyze"
	"github.com/pkg/errors"
)

const separatorWidth = 10

var separatorColor = color.White

func Concatenate(leftPictureUrl, rightPictureUrl string) (*bytes.Buffer, error) {
	images, err := paralyzeTasks(
		func() (interface{}, error) { return downloadImage(leftPictureUrl) },
		func() (interface{}, error) { return downloadImage(rightPictureUrl) },
	)
	if err != nil {
		return nil, errors.Wrap(err, "concatenation failed")
	}
	leftPicture, rightPicture := images[0].(image.Image), images[1].(image.Image)
	leftPicture, rightPicture = ensureSameHeight(leftPicture, rightPicture)
	resultImg := appendHorizontally(leftPicture, rightPicture)
	buf := &bytes.Buffer{}
	err = imaging.Encode(buf, resultImg, imaging.JPEG)
	if err != nil {
		return nil, errors.Wrap(err, "concatenation failed")
	}
	return buf, nil
}

func paralyzeTasks(funcs ...paralyze.Paralyzable) ([]interface{}, error) {
	results, errs := paralyze.Paralyze(funcs...)

	var errorMsgs []string
	for _, e := range errs {
		if e == nil {
			continue
		}
		errorMsgs = append(errorMsgs, e.Error())
	}
	if errorMsgs != nil {
		return nil, errors.New(strings.Join(errorMsgs, "; "))
	}
	return results, nil
}

func appendHorizontally(leftPicture, rightPicture image.Image) image.Image {
	rightPicturePosX := getWidth(leftPicture) + separatorWidth
	resultWidth := rightPicturePosX + getWidth(rightPicture)
	resultHeight := getHeight(leftPicture)
	resultImg := imaging.New(resultWidth, resultHeight, separatorColor)
	resultImg = imaging.Paste(resultImg, leftPicture, image.Point{0, 0})
	resultImg = imaging.Paste(resultImg, rightPicture, image.Point{rightPicturePosX, 0})
	return resultImg
}

func ensureSameHeight(leftPicture, rightPicture image.Image) (image.Image, image.Image) {
	leftPictureHeight := getHeight(leftPicture)
	rightPictureHeight := getHeight(rightPicture)
	if leftPictureHeight == rightPictureHeight {
		return leftPicture, rightPicture
	} else if leftPictureHeight < rightPictureHeight {
		return leftPicture, crop(rightPicture, leftPictureHeight)
	} else {
		return rightPicture, crop(leftPicture, rightPictureHeight)
	}
}

func crop(picture image.Image, resultHeight int) image.Image {
	originalWidth := getWidth(picture)
	return imaging.CropCenter(picture, originalWidth, resultHeight)
}

func getWidth(picture image.Image) int {
	return picture.Bounds().Max.X
}

func getHeight(picture image.Image) int {
	return picture.Bounds().Max.Y
}

func downloadImage(pictureUrl string) (image.Image, error) {
	resp, err := http.Get(pictureUrl)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("download %s: not 200 resp code: %d", pictureUrl, resp.StatusCode)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "download %s", pictureUrl)
	}
	img, err := imaging.Decode(resp.Body)
	return img, errors.Wrapf(err, "download %s: can't open image", pictureUrl)
}
