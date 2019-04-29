package bindec

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Constraint to be checked on a struct field.
type Constraint interface {
	// BeforeRead returns whether the constraint should be checked before
	// reading the content. This is only applicable to slices and strings.
	BeforeRead() bool
	// Validator generates the code to validate the receiver with the current
	// constraint.
	Validator(recv string) string
}

type stringConstraint struct {
	field    string
	isPtr    bool
	template string
}

func (c stringConstraint) BeforeRead() bool { return false }

func (c stringConstraint) Validator(recv string) string {
	var prefix string
	if c.isPtr {
		prefix = "*"
	}
	return fmt.Sprintf(c.template, prefix+recv, c.field)
}

type argConstraint struct {
	field    string
	isPtr    bool
	arg      string
	template string
}

func (c argConstraint) BeforeRead() bool { return false }

func (c argConstraint) Validator(recv string) string {
	var prefix string
	if c.isPtr {
		prefix = "*"
	}
	return fmt.Sprintf(c.template, prefix+recv, c.field, c.arg)
}

type oneOf struct {
	field string
	isPtr bool
	args  []string
}

func (c oneOf) BeforeRead() bool { return false }

func (c oneOf) Validator(recv string) string {
	var prefix string
	if c.isPtr {
		prefix = "*"
	}

	var parts = make([]string, len(c.args))
	for i, a := range c.args {
		parts[i] = fmt.Sprintf("%s%s != %s", prefix, recv, a)
	}

	return fmt.Sprintf(`if %s {
	return fmt.Errorf("field '%s' should have one of these values: %%s", %q)
}`, strings.Join(parts, " && "), c.field, strings.Join(c.args, ", "))
}

type lenConstraint struct {
	field    string
	len      int
	template string
}

func (c lenConstraint) BeforeRead() bool { return true }

func (c lenConstraint) Validator(_ string) string {
	return fmt.Sprintf(c.template, c.len, c.field)
}

func parseConstraint(
	ctx *parseContext,
	name, args string,
	field string, typ Type,
) (Constraint, error) {
	ctx.addImport("fmt")
	for _, i := range constraintImports[name] {
		ctx.addImport(i)
	}

	if d, ok := constraintDecls[name]; ok {
		ctx.addDecl(d)
	}

	ptr := isPtr(typ)
	switch name {
	case "alpha",
		"alphanum",
		"numeric",
		"hexadecimal",
		"email",
		"url",
		"base64",
		"uuid",
		"ip",
		"ipv4",
		"ipv6":
		if !isString(typ) {
			return nil, fmt.Errorf("constraint %q can only be used on string or *string fields", name)
		}

		return stringConstraint{field, ptr, constraintTemplates[name]}, nil
	case "contains",
		"startswith",
		"endswith":
		if !isString(typ) {
			return nil, fmt.Errorf("constraint %q can only be used on string or *string fields", name)
		}

		return argConstraint{field, ptr, toPrintableValue(args, typ), constraintTemplates[name]}, nil
	case "oneof":
		if !isBasic(typ) {
			return nil, fmt.Errorf("oneof can only be used with basic types")
		}

		var options []string
		for _, a := range strings.Split(args, " ") {
			a = strings.TrimSpace(a)
			if a != "" {
				if !isValueOfType(a, typ) {
					return nil, fmt.Errorf("oneof value %q is not a valid value for the field type", a)
				}
				options = append(options, toPrintableValue(a, typ))
			}
		}

		return oneOf{field, ptr, options}, nil
	case "max", "min":
		if !isNumber(typ) {
			return nil, fmt.Errorf("constraint %q can only be used on numeric fields", name)
		}

		if !isValueOfType(args, typ) {
			return nil, fmt.Errorf("%s value %q is not a valid value for the field type", name, args)
		}

		return argConstraint{field, ptr, args, constraintTemplates[name]}, nil
	case "eq", "neq":
		if !isBasic(typ) {
			return nil, fmt.Errorf("constraint %s can only be used with basic types", name)
		}

		if !isValueOfType(args, typ) {
			return nil, fmt.Errorf("%s value %q is not a valid value for the field type", name, args)
		}

		return argConstraint{field, ptr, toPrintableValue(args, typ), constraintTemplates[name]}, nil
	case "maxlen", "minlen":
		if !isString(typ) && !isSlice(typ) {
			return nil, fmt.Errorf("constraint %q can only be used on string and slice fields", name)
		}

		n, err := strconv.Atoi(args)
		if err != nil {
			return nil, fmt.Errorf("constraint %q value %q is not a valid number", name, args)
		}

		return lenConstraint{field, n, constraintTemplates[name]}, nil
	default:
		return nil, fmt.Errorf("constraint not found: %s", name)
	}
}

