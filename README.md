# Project Title: Distributed Systems & gRPC SWIM Implementation

## Overview

This repository contains our implementation and experiments for a multi-part assignment on Docker, gRPC, and the SWIM protocol. The project is organized into several sections corresponding to the questions in the assignment. The following README explains the assignment structure, outlines the development process for each question, and provides instructions on how to compile, run, and test the code.

Authors
Pranav S Damu - 1002158458
Mehul Kanotra - 1002166093

## Assignment Structure

- **Q1: Docker Get Started Guide**
  - **Objective:** Install Docker, complete foundational tutorials, and share Docker images on Docker Hub.
  - **Development:**  
    - Followed the [Docker Get Started guide](https://docs.docker.com/get-started/) and completed all required parts.
    - Categorized Docker commands with a brief explanation for each (e.g., `docker pull`, `docker run`).
    - Shared two personal Docker images: one for the introductory guide and one from Part 7 (Docker Compose).
  - **Deliverables:**  
    - Docker Hub image links (inserted in the final report).

- **Q2: gRPC Quick Start in Two Languages**
  - **Objective:** Get started with gRPC using two programming languages.
  - **Development:**  
    - Selected two languages from Python, Java, and Golang.  
    - Followed the respective Quick Start and Basics Tutorial guides (screenshots and command outputs captured).
    - For every step, recorded the command with a timestamp and username (using `date; whoami;` on Mac/Linux or `echo %date% %time%; whoami;` on Windows).
  - **Deliverables:**  
    - Screenshots for each step with brief explanations.
    - Both Mehul’s and Pranav's screenshots attached in the report.

- **Q3: Server-Client Implementation using gRPC**
  - **Objective:** Implement two basic server-client pairs using the same proto file but in different languages.
  - **Development:**  
    - **Mehul’s Contribution:** Created a Go client and a Python server.
    - **Pranav’s Contribution:** Created a Python client and a Go server.
    - Each implementation is contained in separate subfolders (`Q3/Mehul` and `Q3/Pranav`) with individual README files detailing how to run the code.
    - Containerized each component so that any instance of a server can communicate with any instance of a client.
  - **Deliverables:**  
    - Code for both server-client pairs.
    - Individual READMEs with execution instructions.

- **Q4: Failure Detector Component (SWIM Protocol)**
  - **Objective:** Implement the failure detector in Python.
  - **Development:**  
    - Designed a gRPC proto file defining the necessary data structures and RPC methods.
    - Implemented the failure detection logic:
      - Every T' seconds, each node pings a random node.
      - Upon a successful response, updates the "last heard" timestamp.
      - If a node fails to respond, additional nodes are contacted before marking it as failed.
    - Added logging statements on both the client and server sides in the specified format.
    - Containerized the implementation so that at least five nodes can interact.
  - **Deliverables:**  
    - Python code under the `Q4` folder.
    - A README with instructions on how to run and test the failure detector.

- **Q5: Dissemination Component (SWIM Protocol)**
  - **Objective:** Implement the dissemination component using a language different from Q4.
  - **Development:**  
    - Extended the same proto file to add new gRPC messages and service methods for dissemination.
    - Implemented functionality for:
      - Multicasting failed node information to the membership list.
      - Handling new node join requests via a bootstrap node.
    - Enabled inter-component communication via gRPC within the same container.
    - Containerized the implementation to allow at least five nodes to communicate.
  - **Deliverables:**  
    - Code in the `Q5` folder.
    - A README with detailed instructions and any unusual design choices explained.

- **Q6: Test Cases for SWIM Implementation**
  - **Objective:** Design and implement 5 different test cases.
  - **Development:**  
    - Created test cases to cover key scenarios (e.g., new node joining, node failure detection, and dissemination of failure information).
    - Documented each test case with a brief explanation and captured screenshots of their executions.
    - The test cases are provided in a file named `Q6-test-cases.txt`.
  - **Deliverables:**  
    - Test case documentation and screenshots included in the final report.
    - Instructions on how to execute the tests.

