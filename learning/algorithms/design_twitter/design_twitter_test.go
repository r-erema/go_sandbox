package designtwitter_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDesignTwitter(t *testing.T) {
	t.Parallel()

	twitter := Constructor()
	twitter.PostTweet(1, 5)
	feed := twitter.GetNewsFeed(1)
	assert.Equal(t, []int{5}, feed)
	twitter.Follow(1, 2)
	twitter.PostTweet(2, 6)
	feed = twitter.GetNewsFeed(1)
	assert.Equal(t, []int{6, 5}, feed)
	twitter.Unfollow(1, 2)
	feed = twitter.GetNewsFeed(1)
	assert.Equal(t, []int{5}, feed)
}

type feed []tweet

type tweet struct {
	order, uid, tid int
}

type Twitter struct {
	autoIncrement     int
	userFeeds         map[int]*feed
	tweets            []tweet
	userFollowers     map[int]map[int]struct{}
	userSubscriptions map[int]map[int]struct{}
}

func Constructor() Twitter {
	return Twitter{
		userFeeds:         make(map[int]*feed),
		userFollowers:     make(map[int]map[int]struct{}),
		userSubscriptions: make(map[int]map[int]struct{}),
	}
}

// Time O(M*logN), where M is number of followers and N is tweets in feed(heap insertion is logN)
// Space O(M*N), M is users and N is total tweets.
func (tw *Twitter) PostTweet(userId int, tweetId int) {
	twt := tweet{tw.autoIncrement, userId, tweetId}
	tw.autoIncrement++

	tw.tweets = append(tw.tweets, twt)

	if _, ok := tw.userFeeds[userId]; !ok {
		tw.userFeeds[userId] = &feed{}
		heapify(*tw.userFeeds[userId])
	}

	push(tw.userFeeds[userId], twt)

	for followerID := range tw.userFollowers[userId] {
		push(tw.userFeeds[followerID], twt)
	}
}

// Time O(1), since the heap popping
// Space O(1), since the fixed feed length.
func (tw *Twitter) GetNewsFeed(userId int) []int {
	var res []int

	if tw.userFeeds[userId] == nil {
		return res
	}

	if f := tw.userFeeds[userId]; f != nil {
		tmp := slices.Clone(*f)

		for len(tmp) > 0 && len(res) < 10 {
			t := pop(&tmp)
			res = append(res, t.tid)
		}
	}

	return res
}

// Time O(M*logN), where M is number of followees and N is tweets in feed(heap insertion is logN)
// Space O(M*N), M is followees and N is followee tweets.
func (tw *Twitter) Follow(followerId int, followeeId int) {
	if tw.userFollowers[followeeId] == nil {
		tw.userFollowers[followeeId] = make(map[int]struct{})
	}

	tw.userFollowers[followeeId][followerId] = struct{}{}

	if tw.userSubscriptions[followerId] == nil {
		tw.userSubscriptions[followerId] = make(map[int]struct{})
	}

	tw.userSubscriptions[followerId][followeeId] = struct{}{}

	tw.refreshFeed(followerId)
}

// Time O(M*logN), where M is number of followees and N is tweets in feed(heap insertion is logN)
// Space O(M*N), M is followees and N is followee tweets.
func (tw *Twitter) Unfollow(followerId int, followeeId int) {
	delete(tw.userFollowers[followeeId], followerId)
	delete(tw.userSubscriptions[followerId], followeeId)
	tw.refreshFeed(followerId)
}

func (tw *Twitter) refreshFeed(userId int) {
	subscription := tw.userSubscriptions[userId]

	feed := new(feed)

	for i := range tw.tweets {
		if _, ok := subscription[tw.tweets[i].uid]; ok || tw.tweets[i].uid == userId {
			push(feed, tw.tweets[i])
		}
	}

	tw.userFeeds[userId] = feed
}

func heapify(arr []tweet) {
	for i := range arr {
		percolateUp(arr[:i+1])
	}
}

func push(heap *feed, node tweet) {
	*heap = append(*heap, node)
	percolateUp(*heap)
}

func pop(heap *feed) tweet {
	popped := (*heap)[0]
	(*heap)[0], *heap = (*heap)[len(*heap)-1], (*heap)[:len(*heap)-1]

	percolateDown(*heap)

	return popped
}

func percolateUp(arr []tweet) {
	if len(arr) < 2 {
		return
	}

	parentIndex := (len(arr) - 1 - 1) / 2

	if arr[parentIndex].order > arr[len(arr)-1].order {
		return
	}

	arr[parentIndex], arr[len(arr)-1] = arr[len(arr)-1], arr[parentIndex]

	percolateUp(arr[:parentIndex+1])
}

func percolateDown(arr []tweet) {
	currentNodeIndex := 0

	var bfs func()

	bfs = func() {
		leftChildIndex, rightChildIndex := currentNodeIndex*2+1, currentNodeIndex*2+2

		if leftChildIndex < len(arr) && arr[leftChildIndex].order > arr[currentNodeIndex].order {
			if rightChildIndex < len(arr) && arr[rightChildIndex].order > arr[leftChildIndex].order {
				arr[rightChildIndex], arr[currentNodeIndex] = arr[currentNodeIndex], arr[rightChildIndex]

				currentNodeIndex = rightChildIndex

				bfs()
			} else {
				arr[leftChildIndex], arr[currentNodeIndex] = arr[currentNodeIndex], arr[leftChildIndex]

				currentNodeIndex = leftChildIndex

				bfs()
			}

			return
		}

		if rightChildIndex < len(arr) && arr[rightChildIndex].order > arr[currentNodeIndex].order {
			arr[rightChildIndex], arr[currentNodeIndex] = arr[currentNodeIndex], arr[rightChildIndex]

			currentNodeIndex = rightChildIndex

			bfs()
		}
	}

	bfs()
}
