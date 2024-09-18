package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/solution/choreography"
)

var (
	topicSize  = 5
	workerPool = 3
)

func main() {
	combatTopic := make(choreography.MessageBrokerTopic, topicSize)
	statusTopic := make(choreography.MessageBrokerTopic, topicSize)
	pokeballTopic := make(choreography.MessageBrokerTopic, topicSize)
	captureTopic := make(choreography.MessageBrokerTopic, topicSize)

	combatService := pokemon.NewCombatService()
	statusService := pokemon.NewStatusService()
	pokeballService := pokemon.NewPokeballService()

	statusWorker := choreography.NewStatusWorker(
		statusService,
		statusTopic,
		combatTopic,
	)

	combatWorker := choreography.NewCombatWorker(
		combatService,
		combatTopic,
		pokeballTopic,
	)

	pokeballWorker := choreography.NewPokeballWorker(
		pokeballService,
		pokeballTopic,
		captureTopic,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for range workerPool {
		go statusWorker.Run(ctx)
		go combatWorker.Run(ctx)
		go pokeballWorker.Run(ctx)
	}

	statusTopic <- choreography.Event{
		Type:    choreography.PokemonEncountered,
		Pokemon: pokemon.Mewtwo(),
		Trainer: pokemon.Sacha(),
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Waiting input to exit")
	reader.ReadString('\n')
}
