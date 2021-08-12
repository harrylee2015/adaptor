package main

import (
	"adaptor/client"
	"adaptor/types"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

func main() {
	conf := types.ReadConf("conf.yaml")
	client := client.NewClient(conf.JsonRpc)
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/getBlockHeight", client.GetBlockHeight),
		rest.Get("/getNodeCount", client.GetNodeCount),
		rest.Get("/getTxCountAccepted", client.GetTxAccepted),
		rest.Get("/getTxCountConfirmed", client.GetTxConfirmed),
		rest.Post("/getTxInfo", client.GetTxInfo),
		rest.Get("/getBlockInfo", client.GetBlockInfo),
		rest.Post("/createTx", client.CreateTx),
		rest.Post("/sendTx", client.SendTx),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8999", api.MakeHandler()))

}
