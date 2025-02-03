**LangCompiler** is a cloud-native, high-performance service carefully designed for secure, scalable execution of code across a wide variety of programming languages. It now supports strong features in four of the most widely used programming languages: Python, Java, JavaScript, and C++, with future plans to expand its language support based on community demand. With dynamic resource constraints, advanced error handling, and comprehensive execution metrics, the service is particularly well-suited for educational platforms, technical tests, automated testing platforms, and enterprise applications requiring reliable code execution features.

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

### Supported Languages
Gets all the current supported languages

**Endpoint**: `GET /languages`

**Response**:
```bash
{ "languages": [ "python", "cpp", "javascript", "java" ] }
```

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
