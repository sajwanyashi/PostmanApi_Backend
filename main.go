package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

func GetItemsForMeal(filePath, day, meal string) []string {
	// Open the XLSX file
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}


	var items []string
	// Find the column index for the specified day
	var dayColumnIndex int
	headerRow := xlFile.Sheets[0].Rows[0]
	for i, cell := range headerRow.Cells {
		if strings.ToLower(cell.String()) == strings.ToLower(day) {
			dayColumnIndex = i
			break
		}
	}

	// Find the row index for the specified meal
	var mealRowIndex int
	for i, row := range xlFile.Sheets[0].Rows {
		if i == 0 {
			continue // Skip the header row
		}
		if strings.ToLower(row.Cells[dayColumnIndex].String()) == strings.ToLower(meal) {
			mealRowIndex = i
			break
		}
	}

	if mealRowIndex == 0 {
		fmt.Println("Invalid Input!")
		os.Exit(1)
	}

	// Iterate through cells in the column for the specified day
	// Start from the row after the meal row
	for _, row := range xlFile.Sheets[0].Rows[mealRowIndex+1:] {
		// Append items if the cell's value is not the same as the specified day
		if meal == "dinner" && row.Cells[dayColumnIndex].String() == "" {
			break // Stop appending if the meal is "dinner" and an empty cell is reached
		}
		if strings.ToLower(row.Cells[dayColumnIndex].String()) != strings.ToLower(day) {
			items = append(items, row.Cells[dayColumnIndex].String())
		} else {
			break // Stop appending if the cell's value is the same as the specified day
		}
	}
	return items
}

// Function to return the number of items in a meal
func NumberOfItemsInMeal(filePath, day, meal string) {
	// Call the GetItemsForMeal function to retrieve items for the specified meal
	items := len(GetItemsForMeal(filePath, day, meal))

	// Print the number of items
	fmt.Printf("Number of items in %s on %s: %d\n", meal, day, items)
}

// Function to check if the given item is in a specific meal for the given day
func IsItemInMeal(filePath, day, meal, item string) bool {
	// Get the items for the specified meal and day
	items := GetItemsForMeal(filePath, day, meal)

	// Iterate over the items and check if the given item is present
	for _, i := range items {
		if strings.ToLower(i) == strings.ToLower(item) {
			return true
		}
	}

	// Return false if the item is not found
	return false
}

// Function to convert the menu into JSON and save it to a file
func ConvertMenuToJSON() error {
	filePath := " "
	type MenuItem struct {
		Day   string
		Meal  string
		Items []string
	}

	var menu []MenuItem

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	meals := []string{"Breakfast", "Lunch", "Dinner"}
	for _, meal1 := range meals {
		for _, day := range days {

			items := GetItemsForMeal(filePath, day, meal1)

			// Create MenuItem
			menuItem := MenuItem{
				Day:   day,
				Meal:  meal1,
				Items: items,
			}

			// Append MenuItem to menu
			menu = append(menu, menuItem)

		}

	}
	// Convert menu to JSON
	jsonData, err := json.MarshalIndent(menu, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON data to file
	jsonFileName := strings.TrimSuffix(filePath, ".xlsx") + ".json"
	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Printf("Menu data successfully saved to %s\n", jsonFileName)
	return nil

}

func main() {

	var option int

	fmt.Println(`Choose the operation you would like to perform: 
	1. Input day and meal and get the corresponding items available for that meal.
	2. Input day and meal and get the number of items in that meal.
	3. Input day, meal and item and check if the given item is in that particular meal or not.
	4. Convert menu to a struct and then save is as a JSON file in the same directory`)
	fmt.Scan(&option)
	// Provide the path to the XLSX file
	filePath := " " 

	switch option {
	case 1:
		var day, meal string
		fmt.Print("Enter the day: ")
		fmt.Scanln(&day)
		fmt.Print("Enter the meal: ")
		fmt.Scanln(&meal)
		items := GetItemsForMeal(filePath, day, meal)

		fmt.Printf("Items for %s on %s:\n", meal, day)
		for _, item := range items {
			fmt.Println(item)
		}

	case 2:
		var day, meal string
		fmt.Print("Enter the day: ")
		fmt.Scanln(&day)
		fmt.Print("Enter the meal: ")
		fmt.Scanln(&meal)
		NumberOfItemsInMeal(filePath, day, meal)

	case 3:
		var day, meal, item string
		fmt.Print("Enter the day: ")
		fmt.Scanln(&day)
		fmt.Print("Enter the meal: ")
		fmt.Scanln(&meal)
		fmt.Print("Enter the item: ")
		fmt.Scanln(&item)
		// Check if the item is in the specified meal for the given day
		if IsItemInMeal(filePath, day, meal, item) {
			fmt.Printf("%s is in %s on %s\n", item, meal, day)
		} else {
			fmt.Printf("%s is not in %s on %s\n", item, meal, day)
		}

	case 4:
		ConvertMenuToJSON()

	}
}