func isString(t Type) bool {
	switch t := t.(type) {
	case Basic:
		return t.Kind == String
	case Maybe:
		b, ok := t.Elem.(Basic)
		if !ok {
			return false
		}

		return b.Kind == String
	default:
		return false
	}
}

func isNumber(t Type) bool {
	switch t := t.(type) {
	case Basic:
		return t.Kind >= Int && t.Kind <= Float64
	case Maybe:
		b, ok := t.Elem.(Basic)
		if !ok {
			return false
		}

		return b.Kind >= Int && b.Kind <= Float64
	default:
		return false
	}
}

func isBasic(t Type) bool {
	switch t := t.(type) {
	case Basic:
		return true
	case Maybe:
		_, ok := t.Elem.(Basic)
		return ok
	default:
		return false
	}
}

func isSlice(t Type) bool {
	switch t := t.(type) {
	case Bytes:
		return true
	case Slice:
		return true
	case Maybe:
		switch t.Elem.(type) {
		case Bytes:
			return true
		case Slice:
			return true
		}
	}
	return false
}

func isPtr(typ Type) bool {
	_, ok := typ.(Maybe)
	return ok
}

// constraints is a map between the constraint name and whether it requires
// arguments.
var constraints = map[string]bool{
	"alpha":       false,
	"alphanum":    false,
	"numeric":     false,
	"hexadecimal": false,
	"email":       false,
	"url":         false,
	"base64":      false,
	"contains":    true,
	"startswith":  true,
	"endswith":    true,
	"eq":          true,
	"neq":         true,
	"uuid":        false,
	"ip":          false,
	"ipv4":        false,
	"ipv6":        false,
	"oneof":       true,
	"max":         true,
	"min":         true,
	"maxlen":      true,
	"minlen":      true,
}

var constraintTemplates = map[string]string{
	"alpha":       alphaTpl,
	"alphanum":    alphanumTpl,
	"numeric":     numericTpl,
	"hexadecimal": hexadecimalTpl,
	"email":       emailTpl,
	"url":         urlTpl,
	"base64":      base64Tpl,
	"uuid":        uuidTpl,
	"ip":          ipTpl,
	"ipv4":        ipv4Tpl,
	"ipv6":        ipv6Tpl,
	"max":         maxTpl,
	"min":         minTpl,
	"maxlen":      maxlenTpl,
	"minlen":      minlenTpl,
	"eq":          eqTpl,
	"neq":         neqTpl,
	"contains":    containsTpl,
	"startswith":  startsWithTpl,
	"endswith":    endsWithTpl,
}

