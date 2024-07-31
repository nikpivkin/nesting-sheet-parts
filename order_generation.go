package main

import (
	"fmt"
	"math/rand"
	"sort"
)

const (
	noImprovementLimit    = 10
	noNewIndividualsLimit = 20
)

type Individual struct {
	// gene represents the part number
	chromosome []int

	fitness float32
}

func (i Individual) Order() []int {
	return i.chromosome
}

func (i Individual) Fitness() float32 {
	return i.fitness
}

func (i Individual) Hash() string {
	// TODO: improve hashing
	return fmt.Sprintf("%v", i.chromosome)
}

func NewIndividual(numGenes int) Individual {
	individual := Individual{
		chromosome: rangeSlice(0, numGenes, 1),
	}
	shuffle(individual.chromosome)
	return individual
}

func shuffle(s []int) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

type GeneticAlgorithm struct {
	populationSize int
	numGenes       int
	elitismRate    float32
	mutationRate   float32

	fitnessFn  fitnessFunc
	best       Individual
	history    map[string]Individual
	population []Individual
}

type fitnessFunc func(Individual) float32

type GAOption func(*GeneticAlgorithm)

func WithPopulationSize(populationSize int) GAOption {
	return func(g *GeneticAlgorithm) {
		g.populationSize = populationSize
	}
}

func WithElitismRate(elitismRate float32) GAOption {
	return func(g *GeneticAlgorithm) {
		g.elitismRate = elitismRate
	}
}

func WithMutationRate(mutationRate float32) GAOption {
	return func(g *GeneticAlgorithm) {
		g.mutationRate = mutationRate
	}
}

func NewGeneticAlgorithm(numGenes int, fitnessFn fitnessFunc, options ...GAOption) *GeneticAlgorithm {
	ga := &GeneticAlgorithm{
		numGenes:       numGenes,
		populationSize: 10,
		elitismRate:    0.2,
		mutationRate:   0.1,
		fitnessFn:      fitnessFn,
		history:        make(map[string]Individual),
	}

	for _, option := range options {
		option(ga)
	}

	return ga
}

func (g *GeneticAlgorithm) Best() Individual {
	return g.best
}

func (g *GeneticAlgorithm) Run(numGenerations int) {

	g.population = g.newPopulation(g.populationSize)

	noImprovement := 0

	for generation := 0; generation < numGenerations; generation++ {
		println("Generation: ", generation)

		g.fitness()

		sort.Slice(g.population, func(i, j int) bool {
			return g.population[i].fitness > g.population[j].fitness
		})

		fmt.Printf("Best chromosome: %v, fitness: %f\n", g.population[0].chromosome, g.population[0].fitness)

		if generation == 0 {
			g.best = g.population[0]
		} else if currentBestFitness := g.population[0]; currentBestFitness.fitness > g.best.fitness {
			noImprovement = 0
			g.best = currentBestFitness
			fmt.Printf("Epoch: %d, best: %f\n", generation, g.best.fitness)
		} else {
			noImprovement++
		}

		if noImprovement > noImprovementLimit {
			fmt.Printf("No improvement for %d generations. Exiting.\n", noImprovementLimit)
			break
		}

		for j := 0; j < g.populationSize; j++ {
			g.history[g.population[j].Hash()] = g.population[j]
		}

		elite := g.elitism()

		newPopulation := make([]Individual, g.populationSize)
		copy(newPopulation, elite)

		noChanges := 0
		for count := 0; count < g.populationSize-len(elite); {
			if noChanges > noNewIndividualsLimit {
				fmt.Printf("No new individuals. Exiting.\n")
				return
			}

			var parent Individual

			// TODO: which rate to use?
			if rand.Float32() < 0.05 {
				parent = elite[rand.Intn(len(elite))]
			} else {
				parent = g.population[rand.Intn(len(g.population))]
			}
			child := g.crossover(parent)
			child = g.mutation(child)
			if _, exists := g.history[child.Hash()]; exists {
				noChanges++
				continue
			}

			newPopulation[g.populationSize-count-1] = child
			count++
			noChanges = 0
		}

		g.population = newPopulation
	}
}

func (g *GeneticAlgorithm) fitness() {
	for j := 0; j < len(g.population); j++ {
		g.population[j].fitness = g.fitnessFn(g.population[j])
	}
}

func (g *GeneticAlgorithm) newPopulation(size int) []Individual {
	population := make([]Individual, size)

	for i := 0; i < size; i++ {
		population[i] = NewIndividual(g.numGenes)
	}

	return population
}

func (g *GeneticAlgorithm) elitism() []Individual {
	size := int(float32(g.populationSize) * g.elitismRate)
	elite := make([]Individual, size)
	for i := 0; i < size; i++ {
		elite[i] = g.population[i]
	}
	return elite
}

func (g *GeneticAlgorithm) crossover(parent Individual) Individual {
	if len(parent.chromosome) == 0 {
		panic("the length of the parent cannot be less than 0")
	}
	pointIdx := rand.Intn(len(parent.chromosome))

	// TODO: better way?
	if pointIdx == 0 {
		pointIdx = 1
	} else if pointIdx == len(parent.chromosome)-1 {
		pointIdx = len(parent.chromosome) - 2
	}

	child := Individual{
		chromosome: swapSliceParts(parent.chromosome, pointIdx),
		fitness:    0,
	}

	return child
}

func (g *GeneticAlgorithm) mutation(individual Individual) Individual {
	if rand.Float32() < g.mutationRate {
		prev := make([]int, len(individual.chromosome))
		copy(prev, individual.chromosome)

		pointIdx := rand.Intn(len(individual.chromosome))
		left := rand.Float32() > 0.5

		if left && pointIdx > 0 {
			swap(individual.chromosome, pointIdx, pointIdx-1)
		} else if pointIdx < len(individual.chromosome)-1 {
			swap(individual.chromosome, pointIdx, pointIdx+1)
		}
	}
	return individual
}
