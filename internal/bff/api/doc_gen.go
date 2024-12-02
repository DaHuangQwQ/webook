package api

type OpenAPISpec struct {
	OpenAPI string          `json:"openapi"`
	Info    Info            `json:"info"`
	Paths   map[string]Path `json:"paths"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Path struct {
	Get *Operation `json:"get,omitempty"`
}

type Operation struct {
	Summary     string              `json:"summary"`
	Parameters  []Parameter         `json:"parameters,omitempty"` // 请求参数
	Responses   map[string]Response `json:"responses"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
}

type Parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`       // query, path, header, cookie
	Required bool   `json:"required"` // 是否为必填项
	Schema   Schema `json:"schema"`   // 参数的数据类型
}

type Schema struct {
	Type       string            `json:"type"` // 数据类型，如 string, integer
	Properties map[string]Schema `json:"properties,omitempty"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"` // 响应内容
}

type MediaType struct {
	Schema Schema `json:"schema"` // 响应内容的数据结构
}

type RequestBody struct {
	Content map[string]MediaType `json:"content"`
}
