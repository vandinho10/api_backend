// pwd: /app/logger/logger.go
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	loggerInstance *Logger
	once           sync.Once
)

type Logger struct {
	logPath      string
	logPrefix    string
	logLevel     string
	logToConsole bool
	logFile      *os.File
}

// Init inicializa o logger
func Init() {
	once.Do(func() {
		_ = godotenv.Load("/app/.env")
		loggerInstance = &Logger{
			logPath:      getEnv("LOG_PATH", "logs"),
			logPrefix:    getEnv("LOGFILE_PREFIX", "app"),
			logLevel:     getEnv("LOG_LEVEL", "INFO"),
			logToConsole: getEnv("LOG_TO_CONSOLE", "true") == "true",
		}

		if err := os.MkdirAll(loggerInstance.logPath, 0755); err != nil {
			log.Printf("Erro ao criar diretório de logs: %v. Usando stdout apenas.", err)
			loggerInstance.logFile = nil // Evita tentativa de escrever em um arquivo inexistente
		}

		logFilePath := fmt.Sprintf("%s/%s.log", loggerInstance.logPath, loggerInstance.logPrefix)
		var writers []io.Writer
		// writers = append(writers, fileLogs)
		// if loggerInstance.logToConsole {
		// 	writers = append(writers, os.Stdout)
		// }

		fileLogs, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Erro ao abrir arquivo de log: %v. Usando stdout apenas.", err)
		} else {
			loggerInstance.logFile = fileLogs
			writers = append(writers, fileLogs)
		}

		log.SetOutput(io.MultiWriter(writers...))
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	})
}

// GetLogger retorna a instância do logger
func GetLogger() *Logger {
	if loggerInstance == nil {
		Init()
	}
	return loggerInstance
}

// logMessage é uma função interna para registrar mensagens de log
func (l *Logger) logMessage(level, message string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}

	callerInfo := getCallerInfo()
	formattedMessage := fmt.Sprintf("[%s] %s - %s", level, callerInfo, fmt.Sprintf(message, args...))
	log.Println(formattedMessage)
}

// shouldLog verifica se o nível da mensagem deve ser registrado
func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{"DEBUG": 1, "INFO": 2, "WARN": 3, "ERROR": 4}

	logLevel, exists := levels[strings.ToUpper(l.logLevel)]
	if !exists {
		logLevel = 2 // Fallback para INFO se a configuração for inválida
	}

	return levels[strings.ToUpper(level)] >= logLevel
}

// getCallerInfo retorna o nome do pacote, função e linha de onde o log foi chamado
func getCallerInfo() string {
	for i := 2; i < 10; i++ { // Percorre a pilha até encontrar uma função fora do logger
		pc, _, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcName := runtime.FuncForPC(pc).Name()
		if !strings.Contains(funcName, "logger.") { // Evita capturar funções internas do logger
			return fmt.Sprintf("%s:%d", funcName, line)
		}
	}
	return "Unknown"
}

// Métodos públicos do Logger
func Debug(message string, args ...interface{}) { GetLogger().logMessage("DEBUG", message, args...) }
func Info(message string, args ...interface{})  { GetLogger().logMessage("INFO", message, args...) }
func Warn(message string, args ...interface{})  { GetLogger().logMessage("WARN", message, args...) }
func Error(message string, args ...interface{}) { GetLogger().logMessage("ERROR", message, args...) }
func Fatal(message string, args ...interface{}) { GetLogger().logMessage("FATAL", message, args...) }
func Panic(message string, args ...interface{}) { GetLogger().logMessage("PANIC", message, args...) }

// LoggerMiddleware registra logs de requisições HTTP do Gin
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // Processa a requisição
		duration := time.Since(start)

		clientIP := c.GetHeader("CF-Connecting-IP") // Prioridade para Cloudflare
		if clientIP == "" {
			clientIP = c.GetHeader("X-Forwarded-For") // Verifica proxy reverso
			if clientIP == "" {
				clientIP = c.ClientIP() // Fallback padrão
			}
		}

		msg := fmt.Sprintf("|%d |%s |%s |%s |%s |", c.Writer.Status(), c.Request.Method, clientIP, c.Request.RequestURI, duration)
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				msg += fmt.Sprintf(" | ERROR: %s", err.Error())
			}
			Error("%s", msg)
		} else {
			Info("%s", msg)
		}
	}
}

// Close fecha o arquivo de log
func Close() {
	if loggerInstance != nil && loggerInstance.logFile != nil {
		loggerInstance.logFile.Close()
	}

}

// getEnv retorna a variável de ambiente ou um valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
