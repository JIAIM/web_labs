// package main

// //Імпорт необхідних нам пакетів
// import (
// 	"fmt"
// 	"html/template"
// 	"math"
// 	"net/http"
// 	"strconv"
// )

// // Простий шаблон сторінки калькулятора
// var tmpl = template.Must(template.New("calc").Parse(`
// <!DOCTYPE html>
// <html lang="ru">
// <head>
//     <meta charset="UTF-8">
//     <title>Калькулятор</title>
//     <style>
//         body {
//             font-family: Arial, sans-serif;
//             background-color: #f4f4f4;
//             margin: 50px;
//         }

//         input {
//             padding: 10px;
//             width: 200px;
//             border: 1px solid #ccc;
//             border-radius: 5px;
//             text-align: center;
//         }

// 		div {
// 			display: flex;
//             flex-direction: column;
// 		}

//         button {
//             padding: 10px;
//             margin-top: 15px;
//             width: 110px;
//             border: none;
//             background: #28a745;
//             border-radius: 5px;
//             cursor: pointer;
//         }
//         button:hover {
//             background: #218838;
//         }
//         h3 {
//             color: #333;
//         }
//     </style>
// </head>
// <body>
//     <h2>Калькулятор</h2>
//     <form method="post">
//         <div>
// 			<h3>Середньодобова потужність Pc(MBт)</h3>
//             <input type="text" name="pc" placeholder="" required>
// 			<h3>Сигма1 (MBт)</h3>
//             <input type="text" name="sigma1" placeholder="" required>
// 			<h3>Сигма2 (MBт)</h3>
//             <input type="text" name="sigma2" placeholder="" required>
// 			<h3>Вартість електроенергії (грн/кВт*год)</h3>
// 			<input type="text" name="price" placeholder="" required>
//         </div>
//         <button type="submit">Розрахувати</button>
//     </form>
//     {{if .}}
//         <pre>{{.}}</pre>
//     {{end}}
// </body>
// </html>
// `))

// // Основна ф-ія програми
// func calcHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		// Задаємо значення для різних видів пального
// 		Pc, err1 := strconv.ParseFloat(r.FormValue("pc"), 64)
// 		sigma1, err2 := strconv.ParseFloat(r.FormValue("sigma1"), 64)
// 		sigma2, err3 := strconv.ParseFloat(r.FormValue("sigma2"), 64)
// 		price, err4 := strconv.ParseFloat(r.FormValue("price"), 64)

// 		errors := []error{err1, err2, err3, err4}

// 		for _, err := range errors {
// 			if err != nil {
// 				http.Error(w, "Введіть правильні числа", http.StatusBadRequest)
// 				return
// 			}
// 		}

// 		prib1 := calculatePribytok(Pc, price, sigma1)
// 		prib2 := calculatePribytok(Pc, price, sigma2)

// 		result := fmt.Sprintf("Прибуток для системи 1 = %.2f тис.грн\n", prib1)
// 		result += fmt.Sprintf("Прибуток для вдосконаленої системи 2 = %.2f тис.грн\n", prib2)

// 		tmpl.Execute(w, result)
// 		return
// 	}

// 	tmpl.Execute(w, nil)
// }

// // Функція для обчислення нормального розподілу в точці p
// func normalDistribution(p, Pc, sigma float64) float64 {
// 	return (1 / (sigma * math.Sqrt(2*math.Pi))) * math.Exp(-0.5*math.Pow((p-Pc)/sigma, 2))
// }

// // Функція для обчислення інтегралу нормального розподілу в межах від a до b
// func calculateIntegral(a, b, Pc, sigma float64) float64 {
// 	step := 0.001 // Крок інтегрування
// 	sum := 0.0
// 	for x := a; x <= b; x += step {
// 		sum += normalDistribution(x, Pc, sigma) * step
// 	}
// 	return sum
// }

// // Функція для знаходження частки енергії
// func findPercentEnergy(Pc, sigma float64) float64 {
// 	a := 4.75 // Нижня межа інтегрування
// 	b := 5.25 // Верхня межа інтегрування
// 	return calculateIntegral(a, b, Pc, sigma)
// }

// // Функція для обчислення прибутку
// func calculatePribytok(Pct, vartist, sigma float64) float64 {
// 	Pc := Pct
// 	sigmaW1 := findPercentEnergy(Pc, sigma)
// 	W1 := Pc * 24 * sigmaW1
// 	P1 := W1 * vartist
// 	W2 := Pc * 24 * (1 - sigmaW1)
// 	shtraf := W2 * vartist
// 	prib1 := P1 - shtraf
// 	return prib1
// }

// func main() {
// 	http.HandleFunc("/", calcHandler)
// 	fmt.Println("http://localhost:8080")
// 	http.ListenAndServe(":8080", nil)
// }
