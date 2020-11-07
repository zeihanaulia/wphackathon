package engine

import (
	"fmt"
	"log"

	"github.com/zeihanaulia/instagram-scraper/config"
)

func Start(cfg config.Configurations) {
	sheety := NewSheety(cfg.SHEETY_API)
	instagram := NewInstagram(cfg.INSTAGRAM_USERNAME, cfg.INSTAGRAM_PASSWORD)

	adsPartners, err := sheety.FindAllAdsPartner()
	if err != nil {
		log.Fatal(fmt.Sprintf("error when FindAllAdsPartner: %v", err))
		return
	}

	fullAdsSocial, err := sheety.FindAllAdsSocial()
	if err != nil {
		log.Fatal(fmt.Sprintf("error when FindAllAdsSocial: %v", err))
		return
	}

	hashtags, err := instagram.FindHastagFor(adsPartners)
	if err != nil {
		log.Fatal(fmt.Sprintf("error when FindHastagFor: %v", err))
		return
	}

	inserts, updates, err := inserOrUpdate(fullAdsSocial, hashtags)
	if err != nil {
		log.Fatal(fmt.Sprintf("error when inserOrUpdate: %v", err))
		return
	}

	for _, u := range updates {
		_, err := sheety.UpdateAdsSocial(u)
		if err != nil {
			log.Fatal(fmt.Sprintf("error when UpdateAdsSocial: %v", err))
			return
		}
	}

	for _, i := range inserts {
		_, err := sheety.InsertAdsSocial(i)
		if err != nil {
			log.Fatal(fmt.Sprintf("error when InsertAdsSocial: %v", err))
			return
		}
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
