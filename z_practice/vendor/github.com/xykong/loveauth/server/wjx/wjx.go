package wjx

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"github.com/xykong/loveauth/settings"
	"strconv"
)

var postHandlers = make(map[string]gin.HandlerFunc)
var games = make(map[uint32]map[uint8][]string)

func Start(group *gin.RouterGroup) {

	for key, value := range postHandlers {
		group.POST(key, value)
	}

	g := settings.GetStringMap("wjx", "callback")

	for a, v := range g {
		value, err := strconv.Atoi(a)
		if err != nil {
			continue
		}

		AreaId := uint32(value)

		for p, s := range v.(map[string]interface{}) {

			value, err := strconv.Atoi(p)
			if err != nil {
				continue
			}

			PlatId := uint8(value)

			if _, ok := games[AreaId]; !ok {
				games[AreaId] = make(map[uint8][]string)
			}

			if s != nil {

				for _, u := range s.([]interface{}) {

					games[AreaId][PlatId] = append(games[AreaId][PlatId], u.(string))
				}
			}
			//games[AreaId][PlatId] = s.(string)
		}
	}
}


func randGames(areaId uint32, platId uint8) string {

	count := len(games[areaId][platId])
	if count <= 0 {

		return ""
	}

	return games[areaId][platId][rand.Intn(count)]
}
