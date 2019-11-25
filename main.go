package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"

	"github.com/tarm/serial"
)

var port *serial.Port

const PocsagMaxMsgLen = 78

func main() {
	initCmd()
}

func openPort() {
	c := &serial.Config{
		Name:     "/dev/ttyUSB0",
		Baud:     9600,
		Size:     serial.DefaultSize,
		Parity:   serial.ParityNone,
		StopBits: serial.Stop1,
	}

	var err error
	port, err = serial.OpenPort(c)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func closePort() {
	if err := port.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCmd() {
	app := &cli.App{
		Action: mainActionCmd,
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(err)
	}
}

func mainActionCmd(c *cli.Context) error {
	if c.NArg() < 2 {
		return fmt.Errorf("Pocos argumentos | Uso: beeper <capcode> <mensaje>")
	}

	capcode, err := strconv.Atoi(c.Args().Get(0))

	if err != nil {
		return fmt.Errorf("El capcode debe ser un numero")
	}

	msgcmd := c.Args().Slice()[1:]
	msg := strings.Join(msgcmd, " ")

	psmsg, err := createPocsagMessage(capcode, msg)

	if err != nil {
		return err
	}

	err = sendMsgToPort(psmsg)

	if err != nil {
		return err
	}

	fmt.Println("Mensaje enviado al biper: \"" + msg + "\"")

	return nil
}

func createPocsagMessage(capcode int, msg string) (string, error) {
	strcapcode := strconv.Itoa(capcode)

	switch {
	case len(msg) > PocsagMaxMsgLen:
		return "", fmt.Errorf("El mensaje debe ser menor a 78 caracteres")
	case len(strcapcode) > 7:
		return "", fmt.Errorf("El capcode debe ser igual o menor a 7 digitos")
	}

	for len(strcapcode) < 7 {
		strcapcode = "0" + strcapcode
	}

	pocsagmsg := "P," + strcapcode + ",0,512,N,A," + strings.ToUpper(msg)

	return pocsagmsg, nil
}

func sendMsgToPort(msg string) error {
	openPort()

	_, err := port.Write([]byte(msg))

	if err != nil {
		return err
	}

	closePort()

	return nil
}
