package tests

import (
	"fmt"
	"testing"
)

func TestInsights(t *testing.T) {
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