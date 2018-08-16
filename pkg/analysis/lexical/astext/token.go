package astext

import (
	"bytes"
	"fmt"

	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
)

// tokenName returns a name for a token.
// nolint: gocyclo
func TokenName(token interface{}) string {
	switch t := token.(type) {
	case *ast.Apply:
		return "(apply)"
	case *ast.Array:
		return "(array)"
	case *ast.ArrayComp:
		return "(arraycomp)"
	case *ast.Binary:
		return "(binary)"
	case *ast.Conditional:
		return "(conditional)"
	case *ast.DesugaredObject:
		return "(object)"
	case ast.DesugaredObjectField:
		name := TokenValue(t.Name)
		return fmt.Sprintf("(field) %s", name)
	case ast.ForSpec:
		return "forspec"
	case *ast.Function:
		return fmt.Sprintf("(function)")
	case *ast.LiteralBoolean:
		return "(bool)"
	case *ast.LiteralNull:
		return "(null)"
	case *ast.LiteralNumber:
		return "(number)"
	case *ast.LiteralString:
		return "(string)"
	case ast.Identifier:
		return fmt.Sprintf("(identifier) %s", string(t))
	case *ast.Identifier:
		return fmt.Sprintf("(identifier) %s", string(*t))
	case *ast.Import:
		return fmt.Sprintf("(import) %s", t.File.Value)
	case *ast.ImportStr:
		return fmt.Sprintf("(importstr) %s", t.File.Value)
	case *ast.Index:
		if t.Id == nil {
			return fmt.Sprintf("(array index) [%s]", TokenValue(t.Index))
		}
		return fmt.Sprintf("(index) %s", string(*t.Id))
	case *ast.Local:
		return "(local)"
	case ast.LocalBind:
		return fmt.Sprintf("(local bind) %s", string(t.Variable))
	case ast.NamedParameter:
		val := TokenValue(t.DefaultArg)
		return fmt.Sprintf("(optional parameter) %s=%s", string(t.Name), val)
	case *ast.Object:
		return "(object)"
	case ast.ObjectField:
		return fmt.Sprintf("(field) %s", ObjectFieldName(t))
	case *ast.Self:
		return "(self)"
	case *ast.SuperIndex:
		return fmt.Sprintf("(super index) %s", string(*t.Id))
	case *ast.Var:
		return fmt.Sprintf("(var) %s", string(t.Id))
	case RequiredParameter:
		return fmt.Sprintf("(required parameter) %s", string(t.ID))
	default:
		return fmt.Sprintf("(unknown) %T", t)
	}
}

func TokenValue(token interface{}) string {
	switch t := token.(type) {
	case *ast.LiteralNumber:
		return t.OriginalString
	case *ast.LiteralString:
		return t.Value
	default:
		return fmt.Sprintf("unknown value from %T", t)
	}
}

func ObjectFieldName(f ast.ObjectField) string {
	if f.Id != nil {
		return string(*f.Id)
	}

	if f.Expr1 != nil {
		return TokenValue(f.Expr1)
	}

	panic("object field does not have a name")
}

func ObjectFieldVisibility(f ast.ObjectField) (string, error) {
	switch f.Hide {
	case ast.ObjectFieldHidden:
		return "::", nil
	case ast.ObjectFieldInherit:
		return ":", nil
	case ast.ObjectFieldVisible:
		return ":::", nil
	default:
		return "", errors.Errorf("unknown object visibility %d", f.Hide)
	}
}

func ObjectDescription(o *ast.Object) (string, error) {
	if o == nil {
		return "", errors.New("object is nil")
	}

	var buf bytes.Buffer
	if _, err := buf.WriteString("(object) {"); err != nil {
		return "", err
	}

	// find object fields
	for i, field := range o.Fields {
		if i == 0 {
			if _, err := buf.WriteString("\n"); err != nil {
				return "", err
			}
		}
		fieldName := ObjectFieldName(field)
		visibility, err := ObjectFieldVisibility(field)
		if err != nil {
			return "", err
		}
		if _, err := buf.WriteString(fmt.Sprintf("  (field) %s%s,\n", fieldName, visibility)); err != nil {
			return "", err
		}
	}
	if _, err := buf.WriteString("}"); err != nil {
		return "", err
	}

	return buf.String(), nil
}
