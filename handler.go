package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"
	name := r.URL.Query().Get("name")
	if name != "" {
		message = fmt.Sprintf("Hello, %s. This HTTP triggered function executed successfully.\n", name)
	}
	fmt.Fprint(w, message)
}

func getKeys(w http.ResponseWriter, r *http.Request) {
	keyVaultName := "testbuzzkvmigrate"
	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)
	// create credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	// create azkeys client
	client, err := azkeys.NewClient(keyVaultUrl, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	// create RSA Key
	// rsaKeyParams := azkeys.CreateKeyParameters{
	// 	Kty:     to.Ptr(azkeys.JSONWebKeyTypeRSA),
	// 	KeySize: to.Ptr(int32(2048)),
	// }
	// rsaResp, err := client.CreateKey(context.TODO(), "new-rsa-key", rsaKeyParams, nil)
	// if err != nil {
	// 	log.Fatalf("failed to create rsa key: %v", err)
	// }
	// fmt.Printf("New RSA key ID: %s\n", *rsaResp.Key.KID)

	// // create EC Key
	// ecKeyParams := azkeys.CreateKeyParameters{
	// 	Kty:   to.Ptr(azkeys.JSONWebKeyTypeEC),
	// 	Curve: to.Ptr(azkeys.JSONWebKeyCurveNameP256),
	// }
	// ecResp, err := client.CreateKey(context.TODO(), "new-ec-key", ecKeyParams, nil)
	// if err != nil {
	// 	log.Fatalf("failed to create ec key: %v", err)
	// }
	// fmt.Printf("New EC key ID: %s\n", *ecResp.Key.KID)

	// list all vault keys
	fmt.Println("List all vault keys:")
	pager := client.NewListKeysPager(nil)
	var response []string
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		for _, key := range page.Value {
			fmt.Println(*key.KID)
			response = append(response, key.KID.Name())
		}
	}
	fmt.Fprint(w, response)
	// update key properties to disable key
	// updateParams := azkeys.UpdateKeyParameters{
	// 	KeyAttributes: &azkeys.KeyAttributes{
	// 		Enabled: to.Ptr(false),
	// 	},
	// }
	// // an empty string version updates the latest version of the key
	// version := ""
	// updateResp, err := client.UpdateKey(context.TODO(), "new-rsa-key", version, updateParams, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Key %s Enabled attribute set to: %t\n", *updateResp.Key.KID, *updateResp.Attributes.Enabled)

	// delete the created keys
	// for _, keyName := range []string{"new-rsa-key", "new-ec-key"} {
	// 	// DeleteKey returns when Key Vault has begun deleting the key. That can take several
	// 	// seconds to complete, so it may be necessary to wait before performing other operations
	// 	// on the deleted key.
	// 	delResp, err := client.DeleteKey(context.TODO(), keyName, nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("Successfully deleted key %s", *delResp.Key.KID)
	// }
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/api/HttpExample", helloHandler)
	http.HandleFunc("/api/GetKeys", getKeys)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
