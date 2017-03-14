package main

import _ "net/http/pprof"

import (
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"fmt"
	"github.com/docker/docker/api/types/events"
	"strings"
	"time"
	"log"
	"net/http"
	"encoding/json"
	"golang.org/x/net/websocket"
	"os"
	"os/signal"
)

func main() {
	websocketServer := newWebsocketServer()
	containerRepository := newContainerRepository()

	go func() {
		fs := http.FileServer(http.Dir("static"))
		http.Handle("/", fs)
		http.Handle("/ws", websocket.Handler(websocketServer.Handle))

		http.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(containerRepository.List())
		})
		log.Println(http.ListenAndServe(":8080", nil))
	}()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	context := context.Background()

	go loadExistingContainers(cli, context, containerRepository, websocketServer)

	go startListeningForEvents(cli, context, containerRepository, websocketServer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func loadExistingContainers(cli *client.Client, context context.Context, containerRepository *containerRepository, websocketServer *WebsocketServer) {
	containers, err := cli.ContainerList(context, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		go processExistingContainer(container, cli, context, containerRepository, websocketServer)
	}
}

func processExistingContainer(container types.Container, cli *client.Client, context context.Context, containerRepository *containerRepository, websocketServer *WebsocketServer) {
	containerJson, err := cli.ContainerInspect(context, container.ID)
	if err != nil {
		panic(err)
	}
	processStartEvent(containerJson, cli, context, containerRepository, websocketServer)
}

func startListeningForEvents(cli *client.Client, context context.Context, containerRepository *containerRepository, websocketServer *WebsocketServer) {
	options := types.EventsOptions{}

	eventsChannel, _ := cli.Events(context, options)

	for event :=range eventsChannel {
		if (event.Type == "container") {
			go processContainerEvent(event, cli, context, containerRepository, websocketServer)
		}
	}
}

func processContainerEvent(event events.Message, cli *client.Client, context context.Context, containerRepository *containerRepository, websocketServer *WebsocketServer) {

	container, err := cli.ContainerInspect(context, event.ID)
	if err != nil {
		panic(err)
	}

	if (event.Status == "start") {
		processStartEvent(container, cli, context, containerRepository, websocketServer)
	}
	if (event.Status == "die") {
		processDieEvent(container, containerRepository, websocketServer)
	}
}

func processDieEvent(containerJSON types.ContainerJSON, containerRepository *containerRepository, websocketServer *WebsocketServer) {
	fmt.Printf("Container %s has been stoped\n", containerJSON.Name)

	jsonStr, err:=json.Marshal(Message{Type:"die",Data:containerJSON.ID})
	if(err != nil){
		panic(err)
	}

	websocketServer.Broadcast(string(jsonStr))
	containerRepository.cancel(containerJSON.ID)
}

func processStartEvent(container types.ContainerJSON, cli *client.Client, context context.Context, containerRepository *containerRepository, websocketServer *WebsocketServer) {

	hasMaxAge, maxAge := getMaxAge(container.Config.Env);

	startedAt, err := time.Parse(time.RFC3339Nano, container.State.StartedAt)
	if err != nil {
		panic(err)
	}

	stopDuration := startedAt.Add(maxAge).Sub(time.Now())

	if (hasMaxAge) {
		fmt.Printf("Container %s will be stopped in %s\n", container.Name, stopDuration.String())

		timer := time.AfterFunc(stopDuration, func() {
			fmt.Printf("stopping container name: %s\n", container.Name)

			containerRepository.remove(container.ID)
			timeout := time.Duration(time.Minute)
			err := cli.ContainerStop(context, container.ID, &timeout)
			if err != nil {
				panic(err)
			}
			fmt.Printf("stoped container name: %s\n", container.Name)
		})

		containerTimer := &ContainerTimer{
			Id:container.ID,
			Name:container.Name,
			timer:timer,
			MaxAge:maxAge.Seconds(),
			StartedAt:startedAt.Unix(),
		}

		containerRepository.Add(containerTimer)
		json, err := json.Marshal(Message{Data:containerTimer, Type:"start"})
		if (err != nil) {
			panic(err)
		}
		websocketServer.Broadcast(string(json))

	}
}

func getMaxAge(environment []string) (hasMaxAge bool, duration time.Duration) {

	for _, variable := range environment {
		if (strings.HasPrefix(variable, "MAX_AGE=")) {

			parsedDuration, err := time.ParseDuration(strings.TrimPrefix(variable, "MAX_AGE="))
			if err != nil {
				panic(err)
			}
			hasMaxAge = true
			duration = parsedDuration
		}
	}
	return
}
