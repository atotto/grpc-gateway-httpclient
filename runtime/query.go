package runtime

import (
	"net/url"

	"github.com/gorilla/schema"
)

func toQueryURL(encoder *schema.Encoder, u string, src interface{}) string {
	ur, _ := url.Parse(u)
	q := ur.Query()

	encoder.Encode(src, q)

	ur.RawQuery = q.Encode()
	return ur.String()
}
