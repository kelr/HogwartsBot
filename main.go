package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	token   = ""
	adminID = ""
	hatID   = ""
)

var (
	houseCount = map[string]int{"hufflepuff": 0, "ravenclaw": 0, "gryffindor": 0, "slytherin": 0}
)

// Handle !points command, return map of houseCount values sorted by value
func handlePointsList() string {
	type keyValue struct {
		Key   string
		Value int
	}

	var keys []keyValue
	for house, count := range houseCount {
		keys = append(keys, keyValue{house, count})
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Value > keys[j].Value
	})

	var msg string
	for _, keyVal := range keys {
		msg += keyVal.Key + ": " + strconv.Itoa(keyVal.Value) + "\n"
	}

	return msg
}

// Handle commands that change the point count of a house
func handlePointChange(content string) string {
	// Ensure command has at least 4 splittable parts
	split := strings.SplitN(content, " ", 4)
	if len(split) != 4 {
		return "I couldn't understand that. Try: 5 points to house kybo"
	}

	// Parse point change to int
	points, err := strconv.Atoi(split[0])
	if err != nil {
		return "I can't understand how many points that is. Numbers doofus, do you know how to use them?"
	}

	// Check if the house exists
	house := strings.ToLower(split[3])
	if _, ok := houseCount[house]; !ok {
		houseCount[house] = points
		return "That ain't no house I ever heard of but I'll add it."
	}

	houseCount[house] += points
	writeFile()

	return split[3] + " now has " + strconv.Itoa(houseCount[house]) + " points"
}

// Handle commands that delete a house
func handleHouseDelete(content string) string {
	// Ensure command has at least 4 splittable parts
	split := strings.SplitN(content, " ", 2)
	if len(split) != 2 {
		return "I couldn't understand that. Try: !delete hufflepuff"
	}

	house := strings.ToLower(split[1])
	if _, ok := houseCount[house]; !ok {
		return "I can't delete something that doesn't exist you absolute buffoon."
	}

	delete(houseCount, house)
	writeFile()
	return split[1] + " has been destroyed."
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!points") {
		s.ChannelMessageSend(m.ChannelID, handlePointsList())
	} else if strings.Index(m.Content, "points to") != -1 {
		// Only authorized users can edit the points
		if (m.Author.ID == adminID) || (m.Author.ID == hatID) {
			s.ChannelMessageSend(m.ChannelID, handlePointChange(m.Content))
		}
	} else if strings.HasPrefix(m.Content, "!delete") {
		// Only authorized users can delete houses
		if (m.Author.ID == adminID) || (m.Author.ID == hatID) {
			s.ChannelMessageSend(m.ChannelID, handleHouseDelete(m.Content))
		}
	}
}

func main() {
	loadFile()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	dg.AddHandler(onMessage)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Bot is now running.")

	// Sleeb forever
	select {}
}
