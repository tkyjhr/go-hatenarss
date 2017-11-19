package hatenarss

import (
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	defer os.Exit(code)
}

func TestGet(t *testing.T) {
	feed, err := Get(FeedCategoryAll, nil)
	if err != nil {
		t.Fatal(err)
	}
	if feed == nil {
		t.Fatal("feed == nil")
	}
}

func TestFilterByBookmarkCount(t *testing.T) {
	feed, err := Get(FeedCategoryAll, nil)
	if err != nil {
		t.Fatal(err)
	}
	threshold := 300
	FilterByBookmarkCount(&feed.Items, threshold)
	for _, v := range feed.Items {
		if v.BookmarkCount < threshold {
			t.Errorf("Bookmarkcount(%v) is less than the threshold(%v)", v.BookmarkCount, threshold)
		}
	}
}

func TestFilterByTitle(t *testing.T) {
	feed, err := Get(FeedCategoryAll, nil)
	if err != nil {
		t.Fatal(err)
	}
	word := "あ"
	FilterByTitle(&feed.Items, word)
	for _, v := range feed.Items {
		if strings.Contains(v.Title, word) {
			t.Errorf("%v contains %v", v.Title, word)
		}
	}
}

func TestFilterByLink(t *testing.T) {
	feed, err := Get(FeedCategoryAll, nil)
	if err != nil {
		t.Fatal(err)
	}
	link := "anond.hatelabo.jp"
	FilterByLink(&feed.Items, link)
	for _, v := range feed.Items {
		if strings.Contains(v.Link, link) {
			t.Errorf("%v contains %v", v.Link, link)
		}
	}
}

func TestSortByBookmarkCount(t *testing.T) {
	feed, err := Get(FeedCategoryAll, nil)
	if err != nil {
		t.Fatal(err)
	}
	SortByBookmarkCount(feed.Items)
	pre := 0
	for _, v := range feed.Items {
		if pre > v.BookmarkCount {
			t.Fail()
		}
		pre = v.BookmarkCount
	}
}
