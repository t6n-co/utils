package metrics

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

type Tag struct {
	Name  string
	Value interface{}
}

// attrsFrom Convert Tag slice to OTel attributes.
func attrsFrom(tags ...Tag) []attribute.KeyValue {
	if len(tags) == 0 {
		return nil
	}
	out := make([]attribute.KeyValue, 0, len(tags))
	for _, t := range tags {
		switch v := t.Value.(type) {
		case string:
			out = append(out, attribute.String(t.Name, v))
		case fmt.Stringer:
			out = append(out, attribute.String(t.Name, v.String()))
		case bool:
			out = append(out, attribute.Bool(t.Name, v))
		case int:
			out = append(out, attribute.Int(t.Name, v))
		case int32:
			out = append(out, attribute.Int64(t.Name, int64(v)))
		case int64:
			out = append(out, attribute.Int64(t.Name, v))
		case uint:
			out = append(out, attribute.Int64(t.Name, int64(v)))
		case uint32:
			out = append(out, attribute.Int64(t.Name, int64(v)))
		case uint64:
			out = append(out, attribute.Int64(t.Name, int64(v)))
		case float32:
			out = append(out, attribute.Float64(t.Name, float64(v)))
		case float64:
			out = append(out, attribute.Float64(t.Name, v))
		case []string:
			out = append(out, attribute.StringSlice(t.Name, v))
		default:
			out = append(out, attribute.String(t.Name, fmt.Sprint(v)))
		}
	}
	return out
}
