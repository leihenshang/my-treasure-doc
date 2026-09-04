package seed

import (
	"fmt"
	"net"
	"strings"
	"time"

	"fastduck/treasure-doc/module/blog/data/model"
)

type Options struct {
	Enabled        bool
	AllowRemote    bool
	RestoreDeleted bool
	Release        bool
	Host           string
	Database       string
}

type Result struct {
	Created  int
	Existing int
	Restored int
	Skipped  int
}

type contentSeed struct {
	Key         string
	Title       string
	Summary     string
	Content     string
	Category    string
	Status      string
	PublishedOn time.Time
	PublishedAt time.Time
	Pinned      bool
	Tags        []string
	Deleted     bool
}

func Validate(options Options) error {
	if !options.Enabled {
		return nil
	}
	if options.Release {
		return fmt.Errorf("blog seed is forbidden in release mode")
	}
	if !options.AllowRemote && !isLocalHost(options.Host) {
		return fmt.Errorf("blog seed for remote database %s/%s requires allowRemote=true", options.Host, options.Database)
	}
	return nil
}

func isLocalHost(host string) bool {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "localhost" || host == "::1" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func postSeeds() []contentSeed {
	return contentSeeds("mock-post", "Mock 文章", []string{"tech", "essay", "reading"}, false)
}

func diarySeeds() []contentSeed {
	items := contentSeeds("mock-diary", "Mock 日记", nil, true)
	items[0].Key = "a-day-without-date"
	return items
}

func contentSeeds(prefix, title string, categories []string, diary bool) []contentSeed {
	items := make([]contentSeed, 0, 20)
	for index := 1; index <= 20; index++ {
		status, publishedAt := contentState(index)
		category := ""
		if len(categories) > 0 {
			category = categories[(index-1)%len(categories)]
		}
		tags := []string{"Go"}
		if diary {
			tags = []string{"复盘"}
		}
		if index%2 == 0 {
			tags = append(tags, "Vue")
		}
		if index == 13 {
			tags = []string{"ScheduledOnly"}
		}
		if index >= 15 && index <= 17 {
			tags = []string{"DraftOnly"}
		}
		if index >= 18 {
			tags = []string{"ArchivedOnly"}
		}
		items = append(items, contentSeed{
			Key: fmt.Sprintf("%s-%02d", prefix, index), Title: fmt.Sprintf("%s %02d", title, index),
			Summary: fmt.Sprintf("第 %02d 条固定示例数据。", index), Content: fmt.Sprintf("# %s %02d\n\n包含 Go、Vue、生活和工程实践关键词。", title, index),
			Category: category, Status: status, PublishedOn: dateFor(index), PublishedAt: publishedAt,
			Pinned: index <= 4, Tags: tags, Deleted: index == 14,
		})
	}
	return items
}

func contentState(index int) (string, time.Time) {
	if index <= 12 {
		return model.StatusPublished, pastTime(index)
	}
	if index == 13 {
		return model.StatusPublished, time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if index == 14 {
		return model.StatusPublished, pastTime(index)
	}
	if index <= 17 {
		return model.StatusDraft, time.Time{}
	}
	return model.StatusArchived, pastTime(index)
}

func otherState(index int) (string, time.Time) {
	if index <= 4 {
		return model.StatusPublished, pastTime(index)
	}
	if index == 5 {
		return model.StatusPublished, time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if index == 6 {
		return model.StatusDraft, time.Time{}
	}
	if index == 7 {
		return model.StatusArchived, pastTime(index)
	}
	return model.StatusPublished, pastTime(index)
}

func pastTime(index int) time.Time {
	return time.Date(2025, time.Month((index-1)%12+1), min(index, 28), 8, 0, 0, 0, time.UTC)
}
func dateFor(index int) time.Time {
	return time.Date(2025, time.Month((index-1)%12+1), min(index, 28), 0, 0, 0, 0, time.UTC)
}
