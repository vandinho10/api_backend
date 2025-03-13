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

// Logger estrutura os atributos do logger.
//
// Contém informações sobre o caminho do log, prefixo, nível de log,
// configuração para exibição no console e o arquivo de log em uso.
type Logger struct {
	logPath      string
	logPrefix    string
	logLevel     string
	logToConsole bool
	logFile      *os.File
}

// Init inicializa a instância do logger.
//
// Carrega as variáveis de ambiente, configura o diretório de logs e define o nível de log.
// Se ocorrer erro ao criar o diretório de logs ou abrir o arquivo, ele utiliza stdout.
//
// Não retorna valores.
func Init() {
	once.Do(func() {
		_ = godotenv.Load("/app/.env")
		loggerInstance = &Logger{
			logPath:      getEnv("LOG_PATH", "logs"),
			logPrefix:    getEnv("LOGFILE_PREFIX", "app"),
			logLevel:     getEnv("LOG_LEVEL", "INFO"),
			logToConsole: getEnv("LOG_TO_CONSOLE", "true") == "true",
		}

		// Cria o diretório de logs se não existir
		if err := os.MkdirAll(loggerInstance.logPath, 0755); err != nil {
			log.Printf("Erro ao criar diretório de logs: %v. Usando stdout apenas.", err)
			loggerInstance.logFile = nil
		}

		logFilePath := fmt.Sprintf("%s/%s.log", loggerInstance.logPath, loggerInstance.logPrefix)
		var writers []io.Writer

		// Abre o arquivo de log e adiciona ao writer
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

// GetLogger retorna a instância do logger.
//
// Se o logger ainda não foi inicializado, ele chama Init() antes de retornar.
//
// Retorna:
//   - *Logger: Instância do logger.
func GetLogger() *Logger {
	if loggerInstance == nil {
		Init()
	}
	return loggerInstance
}

// logMessage registra uma mensagem de log com um determinado nível.
//
// Parâmetros:
//   - level (string): Nível do log (DEBUG, INFO, WARN, ERROR).
//   - message (string): Mensagem de log formatada.
//   - args (...interface{}): Argumentos opcionais para formatação da mensagem.
//
// Não retorna valores.
func (l *Logger) logMessage(level, message string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}

	callerInfo := getCallerInfo()
	formattedMessage := fmt.Sprintf("[%s] %s - %s", level, callerInfo, fmt.Sprintf(message, args...))
	log.Println(formattedMessage)
}

// shouldLog verifica se o nível da mensagem deve ser registrado com base na configuração do logger.
//
// Parâmetros:
//   - level (string): Nível da mensagem a ser verificada.
//
// Retorna:
//   - bool: true se a mensagem deve ser registrada, false caso contrário.
func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{"DEBUG": 1, "INFO": 2, "WARN": 3, "ERROR": 4}

	logLevel, exists := levels[strings.ToUpper(l.logLevel)]
	if !exists {
		logLevel = 2 // Fallback para INFO se a configuração for inválida
	}

	return levels[strings.ToUpper(level)] >= logLevel
}

// getCallerInfo obtém informações sobre a função que chamou o logger.
//
// Retorna:
//   - string: Nome do pacote, função e linha de onde o log foi chamado.
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

// Métodos públicos do Logger para registrar logs com diferentes níveis.
func Debug(message string, args ...interface{}) { GetLogger().logMessage("DEBUG", message, args...) }
func Info(message string, args ...interface{})  { GetLogger().logMessage("INFO", message, args...) }
func Warn(message string, args ...interface{})  { GetLogger().logMessage("WARN", message, args...) }
func Error(message string, args ...interface{}) { GetLogger().logMessage("ERROR", message, args...) }
func Fatal(message string, args ...interface{}) { GetLogger().logMessage("FATAL", message, args...) }
func Panic(message string, args ...interface{}) { GetLogger().logMessage("PANIC", message, args...) }

// LoggerMiddleware cria um middleware para registrar logs de requisições HTTP no Gin.
//
// O middleware registra informações como método da requisição, IP do cliente, URI solicitada
// e tempo de duração da requisição. Se houver erros no contexto, eles serão registrados.
//
// Retorna:
//   - gin.HandlerFunc: Função middleware para ser usada no Gin.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // Processa a requisição
		duration := time.Since(start)

		// Obtém o IP do cliente, priorizando cabeçalhos de proxies reversos.
		clientIP := c.GetHeader("CF-Connecting-IP")
		if clientIP == "" {
			clientIP = c.GetHeader("X-Forwarded-For")
			if clientIP == "" {
				clientIP = c.ClientIP()
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

// Close fecha o arquivo de log caso esteja aberto.
//
// Não retorna valores.
func Close() {
	if loggerInstance != nil && loggerInstance.logFile != nil {
		loggerInstance.logFile.Close()
	}
}

// getEnv obtém uma variável de ambiente, retornando um valor padrão se não estiver definida.
//
// Parâmetros:
//   - key (string): Nome da variável de ambiente.
//   - defaultValue (string): Valor padrão caso a variável não esteja definida.
//
// Retorna:
//   - string: Valor da variável de ambiente ou o valor padrão.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
