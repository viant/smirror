package mirror

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/storage"
	"io"
	"io/ioutil"
	"smirror/job"
	"strings"
	"testing"
)

type serviceUseCase struct {
	description    string
	sourceURL      string
	sourceContent  string
	config         *Config
	compress       bool
	expectResponse interface{}
	expectedURLs   map[string]int
}

func TestService_Mirror(t *testing.T) {

	var useCases = []*serviceUseCase{
		{
			description: "no split transfer, folder depth = 0",
			sourceURL:   "mem://localhost/folder/subfolder/file1.txt",
			sourceContent: `line1,
line2,
line3,
line4`,
			config: &Config{
				Routes: Routes{
					&Route{
						Suffix:  ".txt",
						DestURL: "mem://localhost/cloned/data",
					},
				},
			},
			expectedURLs: map[string]int{
				"mem://localhost/cloned/data/file1.txt":      26,
				"mem://localhost/folder/subfolder/file1.txt": 26,
			},
			expectResponse: `{
	"DestURLs": ["mem://localhost/cloned/data/file1.txt"],
	"Status": "ok"
}`,
		},
		{
			description: "no split transfer, folder depth = 2",
			sourceURL:   "mem://localhost/folder/subfolder/file1.txt",
			sourceContent: `line1,
line2,
line3,
line4`,
			config: &Config{
				Routes: Routes{
					&Route{
						FolderDepth: 2,
						Suffix:      ".txt",
						DestURL:     "mem://localhost/cloned/data",
					},
				},
			},
			expectedURLs: map[string]int{
				"mem://localhost/cloned/data/folder/subfolder/file1.txt": 26,
				"mem://localhost/folder/subfolder/file1.txt":             26,
			},
			expectResponse: `{
	"DestURLs": ["mem://localhost/cloned/data/folder/subfolder/file1.txt"],
	"Status": "ok"
}`,
		},

		{
			description: "split transfer",
			sourceURL:   "mem://localhost/folder/subfolder/file1.txt",
			sourceContent: `line1,
line2,
line3,
line4`,
			config: &Config{
				Routes: Routes{
					&Route{
						Suffix:  ".txt",
						DestURL: "mem://localhost/cloned/data",
						Split: &Split{
							MaxLines: 3,
							Template: "%v_%05d",
						},
					},
				},
			},
			expectedURLs: map[string]int{
				"mem://localhost/cloned/data/file1.txt":       26,
				"mem://localhost/cloned/data/file1_00001.txt": 20,
				"mem://localhost/cloned/data/file1_00002.txt": 5,
			},
			expectResponse: `{
	"DestURLs": [
		"mem://localhost/cloned/data/file1_00001.txt",
		"mem://localhost/cloned/data/file1_00002.txt"
	],
	"Status": "ok"
}`,
		},


		{
			description: "compressed transfer",
			sourceURL:   "mem://localhost/folder/subfolder/file1.txt",
			sourceContent: `line1,
line2,
line3,
line4,
line5,
line6,
line7,
line8`,
			config: &Config{
				Routes: Routes{
					&Route{
						Suffix:  ".txt",
						DestURL: "mem://localhost/data",
						Compression: &Compression{
							Codec: GZipCodec,
						},
					},
				},
			},
			expectedURLs: map[string]int{
				"mem://localhost/folder/subfolder/file1.txt": 54,
				"mem://localhost/data/file1.txt.gz":          41,
			},
			expectResponse: `{
	"DestURLs": [
		"mem://localhost/data/file1.txt.gz"
	],
	"Status": "ok"
}`,
		},

		{
			description: "compressed split transfer",
			sourceURL:   "mem://localhost/folder/subfolder/file1.txt",
			sourceContent: `line1,
line2,
line3,
line4,
line5,
line6,
line7,
line8
line9,
line10,
line11
`,
			config: &Config{
				Routes: Routes{
					&Route{
						Suffix:  ".txt",
						DestURL: "mem://localhost/data",
						Split: &Split{
							MaxLines: 10,
							Template: "%v_%05d",
						},
						Compression: &Compression{
							Codec: GZipCodec,
						},
					},
				},
			},
			expectedURLs: map[string]int{
				"mem://localhost/folder/subfolder/file1.txt": 77,
				"mem://localhost/data/file1_00001.txt.gz":    62,
				"mem://localhost/data/file1_00002.txt.gz":    35,
			},
			expectResponse: `{
	"DestURLs": [
		"mem://localhost/data/file1_00001.txt.gz",
		"mem://localhost/data/file1_00002.txt.gz"
	],
	"Status": "ok"
}`,
		},

		{
			description: "compressed split transfer",
			sourceURL:   "mem://localhost/folder/subfolder/file1.txt.gz",
			sourceContent: `line1,
line2,
line3,
line4,
line5,
line6,
line7,
line8
line9,
line10,
line11
`,
			config: &Config{
				Routes: Routes{
					&Route{
						Suffix:  ".txt.gz",
						DestURL: "mem://localhost/data",
						Split: &Split{
							MaxLines: 10,
							Template: "%v_%05d",
						},
						Compression: &Compression{
							Codec: GZipCodec,
						},
					},
				},
			},
			expectedURLs: map[string]int{
				"mem://localhost/folder/subfolder/file1.txt.gz": 64,
				"mem://localhost/data/file1_00001.txt.gz":       62,
				"mem://localhost/data/file1_00002.txt.gz":       35,
			},
			expectResponse: `{
	"DestURLs": [
		"mem://localhost/data/file1_00001.txt.gz",
		"mem://localhost/data/file1_00002.txt.gz"
	],
	"Status": "ok"
}`,
		},

		{
			description: "no split transfer - on success delete ",
			sourceURL:   "mem://localhost/folder/subfolder/file1.txt",
			sourceContent: `line1,
line2,
line3,
line4`,
			config: &Config{
				Routes: Routes{
					&Route{
						Suffix:  ".txt",
						DestURL: "mem://localhost/cloned/data",
						OnCompletion:job.Completion{
							OnSuccess:[]*job.Action{
								{
									Name:job.ActionDelete,
								},
							},
						},
					},
				},
			},
			expectedURLs: map[string]int{
				"mem://localhost/cloned/data/file1.txt":      26,
				"mem://localhost/folder/subfolder/file1.txt": 0,
			},
			expectResponse: `{
	"DestURLs": ["mem://localhost/cloned/data/file1.txt"],
	"Status": "ok"
}`,
		},
		{
			description: "no split transfer - on success move ",
			sourceURL:   "mem://localhost/folder/subfolder/file1.txt",
			sourceContent: `line1,
line2,
line3,
line4`,
			config: &Config{
				Routes: Routes{
					&Route{
						Suffix:  ".txt",
						DestURL: "mem://localhost/cloned/data",
						OnCompletion:job.Completion{
							OnSuccess:[]*job.Action{
								{
									Name: job.ActionMove,
									URL:  "mem://localhost/processed",
								},
							},
						},
					},
				},
			},
			expectedURLs: map[string]int{
				"mem://localhost/cloned/data/file1.txt":      26,
				"mem://localhost/folder/subfolder/file1.txt": 0,
				"mem://localhost/processed/file1.txt":26,
			},
			expectResponse: `{
	"DestURLs": ["mem://localhost/cloned/data/file1.txt"],
	"Status": "ok"
}`,
		},
	}

	memStorage := storage.NewMemoryService()

	for _, useCase := range useCases {

		initUseCase(useCase, memStorage, t)

		service := New(useCase.config)
		response := service.Mirror(&Request{URL: useCase.sourceURL})
		if !assertly.AssertValues(t, useCase.expectResponse, response, useCase.description) {
			_ = toolbox.DumpIndent(response, true)
		}
		if len(useCase.expectedURLs) == 0 {
			continue
		}

		for URL, expectedSize := range useCase.expectedURLs {
			reader, err := memStorage.DownloadWithURL(URL)
			if expectedSize == 0 { //DO NOT EXPECT ASSET IN THAT URL
				if assert.NotNil(t, err, useCase.description) {
					continue
				}
			}
			if ! assert.Nil(t, err, useCase.description + " on " + URL) {
				continue
			}

			data, err := ioutil.ReadAll(reader)
			assert.Nil(t, err, useCase.description)
			assert.Equal(t, expectedSize, len(data), useCase.description)

		}

	}

}

func initUseCase(useCase *serviceUseCase, memStorage storage.Service, t *testing.T) {
	var sourceReader io.Reader = strings.NewReader(useCase.sourceContent)
	if strings.HasSuffix(useCase.sourceURL, GZIPExtension) {
		buffer := new(bytes.Buffer)
		writer := gzip.NewWriter(buffer)
		_, _ = io.Copy(writer, sourceReader)
		_ = writer.Flush()
		_ = writer.Close()
		sourceReader = buffer
	}
	err := memStorage.Upload(useCase.sourceURL, sourceReader)
	assert.Nil(t, err, useCase.description)
}
