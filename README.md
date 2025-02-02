# Langcompiler

A secure and scalable service that executes code snippets in multiple programming languages, Current supporting -  Python, Java, JavaScript, and C++.

## Table of Contents
- [Features](#features)
- [API Documentation](#api-documentation)
- [Security](#security)

## Features
- Supports multiple programming languages:
  - Python
  - Java
  - JavaScript (Node.js)
  - C++
- Code sanitization and security checks
- Configurable execution timeouts
- Resource usage limits
- RESTful API interface
- Docker containerization

## API Documentation

### Execute Code
Executes code in the specified programming language.

**Endpoint:** `POST /execute`

**Request Body:**

```json
{
"language": "string", // Required: "python", "java", "javascript", or "cpp"
"code": "string" // Required: Source code to execute
}
```

**Response:**

```json
{
"output": "string", // Execution output
"error": "string" // Error message (if any)
}
```

**Example Requests:**


```bash

curl -X POST http://localhost:8080/execute \
-H "Content-Type: application/json" \
-d '{
    "language": "java",
    "code": "public class Main {\n    public static void main(String[] args) {\n        System.out.println(\"Hello, World!\");\n        for(int i = 0; i < 3; i++) {\n            System.out.println(\"Count: \" + i);\n        }\n    }\n}"
}

> {
    "output": "Hello, World!\nCount: 0\nCount: 1\nCount: 2\n"
}
```

## Security

### Code Sanitization
The service implements several security measures:

1. Pattern-based code analysis
2. Resource limitations
3. Isolated execution environment
4. Non-root user execution
5. Timeout constraints

### Unsafe Patterns
The service blocks potentially dangerous operations:
- File system operations
- Network access
- System command execution
- Resource-intensive operations