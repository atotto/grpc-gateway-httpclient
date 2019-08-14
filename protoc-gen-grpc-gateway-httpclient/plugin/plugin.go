package plugin

import (
	"fmt"
	"log"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	options "github.com/istio/gogo-genproto/googleapis/google/api"
)

type protoHttp struct {
	*generator.Generator
	generator.PluginImports
	localName string
}

func NewPlugin() *protoHttp {
	return &protoHttp{}
}

func (p *protoHttp) Name() string {
	return "grpc-gateway-httpclient"
}

func (p *protoHttp) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *protoHttp) typeName(s string) string {
	return p.Generator.TypeName(p.Generator.ObjectNamed(s))
}

func extractAPIOptions(meth *descriptor.MethodDescriptorProto) (*options.HttpRule, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, options.E_Http) {
		return nil, nil
	}
	ext, err := proto.GetExtension(meth.Options, options.E_Http)
	if err != nil {
		return nil, err
	}
	opts, ok := ext.(*options.HttpRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}
	return opts, nil
}

func (p *protoHttp) Generate(file *generator.FileDescriptor) {
	if !gogoproto.IsProto3(file.FileDescriptorProto) {
		// not supported
		return
	}
	p.PluginImports = NewPluginImports(p.Generator)
	p.NewImport("context")
	p.NewImport("net/http")
	p.NewImport("net/url")
	p.NewImport("google.golang.org/grpc")
	p.NewImport("github.com/atotto/grpc-gateway-httpclient/runtime")

	p.localName = generator.FileName(file)

	for _, service := range file.Service {
		serviceName := service.GetName()
		serviceClientName := fmt.Sprintf("%sHttpClient", serviceName)

		p.P(``)
		p.P(`func New`, serviceClientName, `(client *http.Client, url *url.URL) *`, serviceClientName, ` {`)
		p.P(`    return &`, serviceClientName, `{Client: &runtime.Client{Client: client}, URL: *url}`)
		p.P(`}`)
		p.P(``)

		p.P(`type `, serviceClientName, ` struct {`)
		p.P(`    Client *runtime.Client`)
		p.P(`    URL url.URL`)
		p.P(`}`)

		for _, method := range service.Method {
			rule, err := extractAPIOptions(method)
			if err != nil {
				log.Printf("%s invalid http rule: %s", method.GetName(), err)
				continue
			}
			p.P(``)
			p.P(`func (c *`, serviceClientName, `)`, fmt.Sprintf("%+v", method.GetName()), `(ctx context.Context, req *`, p.typeName(method.GetInputType()), `, opts ...grpc.CallOption) (*`, p.typeName(method.GetOutputType()), `, error) {`)
			p.In()
			p.P(`var err error`)
			p.P(`res := &`, p.typeName(method.GetOutputType()), `{}`)
			if rule != nil {
				p.P(`u := c.URL`)
				// TODO: support query
				switch httpRule := rule.Pattern.(type) {
				case *options.HttpRule_Get:
					p.P(`u.Path = `, MustParsePattern(httpRule.Get))
					p.P(`err = c.Client.Get(ctx, u.String(), req, res)`)
				case *options.HttpRule_Post:
					p.P(`u.Path = `, MustParsePattern(httpRule.Post))
					p.P(`err = c.Client.Post(ctx, u.String(), req, res)`)
				case *options.HttpRule_Delete:
					p.P(`u.Path = `, MustParsePattern(httpRule.Delete))
					p.P(`err = c.Client.Delete(ctx, u.String(), req, res)`)
				case *options.HttpRule_Patch:
					p.P(`u.Path = `, MustParsePattern(httpRule.Patch))
					p.P(`err = c.Client.Patch(ctx, u.String(), req, res)`)
				case *options.HttpRule_Custom:
					// TODO: implement
				}
				if err != nil {
					log.Printf("%s invalid http rule: %s", method.GetName(), err)
				}
			}
			p.P(`return res, err`)
			p.Out()
			p.P(`}`)
		}
	}
}

func init() {
	generator.RegisterPlugin(NewPlugin())
}
