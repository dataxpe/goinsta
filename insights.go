package goinsta

import (
	"encoding/json"
	"fmt"
)

type Insights struct {
	Data struct {
		Media InsightsMedia `json:"media"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	}
}

func (ins *Insights) LikesCount() int {
	return ins.Data.Media.LikeCount
}

func (ins *Insights) CommentCount() int {
	return ins.Data.Media.CommentCount
}

func (ins *Insights) SaveCount() int {
	return ins.Data.Media.SaveCount
}

func (ins *Insights) DisplayURL() string {
	return ins.Data.Media.DisplayURL
}

func (ins *Insights) IsActive() bool {
	return ins.Data.Media.IsActive()
}

func (ins *Insights) ProfileActionBioLinkClicked() int {
	return ins.Data.Media.ProfileActionBioLinkClicked()
}

func (ins *Insights) ProfileActionCalled() int {
	return ins.Data.Media.ProfileActionCalled()
}

func (ins *Insights) ProfileActionDirection() int {
	return ins.Data.Media.ProfileActionDirection()
}

func (ins *Insights) ProfileActionEmail() int {
	return ins.Data.Media.ProfileActionEmail()
}

func (ins *Insights) ProfileActionText() int {
	return ins.Data.Media.ProfileActionText()
}

func (ins *Insights) ProfileViews() int {
	return ins.Data.Media.ProfileViews()
}

func (ins *Insights) Reached() (reachCount int, nonFollowers int) {
	return ins.Data.Media.Reached()
}

func (ins *Insights) Shares() (shares int) {
	return ins.Data.Media.Shares()
}

func (ins *Insights) Follows() int {
	return ins.Data.Media.Follows()
}

func (ins *Insights) ImpressionsTotal() int {
	return ins.Data.Media.ImpressionsTotal()
}

func (ins *Insights) ImpressionsFromFeed() int {
	return ins.Data.Media.ImpressionsFromFeed()
}

func (ins *Insights) ImpressionsFromProfile() int {
	return ins.Data.Media.ImpressionsFromProfile()
}

func (ins *Insights) ImpressionsFromHashtags() int {
	return ins.Data.Media.ImpressionsFromHashtags()
}

func (ins *Insights) ImpressionsFromExplore() int {
	return ins.Data.Media.ImpressionsFromExplore()
}

func (ins *Insights) ImpressionsFromLocation() int {
	return ins.Data.Media.ImpressionsFromLocation()
}



type InsightsMedia struct {
	ID                    string `json:"id"`
	CreationTime          int    `json:"creation_time"`
	HasProductTags        bool   `json:"has_product_tags"`
	InstagramMediaID      string `json:"instagram_media_id"`
	InstagramMediaOwnerID string `json:"instagram_media_owner_id"`
	InlineInsightsNode struct {
		State   string `json:"state"`
		Metrics struct {
			ShareCount struct {
				Tray struct {
					Nodes []struct {
						Typename string `json:"__typename"`
						Value    int    `json:"value"`
					} `json:"nodes"`
				} `json:"tray"`
				Post struct {
					Value int `json:"value"`
					Nodes []struct {
						Typename string `json:"__typename"`
						Name     string `json:"name"`
						Value    int    `json:"value"`
					} `json:"nodes"`
				} `json:"post"`
				Share struct {
					Value int `json:"value"`
				} `json:"share"`
			} `json:"share_count"`
			OwnerProfileViewsCount int `json:"owner_profile_views_count"`
			ReachCount             int `json:"reach_count"`
			ProfileActions         struct {
				Actions struct {
					Value int `json:"value"`
					Nodes []struct {
						Typename string `json:"__typename"`
						Name     string `json:"name"`
						Value    int    `json:"value"`
					} `json:"nodes"`
					Edges []struct {
						Node struct {
							Typename string `json:"__typename"`
							Name     string `json:"name"`
							Value    int    `json:"value"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"actions"`
			} `json:"profile_actions"`
			ImpressionCount int `json:"impression_count"`
			Impressions     struct {
				Value    int `json:"value"`
				Surfaces struct {
					Nodes []struct {
						Typename string `json:"__typename"`
						Name     string `json:"name"`
						Value    int    `json:"value"`
					} `json:"nodes"`
				} `json:"surfaces"`
			} `json:"impressions"`
			OwnerAccountFollowsCount int `json:"owner_account_follows_count"`
			StoryLinkNavigationCount int `json:"story_link_navigation_count"`
			Reach                    struct {
				Value        int `json:"value"`
				FollowStatus struct {
					Nodes []struct {
						Typename string `json:"__typename"`
						Name     string `json:"name"`
						Value    int    `json:"value"`
					} `json:"nodes"`
				} `json:"follow_status"`
			} `json:"reach"`
			HashtagsImpressions interface{} `json:"hashtags_impressions"`
		} `json:"metrics"`
		Error interface{} `json:"error"`
	} `json:"inline_insights_node"`
	DisplayURL         string `json:"display_url"`
	InstagramMediaType string `json:"instagram_media_type"`
	Image              struct {
		Height int `json:"height"`
		Width  int `json:"width"`
		URI string `json:"uri"`
	} `json:"image"`
	CommentCount               int         `json:"comment_count"`
	LikeCount                  int         `json:"like_count"`
	SaveCount                  int         `json:"save_count"`
	Engagement                 int         `json:"engagement"`
	//AdMedia                    interface{} `json:"ad_media"`
	OrganicInstagramMediaID    string      `json:"organic_instagram_media_id"`
	ShoppingOutboundClickCount int         `json:"shopping_outbound_click_count"`
	ShoppingProductClickCount  int         `json:"shopping_product_click_count"`
	ExitsCount          	   int    	   `json:"exits_count"`
	ImpressionCount     	   int    	   `json:"impression_count"`
	ReachCount          	   int    	   `json:"reach_count"`
	RepliesCount        	   int    	   `json:"replies_count"`
	TapsBackCount       	   int    	   `json:"taps_back_count"`
	TapsForwardCount    	   int    	   `json:"taps_forward_count"`
	StorySwipeAwayCount 	   int    	   `json:"story_swipe_away_count"`
}

func (ins *InsightsMedia) IsActive() bool {
	return ins.InlineInsightsNode.State == "AVAILABLE"
}

func (ins *InsightsMedia) ProfileActionBioLinkClicked() int {
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "BIO_LINK_CLICKED" {
			return a.Value
		}
	}
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Edges {
		if a.Node.Name == "BIO_LINK_CLICKED" {
			return a.Node.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ProfileActionCalled() int {
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "CALL" {
			return a.Value
		}
	}
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Edges {
		if a.Node.Name == "CALL" {
			return a.Node.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ProfileActionDirection() int {
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "DIRECTION" {
			return a.Value
		}
	}
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Edges {
		if a.Node.Name == "DIRECTION" {
			return a.Node.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ProfileActionEmail() int {
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "EMAIL" {
			return a.Value
		}
	}
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Edges {
		if a.Node.Name == "EMAIL" {
			return a.Node.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ProfileActionText() int {
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "TEXT" {
			return a.Value
		}
	}
	for _, a := range ins.InlineInsightsNode.Metrics.ProfileActions.Actions.Edges {
		if a.Node.Name == "TEXT" {
			return a.Node.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ProfileViews() int {
	return ins.InlineInsightsNode.Metrics.OwnerProfileViewsCount
}

func (ins *InsightsMedia) Reached() (reachCount int, nonFollowers int) {
	reachCount = ins.InlineInsightsNode.Metrics.Reach.Value
	for _, a := range ins.InlineInsightsNode.Metrics.Reach.FollowStatus.Nodes {
		if a.Name == "NON_FOLLOWER" {
			nonFollowers = a.Value
		}
	}
	return
}

func (ins *InsightsMedia) Shares() (shares int) {
	shares += ins.InlineInsightsNode.Metrics.ShareCount.Post.Value // posts
	shares += ins.InlineInsightsNode.Metrics.ShareCount.Share.Value // stories
	for _, c := range ins.InlineInsightsNode.Metrics.ShareCount.Tray.Nodes {
		shares += c.Value // tray
	}
	return
}

func (ins *InsightsMedia) Follows() int {
	return ins.InlineInsightsNode.Metrics.OwnerAccountFollowsCount
}

func (ins *InsightsMedia) ImpressionsTotal() int {
	return ins.InlineInsightsNode.Metrics.Impressions.Value
}

func (ins *InsightsMedia) ImpressionsFromFeed() int {
	for _, a := range ins.InlineInsightsNode.Metrics.Impressions.Surfaces.Nodes {
		if a.Name == "FEED" {
			return a.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ImpressionsFromProfile() int {
	for _, a := range ins.InlineInsightsNode.Metrics.Impressions.Surfaces.Nodes {
		if a.Name == "PROFILE" {
			return a.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ImpressionsFromHashtags() int {
	for _, a := range ins.InlineInsightsNode.Metrics.Impressions.Surfaces.Nodes {
		if a.Name == "HASHTAG" {
			return a.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ImpressionsFromExplore() int {
	for _, a := range ins.InlineInsightsNode.Metrics.Impressions.Surfaces.Nodes {
		if a.Name == "EXPLORE" {
			return a.Value
		}
	}
	return 0
}

func (ins *InsightsMedia) ImpressionsFromLocation() int {
	for _, a := range ins.InlineInsightsNode.Metrics.Impressions.Surfaces.Nodes {
		if a.Name == "LOCATION" {
			return a.Value
		}
	}
	return 0
}

// variables came from...
// https://github.com/dilame/instagram-private-api/blob/b5cf84910d0b91151189cdc28e7c12eb8d725f31/src/services/insights.service.ts
// https://github.com/dilame/instagram-private-api/blob/dcceea8dc2f3bb9f4ae4daa31cd5110264fb1ef2/src/types/insights.options.ts
// https://github.com/dilame/instagram-private-api/blob/dcceea8dc2f3bb9f4ae4daa31cd5110264fb1ef2/src/repositories/ads.repository.ts
func (insta *Instagram) InsightsMedia(MediaID string) (ins Insights, err error) {
	variables := map[string]interface{}{
		"query_params": map[string]interface{}{
			"access_token": "",
			"id":      fmt.Sprintf("%s", MediaID),
		},
	}
	vjson, err := json.Marshal(variables)
	if err != nil {
		return
	}

	data := map[string]string{
		"fb_api_caller_class": "RelayModern",
		"fb_api_req_friendly_name": "IgInsightsPostInsightsQuery",
		"variables": fmt.Sprintf("%s",vjson),
		"doc_id": "2009845309144121", // specific hardcoded id
	}

	body, err := insta.sendRequest(
		&reqOptions{
			Endpoint: 	"ads/graphql/?locale=en_US&vc_policy=insights_policy&surface=post",
			Query: 		data,
			IsPost:   	true,
		},
	)
	if err != nil {
		return
	}

	err = json.Unmarshal(body,&ins)
	if err != nil {
		return
	}

	if len(ins.Errors) > 0 {
		err = fmt.Errorf("%s",ins.Errors[0].Message)
	}

	return
}





type InsightsPosts struct {
	inst *Instagram
	Data struct {
		User struct {
			ID string `json:"id"`
			BusinessManager struct {
				TopPostsUnit struct {
					TopPosts struct {
						Edges []struct {
							Node InsightsMedia `json:"node"`
						} `json:"edges"`
						PageInfo struct {
							EndCursor   string `json:"end_cursor"`
							HasNextPage bool   `json:"has_next_page"`
						} `json:"page_info"`
					} `json:"top_posts"`
				} `json:"top_posts_unit"`
			} `json:"business_manager"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	}
}

