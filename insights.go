package goinsta

import (
	"encoding/json"
	"fmt"
)

type Insights struct {
	Data struct {
		Media struct {
			ID                    string `json:"id"`
			CreationTime          int    `json:"creation_time"`
			HasProductTags        bool   `json:"has_product_tags"`
			InstagramMediaID      string `json:"instagram_media_id"`
			InstagramMediaOwnerID string `json:"instagram_media_owner_id"`
			InstagramActor        struct {
				InstagramActorID string `json:"instagram_actor_id"`
				ID               string `json:"id"`
			} `json:"instagram_actor"`
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
			} `json:"image"`
			CommentCount               int         `json:"comment_count"`
			LikeCount                  int         `json:"like_count"`
			SaveCount                  int         `json:"save_count"`
			AdMedia                    interface{} `json:"ad_media"`
			OrganicInstagramMediaID    string      `json:"organic_instagram_media_id"`
			ShoppingOutboundClickCount int         `json:"shopping_outbound_click_count"`
			ShoppingProductClickCount  int         `json:"shopping_product_click_count"`
		} `json:"media"`
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

func (ins *Insights) ProfileActionBioLinkClicked() int {
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "BIO_LINK_CLICKED" {
			return a.Value
		}
	}
	return 0
}

func (ins *Insights) ProfileActionCalled() int {
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "CALL" {
			return a.Value
		}
	}
	return 0
}

func (ins *Insights) ProfileActionDirection() int {
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "DIRECTION" {
			return a.Value
		}
	}
	return 0
}

func (ins *Insights) ProfileActionEmail() int {
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "EMAIL" {
			return a.Value
		}
	}
	return 0
}

func (ins *Insights) ProfileActionText() int {
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.ProfileActions.Actions.Nodes {
		if a.Name == "TEXT" {
			return a.Value
		}
	}
	return 0
}

func (ins *Insights) ProfileViews() int {
	return ins.Data.Media.InlineInsightsNode.Metrics.OwnerProfileViewsCount
}

func (ins *Insights) Reached() (reachCount int, nonFollowers int) {
	reachCount = ins.Data.Media.InlineInsightsNode.Metrics.Reach.Value
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.Reach.FollowStatus.Nodes {
		if a.Name == "NON_FOLLOWER" {
			nonFollowers = a.Value
		}
	}
	return
}

func (ins *Insights) Shares() int {
	return ins.Data.Media.InlineInsightsNode.Metrics.ShareCount.Post.Value
}

func (ins *Insights) Follows() int {
	return ins.Data.Media.InlineInsightsNode.Metrics.OwnerAccountFollowsCount
}

func (ins *Insights) ImpressionsTotal() int {
	return ins.Data.Media.InlineInsightsNode.Metrics.Impressions.Value
}

func (ins *Insights) ImpressionsFromFeed() int {
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.Impressions.Surfaces.Nodes {
		if a.Name == "FEED" {
			return a.Value
		}
	}
	return 0
}

func (ins *Insights) ImpressionsFromProfile() int {
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.Impressions.Surfaces.Nodes {
		if a.Name == "PROFILE" {
			return a.Value
		}
	}
	return 0
}

func (ins *Insights) ImpressionsFromHashtags() int {
	for _, a := range ins.Data.Media.InlineInsightsNode.Metrics.Impressions.Surfaces.Nodes {
		if a.Name == "HASHTAG" {
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

