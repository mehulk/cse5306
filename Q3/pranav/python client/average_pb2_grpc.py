# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

import average_pb2 as average__pb2

GRPC_GENERATED_VERSION = '1.70.0'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in average_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class CalculatorStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Average = channel.unary_unary(
                '/average.Calculator/Average',
                request_serializer=average__pb2.AverageRequest.SerializeToString,# pylint: disable=no-member
                response_deserializer=average__pb2.AverageResponse.FromString,# pylint: disable=no-member
                _registered_method=True)
        self.RunningAverage = channel.stream_unary(
                '/average.Calculator/RunningAverage',
                request_serializer=average__pb2.NumberRequest.SerializeToString,# pylint: disable=no-member
                response_deserializer=average__pb2.AverageResponse.FromString,# pylint: disable=no-member
                _registered_method=True)


class CalculatorServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Average(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def RunningAverage(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_CalculatorServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Average': grpc.unary_unary_rpc_method_handler(
                    servicer.Average,
                    request_deserializer=average__pb2.AverageRequest.FromString,# pylint: disable=no-member
                    response_serializer=average__pb2.AverageResponse.SerializeToString,# pylint: disable=no-member
            ),
            'RunningAverage': grpc.stream_unary_rpc_method_handler(
                    servicer.RunningAverage,
                    request_deserializer=average__pb2.NumberRequest.FromString,# pylint: disable=no-member
                    response_serializer=average__pb2.AverageResponse.SerializeToString,# pylint: disable=no-member
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'average.Calculator', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('average.Calculator', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class Calculator(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Average(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/average.Calculator/Average',
            average__pb2.AverageRequest.SerializeToString,# pylint: disable=no-member
            average__pb2.AverageResponse.FromString,# pylint: disable=no-member
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def RunningAverage(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_unary(
            request_iterator,
            target,
            '/average.Calculator/RunningAverage',
            average__pb2.NumberRequest.SerializeToString,# pylint: disable=no-member
            average__pb2.AverageResponse.FromString,# pylint: disable=no-member
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
