import grpc
import time
from concurrent import futures
import generated.average_pb2 as pb2
import generated.average_pb2_grpc as pb2_grpc

class CalculatorServicer(pb2_grpc.CalculatorServicer):
    def Average(self, request, context):
        result = (request.num1 + request.num2) / 2.0
        return pb2.AverageResponse(result=result)

    def RunningAverage(self, request_iterator, context):
        count, total = 0, 0
        for request in request_iterator:
            count += 1
            total += request.number
            print(f"Received: {request.number}, New Average: {total / count}")
        return pb2.AverageResponse(result=total / count)  # Send final response only once

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_CalculatorServicer_to_server(CalculatorServicer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Server running on port 50051...")
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()
