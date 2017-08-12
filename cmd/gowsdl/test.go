package main


import (
	"fmt"
	"github.com/oshapeman/gowsdl/cmd/gowsdl/myservice"
)

func main()  {
	auth3 := myservice.BasicAuth{Login:"admin",Password:"1234"}
	client3 := myservice.NewWeatherWebServiceSoap("",true, &auth3)
	querycode3 := myservice.GetWeatherbyCityName{TheCityName:"北京"}
	result3,_ := client3.GetWeatherbyCityName(&querycode3)
	fmt.Println(result3)
}
