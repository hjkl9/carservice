package order

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateNumber(t *testing.T) {
	orderNumber := GenerateNumber(time.Now())
	fmt.Println(orderNumber)
}
