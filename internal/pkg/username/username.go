package username

import "fmt"

func GenerateHexById(id uint) string {
	return fmt.Sprintf("%08x", id)
}
