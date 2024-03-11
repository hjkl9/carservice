package captcha

import (
	"fmt"
	"math/rand"
	"time"
)

func PhoneNumberCaptcha() string {
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	return fmt.Sprintf("%06v", rand.Int31n(1000000))
}
