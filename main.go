package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"golang.org/x/time/rate"
)

// Configuration structures
type Config struct {
	App       AppConfig       `json:"app"`
	Recaptcha RecaptchaConfig `json:"recaptcha"`
	LogFile   string          `json:"logFile"`
	Timeout   int             `json:"timeout"`
	Routers   []RouterConfig  `json:"routers"`
	Security  SecurityConfig  `json:"security"`
}

type AppConfig struct {
	Title        string `json:"title"`
	ContactEmail string `json:"contactEmail"`
	LogoURL      string `json:"logoUrl"`
	FaviconURL   string `json:"faviconUrl"`
	Disclaimer   string `json:"disclaimer"`
}

type RecaptchaConfig struct {
	Enabled   bool   `json:"enabled"`
	SiteKey   string `json:"siteKey"`
	SecretKey string `json:"secretKey"`
}

type RouterConfig struct {
	Name        string           `json:"name"`
	Title       string           `json:"title"`
	OSType      string           `json:"osType"`
	Location    string           `json:"location"`
	IPv4Enabled bool             `json:"ipv4Enabled"`
	IPv6Enabled bool             `json:"ipv6Enabled"`
	Connection  ConnectionConfig `json:"connection"`
}

type ConnectionConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
}

type SecurityConfig struct {
	RateLimit      RateLimitConfig `json:"rateLimit"`
	SecureMode     bool            `json:"secureMode"`
	AllowedOrigins []string        `json:"allowedOrigins"`
}

type RateLimitConfig struct {
	WindowMs int `json:"windowMs"`
	Max      int `json:"max"`
}

