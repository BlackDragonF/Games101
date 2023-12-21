package triangle

import (
	"assignment1/common"
	"testing"
)

func TestGet(t *testing.T) {
	triangle := NewTriangle()
	expected := common.NewVec3f()

	a := triangle.GetA()
	if !a.Equal(expected) {
		t.Errorf("test triangle: Wrong value of GetA(): Got: %v, Expected: %v", a, expected)
	}
	b := triangle.GetB()
	if !b.Equal(expected) {
		t.Errorf("test triangle: Wrong value of GetB(): Got: %v, Expected: %v", b, expected)
	}
	c := triangle.GetC()
	if !c.Equal(expected) {
		t.Errorf("test triangle: Wrong value of GetC(): Got: %v, Expected: %v", c, expected)
	}
}

func TestSetVertex(t *testing.T) {
	triangle := NewTriangle()
	if err := triangle.SetVertex(0, common.Vec3f{1., 2., 3.}); err != nil {
		t.Errorf("test triangle: %v", err)
	}
	if err := triangle.SetVertex(1, common.Vec3f{4., 5., 6.}); err != nil {
		t.Errorf("test triangle: %v", err)
	}
	if err := triangle.SetVertex(2, common.Vec3f{7., 8., 9.}); err != nil {
		t.Errorf("test triangle: %v", err)
	}

	if err := triangle.SetVertex(-1, common.Vec3f{0., 0., 0.}); err == nil {
		t.Errorf("test triangle: Wrong err value of SetVertex(): Got: nil")
	}
	if err := triangle.SetVertex(3, common.Vec3f{0., 0., 0.}); err == nil {
		t.Errorf("test triangle: Wrong err value of SetVertex(): Got: nil")
	}

	expected := common.Vec3f{1., 2., 3.}
	a := triangle.GetA()
	if !a.Equal(expected) {
		t.Errorf("test triangle: Wrong value of GetA(): Got: %v, Expected: %v", a, expected)
	}
	expected = common.Vec3f{4., 5., 6.}
	b := triangle.GetB()
	if !b.Equal(expected) {
		t.Errorf("test triangle: Wrong value of GetB(): Got: %v, Expected: %v", b, expected)
	}
	expected = common.Vec3f{7., 8., 9.}
	c := triangle.GetC()
	if !c.Equal(expected) {
		t.Errorf("test triangle: Wrong value of GetC(): Got: %v, Expected: %v", c, expected)
	}
}

func TestSetNormal(t *testing.T) {
	//TODO
}

func TestSetColor(t *testing.T) {
	//TODO
}

func TestSetTexCoord(t *testing.T) {
	//TODO
}

func TestToVector4(t *testing.T) {
	triangle := NewTriangle()
	expected := common.NewVec4f()
	expected[3] = 1.

	vec4s := triangle.ToVec4()
	for i := 0; i < len(vec4s); i++ {
		if !vec4s[i].Equal(expected) {
			t.Errorf("test triangle: Wrong value of GetA(): Got: %v, Expected: %v", vec4s[i], expected)
		}
	}
}
