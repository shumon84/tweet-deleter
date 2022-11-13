package main

import (
	"bufio"
	"context"
	"github.com/sivchari/gotwtr"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type TweetDeleter struct {
	client       *gotwtr.Client
	targetIDs    *Queue
	numOfTargets int
}

func NewTweetDeleter(listFilePath string, httpClient *http.Client) (*TweetDeleter, error) {
	listFiles, err := os.Open(listFilePath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(listFiles)
	idList := make([]string, 0, 256)
	for scanner.Scan() {
		idList = append(idList, scanner.Text())
	}
	deleteTargetTweetIDs := NewQueue()
	deleteTargetTweetIDs.BulkPush(idList)
	return &TweetDeleter{
		client:       gotwtr.New(bearerToken, gotwtr.WithHTTPClient(httpClient)),
		targetIDs:    deleteTargetTweetIDs,
		numOfTargets: len(idList),
	}, nil
}

func (t *TweetDeleter) delete(ctx context.Context) {
	go func() {
		rateLimit := time.Tick(time.Minute * 15)
		for {
			wg := sync.WaitGroup{}
			wg.Add(50)
			for offset := 0; offset < 50; offset++ {
				tweetID, err := t.targetIDs.Pop()
				if err != nil {
					break
				}
				go func(ctx context.Context, tweetID string) {
					defer func() {
						wg.Done()
						log.Println("[Info] TweetID =", tweetID, ": DONE")
					}()
					resp, err := t.client.DeleteTweet(ctx, tweetID)
					if err != nil {
						log.Println("[Error] TweetID =", tweetID, ":", err)
					}
					if !resp.Data.Deleted {
						log.Println("[Info] TweetID =", tweetID, ": 削除に失敗したため、再度エンキューします")
						t.targetIDs.Push(tweetID)
					} else {
						log.Println("[Info] TweetID =", tweetID, ": 削除成功")
					}
				}(ctx, tweetID)
			}
			wg.Wait()
			if t.targetIDs.IsEmpty() {
				break
			}
			<-rateLimit
		}
	}()
}

func (t *TweetDeleter) getProgress() (num int, denom int) {
	denom = t.numOfTargets
	num = denom - len(t.targetIDs.q)
	return
}
