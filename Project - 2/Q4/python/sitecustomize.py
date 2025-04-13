"""Polyfill for generated stubs that expect protobuf‑5.x runtime."""
import types, google.protobuf as _pb

if not hasattr(_pb, "runtime_version"):
    _shim = types.ModuleType("runtime_version")

    class _Domain:
        PUBLIC = 4
        INTERNAL = 3
    _shim.Domain = _Domain

    def _validate(*_, **__):
        return None  # no‑op
    _shim.ValidateProtobufRuntimeVersion = _validate

    _pb.runtime_version = _shim