const (
	alphaTpl = `for _, ru := range %[1]s {
	if !unicode.IsLetter(ru) {
		return fmt.Errorf("field '%%v' contains non alpha characters", %[2]q)
	}
}`
	alphanumTpl = `for _, ru := range %[1]s {
	if !unicode.IsLetter(ru) && !unicode.IsDigit(ru) {
		return fmt.Errorf("field '%%v' contains non alphanumeric characters", %[2]q)
	}
}`
	numericTpl = `for _, ru := range %[1]s {
	if !unicode.IsDigit(ru) {
		return fmt.Errorf("field '%%v' contains non numeric characters", %[2]q)
	}
}`
	hexadecimalTpl = `if !hexadecimalConstraintRegex.MatchString(%[1]s) {
	return fmt.Errorf("field '%%v' is not a valid hexadecimal string", %[2]q)
}`
	emailTpl = `if !emailConstraintRegex.MatchString(%[1]s) {
	return fmt.Errorf("field '%%v' is not a valid email", %[2]q)
}`
	urlTpl = `if url, err := url.ParseRequestURI(%[1]s); err != nil || url.Scheme == "" {
	return fmt.Errorf("field '%%v' is not a valid URL", %[2]q)
}`
	base64Tpl = `if !base64ConstraintRegex.MatchString(%[1]s) {
	return fmt.Errorf("field '%%v' is not a valid base64 string", %[2]q)
}`
	uuidTpl = `if !uuidConstraintRegex.MatchString(%[1]s) {
	return fmt.Errorf("field '%%v' is not a valid UUID", %[2]q)
}`
	ipTpl = `if net.ParseIP(%[1]s) == nil {
	return fmt.Errorf("field '%%v' is not a valid IP address", %[2]q)
}`
	ipv4Tpl = `if ip := net.ParseIP(%[1]s); ip == nil || ip.To4() == nil {
	return fmt.Errorf("field '%%v' is not a valid IPv4", %[2]q)
}`
	ipv6Tpl = `if ip := net.ParseIP(%[1]s); ip == nil || ip.To4() != nil {
	return fmt.Errorf("field '%%v' is not a valid IPv6", %[2]q)
}`
	containsTpl = `if !strings.Contains(%[1]s, %[3]s) {
	return fmt.Errorf("field '%%v' does not contain '%%v'", %[2]q, %[3]s)
}`
	startsWithTpl = `if !strings.HasPrefix(%[1]s, %[3]s) {
	return fmt.Errorf("field '%%v' does not start with '%%v'", %[2]q, %[3]s)
}`
	endsWithTpl = `if !strings.HasSuffix(%[1]s, %[3]s) {
	return fmt.Errorf("field '%%v' does not end with '%%v'", %[2]q, %[3]s)
}`
	eqTpl = `if %[1]s != %[3]s {
	return fmt.Errorf("field '%%v' does not equal %%v", %[2]q, %[3]s)
}`
	neqTpl = `if %[1]s == %[3]s {
	return fmt.Errorf("field '%%v' should not be equal to %%v", %[2]q, %[3]s)
}`
	minTpl = `if %[1]s < %[3]s {
	return fmt.Errorf("field '%%v' has a minimum value of %%v", %[2]q, %[3]s)
}`
	maxTpl = `if %[1]s > %[3]s {
	return fmt.Errorf("field '%%v' has a maximum value of %%v", %[2]q, %[3]s)
}`
	maxlenTpl = `if sz > %[1]d {
		return fmt.Errorf("field '%%v' has a maximum length of %%v", %[2]q, %[1]d)
	}`
	minlenTpl = `if sz < %[1]d {
	return fmt.Errorf("field '%%v' has a minimum length of %%v", %[2]q, %[1]d)
}`
)

var constraintImports = map[string][]string{
	"alpha":      []string{"unicode"},
	"alphanum":   []string{"unicode"},
	"numeric":    []string{"regexp"},
	"email":      []string{"regexp"},
	"url":        []string{"net/url"},
	"base64":     []string{"regexp"},
	"ip":         []string{"net"},
	"ipv4":       []string{"net"},
	"ipv6":       []string{"net"},
	"contains":   []string{"strings"},
	"startswith": []string{"strings"},
	"endswith":   []string{"strings"},
}

var constraintDecls = map[string]string{
	"numeric":     `var numericConstraintRegex = regexp.MustCompile("^[-+]?[0-9]+(?:\\.[0-9]+)?$")`,
	"hexadecimal": `var hexadecimalConstraintRegex = regexp.MustCompile("^[0-9a-fA-F]+$")`,
	"email":       `var emailConstraintRegex = regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_` + "`" + `{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_` + "`" + `{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")`,
	"base64":      `var base64ConstraintRegex = regexp.MustCompile("^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$")`,
	"uuid":        `var uuidConstraintRegex = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")`,
}

func isValueOfType(v string, typ Type) bool {
	if m, ok := typ.(Maybe); ok {
		typ = m.Elem
	}

	switch typ.(Basic).Kind {
	case String:
		return true
	case Bool:
		return v == "true" || v == "false"
	case Int, Int8, Int16, Int32, Int64:
		_, err := strconv.ParseInt(v, 10, 64)
		return err == nil
	case Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		_, err := strconv.ParseUint(v, 10, 64)
		return err == nil
	case Float32, Float64:
		_, err := strconv.ParseFloat(v, 64)
		return err == nil
	default:
		return false
	}
}

func toPrintableValue(v string, typ Type) string {
	if isString(typ) {
		return fmt.Sprintf(`"%s"`, v)
	}
	return v
}

func constraintsForTpl(cs []Constraint, recv string) (before, after string) {
	csb, csa := splitConstraints(cs)
	return constraintsToCode(csb, recv), constraintsToCode(csa, recv)
}

func constraintsToCode(cs []Constraint, recv string) string {
	var buf bytes.Buffer
	for _, c := range cs {
		buf.WriteString(c.Validator(recv))
	}
	return buf.String()
}

func splitConstraints(
	cs []Constraint,
) (before []Constraint, after []Constraint) {
	for _, c := range cs {
		if c.BeforeRead() {
			before = append(before, c)
		} else {
			after = append(after, c)
		}
	}
	return before, after
}
