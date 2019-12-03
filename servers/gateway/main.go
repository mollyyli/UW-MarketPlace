package main

import (
	"UW-Marketplace/servers/gateway/handlers"
	"UW-Marketplace/servers/gateway/models/users"
	"UW-Marketplace/servers/gateway/sessions"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

//main is the main entry point for the server
//this functions first gets the vlaue of the ADDR environment variable
//if it's blank default to 80
//then creates a new router by calling NewServeMux on http
//tells the mux to call handlers.SummaryHandler function
//when the user requests the resource path /v1/summary
//then uses the nux to start the web server
//if error occurs uses fatal to report the error
//call ListenAndServe function on http to block and to run until killed

type Director func(r *http.Request)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func CustomDirector(targets []*url.URL, signingKey string, rc *sessions.RedisStore) Director {
	var counter int32
	counter = 0
	return func(r *http.Request) {
		// sid := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer")
		var currentSession handlers.SessionState
		_, err := sessions.GetState(r, signingKey, rc, &currentSession)
		// err := rc.Get(sessions.SessionID(sid), &currentSession)
		if &currentSession.User == nil || err != nil {
			r.Header.Del("X-User")
		} else {
			encodedUser, err := json.Marshal(currentSession.User)
			if err != nil {
				log.Println("error json marshal")
			}
			stringUser := string(encodedUser)
			r.Header.Set("X-User", stringUser)
		}
		targ := targets[int(counter)%len(targets)]
		atomic.AddInt32(&counter, 1)
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = "http"
	}
}

func main() {
	TLSCERT := os.Getenv("TLSCERT")
	TLSKEY := os.Getenv("TLSKEY")
	SESSIONKEY := os.Getenv("SESSIONKEY")
	REDISADDR := os.Getenv("REDISADDR")
	DSN := os.Getenv("DSN")
	RABBITMQADDR := os.Getenv("RABBITMQADDR")
	listingAddrs := strings.Split(os.Getenv("listingAddrs"), ", ")
	summaryAddrs := strings.Split(os.Getenv("summaryAddrs"), ", ")

	LISTINGADDR := make([]*url.URL, len(listingAddrs))
	for i := range LISTINGADDR {
		urlAddress, _ := url.Parse(listingAddrs[i])

		LISTINGADDR[i] = urlAddress
	}

	SUMMARYADDR := make([]*url.URL, len(summaryAddrs))
	for i := range summaryAddrs {
		urlAddress, _ := url.Parse(summaryAddrs[i])
		SUMMARYADDR[i] = urlAddress
	}

	if len(TLSCERT) == 0 {
		log.Fatal("TLSCERT is not set")
		os.Exit(1)
	}
	if len(TLSKEY) == 0 {
		log.Fatal("TLSKEY is not set")
		os.Exit(1)
	}
	if len(REDISADDR) == 0 {
		REDISADDR = "172.17.0.1:6379"
	}
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}
	red := redis.NewClient(&redis.Options{
		Addr: REDISADDR,
	})
	redisStore := sessions.NewRedisStore(red, time.Hour)
	listingProxy := &httputil.ReverseProxy{Director: CustomDirector(LISTINGADDR, SESSIONKEY, redisStore)}
	summaryProxy := &httputil.ReverseProxy{Director: CustomDirector(SUMMARYADDR, SESSIONKEY, redisStore)}
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Println("error opening database")
	}
	db.SetMaxIdleConns(0)
	mysql := &users.MySQLConnection{Client: db}
	ctx := handlers.Context{
		SigningKey:   SESSIONKEY,
		SessionStore: redisStore,
		UserStore:    mysql,
		SockStore:    *new(handlers.SocketStore),
	}

	conn, err := amqp.Dial("amqp://" + RABBITMQADDR + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message-queue", // name
		true,            // durable
		false,           // delete when usused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for message := range msgs {
			var msgObj struct {
				userIDs []int
			}
			json.Unmarshal(message.Body, &msgObj)
			if len(msgObj.userIDs) == 0 {
				ctx.SockStore.WriteToAllConnections(1, message.Body)
			} else {
				for _, userID := range msgObj.userIDs {
					ctx.SockStore.Connections[int64(userID)].WriteMessage(1, message.Body)
				}
			}
			// if message.Body {

			// }
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionHandler)
	mux.Handle("/v1/listings", listingProxy)
	mux.Handle("/v1/listings/", listingProxy)
	// mux.Handle("/v1/channels", messageeProxy)
	mux.Handle("/v1/summary", summaryProxy)
	mux.HandleFunc("/v1/ws", ctx.SocketHandler)
	log.Printf("server is listening at http://%s", addr)
	wrap := handlers.NewCors(&handlers.CorsHandler{Handler: mux})
	log.Fatal(http.ListenAndServeTLS(addr, TLSCERT, TLSKEY, wrap))
}
