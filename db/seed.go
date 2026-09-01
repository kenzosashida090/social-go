package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/kenzosashida090/social/store"
)

var usernames = []string{
	"ShadowByte", "NovaStrike", "GhostRunner", "CyberWolf", "FrostViper",
	"DarkOrbit", "NeonPhantom", "VoidWalker", "ApexShadow", "LunarFang",
	"SilentRogue", "IronSpecter", "BlazeKnight", "CrimsonByte", "EchoHunter",
	"StormVex", "NightCipher", "ToxicNova", "RoguePulse", "PhantomCore",
	"ZeroFrost", "VenomRush", "StaticGhost", "QuantumWolf", "InfernoX",
	"ArcticRogue", "ChaosByte", "MysticFury", "TurboPhantom", "CyberFang",
	"DarkVortex", "NovaReaper", "GhostByte", "ShadowPulse", "FrostNova",
	"VoidHunter", "NeonViper", "SilentFang", "LunarGhost", "StormCipher",
	"BlazeVortex", "CrimsonWolf", "EchoRogue", "IronViper", "NightFury",
	"ToxicShadow", "RogueNova", "PhantomRush", "ZeroVortex", "VenomByte",
}
var titles = []string{
	"Best Gaming Setup", "New Game Update", "Looking for Teammates", "My First Win", "Crazy Match",
	"Game Night", "Tips for Beginners", "Favorite Character", "Epic Comeback", "Ranked Grind",
	"New Strategy", "What Do You Think?", "Funny Game Moment", "Almost Won", "Finally Ranked",
	"Best Weapon", "Need Some Advice", "That Was Close", "Weekend Tournament", "Game Recommendations",
}

var contents = []string{
	"Just had an amazing game with the team.",
	"Anyone else enjoying the latest update?",
	"I finally reached a new rank today.",
	"Looking for some people to play with tonight.",
	"That was probably my craziest match yet.",
	"Does anyone have tips for improving at this game?",
	"I tried a new strategy and it actually worked.",
	"What's everyone's favorite character right now?",
	"We almost lost but somehow managed to make a comeback.",
	"Just wanted to share this funny moment from my game.",
	"I've been grinding ranked all week and finally made it.",
	"What strategy do you usually use when playing?",
	"Anyone interested in joining a tournament this weekend?",
	"I need some advice on improving my gameplay.",
	"That last match was incredibly close.",
	"What's the best setup for this game?",
	"I just discovered a really useful trick.",
	"Who wants to play a few matches tonight?",
	"Finally got the achievement I've been trying to unlock.",
	"Drop your favorite gaming tips below.",
}

var tags = []string{
	"gaming", "fun", "multiplayer", "strategy", "community",
	"games", "ranked", "tips", "tournament", "chat",
}
var commentsExample = []string{
	"This is such a great post!",
	"Really interesting perspective.",
	"I completely agree with this.",
	"This is definitely worth thinking about.",
	"Great insights here!",
	"Thanks for sharing this.",
	"This caught my attention right away.",
	"Absolutely love this!",
	"Very well said.",
	"This makes a lot of sense.",
	"I hadn't thought about it this way before.",
	"Such a valuable perspective.",
	"Couldn't agree more!",
	"This is really interesting.",
	"Great point!",
	"Definitely something more people should see.",
	"Appreciate you sharing this.",
	"This is a good reminder.",
	"Really enjoyed reading this.",
	"Looking forward to seeing more posts like this.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Println("Error", err.Error())
			return
		}
	}
	generatedPosts := generatePosts(20, users)
	for _, post := range generatedPosts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error posts", err.Error())
			return
		}
	}
	generateComments := generateComments(20, generatedPosts)
	for _, comment := range generateComments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("ERror comment", err.Error())
			return
		}
	}
}

func generateUsers(size int) []*store.User {
	users := make([]*store.User, size)
	for i := range size {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}

func generatePosts(size int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, size)
	for i := range size {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(size int, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, size)
	for i := range posts {
		post := posts[rand.Intn(len(posts))]
		comments[i] = &store.Comment{
			Content: commentsExample[rand.Intn(len(commentsExample))],
			PostID:  post.ID,
		}
	}
	return comments
}
