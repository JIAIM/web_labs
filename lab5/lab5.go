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
// 			margin-bottom: 5px;
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
//             width: 200px;
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
// 			<h3>Частота відмов одноколкової системи</h3>
//             <input type="text" name="W" placeholder="" required>
// 			<h3>Ціна за збитки від аварійних вимкнень електропостачання у грн/кВт * год</h3>
//             <input type="text" name="Za" placeholder="" required>
// 			<h3>Ціна за збитки від планових вимкнень електропостачання у грн/кВт * год</h3>
// 			<input type="text" name="Zp" placeholder="" required>
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
// 		W, err1 := strconv.ParseFloat(r.FormValue("W"), 64)
// 		Za, err2 := strconv.ParseFloat(r.FormValue("Za"), 64)
// 		Zp, err3 := strconv.ParseFloat(r.FormValue("Zp"), 64)

// 		errors := []error{err1, err2, err3}

// 		for _, err := range errors {
// 			if err != nil {
// 				http.Error(w, "Введіть правильні числа", http.StatusBadRequest)
// 				return
// 			}
// 		}

// 		result := fmt.Sprintln(findW2oc(W))
// 		zbitky := findZbitky(Za, Zp)
// 		result += fmt.Sprintf("Збитки від переривання електропостачання %.2f грн\n", zbitky)

// 		tmpl.Execute(w, result)
// 		return
// 	}

// 	tmpl.Execute(w, nil)
// }

// // Функція для знаходження частоти відмов двоколкової системи
// func findW2oc(Woc float64) string {
// 	Wcv := 0.02
// 	tvoc := 10.7
// 	kpmax := 43.0 / 8760.0
// 	kaoc := math.Round(Woc*tvoc*10000.0/8760.0) / 10000.0
// 	kpoc := math.Round(1.2*kpmax*10000.0) / 10000.0
// 	Wdk := 2 * Woc * (kaoc + kpoc)
// 	Wdc := math.Round((Wdk+Wcv)*10000.0) / 10000.0

// 	result := ""

// 	if Woc == Wdc {
// 		result += "Надійність систем однакова\n"
// 	} else if Woc < Wdc {
// 		result += "Надійність одноколкової системи електропередачі є вищою ніж двоколкової\n"
// 	} else {
// 		result += "Надійність двоколкової системи електропередачі є вищою ніж одноколкової\n"
// 	}

// 	result += fmt.Sprintf("Частота відмов одноколкової системи = %.3f\n", Woc)
// 	result += fmt.Sprintf("Частота відмов двоколкової системи = %.4f\n", Wdc)

// 	return result
// }

// // Функція для розрахунку збитків від переривання електропостачання
// func findZbitky(Za, Zp float64) float64 {
// 	Ma := 14900.0
// 	Mp := 132400.0
// 	return Za*Ma + Zp*Mp
// }

// func main() {
// 	http.HandleFunc("/", calcHandler)
// 	fmt.Println("http://localhost:8080")
// 	http.ListenAndServe(":8080", nil)
// }
