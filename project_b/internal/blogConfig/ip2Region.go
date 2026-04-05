package blogConfig

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// Ip2Region ip to region v4 and v6
type Ip2Region struct {
	v4 *xdb.Searcher
	v6 *xdb.Searcher
}

func loadSearcher(versionName, dbPath string) (*xdb.Searcher, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, err
	}
	version, err := xdb.VersionFromName(versionName)
	if err != nil {
		return nil, err
	}
	searcher, err := xdb.NewWithFileOnly(version, dbPath)
	if err != nil {
		return nil, err
	}

	return searcher, nil
}

func NewIp2Region() *Ip2Region {
	v4, err := loadSearcher("V4", "data/ip2region_v4.xdb")
	if err != nil {
		log.Println("ip2region v4 error: ", err)
		return nil
	}
	v6, err := loadSearcher("V6", "data/ip2region_v6.xdb")
	if err != nil {
		log.Println("ip2region v6 error: ", err)
		return nil
	}

	return &Ip2Region{
		v4: v4,
		v6: v6,
	}
}

func (i *Ip2Region) Close() {
	if i.v4 != nil {
		i.v4.Close()
	}
	if i.v6 != nil {
		i.v6.Close()
	}
}

func (i *Ip2Region) SearchByStr(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("invalid ip")
	}

	if ip.To4() != nil {
		if i.v4 != nil {
			return i.v4.SearchByStr(ipStr)
		}
	} else {
		if i.v6 != nil {
			return i.v6.SearchByStr(ipStr)
		}
	}
	return "", errors.New("ip2region not initialized")
}
