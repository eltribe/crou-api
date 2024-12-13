package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var fronts = []string{
	"예쁜",
	"귀여운",
	"배고픈",
	"철학적인",
	"현학적인",
	"슬픈",
	"푸른",
	"행복한",
	"밝은",
	"분주한",
	"빛나는",
	"건강한",
	"한가한",
	"달콤한",
	"착한",
	"관대한",
	"시원한",
	"낙천적인",
	"차가운",
	"멋쩍은",
	"흐믓한",
	"어린",
	"영리한",
	"화려한",
	"똑똑한",
	"든든한",
	"고요한",
}

var animals = []string{
	"호랑이",
	"비버",
	"부엉이",
	"여우",
	"치타",
	"문어",
	"고양이",
	"미어캣",
	"다람쥐",
	"까치",
	"두루미",
	"상어",
	"독수리",
	"치와와",
	"리트리버",
	"푸들",
	"비글",
	"원숭이",
	"라마",
	"바위",
	"나무",
	"야자수",
	"민들레",
	"꿀벌",
}

func CreateRandomNickName() string {
	p := fronts[rand.Intn(len(fronts))]
	a := animals[rand.Intn(len(animals))]
	return p + " " + a
}
