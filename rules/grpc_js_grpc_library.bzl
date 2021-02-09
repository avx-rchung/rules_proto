load(
    "@build_stack_rules_proto//rules:grpc_js_grpc_compile.bzl",
    "grpc_js_grpc_compile",
)

load("@build_stack_rules_proto//rules:closure.bzl", "closure_proto_compile")

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

PROTO_DEPS = ["@io_bazel_rules_closure//closure/protobuf:jspb"]

def grpc_js_grpc_library(**kwargs):

    name = kwargs.get("name")
    name_pb = name + "_pb"
    name_pb_lib = name + "_pb_lib"
    name_pb_grpc = name + "_pb_grpc"

    closure_deps = kwargs.pop("closure_deps", [])
    deps = kwargs.get("deps", [])
    verbose = kwargs.pop("verbose", 0)
    visibility = kwargs.get("visibility", [])
    tags = kwargs.get("tags", [])

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
        tags = tags,
    )

    grpc_js_grpc_compile(
        name = name_pb_grpc,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
        tags = tags,
    )

    closure_js_library(
        name = name_pb_lib,
        srcs = [name_pb],
        deps = [
            "@io_bazel_rules_closure//closure/protobuf:jspb",
        ] + closure_deps,
        internal_descriptors = [
            name_pb + "/descriptor.source.bin",
        ],
        suppress = [
            "JSC_LATE_PROVIDE_ERROR",
            "JSC_UNDEFINED_VARIABLE",
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            "JSC_STRICT_INEXISTENT_PROPERTY",
            "JSC_POSSIBLE_INEXISTENT_PROPERTY",
            "JSC_UNRECOGNIZED_TYPE_ERROR",
            # "stricterMissingRequireType",
            # "analyzerChecks",
            # "analyzerChecksInternal",
        ],
        tags = tags,
        visibility = visibility,
    )

    closure_js_library(
        name = name,
        srcs = [name_pb_grpc],
        deps = [
            name_pb_lib,
            "@io_bazel_rules_closure//closure/library/promise",
            "@com_github_stackb_grpc_js//js/grpc/stream/observer:call",
            "@com_github_stackb_grpc_js//js/grpc",
            "@com_github_stackb_grpc_js//js/grpc:api",
            "@com_github_stackb_grpc_js//js/grpc:options",
        ] + closure_deps,
        internal_descriptors = [
            name_pb + "/descriptor.source.bin",
            name_pb_grpc + "/descriptor.source.bin",
        ],
        exports = [
            name_pb_lib,
        ],
        suppress = [
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            # "stricterMissingRequireType",
            # "analyzerChecks",
            # "analyzerChecksInternal",
        ],
        tags = tags,
        visibility = visibility,
    )