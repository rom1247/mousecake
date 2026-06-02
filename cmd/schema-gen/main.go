// cmd/schema-gen 从 config.Config 生成 JSON Schema（Draft 2020-12），
// 输出到 config/app.schema.json，供 IDE yaml-language-server 使用。
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/mousecake-go/mousecake-go/config"
)

// durationType 用于 reflect.TypeOf 比较。
var durationType = reflect.TypeOf(time.Duration(0))

func main() {
	r := &jsonschema.Reflector{
		// Mapper 将 time.Duration 映射为 string + Go duration pattern
		Mapper: func(t reflect.Type) *jsonschema.Schema {
			if t == durationType {
				return &jsonschema.Schema{
					Type:        "string",
					Pattern:     "^[0-9]+(ns|us|ms|s|m|h)$",
					Description: "Go duration 格式，如 30m、200ms、10s",
				}
			}
			return nil
		},
	}

	r.AddGoComments("github.com/mousecake-go/mousecake-go", "./config")

	s := r.Reflect(&config.Config{})
	s.Title = "MouseCake Go 配置"
	s.Description = "mousecake-go 应用的 YAML 配置文件 Schema"
	s.Version = "https://json-schema.org/draft/2020-12/schema"

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "序列化 schema 失败: %v\n", err)
		os.Exit(1)
	}

	outPath := "app.schema.json"
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "写入 %s 失败: %v\n", outPath, err)
		os.Exit(1)
	}

	fmt.Printf("已生成 %s\n", outPath)
}
