package refmt_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/obj/atlas"
)

func ExampleJsonEncodeAtlasDefaults() {
	type MyType struct {
		X string
		Y int
	}

	MyType_AtlasEntry := atlas.BuildEntry(MyType{}).
		StructMap().Autogenerate().
		Complete()

	atl := atlas.MustBuild(
		MyType_AtlasEntry,
		// this is a vararg... stack more entries here!
	)

	var buf bytes.Buffer
	encoder := refmt.NewAtlasedJsonEncoder(&buf, atl)
	err := encoder.Marshal(MyType{"a", 1})
	fmt.Println(buf.String())
	fmt.Printf("%v\n", err)

	// Output:
	// {"x":"a","y":1}
	// <nil>
}

func ExampleJsonEncodeAtlasCustom() {
	type MyType struct {
		X string
		Y int
	}

	MyType_AtlasEntry := atlas.BuildEntry(MyType{}).
		StructMap().
		AddField("X", atlas.StructMapEntry{SerialName: "overrideName"}).
		// and no "Y" mapping at all!
		Complete()

	atl := atlas.MustBuild(
		MyType_AtlasEntry,
		// this is a vararg... stack more entries here!
	)

	var buf bytes.Buffer
	encoder := refmt.NewAtlasedJsonEncoder(&buf, atl)
	err := encoder.Marshal(MyType{"a", 1})
	fmt.Println(buf.String())
	fmt.Printf("%v\n", err)

	// Output:
	// {"overrideName":"a"}
	// <nil>
}

func ExampleJsonEncodeAtlas() {
	type MyType struct {
		X string
		Y string
		Z string
	}

	MyType_AtlasEntry := atlas.BuildEntry(MyType{}).
		Transform().
		TransformMarshal(atlas.MakeMarshalTransformFunc(
			func(x MyType) (string, error) {
				return string(x.X) + ":" + string(x.Y) + ":" + string(x.Z), nil
			})).
		TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
			func(x string) (MyType, error) {
				ss := strings.Split(x, ":")
				if len(ss) != 3 {
					return MyType{}, fmt.Errorf("parsing MyType: string must have 3 parts, separated by colon")
				}
				return MyType{ss[0], ss[1], ss[2]}, nil
			})).
		Complete()

	atl := atlas.MustBuild(
		MyType_AtlasEntry,
		// this is a vararg... stack more entries here!
	)

	var buf bytes.Buffer
	encoder := refmt.NewAtlasedJsonEncoder(&buf, atl)

	err := encoder.Marshal(MyType{"serializes", "as", "string!"})
	fmt.Println(buf.String())
	fmt.Printf("%v\n", err)

	// Output:
	// "serializes:as:string!"
	// <nil>
}
