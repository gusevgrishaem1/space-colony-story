package server

import (
	"fmt"
	"html/template"
	"net/http"
)

type StoryPart struct {
	Text    string
	Choices []Choice
}

type Choice struct {
	Text string
	Next string
}

var storyData = map[string]StoryPart{
	"start": {
		Text: "Welcome to the colony, crewmember. The power systems have failed, and the station is on the verge of collapse. You have to make decisions that will determine your fate.",
		Choices: []Choice{
			{"Inspect the power core.", "core"},
			{"Repair the communication array.", "array"},
			{"Investigate the mysterious signals.", "signals"},
		},
	},
	"core": {
		Text: "You approach the power core. Sparks are flying, and it’s a dangerous place. Do you attempt to fix it?",
		Choices: []Choice{
			{"Try to repair it.", "repair_core"},
			{"Leave it and check the communication array.", "array"},
		},
	},
	"array": {
		Text: "You head to the communication array. It’s barely functioning, but you think you can fix it. What will you do?",
		Choices: []Choice{
			{"Try to repair the array.", "repair_array"},
			{"Check the mysterious signals instead.", "signals"},
		},
	},
	"signals": {
		Text: "You decode the signals and discover they might be coming from an alien ship. What do you do?",
		Choices: []Choice{
			{"Try to communicate with them.", "alien_communication"},
			{"Ignore the signals and focus on repairs.", "core"},
		},
	},
	"repair_core": {
		Text: "You carefully attempt to repair the power core. After a tense few moments, you manage to stabilize it, restoring power to critical systems. The station is saved—for now. What’s your next move?",
		Choices: []Choice{
			{"Check the communication array.", "array"},
			{"Investigate the mysterious signals.", "signals"},
		},
	},
	"repair_array": {
		Text: "You repair the communication array, restoring full functionality. You manage to send a distress signal to the nearest space station, hoping for rescue. The crisis is over, but you're still stranded. What now?",
		Choices: []Choice{
			{"Check the power core.", "core"},
			{"Investigate the mysterious signals.", "signals"},
		},
	},
	"alien_communication": {
		Text: "You initiate communication with the alien ship. They respond, but their message is cryptic and seems to indicate they are aware of the station's collapse. Suddenly, a small alien vessel docks with the station. Do you approach it?",
		Choices: []Choice{
			{"Approach the alien vessel.", "alien_encounter"},
			{"Ignore the alien vessel and focus on repairs.", "repair_array"},
		},
	},
	"alien_encounter": {
		Text: "You approach the alien vessel cautiously. The door opens, and a humanoid alien steps out, offering a device that could help repair the station permanently. However, it comes with a mysterious warning about future consequences. Do you accept their help?",
		Choices: []Choice{
			{"Accept the alien's help.", "accept_help"},
			{"Decline and try to fix it yourself.", "repair_core"},
		},
	},
	"accept_help": {
		Text: "The alien provides a powerful energy source, and with it, you manage to repair the station's core. The alien ship departs, leaving you with a sense of both relief and unease. The station is saved, but what will the future hold? You wonder if this alien encounter will have consequences down the line.",
		Choices: []Choice{
			{"Live to fight another day.", "end_good"},
		},
	},
	"decline_help": {
		Text: "You decline the alien's offer and decide to fix the station yourself. It takes several tense hours, but you manage to restore power and communication systems. The crisis is over, but you remain uneasy about the missed opportunity.",
		Choices: []Choice{
			{"Live to fight another day.", "end_neutral"},
		},
	},
	"end_good": {
		Text: "With the station repaired and a future uncertain, you breathe a sigh of relief. The rescue ship arrives, and you're saved. Your decision to trust the alien might have saved the colony in ways you’ll never fully understand. You’re safe, but your journey is far from over.",
		Choices: []Choice{
			{"End the story.", "end"},
		},
	},
	"end_neutral": {
		Text: "The station is back online, and you survive. The distress signal has gone out, but you can’t shake the feeling that there might have been a better way to handle things. Perhaps the mysterious alien race will be watching you in the future.",
		Choices: []Choice{
			{"End the story.", "end"},
		},
	},
	"end": {
		Text:    "The adventure ends here. Will you ever be called back to adventure? Who knows. But for now, you can rest knowing you survived the collapse of the colony. Until next time.",
		Choices: []Choice{},
	},
}

// Serve the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

// Serve the dynamic story content based on the choice
func storyHandler(w http.ResponseWriter, r *http.Request) {
	choice := r.URL.Query().Get("choice")
	if choice == "" {
		choice = "start"
	}

	part := storyData[choice]

	var choicesHTML string
	for _, c := range part.Choices {
		choicesHTML += fmt.Sprintf(`<div class="choice p-4 m-2 bg-gray-800 text-white rounded-lg cursor-pointer hover:bg-gray-600 transition" hx-get="/next?choice=%s" hx-target="#content">%s</div>`, c.Next, c.Text)
	}

	contentHTML := fmt.Sprintf(`<p class="text-lg mb-4">%s</p>%s`, part.Text, choicesHTML)
	w.Write([]byte(contentHTML))
}

func StartServer() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/next", storyHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", nil)
}
