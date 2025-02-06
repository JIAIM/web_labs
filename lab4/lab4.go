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
//             <input type="text" name="strumKZ" placeholder="Струм КЗ(кА)" required>
//             <input type="text" name="potyzhnist" placeholder="Потужність станції(МВа)" required>
//             <input type="text" name="Ukmax" placeholder="Ukmax(%)" required>
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
// 		strumKZ, err1 := strconv.ParseFloat(r.FormValue("strumKZ"), 64)
// 		potyzhnist, err2 := strconv.ParseFloat(r.FormValue("potyzhnist"), 64)
// 		Ukmax, err3 := strconv.ParseFloat(r.FormValue("Ukmax"), 64)

// 		errors := []error{err1, err2, err3}

// 		for _, err := range errors {
// 			if err != nil {
// 				http.Error(w, "Введіть правильні числа", http.StatusBadRequest)
// 				return
// 			}
// 		}

// 		minDiameter := findDiameterOfWire(strumKZ)
// 		result := fmt.Sprintf("Mінімальний діаметр провода повинен бути >= %.2f мм^2\n", minDiameter)

// 		startIp := findStartKZ(potyzhnist)
// 		result += fmt.Sprintf("Початкове значення струму КЗ на шинах 10В = %.2f кА\n", startIp)

// 		result += findKzKhmelnytsk(Ukmax)

// 		tmpl.Execute(w, result)
// 		return
// 	}

// 	tmpl.Execute(w, nil)
// }

// // Функція для знаходження діаметра кабеля
// func findDiameterOfWire(Ik float64) float64 {
// 	return Ik * math.Sqrt(2.5) / 92
// }

// // Функція для знаходження початкового струму короткого замикання
// func findStartKZ(Sk float64) float64 {
// 	Usn := 10.5
// 	Uk := 10.5
// 	Snomt := 6.3
// 	Xc := Usn * Usn / Sk
// 	Xt := (Uk * Usn * Usn) / (100 * Snomt)
// 	Xs := Xc + Xt
// 	Ip := Usn / (math.Sqrt(3.0) * Xs)
// 	return Ip
// }

// // Функція для розрахунку короткого замикання в Хмельницькому
// func findKzKhmelnytsk(Ukmax float64) string {
// 	Uvn := 115.0
// 	Snomt := 6.3
// 	Xt := (Ukmax * Uvn * Uvn) / (100 * Snomt)

// 	Rsn := 10.65
// 	Rsh := Rsn
// 	Xsn := 24.02
// 	Xsh := Xsn + Xt
// 	Zsh := math.Sqrt(Rsh*Rsh + Xsh*Xsh)

// 	Xcmin := 65.68
// 	Rcmin := 34.88
// 	Rshmin := Rcmin
// 	Xshmin := Xcmin + Xt
// 	Zshmin := math.Sqrt(Rshmin*Rshmin + Xshmin*Xshmin)

// 	I3sh := Uvn * 1000 / (math.Sqrt(3.0) * Zsh)
// 	I2sh := I3sh * math.Sqrt(3.0) / 2.0
// 	I3shmin := Uvn * 1000 / (math.Sqrt(3.0) * Zshmin)
// 	I2shmin := I3shmin * math.Sqrt(3.0) / 2.0

// 	Unn := 11.0
// 	kpr := (Unn * Unn) / (Uvn * Uvn)

// 	Rshn := Rsh * kpr
// 	Xshn := Xsh * kpr
// 	Zshn := math.Sqrt(Rshn*Rshn + Xshn*Xshn)
// 	Rshnmin := Rshmin * kpr
// 	Xshnmin := Xshmin * kpr
// 	Zshnmin := math.Sqrt(Rshnmin*Rshnmin + Xshnmin*Xshnmin)

// 	I3shn := Unn * 1000 / (math.Sqrt(3.0) * Zshn)
// 	I2shn := I3shn * math.Sqrt(3.0) / 2.0
// 	I3shnmin := Unn * 1000 / (math.Sqrt(3.0) * Zshnmin)
// 	I2shnmin := I3shnmin * math.Sqrt(3.0) / 2.0

// 	result := fmt.Sprintf("Нормальний режим - Трифазне КЗ для 110кв = %.2f А\n", I3sh)
// 	result += fmt.Sprintf("Мінімальний режим - Трифазне КЗ для 110кв = %.2f А\n", I3shmin)
// 	result += fmt.Sprintf("Нормальний режим - Двофазне КЗ для 110кв = %.2f А\n", I2sh)
// 	result += fmt.Sprintf("Мінімальний режим - Двофазне КЗ для 110кв = %.2f А\n", I2shmin)
// 	result += fmt.Sprintf("Нормальний режим - Трифазне КЗ для 10кв = %.2f А\n", I3shn)
// 	result += fmt.Sprintf("Мінімальний режим - Трифазне КЗ для 10кв = %.2f А\n", I3shnmin)
// 	result += fmt.Sprintf("Нормальний режим - Двофазне КЗ для 10кв = %.2f А\n", I2shn)
// 	result += fmt.Sprintf("Мінімальний режим - Двофазне КЗ для 10кв = %.2f А\n", I2shnmin)
// 	return result
// }

// func main() {
// 	http.HandleFunc("/", calcHandler)
// 	fmt.Println("http://localhost:8080")
// 	http.ListenAndServe(":8080", nil)
// }
