// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type CodeResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

type LanguageConfig struct {
    Compiler  string
    Filename  string
    RunCmd    string
    BuildArgs []string
    MainClass string   // Added for Java
}

var languageConfigs = map[string]LanguageConfig{
	"python": {
		Compiler:  "python3",
		Filename:  "main.py",
		RunCmd:    "python3",
		BuildArgs: []string{"-c", "import py_compile; py_compile.compile('{}')"},
	},
	"cpp": {
		Compiler:  "g++",
		Filename:  "main.cpp",
		RunCmd:    "./a.out",
		BuildArgs: []string{"-o", "a.out"},
	},
	"javascript": {
		Compiler:  "node",
		Filename:  "main.js",
		RunCmd:    "node",
		BuildArgs: []string{}, // Syntax check before execution
	},
	"java": {
        Compiler:  "javac",
        Filename:  "Main.java",    // Java file must match class name
        RunCmd:    "java",
        BuildArgs: []string{},     // Will be filled with the filename
        MainClass: "Main",         // Main class name
    },
}

type Server struct {
	router *chi.Mux
}

func NewServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	// Middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(30 * time.Second))

	s.router.Post("/execute", s.handleCodeExecution)
	//s.router.Get("/health", s.handleHealthCheck)
	s.router.Get("/languages", s.handleListLanguages) 
}

func (s *Server) handleListLanguages(w http.ResponseWriter, r *http.Request) {
    languages := make([]string, 0, len(languageConfigs))
    for lang := range languageConfigs {
        languages = append(languages, lang)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string][]string{
        "languages": languages,
    })
}

/*func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}*/

func (s *Server) handleCodeExecution(w http.ResponseWriter, r *http.Request) {
	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Println("req", req)
	response := executeCode(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func executeCode(code CodeRequest) CodeResponse {
	execDir := filepath.Join("/tmp", fmt.Sprintf("execution_%d", time.Now().UnixNano()))
	if err := os.MkdirAll(execDir, 0755); err != nil {
		return CodeResponse{Error: "Failed to create execution directory"}
	}
	defer os.RemoveAll(execDir)

	config, ok := languageConfigs[code.Language]
	if !ok {
		return CodeResponse{Error: "Unsupported language"}
	}

	codePath := filepath.Join(execDir, config.Filename)
	if err := os.WriteFile(codePath, []byte(code.Code), 0644); err != nil {
		return CodeResponse{Error: "Failed to write code file"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if code.Language == "javascript" {
		runCmd := exec.CommandContext(ctx, config.RunCmd, codePath)
		runCmd.Dir = execDir

		output, err := runCmd.CombinedOutput()
		if err != nil {
			return CodeResponse{Error: fmt.Sprintf("Execution error: %v\n%s", err, string(output))}
		}
		return CodeResponse{Output: string(output)}
	}

	// For Python, we don't need to compile, just run the script directly
	if code.Language == "python" {
		runCmd := exec.CommandContext(ctx, config.RunCmd, codePath)
		runCmd.Dir = execDir

		// Capture both stdout and stderr for better error reporting
		output, err := runCmd.CombinedOutput()
		if err != nil {
			return CodeResponse{Error: fmt.Sprintf("Execution error: %v\n%s", err, string(output))}
		}

		return CodeResponse{Output: string(output)}
	}

	 // Special handling for Java
	 if code.Language == "java" {
        // Compile Java code
        buildCmd := exec.CommandContext(ctx, config.Compiler, codePath)
        buildCmd.Dir = execDir
        
        output, err := buildCmd.CombinedOutput()
        if err != nil {
            return CodeResponse{Error: fmt.Sprintf("Compilation error: %v\n%s", err, string(output))}
        }

        // Run Java class
        runCmd := exec.CommandContext(ctx, config.RunCmd, config.MainClass)
        runCmd.Dir = execDir
		runCmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid: 0,
		}
        runOutput, err := runCmd.CombinedOutput()
        if err != nil {
            return CodeResponse{Error: fmt.Sprintf("Execution error: %v\n%s", err, string(runOutput))}
        }

        return CodeResponse{Output: string(runOutput)}
    }

	// Handle other languages (e.g., C++) as before
	buildArgs := make([]string, len(config.BuildArgs))
	copy(buildArgs, config.BuildArgs)
	for i, arg := range buildArgs {
		buildArgs[i] = strings.Replace(arg, "{}", codePath, -1)
	}

	if code.Language == "cpp" {
		buildArgs = append([]string{codePath}, buildArgs...)
	}

	buildCmd := exec.CommandContext(ctx, config.Compiler)
	buildCmd.Dir = execDir
	buildCmd.Args = append([]string{config.Compiler}, buildArgs...)

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return CodeResponse{Error: fmt.Sprintf("Compilation error: %v\n%s", err, string(output))}
	}


	runCmd := exec.CommandContext(ctx, filepath.Join(execDir, config.RunCmd))
	runCmd.Dir = execDir
	runCmd.SysProcAttr = &syscall.SysProcAttr{
        Setpgid: true,
        Pgid: 0,
    }
	runOutput, err := runCmd.CombinedOutput()
	if err != nil {
		return CodeResponse{Error: fmt.Sprintf("Execution error: %v\n%s", err, string(runOutput))}
	}

	return CodeResponse{Output: string(runOutput)}
}




func main() {
	server := NewServer()
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", server.router))
}
