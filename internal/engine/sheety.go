package engine

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/zeihanaulia/instagram-scraper/internal/common"
)

type Sheety struct {
	API string
}

func NewSheety(api string) Sheety {
	return Sheety{API: api}
}

type RequestAdsSocial struct {
	AdsSocial AdsSocial `json:"ads_social,omitempty"`
}

func (s *Sheety) UpdateAdsSocial(adsSocial AdsSocial) (refcode []string, err error) {
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(RequestAdsSocial{AdsSocial: adsSocial}).
		Put(fmt.Sprintf("%s%s%d", s.API, common.SHEETY_ADSSOCIAL, adsSocial.ID))
	return
}

func (s *Sheety) InsertAdsSocial(adsSocial AdsSocial) (refcode []string, err error) {
	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(RequestAdsSocial{AdsSocial: adsSocial}).
		Post(fmt.Sprintf("%s%s", s.API, common.SHEETY_ADSSOCIAL))
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

func (s *Sheety) FindAllReferalCode() (refcode []string, err error) {
	client := resty.New()
	resp, err := client.R().EnableTrace().
		Get(fmt.Sprintf("%s%s?filter[active]=%d", s.API, common.SHEETY_ADSPARTNER, 1))
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

func (s *Sheety) FindAllAdsSocial() (adsSocial []AdsSocial, err error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf("%s%s", s.API, common.SHEETY_ADSSOCIAL))
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
