package main

//Імпорт необхідних нам пакетів
import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
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
            margin: 5px;
            width: 80px;
            border: 1px solid #ccc;
            border-radius: 5px;
            text-align: center;
        }
        button {
            padding: 10px;
            margin: 5px;
            width: 110px;  
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
            <h3>Task 1</h3>
            <input type="text" name="Hp" placeholder="Введіть Hp" required>
            <input type="text" name="Cp" placeholder="Введіть Cp" required>
            <input type="text" name="Sp" placeholder="Введіть Sp" required>
            <input type="text" name="Np" placeholder="Введіть Np" required>
            <input type="text" name="Op" placeholder="Введіть Op" required>
            <input type="text" name="W" placeholder="Введіть W" required>
            <input type="text" name="Ap" placeholder="Введіть Ap" required>
        </div>
        <div>
            <h3>Task 2</h3>
            <input type="text" name="Hg" placeholder="Введіть Hg" required>
            <input type="text" name="Cg" placeholder="Введіть Cg" required>
            <input type="text" name="Og" placeholder="Введіть Og" required>
            <input type="text" name="Sg" placeholder="Введіть Sg" required>
            <input type="text" name="QgH" placeholder="Введіть QgH" required>
            <input type="text" name="W2" placeholder="Введіть W" required>
            <input type="text" name="A" placeholder="Введіть A" required>
            <input type="text" name="Vg" placeholder="Введіть Vg" required>
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
		// Дані про компоненти палива
		Hp, err1 := strconv.ParseFloat(r.FormValue("Hp"), 64)
		Cp, err2 := strconv.ParseFloat(r.FormValue("Cp"), 64)
		Sp, err3 := strconv.ParseFloat(r.FormValue("Sp"), 64)
		Np, err4 := strconv.ParseFloat(r.FormValue("Np"), 64)
		Op, err5 := strconv.ParseFloat(r.FormValue("Op"), 64)
		W, err6 := strconv.ParseFloat(r.FormValue("W"), 64)
		Ap, err7 := strconv.ParseFloat(r.FormValue("Ap"), 64)
		// Дані про склад мазуту
		Hg, err8 := strconv.ParseFloat(r.FormValue("Hg"), 64)
		Cg, err9 := strconv.ParseFloat(r.FormValue("Cg"), 64)
		Og, err10 := strconv.ParseFloat(r.FormValue("Og"), 64)
		Sg, err11 := strconv.ParseFloat(r.FormValue("Sg"), 64)
		QgH2, err12 := strconv.ParseFloat(r.FormValue("QgH"), 64)
		W2, err13 := strconv.ParseFloat(r.FormValue("W2"), 64)
		A, err14 := strconv.ParseFloat(r.FormValue("A"), 64)
		Vg, err15 := strconv.ParseFloat(r.FormValue("Vg"), 64)
		// Обробка помилок
		errors := []error{err1, err2, err3, err4, err5, err6, err7, err8, err9, err10, err11, err12, err13, err14, err15}

		for _, err := range errors {
			if err != nil {
				http.Error(w, "Введіть правильні числа", http.StatusBadRequest)
				return
			}
		}

		// Розрахунок коефіцієнтів переходу
		Kpc := 100 / (100 - W)
		Kpg := 100 / (100 - W - Ap)

		// Перевірка коефіцієнтів
		Hc, Cc, Sc, Nc, Oc, Ac := Hp*Kpc, Cp*Kpc, Sp*Kpc, Np*Kpc, Op*Kpc, Ap*Kpc
		if math.Round(Hc+Cc+Sc+Nc+Oc+Ac) != 100 {
			panic("Сума повинна дорівнювати 100%. Коефіцієнт переходу від робочої до сухої маси обраховано неправильно")
		}

		Hg, Cg, Sg, Ng, Og := Hp*Kpg, Cp*Kpg, Sp*Kpg, Np*Kpg, Op*Kpg
		if math.Round(Hg+Cg+Sg+Ng+Og) != 100 {
			panic("Сума повинна дорівнювати 100%. Коефіцієнт переходу від робочої до горючої маси обраховано неправильно")
		}

		// Розрахунок нижчої теплоти згоряння
		QpH := math.Round(339*Cp+1030*Hp-108.8*(Op-Sp)-25*W) / 1000 //Округлюємо значення
		QcH := (QpH + 0.025*W) * 100 / (100 - W)
		QgH := (QpH + 0.025*W) * 100 / (100 - W - Ap)

		// Обчислення коефіцієнта
		kef := (100 - W2 - A) / 100

		// Вирахування нових даних для елементарного складу
		Hp2, Cp2, Sp2, Op2, Ap2, Vp2 := Hg*kef, Cg*kef, Sg*kef, Og*kef, A*kef, Vg*kef

		// Перерахунок нижчої теплоти згоряння мазуту
		QpH2 := QgH2*kef - 0.025*W2

		// Вивід результатів
		result := fmt.Sprintf("Task #1\nQpH: %.2f МДж, QcH: %.2f МДж, QgH: %.2f МДж\n\n", QpH, QcH, QgH)
		result += fmt.Sprintf("Task #2\nQgh: %.2f МДж, QpH: %.2f МДж\n", QgH2, QpH2)
		result += fmt.Sprintf("Елементарний склад:\nВуглець: %.2f\nВодень: %.2f\nКисень: %.2f\nСірка: %.2f\nЗола: %.2f\nВанадій: %.2f\n",
			Hp2, Cp2, Op2, Sp2, Ap2, Vp2)

		tmpl.Execute(w, result)
		return
	}

	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", calcHandler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
