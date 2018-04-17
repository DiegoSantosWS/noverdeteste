package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DiegoSantosWS/gonoverde/generate"
	"github.com/DiegoSantosWS/gonoverde/uteis"
)

func main() {
	params := os.Args[1:]
	generate.SaldoContaCliente(params[0], params[1])
}

func generateFiles(params []string) {
	fmt.Println("...Gerando arquivos...")
	if uteis.CheckNameFile(params) {
		uteis.RemoveFile(uteis.FCONTAS)

		uteis.RemoveFile(uteis.FTRANSACOES)
		//se o arquivo de contas não existir cria e faz a o append no conteudo
		cs, err := os.OpenFile(params[0], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err.Error())
		}
		//se o arquivo de transacoes não existir cria e faz a o append no conteudo
		ts, err1 := os.OpenFile(params[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err1 != nil {
			log.Println(err.Error())
		}

		//dados da conta
		var idconta, k, qtdTransacoes, countConta, countTrans int

		k = 1
		countConta = 0
		countTrans = 0

		qtdTransacoes = Random(2, 100)

		idconta = 10

		var saldoi, vtransacao float64

		var idcontaNow, stringsaldoi, stringtransacao, stringValores string

		var replace = strings.NewReplacer(".", "")

		for i := 1; i <= uteis.LINHAS; i++ {
			//saldo conta
			randSaldo := Random(1, 100000)
			//saldo transacao
			randSaldoTr := Random(1, 100000)
			//saldo transacao negativa
			randSaldoTrNeg := Random(1, 800000)

			//gerando saldo aleatoriamente
			saldoi = RandomD() * float64(randSaldo)

			// positivo
			if randSaldoTr%2 == 0 {
				vtransacao = RandomD() * float64(randSaldoTr)
			} else {
				vtransacao = RandomD() * float64(randSaldoTrNeg) * -1.0
			}
			//conver float for string
			stringtransacao = strconv.FormatFloat(vtransacao, 'f', 2, 64)
			//remove o . (ponto) data string
			stringtransacao = replace.Replace(stringtransacao)
			//conver float for string
			stringsaldoi = strconv.FormatFloat(saldoi, 'f', 2, 64)
			//remove o . (ponto) data string
			stringsaldoi = replace.Replace(stringsaldoi)
			// converte inteiro para string
			idcontaNow = strconv.Itoa(idconta)

			if k == 1 {
				// concatena os dados para gerar o arquivo conforme o modelo solicitado
				stringValores = idcontaNow + "," + stringsaldoi + "\n"

				//gera uma conta e uma transacao somente uma vez salva as contas
				if _, err := cs.Write([]byte(stringValores)); err != nil {
					log.Println(err)
				}

				// concatena os dados para gerar o arquivo conforme o modelo solicitado
				stringValores = idcontaNow + "," + stringtransacao + "\n"

				//gera uma conta e uma transacao somente uma vez salva as contas
				if _, err := ts.Write([]byte(stringValores)); err != nil {
					log.Println(err)
				}
				countConta++
				countTrans++
			} else {
				if k == qtdTransacoes {
					stringValores = idcontaNow + "," + stringtransacao + "\n"
					//salvando a transacao
					if _, err := ts.Write([]byte(stringValores)); err != nil {
						log.Println(err)
					}

					qtdTransacoes = Random(2, 100)
					idconta++
					k = 0
					countTrans++
				} else {
					stringValores = idcontaNow + "," + stringtransacao + "\n"
					//salvando a transacao
					if _, err := ts.Write([]byte(stringValores)); err != nil {
						log.Println(err)
					}

					countTrans++
				}
			}

			k++
		}
		if err := cs.Close(); err != nil {
			log.Println(err)
		}

		fmt.Printf("\nArquivo %s gerado com sucesso, total de [%d] linhas\n", params[0], countConta)
		fmt.Printf("Arquivo %s gerado com sucesso, total de [%d] linhas\n", params[1], countTrans)
		fmt.Printf("Foi gerado total de [%d] linhas\n", uteis.LINHAS)
	}
}

//Random gera numero
func Random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//RandomD gera numero
func RandomD() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}
