package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Expense struct {
	User     string
	Amount   float64
	Category string
	Date     string
}

var dataFile = "expenses.json"

func addFromInput(user, amtStr, category string, w io.Writer) {
	amount, err := strconv.ParseFloat(amtStr, 64)
	if err != nil {
		fmt.Fprintln(w, "Invalid amount")
		return
	}

	f, _ := os.ReadFile(dataFile)
	var expenses []Expense
	json.Unmarshal(f, &expenses)

	expenses = append(expenses, Expense{user, amount, category, time.Now().Format("2006-01-02")})
	out, _ := json.MarshalIndent(expenses, "", "  ")
	os.WriteFile(dataFile, out, 0644)
	fmt.Fprintln(w, "Expense added.")
}

func sendReport(user string, w io.Writer) {
	f, _ := os.ReadFile(dataFile)
	var expenses []Expense
	json.Unmarshal(f, &expenses)

	total := 0.0
	for _, e := range expenses {
		if e.User == user {
			fmt.Fprintf(w, "$%.2f - %s (%s)\n", e.Amount, e.Category, e.Date)
			total += e.Amount
		}
	}
	fmt.Fprintf(w, "Total: $%.2f\n", total)
}

func exportCSV(user string, w io.Writer) {
	f, _ := os.ReadFile(dataFile)
	var expenses []Expense
	json.Unmarshal(f, &expenses)

	out, _ := os.Create(user + "_report.csv")
	defer out.Close()
	fmt.Fprintln(out, "Amount,Category,Date")

	for _, e := range expenses {
		if e.User == user {
			fmt.Fprintf(out, "%.2f,%s,%s\n", e.Amount, e.Category, e.Date)
		}
	}
	fmt.Fprintln(w, "CSV exported.")
}
