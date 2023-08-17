Пример использования:

```
package main

import (
	"context"
	"fmt"
	"log"

	twitch "github.com/atasik/twitch-sdk"
)

func main() {
	ctx := context.Background()

	//Создаю новый клиент
	client, err := twitch.NewClient(<client-id>, <client-secret>)
	if err != nil {
		log.Fatalf("Failed to create client: %s", err.Error())
	}

	//Получаю url для авторизации
	url, err := client.GetAuthorizationURL(<redirect-url>)
	if err != nil {
		log.Fatalf("Failed to get authorization  url: %s", err.Error())
	}
	fmt.Println(url)

	//Получаю access token
	authResp, err := client.GetAccessToken(ctx)

	if err != nil {
		log.Fatalf("failed to authorize: %s", err)
	}

	fmt.Println(authResp.AccessToken)
	
	userResp, err := client.GetUser(ctx, twitch.UserRequest{
		UserLogin: "xqc",
	}, *authResp)

	if err != nil {
		log.Fatalf("failed to get user: %s", err)
	}

	subStatus, err := client.GetSubscriptions(ctx, *authResp)
	if err != nil {
		log.Fatalf("failed to get status: %s", err)
	}

	fmt.Println(subStatus)

	subResp, err := client.Subscribe(ctx, twitch.SubRequest{
		Type: "channel.follow",
		Condition: twitch.Condition{
			Id: userResp.Data[0].Id,
		},
		Transport: twitch.Transport{
			Callback: <callback-url>,
			Secret:   <secret-code>,
		}}, *authResp)

	if err != nil {
		log.Fatalf("failed to subscribe: %s", err)
	}

	fmt.Println(subResp)
}

```
