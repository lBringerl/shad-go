//go:build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type operation func(stack []int) ([]int, error)

func sum(a int, b int) (int, error) {
	return a + b, nil
}

func sub(a int, b int) (int, error) {
	return a - b, nil
}

func mult(a int, b int) (int, error) {
	return a * b, nil
}

func div(a int, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

func arithmetic(symbol string) operation {
	var op func(int, int) (int, error)
	switch symbol {
	case "+":
		op = sum
	case "-":
		op = sub
	case "*":
		op = mult
	case "/":
		op = div
	default:
		panic(fmt.Sprintf("unknown symbol %s", symbol))
	}

	return func(stack []int) ([]int, error) {
		if len(stack) < 2 {
			return nil, fmt.Errorf("need at least 2 values on stack for operation %s, got %d", symbol, len(stack))
		}
		result, err := op(stack[len(stack)-2], stack[len(stack)-1])
		if err != nil {
			return nil, err
		}
		stack = stack[:len(stack)-2]
		stack = append(stack, result)

		return stack, nil
	}
}

func dup(stack []int) ([]int, error) {
	if len(stack) < 1 {
		return nil, fmt.Errorf("need at least 1 value on stack for dup, got %d", len(stack))
	}
	stack = append(stack, stack[len(stack)-1])

	return stack, nil
}

func over(stack []int) ([]int, error) {
	if len(stack) < 2 {
		return nil, fmt.Errorf("need at least 2 values on stack for over, got %d", len(stack))
	}
	stack = append(stack, stack[len(stack)-2])

	return stack, nil
}

func drop(stack []int) ([]int, error) {
	if len(stack) < 1 {
		return nil, fmt.Errorf("need at least 1 value on stack for drop, got %d", len(stack))
	}
	stack = stack[:len(stack)-1]

	return stack, nil
}

func swap(stack []int) ([]int, error) {
	if len(stack) < 2 {
		return nil, fmt.Errorf("need at least 2 values on stack for swap, got %d", len(stack))
	}
	stack[len(stack)-1], stack[len(stack)-2] = stack[len(stack)-2], stack[len(stack)-1]

	return stack, nil
}

type Evaluator struct {
	stack       []int
	operations  map[string]operation
	definitions map[string][]string
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		operations: map[string]operation{
			"+":    arithmetic("+"),
			"-":    arithmetic("-"),
			"*":    arithmetic("*"),
			"/":    arithmetic("/"),
			"dup":  dup,
			"over": over,
			"drop": drop,
			"swap": swap,
		},
		definitions: make(map[string][]string),
	}
}

func (e *Evaluator) putOnStack(number int) {
	e.stack = append(e.stack, number)
}

func (e *Evaluator) addDefinition(splitedRow []string) error {
	if len(splitedRow) < 2 {
		return fmt.Errorf("need at least 2 words for definition, got %d", len(splitedRow))
	}

	definitionWord := splitedRow[0]
	_, err := strconv.Atoi(definitionWord)
	if err == nil {
		return fmt.Errorf("cannot redefine numbers")
	}
	definitionValue := e.extractDefinitions(splitedRow[1:])
	e.definitions[definitionWord] = definitionValue

	return nil
}

func (e *Evaluator) extractDefinitions(splitedRow []string) []string {
	extractedRow := make([]string, 0)
	for _, word := range splitedRow {
		definition, exists := e.definitions[word]
		if !exists {
			extractedRow = append(extractedRow, word)
			continue
		}
		extractedRow = append(extractedRow, definition...)
	}

	return extractedRow
}

func (e *Evaluator) processRow(row string) error {
	var err error

	splitedRow := strings.Split(strings.ToLower(row), " ")

	if len(splitedRow) == 0 {
		return nil
	}
	if len(splitedRow) > 2 && splitedRow[0] == ":" && splitedRow[len(splitedRow)-1] == ";" {
		err = e.addDefinition(splitedRow[1 : len(splitedRow)-1])
		if err != nil {
			return fmt.Errorf("e.addDefinition: %w", err)
		}
		return nil
	}

	extractedRow := e.extractDefinitions(splitedRow)

	for _, word := range extractedRow {
		operation, exists := e.operations[word]
		if !exists {
			number, err := strconv.Atoi(word)
			if err != nil {
				return fmt.Errorf("strconv.Atoi: %w", err)
			}
			e.putOnStack(number)
			continue
		}

		e.stack, err = operation(e.stack)
		if err != nil {
			return fmt.Errorf("error executing operation [%s]: %w", word, err)
		}
	}

	return nil
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	err := e.processRow(row)
	if err != nil {
		return nil, err
	}

	return e.stack, nil
}
