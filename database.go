package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Test with mySQL database")

	_, err := sql.Open("mysql", "root:Ls524454@tcp(localhost:3306)/words")

	if err != nil {
		fmt.Println("Error while connecting to database: %v", err)
	}

	fmt.Println("Success!")

	words := getWordsFromFile()

	splittedWords := splitWords(words)

	writeWordsToFile(splittedWords)

	// insertInDB(*db, splittedWords)
}

func getWordsFromFile() []string {
	file, _ := os.Open("spanish_copy.lex")

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	return text
}

type tableDic struct {
	name   string
	values []string
}

func splitWords(words []string) map[int]tableDic {
	var twoLetterWords []string
	var threeLetterWords []string
	var fourLetterWords []string
	var fiveLetterWords []string
	var sixLetterWords []string
	var sevenLetterWords []string
	var eightLetterWords []string

	for _, word := range words {
		switch len(word) {
		case 2:
			twoLetterWords = append(twoLetterWords, word)
		case 3:
			threeLetterWords = append(threeLetterWords, word)
		case 4:
			fourLetterWords = append(fourLetterWords, word)
		case 5:
			fiveLetterWords = append(fiveLetterWords, word)
		case 6:
			sixLetterWords = append(sixLetterWords, word)
		case 7:
			sevenLetterWords = append(sevenLetterWords, word)
		case 8:
			eightLetterWords = append(eightLetterWords, word)
		}
	}

	return map[int]tableDic{
		2: tableDic{
			name:   "twoletters",
			values: twoLetterWords,
		},
		3: tableDic{
			name:   "threeletters",
			values: threeLetterWords,
		},
		4: tableDic{
			name:   "fourletters",
			values: fourLetterWords,
		},
		5: tableDic{
			name:   "fiveletters",
			values: fiveLetterWords,
		},
		6: tableDic{
			name:   "sixletters",
			values: sixLetterWords,
		},
		7: tableDic{
			name:   "sevenletters",
			values: sevenLetterWords,
		},
		8: tableDic{
			name:   "eightletters",
			values: eightLetterWords,
		},
	}
}

func insertInDB(db sql.DB, splittedWords map[int]tableDic) {
	for _, wordInfo := range splittedWords {
		for _, value := range wordInfo.values {
			time.Sleep(100 * time.Millisecond)
			insert, err := db.Query("INSERT INTO " + wordInfo.name + " VALUES('" + value + "')")

			if err != nil {
				fmt.Println("INSERT INTO " + wordInfo.name + " VALUES('" + value + "')")
				fmt.Println("Error while inserting")
				fmt.Println(err.Error())
				if strings.Contains(err.Error(), "Duplicate") {
					continue
				}
				os.Exit(1)
			}

			defer insert.Close()
		}
	}

	fmt.Println("Successfully inserted values in db")
}

func writeWordsToFile(splittedWords map[int]tableDic) {
	for _, wordInfo := range splittedWords {
		file, _ := os.Create(wordInfo.name + ".csv")
		lineWriter := bufio.NewWriter(file)

		lineWriter.WriteString("\"word\"" + "\n")
		for _, value := range wordInfo.values {
			if strings.Contains(value, "ñ") || strings.Contains(value, " ") {
				continue
			}
			lineWriter.WriteString("\"" + value + "\"" + "\n")
		}

		file.Close()
	}
}
