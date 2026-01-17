package assets

import (
	"math/rand"
	"time"
)

func GetThumbnail() string {
	photos := []string{
		"https://image-assets.soutushenqi.com/UserUploadWallpaper_upload/1754411325048.png",
		"https://image-assets.soutushenqi.com/UserUploadWallpaper_upload/1764199527922.jpg",
		"https://image-assets.soutushenqi.com/UserUploadWallpaper_upload/1746637396204.jpg",
		"http://gd-hbimg.huaban.com/70e0c3d1f55bcf9b0175b92048dae628f58948e6258832-iB6FHT",
		"http://gd-hbimg.huaban.com/42516db95a2b0a5edd97b45af3bb9e89b2373c996be187-UUWyEL",
		"http://i0.hdslb.com/bfs/archive/7ce9f9d1678f0b07c0e607f850877a3e64bac8e9.jpg",
		"http://i0.hdslb.com/bfs/archive/adbed2750f5c0439554314335eca600f4b3cecf2.jpg",
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := rand.Intn(len(photos))
	return photos[seed]
}
