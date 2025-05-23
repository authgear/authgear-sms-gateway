package config_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	goyaml "gopkg.in/yaml.v2"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

func TestRootConfigValidation(t *testing.T) {
	Convey("RootConfig", t, func() {
		f, err := os.Open("testdata/authgear-sms-gateway.tests.yaml")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		type TestCase struct {
			Name   string      `json:"name"`
			Error  *string     `json:"error"`
			Config interface{} `json:"config"`
		}

		decoder := goyaml.NewDecoder(f)
		for {
			var testCase TestCase
			err := decoder.Decode(&testCase)
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				panic(err)
			}

			Convey(testCase.Name, func() {
				inputYAML, err := goyaml.Marshal(testCase.Config)
				if err != nil {
					panic(err)
				}

				ctx := context.Background()
				_, err = config.ParseRootConfigFromYAML(ctx, []byte(inputYAML))
				if testCase.Error != nil {
					So(err, ShouldBeError, *testCase.Error)
				} else {
					So(err, ShouldBeNil)
				}
			})
		}
	})
}
