package emailsender

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/smtp"
	"regexp"

	amqp "github.com/rabbitmq/amqp091-go"
)

var emailPattern = regexp.MustCompile(`(?:[a-z0-9!#$%&'*+=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+=?^_{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`)

type Option func(s *Service) error

func WithSMTP(host, port, from, password string) Option {
	return func(svc *Service) error {
		svc.smtpcreds = &smtpcreds{
			Addr:   fmt.Sprintf("%s:%s", host, port),
			Auth:   smtp.PlainAuth("", from, password, host),
			Sender: from,
		}
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

func WithLogger(log *slog.Logger) Option {
	return func(svc *Service) error {
		if log == nil {
			return errors.New("logger is nil")
		}
		svc.log = log
		return nil
	}
}

type Service struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue

	log *slog.Logger

	smtpcreds *smtpcreds
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

func (svc *Service) Start() error {
	msgs, err := svc.ch.Consume(
		svc.q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	log := svc.log.With(
		slog.String("queue", svc.q.Name),
	)
	for msg := range msgs {
		buff := bytes.NewBuffer(msg.Body)
		log.Info("recieve message", slog.String("message", buff.String()))

		emailToSend := emailPattern.FindString(buff.String())

		log.Info("fetch email to send", slog.String("email", emailToSend))

		//err := smtp.SendMail(
		//	svc.smtpcreds.Addr,
		//	svc.smtpcreds.Auth,
		//	svc.smtpcreds.Sender,
		//	[]string{emailToSend},
		//	buff.Bytes(),
		//)
		//err := false
		//if err != nil {
		//	log.Error("failed to send email",
		//		slog.String("email", emailToSend),
		//		slog.String("error", err.Error()),
		//	)
		//}
	}
	return nil
}

func (svc *Service) Close() error {
	return svc.conn.Close()
}

type smtpcreds struct {
	Addr   string
	Auth   smtp.Auth
	Sender string
}
