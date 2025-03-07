FROM python:3.9-slim

ENV PYTHONUNBUFFERED=1

WORKDIR /app

COPY proto/ /app/proto/
COPY src/ /app/src/
COPY requirements.txt /app/

RUN pip install -r requirements.txt

RUN python -m grpc_tools.protoc -I./proto --python_out=./src --grpc_python_out=./src ./proto/swim.proto

CMD ["python", "-u", "src/main.py"]

