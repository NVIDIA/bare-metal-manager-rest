package util

import (
	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

// GetStrPtr returns a pointer to the string passed in
func GetStrPtr(s string) *string {
	return &s
}

func ProtobufUUIDListToStringList(ids []*cwssaws.UUID) []string {
	s := make([]string, len(ids))

	for i, u := range ids {
		if u == nil {
			s[i] = ""
		} else {
			s[i] = u.Value
		}
	}

	return s
}

func StringsToProtobufUUIDList(ids []string) []*cwssaws.UUID {
	s := make([]*cwssaws.UUID, len(ids))

	for i, u := range ids {
		s[i] = &cwssaws.UUID{Value: u}
	}

	return s
}
