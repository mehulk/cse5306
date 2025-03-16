# Use an official Python runtime as a parent image.
FROM python:3.9-slim

# Set the working directory.
WORKDIR /app

# Install necessary packages.
RUN pip install grpcio grpcio-tools

# Copy the current directory contents into the container.
COPY . .

# Expose ports (adjust if needed; here we expose a range for all nodes).
EXPOSE 5001 5002 5003 5004 5005

# Default command (can be overridden by docker-compose).
CMD ["python", "-u", "raft-election.py", "-id=1", "-port=5001", "-leader", "-peers=node2:5002,node3:5003,node4:5004,node5:5005"]
