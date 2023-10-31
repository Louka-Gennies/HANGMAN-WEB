package main

import (
	"bufio"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
	"strings"
)

type ContactDetails struct {
	LettersGood []string
	LettersWrong []string
	RandomWord string
	WordFind string
	Status string
	Try int
}

var details ContactDetails

func letterExists(letters []string, letter string) bool {
	for _, l := range letters {
		if l == letter {
			return true
		}
	}
	return false
}

func WordList(textFile string) (string, error) {
	file, err := os.Open(textFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var wordList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordList = append(wordList, scanner.Text())
	}

	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(wordList))
	randomWord := wordList[randomIndex]

	return randomWord, nil
}

func Verify(word, letter string) bool {
    WordTab := []rune(word)      // Convert the target word to a rune slice for character comparison
    RuneLetter := []rune(letter) // Convert the input letter to a rune slice for comparison
	correct := false

    // Iterate through the target word to find occurrences of the input letter
    for i := 0; i < len(WordTab); i++ {
        if RuneLetter[0] == WordTab[i] {
            correct = true
			break
        } 
    }
    return correct
}

func VerifyIndice(word, letter string) []int {
    WordTab := []rune(word)      // Convert the target word to a rune slice for character comparison
    RuneLetter := []rune(letter) // Convert the input letter to a rune slice for comparison
    var indices []int            // Initialize a slice to store indices where the letter is found

    // Iterate through the target word to find occurrences of the input letter
    for i := 0; i < len(WordTab); i++ {
        if RuneLetter[0] == WordTab[i] {
            indices = append(indices, i) // Add the index to the slice if the letter is found
        }
    }

    // If no occurrences of the letter are found, return nil
    if len(indices) == 0 {
        return nil
    }

    return indices
}

func PrintWord(word string) string {
    rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time.

    // Calculate the number of letters to reveal (between 1 and len(word)/2 - 1)
    revealedCount := len(word)/2 - 1

    // Generate a random set of indices to reveal
    revealedIndices := make([]int, revealedCount)
    for i := 0; i < revealedCount; i++ {
        randomIndex := rand.Intn(len(word))
        revealedIndices[i] = randomIndex
    }

    var str string

    for i := 0; i < len(word); i++ {
        revealed := false
        for _, index := range revealedIndices {
            if i == index {
                str += string(word[i])
                revealed = true
                break
            }
        }
        if !revealed {
            str += "_"
        }
    }

    return str
}

func RevealLetters(word string, indices []int, revealedWord string) string {
    revealed := []rune(revealedWord) // Convert the revealed word to a rune slice for modification
    WordTab := []rune(word) // Convert the target word to a rune slice for access

    // Iterate through the provided indices and update the revealed word
    for _, index := range indices {
        if index >= 0 && index < len(WordTab) {
            revealed[index] = WordTab[index]
        }
    }

    return string(revealed) // Convert the updated revealed word back to a string
}


func main() {
	tmpl := template.Must(template.ParseFiles("template/forms.html"))

	details.RandomWord, _ = WordList("words.txt")

	details.WordFind = PrintWord(details.RandomWord)

	details.Try = 10

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, details)
			return
		}

		letter := strings.ToUpper(r.FormValue("letter"))

		correctLetter := Verify(details.RandomWord, letter)

		if len(letter) == 1 && !letterExists(details.LettersGood, letter) && !letterExists(details.LettersWrong, letter) {
			if correctLetter == true {
				details.LettersGood = append(details.LettersGood, letter)
			} else {
				details.LettersWrong = append(details.LettersWrong, letter)
				details.Try -= 1
			}
		}

		indice := VerifyIndice(details.RandomWord, letter)

		details.WordFind = RevealLetters(details.RandomWord, indice, details.WordFind)

		if details.RandomWord == details.WordFind{
			details.Status = "Bravo tu as trouvé le mot !!"
		}

		if details.Try == 0 {
			details.Status = "Mince tu as perdu..."
		}

		tmpl.Execute(w, details)
	})

	http.ListenAndServe(":8080", nil)
}
