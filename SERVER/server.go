package main

import (
	"html/template"
	"net/http"
	"strings"
	"ytrack.learn.ynov.com/git/glouka/HANGMAN.git"
)

type ContactDetails struct {
	LettersGood  []string
	LettersWrong []string
	RandomWord   string
	WordFind     string
	Status       string
	Try          int
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

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

    tmpl := template.Must(template.ParseFiles("template/forms.html"))

    details.RandomWord, _ = WordList("words.txt")

    details.WordFind = PrintWord(details.RandomWord)

    details.Try = 10

    http.HandleFunc("/hangman", func(w http.ResponseWriter, r *http.Request) {
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

        if details.RandomWord == details.WordFind {
            http.Redirect(w, r, "/victory", http.StatusSeeOther)
            return
        }

        if details.Try == 0 {
            http.Redirect(w, r, "/defeat", http.StatusSeeOther)
            return
        }

        tmpl.Execute(w, details)
    })

    http.HandleFunc("/victory", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("template/victory.html"))
		
		// Generate a new random word
		newRandomWord, err := WordList("words.txt")
		if err != nil {
			http.Error(w, "Error generating a new random word", http.StatusInternalServerError)
			return
		}
	
		// Reset all variables to their initial state with the new random word
		details.RandomWord = newRandomWord
		details.LettersGood = []string{}
		details.LettersWrong = []string{}
		details.WordFind = PrintWord(details.RandomWord)
		details.Try = 10

		tmpl.Execute(w, nil)
    })

    http.HandleFunc("/defeat", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("template/defeat.html"))

		// Generate a new random word
		newRandomWord, err := WordList("words.txt")
		if err != nil {
			http.Error(w, "Error generating a new random word", http.StatusInternalServerError)
			return
		}
	
		// Reset all variables to their initial state with the new random word
		details.RandomWord = newRandomWord
		details.LettersGood = []string{}
		details.LettersWrong = []string{}
		details.WordFind = PrintWord(details.RandomWord)
		details.Try = 10

		tmpl.Execute(w, nil)
    })

    http.ListenAndServe(":8080", nil)
}
