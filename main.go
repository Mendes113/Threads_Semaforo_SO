package main

import (
	"fmt"
	"sync"
)

type conta struct {
	saldo float64
	nome  string
}

var wg sync.WaitGroup
var semaforo = make(chan struct{}, 3) // Semáforo com capacidade 3

func (c *conta) saque(valor float64) {
	semaforo <- struct{}{}
	defer func() {
		<-semaforo
	}()

	if c.saldo >= valor {
		c.saldo -= valor
		fmt.Printf("%s fez um saque de %.2f\n", c.nome, valor)
	} else {
		fmt.Printf("%s não tem saldo suficiente para o saque\n", c.nome)
	}
}

func (c *conta) transferencia(valor float64, destino *conta) {
	semaforo <- struct{}{}
	defer func() {
		<-semaforo
	}()

	if c.saldo >= valor {
		c.saldo -= valor
		destino.saldo += valor
		fmt.Printf("%s transferiu %.2f para %s\n", c.nome, valor, destino.nome)
	} else {
		fmt.Printf("%s não tem saldo suficiente para a transferência\n", c.nome)
	}
}

func (C *conta) printTabelaDeSaldos() {
	fmt.Print("-------------------------------------------------------\n")
	fmt.Printf("Saldo da conta %s: %.2f\n", C.nome, C.consultaSaldo())
	fmt.Print("-------------------------------------------------------\n")
}

func (c *conta) consultaSaldo() float64 {
	return c.saldo
}

func main() {
	conta1 := conta{saldo: 1000, nome: "João"}
	conta2 := conta{saldo: 4000, nome: "Maria"}
	conta3 := conta{saldo: 5000, nome: "José"}
	conta4 := conta{saldo: 6000, nome: "Ana"}
	contas := []conta{conta1, conta2, conta3, conta4}

	// Print initial balances
	for _, c := range contas {
		c.printTabelaDeSaldos()
	}

	// Additional transactions
	for i := 1; i <= 5; i++ {
		finalWg := sync.WaitGroup{}
		finalWg.Add(1)
		go func() {
			c1 := conta{saldo: conta1.consultaSaldo(), nome: conta1.nome}
			c1.transferencia(500, &conta2)
			
		}()

		finalWg.Add(1)
		go func() {
			c3 := conta{saldo: conta3.consultaSaldo(), nome: conta3.nome}
			c3.transferencia(500, &conta4)
			finalWg.Done()
			}()

		finalWg.Add(1)
		go func() {
			c1 := conta{saldo: conta1.consultaSaldo(), nome: conta1.nome}
			c1.saque(100)
			finalWg.Done()
		}()

		finalWg.Add(1)
		go func() {
			c2 := conta{saldo: conta2.consultaSaldo(), nome: conta2.nome}
			c2.saque(200)
			finalWg.Done()
		}()
			
		finalWg.Wait() // Wait for this round of transactions to finish
		// Print final balances
		fmt.Println("---------------------------------------------------")
		fmt.Printf("Saldo final da conta %s: %.2f\n", conta1.nome, conta1.consultaSaldo())
		fmt.Printf("Saldo final da conta %s: %.2f\n", conta2.nome, conta2.consultaSaldo())
		fmt.Printf("Saldo final da conta %s: %.2f\n", conta3.nome, conta3.consultaSaldo())
		fmt.Printf("Saldo final da conta %s: %.2f\n", conta4.nome, conta4.consultaSaldo())
		

	}

	// // Print the final account balances
	// fmt.Printf("Saldo final da conta %s: %.2f\n", conta1.nome, conta1.consultaSaldo())
	// fmt.Printf("Saldo final da conta %s: %.2f\n", conta2.nome, conta2.consultaSaldo())

	
}
