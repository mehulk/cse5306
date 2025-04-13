FROM python:3.9-slim
WORKDIR /app
ENV PYTHONPATH=/app:/app/proto         

# ----- deps -----
COPY requirements.txt ./
RUN pip install --upgrade pip && pip install --no-cache-dir -r requirements.txt

# ----- source -----
COPY . .

# ----- (re)generate stubs -----
RUN rm -f proto/raft_pb2.py proto/raft_pb2_grpc.py && \
    python -m grpc_tools.protoc -I ./proto \
        --python_out=./proto --grpc_python_out=./proto \
        ./proto/raft.proto && \
    touch proto/__init__.py

# sitecustomize shim already copied with source

EXPOSE 5001 5002 5003 5004 5005
CMD ["sleep","infinity"]