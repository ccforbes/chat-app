package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/handlers"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/gateway/sessions"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

//Director a
type Director func(r *http.Request)

//CustomDirector a
func CustomDirector(targets []*url.URL, ctx *handlers.HandlerCtx) Director {
	var counter int32
	counter = 0
	return func(r *http.Request) {
		targ := targets[rand.Int()%len(targets)] // how to load balance randomly
		//_targets, _ := rc.Get("ChatAddresses").Result()
		//targets := strings.Split(_targets, ",")
		//targ, _ := url.Parse(targets[int(counter)%len(targets)])
		atomic.AddInt32(&counter, 1) // note, to be extra safe, we'll need to use mutexes
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme
		var sessionState handlers.SessionState
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &sessionState)
		if err == nil {
			buffer, _ := json.Marshal(sessionState.User)
			r.Header.Set("X-User", string(buffer))
		}
	}
}

//main is the main entry point for the server
func main() {
	// get env variables
	addr := os.Getenv("ADDR")
	dsn := os.Getenv("DSN")
	messageAddr := strings.Split(os.Getenv("MESSAGESADDR"), ",")
	redisAddr := os.Getenv("REDISADDR")
	sessionKey := os.Getenv("SESSIONKEY")
	summaryAddr := strings.Split(os.Getenv("SUMMARYADDR"), ",")
	tlsCertPath := os.Getenv("TLSCERT")
	tlsKeyPath := os.Getenv("TLSKEY")

	if len(addr) == 0 {
		addr = ":443"
	}
	if tlsKeyPath == "" {
		fmt.Print("Missing Key")
		os.Exit(100)
	}
	if tlsCertPath == "" {
		fmt.Print("Missing Cert")
		os.Exit(100)
	}

	// connect rabbit
	client, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer client.Close()

	ch, err := client.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"msgs_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// connect redis
	redisStore := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	test, err := redisStore.Ping().Result()
	if err != nil {
		fmt.Printf("error pinging redis: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("successfuly connected redis: %s\n", test)
	}

	// connect mysql
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error initializing database: %v", err)
	}

	if err = db.Ping(); err != nil {
		fmt.Printf("error pinging database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("successfully connected database!:\n")
	}

	context := handlers.HandlerCtx{
		SigningKey:   sessionKey,
		SessionStore: sessions.NewRedisStore(redisStore, time.Hour),
		UserStore:    users.NewSQLStore(db),
		Notifier:     handlers.NewNotifier(),
	}
	context.SearchTree, err = context.UserStore.LoadExistingUsers()
	if err != nil {
		fmt.Printf("error adding users: %v", err)
	}

	defer db.Close()

	// start go routine to consume messages from rabbit mq
	go context.Notifier.WriteMessagesToSockets(msgs)

	summaryAddresses := []*url.URL{}
	for _, addr := range summaryAddr {
		summaryAddresses = append(summaryAddresses, &url.URL{Scheme: "http", Host: addr})
	}

	messageAddresses := []*url.URL{}
	for _, addr := range messageAddr {
		messageAddresses = append(messageAddresses, &url.URL{Scheme: "http", Host: addr})
	}

	summaryProx := &httputil.ReverseProxy{Director: CustomDirector(summaryAddresses, &context)}
	message1Prox := &httputil.ReverseProxy{Director: CustomDirector(messageAddresses, &context)}
	message2Prox := &httputil.ReverseProxy{Director: CustomDirector(messageAddresses, &context)}
	message3Prox := &httputil.ReverseProxy{Director: CustomDirector(messageAddresses, &context)}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)
	mux.HandleFunc("/v1/ws", context.WebSocketConnectionHandler)
	mux.Handle("/v1/channels", message1Prox)
	mux.Handle("/v1/channels/", message2Prox)
	mux.Handle("/v1/messages/", message3Prox)
	mux.Handle("/v1/summary", summaryProx)
	wrappedMux := handlers.NewCORSHeader(mux)
	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}
