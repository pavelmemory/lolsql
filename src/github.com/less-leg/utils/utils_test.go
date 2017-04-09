package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
)

func Test1(t *testing.T) {
	t.Log("Suck my dick1")
	var err error
	os.Stdout, err = os.Create("log")
	os.Stderr = os.Stdout
	if err != nil {
		t.FailNow()
	}
	t.Log("Suck my dick1.1")
}

func Test2(t *testing.T) {
	t.Log("Suck my dick2")
}

func Test3(t *testing.T) {
	t.Log("Suck my dick3")
	os.Stdout.Sync()
	os.Stdout.Close()
}

type CommandManager struct {
	commands []Command
}

func (cm *CommandManager) Add(c Command) *CommandManager {
	cm.commands = append(cm.commands, c)
	return cm
}

func (cm *CommandManager) Run() []error {
	var errs []error
	for i := 0; i < len(cm.commands); i++ {
		if erri := cm.commands[i].Do(); erri != nil {
			errs = append(errs, erri)
			for j := i; j >= 0; j-- {
				command := cm.commands[j]
				if errj := command.Revert(); errj != nil {
					errs = append(errs, errj)
				}
				cm.commands = cm.commands[:j]
			}
			break
		}
	}
	return errs
}

type Command struct {
	Do     func() error
	Revert func() error
}

type CloudController struct {
}

func (cc *CloudController) ProvisionInstance(index int) error {
	fmt.Println("Provisioning", index)
	if index == 2 {
		return errors.New("Can't prvision")
	}
	fmt.Println("Provisioned", index)
	return nil
}

func (cc *CloudController) Provision(amount int) int {
	cm := &CommandManager{}
	successfully := 0
	for i := 0; i < amount; i++ {
		x := i
		cm.Add(Command{
			Do: func() error {
				err := cc.ProvisionInstance(x)
				if err == nil {
					successfully++
				}
				return err
			},
			Revert: func() error {
				fmt.Println("Reverted", x)
				successfully--
				return nil
			},
		})
	}
	errs := cm.Run()
	if len(errs) > 0 {
		successfully++
	}
	log.Println("end provison")
	return successfully
}

func Test4(t *testing.T) {
	cc := new(CloudController)
	expected := 0
	actual := cc.Provision(10)
	if actual != expected {
		t.Log("Provisioned only", actual, "of", expected)
		t.Fail()
	}
}
