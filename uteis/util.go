package uteis

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//removendo arquivos
func RemoveFile(file string) (err error) {
	//Verifica se o arquivo existe, se existir apaga para gerar um novo
	if FileExist(file) {
		err = os.Remove(file)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

// Verifica se arquivos existe
func FileExist(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CheckNameFile(name []string) bool {
	for _, j := range name {
		if j != FCONTAS && j != FTRANSACOES {
			WriteLog("Nome dos arquivos não confere com o nome recebido: " + j)
			return false
		}
	}
	return true
}

//WriteLog Gera arquivo de log
func WriteLog(errostr string) {
	t := time.Now()

	errString := t.Format("Mon Jan _2 15:04:05 2006")
	errString = errString + " " + errostr + "\n"

	//
	ft, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	if _, err := ft.Write([]byte(errString)); err != nil {
		log.Println(err)
	}
}

//WriteLogClear Remove o arquivo de log.
func WriteLogClear() {
	RemoveFile(LOGFILE)
}

//FloatToString converte de float para string
func FloatToString(num float64, decimais int) string {
	//converte numero para string mantendo as casas decimais
	return strconv.FormatFloat(num, 'f', decimais, 64)
}

//FloatToStringClean converte de float para string
func FloatToStringClean(num float64, decimais int) string {
	//converte numero para string mantendo as casas decimais
	val := strconv.FormatFloat(num, 'f', decimais, 64)
	return strings.Replace(val, ".", "", -1)
}

//StringToSaldoComDecimal convertendo para float64 uma string
func StringToSaldoComDecimal(stringSaldoBody, stringSaldoDecimal string) (resultado float64) {
	resultado, err := strconv.ParseFloat(stringSaldoBody+"."+stringSaldoDecimal, 64)
	if err != nil {
		log.Println(err)
	}
	return
}

// fazendo um substr da string
func Substr(s string, pos, length int) string {

	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}

// convertendo string int para transformar em float e com casas decimais
// retorna uma string com o formato correto do valor que esta em string
func IdContaSaldoString(linha string) (idConta, SaldoFloatString, SaldoString string) {

	// linha
	// id conta
	// e saldo
	if linha != "" {

		// vamos transformar o valor em decimal
		// sera um float para que possamos
		// fazer os calculos
		vetorConta := strings.Split(linha, ",")

		// get conta
		idConta = vetorConta[0]

		// get saldo
		SaldoString := vetorConta[1]

		// colocar ponto nas duas ultimas posicoes
		// gerando casas decimais da string

		// somente o decimal sem as casas decimais
		stringSaldoBody := SaldoString[:len(SaldoString)-2]

		// somente o algarismo apos a virgula casas decimais
		stringSaldoDecimal := Substr(SaldoString, len(SaldoString)-2, 2)

		// gerando o saldo em float, com as casas decimais para efetuar os calculos
		Saldo := StringToSaldoComDecimal(stringSaldoBody, stringSaldoDecimal)

		// convertendo float para string,
		// float com duas casas decimais
		SaldoFloatString = FloatToString(Saldo, 2)

		// removendo negativo do registro saldo ou transacao
		valorClean := strings.Replace(SaldoString, "-", "", -1)

		//retornando dados da linha
		return idConta, SaldoFloatString, valorClean

	} else {

		return "", "", ""
	}
}

// convertendo para float64 uma string
func StringToFloat(valorString string) (Resultado float64) {

	Resultado, errs := strconv.ParseFloat(valorString, 64)

	if errs != nil {

		log.Println(errs)
	}

	return
}
