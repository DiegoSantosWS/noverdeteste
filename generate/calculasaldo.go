package generate

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/DiegoSantosWS/gonoverde/bdb"
	"github.com/DiegoSantosWS/gonoverde/uteis"
)

var errs error

func SaldoContaCliente(ContasFile, TransFile string) {
	// apagando a base
	// de dados para
	// gerar uma nova
	bdb.DropDatabase()

	// Testing boltdb database
	// Start ping database
	// Creating ping ok
	bdb.Save("Ping", "ok")

	// Testing whether it was recorded
	// and read on the boltdb, we
	// recorded a Ping and then
	// read it back.
	if bdb.Get("Ping") != "ok" {

		log.Println("Services Error Data Base!")
		os.Exit(0)
	}

	errs = LerArquivoSaveDb(ContasFile)

	if errs != nil {

		fmt.Println("Error ao ler aquivo " + ContasFile + " não poderemos continuar!")
		log.Println(errs)
		os.Exit(0)
	}

	uteis.WriteLogClear()

	errs = CalcularSaldoTransacoes(TransFile)

	if errs != nil {

		fmt.Println("Error ao ler aquivo " + TransFile + ", isto impossibilita de fazer os calculos de saldo!")
		log.Println(errs)
		os.Exit(0)
	}
}

func LerArquivoSaveDb(ContasFile string) error {

	// variaveis declaradas para evitar declaracoes dentro do loop
	var linha, idConta, SaldoInicialFloatString, SaldoInicialString string

	// Abre o arquivo
	arquivo, err := os.Open(ContasFile)
	if err != nil {
		return err
	}
	defer arquivo.Close()
	scanner := bufio.NewScanner(arquivo)
	// varrendo o arquivo
	for scanner.Scan() {
		linha = scanner.Text()
		if linha != "" {
			quantVirg := strings.Count(linha, ",")
			if quantVirg == 1 {
				// retornando dados da linha em string
				idConta, SaldoInicialFloatString, SaldoInicialString = uteis.IdContaSaldoString(linha)
				if idConta != "" && SaldoInicialString != "" {
					re, _ := regexp.Compile(`[^0-9]`)
					if !re.MatchString(idConta) && !re.MatchString(SaldoInicialString) {
						// salvar nova banco idConta => Saldo
						bdb.Save(idConta, SaldoInicialFloatString)
					} else {
						// gerar log de erro
						uteis.WriteLog("O Arquivo " + ContasFile + ", foi encontrado o Idconta ou Saldo errados => idConta: " + idConta + " Saldo: " + SaldoInicialString)
					}
				} else {
					//gera log
					uteis.WriteLog("O Arquivo " + ContasFile + ", foi encontrado Idconta ou Saldo vazios => idConta: " + idConta + " Saldo: " + SaldoInicialString)
				}
			} else {
				uteis.WriteLog("O Arquivo " + ContasFile + " contém varias virgulas, isto não é permitido!")
			}
		} else {
			// gerar log
			uteis.WriteLog("O Arquivo " + ContasFile + " contém linha vazia!")
		}
	}
	return scanner.Err()
}

