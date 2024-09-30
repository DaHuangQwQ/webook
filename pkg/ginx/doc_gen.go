package ginx

import (
	"reflect"
	"webook/bff/api"
)

var Paths = make(map[string]api.Path)

//func DocGen(req any) {
//	var (
//		path       string
//		title      string
//		oper       *api.Operation
//		parameters []api.Parameter
//	)
//	t := reflect.TypeOf(req)
//	for i := 0; i < t.NumField(); i++ {
//		field := t.Field(i)
//		if field.Name == "Meta" {
//			path = field.Tag.Get("path")
//			title = field.Tag.Get("title")
//		} else {
//			parameters = append(parameters, api.Parameter{
//				Name:     field.Tag.Get("json"),
//				In:       getParameterLocation(field),
//				Required: field.Tag.Get("required") == "true",
//				Schema: api.Schema{
//					Type: getJSONType(field.Type),
//				},
//			})
//		}
//
//		// 打印字段名和对应的 JSON 标签
//		//fmt.Printf("Field Name: %s， Type : %s, JSON Tag: %s\n", field.Name, field.Type, field.Tag.Get("json"))
//	}
//	oper = &api.Operation{
//		Summary:    title,
//		Parameters: parameters,
//		Responses: map[string]api.Response{
//			"200": {
//				Description: "A list of users",
//				Content: map[string]api.MediaType{
//					"application/json": {
//						Schema: api.Schema{
//							Type: "array", // 返回的是一个用户数组
//						},
//					},
//				},
//			},
//		},
//	}
//	Paths[path] = api.Path{
//		Get: oper,
//	}
//}

func DocGen(req any) {
	var (
		path       string
		title      string
		operation  *api.Operation
		parameters []api.Parameter
		bodySchema api.Schema // 用于存储请求体 schema
		hasBody    bool       // 用于判断是否存在请求体
	)
	t := reflect.TypeOf(req)

	// 遍历结构体的字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Meta" {
			// 提取路径和标题
			path = field.Tag.Get("path")
			title = field.Tag.Get("title")
		} else {
			// 判断字段的参数位置
			location := getParameterLocation(field)
			if location == "body" {
				// 如果是 body 参数，构建请求体的 schema
				hasBody = true
				if bodySchema.Type == "" {
					bodySchema = api.Schema{
						Type: "object", // 请求体是一个对象
						Properties: map[string]api.Schema{
							field.Tag.Get("json"): {
								Type: getJSONType(field.Type),
							},
						},
					}
				} else {
					// 添加更多字段到请求体
					bodySchema.Properties[field.Tag.Get("json")] = api.Schema{
						Type: getJSONType(field.Type),
					}
				}
			} else {
				// 收集非 body 的请求参数（如 query, path 等）
				parameters = append(parameters, api.Parameter{
					Name:     field.Tag.Get("json"),
					In:       location,
					Required: field.Tag.Get("required") == "true",
					Schema: api.Schema{
						Type: getJSONType(field.Type),
					},
				})
			}
		}
	}

	// 构建操作对象
	operation = &api.Operation{
		Summary:    title,
		Parameters: parameters,
		Responses: map[string]api.Response{
			"200": {
				Description: "Successful Response",
				Content: map[string]api.MediaType{
					"application/json": {
						Schema: api.Schema{
							Type: "array", // 假设返回一个数组
						},
					},
				},
			},
		},
	}

	// 如果有请求体，添加 requestBody
	if hasBody {
		operation.RequestBody = &api.RequestBody{
			Content: map[string]api.MediaType{
				"application/json": {
					Schema: bodySchema,
				},
			},
		}
	}

	// 存储路径对应的操作
	Paths[path] = api.Path{
		Get: operation, // 假设仅支持 GET 请求
	}
}

func getJSONType(v reflect.Type) string {
	switch v.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Struct:
		return "object"
	case reflect.Map:
		if v.Key().Kind() == reflect.String {
			return "object" // 在 JSON 中，map 的 key 只能是字符串
		}
	}
	return "unknown"
}

func getParameterLocation(field reflect.StructField) string {
	// 根据 Tag 中 "method" 的值来判断
	method := field.Tag.Get("method")
	if method == "GET" {
		return "query"
	}
	return "body"
}
