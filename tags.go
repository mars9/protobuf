package protobuf

import "strings"

const (
	optional = iota + 1
	required

	ftypeStart
	sfixed32
	sfixed64
	fixed32
	fixed64
	sint32
	sint64
	ftypeEnd
)

// `protobuf:"fixed64,required"`
// `protobuf:"fixed64,optional"`
func parseTag(tag string) (typ, field int) {
	fields := strings.Split(tag, ",")
	for i := range fields {
		switch fields[i] {
		case "sint32":
			typ = sint32
		case "sint64":
			typ = sint64
		case "fixed64":
			typ = fixed64
		case "sfixed64":
			typ = sfixed64
		case "fixed32":
			typ = fixed32
		case "sfixed32":
			typ = sfixed32
		case "optional":
			field = optional
		case "required":
			field = required
		}
	}
	return typ, field
}
