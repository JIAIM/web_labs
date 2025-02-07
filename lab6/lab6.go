package main

//Імпорт необхідних нам пакетів
import (
	"fmt"
	"html/template"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

// Простий шаблон сторінки калькулятора
var tmpl = template.Must(template.New("calc").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Калькулятор</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 50px;
        }

        input {
            padding: 10px;
			margin-bottom: 5px;
            width: 200px;
            border: 1px solid #ccc;
            border-radius: 5px;
            text-align: center;
        }
		
		div {
			display: flex;
            flex-direction: column;
		}

        button {
            padding: 10px;
            margin-top: 15px;
            width: 200px; 
            border: none;
            background: #28a745;
            border-radius: 5px;
            cursor: pointer;
        }
        button:hover {
            background: #218838;
        }
        h3 {
            color: #333;
        }
    </style>
</head>
<body>
    <h2>Калькулятор</h2>
    <form method="post">
        <div>
            <input type="text" name="Pn" placeholder="Pn" required>
            <input type="text" name="Kv" placeholder="Kv" required>
			<input type="text" name="tg(phi)" placeholder="tg(phi)" required>
        </div>
        <button type="submit">Розрахувати</button>
    </form>
    {{if .}}
        <pre>{{.}}</pre>
    {{end}}
</body>
</html>
`))

// Основна ф-ія програми
func calcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Задаємо значення для різних видів пального
		Pn, err1 := strconv.Atoi(r.FormValue("Pn"))
		Kv, err2 := strconv.ParseFloat(r.FormValue("Kv"), 64)
		tgphi, err3 := strconv.ParseFloat(r.FormValue("tg(phi)"), 64)

		errors := []error{err1, err2, err3}

		for _, err := range errors {
			if err != nil {
				http.Error(w, "Введіть правильні числа", http.StatusBadRequest)
				return
			}
		}

		result := table(Pn, Kv, tgphi)

		tmpl.Execute(w, result)
		return
	}

	tmpl.Execute(w, nil)
}

// Функція для оцінки струму першого рівня
func estimatedCurrents1level(data int) float64 {
	cosphin := 0.9
	etan := 0.92
	Un := 0.38
	return float64(data) / (math.Sqrt(3) * Un * cosphin * etan)
}

// Функція для генерації таблиці розрахунків
func table(Pn int, Kv float64, tgphi float64) string {
	Un := 0.38
	nData := []int{4, 2, 4, 1, 1, 1, 2, 1}
	PnData := []int{Pn, 14, 42, 36, 20, 40, 32, 20}
	KvData := []float64{0.15, 0.12, 0.15, 0.3, 0.5, Kv, 0.2, 0.65}
	tgphiData := []float64{1.33, 1.0, 1.33, tgphi, 0.75, 1.0, 1.0, 0.75}

	nPnData := make([]int, len(nData))
	nPnKvData := make([]float64, len(nData))
	nPnKvTgphiData := make([]float64, len(nData))
	nPnPnData := make([]int, len(nData))
	IpData := make([]float64, len(nData))

	for i := range nData {
		nPnData[i] = nData[i] * PnData[i]
		nPnKvData[i] = float64(nPnData[i]) * KvData[i]
		nPnKvTgphiData[i] = nPnKvData[i] * tgphiData[i]
		nPnPnData[i] = nData[i] * PnData[i] * PnData[i]
		IpData[i] = estimatedCurrents1level(nPnData[i])
	}

	SHR1_n := sum(nData)
	SHR1_n_Pn := sum(nPnData)
	SHR1_Kv := new(big.Float).Quo(new(big.Float).SetFloat64(sumFloat(nPnKvData)), new(big.Float).SetFloat64(float64(SHR1_n_Pn)))
	SHR1_n_Pn_Kv := sumFloat(nPnKvData)
	SHR1_n_Pn_Kv_tgphi := sumFloat(nPnKvTgphiData)
	SHR1_n_Pn_Pn := sum(nPnPnData)
	SHR1_ne := (SHR1_n_Pn*SHR1_n_Pn)/SHR1_n_Pn_Pn + 1
	SHR1_Kp := 1.25
	SHR1_Pp := SHR1_Kp * SHR1_n_Pn_Kv
	SHR1_Qp := SHR1_n_Pn_Kv_tgphi
	SHR1_Sp := math.Sqrt(SHR1_Pp*SHR1_Pp + SHR1_Qp*SHR1_Qp)
	SHR1_Ip := SHR1_Pp / Un

	transfor_n, transfor_Pn := 2, 100
	transfor_n_Pn := transfor_n * transfor_Pn
	transfor_Kv := 0.2
	transfor_tgphi := 3.0
	transfor_n_Pn_Kv := float64(transfor_n_Pn) * transfor_Kv
	transfor_n_Pn_Kv_tgphi := transfor_n_Pn_Kv * transfor_tgphi
	transfor_n_Pn_Pn := transfor_n_Pn * transfor_Pn
	transfor_Ip := estimatedCurrents1level(transfor_n_Pn)

	sushi_n, sushi_Pn := 2, 120
	sushi_n_Pn := sushi_n * sushi_Pn
	sushi_Kv := 0.8
	sushi_n_Pn_Kv := float64(sushi_n_Pn) * sushi_Kv
	sushi_n_Pn_Pn := sushi_n_Pn * sushi_Pn
	sushi_Ip := estimatedCurrents1level(sushi_n_Pn)

	all_n, all_n_Pn := 81, 2330
	all_n_Pv_Kv := 752
	all_Kv := new(big.Float).Quo(new(big.Float).SetFloat64(float64(all_n_Pv_Kv)), new(big.Float).SetFloat64(float64(all_n_Pn)))
	all_n_Pv_Kv_tgphi := 657
	all_n_Pn_Pn := 96399
	all_ne := all_n_Pn * all_n_Pn / all_n_Pn_Pn
	all_Kp := 0.7
	all_Pp := all_Kp * float64(all_n_Pv_Kv)
	all_Qp := all_Kp * float64(all_n_Pv_Kv_tgphi)
	all_Sp := math.Sqrt(all_Pp*all_Pp + all_Qp*all_Qp)
	all_Ip := all_Pp / Un

	result := ""

	result += fmt.Sprintln(strings.Repeat("-", 127))

	result += fmt.Sprintln("|name      |etan|cosu| Un  | n | Pn | n*Pn | Kv  | tgφ | n*Pn*Kv |n*Pn*Kv*tgφ| n*Pn*Pn| Ne |  Kp |  Pp  |  Qp  |   Sp  |  Ip  |")

	result += fmt.Sprintln(strings.Repeat("-", 127))

	for i := range nData {
		result += fmt.Sprintf("|SHR1      |%.2f|%.1f |%.2f |%2d |%3d |%5d |%.2f |%.2f |%8.2f |%10.3f |%7d |    |     |      |      |       |%5.1f |\n",
			0.92, 0.9, Un, nData[i], PnData[i], nPnData[i], KvData[i], tgphiData[i],
			nPnKvData[i], nPnKvTgphiData[i], nPnPnData[i], IpData[i])
	}

	result += fmt.Sprintln(strings.Repeat("-", 127))

	for i := 1; i <= 3; i++ {
		str := fmt.Sprintf("SHR%d", i)
		result += fmt.Sprintf("|all %s  |    |    |     | %2d|    |%5d |%.2f |     |%8.2f |%10.3f |%7d | %d |%5.2f|%5.2f|%5.2f|%7.2f|%5.1f |\n",
			str, SHR1_n, SHR1_n_Pn, SHR1_Kv, SHR1_n_Pn_Kv, SHR1_n_Pn_Kv_tgphi, SHR1_n_Pn_Pn, SHR1_ne, SHR1_Kp, SHR1_Pp, SHR1_Qp, SHR1_Sp, SHR1_Ip)
	}
	result += fmt.Sprintln(strings.Repeat("-", 127))

	result += fmt.Sprintf("|transforma|%.2f|%.1f |%.2f |%2d |%3d |%5d |%.2f |%.2f |%8.2f |%10.3f |%7d |    |     |      |      |       |%5.1f |\n",
		0.92, 0.9, Un, transfor_n, transfor_Pn, transfor_n_Pn, transfor_Kv, transfor_tgphi, transfor_n_Pn_Kv, transfor_n_Pn_Kv_tgphi, transfor_n_Pn_Pn, transfor_Ip)

	result += fmt.Sprintf("|shafa     |%.2f|%.1f |%.2f |%2d |%3d |%5d |%.2f |     |%8.2f |           |%7d |    |     |      |      |       |%5.1f |\n",
		0.92, 0.9, Un, sushi_n, sushi_Pn, sushi_n_Pn, sushi_Kv, sushi_n_Pn_Kv, sushi_n_Pn_Pn, sushi_Ip)

	result += fmt.Sprintln(strings.Repeat("-", 127))

	result += fmt.Sprintf("|All       |    |    |     | %2d|    |%5d |%.2f |     |%8d |%10d |%7d | %d |%5.2f|%5.2f|%5.2f|%7.2f|%5.1f|\n",
		all_n, all_n_Pn, all_Kv, all_n_Pv_Kv, all_n_Pv_Kv_tgphi, all_n_Pn_Pn, all_ne, all_Kp, all_Pp, all_Qp, all_Sp, all_Ip)

	result += fmt.Sprintln(strings.Repeat("-", 127))
	return result
}

func sum(arr []int) int {
	sum := 0
	for _, val := range arr {
		sum += val
	}
	return sum
}

func sumFloat(arr []float64) float64 {
	sum := 0.0
	for _, val := range arr {
		sum += val
	}
	return sum
}

func main() {
	http.HandleFunc("/", calcHandler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
