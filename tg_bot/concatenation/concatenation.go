package concatenation

import (
	"bytes"
	"context"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"image"
	"image/color"
	"net/http"
	"time"
)

const separatorWidth = 10

var (
	Version        = "v1"
	separatorColor = color.White
	httpClient     = &http.Client{Timeout: 3 * time.Second}
)

func Concatenate(ctx context.Context, leftPictureUrl, rightPictureUrl string) (*bytes.Buffer, error) {
	leftPicture, err := downloadImage(leftPictureUrl)
	if err != nil {
		return nil, err
	}
	rightPicture, err := downloadImage(rightPictureUrl)
	if err != nil {
		return nil, err
	}
	leftPicture, rightPicture = ensureSameHeight(leftPicture, rightPicture)
	resultImg := appendHorizontally(leftPicture, rightPicture)
	buf := &bytes.Buffer{}
	err = imaging.Encode(buf, resultImg, imaging.JPEG)
	if err != nil {
		return nil, errors.Wrap(err, "concatenation failed")
	}
	return buf, nil
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
		return crop(leftPicture, rightPictureHeight), rightPicture
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
	resp, err := httpClient.Get(pictureUrl)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("download %s: not 200 resp code: %d", pictureUrl, resp.StatusCode)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "download %s", pictureUrl)
	}
	img, err := imaging.Decode(resp.Body)
	return img, errors.Wrapf(err, "download %s: can't open image", pictureUrl)
}
