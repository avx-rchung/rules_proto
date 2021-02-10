load(
    "@build_stack_rules_proto//rules:proto_aspect.bzl",
    "proto_compile_aspect",
    "proto_compile_rule_macro",
    "proto_compile_rule",
)

_default_plugins = [
    str(Label("//plugins/java/proto:proto")),
]

_java_proto_compile_aspect = proto_compile_aspect(_default_plugins, "java_proto_compile_aspect")

_java_proto_compile_rule = proto_compile_rule(_java_proto_compile_aspect)

def java_proto_compile(**kwargs):
    proto_compile_rule_macro(_java_proto_compile_rule, **kwargs)