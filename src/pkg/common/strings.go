package common

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alielmi98/golang-otp-auth/pkg/config"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// To snake case : CountryId -> country_id
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
func GenerateOtp() string {
	cfg := config.GetConfig()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	min := int(math.Pow(10, float64(cfg.Otp.Digits-1)))   // 10^d-1 100000
	max := int(math.Pow(10, float64(cfg.Otp.Digits)) - 1) // 999999 = 1000000 - 1 (10^d) -1

	var num = r.Intn(max-min) + min
	return strconv.Itoa(num)
}
