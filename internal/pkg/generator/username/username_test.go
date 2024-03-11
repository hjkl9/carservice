package username

import (
	"fmt"
	"testing"
)

func TestGenerateHexById(t *testing.T) {
	result := GenerateHexById(1)
	if result != "00000001" {
		t.Fatalf("incorrect result")
	}
	fmt.Println("It's correct")
}
