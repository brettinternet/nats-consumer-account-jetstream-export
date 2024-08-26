package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const streamName = "ORDERS"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var wg sync.WaitGroup
	nc1, js1, err := centralAccountStream(context.Background())
	if err != nil {
		return err
	}
	nc2, err := remoteAccountConsumer(context.Background(), &wg)
	if err != nil {
		return err
	}
	if err := publish(&wg, js1); err != nil {
		return err
	}

	wg.Wait()
	defer nc1.Close()
	defer nc2.Close()
	return errors.Join(nc1.LastError(), nc2.LastError())
}

func centralAccountStream(ctx context.Context) (*nats.Conn, jetstream.JetStream, error) {
	nc, err := nats.Connect(nats.DefaultURL, nats.UserCredentials("./nats/u1.creds"))
	if err != nil {
		return nil, nil, err
	}
	if nc.IsConnected() {
		fmt.Println("u1 connected to nats")
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, err
	}
	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  []string{"ORDERS.*"},
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		fmt.Println("stream creation error", err)
		return nil, nil, err
	}
	return nc, js, nil
}

func publish(wg *sync.WaitGroup, js jetstream.JetStream) error {
	msgCount := 10
	wg.Add(msgCount)
	for i := 1; i <= msgCount; i++ {
		payload := fmt.Sprintf("order %d", i)
		if _, err := js.PublishAsync("ORDERS.new", []byte(payload)); err != nil {
			return fmt.Errorf("message %d publish failed: %w", i, err)
		} else {
			fmt.Println("publishing", time.Now().Format(time.RFC3339), i)
		}
	}
	return nil
}

func remoteAccountConsumer(ctx context.Context, wg *sync.WaitGroup) (*nats.Conn, error) {
	nc, err := nats.Connect(nats.DefaultURL, nats.UserCredentials("./nats/u2.creds"))
	if err != nil {
		return nil, err
	}
	if nc.IsConnected() {
		fmt.Println("u2 connected to nats")
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	name := "orders-consumer"
	cfg := jetstream.ConsumerConfig{
		Name:          name,
		Durable:       name,
		FilterSubject: "ORDERS.new",
	}
	cons, err := js.CreateOrUpdateConsumer(ctx, streamName, cfg)
	if err != nil {
		return nil, err
	}
	handle := func(jmsg jetstream.Msg) {
		fmt.Println("consuming", time.Now().Format(time.RFC3339), string(jmsg.Data()))
		jmsg.Ack()
		defer wg.Done()
	}
	if _, err := cons.Consume(handle, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		fmt.Println("consume error", err)
	})); err != nil {
		return nil, err
	}

	go func() {
		for {
			if info, err := cons.Info(ctx); err != nil {
				fmt.Println("consumer info error", err)
			} else {
				fmt.Println("__consumer__",
					"delivered-sequence:", info.Delivered.Consumer, info.Delivered.Stream,
					"ackfloor-sequence:", info.Delivered.Consumer, info.Delivered.Stream,
					"sent:", info.AckFloor.Consumer,
					"redelivered:", info.NumRedelivered,
					"pending:", info.NumPending,
					"ack-pending:", info.NumAckPending,
					"pulls-waiting:", info.NumWaiting,
				)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	return nc, nil
}
