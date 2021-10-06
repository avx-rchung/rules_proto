package rules_scala

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/plugin/scalapb/scalapb"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	GrpcscalaLibraryRuleName  = "grpc_scala_library"
	ProtoscalaLibraryRuleName = "proto_scala_library"
	scalaLibraryRuleSuffix    = "_scala_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:"+ProtoscalaLibraryRuleName,
		&scalaLibrary{
			kindName: ProtoscalaLibraryRuleName,
			shouldProvideRule: func(library protoc.ProtoLibrary, plugin *protoc.PluginConfiguration) bool {
				return !hasServicesAndGrpcOption(library, plugin)
			},
		})
	protoc.Rules().MustRegisterRule("stackb:rules_proto:"+GrpcscalaLibraryRuleName,
		&scalaLibrary{
			kindName:          GrpcscalaLibraryRuleName,
			shouldProvideRule: hasServicesAndGrpcOption,
		})
}

func hasServicesAndGrpcOption(library protoc.ProtoLibrary, plugin *protoc.PluginConfiguration) bool {
	// if any of the proto_library files have grpc service definitions AND the
	// grpc option is configured, emit a grpc_scala_library rule instead.
	if !protoc.HasServices(library.Files()...) {
		return false
	}
	for option, want := range plugin.Config.Options {
		if option == "grpc" && want {
			return true
		}
	}
	return false
}

// scalaLibrary implements LanguageRule for the '{proto|grpc}_scala_library' rule from
// @rules_proto.
type scalaLibrary struct {
	kindName          string
	shouldProvideRule func(library protoc.ProtoLibrary, plugin *protoc.PluginConfiguration) bool
}

// Name implements part of the LanguageRule interface.
func (s *scalaLibrary) Name() string {
	return s.kindName
}

// KindInfo implements part of the LanguageRule interface.
func (s *scalaLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs":       true,
			"deps":       true,
			"exports":    true,
			"visibility": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *scalaLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    fmt.Sprintf("@build_stack_rules_proto//rules/scala:%s.bzl", s.kindName),
		Symbols: []string{s.kindName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *scalaLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	plugin := pc.GetPluginConfiguration(scalapb.ScalaPBPluginName)
	if plugin == nil {
		log.Fatalf("expected plugin configuration for %q to be defined", scalapb.ScalaPBPluginName)
	}
	if len(plugin.Outputs) == 0 {
		return nil
	}
	if !s.shouldProvideRule(pc.Library, plugin) {
		return nil
	}

	return &scalaLibraryRule{
		kindName:       s.kindName,
		ruleNameSuffix: scalaLibraryRuleSuffix,
		outputs:        plugin.Outputs,
		ruleConfig:     cfg,
		config:         pc,
		resolver: func(impl protoc.DepsProvider, pc *protoc.ProtocConfiguration, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
			protoc.ResolveDepsWithSuffix(scalaLibraryRuleSuffix)(impl, pc, c, r, importsRaw, from)
			r.SetAttr("exports", r.Attr("deps"))
		},
	}
}

// scalaLibraryRule implements RuleProvider for 'scala_library'-derived rules.
type scalaLibraryRule struct {
	kindName       string
	ruleNameSuffix string
	outputs        []string
	config         *protoc.ProtocConfiguration
	ruleConfig     *protoc.LanguageRuleConfig
	resolver       protoc.DepsResolver
}

// Kind implements part of the ruleProvider interface.
func (s *scalaLibraryRule) Kind() string {
	return s.kindName
}

// Name implements part of the ruleProvider interface.
func (s *scalaLibraryRule) Name() string {
	return s.config.Library.BaseName() + s.ruleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *scalaLibraryRule) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.outputs {
		if strings.HasSuffix(output, ".srcjar") {
			srcs = append(srcs, protoc.StripRel(s.config.Rel, output))
		}
	}
	return srcs
}

// Deps computes the deps list for the rule.
func (s *scalaLibraryRule) Deps() []string {
	return s.ruleConfig.GetDeps()
}

// Visibility implements part of the ruleProvider interface.
func (s *scalaLibraryRule) Visibility() []string {
	visibility := make([]string, 0)
	for k, want := range s.ruleConfig.Visibility {
		if !want {
			continue
		}
		visibility = append(visibility, k)
	}
	sort.Strings(visibility)
	return visibility
}

// Rule implements part of the ruleProvider interface.
func (s *scalaLibraryRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *scalaLibraryRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
	if s.resolver == nil {
		return
	}
	s.resolver(s, s.config, c, r, importsRaw, from)
}