func (insta *Instagram) InsightsPosts() (is InsightsPosts) {
	is.inst = insta
	return
}

func (ip *InsightsPosts) Posts() (Stories []InsightsMedia) {
	for _, s := range ip.Data.User.BusinessManager.TopPostsUnit.TopPosts.Edges {
		Stories = append(Stories,s.Node)
	}
	return
}

// https://github.com/dilame/instagram-private-api/blob/dcceea8dc2f3bb9f4ae4daa31cd5110264fb1ef2/src/feeds/posts-insights.feed.ts
func (ip *InsightsPosts) Next() bool {
	variables := map[string]interface{}{
		"count": 15,
		"cursor": ip.Data.User.BusinessManager.TopPostsUnit.TopPosts.PageInfo.EndCursor,
		"IgInsightsGridMediaImage_SIZE": 256,
		"timeframe": "TWO_WEEKS", // ONE_DAY | ONE_WEEK | TWO_WEEKS
		"queryParams": map[string]interface{}{
			"access_token": "",
			"id":      fmt.Sprintf("%d", ip.inst.Account.ID),
		},
	}
	vjson, err := json.Marshal(variables)
	if err != nil {
		return false
	}

	data := map[string]string{
		"fb_api_caller_class": "RelayModern",
		"fb_api_req_friendly_name": "IgInsightsPostGridSurfaceQuery",
		"variables": fmt.Sprintf("%s",vjson),
		"doc_id": "1981884911894608", // specific hardcoded id
	}

	body, err := ip.inst.sendRequest(
		&reqOptions{
			Endpoint: 	"ads/graphql/?locale=en_US&vc_policy=insights_policy",
			Query: 		data,
			IsPost:   	true,
		},
	)
	if err != nil {
		return false
	}

	err = json.Unmarshal(body,&ip)
	if err != nil {
		return false
	}

	return ip.Data.User.BusinessManager.TopPostsUnit.TopPosts.PageInfo.HasNextPage
}

