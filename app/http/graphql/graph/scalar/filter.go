package scalar

import (
	"encoding/json"
	"io"
	"log"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalFilter :
func MarshalFilter(v json.RawMessage) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		// fmt.Fprint(w, strconv.FormatUint(uint64(v), 10))
	})
}

// UnmarshalFilter :
func UnmarshalFilter(v interface{}) (json.RawMessage, error) {
	switch vi := v.(type) {
	case map[string]interface{}:
		log.Println("filter =>", vi)
		return nil, nil

	default:
		return json.RawMessage(`{}`), nil
	}
}
