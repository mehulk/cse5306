import grpc
import average_pb2
import average_pb2_grpc
import sys

def run():
    with grpc.insecure_channel('localhost:50051') as channel: #grpc correction
        stub = average_pb2_grpc.CalculatorStub(channel)

        # Unary RPC: Average of two numbers
        if len(sys.argv) < 3:
            print("Usage: python calculator_client.py <num1> <num2> [numbers...]")
            return

        num1 = int(sys.argv[1])
        num2 = int(sys.argv[2])
        response = stub.Average(average_pb2.AverageRequest(num1=num1, num2=num2))  # pylint: disable=no-member
        print(f"Average: {response.result}")

        # Server Streaming RPC: Running average of a stream of numbers
        numbers = map(int, sys.argv[3:])
        number_iterator = (average_pb2.NumberRequest(number=num) for num in numbers)  # pylint: disable=no-member
        response = stub.RunningAverage(number_iterator)
        print(f"Running Average: {response.result}")

if __name__ == '__main__':
    run()
