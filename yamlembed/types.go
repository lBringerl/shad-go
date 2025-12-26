package yamlembed

import (
	"strings"
)

type Foo struct {
	A string `yaml:"aa"`
	p int64  `yaml:"-"`
}

type Bar struct {
	I      int64    `yaml:"-"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"oi,omitempty"`
	F      []any    `yaml:"f,flow"`
}

type barUnmarshall struct {
	I      int64    `yaml:"-"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"oi,omitempty"`
	F      []any    `yaml:"f"`
}

func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var inner barUnmarshall
	unmarshal(&inner)

	b.I = inner.I
	b.B = inner.B
	b.UpperB = strings.ToUpper(inner.B)
	b.OI = inner.OI
	b.F = inner.F

	return nil
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}

func (b *Baz) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&b.Foo)
	if err != nil {
		return err
	}

	err = b.Bar.UnmarshalYAML(unmarshal)
	if err != nil {
		return err
	}

	return nil
}
