package harness

import "time"

type State struct {
	NextGoalID   int    `json:"next_goal_id"`
	NextTaskID   int    `json:"next_task_id"`
	ActiveGoalID int    `json:"active_goal_id"`
	Goals        []Goal `json:"goals"`
	Tasks        []Task `json:"tasks"`
	UpdatedAt    string `json:"updated_at"`
}

type Goal struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID        int        `json:"id"`
	GoalID    int        `json:"goal_id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	CreatedAt time.Time  `json:"created_at"`
	DoneAt    *time.Time `json:"done_at,omitempty"`
}

func (s *State) AddGoal(title string) Goal {
	goal := Goal{
		ID:        s.NextGoalID,
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}
	s.Goals = append(s.Goals, goal)
	s.NextGoalID++
	s.ActiveGoalID = goal.ID
	s.touch()
	return goal
}

func (s *State) SetActiveGoal(id int) bool {
	if _, ok := s.GoalByID(id); !ok {
		return false
	}
	s.ActiveGoalID = id
	s.touch()
	return true
}

func (s *State) AddTask(title string) (Task, bool) {
	return s.AddTaskToGoal(s.ActiveGoalID, title)
}

func (s *State) GoalByID(id int) (Goal, bool) {
	for _, goal := range s.Goals {
		if goal.ID == id {
			return goal, true
		}
	}
	return Goal{}, false
}

func (s *State) ActiveGoal() (Goal, bool) {
	return s.GoalByID(s.ActiveGoalID)
}

func (s *State) AddTaskToGoal(goalID int, title string) (Task, bool) {
	if _, ok := s.GoalByID(goalID); !ok {
		return Task{}, false
	}
	task := Task{
		ID:        s.NextTaskID,
		GoalID:    goalID,
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}
	s.Tasks = append(s.Tasks, task)
	s.NextTaskID++
	s.touch()
	return task, true
}

func (s *State) NextOpenTask() (Task, bool) {
	for _, task := range s.Tasks {
		if !task.Done && task.GoalID == s.ActiveGoalID {
			return task, true
		}
	}
	return Task{}, false
}

func (s *State) MarkDone(id int) bool {
	now := time.Now().UTC()
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			s.Tasks[i].Done = true
			s.Tasks[i].DoneAt = &now
			s.touch()
			return true
		}
	}
	return false
}

func (s *State) GoalsInListOrder() []Goal {
	goals := make([]Goal, 0, len(s.Goals))
	for _, goal := range s.Goals {
		if goal.ID != s.ActiveGoalID {
			goals = append(goals, goal)
		}
	}
	if active, ok := s.ActiveGoal(); ok {
		goals = append(goals, active)
	}
	return goals
}

func (s *State) TasksForGoal(goalID int) []Task {
	tasks := make([]Task, 0)
	for _, task := range s.Tasks {
		if task.GoalID == goalID {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (s *State) OpenTaskCount() int {
	count := 0
	for _, task := range s.Tasks {
		if !task.Done {
			count++
		}
	}
	return count
}

func (s *State) DoneTaskCount() int {
	count := 0
	for _, task := range s.Tasks {
		if task.Done {
			count++
		}
	}
	return count
}

func (s *State) touch() {
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}
