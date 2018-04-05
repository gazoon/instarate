package instagram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	apiUrl      = "https://www.instagram.com/"
	mediaPath   = "p/"
	magicSuffix = "/?__a=1"
)

var (
	MediaForbidden = errors.New("it's a private account or the media doesn't exist")
)

var (
	httpClient = http.Client{Timeout: time.Second * 3}
)

type Media struct {
	Owner   string
	Url     string
	IsPhoto bool
}

type mediaResponse struct {
	GraphQL struct {
		Media struct {
			Owner struct {
				Username string
			}
			DisplayResources []struct {
				Src string
			} `json:"display_resources"`
			IsVideo *bool `json:"is_video"`
			Sidecar struct {
				Edges []struct {
					Node struct {
						IsVideo *bool `json:"is_video"`
					}
				}
			} `json:"edge_sidecar_to_children"`
		} `json:"shortcode_media"`
	} `json:"graphql"`
}

func GetMediaInfo(ctx context.Context, mediaCode string) (*Media, error) {
	mediaUrl := apiUrl + mediaPath + mediaCode + magicSuffix
	resp, err := httpClient.Get(mediaUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, MediaForbidden
	}
	responseData := &mediaResponse{}
	err = json.NewDecoder(resp.Body).Decode(responseData)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse media response")
	}
	err = validateMediaResponse(responseData)
	if err != nil {
		return nil, errors.Wrap(err, "invalid media response")
	}
	mediaData := responseData.GraphQL.Media
	username := mediaData.Owner.Username
	displayUrl := mediaData.DisplayResources[0].Src
	isPhoto := true
	if *mediaData.IsVideo {
		isPhoto = false
	} else {
		edges := mediaData.Sidecar.Edges
		if len(edges) != 0 {
			if *edges[0].Node.IsVideo {
				isPhoto = false
			}
		}
	}

	return &Media{Owner: username, Url: displayUrl, IsPhoto: isPhoto}, nil
}

func validateMediaResponse(data *mediaResponse) error {
	media := data.GraphQL.Media
	if media.Owner.Username == "" {
		return errors.New("no username")
	}
	if len(media.DisplayResources) == 0 {
		return errors.New("no display resources")
	}
	if media.DisplayResources[0].Src == "" {
		return errors.New("no display resource src")
	}
	if media.IsVideo == nil {
		return errors.New("no is video")
	}
	if len(media.Sidecar.Edges) != 0 && media.Sidecar.Edges[0].Node.IsVideo == nil {
		return errors.New("no sidecar is video")
	}
	return nil
}

type userResponse struct {
	User struct {
		FollowedBy struct {
			Count *int
		} `json:"followed_by"`
	}
}

func GetFollowersNumber(ctx context.Context, username string) (int, error) {
	profileUrl := apiUrl + username + magicSuffix
	resp, err := httpClient.Get(profileUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	responseData := &userResponse{}
	err = json.NewDecoder(resp.Body).Decode(responseData)
	if err != nil {
		return 0, errors.Wrap(err, "can't parse user response")
	}
	if responseData.User.FollowedBy.Count == nil {
		return 0, errors.New("invalid user response: no followers info")
	}
	return *responseData.User.FollowedBy.Count, nil
}

func BuildProfileUrl(username string) string {
	return apiUrl + username + "/"
}

func ExtractUsername(profileUrl string) (string, error) {
	username, err := extractLastPathPart(profileUrl)
	return username, errors.Wrap(err, "cant extract username")
}

func ExtractMediaCode(mediaUrl string) (string, error) {
	mediaCode, err := extractLastPathPart(mediaUrl)
	return mediaCode, errors.Wrap(err, "cant extract media code")
}

func extractLastPathPart(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	urlPath := strings.TrimSuffix(u.Path, "/")
	_, lastPart := path.Split(urlPath)
	if lastPart == "" {
		return "", errors.New("url has empty path")
	}
	return lastPart, nil
}