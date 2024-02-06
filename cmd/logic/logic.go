package logic

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hacker301et/wifipen/utils"
)

type CrackWifies struct {
	P     *tea.Program
	Iface string
}

func Init(iface string) *CrackWifies {
	return &CrackWifies{
		Iface: iface,
	}
}
func (c *CrackWifies) checkForRequiredTools(m *utils.Model) error {
	var err error
	tools := []string{"aircrack-ng", "ifconfig"}
	for _, tool := range tools {
		cmd := exec.Command("which", tool)
		if nerr := cmd.Run(); nerr != nil {
			err = nerr
			m.Sub <- utils.Row{tool, "missing ⭕"}
			time.Sleep(time.Second)
			continue
		}
		m.Sub <- utils.Row{tool, "present ✅"}
		time.Sleep(time.Second)
	}

	return err
}
func (c *CrackWifies) clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func (c *CrackWifies) setupWifiInterface() error {
	cmd := exec.Command("ifconfig", c.Iface)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to find inteface named %s", c.Iface)
	}
	cmd = exec.Command("ifconfig", c.Iface, "down")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to access inteface  down command for %s make sure to run the tool as root (sudo)", c.Iface)
	}
	time.Sleep(time.Second * 2)

	cmd = exec.Command("airmon-ng", "start", c.Iface)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to change inteface %s to monitor mode ( inteface may not support monitor mode)", c.Iface)
	}
	c.Iface = c.Iface + "mon"
	time.Sleep(time.Second * 2)
	cmd = exec.Command("airmon-ng", "check", "kill")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to prepare inteface %s", c.Iface)
	}
	time.Sleep(time.Second * 2)
	cmd = exec.Command("ifconfig", c.Iface, "up")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to access inteface  up command for %s ", c.Iface)
	}

	return nil
}
func (c *CrackWifies) MonitorWifies() error {
	file, err := os.OpenFile("store/wifi.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		c.P.Quit()
		return err
	}
	defer file.Close()

	cmd := exec.Command("airodump-ng", c.Iface)
	cmd.Stdout = file
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return err
	}

	// err = cmd.Wait()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (c *CrackWifies) Run(m *utils.Model, p *tea.Program) {
	if err := c.checkForRequiredTools(m); err != nil {
		p.Quit()
		return
	}
	c.clearScreen()
	fmt.Println("preparting  your interface...")
	if err := c.setupWifiInterface(); err != nil {
		fmt.Println(err)
		p.Quit()

	}
	c.clearScreen()
	c.MonitorWifies()

}

func (c *CrackWifies) Start() {
	m := utils.NewView([]table.Column{
		{Title: "tool", Width: 100},
		{Title: "status", Width: 10},
	})
	p := tea.NewProgram(m)

	go c.Run(m, p)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}

}
