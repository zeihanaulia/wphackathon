package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/ahmdrz/goinsta/v2"
)

type Instagram struct {
	Username string
	Password string
}

func NewInstagram(username, password string) Instagram {
	return Instagram{Username: username, Password: password}
}

func (i *Instagram) FindHastagFor(ads []AdsPartner) (adsSocials []AdsSocial, err error) {
	insta := goinsta.New(i.Username, i.Password)
	err = insta.Login()
	if err != nil {
		log.Println(fmt.Sprintf("cannot connect instagram, err: %v", err))
		return
	}

	for _, rc := range ads {
		log.Println("start sync hashtag: ", rc.RefCode)
		hashtag := insta.NewHashtag(rc.RefCode)
		_ = hashtag.Sync()

		pages := 100
		for page := 0; page < pages && hashtag.Next(); page++ {
			for _, section := range hashtag.Sections {
				for _, media := range section.LayoutContent.Medias {
					takenAt := time.Unix(media.Item.TakenAt, 0).Format("2006-01-02")
					takenAtd, _ := time.Parse("2006-01-02", takenAt)
					startDate, _ := time.Parse("2006-01-02", rc.StartDate)
					endDate, _ := time.Parse("2006-01-02", rc.EndDate)

					if takenAtd.Equal(startDate) || takenAtd.After(startDate) {
						if takenAtd.Equal(endDate) || takenAtd.Before(endDate) {
							adsSocials = append(adsSocials, AdsSocial{
								Code:           media.Item.Code,
								RefCode:        rc.RefCode,
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
		}
	}

	defer func() {
		_ = insta.Logout()
	}()

	return
}
