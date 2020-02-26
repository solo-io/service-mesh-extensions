package main

import (
	"context"
	"log"

	"github.com/solo-io/anyvendor/anyvendor"
	"github.com/solo-io/anyvendor/pkg/manager"
)

func main() {
	ctx := context.Background()
	mgr, err := manager.NewManager(ctx, ".")
	if err != nil {
		log.Fatal(err)
	}
	config := &anyvendor.Config{
		Local: &anyvendor.Local{
			Patterns: []string{"api/**/*.proto"},
		},
		Imports: []*anyvendor.Import{
			{
				ImportType: &anyvendor.Import_GoMod{
					GoMod: &anyvendor.GoModImport{
						Patterns: []string{"api/v1/**.proto"},
						Package:  "github.com/solo-io/solo-kit",
					},
				},
			},
			{
				ImportType: &anyvendor.Import_GoMod{
					GoMod: &anyvendor.GoModImport{
						Package:  "github.com/gogo/protobuf",
						Patterns: []string{"gogoproto/*.proto"},
					},
				},
			},
			{
				ImportType: &anyvendor.Import_GoMod{
					GoMod: &anyvendor.GoModImport{
						Package:  "github.com/solo-io/protoc-gen-ext",
						Patterns: []string{"extproto/*.proto"},
					},
				},
			},
		},
	}
	if err = mgr.Ensure(ctx, config); err != nil {
		log.Fatal(err)
	}

}
