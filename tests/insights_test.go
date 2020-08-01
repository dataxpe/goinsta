package tests

import (
	"fmt"
	"testing"
)

func TestInsightsMedia(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	feed := insta.Account.Feed(nil)
	feed.Next()
	if feed.Error() != nil {
		t.Fatalf("%s",feed.Error())
	}
	fmt.Printf("got %d items\n",len(feed.Items))
	for _, item := range feed.Items {
		insights, err := insta.InsightsMedia(item.ID)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("\n\n%s\n",insights.DisplayURL())
		fmt.Printf("%d saves %d\n",item.Pk,insights.SaveCount())
		return
	}
}

func TestInsightsStories(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}

	InsightsStories := insta.InsightsStories()
	for {
		InsightsStories.Next()
		if len(InsightsStories.Errors) > 0 {
			t.Fatal(InsightsStories.Errors[0])
		}

		for i, st := range InsightsStories.Stories() {
			fmt.Printf("%d: %s %d %d\n", i, st.ID, st.ImpressionCount, st.ProfileActionBioLinkClicked())
		}

		if !InsightsStories.HasMore() || len(InsightsStories.Errors) > 0 {
			break
		}
	}
}

func TestInsightsPosts(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}

	InsightsPosts := insta.InsightsPosts()
	for {
		InsightsPosts.Next()
		if len(InsightsPosts.Errors) > 0 {
			t.Fatal(InsightsPosts.Errors[0])
		}

		for i, st := range InsightsPosts.Posts() {
			fmt.Printf("%d: %s %d %d\n", i, st.ID, st.ImpressionCount, st.ProfileActionBioLinkClicked())
		}

		if !InsightsPosts.HasMore() || len(InsightsPosts.Errors) > 0 {
			break
		}
	}
}