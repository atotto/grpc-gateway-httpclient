package plugin

import (
	"path"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type pluginImports struct {
	generator *generator.Generator
	singles   []generator.Single
}

func NewPluginImports(g *generator.Generator) *pluginImports {
	return &pluginImports{g, make([]generator.Single, 0)}
}

func (p *pluginImports) NewImport(pkg string) generator.Single {
	imp := newImportedPackage(pkg)
	p.singles = append(p.singles, imp)
	return imp
}

func (p *pluginImports) GenerateImports(file *generator.FileDescriptor) {
	for _, s := range p.singles {
		if s.IsUsed() {
			p.generator.PrintImport(generator.GoPackageName(s.Name()), generator.GoImportPath(s.Location()))
		}
	}
}

type importedPackage struct {
	pkg  string
	name string
}

func newImportedPackage(pkg string) *importedPackage {
	_, name := path.Split(pkg)
	return &importedPackage{
		pkg:  pkg,
		name: name,
	}
}

func (im *importedPackage) Use() string      { return im.name }
func (im *importedPackage) IsUsed() bool     { return true }
func (im *importedPackage) Name() string     { return im.name }
func (im *importedPackage) Location() string { return im.pkg }
