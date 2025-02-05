package main

//Імпорт необхідних нам пакетів
import (
	"fmt"
	"html/template"
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
            width: 200px;
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
            <input type="text" name="vyhol" placeholder="Вугілля, т" required>
            <input type="text" name="mazyt" placeholder="Мазут, т" required>
            <input type="text" name="gaz" placeholder="Природний газ, м3" required>
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
		vyhol, err1 := strconv.ParseFloat(r.FormValue("vyhol"), 64)
		mazyt, err2 := strconv.ParseFloat(r.FormValue("mazyt"), 64)
		gaz, err3 := strconv.ParseFloat(r.FormValue("gaz"), 64)

		errors := []error{err1, err2, err3}

		for _, err := range errors {
			if err != nil {
				http.Error(w, "Введіть правильні числа", http.StatusBadRequest)
				return
			}
		}

		// Параметри для вугілля
		dataForVyhol := []float64{20.47, 0.8, 25.20, 1.5, 0.985, 0.0}
		// Параметри для мазуту
		dataForMazyt := []float64{39.48, 1.0, 0.15, 0.0, 0.985, 0.0}
		// Параметри для газу (в даному випадку всі нулі)
		//dataForGaz := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

		// Розрахунки для вугілля
		result := fmt.Sprintf("Викиди для вугілля: %.2f т\n", vyhol)
		ktvVyhol := findEmisia(dataForVyhol)
		result += fmt.Sprintf("KTV: %.2f г/ГДЖ\n", ktvVyhol)
		EtvVyhol := findVikid(ktvVyhol, dataForVyhol[0], vyhol)
		result += fmt.Sprintf("Etv: %.2f т\n\n", EtvVyhol)

		// Розрахунки для мазуту
		result += fmt.Sprintf("Викиди для мазуту: %.2f т\n", mazyt)
		ktvMazyt := findEmisia(dataForMazyt)
		result += fmt.Sprintf("KTV: %.2f г/ГДЖ\n", ktvMazyt)
		EtvMazyt := findVikid(ktvMazyt, dataForMazyt[0], mazyt)
		result += fmt.Sprintf("Etv: %.2f т\n\n", EtvMazyt)

		// Розрахунки для газу
		result += fmt.Sprintf("Викиди для газу: %.2f метрів кубічних\n", gaz)
		result += fmt.Sprintln("KTV: 0 г/ГДЖ") // Значення KTV для газу
		result += fmt.Sprintln("Etv: 0 т\n")   // Викиди для газу

		tmpl.Execute(w, result)
		return
	}

	tmpl.Execute(w, nil)
}

// Функція для розрахунку емісії твердих частинок
func findEmisia(data []float64) float64 {
	if len(data) != 6 {
		panic("Недостатньо параметрів.")
	}

	Qri := data[0]  // Нижня теплота згоряння
	avin := data[1] // Ефективність згоряння
	Ar := data[2]   // Масова частка золи
	Gvin := data[3] // Процент втрат тепла
	nzy := data[4]  // Ефективність очистки
	ktvs := data[5] // Специфічний показник

	// Обчислення результату з округленням
	result := (1000000/Qri)*avin*(Ar/(100-Gvin))*(1-nzy) + ktvs
	return result
}

// Функція для розрахунку викидів
func findVikid(ktv, Qri, B float64) float64 {
	// Обчислення викидів з округленням
	result := ktv * Qri * B / 1000000
	return result
}

func main() {
	http.HandleFunc("/", calcHandler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
