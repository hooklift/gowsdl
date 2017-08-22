package main

import (
	"fmt"
	"github.com/oshapeman/gowsdl/cmd/gowsdl/tianqi"
	"github.com/oshapeman/gowsdl/cmd/gowsdl/dream"
)

func main()  {
	auth2 := dream.BasicAuth{Login:"admin",Password:"1234"}
	client2 := dream.NewDreamSoap("",true, &auth2)
	querycode2 := dream.SearchDreamInfo{Dream:"财富"}
	result2,_ := client2.SearchDreamInfo(&querycode2)
	fmt.Println(result2)


	auth3 := tianqi.BasicAuth{Login:"admin",Password:"1234"}
	client3 := tianqi.NewWeatherWebServiceSoap("",true, &auth3)
	querycode3 := tianqi.GetWeatherbyCityName{TheCityName:"上海"}
	result3,_ := client3.GetWeatherbyCityName(&querycode3)
	fmt.Println(result3)
}
