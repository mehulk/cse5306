import grpc
import average_pb2
import average_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = average_pb2_grpc.CalculatorStub(channel)

        # Unary RPC: Average of two numbers
        print("Enter two numbers to calculate their average:")
        num1 = int(input("Enter first number: "))
        num2 = int(input("Enter second number: "))
        response = stub.Average(average_pb2.AverageRequest(num1=num1, num2=num2))# pylint: disable=no-member
        print(f"Average: {response.result}")

        # Server Streaming RPC: Running average of a stream of numbers
        print("\nEnter numbers one by one to calculate the running average (type 'done' to finish):")
        numbers = []
        while True:
            user_input = input("Enter a number (or 'done' to finish): ")
            if user_input.lower() == 'done':
                break
            try:
                number = int(user_input)
                numbers.append(number)
            except ValueError:
                print("Invalid input. Please enter a valid number.")

        # Send the stream of numbers to the server
        number_iterator = (average_pb2.NumberRequest(number=num) for num in numbers)# pylint: disable=no-member
        response = stub.RunningAverage(number_iterator)
        print(f"Running Average: {response.result}")

if __name__ == '__main__':
    run()
