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
		"http://i0.hdslb.com/bfs/archive/7ce9f9d1678f0b07c0e607f850877a3e64bac8e9.jpg",
		"http://i0.hdslb.com/bfs/archive/adbed2750f5c0439554314335eca600f4b3cecf2.jpg",
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := rand.Intn(len(photos))
	return photos[seed]
}
