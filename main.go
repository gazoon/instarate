package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func init() {
	raven.SetDSN("https://dd153cdf7a6e4e5a90b8368142cc5e42:432f707ddc87488ca451603468d7b740@sentry.edwin.ai/7")
}

func f1() error {
	//var a *string
	//println(*a)
	return errors.New("ffff")
}
func f3() error {
	return f1()
}
func f4() error {
	return f3()
}
func f5() error {
	return f4()
}
func f6() error {
	var a = 1
	fmt.Println(a)
	return f5()
}
func f7() error {
	return f6()
}
func f8() error {
	return f7()
}
func f9() error {
	return f8()
}
func f2() {
	defer func() {
		if err, ok := recover().(error); ok {
			handle(err)
		}
	}()
	err := f9()
	panic(err)
}

func handle(err error) {
	log.Errorf("%+v", err)
	debug.PrintStack()
	raven.CaptureErrorAndWait(err, map[string]string{"request_id": "2228"})
}

type A struct {
	X string
	Y int
}

//func (self A) String() string {
//	return utils.ObjToString(&self)
//}

func main() {
	a := &A{"sss", 2}
	fmt.Printf("%v\n", a)
	log.WithField("struct", a).Info("kek")
	fmt.Println(2*time.Hour - 2*time.Minute)
	fmt.Println(len(strings.Fields("   ")))
	fmt.Println(strings.Fields("https://www.instagram.com/p/BkISd3KlaEyey8yrI1dR-n4d9Jf5Uf6mRc50dk0/?hl=ru&taken-by=nik_nik0n0v"))
	fmt.Println(len(strings.Split("s", "s")))
	_, err := strconv.Atoi("")
	fmt.Println(err)

}
