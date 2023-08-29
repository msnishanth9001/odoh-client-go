package commands

import (
	"fmt"
	"io/ioutil"
	"net/http"
	// "encoding/binary"
	// "reflect"
	// "encoding/pem"
	"runtime"
	
	odoh "github.com/cloudflare/odoh-go"
	"github.com/urfave/cli"
)

func fetchTargetConfigsFromWellKnown(url string) (odoh.ObliviousDoHConfigs, error) {
	_, file, no, ok := runtime.Caller(1)
    if ok {
        fmt.Printf("\n>>>>> == called from %s#%d", file, no)
    }

	fmt.Println(">>>>>  url fetchTargetConfigs", url, "\n")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return odoh.ObliviousDoHConfigs{}, err
	}

    client := http.Client{}
	resp, err := client.Do(req)
	
	if err != nil {
		return odoh.ObliviousDoHConfigs{}, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return odoh.ObliviousDoHConfigs{}, err
	}
	fmt.Println(">>>>> response", resp)
	// fmt.Println(reflect.TypeOf(bodyBytes))
	fmt.Println("\n>>>>> == bodyBytes", bodyBytes)

	// kemId := binary.BigEndian.Uint16(bodyBytes[0:])
	// kdfId := binary.BigEndian.Uint16(bodyBytes[2:])
	// aeadId := binary.BigEndian.Uint16(bodyBytes[4:])
	// publicKeyLength := binary.BigEndian.Uint16(bodyBytes[6:])
	// publicKey := bodyBytes[14:14+publicKeyLength]
	// fmt.Println(reflect.TypeOf(publicKey))

	// fmt.Println(">>>>> == here kemId, kdfId, aeadId, publicKeyLength, publicKey", kemId, kdfId, aeadId, publicKeyLength, publicKey)
	// fmt.Println(">>>>> == pubkey",publicKey)
	return odoh.UnmarshalObliviousDoHConfigs(bodyBytes)
}

func fetchTargetConfigs(targetName string) (odoh.ObliviousDoHConfigs, error) {
	_, file, no, ok := runtime.Caller(1)
    if ok {
        fmt.Printf("called from %s#%d\n", file, no)
    }    
	return fetchTargetConfigsFromWellKnown(buildOdohConfigURL(targetName).String())
}

func getTargetConfigs(c *cli.Context) error {
	fmt.Println(">>>>> == getTargetConfigs")
	targetName := c.String("target")
	pretty := c.Bool("pretty")

	odohConfigs, err := fetchTargetConfigs(targetName)
	if err != nil {
		return err
	}
	fmt.Println("odohConfigs", odohConfigs)
	if pretty {
		fmt.Println("ObliviousDoHConfigs:")
		for i, config := range odohConfigs.Configs {
			configContents := config.Contents
			fmt.Printf("  Config %d: Version(0x%04x), KEM(0x%04x), KDF(0x%04x), AEAD(0x%04x) KeyID(%x)\n", (i + 1), config.Version, configContents.KemID, configContents.KdfID, configContents.AeadID, configContents.KeyID())
		}
	} else {
		fmt.Printf("%x", odohConfigs.Marshal())
	}
	return nil
}
