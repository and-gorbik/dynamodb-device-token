package region

var regions = map[string]struct{}{
	"eu-central-1": {},
	"us-east-1":    {},
	"sa-east-1":    {},
}

func In(val string) bool {
	_, found := regions[val]
	return found
}

func Default() string {
	return "us-east-1"
}
