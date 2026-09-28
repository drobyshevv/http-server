package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Server представляет HTTP-сервер и управляет его жизненным циклом,
// TCP-соединениями и их корректным завершением.
type Server struct {
	cfg      Config
	listener net.Listener
	router   *Router
	// wg отслеживает все активные соединения сервера.
	wg sync.WaitGroup
	// shutdownCtx используется для уведомления активных соединений
	// о необходимости завершить работу.
	shutdownCtx context.Context
	// cancelShutdown отменяет shutdownCtx и запускает завершение
	// всех активных соединений.
	cancelShutdown context.CancelFunc
}

// Config содержит настройки TCP-сервера.
type Config struct {
	Network string
	Port    int
}

// NewServer создаёт новый Server с указанными настройками и роутером.
func NewServer(cfg Config, router *Router) *Server {
	// Создаём контекст жизненного цикла сервера.
	// Он будет отменён при graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		cfg:            cfg,
		router:         router,
		wg:             sync.WaitGroup{},
		shutdownCtx:    ctx,
		cancelShutdown: cancel,
	}
}

// Run запускает сервер и начинает принимать входящие TCP-соединения.
//
// Метод блокируется до тех пор, пока listener не будет закрыт.
func (s *Server) Run() error {
	listener, err := net.Listen(s.cfg.Network, fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener

	// Ожидаем новые TCP-соединения.
	// Цикл завершается после закрытия listener во время shutdown.
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			continue
		}
		s.wg.Go(func() { s.serveConn(s.shutdownCtx, conn) })
	}

	return nil
}

// Shutdown выполняет graceful shutdown сервера.
//
// Сервер перестаёт принимать новые соединения и ожидает завершения
// уже запущенных соединений. Если они не завершились за 5 секунд,
// возвращается ошибка тайм-аута.
func (s *Server) Shutdown() error {
	// Закрываем listener, чтобы перестать принимать новые соединения.
	// Это также разблокирует Accept() в Run().
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	// Ждём завершения всех активных соединений в отдельной горутине,
	// чтобы можно было установить максимальное время ожидания.
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	// Отменяем контекст сервера и уведомляем активные соединения
	// о необходимости завершиться.
	s.cancelShutdown()

	select {
	case <-done:
		fmt.Println("server stopped")
	case <-time.After(time.Second * 5):
		return fmt.Errorf("graceful shutdown timeout")
	}
	return nil
}
