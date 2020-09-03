package tests

import (
	"fmt"
	"testing"
)

func TestStories(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}

	stories := insta.Account.Stories()
	stories.Next()
	if stories.Error() != nil {
		//t.Fatalf("%s",stories.Error())
	}
	fmt.Printf("got %d items\n",len(stories.Stories))


	st, err := insta.Timeline.Stories()
	if err != nil {
		fmt.Printf("err: %s\n",err)
	}
	fmt.Printf("got %d items\n",len(st.Stories))
}
