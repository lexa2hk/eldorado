package statistic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
	"github.com/romankravchuk/eldorado/internal/services"
	"github.com/romankravchuk/eldorado/internal/storages"
	"github.com/romankravchuk/eldorado/internal/storages/tasks"
	"github.com/romankravchuk/eldorado/internal/storages/tasks/pg"
)

const (
	mimeHeaders  = "MIME-version: 1.0;\nContent-Type: text/html; chartset=\"utf-8\";\n\n"
	headerFormat = "To: %s\nSubject: Welcome To Our Community\n%s\n\n"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("/email.html")
	if err != nil {
		slog.Error("failed to parse template", err)
		os.Exit(1)
	}
}

type Option func(s *Service) error

func WithLogger(log *slog.Logger) Option {
	return func(s *Service) error {
		if log == nil {
			return errors.New("logger is nil")
		}
		s.log = log
		return nil
	}
}

func WithRabbitMQ(amqpURI, queueName string) Option {
	return func(svc *Service) error {
		var err error
		svc.conn, err = amqp.Dial(amqpURI)
		if err != nil {
			return err
		}
		svc.ch, err = svc.conn.Channel()
		if err != nil {
			return err
		}
		svc.q, err = svc.ch.QueueDeclare(
			queueName,
			false,
			false,
			false,
			false,
			nil,
		)
		return err
	}
}

func WithTasksStorage(tasks tasks.Storage) Option {
	return func(s *Service) error {
		if tasks == nil {
			return services.ErrNilTasksStorage
		}

		s.tasks = tasks
		return nil
	}
}

func WithPosgresTasksStorage(url string) Option {
	return func(s *Service) error {
		db, err := storages.NewDBPool("postgres", url)
		if err != nil {
			return err
		}

		tasks, err := pg.New(db)
		if err != nil {
			return err
		}

		return WithTasksStorage(tasks)(s)
	}
}

func WithCron(schedule string) Option {
	return func(s *Service) error {
		s.cron = cron.New()
		s.schedule = schedule
		return nil
	}
}

type Service struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue

	log *slog.Logger

	cron     *cron.Cron
	schedule string

	tasks tasks.Storage
}

func New(opts ...Option) (*Service, error) {
	s := &Service{}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Service) Start() error {
	_, err := s.cron.AddFunc(s.schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.sendStatistic(ctx); err != nil {
			s.log.Error("failed to send statistic", slog.String("error", err.Error()))
			return
		}
		s.log.Info("statistic sent")
	})
	if err != nil {
		return err
	}
	s.cron.Run()
	return nil
}

func (s *Service) sendStatistic(ctx context.Context) error {
	tasks, err := s.tasks.UncompletedStatistic(ctx)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return nil
	}

	var buff bytes.Buffer
	tasksMap := make(map[string][]string, 0)

	for _, task := range tasks {
		tasksMap[task.Email] = append(tasksMap[task.Email], task.Title)
	}

	for e, tt := range tasksMap {
		buff.Write([]byte(fmt.Sprintf(
			headerFormat,
			e,
			mimeHeaders,
		)))

		if err := tmpl.Execute(&buff, tt); err != nil {
			return err
		}
	}

	err = s.ch.PublishWithContext(ctx,
		"",
		s.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/html",
			Body:        buff.Bytes(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
