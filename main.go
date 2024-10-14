package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"crypto.com/wallet"
)

func printHelp() {
	fmt.Println(
		"usage: \n" +
			fmt.Sprintf("\t%-40s\t%-30s\n", "deposits money into userA's wallet", "-> a 100.00") +
			fmt.Sprintf("\t%-40s\t%-30s\n", "withdraws money from userA's wallet", "a -> 100.00") +
			fmt.Sprintf("\t%-40s\t%-30s\n", "userA sends money to userB", "a -> 100.00 -> b") +
			fmt.Sprintf("\t%-40s\t%-30s\n", "check userA's balance", "a ?") +
			fmt.Sprintf("\t%-40s\t%-30s\n", "check userA's transaction history", "a ??") +
			fmt.Sprintf("\t%-40s\t%-30s\n", "exit", "exit"),
	)
}

func main() {
	flag.BoolVar(&OnlyClient, "client", false, "only run the client side")
	flag.Parse()

	if !OnlyClient {
		go func() {
			s.Handler = wallet.NewRouter()
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				if strings.Contains(err.Error(), "bind") { // port already in use
					OnlyClient = true
				} else {
					log.Fatalf("s.ListenAndServe err: %v", err)
				}
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // watch for SIGINT (Ctrl+C)
	reader := bufio.NewReader(os.Stdin)
	printHelp()

	for {
		select {
		case <-quit:
			exit()
		default:
			fmt.Print(">> ")

			input, err := reader.ReadString('\n')
			if err != nil && err == io.EOF {
				exit()
			}

			handleCommand(input)
		}
	}
}

func handleCommand(command string) {
	client := wallet.NewClient(ServerAddress)
	parts := strings.Fields(command)
	lenParts := len(parts)

	if lenParts == 0 {
		printHelp()
	} else if lenParts == 1 {
		switch parts[0] {
		case "help", "?":
			printHelp()
		case "exit":
			exit()
		default:
			fmt.Printf("Unknown command: %s\n", parts[0])
		}
	} else if lenParts == 2 && parts[1] == "?" {
		p := &wallet.ParamUser{UserName: parts[0]}
		if result, err := client.Balance(p); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("UserName: %s\nBalance: %s\n", result.UserName, wallet.Cent2String(result.Balance))
		}
	} else if lenParts == 2 && parts[1] == "??" {
		p := &wallet.ParamUser{UserName: parts[0]}
		if result, err := client.History(p); err != nil {
			fmt.Println(err)
		} else if len(result.List) == 0 {
			fmt.Println("no transaction")
		} else {
			fmt.Printf("%-20s\t%-20s\t%-10s\t%10s\t%10s\t%-15s\n", "Sender", "Receiver", "Action", "Money", "Balance", "Date Time")
			for i := len(result.List); i > 0; i-- {
				v := result.List[i-1]
				vTime := time.Unix(v.CreateDate, 0).Format("2006-01-02 15:04:05")
				fmt.Printf("%-20s\t%-20s\t%-10s\t%10s\t%10s\t%10s\n", v.Sender, v.Receiver, v.Action, wallet.Cent2String(v.Money), wallet.Cent2String(v.Balance), vTime)
			}
		}
	} else if lenParts == 3 && parts[0] == "->" {
		p := &wallet.ParamDeposit{}
		p.UserName = parts[1]
		p.Money = wallet.String2Cent(parts[2])
		if result, err := client.Deposit(p); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("UserName: %s\nBalance: %s\n", result.UserName, wallet.Cent2String(result.Balance))
		}
	} else if lenParts == 3 && parts[1] == "->" {
		p := &wallet.ParamWithdraw{}
		p.UserName = parts[0]
		p.Money = wallet.String2Cent(parts[2])
		if result, err := client.Withdraw(p); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("UserName: %s\nBalance: %s\n", result.UserName, wallet.Cent2String(result.Balance))
		}
	} else if lenParts == 5 && parts[1] == "->" && parts[3] == "->" {
		p := &wallet.ParamTransfer{}
		p.UserName = parts[0]
		p.Money = wallet.String2Cent(parts[2])
		p.Receiver = wallet.ParamUser{UserName: parts[4]}
		if result, err := client.Transfer(p); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Sender: %s\nBalance: %s\nReceiver: %s\nBalance: %s\n",
				result.Sender.UserName,
				wallet.Cent2String(result.Sender.Balance),
				result.Receiver.UserName,
				wallet.Cent2String(result.Receiver.Balance),
			)
		}
	} else {
		fmt.Println("Unknown command")
	}
}

func exit() {
	if !OnlyClient {
		s.Close() // shutdown the server
	}
	fmt.Printf("\n\nBye!\n\n")
	os.Exit(0)
}

var (
	ServerAddress = "localhost:8686"
	s             = &http.Server{
		Addr:           ServerAddress,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
)

var OnlyClient bool // only run the client side
