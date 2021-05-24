package plugin

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocRubyPlugin{})
}

// ProtocRubyPlugin implements Plugin for the built-in protoc C++ plugin.
type ProtocRubyPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocRubyPlugin) Name() string {
	return "protoc:ruby"
}

// Configure implements part of the Plugin interface.
func (p *ProtocRubyPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	cfg.Label = label.New("build_stack_rules_proto", "plugin/protoc", "ruby")
	cfg.Outputs = protoc.FlatMapFiles(
		protoc.RelativeFileNameWithExtensions(ctx.Rel, "_pb.rb"),
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)
}
