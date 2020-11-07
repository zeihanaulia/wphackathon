package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ahmdrz/goinsta/v2"
	"github.com/go-resty/resty/v2"
)

const (
	SHEETYAPI  = "https://api.sheety.co/"
	SHEETYHASH = "92a7688cc908c0eda4eacc8247388e90"
)

var API = fmt.Sprintf("%s%s", SHEETYAPI, SHEETYHASH)

func main() {
	refCodes, err := findAllReferalCode()
	if err != nil {
		log.Fatal(err)
	}

	fullAdsSocial, err := findAllAdsSocial()
	if err != nil {
		log.Fatal(err)
	}

	hashtags, err := findHastagFor(refCodes)
	if err != nil {
		log.Fatal(err)
	}

	inserts, updates, _ := inserOrUpdate(fullAdsSocial, hashtags)

	for _, u := range updates {
		_, _ = updateAdsSocial(u)
	}

	for _, i := range inserts {
		_, _ = insertAdsSocial(i)
	}
}

func inserOrUpdate(fullAdsSocial, hashtags []AdsSocial) (insert, update []AdsSocial, err error) {
	var fasmap = make(map[string]AdsSocial)
	for _, fas := range fullAdsSocial {
		ids := fmt.Sprintf("%s%s", fas.Username, fas.URL)
		if _, ok := fasmap[ids]; !ok {
			fasmap[ids] = fas
		}
	}

	for _, hs := range hashtags {
		ids := fmt.Sprintf("%s%s", hs.Username, hs.URL)
		if v, ok := fasmap[ids]; ok {
			hs.ID = v.ID
			update = append(update, hs)
		} else {
			insert = append(insert, hs)
		}
	}

	return
}

type RequestAdsSocial struct {
	AdsSocial AdsSocial `json:"ads_social,omitempty"`
}

func updateAdsSocial(adsSocial AdsSocial) (refcode []string, err error) {
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(RequestAdsSocial{AdsSocial: adsSocial}).
		Put(fmt.Sprintf("%s/dataKpr/adsSocial/%d", API, adsSocial.ID))
	return
}

func insertAdsSocial(adsSocial AdsSocial) (refcode []string, err error) {
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(RequestAdsSocial{AdsSocial: adsSocial}).
		Post(fmt.Sprintf("%s/dataKpr/adsSocial/", API))
	return
}

type ResponseAdsPartner struct {
	AdsPartner []struct {
		AdsID       int     `json:"adsId,omitempty"`
		AdsName     string  `json:"adsName,omitempty"`
		RefCode     string  `json:"refCode,omitempty"`
		WPCode      string  `json:"wpCode,omitempty"`
		IGLike      int     `json:"igLike,omitempty"`
		IGComment   int     `json:"igComment,omitempty"`
		IGVideoView int     `json:"igVideoView,omitempty"`
		Score       int     `json:"score,omitempty"`
		Amount      float64 `json:"amount,omitempty"`
		State       string  `json:"state,omitempty"`
		ID          int     `json:"id,omitempty"`
	} `json:"adsPartner,omitempty"`
}

func findAllReferalCode() (refcode []string, err error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf("%s/dataKpr/adsPartner/?filter[active]=%d", API, 1))
	if err != nil {
		return
	}

	response := ResponseAdsPartner{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return
	}

	for _, adp := range response.AdsPartner {
		refcode = append(refcode, adp.RefCode)
	}

	return
}

type ResponseAdsSocial struct {
	AdsSocial []AdsSocial `json:"adsSocial,omitempty"`
}

type AdsSocial struct {
	UserID         int64   `json:"userId"`
	Username       string  `json:"userName"`
	AdsPartnerID   string  `json:"adsPartnerId"`
	RefCode        string  `json:"refCode"`
	URL            string  `json:"url"`
	Caption        string  `json:"caption"`
	Comments       int     `json:"comment"`
	Likes          int     `json:"likes"`
	VideoViewCount float64 `json:"videoViewCount"`
	ID             int     `json:"id"`
	Code           string  `json:"code"`
	TakenAt        string  `json:"takenAt"`
}

func findAllAdsSocial() (adsSocial []AdsSocial, err error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf("%s/dataKpr/adsSocial/", API))
	if err != nil {
		err = fmt.Errorf("error get findAllAdsSocial", err)
		return
	}

	response := ResponseAdsSocial{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		err = fmt.Errorf("error unmarshal", err)
		return
	}

	adsSocial = response.AdsSocial
	return
}

func findHastagFor(refcodes []string) (adsSocials []AdsSocial, err error) {
	insta := goinsta.New("hackhatonwarpinanang", "Warpin123!")
	_ = insta.Login()

	for _, rc := range refcodes {
		log.Println("start sync hashtag: ", rc)
		hashtag := insta.NewHashtag(rc)
		_ = hashtag.Sync()

		pages := 100
		for page := 0; page < pages && hashtag.Next(); page++ {
			for _, section := range hashtag.Sections {
				for _, media := range section.LayoutContent.Medias {
					takenAt := time.Unix(media.Item.TakenAt, 0).Format("2006-01-02")
					adsSocials = append(adsSocials, AdsSocial{
						Code:           media.Item.Code,
						RefCode:        rc,
						UserID:         media.Item.User.ID,
						Username:       media.Item.User.Username,
						Likes:          media.Item.Likes,
						URL:            fmt.Sprintf("https://www.instagram.com/p/%s", media.Item.Code),
						Comments:       media.Item.CommentCount,
						VideoViewCount: media.Item.ViewCount,
						Caption:        media.Item.Caption.Text,
						TakenAt:        takenAt,
					})
				}
			}
		}
	}

	defer func() {
		_ = insta.Logout()
	}()

	return
}
