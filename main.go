package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dotchev/sm-plugins/broker"
	"github.com/dotchev/sm-plugins/plugins"
	"github.com/dotchev/sm-plugins/sm"
	"github.com/dotchev/sm-plugins/sm/plugin/rest"
)

func main() {
	// start dummy broker
	broker := broker.Start()
	defer broker.Close()

	m := sm.NewServiceManager(&sm.Options{
		Plugins: []*rest.Plugin{
			plugins.DescriptionSetter(),
			plugins.CounterPlugin(),
		},
		BrokerURL: broker.URL,
	})

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", m))
}
