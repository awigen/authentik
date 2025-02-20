package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nmcclain/ldap"
	ldapConstants "goauthentik.io/internal/outpost/ldap/constants"
)

func ldapResolveTypeSingle(in interface{}) *string {
	switch t := in.(type) {
	case string:
		return &t
	case *string:
		return t
	case bool:
		s := strconv.FormatBool(t)
		return &s
	case float32:
		s := strconv.FormatFloat(float64(t), 'f', -1, 64)
		return &s
	case float64:
		s := strconv.FormatFloat(t, 'f', -1, 64)
		return &s
	case int:
		s := strconv.FormatInt(int64(t), 10)
		return &s
	default:
		if in != nil {
			s := fmt.Sprintf("%s", in)
			return &s
		}
		return nil
	}
}

func AKAttrsToLDAP(attrs map[string]interface{}) []*ldap.EntryAttribute {
	attrList := []*ldap.EntryAttribute{}
	if attrs == nil {
		return attrList
	}
	for attrKey, attrValue := range attrs {
		entry := &ldap.EntryAttribute{Name: attrKey}
		switch t := attrValue.(type) {
		case []string:
			entry.Values = t
		case *[]string:
			entry.Values = *t
		case []interface{}:
			entry.Values = make([]string, len(t))
			for idx, v := range t {
				v := ldapResolveTypeSingle(v)
				if v != nil {
					entry.Values[idx] = *v
				}
			}
		default:
			v := ldapResolveTypeSingle(t)
			if v != nil {
				entry.Values = []string{*v}
			}
		}
		attrList = append(attrList, entry)
	}
	return attrList
}

func EnsureAttributes(attrs []*ldap.EntryAttribute, shouldHave map[string][]string) []*ldap.EntryAttribute {
	for name, values := range shouldHave {
		attrs = MustHaveAttribute(attrs, name, values)
	}
	return attrs
}

func MustHaveAttribute(attrs []*ldap.EntryAttribute, name string, value []string) []*ldap.EntryAttribute {
	shouldSet := true
	for _, attr := range attrs {
		if attr.Name == name {
			shouldSet = false
		}
	}
	if shouldSet {
		return append(attrs, &ldap.EntryAttribute{
			Name:   name,
			Values: value,
		})
	}
	return attrs
}

func IncludeObjectClass(searchOC string, ocs map[string]bool) bool {
	if searchOC == "" {
		return true
	}

	return ocs[searchOC]
}

func GetContainerEntry(filterOC string, dn string, ou string) *ldap.Entry {
	if IncludeObjectClass(filterOC, ldapConstants.GetContainerOCs()) {
		return &ldap.Entry{
			DN: dn,
			Attributes: []*ldap.EntryAttribute{
				{
					Name:   "distinguishedName",
					Values: []string{dn},
				},
				{
					Name:   "objectClass",
					Values: []string{"top", "nsContainer"},
				},
				{
					Name:   "commonName",
					Values: []string{ou},
				},
			},
		}
	}

	return nil
}

func HasSuffixNoCase(s1 string, s2 string) bool {
	return strings.HasSuffix(strings.ToLower(s1), strings.ToLower(s2))
}
