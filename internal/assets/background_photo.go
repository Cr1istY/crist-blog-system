package assets

import (
	"math/rand"
	"time"
)

func GetThumbnail() string {
	photos := []string{
		"https://image-assets.soutushenqi.com/UserUploadWallpaper_upload/1754411325048.png",
		"https://image-assets.soutushenqi.com/UserUploadWallpaper_upload/1764199527922.jpg",
		"http://gd-hbimg.huaban.com/3d279e012e9c4e7e62aa6d66c9de962b6038df6a30b7a1-SCpDAX",
		"https://image-assets.soutushenqi.com/UserUploadWallpaper_upload/1746637396204.jpg",
		"http://gd-hbimg.huaban.com/70e0c3d1f55bcf9b0175b92048dae628f58948e6258832-iB6FHT",
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := rand.Intn(len(photos))
	return photos[seed]
}
