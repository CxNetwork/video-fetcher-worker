package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const developerKey = "[REDACTED]"
const projectID = "cx-network-204116"

type Video struct {
	id          string
	title       string
	publishedAt string
}

type Streamer struct {
	id      string
	channel string
}

func main() {
	ctx := context.Background()

	gClient := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(gClient)
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	opt := option.WithCredentialsFile("./[REDACTED].json")
	fClient, err := firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		log.Fatalf("Error creating new Firestore client: %v", err)
	}

	streamers := []Streamer{}

	iter := fClient.Collection("streamers").Where("enabled", "==", true).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		ytId := doc.Data()["socials"].(map[string]interface{})["yt"].(string)

		streamers = append(streamers, Streamer{id: doc.Ref.ID, channel: ytId})
	}

	ticker := time.NewTicker(time.Hour)

	func() {
		log.Println("Starting monitor...")
		for {
			select {
			case <-ticker.C:
				batch := fClient.Batch()
				for _, streamer := range streamers {
					strRef := fClient.Collection("streamers").Doc(streamer.id)
					video := getVideos(streamer, service)
					batch.Set(strRef, map[string]interface{}{
						"latestVideo": map[string]interface{}{
							"title":       video.title,
							"id":          video.id,
							"publishedAt": video.publishedAt,
						},
					}, firestore.MergeAll)
				}
				_, err := batch.Commit(ctx)

				log.Println("Firestore -> Batch update complete")

				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}()
}

func getVideos(streamer Streamer, service *youtube.Service) Video {
	call := service.Search.List("id,snippet").
		ChannelId(streamer.channel).
		Order("date").
		MaxResults(3)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	latestVideo := Video{id: response.Items[0].Id.VideoId, title: response.Items[0].Snippet.Title, publishedAt: response.Items[0].Snippet.PublishedAt}

	return latestVideo
}

func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}
