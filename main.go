package main

import (
	"encoding/json"
	"fmt"
	"os"
	
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Expense struct {
	Username string  `json:"username"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

var usersFile = "users.json"
var expensesFile = "expenses.json"

// --- File helpers ---
func loadUsers() []User {
	data, err := os.ReadFile(usersFile)
	if err != nil {
		return []User{}
	}
	var users []User
	json.Unmarshal(data, &users)
	return users
}

func saveUsers(users []User) {
	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile(usersFile, data, 0644)
}

func loadExpenses() []Expense {
	data, err := os.ReadFile(expensesFile)
	if err != nil {
		return []Expense{}
	}
	var expenses []Expense
	json.Unmarshal(data, &expenses)
	return expenses
}

func saveExpenses(expenses []Expense) {
	data, _ := json.MarshalIndent(expenses, "", "  ")
	os.WriteFile(expensesFile, data, 0644)
}

// --- Auth functions ---
func userExists(username string, users []User) bool {
	for _, u := range users {
		if u.Username == username {
			return true
		}
	}
	return false
}

func register() {
	users := loadUsers()
	var username, password string
	fmt.Println("=== Register ===")
	fmt.Print("Enter username: ")
	fmt.Scan(&username)

	if userExists(username, users) {
		fmt.Println("User already exists. Please login.")
		return
	}

	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	users = append(users, User{Username: username, Password: password})
	saveUsers(users)
	fmt.Println("Registration successful! Please login now.")
}

func login() *User {
	users := loadUsers()
	if len(users) == 0 {
		fmt.Println("No users found. Please register first.")
		register()
		return nil
	}

	var username, password string
	fmt.Println("=== Login ===")
	fmt.Print("Enter username: ")
	fmt.Scan(&username)
	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	for _, u := range users {
		if u.Username == username && u.Password == password {
			fmt.Println("Login successful!")
			return &u
		}
	}
	fmt.Println("Invalid credentials.")
	return nil
}

// --- Expense Tracker ---
func addExpense(username string) {
	expenses := loadExpenses()
	categories := []string{"Travel", "Food", "Fun", "Work"}
	fmt.Println("=== Add Expense ===")
	fmt.Println("Choose category:")
	for i, c := range categories {
		fmt.Printf("%d. %s\n", i+1, c)
	}
	var choice int
	fmt.Print("Enter choice: ")
	fmt.Scan(&choice)
	if choice < 1 || choice > len(categories) {
		fmt.Println("Invalid category.")
		return
	}

	var amount float64
	fmt.Print("Enter amount: ")
	fmt.Scan(&amount)

	expenses = append(expenses, Expense{
		Username: username,
		Category: categories[choice-1],
		Amount:   amount,
	})
	saveExpenses(expenses)
	fmt.Println("Expense added successfully!")
}

func viewReport(username string) {
	expenses := loadExpenses()
	fmt.Println("=== Expense Report ===")
	total := 0.0
	categoryTotals := map[string]float64{}

	for _, e := range expenses {
		if e.Username == username {
			fmt.Printf("%s - %.2f\n", e.Category, e.Amount)
			total += e.Amount
			categoryTotals[e.Category] += e.Amount
		}
	}

	if total == 0 {
		fmt.Println("No expenses found.")
		return
	}

	fmt.Printf("Total spent: %.2f\n", total)

	// find max category
	maxCategory := ""
	maxAmount := 0.0
	for cat, amt := range categoryTotals {
		if amt > maxAmount {
			maxAmount = amt
			maxCategory = cat
		}
	}
	fmt.Printf("You spend the most on: %s (%.2f)\n", maxCategory, maxAmount)
}

func expenseMenu(user *User) {
	for {
		fmt.Println("\n=== Expense Tracker Menu ===")
		fmt.Println("1. Add Expense")
		fmt.Println("2. View Report")
		fmt.Println("3. Logout")
		var choice int
		fmt.Print("Enter choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			addExpense(user.Username)
		case 2:
			viewReport(user.Username)
		case 3:
			fmt.Println("Logged out successfully.")
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}

func main() {
	for {
		fmt.Println("\n=== Main Menu ===")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Exit")
		var choice int
		fmt.Print("Enter choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			user := login()
			if user != nil {
				expenseMenu(user)
			}
		case 2:
			register()
		case 3:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