func CalcularSaldoTransacoes(TransFile string) error {
	var linha, idConta, ValorTransacaoStr, idContaTemp, SaldoIString, ValorTransacaoNotFloat string

	var SaldoInicialFloat, SaldoFloatTotal, ValorTransacaoFloat float64

	var VetorTransacao []float64

	var j int

	// set inicio
	idContaTemp = ""
	j = 0
	SaldoFloatTotal = 0

	// Abre o arquivo
	arquivo, err := os.Open(TransFile)

	// Caso tenha encontrado algum erro ao tentar abrir o arquivo retorne o erro encontrado
	if err != nil {
		return err
	}

	// Garante que o arquivo sera fechado apos o uso
	defer arquivo.Close()

	// Cria um scanner que le cada linha do arquivo
	scanner := bufio.NewScanner(arquivo)

	for scanner.Scan() {

		linha = scanner.Text()

		if linha != "" {
			quantVirg := strings.Count(linha, ",")
			if quantVirg == 1 {
				// retornando IdConta e o valor da transacao
				idConta, ValorTransacaoStr, ValorTransacaoNotFloat = uteis.IdContaSaldoString(linha)
				if idConta != "" && ValorTransacaoNotFloat != "" {
					// validar o conteudo do arquivo
					re, _ := regexp.Compile(`[^0-9]`)
					if !re.MatchString(idConta) && !re.MatchString(ValorTransacaoNotFloat) {
						ValorTransacaoFloat = uteis.StringToFloat(ValorTransacaoStr)
						SaldoIString = bdb.Get(idConta)
						if SaldoIString == "" {
							// mensagem de erro caso nao encontre o id da conta para pegar o saldo inicial
							textError := "O Arquivo [" + TransFile + "] não foi encontrado o saldo da conta [" + idConta + "] não foi encontrado"
							// gerar log de erro
							uteis.WriteLog(textError)
							continue
						}
						if j == 0 {

							// coloca no vetor
							VetorTransacao = append(VetorTransacao, ValorTransacaoFloat)

							// seta e
							// nao ira entrar
							// mais nesta condicao
							j = 1

						} else {
							if idContaTemp == idConta {
								VetorTransacao = append(VetorTransacao, ValorTransacaoFloat)
							} else {
								SaldoInicialFloat = uteis.StringToFloat(SaldoIString)

								SaldoFloatTotal = SaldoFloatTotal + SaldoInicialFloat

								CalculaSaldoBalanco(idContaTemp, VetorTransacao, SaldoFloatTotal)
								VetorTransacao = []float64{}
								SaldoFloatTotal = 0
								VetorTransacao = append(VetorTransacao, ValorTransacaoFloat)
							}
						}
						// pegar o idConta
						idContaTemp = idConta
					} else {
						uteis.WriteLog("O Arquivo [" + TransFile + "], foi encontrado o Idconta ou Transacao errados => idConta: [" + idConta + "] Transacao: " + ValorTransacaoNotFloat)
					}
				} else {

					uteis.WriteLog("O Arquivo [" + TransFile + "] não conseguimos ler o id Conta e o Saldo !")
				}
			} else { // varias virgulas ou uma
				uteis.WriteLog("O Arquivo [" + TransFile + "] contém varias virgulas isto não é permitido!")
			}
		} else { // linha vazia
			uteis.WriteLog("O Arquivo [" + TransFile + "] contém linha vazia!")
		}

	} // quando ele quebrar o laco precisará fazer o ultimo registro

	// fazendo a ultima posicao do vetor
	if len(VetorTransacao) > 0 {
		SaldoIString = bdb.Get(idContaTemp)

		if SaldoIString == "" {
			// mensagem de erro caso nao encontre o id da conta para pegar o saldo inicial
			textError := "O Arquivo [" + TransFile + "] não foi encontrado o saldo da conta [" + idContaTemp + "] não foi encontrado no banco de dados!"
			// err := errors.New(textError)

			// gerar log de erro
			uteis.WriteLog(textError)
		} else {
			// trasforma string em float do saldo
			SaldoInicialFloat = uteis.StringToFloat(SaldoIString)

			// saldo total
			SaldoFloatTotal = SaldoFloatTotal + SaldoInicialFloat

			// fazendo o calculo para escrever na tela
			CalculaSaldoBalanco(idContaTemp, VetorTransacao, SaldoFloatTotal)
		}
	}

	return scanner.Err()
}

func LerArquivoTransacao(TransFile string) error {
	var linha, idconta, SaldoInicialFloatString, SaldoInicialString string

	arquivo, err := os.Open(TransFile)
	if err != nil {
		return err
	}

	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		linha = scanner.Text()

		if linha != "" {
			idconta, SaldoInicialFloatString, SaldoInicialString = uteis.IdContaSaldoString(linha)
			if idconta != "" && SaldoInicialString != "" {
				re, _ := regexp.Compile(`[^0-9]`)
				if !re.MatchString(idconta) && !re.MatchString(SaldoInicialString) {
					keyT := "trans_" + idconta
					// perguntar se existe primeiro
					stringValores := bdb.Get(keyT, bdb.BDTrans)
					stringValores = stringValores + ";" + SaldoInicialFloatString
					if stringValores == "" {
						// salvar nova banco idConta => Saldo
						bdb.Save(keyT, SaldoInicialFloatString, bdb.BDTrans)
					} else {
						// concatenando os values para o id correspondente..
						stringValores = stringValores + ";" + SaldoInicialFloatString
						// salvando no banco
						bdb.Save(keyT, stringValores, bdb.BDTrans)
					}
				} else {
					uteis.WriteLog("O aquivo " + TransFile + " foi encontrado o idconta ou saldos errados -> idconta:" + idconta + " Saldo:" + SaldoInicialString)
				}
			} else {
				uteis.WriteLog("O aquivo " + TransFile + " foi encontrado o idconta ou saldos vazios -> idconta:" + idconta + " Saldo:" + SaldoInicialString)
			}
		} else {
			uteis.WriteLog("O aquivo " + TransFile + " contem linas em branco")
		}
	}
	return scanner.Err()
}

func CalculaSaldoBalanco(idContaTemp string, ArrayTranscao []float64, SaldoFloatTotal float64) {
	for _, Tvalor := range ArrayTranscao {
		SaldoFloatTotal = SaldoFloatTotal + Tvalor
		if SaldoFloatTotal < 0 && Tvalor < 0 {
			SaldoFloatTotal = SaldoFloatTotal - 5
		}

	}

	fmt.Println(idContaTemp, ", ", uteis.FloatToStringClean(SaldoFloatTotal, 2))
}