// Request/Response structures
type ExecuteRequest struct {
	Query    string `json:"query" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
	Addr     string `json:"addr"`
	Router   string `json:"router" binding:"required"`
	Token    string `json:"token"`
}

type ExecuteResponse struct {
	Success   bool   `json:"success"`
	Router    string `json:"router"`
	Command   string `json:"command"`
	Output    string `json:"output"`
	Timestamp string `json:"timestamp"`
}

type RouterInfo struct {
	Value       string `json:"value"`
	Text        string `json:"text"`
	Location    string `json:"location"`
	IPv4Enabled bool   `json:"ipv4Enabled"`
	IPv6Enabled bool   `json:"ipv6Enabled"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// NEW: Streaming response structure
type StreamResponse struct {
	Type    string `json:"type"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Command string `json:"command,omitempty"`
	Router  string `json:"router,omitempty"`
}

// Global variables
var (
	config      Config
	rateLimiter *rate.Limiter
	logMutex    sync.Mutex
)

// Helper function per saltare righe inutili
func shouldSkipLine(line string) bool {
	return strings.HasSuffix(line, "@juno01>") ||
		strings.HasSuffix(line, "@juno01-mx204>") ||
		strings.HasSuffix(line, "ellegi@juno01>") ||
		strings.HasSuffix(line, "kappa@juno01-mx204>") ||
		strings.HasSuffix(line, "netengine01>") ||
		strings.HasPrefix(line, "<netengine01") ||
		strings.Contains(line, "Info: The max number of VTY users") ||
		strings.Contains(line, "The current login time is") ||
		strings.Contains(line, "The last login time is") ||
		strings.Contains(line, "through SSH") ||
		strings.Contains(line, "from 2A02:4460:1:1::19 through SSH") ||
		strings.Contains(line, "and the number of current VTY users on line is")
}

// NEW: Funzione per streaming SSH con output in tempo reale
func executeSSHCommandStreaming(host string, port int, username, password, command string, w http.ResponseWriter) {
	// Set headers per streaming
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Funzione helper per inviare dati
	sendData := func(resp StreamResponse) {
		data, _ := json.Marshal(resp)
		fmt.Fprintf(w, "%s\n", data)
		flusher.Flush()
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
		Config: ssh.Config{
			KeyExchanges: []string{
				"diffie-hellman-group-exchange-sha256",
				"diffie-hellman-group-exchange-sha1",
				"diffie-hellman-group14-sha256",
				"diffie-hellman-group14-sha1",
				"diffie-hellman-group1-sha1",
				"curve25519-sha256",
				"curve25519-sha256@libssh.org",
				"ecdh-sha2-nistp256",
				"ecdh-sha2-nistp384",
				"ecdh-sha2-nistp521",
			},
			Ciphers: []string{
				"aes128-ctr", "aes192-ctr", "aes256-ctr",
				"aes128-cbc", "aes192-cbc", "aes256-cbc",
				"3des-cbc", "blowfish-cbc", "cast128-cbc",
				"arcfour256", "arcfour128", "arcfour",
			},
			MACs: []string{
				"hmac-sha2-256", "hmac-sha2-512", "hmac-sha1",
				"hmac-sha1-96", "hmac-md5", "hmac-md5-96", "hmac-ripemd160",
			},
		},
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Connecting to %s for streaming...", addr)

	// Invia messaggio di inizio
	sendData(StreamResponse{Type: "start", Command: command})

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		sendData(StreamResponse{Type: "error", Error: fmt.Sprintf("SSH connection failed: %v", err)})
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		sendData(StreamResponse{Type: "error", Error: fmt.Sprintf("SSH session failed: %v", err)})
		return
	}
	defer session.Close()

	// Setup pipes per lettura real-time
	stdout, err := session.StdoutPipe()
	if err != nil {
		sendData(StreamResponse{Type: "error", Error: fmt.Sprintf("Stdout pipe failed: %v", err)})
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		sendData(StreamResponse{Type: "error", Error: fmt.Sprintf("Stderr pipe failed: %v", err)})
		return
	}

	// Start command
	if err := session.Start(command); err != nil {
		sendData(StreamResponse{Type: "error", Error: fmt.Sprintf("Command start failed: %v", err)})
		return
	}

	// Channel per sincronizzare la fine della lettura
	done := make(chan bool, 2)

	// Leggi stdout in real-time
	go func() {
		defer func() { done <- true }()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !shouldSkipLine(line) {
				sendData(StreamResponse{Type: "data", Data: line})
			}
		}
	}()

	// Leggi stderr in real-time
	go func() {
		defer func() { done <- true }()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !shouldSkipLine(line) {
				sendData(StreamResponse{Type: "data", Data: line})
			}
		}
	}()

	// Aspetta che il comando finisca
	go func() {
		session.Wait()
		done <- true
	}()

	// Timeout di 5 minuti
	timeout := time.After(5 * time.Minute)
	completedReaders := 0

	for completedReaders < 2 {
		select {
		case <-done:
			completedReaders++
		case <-timeout:
			sendData(StreamResponse{Type: "error", Error: "Command timeout after 5 minutes"})
			session.Signal(ssh.SIGTERM)
			return
		}
	}

	// Invia messaggio di completamento
	sendData(StreamResponse{Type: "complete"})
}

// SSH Client with improved router detection and command execution (ORIGINAL)
func executeSSHCommand(host string, port int, username, password, command string) (string, error) {
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
		Config: ssh.Config{
			KeyExchanges: []string{
				"diffie-hellman-group-exchange-sha256",
				"diffie-hellman-group-exchange-sha1",
				"diffie-hellman-group14-sha256",
				"diffie-hellman-group14-sha1",
				"diffie-hellman-group1-sha1",
				"curve25519-sha256",
				"curve25519-sha256@libssh.org",
				"ecdh-sha2-nistp256",
				"ecdh-sha2-nistp384",
				"ecdh-sha2-nistp521",
			},
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr",
				"aes256-ctr",
				"aes128-cbc",
				"aes192-cbc",
				"aes256-cbc",
				"3des-cbc",
				"blowfish-cbc",
				"cast128-cbc",
				"arcfour256",
				"arcfour128",
				"arcfour",
			},
			MACs: []string{
				"hmac-sha2-256",
				"hmac-sha2-512",
				"hmac-sha1",
				"hmac-sha1-96",
				"hmac-md5",
				"hmac-md5-96",
				"hmac-ripemd160",
			},
		},
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Connecting to %s with legacy SSH algorithms...", addr)

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return "", fmt.Errorf("SSH connection failed: %v", err)
	}
	defer conn.Close()

	log.Printf("SSH connected successfully to %s", host)

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("SSH session failed: %v", err)
	}
	defer session.Close()

	// Gestione paging migliorata per entrambi i router types
	if !strings.Contains(command, "ping") && !strings.Contains(command, "trace") {
		command = command + " | no-more"
	}

	log.Printf("Executing command: %s", command)

	done := make(chan error, 1)
	var output []byte

	go func() {
		output, err = session.CombinedOutput(command)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Command error: %v", err)
			if len(output) > 0 {
				log.Printf("Returning partial output despite error")
				return cleanSSHOutput(string(output), host), nil
			}
			return "", fmt.Errorf("command failed: %v", err)
		}
		log.Printf("Command completed successfully, %d bytes output", len(output))
		rawOutput := string(output)
		log.Printf("Raw output first 200 chars: %q", rawOutput[:min(200, len(rawOutput))])
		cleaned := cleanSSHOutput(rawOutput, host)
		log.Printf("Cleaned output first 200 chars: %q", cleaned[:min(200, len(cleaned))])
		return cleaned, nil
	case <-time.After(60 * time.Second):
		log.Printf("Command timeout after 60 seconds")
		session.Close()
		if len(output) > 0 {
			return cleanSSHOutput(string(output), host), nil
		}
		return "", fmt.Errorf("command timeout")
	}
}

// ULTRA minimal cleaning - only ANSI codes and line endings
func cleanSSHOutput(output string, host string) string {
	// Remove ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cleaned := ansiRegex.ReplaceAllString(output, "")

	// Normalize line endings
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "\n")

	lines := strings.Split(cleaned, "\n")
	var result []string

	for _, line := range lines {
		// Keep original line, just trim whitespace
		line = strings.TrimSpace(line)

		// Skip ONLY these very specific patterns:
		if line == "" ||
			strings.HasSuffix(line, "@juno01>") ||
			strings.HasSuffix(line, "@juno01-mx204>") ||
			strings.HasSuffix(line, "ellegi@juno01>") ||
			strings.HasSuffix(line, "kappa@juno01-mx204>") ||
			strings.HasSuffix(line, "netengine01>") ||
			strings.HasPrefix(line, "<netengine01") ||
			strings.Contains(line, "Info: The max number of VTY users") ||
			strings.Contains(line, "The current login time is") ||
			strings.Contains(line, "The last login time is") ||
			strings.Contains(line, "through SSH") ||
			strings.Contains(line, "from 2A02:4460:1:1::19 through SSH") ||
			strings.Contains(line, "and the number of current VTY users on line is") {
			continue
		}

		result = append(result, line)
	}

	return strings.TrimSpace(strings.Join(result, "\n"))
}

// Corrected command generation with proper Junos syntax
func generateCommand(query, protocol, addr string, routerConfig RouterConfig) (string, error) {
	isIPv6 := protocol == "IPv6"

	switch routerConfig.OSType {
	case "huawei":
		switch query {
		case "bgp":
			if isIPv6 {
				return fmt.Sprintf("display bgp ipv6 routing-table %s", addr), nil
			}
			return fmt.Sprintf("display bgp routing-table %s", addr), nil
		case "advertised-routes":
			if isIPv6 {
				return fmt.Sprintf("display bgp ipv6 routing-table peer %s advertised-routes", addr), nil
			}
			return fmt.Sprintf("display bgp routing-table peer %s advertised-routes", addr), nil
		case "unicast neighbors":
			if isIPv6 {
				return "display bgp ipv6 peer verbose", nil
			}
			return "display bgp peer verbose", nil
		case "summary":
			if isIPv6 {
				return "display bgp ipv6 peer", nil
			}
			return "display bgp peer", nil
		case "ping":
			if isIPv6 {
				return fmt.Sprintf("ping ipv6 %s", addr), nil
			}
			return fmt.Sprintf("ping %s", addr), nil
		case "trace":
			if isIPv6 {
				return fmt.Sprintf("tracert ipv6 %s", addr), nil
			}
			return fmt.Sprintf("tracert %s", addr), nil
		}

	case "junos":
		switch query {
		case "bgp":
			// CORRECTED: Use show route X instead of show route protocol bgp X
			return fmt.Sprintf("show route %s", addr), nil
		case "advertised-routes":
			// Advertised routes command correct as is
			return fmt.Sprintf("show route advertising-protocol bgp %s", addr), nil
		case "unicast neighbors":
			return "show bgp neighbor", nil
		case "summary":
			return "show bgp summary", nil
		case "ping":
			return fmt.Sprintf("ping count 5 %s", addr), nil
		case "trace":
			if isIPv6 {
				return fmt.Sprintf("traceroute %s", addr), nil
			}
			return fmt.Sprintf("traceroute %s as-number-lookup", addr), nil
		}
	}

	return "", fmt.Errorf("unsupported query type or router OS")
}

// Logging
func logCommand(ip, router, command string, success bool) {
	logMutex.Lock()
	defer logMutex.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("%s - IP: %s - Router: %s - Command: %s - Success: %t\n",
		timestamp, ip, router, command, success)

	log.Print(strings.TrimSpace(logEntry))

	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open log file: %v", err)
			return
		}
		defer file.Close()
		file.WriteString(logEntry)
	}
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rateLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error: "Too many requests, please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func loadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&config)
}

func verifyRecaptcha(token string) (bool, error) {
	return !config.Recaptcha.Enabled || token != "", nil
}

// API Handlers
func getRoutersHandler(c *gin.Context) {
	routers := make([]RouterInfo, len(config.Routers))
	for i, router := range config.Routers {
		routers[i] = RouterInfo{
			Value:       router.Name,
			Text:        router.Title,
			Location:    router.Location,
			IPv4Enabled: router.IPv4Enabled,
			IPv6Enabled: router.IPv6Enabled,
		}
	}
	c.JSON(http.StatusOK, routers)
}

// NEW: Streaming handler
func executeStreamingHandler(c *gin.Context) {
	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	clientIP := c.ClientIP()

	// Validazione
	needsAddress := !contains([]string{"summary", "unicast neighbors"}, req.Query)
	if needsAddress && req.Addr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Address is required for this query type"})
		return
	}

	valid, err := verifyRecaptcha(req.Token)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "reCAPTCHA verification failed"})
		return
	}

	var routerConfig RouterConfig
	found := false
	for _, router := range config.Routers {
		if router.Name == req.Router {
			routerConfig = router
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid router selection"})
		return
	}

	command, err := generateCommand(req.Query, req.Protocol, req.Addr, routerConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Starting streaming command from %s: %s", clientIP, command)

	// Esegui comando in streaming
	executeSSHCommandStreaming(
		routerConfig.Connection.Host,
		routerConfig.Connection.Port,
		routerConfig.Connection.Username,
		routerConfig.Connection.Password,
		command,
		c.Writer,
	)

	logCommand(clientIP, req.Router, command, true)
}

// ORIGINAL execute handler
func executeHandler(c *gin.Context) {
	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	clientIP := c.ClientIP()

	needsAddress := !contains([]string{"summary", "unicast neighbors"}, req.Query)
	if needsAddress && req.Addr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Address is required for this query type"})
		return
	}

	valid, err := verifyRecaptcha(req.Token)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "reCAPTCHA verification failed"})
		return
	}

	var routerConfig RouterConfig
	found := false
	for _, router := range config.Routers {
		if router.Name == req.Router {
			routerConfig = router
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid router selection"})
		return
	}

	command, err := generateCommand(req.Query, req.Protocol, req.Addr, routerConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Executing command on %s: %s", routerConfig.Name, command)
	output, err := executeSSHCommand(
		routerConfig.Connection.Host,
		routerConfig.Connection.Port,
		routerConfig.Connection.Username,
		routerConfig.Connection.Password,
		command,
	)

	if err != nil {
		logCommand(clientIP, req.Router, command, false)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Command execution failed: %v", err)})
		return
	}

	logCommand(clientIP, req.Router, command, true)

	response := ExecuteResponse{
		Success:   true,
		Router:    routerConfig.Title,
		Command:   command,
		Output:    output,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"timestamp":  time.Now().Format(time.RFC3339),
		"version":    "2.0.15-simple-streaming",
		"routers":    len(config.Routers),
		"protocol":   "SSH",
		"algorithms": "Legacy Compatible",
		"features":   "Simple JSON Streaming, Real-time output",
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	rateValue := rate.Limit(float64(config.Security.RateLimit.Max) / float64(config.Security.RateLimit.WindowMs/1000))
	rateLimiter = rate.NewLimiter(rateValue, config.Security.RateLimit.Max)

	if config.Security.SecureMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Security.AllowedOrigins
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	r.GET("/", func(c *gin.Context) {
		c.File("./index.html")
	})
	r.Static("/public", "./public")

	api := r.Group("/api")
	api.Use(rateLimitMiddleware())
	{
		api.GET("/routers", getRoutersHandler)
		api.POST("/execute", executeHandler)           // Original endpoint
		api.POST("/execute-stream", executeStreamingHandler) // NEW: Streaming endpoint
		api.GET("/health", healthHandler)
	}

	srv := &http.Server{
		Addr:         ":3002",
		Handler:      r,
		ReadTimeout:  10 * time.Minute, // Long timeout for streaming
		WriteTimeout: 10 * time.Minute, // Long timeout for streaming
		IdleTimeout:  2 * time.Minute,
	}

	go func() {
		log.Printf("GoLine Looking Glass Server v2.0.15-SIMPLE-STREAMING")
		log.Printf("Server running on http://localhost:3002")
		log.Printf("Configured routers: %d", len(config.Routers))
		log.Printf("Protocol: SSH with LEGACY algorithm support")
		log.Printf("Features: Simple JSON Streaming for real-time output")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
