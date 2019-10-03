package smirror

import (
	"context"
	"fmt"
	"github.com/viant/afs/matcher"
	"github.com/viant/toolbox"
	"log"
	"os"
	"path"
	"smirror/config"
	"testing"
)

func Test_Service(t *testing.T) {



	config:= &Config{
	Mirrors: config.Routes{
		Rules: []*config.Route{
			{
				Source: &config.Resource{
					Basic: matcher.Basic{
						Suffix: ".json",
					},
				},
				Dest: &config.Resource{
					URL: "mem://localhost/data",
				},
				Compression: &config.Compression{
					Codec: config.GZipCodec,
				},
				Replace:[]*config.Replace{
				    	{From:`""`, To:`"`},
				},
			},
		},
	}}

	ctx := context.Background()
	service, err := New(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	homeDir := os.Getenv("HOME")

	sourceURL:= fmt.Sprintf("file://%v", path.Join(homeDir, "Downloads", "feed.json"))

	response := service.Mirror(ctx, &Request{URL: sourceURL})
	//feed.json
	toolbox.DumpIndent(response, true)



}
