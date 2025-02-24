package application

import (
	"errors"
	"log"
)

type ApplicationStack struct {
	stack []ApplicationFrame
}

func NewApplicationStack() *ApplicationStack {
	return &ApplicationStack{
		stack: []ApplicationFrame{},
	}
}

func (af *ApplicationStack) Push(frame ApplicationFrame) error {
	log.Println("New frame pushed")
	af.stack = append(af.stack, frame)
	return nil
}

func (af *ApplicationStack) Pop() (ApplicationFrame, error) {
	length := len(af.stack)

	// This is the last frame, do not throw an error
	// Failsafe is Engineer accidentally allows this to happen
	if length <= 1 {
		return af.stack[0], nil
	}

	frame := af.stack[length-1]
	af.stack = af.stack[:length-1]

	return frame, nil
}

func (af *ApplicationStack) ClearStack() error {
	af.stack = nil
	return nil
}

func (af *ApplicationStack) Peek() (ApplicationFrame, error) {
	length := len(af.stack)
	if length < 1 {
		return nil, errors.New("There are no frames to view, something went wrong in the stack module")
	}

	return af.stack[length-1], nil
}

func (af *ApplicationStack) Size() int {
	return len(af.stack)
}
