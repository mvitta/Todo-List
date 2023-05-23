package car

import "fmt"

var _id int = 1

type Task struct {
	Description string
	IsComplete  bool
}

type Person struct {
	id        int
	firstName string
	lastName  string
	tasks     []Task
}

func NewPerson(_firstName string, _lastName string, tasks ...Task) *Person {
	p := &Person{
		id:        _id,
		firstName: _firstName,
		lastName:  _lastName,
		tasks:     make([]Task, 0),
	}

	p.tasks = append(p.tasks, tasks...)
	fmt.Println(&p)
	return p
}

func (p *Person) AddTask(newTask Task) {
	p.tasks = append(p.tasks, newTask)
}

func (p *Person) MarkCompleted(status bool, index int) {
	p.tasks[index].IsComplete = status
}

func (p *Person) CountTaskCompleted() {
	count := new(int)
	for i := 0; i < len(p.tasks); i++ {
		isComplete := p.tasks[i]
		if isComplete.IsComplete {
			*count++
		}
	}
	fmt.Println("Total completadas: ", *count)
}