func (ip *InsightsPosts) HasMore() bool {
	return ip.Data.User.BusinessManager.TopPostsUnit.TopPosts.PageInfo.HasNextPage
}











type InsightsStories struct {
	inst *Instagram
	Data struct {
		User struct {
			ID string `json:"id"`
			BusinessManager struct {
				StoriesUnit struct {
					Stories struct {
						Edges []struct {
							Node InsightsMedia `json:"node"`
						} `json:"edges"`
						PageInfo struct {
							EndCursor   string `json:"end_cursor"`
							HasNextPage bool   `json:"has_next_page"`
						} `json:"page_info"`
					} `json:"stories"`
				} `json:"stories_unit"`
			} `json:"business_manager"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	}
}

func (insta *Instagram) InsightsStories() (is InsightsStories) {
	is.inst = insta
	return
}

func (is *InsightsStories) Stories() (Stories []InsightsMedia) {
	for _, s := range is.Data.User.BusinessManager.StoriesUnit.Stories.Edges {
		Stories = append(Stories,s.Node)
	}
	return
}

// https://github.com/dilame/instagram-private-api/blob/dcceea8dc2f3bb9f4ae4daa31cd5110264fb1ef2/src/feeds/stories-insights.feed.ts
func (is *InsightsStories) Next() bool {
	variables := map[string]interface{}{
		"count": 15,
		"cursor": is.Data.User.BusinessManager.StoriesUnit.Stories.PageInfo.EndCursor,
		"IgInsightsGridMediaImage_SIZE": 256,
		"timeframe": "TWO_WEEKS", // ONE_DAY | ONE_WEEK | TWO_WEEKS
		"queryParams": map[string]interface{}{
			"access_token": "",
			"id":      fmt.Sprintf("%d", is.inst.Account.ID),
		},
	}
	vjson, err := json.Marshal(variables)
	if err != nil {
		return false
	}

	data := map[string]string{
		"fb_api_caller_class": "RelayModern",
		"fb_api_req_friendly_name": "IgInsightsStoryGridSurfaceQuery",
		"variables": fmt.Sprintf("%s",vjson),
		"doc_id": "1995528257207653", // specific hardcoded id
	}

	body, err := is.inst.sendRequest(
		&reqOptions{
			Endpoint: 	"ads/graphql/?locale=en_US&vc_policy=insights_policy",
			Query: 		data,
			IsPost:   	true,
		},
	)
	if err != nil {
		return false
	}

	err = json.Unmarshal(body,&is)
	if err != nil {
		return false
	}

	return is.Data.User.BusinessManager.StoriesUnit.Stories.PageInfo.HasNextPage
}

func (is *InsightsStories) HasMore() bool {
	return is.Data.User.BusinessManager.StoriesUnit.Stories.PageInfo.HasNextPage
}

