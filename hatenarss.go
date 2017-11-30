package hatenarss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

// FeedCategory は RSS フィードのカテゴリを表します。
type FeedCategory int

const (
	// FeedCategoryAll はホットエントリー総合を表します。
	FeedCategoryAll FeedCategory = iota
	// FeedCategorySocial はホットエントリー「世の中」を表します。
	FeedCategorySocial
	// FeedCategoryEconomics はホットエントリー「政治と経済」を表します。
	FeedCategoryEconomics
	// FeedCategoryLife はホットエントリー「暮らし」を表します。
	FeedCategoryLife
	// FeedCategoryKnowledge はホットエントリー「学び」を表します。
	FeedCategoryKnowledge
	// FeedCategoryIt はホットエントリー「テクノロジー」を表します。
	FeedCategoryIt
	// FeedCategoryEntertainment はホットエントリー「エンタメ」
	FeedCategoryEntertainment
	// FeedCategoryGame はホットエントリー「アニメとゲーム」を表します。
	FeedCategoryGame
	// FeedCategoryFun はホットエントリー「おもしろ」を表します。
	FeedCategoryFun
	// FeedCategoryVideo はホットエントリー「動画」を表します。
	FeedCategoryVideo
)

// GetFeedCategoryList は定義された全ての FeedCategory を返します。
func GetFeedCategoryList() []FeedCategory {
	return []FeedCategory{
		FeedCategoryAll,
		FeedCategorySocial,
		FeedCategoryEconomics,
		FeedCategoryLife,
		FeedCategoryKnowledge,
		FeedCategoryIt,
		FeedCategoryEntertainment,
		FeedCategoryGame,
		FeedCategoryFun,
		FeedCategoryVideo,
	}
}

// Title は RSS フィードの名前を返します。
func (i FeedCategory) Title() string {
	switch i {
	case FeedCategoryAll:
		return "総合"
	case FeedCategorySocial:
		return "世の中"
	case FeedCategoryEconomics:
		return "政治と経済"
	case FeedCategoryLife:
		return "暮らし"
	case FeedCategoryKnowledge:
		return "学び"
	case FeedCategoryIt:
		return "テクノロジー"
	case FeedCategoryEntertainment:
		return "エンタメ"
	case FeedCategoryGame:
		return "アニメとゲーム"
	case FeedCategoryFun:
		return "おもしろ"
	case FeedCategoryVideo:
		return "動画"
	default:
		log.Fatal("Unsupported HatenaCategory")
		return ""
	}
}

// URL は RSS フィードの URL を返します。
func (i FeedCategory) URL() string {
	switch i {
	case FeedCategoryAll:
		return "http://b.hatena.ne.jp/hotentry.rss"
	case FeedCategorySocial:
		return "http://b.hatena.ne.jp/hotentry/social.rss"
	case FeedCategoryEconomics:
		return "http://b.hatena.ne.jp/hotentry/economics.rss"
	case FeedCategoryLife:
		return "http://b.hatena.ne.jp/hotentry/life.rss"
	case FeedCategoryKnowledge:
		return "http://b.hatena.ne.jp/hotentry/knowledge.rss"
	case FeedCategoryIt:
		return "http://b.hatena.ne.jp/hotentry/it.rss"
	case FeedCategoryEntertainment:
		return "http://b.hatena.ne.jp/hotentry/entertainment.rss"
	case FeedCategoryGame:
		return "http://b.hatena.ne.jp/hotentry/game.rss"
	case FeedCategoryFun:
		return "http://b.hatena.ne.jp/hotentry/fun.rss"
	case FeedCategoryVideo:
		return "http://b.hatena.ne.jp/video.rss"
	default:
		log.Fatal("Unsupported HatenaCategory")
		return ""
	}
}

// Feed は、はてなの RSS フィード全体の情報を表す構造体です。
type Feed struct {
	Channel Channel `xml:"channel"`
	Items   []Item  `xml:"item"`
}

func (feed Feed) String() string {
	buf := bytes.Buffer{}
	for i, item := range feed.Items {
		buf.WriteString(fmt.Sprintf("[%d]\n%s\n", i, item))
	}
	return buf.String()
}

// Channel は、はてなの RSS フィードのチャンネル情報を表す構造体です。
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

// Item は はてなの RSS フィードの含まれた個別の記事情報を表す構造体です。
type Item struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Content       string    `xml:"encoded"`
	Date          time.Time `xml:"date"`
	Subject       string    `xml:"subject"`
	BookmarkCount int       `xml:"bookmarkcount"`
}

func (item Item) String() string {
	return fmt.Sprintf(`Title 	        : %s
Link            : %s
Description     : %s
Content         : %s
Date            : %s
Subject         : %s
BookmarkCount   : %d
`, item.Title, item.Link, item.Description, item.Content, item.Date, item.Subject, item.BookmarkCount)
}

// Get は指定したカテゴリの RSS フィードを返します。
func Get(category FeedCategory, client *http.Client) (*Feed, error) {
	if client == nil {
		client = &http.Client{}
	}
	resp, err := client.Get(category.URL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resp.StatusCode != http.StatusOK : %v", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var feed Feed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}
	return &feed, nil
}

// Filter  は入力の items から条件に一致する項目を取り除きます。
func Filter(items *[]Item, reject func(item Item) bool) {
	filteredItem := (*items)[:0]
	for _, item := range *items {
		if !reject(item) {
			filteredItem = append(filteredItem, item)
		}
	}
	*items = filteredItem
}

// FilterByBookmarkCount は入力の items からはてなブックマークの数が threshold 未満の項目を取り除きます。
func FilterByBookmarkCount(items *[]Item, threshold int) {
	Filter(items, func(item Item) bool {
		return item.BookmarkCount < threshold
	})
}

// FilterByTitle は入力の items から記事のタイトルに patterns が含まれる項目を取り除きます。
func FilterByTitle(items *[]Item, patterns ...string) {
	Filter(items, func(item Item) bool {
		for _, p := range patterns {
			if strings.Contains(item.Title, p) {
				return true
			}
		}
		return false
	})
}

// FilterByLink は入力の items から記事のリンクに patterns が含まれる項目を取り除きます。
func FilterByLink(items *[]Item, patterns ...string) {
	Filter(items, func(item Item) bool {
		for _, p := range patterns {
			if strings.Contains(item.Link, p) {
				return true
			}
		}
		return false
	})
}

type hatenaItems []Item

func (items hatenaItems) Len() int {
	return len(items)
}

func (items hatenaItems) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items hatenaItems) Less(i, j int) bool {
	return items[i].BookmarkCount < items[j].BookmarkCount
}

// SortByBookmarkCount は入力の items をブックマーク数で昇順に並び替えます。
func SortByBookmarkCount(items []Item) {
	var h hatenaItems = items
	sort.Sort(h)
}
