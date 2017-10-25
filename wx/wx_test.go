package wx

import (
	"testing"
	"regexp"
	"fmt"
)

func TestSync(t *testing.T){
	body := `window.synccheck={retcode:"0",selector:"2"}`
	reg := regexp.MustCompile(`window\.synccheck\=\{retcode:\"(\d+)\",selector\:\"(\d+)\"\}`)
	td := reg.FindStringSubmatch(body)
	fmt.Println(td)

}


