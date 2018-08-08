package instagram

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"io/ioutil"

	"github.com/pkg/errors"
)

const (
	apiUrl      = "https://www.instagram.com/"
	mediaPath   = "p/"
	httpTimeout = time.Second * 3
)

var (
	MediaForbidden = errors.New("it's a private account or the media doesn't exist")
)

var (
	httpClient = http.Client{Timeout: httpTimeout}
)

type Media struct {
	Owner   string
	Url     string
	IsPhoto bool
}

type mediaResponse struct {
	EntryData struct {
		PostPage []struct {
			GraphQL struct {
				Media struct {
					Owner struct {
						Username string `json:"username"`
					} `json:"owner"`
					DisplayResources []struct {
						Src string `json:"src"`
					} `json:"display_resources"`
					IsVideo *bool `json:"is_video"`
					Sidecar struct {
						Edges []struct {
							Node struct {
								IsVideo *bool `json:"is_video"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_sidecar_to_children"`
				} `json:"shortcode_media"`
			} `json:"graphql"`
		} `json:"PostPage"`
	} `json:"entry_data"`
}

func GetMediaInfo(ctx context.Context, mediaCode string) (*Media, error) {
	mediaUrl := apiUrl + mediaPath + mediaCode
	jsonContent, httpStatus, err := requestResourceDataWithStatusCode(mediaUrl)
	if err != nil {
		if httpStatus == http.StatusForbidden {
			return nil, MediaForbidden
		}
		return nil, err
	}
	responseData := &mediaResponse{}
	err = json.Unmarshal([]byte(jsonContent), responseData)
	if err != nil {
		return nil, errors.Wrap(err, "parse instagram media response into json")
	}
	err = validateMediaResponse(responseData)
	if err != nil {
		return nil, errors.Wrapf(err, "media_url=%s; json_content=%s; invalid media response", mediaUrl, jsonContent)
	}
	mediaData := responseData.EntryData.PostPage[0].GraphQL.Media
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
	if len(data.EntryData.PostPage) != 1 {
		return errors.Errorf("expected only one post, got: %d", len(data.EntryData.PostPage))
	}
	media := data.EntryData.PostPage[0].GraphQL.Media
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
	EntryData struct {
		ProfilePage []struct {
			GraphQL struct {
				User struct {
					FollowedBy struct {
						Count *int `json:"count"`
					} `json:"edge_followed_by"`
				} `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}

func GetFollowersNumber(ctx context.Context, username string) (int, error) {
	profileUrl := apiUrl + username
	jsonContent, err := RequestResourceData(profileUrl)
	if err != nil {
		return 0, err
	}
	responseData := &userResponse{}
	err = json.Unmarshal([]byte(jsonContent), responseData)
	if err != nil {
		return 0, errors.Wrap(err, "parse instagram profile response into json")
	}
	err = validateUserResponse(responseData)
	if err != nil {
		return 0, errors.Wrapf(err, "profile_url=%s; json_content=%s; invalid user response", profileUrl, jsonContent)
	}
	return *responseData.EntryData.ProfilePage[0].GraphQL.User.FollowedBy.Count, nil
}

func RequestResourceData(url string) (string, error) {
	data, _, err := requestResourceDataWithStatusCode(url)
	return data, err
}

func requestResourceDataWithStatusCode(url string) (string, int, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", 0, errors.Wrapf(err, "url=%s; http get instagram data", url)
	}
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.Errorf("url=%s; instagram API unexpected status: %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, errors.Wrapf(err, "url=%s; read instagram http response content", url)
	}
	text := string(b)
	firstMarker := "window._sharedData = "
	secondMarker := ";</script>"
	dataStart := strings.Index(text, firstMarker) + len(firstMarker)
	dataEnd := strings.Index(text, secondMarker)
	data := text[dataStart:dataEnd]
	return data, resp.StatusCode, nil
}

func validateUserResponse(data *userResponse) error {
	if len(data.EntryData.ProfilePage) != 1 {
		return errors.Errorf("expected only one profile, got: %d", len(data.EntryData.ProfilePage))
	}
	if data.EntryData.ProfilePage[0].GraphQL.User.FollowedBy.Count == nil {
		return errors.New("no followers info")
	}
	return nil
}

func BuildProfileUrl(username string) string {
	return apiUrl + username + "/"
}
