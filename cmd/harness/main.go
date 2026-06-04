package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go-towerdefense/internal/harness"
)

func main() {
	statePath := ".tdharness/state.json"
	store := harness.NewStore(statePath)
	state, err := store.Load()
	if err != nil {
		fatal(err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "goal":
		if len(args) < 2 {
			fatal(fmt.Errorf("usage: goal <text>"))
		}
		goal := state.AddGoal(strings.Join(args[1:], " "))
		must(store.Save(state))
		fmt.Printf("goal created and activated: #%d %s\n", goal.ID, goal.Title)
	case "switch":
		if len(args) != 2 {
			fatal(fmt.Errorf("usage: switch <goal-id>"))
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fatal(fmt.Errorf("invalid goal id: %w", err))
		}
		if !state.SetActiveGoal(id) {
			fatal(fmt.Errorf("goal not found: %d", id))
		}
		must(store.Save(state))
		fmt.Printf("active goal set to #%d\n", id)
	case "add":
		if len(args) < 2 {
			fatal(fmt.Errorf("usage: add <task>"))
		}
		task, ok := state.AddTask(strings.Join(args[1:], " "))
		if !ok {
			fatal(fmt.Errorf("no active goal set; create or switch to a goal first"))
		}
		must(store.Save(state))
		fmt.Printf("task added to goal #%d: #%d %s\n", task.GoalID, task.ID, task.Title)
	case "list":
		printTasks(state)
	case "next":
		task, ok := state.NextOpenTask()
		if !ok {
			fmt.Println("no open tasks")
			return
		}
		fmt.Printf("#%d %s\n", task.ID, task.Title)
	case "done":
		if len(args) != 2 {
			fatal(fmt.Errorf("usage: done <id>"))
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fatal(fmt.Errorf("invalid task id: %w", err))
		}
		if !state.MarkDone(id) {
			fatal(fmt.Errorf("task not found: %d", id))
		}
		must(store.Save(state))
		fmt.Println("task completed")
	case "summary":
		printSummary(state)
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("basic harness")
	fmt.Println("commands:")
	fmt.Println("  goal <text>     create and activate a new goal")
	fmt.Println("  switch <id>     switch to an existing goal")
	fmt.Println("  add <task>      add a task to the active goal")
	fmt.Println("  list            list tasks grouped by goal")
	fmt.Println("  next            show the next open task for the active goal")
	fmt.Println("  done <id>       mark a task complete")
	fmt.Println("  summary         show active goal and task counts")
}

func printTasks(state *harness.State) {
	if len(state.Goals) == 0 {
		fmt.Println("no goals yet")
		return
	}
	for _, goal := range state.GoalsInListOrder() {
		active := ""
		if goal.ID == state.ActiveGoalID {
			active = " (active)"
		}
		fmt.Printf("goal #%d%s: %s\n", goal.ID, active, goal.Title)
		tasks := state.TasksForGoal(goal.ID)
		if len(tasks) == 0 {
			fmt.Println("  (no tasks)")
			continue
		}
		for _, task := range tasks {
			status := "open"
			if task.Done {
				status = "done"
			}
			fmt.Printf("  #%d [%s] %s\n", task.ID, status, task.Title)
		}
	}
}

func printSummary(state *harness.State) {
	activeGoal, ok := state.ActiveGoal()
	if ok {
		fmt.Printf("active goal: #%d %s\n", activeGoal.ID, activeGoal.Title)
	} else {
		fmt.Println("active goal: (unset)")
	}
	fmt.Printf("goals: %d\n", len(state.Goals))
	fmt.Printf("open tasks: %d\n", state.OpenTaskCount())
	fmt.Printf("done tasks: %d\n", state.DoneTaskCount())
}

func must(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
