package log

import (
	"log"

	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/color"
)

func Fatalf(format string, v ...any) {
	log.Fatalf(color.Red+format+color.Reset, v...)
}

func Fatalln(v ...any) {
	log.Fatalln([]any{color.Red, v, color.Reset}...)
}

func Errorf(format string, v ...any) {
	log.Printf(color.Red+format+color.Reset, v...)
}

func Errorln(v ...any) {
	log.Println([]any{color.Red, v, color.Reset}...)
}

func Warnf(format string, v ...any) {
	log.Printf(color.Yellow+format+color.Reset, v...)
}

func Warnln(v ...any) {
	log.Println([]any{color.Yellow, v, color.Reset}...)
}

func Infof(format string, v ...any) {
	log.Printf(format, v...)
}

func Infoln(v ...any) {
	log.Println(v...)
}
