package tengo

import (
	"fmt"

	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

type TengoObj struct {
	Name string
}

func (t *TengoObj) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {

	panic("implement me")
}

func (t *TengoObj) IsFalsy() bool {
	panic("implement me")
}

func (t *TengoObj) Equals(another objects.Object) bool {
	panic("implement me")
}

func (t *TengoObj) Copy() objects.Object {
	panic("implement me")
}

func (t *TengoObj) TypeName() string {

	return "tengo_object"
}

func (t *TengoObj) String() string {

	return fmt.Sprintf("TengoObject:%s", t.Name)
}
